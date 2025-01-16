package queues

import (
	"context"

	"github.com/gemyago/aws-messaging-boilerplate-go/internal/app"
	"github.com/gemyago/aws-messaging-boilerplate-go/internal/services"
	"github.com/gemyago/aws-messaging-boilerplate-go/internal/services/awsapi"
	"go.uber.org/dig"
)

type messagePoller interface {
	RegisterQueue(queue awsapi.MessagesPollerQueue)
	Start(ctx context.Context) error
}

type Deps struct {
	dig.In

	// config
	DummyMessagesQueueURL                  string `name:"config.aws.sqs.dummyMessagesQueueUrl"`
	DummyMessagesQueueVisibilityTimeoutSec int32  `name:"config.aws.sqs.dummyMessagesQueueVisibilityTimeoutSec"`

	// app layer
	Commands *app.Commands

	// services
	MessagesPoller messagePoller
	*services.ShutdownHooks
}

func StartPolling(ctx context.Context, deps Deps) error {
	queuesCtx, cancel := context.WithCancel(ctx)
	deps.ShutdownHooks.RegisterNoCtx("queues", func() error {
		cancel()
		return nil
	})

	poller := deps.MessagesPoller
	poller.RegisterQueue(awsapi.MessagesPollerQueue{
		QueueURL:          deps.DummyMessagesQueueURL,
		Handler:           awsapi.NewRawMessageHandler(deps.Commands.ProcessMessage),
		VisibilityTimeout: deps.DummyMessagesQueueVisibilityTimeoutSec,
	})

	return poller.Start(queuesCtx)
}
