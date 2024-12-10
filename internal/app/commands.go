package app

import (
	"context"
	"fmt"

	"github.com/gemyago/aws-sqs-boilerplate-go/internal/api/http/v1routes/handlers"
	"github.com/gemyago/aws-sqs-boilerplate-go/internal/services"
	"go.uber.org/dig"
)

type CommandsDeps struct {
	dig.In

	SendMessage services.MessageSender
}

type Commands struct {
	deps CommandsDeps
}

func (c *Commands) PublishMessage(ctx context.Context, req *handlers.MessagesPublishMessageRequest) error {
	msg := services.Message(*req.Payload)
	if err := c.deps.SendMessage(ctx, &msg); err != nil {
		return fmt.Errorf("failed to send message, %w", err)
	}
	return nil
}

func NewCommands(deps CommandsDeps) *Commands {
	return &Commands{deps: deps}
}
