package awsapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"runtime"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/gemyago/aws-messaging-boilerplate-go/internal/diag"
	"github.com/gofrs/uuid/v5"
	"go.uber.org/dig"
	"golang.org/x/sync/errgroup"
)

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
			return fmt.Errorf("failed to handle message with target handler: %w", err)
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
				handlerCtx := diag.SetLogAttributesToContext(ctx, diag.LogAttributes{
					CorrelationID: slog.StringValue(uuid.Must(uuid.NewV4()).String()),
				})
				messageID := *data.rawMessage.MessageId
				if err := data.handler(handlerCtx, data.rawMessage); err != nil {
					p.logger.ErrorContext(ctx,
						"Failed to process message",
						diag.ErrAttr(err),
						slog.String("messageID", messageID),
						slog.Any("messageAttributes", data.rawMessage.MessageAttributes),
						slog.Any("systemAttributes", data.rawMessage.Attributes),
					)
				} else {
					p.logger.InfoContext(ctx,
						"Message processed",
						slog.String("messageID", messageID),
						slog.Any("messageAttributes", data.rawMessage.MessageAttributes),
						slog.Any("systemAttributes", data.rawMessage.Attributes),
					)
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
			p.logger.DebugContext(ctx,
				"Polling messages from queue",
				slog.String("queueURL", queue.QueueURL),
				slog.String("visibilityTimeout", fmt.Sprintf("%ds", queue.VisibilityTimeout)),
				slog.Int64("pollWaitTimeSec", int64(p.deps.MaxPollWaitTimeSec)),
			)
			for {
				gotMessages, err := p.deps.SqsClient.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
					QueueUrl:                    &queue.QueueURL,
					MaxNumberOfMessages:         p.maxProcessingWorkers,
					WaitTimeSeconds:             p.deps.MaxPollWaitTimeSec,
					VisibilityTimeout:           queue.VisibilityTimeout,
					MessageAttributeNames:       []string{"All"},
					MessageSystemAttributeNames: []types.MessageSystemAttributeName{"All"},
				})
				var shouldContinue bool
				if shouldContinue, err = p.handleReceiveError(queue.QueueURL, err); !shouldContinue {
					return err
				}
				p.logger.DebugContext(ctx,
					"Processing queue messages",
					slog.String("queueURL", queue.QueueURL),
					slog.Int("messagesCount", len(gotMessages.Messages)),
					slog.Any("resultMetadata", gotMessages.ResultMetadata),
				)
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
