package controllers

import (
	"go-rest-boilerplate/presentation/render"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (c HealthCheckController) HealthCheck(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	render.JSON(w, r, map[string]string{"message": "OK"})
}
