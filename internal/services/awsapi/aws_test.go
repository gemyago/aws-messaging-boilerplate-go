package awsapi

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/gemyago/aws-sqs-boilerplate-go/internal/config"
	"github.com/gemyago/aws-sqs-boilerplate-go/internal/diag"
	"github.com/go-faker/faker/v4"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newRandomMessage() *Message {
	return &Message{
		Id:       faker.UUIDHyphenated(),
		Name:     faker.Name(),
		Comments: faker.Sentence(),
	}
}

func TestAWSMessageSender(t *testing.T) {
	appCfg := config.LoadTestConfig()
	ctx := context.Background()
	awsCfg := lo.Must(newAWSConfigFactory(ctx)(AWSConfigDeps{
		Region:       appCfg.GetString("aws.region"),
		BaseEndpoint: appCfg.GetString("aws.baseEndpoint"),
	}))
	sqsClient := sqs.NewFromConfig(awsCfg)
	snsClient := sns.NewFromConfig(awsCfg)
	topicARN := appCfg.GetString("aws.sns.messagesTopicArn")
	queueURL := appCfg.GetString("aws.sqs.messagesQueueUrl")
	sender := NewMessageSender(MessageSenderDeps{
		SnsClient:        snsClient,
		RootLogger:       diag.RootTestLogger(),
		MessagesTopicARN: topicARN,
	})

	t.Run("should send the message to sns", func(t *testing.T) {
		message := newRandomMessage()
		lo.Must(sqsClient.PurgeQueue(ctx, &sqs.PurgeQueueInput{
			QueueUrl: aws.String(queueURL),
		}))
		err := sender(context.Background(), message)
		require.NoError(t, err)

		res, err := sqsClient.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(queueURL),
			MaxNumberOfMessages: 1,
			WaitTimeSeconds:     1,
		})
		require.NoError(t, err)
		require.Len(t, res.Messages, 1)
		var receivedMessage Message
		require.NoError(t, json.Unmarshal([]byte(*res.Messages[0].Body), &receivedMessage))
		assert.Equal(t, message, &receivedMessage)
	})
}

func TestAWSMessagesPoller(t *testing.T) {
	appCfg := config.LoadTestConfig()
	rootCtx := context.Background()
	awsCfg := lo.Must(newAWSConfigFactory(rootCtx)(AWSConfigDeps{
		Region:       appCfg.GetString("aws.region"),
		BaseEndpoint: appCfg.GetString("aws.baseEndpoint"),
	}))
	sqsClient := sqs.NewFromConfig(awsCfg)
	snsClient := sns.NewFromConfig(awsCfg)
	queueURL := appCfg.GetString("aws.sqs.messagesQueueUrl")
	topicARN := appCfg.GetString("aws.sns.messagesTopicArn")
	sender := NewMessageSender(MessageSenderDeps{
		SnsClient:        snsClient,
		RootLogger:       diag.RootTestLogger(),
		MessagesTopicARN: topicARN,
	})

	t.Run("should receive the message from the queue", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		message := newRandomMessage()
		lo.Must(sqsClient.PurgeQueue(ctx, &sqs.PurgeQueueInput{
			QueueUrl: aws.String(queueURL),
		}))
		poller := NewMessagesPoller(MessagesPollerDeps{
			SqsClient:  sqsClient,
			RootLogger: diag.RootTestLogger(),
		})
		handledMessage := make(chan *Message)
		poller.RegisterQueue(MessagesPollerQueue{
			QueueURL: queueURL,
			Handler: NewRawMessageHandler(func(_ context.Context, message *Message) error {
				handledMessage <- message
				return nil
			}),
		})

		startComplete := make(chan error)
		go func() {
			startComplete <- poller.Start(ctx)
		}()

		err := sender(ctx, message)
		require.NoError(t, err)

		receivedMessage := <-handledMessage
		assert.Equal(t, message, receivedMessage)

		cancel()
		require.NoError(t, <-startComplete)
	})

	t.Run("should stop the poller if the context is canceled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		lo.Must(sqsClient.PurgeQueue(ctx, &sqs.PurgeQueueInput{
			QueueUrl: aws.String(queueURL),
		}))
		poller := NewMessagesPoller(MessagesPollerDeps{
			SqsClient:  sqsClient,
			RootLogger: diag.RootTestLogger(),
		})
		poller.RegisterQueue(MessagesPollerQueue{
			QueueURL: queueURL,
			Handler: NewRawMessageHandler(func(_ context.Context, _ *Message) error {
				return nil
			}),
		})

		startComplete := make(chan error)
		go func() {
			startComplete <- poller.Start(ctx)
		}()

		cancel()
		require.NoError(t, <-startComplete)
	})
}
