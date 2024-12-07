package app

import (
	"context"

	"github.com/gemyago/aws-sqs-boilerplate-go/internal/api/http/v1routes/handlers"
)

type Commands struct{}

func (c *Commands) PublishMessage(_ context.Context, _ *handlers.MessagesPublishMessageRequest) error {
	// noop
	return nil
}

func NewCommands() *Commands {
	return &Commands{}
}
