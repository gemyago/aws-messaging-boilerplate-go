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
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSNSMessageSender(t *testing.T) {
	appCfg := config.LoadTestConfig()
	ctx := context.Background()
	awsCfg := newTestAWSConfig(ctx, appCfg)
	sqsClient := sqs.NewFromConfig(awsCfg)
	topicARN := appCfg.GetString("aws.sns.dummyMessagesTopicArn")
	queueURL := appCfg.GetString("aws.sqs.dummyMessagesQueueUrl")
	sender := NewSNSMessageSender[testMessage](topicARN, SNSMessageSenderDeps{
		SnsClient:  sns.NewFromConfig(awsCfg),
		RootLogger: diag.RootTestLogger(),
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
		var receivedMessage testMessage
		require.NoError(t, json.Unmarshal([]byte(*res.Messages[0].Body), &receivedMessage))
		assert.Equal(t, message, &receivedMessage)
	})
}
