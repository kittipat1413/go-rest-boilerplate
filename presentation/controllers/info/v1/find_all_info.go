package controllers

import (
	"net/http"

	"go-rest-boilerplate/presentation/render"

	"github.com/julienschmidt/httprouter"
)

func (c InfoController) FindAllInfo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if infos, err := c.infoUsecase.FindAllInfo(r.Context()); err != nil {
		render.Error(w, r, err)
	} else {
		render.JSON(w, r, infos)
	}
}
