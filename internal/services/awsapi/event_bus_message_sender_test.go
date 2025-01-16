package awsapi

import (
	"context"
	"encoding/json"
	"net"
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

	t.Run("send message to event bus", func(t *testing.T) {
		gotMessages := make(chan testMessage, 1)
		testListener, err := net.Listen("tcp", "0.0.0.0:48080")
		require.NoError(t, err)

		testSrv := httptest.NewUnstartedServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var message testMessage
				if err = json.NewDecoder(r.Body).Decode(&message); err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				gotMessages <- message
				w.WriteHeader(http.StatusOK)
			}),
		)
		testSrv.Listener.Close()
		testSrv.Listener = testListener
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
			require.Equal(t, *wantMsg, gotMsg)
		case <-time.After(1 * time.Second):
			require.Fail(t, "timeout waiting for message")
		}
	})
}
