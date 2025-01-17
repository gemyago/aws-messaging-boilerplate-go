package v1controllers

import (
	"context"

	"github.com/gemyago/aws-messaging-boilerplate-go/internal/api/http/v1routes/handlers"
	"github.com/gemyago/aws-messaging-boilerplate-go/internal/api/http/v1routes/models"
	"go.uber.org/dig"
)

func healthCheck(_ context.Context) error {
	// noop
	return nil
}

type dummyMessagesCommands interface {
	PublishMessage(context.Context, *handlers.MessagesPublishDummyMessageRequest) error
	ProcessMessage(context.Context, *models.DummyMessage) error
}

type MessagesControllerDeps struct {
	dig.In

	Commands dummyMessagesCommands
}

func NewMessagesController(deps MessagesControllerDeps) *handlers.MessagesController {
	return handlers.BuildMessagesController().
		HandleHealthCheck.With(healthCheck).
		HandlePublishDummyMessage.With(deps.Commands.PublishMessage).
		HandleProcessDummyMessage.With(
		func(ctx context.Context, mpmr *handlers.MessagesProcessDummyMessageRequest) error {
			return deps.Commands.ProcessMessage(ctx, mpmr.Payload)
		}).
		Finalize()
}
