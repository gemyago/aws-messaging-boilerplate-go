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

type testMessage struct {
	Id       string `json:"id"` //nolint:revive,stylecheck // Id is used to match apigen generated code
	Name     string `json:"name"`
	Comments string `json:"comments,omitempty"`
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
	sender := NewSNSMessageSender[testMessage](topicARN, SNSMessageSenderDeps{
		SnsClient:  snsClient,
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
