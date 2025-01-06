// Code generated by apigen DO NOT EDIT.

package internal

import (
	"time"
	. "github.com/gemyago/aws-sqs-boilerplate-go/internal/api/http/v1routes/models"
)

// Below is to workaround unused imports.
var _ = time.Time{}

func NewDummyMessageValidator() FieldValidator[*DummyMessage] {
	validateMessage := NewSimpleFieldValidator[string](
		EnsureNonDefault[string],
	)
	validateFailProcessing := NewSimpleFieldValidator[bool](
	)
	
	return func(bindingCtx *BindingContext, value *DummyMessage) {
		validateMessage(bindingCtx.Fork("message"), value.Message)
		validateFailProcessing(bindingCtx.Fork("failProcessing"), value.FailProcessing)
	}
}
