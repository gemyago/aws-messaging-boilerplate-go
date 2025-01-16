package v1controllers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gemyago/aws-messaging-boilerplate-go/internal/api/http/server"
	"github.com/gemyago/aws-messaging-boilerplate-go/internal/api/http/v1routes/handlers"
	"github.com/stretchr/testify/assert"
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
}
