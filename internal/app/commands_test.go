package app

import (
	"context"
	"errors"
	"testing"

	"github.com/gemyago/aws-sqs-boilerplate-go/internal/api/http/v1routes/handlers"
	"github.com/gemyago/aws-sqs-boilerplate-go/internal/api/http/v1routes/models"
	"github.com/gemyago/aws-sqs-boilerplate-go/internal/services"
	"github.com/go-faker/faker/v4"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
)

func TestCommands(t *testing.T) {
	randomMessagesPublishMessageRequest := func() *handlers.MessagesPublishMessageRequest {
		return &handlers.MessagesPublishMessageRequest{
			Payload: &models.Message{
				Id:       faker.UUIDHyphenated(),
				Name:     faker.Name(),
				Comments: faker.Sentence(),
			},
		}
	}

	type mockCommandDeps struct {
		mockMessageSender *services.MockMessageSender
		deps              CommandsDeps
	}

	makeMockDeps := func(t *testing.T) mockCommandDeps {
		mockMessageSender := services.NewMockMessageSender(t)
		return mockCommandDeps{
			mockMessageSender: mockMessageSender,
			deps: CommandsDeps{
				SendMessage: mockMessageSender.Execute,
			},
		}
	}

	t.Run("PublishMessage", func(t *testing.T) {
		t.Run("should send message", func(t *testing.T) {
			ctx := context.Background()
			deps := makeMockDeps(t)
			commands := NewCommands(deps.deps)
			req := randomMessagesPublishMessageRequest()
			deps.mockMessageSender.EXPECT().Execute(ctx, lo.ToPtr(services.Message(*req.Payload))).Return(nil)
			err := commands.PublishMessage(ctx, req)
			require.NoError(t, err)
		})

		t.Run("should return error if failed to send message", func(t *testing.T) {
			ctx := context.Background()
			deps := makeMockDeps(t)
			commands := NewCommands(deps.deps)
			req := randomMessagesPublishMessageRequest()
			expectedErr := errors.New(faker.Sentence())
			deps.mockMessageSender.EXPECT().Execute(ctx, lo.ToPtr(services.Message(*req.Payload))).Return(expectedErr)
			err := commands.PublishMessage(ctx, req)
			require.ErrorIs(t, err, expectedErr)
		})
	})
}