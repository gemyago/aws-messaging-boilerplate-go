package awsapi

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/go-faker/faker/v4"
	"github.com/samber/lo"
	"github.com/spf13/viper"
)

type testMessage struct {
	Id       string `json:"id"` //nolint:revive,stylecheck // Id is used to match apigen generated code
	Name     string `json:"name"`
	Comments string `json:"comments,omitempty"`
}

func newRandomMessage() *testMessage {
	return &testMessage{
		Id:       faker.UUIDHyphenated(),
		Name:     faker.Name(),
		Comments: faker.Sentence(),
	}
}

func newTestAWSConfig(ctx context.Context, appCfg *viper.Viper) aws.Config {
	return lo.Must(newAWSConfigFactory(ctx)(AWSConfigDeps{
		Region:       appCfg.GetString("aws.region"),
		BaseEndpoint: appCfg.GetString("aws.baseEndpoint"),
	}))
}
