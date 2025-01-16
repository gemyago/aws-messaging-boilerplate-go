package queues

import (
	"context"
	"math/rand/v2"
	"testing"

	"github.com/gemyago/aws-messaging-boilerplate-go/internal/services"
	"github.com/gemyago/aws-messaging-boilerplate-go/internal/services/awsapi"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestStartPolling(t *testing.T) {
	makeDeps := func(t *testing.T) Deps {
		return Deps{
			DummyMessagesQueueURL:                  faker.URL(),
			DummyMessagesQueueVisibilityTimeoutSec: rand.Int32(),
			MessagesPoller:                         newMockMessagePoller(t),
			ShutdownHooks:                          services.NewTestShutdownHooks(),
		}
	}

	t.Run("should register queues and start polling", func(t *testing.T) {
		deps := makeDeps(t)

		mockPoller, _ := deps.MessagesPoller.(*mockMessagePoller)
		mockPoller.EXPECT().RegisterQueue(mock.MatchedBy(
			func(queue awsapi.MessagesPollerQueue) bool {
				return assert.Equal(t, deps.DummyMessagesQueueURL, queue.QueueURL) &&
					assert.NotNil(t, queue.Handler) &&
					assert.Equal(t, deps.DummyMessagesQueueVisibilityTimeoutSec, queue.VisibilityTimeout)
			},
		))
		mockPoller.EXPECT().Start(mock.Anything).Return(nil)

		ctx := context.Background()
		err := StartPolling(ctx, deps)
		require.NoError(t, err)
		require.NoError(t, deps.ShutdownHooks.PerformShutdown(ctx))
	})
}
