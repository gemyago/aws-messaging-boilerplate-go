package v1controllers

import (
	"context"

	"github.com/gemyago/aws-sqs-boilerplate-go/internal/api/http/v1routes/handlers"
	"go.uber.org/dig"
)

func healthCheck(_ context.Context) error {
	// noop
	return nil
}

type messagesCommands interface {
	PublishMessage(context.Context, *handlers.MessagesPublishMessageRequest) error
}

type MessagesControllerDeps struct {
	dig.In

	Commands messagesCommands
}

func NewMessagesController(deps MessagesControllerDeps) *handlers.MessagesController {
	return handlers.BuildMessagesController().
		HandleHealthCheck.With(healthCheck).
		HandlePublishMessage.With(deps.Commands.PublishMessage).
		Finalize()
}
