// Code generated by apigen DO NOT EDIT.

package handlers

import (
	"net/http"
	"time"

	. "github.com/gemyago/aws-sqs-boilerplate-go/internal/api/http/v1routes/models"
	. "github.com/gemyago/aws-sqs-boilerplate-go/internal/api/http/v1routes/internal"
)

// Below is to workaround unused imports.
var _ = time.Time{}
type _ func() Error

type paramsParserMessagesProcessMessage struct {
	bindPayload requestParamBinder[*http.Request, *Message]
}

func (p *paramsParserMessagesProcessMessage) parse(router httpRouter, req *http.Request) (*MessagesProcessMessageRequest, error) {
	bindingCtx := BindingContext{}
	reqParams := &MessagesProcessMessageRequest{}
	// body params
	p.bindPayload(bindingCtx.Fork("body"), readRequestBodyValue(req), &reqParams.Payload)
	return reqParams, bindingCtx.AggregatedError()
}

func newParamsParserMessagesProcessMessage(app *HTTPApp) paramsParser[*MessagesProcessMessageRequest] {
	return &paramsParserMessagesProcessMessage{
		bindPayload: newRequestParamBinder(binderParams[*http.Request, *Message]{
			required: true,
			parseValue: parseSoloValueParamAsSoloValue(
				parseJSONPayload[*Message],
			),
			validateValue: NewMessageValidator(),
		}),
	}
}

type paramsParserMessagesPublishMessage struct {
	bindPayload requestParamBinder[*http.Request, *Message]
}

func (p *paramsParserMessagesPublishMessage) parse(router httpRouter, req *http.Request) (*MessagesPublishMessageRequest, error) {
	bindingCtx := BindingContext{}
	reqParams := &MessagesPublishMessageRequest{}
	// body params
	p.bindPayload(bindingCtx.Fork("body"), readRequestBodyValue(req), &reqParams.Payload)
	return reqParams, bindingCtx.AggregatedError()
}

func newParamsParserMessagesPublishMessage(app *HTTPApp) paramsParser[*MessagesPublishMessageRequest] {
	return &paramsParserMessagesPublishMessage{
		bindPayload: newRequestParamBinder(binderParams[*http.Request, *Message]{
			required: true,
			parseValue: parseSoloValueParamAsSoloValue(
				parseJSONPayload[*Message],
			),
			validateValue: NewMessageValidator(),
		}),
	}
}
