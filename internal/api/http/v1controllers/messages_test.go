package v1controllers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessages(t *testing.T) {
	newDeps := func(t *testing.T) *MessagesControllerDeps {
		return &MessagesControllerDeps{
			Commands: newMockMessagesCommands(t),
		}
	}

	t.Run("NewMessagesController", func(t *testing.T) {
		t.Run("should initialize controller", func(t *testing.T) {
			deps := newDeps(t)
			ctrl := NewMessagesController(*deps)
			assert.NotNil(t, ctrl)
		})
	})
}
