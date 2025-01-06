package app

import (
	"context"
	"errors"
	"testing"

	"github.com/gemyago/aws-sqs-boilerplate-go/internal/api/http/v1routes/handlers"
	"github.com/gemyago/aws-sqs-boilerplate-go/internal/api/http/v1routes/models"
	"github.com/gemyago/aws-sqs-boilerplate-go/internal/diag"
	"github.com/gemyago/aws-sqs-boilerplate-go/internal/services/awsapi"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/require"
)

func TestCommands(t *testing.T) {
	randomMessage := func() *models.DummyMessage {
		return &models.DummyMessage{
			Message: faker.Sentence(),
		}
	}
	randomMessagesPublishMessageRequest := func() *handlers.MessagesPublishDummyMessageRequest {
		return &handlers.MessagesPublishDummyMessageRequest{
			Target:  handlers.MessagesPublishDummyMessageTargetSNS,
			Payload: randomMessage(),
		}
	}

	type mockCommandDeps struct {
		mockSNSMessageSender         *awsapi.MockMessageSender[models.DummyMessage]
		mockEventBridgeMessageSender *awsapi.MockMessageSender[models.DummyMessage]
		deps                         CommandsDeps
	}

	makeMockDeps := func(t *testing.T) mockCommandDeps {
		mockSNSMessageSender := awsapi.NewMockMessageSender[models.DummyMessage](t)
		mockEventBridgeMessageSender := awsapi.NewMockMessageSender[models.DummyMessage](t)
		return mockCommandDeps{
			mockSNSMessageSender:         mockSNSMessageSender,
			mockEventBridgeMessageSender: mockEventBridgeMessageSender,
			deps: CommandsDeps{
				RootLogger:                  diag.RootTestLogger(),
				SendDummySNSMessage:         mockSNSMessageSender.Execute,
				SendDummyEventBridgeMessage: mockEventBridgeMessageSender.Execute,
			},
		}
	}

	t.Run("PublishMessage", func(t *testing.T) {
		t.Run("should send SNS message", func(t *testing.T) {
			ctx := context.Background()
			deps := makeMockDeps(t)
			commands := NewCommands(deps.deps)
			req := randomMessagesPublishMessageRequest()
			deps.mockSNSMessageSender.EXPECT().Execute(ctx, req.Payload).Return(nil)
			err := commands.PublishMessage(ctx, req)
			require.NoError(t, err)
		})
		t.Run("should send EventBridge message", func(t *testing.T) {
			ctx := context.Background()
			deps := makeMockDeps(t)
			commands := NewCommands(deps.deps)
			req := randomMessagesPublishMessageRequest()
			req.Target = handlers.MessagesPublishDummyMessageTargetEVENTBRIDGE
			deps.mockEventBridgeMessageSender.EXPECT().Execute(ctx, req.Payload).Return(nil)
			err := commands.PublishMessage(ctx, req)
			require.NoError(t, err)
		})

		t.Run("should return error if failed to send message", func(t *testing.T) {
			ctx := context.Background()
			deps := makeMockDeps(t)
			commands := NewCommands(deps.deps)
			req := randomMessagesPublishMessageRequest()
			expectedErr := errors.New(faker.Sentence())
			deps.mockSNSMessageSender.EXPECT().Execute(ctx, req.Payload).Return(expectedErr)
			err := commands.PublishMessage(ctx, req)
			require.ErrorIs(t, err, expectedErr)
		})
	})

	t.Run("ProcessMessage", func(t *testing.T) {
		t.Run("should be noop", func(t *testing.T) {
			ctx := context.Background()
			deps := makeMockDeps(t)
			commands := NewCommands(deps.deps)
			msg := randomMessage()
			err := commands.ProcessMessage(ctx, msg)
			require.NoError(t, err)
		})
	})
}
