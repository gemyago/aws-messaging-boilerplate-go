// Code generated by apigen DO NOT EDIT.

package internal

import (
	"time"
	. "github.com/gemyago/aws-sqs-boilerplate-go/internal/api/http/v1routes/models"
)

// Below is to workaround unused imports.
var _ = time.Time{}

func NewDummyMessageValidator() FieldValidator[*DummyMessage] {
	validateId := NewSimpleFieldValidator[string](
		EnsureNonDefault[string],
	)
	validateName := NewSimpleFieldValidator[string](
		EnsureNonDefault[string],
	)
	validateComments := NewSimpleFieldValidator[string](
	)
	
	return func(bindingCtx *BindingContext, value *DummyMessage) {
		validateId(bindingCtx.Fork("id"), value.Id)
		validateName(bindingCtx.Fork("name"), value.Name)
		validateComments(bindingCtx.Fork("comments"), value.Comments)
	}
}
