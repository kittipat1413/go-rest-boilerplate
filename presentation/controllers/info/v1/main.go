package controllers

import (
	infousecase "go-rest-boilerplate/domain/usecase/info"
	"go-rest-boilerplate/presentation/controllers"
	"go-rest-boilerplate/presentation/middleware"

	"github.com/julienschmidt/httprouter"
)

type InfoController struct {
	infoUsecase infousecase.InfoUsecase
}

func NewInfoController(infoUsecase infousecase.InfoUsecase) InfoController {
	return InfoController{
		infoUsecase: infoUsecase,
	}
}

func (c InfoController) Route(router *httprouter.Router) []controllers.RouteBuilder {
	infoController := []controllers.RouteBuilder{
		{Router: router.GET, URI: "/info", Middlewares: []middleware.Middleware{}, Handler: c.FindAllInfo},
	}
	return infoController
}
