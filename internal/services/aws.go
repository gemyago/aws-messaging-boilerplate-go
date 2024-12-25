package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"runtime"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/gemyago/aws-sqs-boilerplate-go/internal/diag"
	"go.uber.org/dig"
	"golang.org/x/sync/errgroup"
)

//go:generate mockery --name=MessageSender --filename=mock_message_sender.go --config ../../.mockery-funcs.yaml

type AWSConfigDeps struct {
	dig.In

	Region       string `name:"config.aws.region"`
	BaseEndpoint string `name:"config.aws.baseEndpoint" optional:"true"`
}

func newAWSConfigFactory(ctx context.Context) func(deps AWSConfigDeps) (aws.Config, error) {
	return func(deps AWSConfigDeps) (aws.Config, error) {
		opts := []func(*config.LoadOptions) error{
			config.WithRegion(deps.Region),
		}
		// BaseEndpoint is defined for local/test modes only and points on localstack instance
		if deps.BaseEndpoint != "" {
			opts = append(opts,
				config.WithBaseEndpoint(deps.BaseEndpoint),
				config.WithCredentialsProvider(aws.AnonymousCredentials{}),
			)
		}
		cfg, err := config.LoadDefaultConfig(ctx, opts...)
		if err != nil {
			return aws.Config{}, fmt.Errorf("failed to load aws configuration, %w", err)
		}
		return cfg, nil
	}
}

type Message struct {
	Id       string `json:"id"` //nolint:revive,stylecheck // Id is used to match apigen generated code
	Name     string `json:"name"`
	Comments string `json:"comments,omitempty"`
}

type MessageSender func(ctx context.Context, message *Message) error

type MessageSenderDeps struct {
	dig.In

	RootLogger       *slog.Logger
	SnsClient        *sns.Client
	MessagesTopicARN string `name:"config.aws.sqs.messagesTopicARN"`
}

func NewMessageSender(deps MessageSenderDeps) MessageSender {
	logger := deps.RootLogger.WithGroup("services.message-sender")
	return func(ctx context.Context, message *Message) error {
		body, err := json.Marshal(message)
		if err != nil {
			return fmt.Errorf("failed to marshal message, %w", err)
		}
		res, err := deps.SnsClient.Publish(ctx, &sns.PublishInput{
			Message:  aws.String(string(body)),
			TopicArn: aws.String(deps.MessagesTopicARN),
		})
		if err != nil {
			return fmt.Errorf("failed send message to sqs queue, %w", err)
		}
		logger.InfoContext(ctx, "Message sent", slog.String("messageId", *res.MessageId))
		return nil
	}
}

type RawMessageHandler func(ctx context.Context, rawMessage types.Message) error

func NewRawMessageHandler[TMessage any](
	handler func(ctx context.Context, msg *TMessage) error,
) RawMessageHandler {
	return func(ctx context.Context, rawMessage types.Message) error {
		var message TMessage
		if err := json.Unmarshal([]byte(*rawMessage.Body), &message); err != nil {
			return fmt.Errorf("failed to unmarshal message, %w", err)
		}
		return handler(ctx, &message)
	}
}

type MessagesPollerQueue struct {
	QueueURL string
	Handler  RawMessageHandler

	// The duration (in seconds) that the received messages are hidden
	// from subsequent retrieve requests after being retrieved
	// by a ReceiveMessage request.
	VisibilityTimeout int32
}

type MessagesPollerDeps struct {
	dig.In

	// config
	MaxPollWaitTimeSec int32 `name:"config.aws.sqs.maxPollWaitTimeSec"`

	RootLogger *slog.Logger
	SqsClient  *sqs.Client
}

type MessagesPoller struct {
	queues               []MessagesPollerQueue
	deps                 MessagesPollerDeps
	maxProcessingWorkers int32
	logger               *slog.Logger
}

func NewMessagesPoller(deps MessagesPollerDeps) *MessagesPoller {
	return &MessagesPoller{
		deps:                 deps,
		maxProcessingWorkers: int32(math.Min(math.MaxInt32, float64(runtime.NumCPU()))),
		logger:               deps.RootLogger.WithGroup("services.messages-poller"),
	}
}

func (p *MessagesPoller) RegisterQueue(queue MessagesPollerQueue) {
	p.queues = append(p.queues, queue)
}

func (p *MessagesPoller) wrapRawMessageHandlerWithDeleteOnSuccess(
	queueURL string,
	handler RawMessageHandler,
) RawMessageHandler {
	return func(ctx context.Context, rawMessage types.Message) error {
		err := handler(ctx, rawMessage)
		if err != nil {
			return fmt.Errorf("failed to handle message with target handler, %w", err)
		}
		if _, err = p.deps.SqsClient.DeleteMessage(ctx, &sqs.DeleteMessageInput{
			QueueUrl:      aws.String(queueURL),
			ReceiptHandle: rawMessage.ReceiptHandle,
		}); err != nil {
			return fmt.Errorf("failed to acknowledge message, %w", err)
		}
		return nil
	}
}

type processingData struct {
	rawMessage types.Message
	handler    RawMessageHandler
}

func (p *MessagesPoller) startProcessingWorkers(
	ctx context.Context,
) []chan processingData {
	processingWorkerChannels := make([]chan processingData, len(p.queues))
	for i := range processingWorkerChannels {
		processingWorkerChannels[i] = make(chan processingData)
		go func(ch <-chan processingData) {
			for data := range ch {
				if err := data.handler(ctx, data.rawMessage); err != nil {
					p.logger.ErrorContext(ctx, "Failed to process message", diag.ErrAttr(err))
				}
			}
		}(processingWorkerChannels[i])
	}
	return processingWorkerChannels
}

// Returns true if polling should be continued
// otherwise returns false and error if polling should be stopped.
func (p *MessagesPoller) handleReceiveError(
	queueURL string,
	err error,
) (bool, error) {
	if err == nil {
		return true, nil
	}
	if errors.Is(err, context.Canceled) {
		return false, nil
	}
	// TODO: Retry logic can be added here
	return false, fmt.Errorf("failed to receive messages from queue %s: %w", queueURL, err)
}

func (p *MessagesPoller) Start(ctx context.Context) error {
	p.logger.InfoContext(ctx, "Starting messages poller",
		slog.Int("queues", len(p.queues)),
		slog.Int64("maxProcessingWorkers", int64(p.maxProcessingWorkers)),
	)
	processingWorkerChannels := p.startProcessingWorkers(ctx)

	grp := errgroup.Group{}

	for _, queue := range p.queues {
		handler := p.wrapRawMessageHandlerWithDeleteOnSuccess(
			queue.QueueURL,
			queue.Handler,
		)
		grp.Go(func() error {
			for {
				gotMessages, err := p.deps.SqsClient.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
					QueueUrl:            &queue.QueueURL,
					MaxNumberOfMessages: p.maxProcessingWorkers,
					WaitTimeSeconds:     p.deps.MaxPollWaitTimeSec,
					VisibilityTimeout:   queue.VisibilityTimeout,
				})
				var shouldContinue bool
				if shouldContinue, err = p.handleReceiveError(queue.QueueURL, err); !shouldContinue {
					return err
				}
				for _, rawMessage := range gotMessages.Messages {
					for _, ch := range processingWorkerChannels {
						ch <- processingData{
							rawMessage: rawMessage,
							handler:    handler,
						}
					}
				}
			}
		})
	}

	return grp.Wait()
}
