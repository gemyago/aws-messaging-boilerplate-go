package v1controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gemyago/aws-messaging-boilerplate-go/internal/api/http/server"
	"github.com/gemyago/aws-messaging-boilerplate-go/internal/api/http/v1routes/handlers"
	"github.com/gemyago/aws-messaging-boilerplate-go/internal/api/http/v1routes/models"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestMessages(t *testing.T) {
	makeDeps := func(t *testing.T) MessagesControllerDeps {
		return MessagesControllerDeps{
			Commands: newMockDummyMessagesCommands(t),
		}
	}

	newRouter := func(deps MessagesControllerDeps) http.Handler {
		mux := http.NewServeMux()
		app := handlers.NewHTTPApp(server.NewMuxRouterAdapter(mux))
		ctrl := NewMessagesController(deps)
		handlers.RegisterMessagesRoutes(ctrl, app)
		return mux
	}

	t.Run("NewMessagesController", func(t *testing.T) {
		t.Run("should initialize controller", func(t *testing.T) {
			deps := makeDeps(t)
			ctrl := NewMessagesController(deps)
			assert.NotNil(t, ctrl)
		})
	})

	t.Run("GET /health", func(t *testing.T) {
		t.Run("should return 204", func(t *testing.T) {
			deps := makeDeps(t)
			router := newRouter(deps)
			req := httptest.NewRequest(http.MethodGet, "/health", http.NoBody)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, 204, w.Code)
		})
	})

	t.Run("POST /dummy-messages/process", func(t *testing.T) {
		randomMessage := func() *models.DummyMessage {
			return &models.DummyMessage{
				Message: faker.Sentence(),
			}
		}
		t.Run("should process dummy message", func(t *testing.T) {
			deps := makeDeps(t)
			router := newRouter(deps)
			dummyMessage := randomMessage()

			mockCommands, _ := deps.Commands.(*mockDummyMessagesCommands)
			mockCommands.EXPECT().ProcessMessage(mock.AnythingOfType("context.backgroundCtx"), dummyMessage).Return(nil)

			body, err := json.Marshal(dummyMessage)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/dummy-messages/process", bytes.NewBuffer(body))
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, 204, w.Code)
		})
	})
}
