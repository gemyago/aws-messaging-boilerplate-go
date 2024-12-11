package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/gemyago/aws-sqs-boilerplate-go/internal/config"
	"github.com/gemyago/aws-sqs-boilerplate-go/internal/diag"
	"github.com/go-faker/faker/v4"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMessageSender(t *testing.T) {
	appCfg := config.LoadTestConfig()
	ctx := context.Background()
	awsCfg := lo.Must(newAWSConfigFactory(ctx)(AWSConfigDeps{
		Region:       appCfg.GetString("aws.region"),
		BaseEndpoint: appCfg.GetString("aws.baseEndpoint"),
	}))
	sqsClient := sqs.NewFromConfig(awsCfg)
	queueURL := appCfg.GetString("aws.sqs.messagesQueueUrl")
	sender := NewMessageSender(MessageSenderDeps{
		SqsClient:        sqsClient,
		RootLogger:       diag.RootTestLogger(),
		MessagesQueueURL: queueURL,
	})

	newRandomMessage := func() *Message {
		return &Message{
			Id:       faker.UUIDHyphenated(),
			Name:     faker.Name(),
			Comments: faker.Sentence(),
		}
	}

	t.Run("should send the message to the queue", func(t *testing.T) {
		message := newRandomMessage()
		lo.Must(sqsClient.PurgeQueue(ctx, &sqs.PurgeQueueInput{
			QueueUrl: aws.String(queueURL),
		}))
		err := sender(context.Background(), message)
		require.NoError(t, err)

		res, err := sqsClient.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl: aws.String(queueURL),
		})
		require.NoError(t, err)
		require.Len(t, res.Messages, 1)
		var receivedMessage Message
		require.NoError(t, json.Unmarshal([]byte(*res.Messages[0].Body), &receivedMessage))
		assert.Equal(t, message, &receivedMessage)
	})
}
