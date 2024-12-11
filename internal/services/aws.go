package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"go.uber.org/dig"
)

//go:generate mockery --name=MessageSender --filename=mock_message_sender.go --config ../../.mockery-funcs.yaml

type AWSConfigDeps struct {
	Region string `config:"aws.region"`
}

func newAWSConfigFactory(ctx context.Context) func(deps AWSConfigDeps) (aws.Config, error) {
	return func(deps AWSConfigDeps) (aws.Config, error) {
		cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(deps.Region))
		if err != nil {
			return aws.Config{}, fmt.Errorf("failed to load aws configuration, %w", err)
		}
		return cfg, nil
	}
}

type sqsClient interface {
	SendMessage(
		ctx context.Context,
		params *sqs.SendMessageInput,
		optFns ...func(*sqs.Options),
	) (*sqs.SendMessageOutput, error)
}

type Message struct {
	Id       string `json:"id"` //nolint:revive,stylecheck // Id is used to match apigen generated code
	Name     string `json:"name"`
	Comments string `json:"comments,omitempty"`
}

type MessageSender func(ctx context.Context, message *Message) error

type MessageSenderDeps struct {
	dig.In

	RootLogger       *slog.Logger
	SqsClient        *sqs.Client
	MessagesQueueURL string `config:"aws.sqs.messagesQueueUrl"`
}

func NewMessageSender(deps MessageSenderDeps) MessageSender {
	logger := deps.RootLogger.WithGroup("services.message-sender")
	return func(ctx context.Context, message *Message) error {
		body, err := json.Marshal(message)
		if (err != nil) {
			return fmt.Errorf("failed to marshal message, %w", err)
		}
		res, err := deps.SqsClient.SendMessage(ctx, &sqs.SendMessageInput{
			MessageBody: aws.String(string(body)),
			QueueUrl:    aws.String(deps.MessagesQueueURL),
		})
		if err != nil {
			return fmt.Errorf("failed send message to sqs queue, %w", err)
		}
		logger.InfoContext(ctx, "message sent", slog.String("messageId", *res.MessageId))
		return nil
	}
}
