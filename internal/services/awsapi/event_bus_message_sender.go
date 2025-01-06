package awsapi

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
	"go.uber.org/dig"
)

type EventBusMessageSenderDeps struct {
	dig.In

	EventBusName   string `name:"config.aws.eventBus.name"`
	EventBusSource string `name:"config.aws.eventBus.source"`

	RootLogger *slog.Logger
	Client     *eventbridge.Client
}

func NewEventBusMessageSender[TMessage any](detailType string, deps EventBusMessageSenderDeps) MessageSender[TMessage] {
	logger := deps.RootLogger.WithGroup("services.event-bus-message-sender")
	return func(ctx context.Context, message *TMessage) error {
		detail, err := json.Marshal(message)
		if err != nil {
			return fmt.Errorf("failed to marshal message, %w", err)
		}
		res, err := deps.Client.PutEvents(ctx, &eventbridge.PutEventsInput{
			Entries: []types.PutEventsRequestEntry{
				{
					EventBusName: &deps.EventBusName,
					Source:       &deps.EventBusSource,
					DetailType:   &detailType,
					Detail:       aws.String(string(detail)),
				},
			},
		})
		if err != nil {
			return fmt.Errorf("failed send message to event bus, %w", err)
		}
		logger.DebugContext(ctx, "PutEvents response", slog.Any("response", res))
		logger.InfoContext(ctx, "Message sent",
			slog.String("eventBusName", deps.EventBusName),
			slog.String("eventBusSource", deps.EventBusSource),
			slog.String("messageId", "n/a"), // TODO
		)
		return nil
	}
}
