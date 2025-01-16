package awsapi

import (
	"context"
	"encoding/json"

	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/gemyago/aws-messaging-boilerplate-go/internal/config"
	"github.com/gemyago/aws-messaging-boilerplate-go/internal/diag"
	"github.com/go-faker/faker/v4"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

func TestNewRawMessageHandler(t *testing.T) {
	type TestMessage struct {
		Data string `json:"data"`
	}

	t.Run("should unmarshal the message and delegate to target handler", func(t *testing.T) {
		wantMessage := &TestMessage{Data: faker.Sentence()}
		handlerInvoked := false
		handler := NewRawMessageHandler(func(_ context.Context, message *TestMessage) error {
			assert.Equal(t, wantMessage, message)
			handlerInvoked = true
			return nil
		})
		rawMessage := types.Message{
			Body: lo.ToPtr(string(lo.Must(json.Marshal(wantMessage)))),
		}
		gotErr := handler(context.Background(), rawMessage)
		require.NoError(t, gotErr)
		assert.True(t, handlerInvoked)
	})

	t.Run("should return an error if the message is not valid JSON", func(t *testing.T) {
		handler := NewRawMessageHandler(func(_ context.Context, _ *TestMessage) error {
			return nil
		})
		rawMessage := types.Message{
			Body: lo.ToPtr(faker.Word()),
		}
		gotErr := handler(context.Background(), rawMessage)
		assert.Error(t, gotErr)
	})
}
