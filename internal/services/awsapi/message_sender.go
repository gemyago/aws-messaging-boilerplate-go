package awsapi

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"go.uber.org/dig"
)

type MessageSender[TMessage any] func(
	ctx context.Context,
	message *TMessage,
) error

type MessageSenderDeps struct {
	dig.In

	RootLogger       *slog.Logger
	SnsClient        *sns.Client
	MessagesTopicARN string `name:"config.aws.sns.messagesTopicARN"`
}

func NewMessageSender[TMessage any](deps MessageSenderDeps) MessageSender[TMessage] {
	logger := deps.RootLogger.WithGroup("services.message-sender")
	return func(ctx context.Context, message *TMessage) error {
		body, err := json.Marshal(message)
		if err != nil {
			return fmt.Errorf("failed to marshal message, %w", err)
		}
		res, err := deps.SnsClient.Publish(ctx, &sns.PublishInput{
			Message:  aws.String(string(body)),
			TopicArn: aws.String(deps.MessagesTopicARN),
		})
		if err != nil {
			return fmt.Errorf("failed send message to sqs queue, %w", err)
		}
		logger.InfoContext(ctx, "Message sent", slog.String("messageId", *res.MessageId))
		return nil
	}
}
