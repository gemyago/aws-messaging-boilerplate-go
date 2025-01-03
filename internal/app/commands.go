package app

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/gemyago/aws-sqs-boilerplate-go/internal/api/http/v1routes/handlers"
	"github.com/gemyago/aws-sqs-boilerplate-go/internal/api/http/v1routes/models"
	"go.uber.org/dig"
)

type messageSender[TMessage any] func(
	ctx context.Context,
	message *TMessage,
) error

type CommandsDeps struct {
	dig.In

	RootLogger *slog.Logger

	SendDummySNSMessage messageSender[models.DummyMessage]
}

type Commands struct {
	logger *slog.Logger
	deps   CommandsDeps
}

func (c *Commands) PublishMessage(ctx context.Context, req *handlers.MessagesPublishDummyMessageRequest) error {
	if err := c.deps.SendDummySNSMessage(ctx, req.Payload); err != nil {
		return fmt.Errorf("failed to send message, %w", err)
	}
	return nil
}

func (c *Commands) ProcessMessage(ctx context.Context, msg *models.DummyMessage) error {
	if c.logger.Enabled(ctx, slog.LevelDebug) {
		c.logger.DebugContext(ctx, "Processing message", slog.Any("message", msg))
	}

	// TODO: Some processing logic. Probably just simulate CPU-bound work.

	return nil
}

func NewCommands(deps CommandsDeps) *Commands {
	return &Commands{
		logger: deps.RootLogger.WithGroup("app.commands"),
		deps:   deps,
	}
}
