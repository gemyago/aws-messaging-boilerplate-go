package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"runtime"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/gemyago/aws-sqs-boilerplate-go/internal/diag"
	"go.uber.org/dig"
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
	SqsClient        *sqs.Client
	MessagesQueueURL string `name:"config.aws.sqs.messagesQueueUrl"`
}

func NewMessageSender(deps MessageSenderDeps) MessageSender {
	logger := deps.RootLogger.WithGroup("services.message-sender")
	return func(ctx context.Context, message *Message) error {
		body, err := json.Marshal(message)
		if err != nil {
			return fmt.Errorf("failed to marshal message, %w", err)
		}
		res, err := deps.SqsClient.SendMessage(ctx, &sqs.SendMessageInput{
			MessageBody: aws.String(string(body)),
			QueueUrl:    aws.String(deps.MessagesQueueURL),
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

type pollerQueue struct {
	queueURL string
	handler  RawMessageHandler
}

type MessagesPollerDeps struct {
	dig.In

	// config
	MaxPollWaitTimeSec int32 `name:"config.aws.sqs.maxPollWaitTimeSec"`

	RootLogger *slog.Logger
	SqsClient  *sqs.Client
}

type MessagesPoller struct {
	queues               []pollerQueue
	deps                 MessagesPollerDeps
	maxProcessingWorkers int
	logger               *slog.Logger
}

func NewMessagesPoller(deps MessagesPollerDeps) *MessagesPoller {
	return &MessagesPoller{
		deps:                 deps,
		maxProcessingWorkers: runtime.NumCPU(),
		logger:               deps.RootLogger.WithGroup("services.messages-poller"),
	}
}

func (p *MessagesPoller) RegisterHandler(
	queueURL string,
	handler RawMessageHandler,
) {
	p.queues = append(p.queues, pollerQueue{
		queueURL: queueURL,
		handler:  handler,
	})
}

func (p *MessagesPoller) Start(ctx context.Context) error {
	type processingData struct {
		rawMessage types.Message
		handler    RawMessageHandler
	}
	p.logger.InfoContext(ctx, "Starting messages poller",
		slog.Int("queues", len(p.queues)),
		slog.Int("maxProcessingWorkers", p.maxProcessingWorkers),
	)
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

	for _, queue := range p.queues {
		handleRawMessagesWithDeleteSucceeded := func(ctx context.Context, rawMessage types.Message) error {
			err := queue.handler(ctx, rawMessage)
			if err != nil {
				return fmt.Errorf("failed to handle message with target handler, %w", err)
			}
			if _, err = p.deps.SqsClient.DeleteMessage(ctx, &sqs.DeleteMessageInput{
				QueueUrl:      aws.String(queue.queueURL),
				ReceiptHandle: rawMessage.ReceiptHandle,
			}); err != nil {
				return fmt.Errorf("failed to acknowledge message, %w", err)
			}
			return nil
		}
		go func() {
			for {
				gotMessages, err := p.deps.SqsClient.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
					QueueUrl:            aws.String(queue.queueURL),
					MaxNumberOfMessages: int32(p.maxProcessingWorkers),
					WaitTimeSeconds:     p.deps.MaxPollWaitTimeSec,
					VisibilityTimeout:   1, // configurable per queue
				})
				if err != nil {
					// TODO: Something like max retries or exponential backoff or some other strategy is required
					p.logger.ErrorContext(ctx, "Failed to receive messages. Will retry", diag.ErrAttr(err))
					continue
				}
				for _, rawMessage := range gotMessages.Messages {
					for _, ch := range processingWorkerChannels {
						ch <- processingData{
							rawMessage: rawMessage,
							handler:    handleRawMessagesWithDeleteSucceeded,
						}
					}
				}
			}
		}()
	}
	return nil
}
