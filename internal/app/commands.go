package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/gemyago/aws-sqs-boilerplate-go/internal/api/http/v1routes/handlers"
	"github.com/gemyago/aws-sqs-boilerplate-go/internal/api/http/v1routes/models"
	"github.com/samber/lo"
	"go.uber.org/dig"
)

type messageSender[TMessage any] func(
	ctx context.Context,
	message *TMessage,
) error

type CommandsDeps struct {
	dig.In

	RootLogger *slog.Logger

	SendDummySNSMessage         messageSender[models.DummyMessage] `name:"dummy-sns-message-sender"`
	SendDummyEventBridgeMessage messageSender[models.DummyMessage] `name:"dummy-eventbridge-message-sender"`
}

type Commands struct {
	logger *slog.Logger
	deps   CommandsDeps
}

func (c *Commands) PublishMessage(ctx context.Context, req *handlers.MessagesPublishDummyMessageRequest) error {
	sendDummyMessage := lo.If(
		req.Target == handlers.MessagesPublishDummyMessageTargetSNS,
		c.deps.SendDummySNSMessage,
	).Else(c.deps.SendDummyEventBridgeMessage)
	if err := sendDummyMessage(ctx, req.Payload); err != nil {
		return fmt.Errorf("failed to send message, %w", err)
	}
	return nil
}

func (c *Commands) ProcessMessage(ctx context.Context, msg *models.DummyMessage) error {
	c.logger.DebugContext(ctx, "Processing message", slog.Any("message", msg))

	if msg.FailProcessing {
		return errors.New("simulated processing error")
	}

	return nil
}

func NewCommands(deps CommandsDeps) *Commands {
	return &Commands{
		logger: deps.RootLogger.WithGroup("app.commands"),
		deps:   deps,
	}
}
