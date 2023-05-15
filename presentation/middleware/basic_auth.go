package middleware

import (
	"crypto/sha256"
	"crypto/subtle"
	"net/http"

	"go-rest-boilerplate/config"
	"go-rest-boilerplate/domain/domainerror"
	"go-rest-boilerplate/presentation/render"

	"github.com/julienschmidt/httprouter"
)

func BasicAuth() Middleware {
	return func(next httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			backofficeToken := r.Header.Get("Authorization")
			if backofficeToken == "" {
				render.Error(w, r, new(domainerror.UnauthorizedError))
				return
			}
			key, secret, ok := r.BasicAuth()
			if ok {
				usernameHash := sha256.Sum256([]byte(key))
				passwordHash := sha256.Sum256([]byte(secret))

				cfg := config.FromRequest(r)
				expectedUsernameHash := sha256.Sum256([]byte(cfg.BackofficeApiKey()))
				expectedPasswordHash := sha256.Sum256([]byte(cfg.BackofficeApiSecret()))

				usernameMatch := (subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1)
				passwordMatch := (subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1)

				if usernameMatch && passwordMatch {
					next(w, r, p)
				} else {
					render.Error(w, r, new(domainerror.UnauthorizedError))
				}
			} else {
				render.Error(w, r, new(domainerror.UnauthorizedError))
			}
		}
	}
}
