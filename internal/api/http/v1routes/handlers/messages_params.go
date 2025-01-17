// Code generated by apigen DO NOT EDIT.

package handlers

import (
	"net/http"
	"time"

	. "github.com/gemyago/aws-messaging-boilerplate-go/internal/api/http/v1routes/models"
	. "github.com/gemyago/aws-messaging-boilerplate-go/internal/api/http/v1routes/internal"
)

// Below is to workaround unused imports.
var _ = time.Time{}
type _ func() DummyMessage

type paramsParserMessagesProcessDummyMessage struct {
	bindPayload requestParamBinder[*http.Request, *DummyMessage]
}

func (p *paramsParserMessagesProcessDummyMessage) parse(router httpRouter, req *http.Request) (*MessagesProcessDummyMessageRequest, error) {
	bindingCtx := BindingContext{}
	reqParams := &MessagesProcessDummyMessageRequest{}
	// body params
	p.bindPayload(bindingCtx.Fork("body"), readRequestBodyValue(req), &reqParams.Payload)
	return reqParams, bindingCtx.AggregatedError()
}

func newParamsParserMessagesProcessDummyMessage(app *HTTPApp) paramsParser[*MessagesProcessDummyMessageRequest] {
	return &paramsParserMessagesProcessDummyMessage{
		bindPayload: newRequestParamBinder(binderParams[*http.Request, *DummyMessage]{
			required: true,
			parseValue: parseSoloValueParamAsSoloValue(
				parseJSONPayload[*DummyMessage],
			),
			validateValue: NewDummyMessageValidator(),
		}),
	}
}

type paramsParserMessagesPublishDummyMessage struct {
	bindTarget requestParamBinder[[]string, MessagesPublishDummyMessageTarget]
	bindPayload requestParamBinder[*http.Request, *DummyMessage]
}

func (p *paramsParserMessagesPublishDummyMessage) parse(router httpRouter, req *http.Request) (*MessagesPublishDummyMessageRequest, error) {
	bindingCtx := BindingContext{}
	reqParams := &MessagesPublishDummyMessageRequest{}
	// query params
	query := req.URL.Query()
	queryParamsCtx := bindingCtx.Fork("query")
	p.bindTarget(queryParamsCtx.Fork("target"), readQueryValue("target", query), &reqParams.Target)
	// body params
	p.bindPayload(bindingCtx.Fork("body"), readRequestBodyValue(req), &reqParams.Payload)
	return reqParams, bindingCtx.AggregatedError()
}

func newParamsParserMessagesPublishDummyMessage(app *HTTPApp) paramsParser[*MessagesPublishDummyMessageRequest] {
	return &paramsParserMessagesPublishDummyMessage{
		bindTarget: newRequestParamBinder(binderParams[[]string, MessagesPublishDummyMessageTarget]{
			required: true,
			parseValue: parseMultiValueParamAsSoloValue(
				ParseMessagesPublishDummyMessageTarget,
			),
			validateValue: NewSimpleFieldValidator[MessagesPublishDummyMessageTarget](
			),
		}),
		bindPayload: newRequestParamBinder(binderParams[*http.Request, *DummyMessage]{
			required: true,
			parseValue: parseSoloValueParamAsSoloValue(
				parseJSONPayload[*DummyMessage],
			),
			validateValue: NewDummyMessageValidator(),
		}),
	}
}
