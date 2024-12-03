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

type paramsParserPetsCreatePet struct {
	bindPayload requestParamBinder[*http.Request, *Pet]
}

func (p *paramsParserPetsCreatePet) parse(router httpRouter, req *http.Request) (*PetsCreatePetRequest, error) {
	bindingCtx := BindingContext{}
	reqParams := &PetsCreatePetRequest{}
	// body params
	p.bindPayload(bindingCtx.Fork("body"), readRequestBodyValue(req), &reqParams.Payload)
	return reqParams, bindingCtx.AggregatedError()
}

func newParamsParserPetsCreatePet(app *HTTPApp) paramsParser[*PetsCreatePetRequest] {
	return &paramsParserPetsCreatePet{
		bindPayload: newRequestParamBinder(binderParams[*http.Request, *Pet]{
			required: true,
			parseValue: parseSoloValueParamAsSoloValue(
				parseJSONPayload[*Pet],
			),
			validateValue: NewPetValidator(),
		}),
	}
}

type paramsParserPetsGetPetById struct {
	bindPetId requestParamBinder[string, int64]
}

func (p *paramsParserPetsGetPetById) parse(router httpRouter, req *http.Request) (*PetsGetPetByIdRequest, error) {
	bindingCtx := BindingContext{}
	reqParams := &PetsGetPetByIdRequest{}
	// path params
	pathParamsCtx := bindingCtx.Fork("path")
	p.bindPetId(pathParamsCtx.Fork("petId"), readPathValue("petId", router, req), &reqParams.PetId)
	return reqParams, bindingCtx.AggregatedError()
}

func newParamsParserPetsGetPetById(app *HTTPApp) paramsParser[*PetsGetPetByIdRequest] {
	return &paramsParserPetsGetPetById{
		bindPetId: newRequestParamBinder(binderParams[string, int64]{
			required: true,
			parseValue: parseSoloValueParamAsSoloValue(
				app.knownParsers.int64Parser,
			),
			validateValue: NewSimpleFieldValidator[int64](
			),
		}),
	}
}

type paramsParserPetsListPets struct {
	bindLimit requestParamBinder[[]string, int64]
	bindOffset requestParamBinder[[]string, int64]
}

func (p *paramsParserPetsListPets) parse(router httpRouter, req *http.Request) (*PetsListPetsRequest, error) {
	bindingCtx := BindingContext{}
	reqParams := &PetsListPetsRequest{}
	// query params
	query := req.URL.Query()
	queryParamsCtx := bindingCtx.Fork("query")
	p.bindLimit(queryParamsCtx.Fork("limit"), readQueryValue("limit", query), &reqParams.Limit)
	p.bindOffset(queryParamsCtx.Fork("offset"), readQueryValue("offset", query), &reqParams.Offset)
	return reqParams, bindingCtx.AggregatedError()
}

func newParamsParserPetsListPets(app *HTTPApp) paramsParser[*PetsListPetsRequest] {
	return &paramsParserPetsListPets{
		bindLimit: newRequestParamBinder(binderParams[[]string, int64]{
			required: true,
			parseValue: parseMultiValueParamAsSoloValue(
				app.knownParsers.int64Parser,
			),
			validateValue: NewSimpleFieldValidator[int64](
				NewMinMaxValueValidator[int64](1, false, true),
				NewMinMaxValueValidator[int64](100, false, false),
			),
		}),
		bindOffset: newRequestParamBinder(binderParams[[]string, int64]{
			required: false,
			parseValue: parseMultiValueParamAsSoloValue(
				app.knownParsers.int64Parser,
			),
			validateValue: NewSimpleFieldValidator[int64](
				NewMinMaxValueValidator[int64](1, false, true),
			),
		}),
	}
}
