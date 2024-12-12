package app

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/gemyago/aws-sqs-boilerplate-go/internal/api/http/v1routes/handlers"
	"github.com/gemyago/aws-sqs-boilerplate-go/internal/services"
	"go.uber.org/dig"
)

type CommandsDeps struct {
	dig.In

	RootLogger *slog.Logger

	SendMessage services.MessageSender
}

type Commands struct {
	logger *slog.Logger
	deps   CommandsDeps
}

func (c *Commands) PublishMessage(ctx context.Context, req *handlers.MessagesPublishMessageRequest) error {
	msg := services.Message(*req.Payload)
	if err := c.deps.SendMessage(ctx, &msg); err != nil {
		return fmt.Errorf("failed to send message, %w", err)
	}
	return nil
}

func (c *Commands) ProcessMessage(ctx context.Context, msg *services.Message) error {
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
