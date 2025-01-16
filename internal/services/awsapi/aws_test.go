package awsapi

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/gemyago/aws-messaging-boilerplate-go/internal/config"
	"github.com/go-faker/faker/v4"
	"github.com/samber/lo"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestNewAWSConfigFactory(t *testing.T) {
	t.Run("should return new config", func(t *testing.T) {
		ctx := context.Background()
		cfg := config.LoadTestConfig()
		awsCfg := newTestAWSConfig(ctx, cfg)
		assert.NotNil(t, awsCfg)
	})

	t.Run("should fail if load fails", func(t *testing.T) {
		ctx := context.Background()
		wantErr := errors.New(faker.Sentence())
		_, err := newAWSConfigFactory(ctx)(AWSConfigDeps{
			loadOpts: []func(*awsConfig.LoadOptions) error{
				func(_ *awsConfig.LoadOptions) error {
					return wantErr
				},
			},
		})
		require.Error(t, err)
		assert.ErrorIs(t, err, wantErr)
	})
}
