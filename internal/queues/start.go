package queues

import (
	"context"

	"github.com/gemyago/aws-sqs-boilerplate-go/internal/app"
	"github.com/gemyago/aws-sqs-boilerplate-go/internal/services"
	"go.uber.org/dig"
)

type Deps struct {
	dig.In

	// config
	MessagesQueueURL                  string `name:"config.aws.sqs.messagesQueueURL"`
	MessagesQueueVisibilityTimeoutSec int32  `name:"config.aws.sqs.messagesQueueVisibilityTimeoutSec"`

	// app layer
	Commands *app.Commands

	// services
	*services.MessagesPoller
	*services.ShutdownHooks
}

func StartPolling(ctx context.Context, deps Deps) error {
	queuesCtx, cancel := context.WithCancel(ctx)
	deps.ShutdownHooks.RegisterNoCtx("queues", func() error {
		cancel()
		return nil
	})

	poller := deps.MessagesPoller
	poller.RegisterQueue(services.MessagesPollerQueue{
		QueueURL:          deps.MessagesQueueURL,
		Handler:           services.NewRawMessageHandler(deps.Commands.ProcessMessage),
		VisibilityTimeout: deps.MessagesQueueVisibilityTimeoutSec,
	})

	return poller.Start(queuesCtx)
}
