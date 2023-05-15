package controllers

import (
	"go-rest-boilerplate/presentation/controllers"
	"go-rest-boilerplate/presentation/middleware"

	"github.com/julienschmidt/httprouter"
)

type HealthCheckController struct{}

func (c HealthCheckController) Route(router *httprouter.Router) []controllers.RouteBuilder {
	healthCheckControllers := []controllers.RouteBuilder{
		{Router: router.GET, URI: "/health/liveness", Middlewares: []middleware.Middleware{}, Handler: c.HealthCheck},
	}
	return healthCheckControllers
}
