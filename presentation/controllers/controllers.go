package controllers

import (
	"encoding/json"
	"io"

	"go-rest-boilerplate/domain/domainerror"
	"go-rest-boilerplate/presentation/middleware"

	"github.com/julienschmidt/httprouter"
)

type ControllerInterface interface {
	Route(router *httprouter.Router) []RouteBuilder
}

type RouteBuilder struct {
	Router      func(string, httprouter.Handle)
	URI         string
	Middlewares []middleware.Middleware
	Handler     httprouter.Handle
}

func (b *RouteBuilder) Build() (string, httprouter.Handle) {
	handler := b.buildChain(b.Handler, b.Middlewares...)
	return b.URI, handler
}

func (b *RouteBuilder) buildChain(final httprouter.Handle, middlewares ...middleware.Middleware) httprouter.Handle {
	if len(middlewares) == 0 {
		return final
	}
	return middlewares[0](b.buildChain(final, middlewares[1:]...))
}

func MountAll(router *httprouter.Router, controllers []ControllerInterface) error {
	for _, controller := range controllers {
		for _, builder := range controller.Route(router) {
			uri, handler := builder.Build()
			builder.Router(uri, handler)
		}
	}
	return nil
}

func ReadJSON(b io.ReadCloser, obj interface{}) error {
	err := json.NewDecoder(b).Decode(obj)
	if err != nil {
		return new(domainerror.RequiredParametersMissing)
	}
	return nil
}
