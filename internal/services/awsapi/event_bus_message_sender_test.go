package awsapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/gemyago/aws-sqs-boilerplate-go/internal/config"
	"github.com/gemyago/aws-sqs-boilerplate-go/internal/diag"
	"github.com/stretchr/testify/require"
)

func TestEventBusMessageSender(t *testing.T) {
	appCfg := config.LoadTestConfig()
	ctx := context.Background()
	awsCfg := newTestAWSConfig(ctx, appCfg)

	gotMessages := make(chan testMessage, 1)
	testSrv := httptest.NewUnstartedServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var message testMessage
			if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			gotMessages <- message
			w.WriteHeader(http.StatusOK)
		}),
	)
	testSrv.Config.Addr = "localhost:418080"
	testSrv.Start()
	defer testSrv.Close()

	wantMsg := newRandomMessage()
	sender := NewEventBusMessageSender[testMessage](
		appCfg.GetString("aws.eventBus.dummyMessagesDetailType"),
		EventBusMessageSenderDeps{
			EventBusName:   appCfg.GetString("aws.eventBus.name"),
			EventBusSource: appCfg.GetString("aws.eventBus.source"),
			RootLogger:     diag.RootTestLogger(),
			Client:         eventbridge.NewFromConfig(awsCfg),
		})
	require.NoError(t, sender(ctx, wantMsg))

	select {
	case gotMsg := <-gotMessages:
		require.Equal(t, wantMsg, gotMsg)
	case <-time.After(1 * time.Second):
		require.Fail(t, "timeout waiting for message")
	}
}
