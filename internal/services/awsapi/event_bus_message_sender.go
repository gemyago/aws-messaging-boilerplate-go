package awsapi

import (
	"context"
	"log/slog"

	"go.uber.org/dig"
)

type EventBusMessageSenderDeps struct {
	dig.In

	EventBusName   string `name:"config.aws.eventBus.name"`
	EventBusSource string `name:"config.aws.eventBus.source"`

	RootLogger *slog.Logger
}

func NewEventBusMessageSender[TMessage any](detailType string, deps EventBusMessageSenderDeps) MessageSender[TMessage] {
	logger := deps.RootLogger.WithGroup("services.event-bus-message-sender")
	return func(ctx context.Context, message *TMessage) error {
		logger.InfoContext(ctx, "Message sent", slog.String("messageId", "n/a"))
		return nil
	}
}
