// Code generated by apigen DO NOT EDIT.

package handlers

import (
	"encoding/json"
	"fmt"
	"time"

	. "github.com/gemyago/aws-sqs-boilerplate-go/internal/api/http/v1routes/models"
)

// Below is to workaround unused imports.
var _ = time.Time{}
var _ = json.Unmarshal
var _ = fmt.Sprint



// MessagesPublishMessageRequest represents params for publishMessage operation
//
// Request: POST /messages/publish.
type MessagesPublishMessageRequest struct {
	// Payload is parsed from request body and declared as payload.
	Payload *Message
}

type MessagesController struct {
	// GET /health
	//
	// Request type: none
	//
	// Response type: none
	HealthCheck httpHandlerFactory

	// POST /messages/publish
	//
	// Request type: MessagesPublishMessageRequest,
	//
	// Response type: none
	PublishMessage httpHandlerFactory
}

type MessagesControllerBuilder struct {
	// GET /health
	//
	// Request type: none
	//
	// Response type: none
	HandleHealthCheck actionBuilderNoParamsVoidResult[*MessagesControllerBuilder]

	// POST /messages/publish
	//
	// Request type: MessagesPublishMessageRequest,
	//
	// Response type: none
	HandlePublishMessage actionBuilderVoidResult[*MessagesControllerBuilder, *MessagesPublishMessageRequest]
}

func (c *MessagesControllerBuilder) Finalize() *MessagesController {
	return &MessagesController{
		HealthCheck: mustInitializeAction("healthCheck", c.HandleHealthCheck.httpHandlerFactory),
		PublishMessage: mustInitializeAction("publishMessage", c.HandlePublishMessage.httpHandlerFactory),
	}
}

func BuildMessagesController() *MessagesControllerBuilder {
	controllerBuilder := &MessagesControllerBuilder{}

	// GET /health
	controllerBuilder.HandleHealthCheck.controllerBuilder = controllerBuilder
	controllerBuilder.HandleHealthCheck.defaultStatusCode = 204
	controllerBuilder.HandleHealthCheck.voidResult = true

	// POST /messages/publish
	controllerBuilder.HandlePublishMessage.controllerBuilder = controllerBuilder
	controllerBuilder.HandlePublishMessage.defaultStatusCode = 202
	controllerBuilder.HandlePublishMessage.voidResult = true
	controllerBuilder.HandlePublishMessage.paramsParserFactory = newParamsParserMessagesPublishMessage

	return controllerBuilder
}

func RegisterMessagesRoutes(controller *MessagesController, app *HTTPApp) {
	app.router.HandleRoute("GET", "/health", controller.HealthCheck(app))
	app.router.HandleRoute("POST", "/messages/publish", controller.PublishMessage(app))
}