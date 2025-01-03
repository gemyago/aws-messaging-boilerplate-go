package awsapi

import (
	"context"

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

func newRandomMessage() *testMessage {
	return &testMessage{
		Id:       faker.UUIDHyphenated(),
		Name:     faker.Name(),
		Comments: faker.Sentence(),
	}
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
	queueURL := appCfg.GetString("aws.sqs.dummyMessagesQueueUrl")
	topicARN := appCfg.GetString("aws.sns.dummyMessagesTopicArn")
	sender := NewSNSMessageSender[testMessage](topicARN, SNSMessageSenderDeps{
		SnsClient:  snsClient,
		RootLogger: diag.RootTestLogger(),
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
		handledMessage := make(chan *testMessage)
		poller.RegisterQueue(MessagesPollerQueue{
			QueueURL: queueURL,
			Handler: NewRawMessageHandler(func(_ context.Context, message *testMessage) error {
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
			Handler: NewRawMessageHandler(func(_ context.Context, _ *testMessage) error {
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
