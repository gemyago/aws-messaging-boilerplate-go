package v1controllers

import (
	"context"

	"github.com/gemyago/aws-sqs-boilerplate-go/internal/api/http/v1routes/handlers"
	"github.com/gemyago/aws-sqs-boilerplate-go/internal/api/http/v1routes/models"
	"go.uber.org/dig"
)

func healthCheck(_ context.Context) error {
	// noop
	return nil
}

type messagesCommands interface {
	PublishMessage(context.Context, *handlers.MessagesPublishMessageRequest) error
	ProcessMessage(context.Context, *models.Message) error
}

type MessagesControllerDeps struct {
	dig.In

	Commands messagesCommands
}

func NewMessagesController(deps MessagesControllerDeps) *handlers.MessagesController {
	return handlers.BuildMessagesController().
		HandleHealthCheck.With(healthCheck).
		HandlePublishMessage.With(deps.Commands.PublishMessage).
		HandleProcessMessage.With(
		func(ctx context.Context, mpmr *handlers.MessagesProcessMessageRequest) error {
			return deps.Commands.ProcessMessage(ctx, mpmr.Payload)
		}).
		Finalize()
}
