package middleware

import (
	"context"
	"go-rest-boilerplate/config"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

func TestBasicAuthMiddleware(t *testing.T) {
	var (
		CorrectAuthUser   = "correct_username"
		CorrectAuthPass   = "correct_password"
		IncorrectAuthUser = "incorrect_username"
		IncorrectAuthPass = "incorrect_password"
	)
	// helper function
	mockHandler := func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Mock Handler"))
	}
	setBasicAuth := func(req *http.Request, user *string, pass *string) {
		if user == nil || pass == nil {
			return
		}
		req.SetBasicAuth(*user, *pass)
	}

	// setup context with config
	cfg := config.MustConfigure()
	cfg.Viper.Set("BACKOFFICE_API_KEY", CorrectAuthUser)
	cfg.Viper.Set("BACKOFFICE_API_SECRET", CorrectAuthPass)
	ctx := config.NewContext(context.Background(), cfg)

	// tests
	type handlerArgs struct {
		BasicAuthUser *string
		BasicAuthPass *string
	}
	tests := []struct {
		name           string
		args           handlerArgs
		wantStatusCode int
	}{
		{
			name: "happy path",
			args: handlerArgs{
				BasicAuthUser: &CorrectAuthUser,
				BasicAuthPass: &CorrectAuthPass,
			},
			wantStatusCode: 200,
		},
		{
			name: "authorization header is empty ",
			args: handlerArgs{
				BasicAuthUser: nil,
				BasicAuthPass: nil,
			},
			wantStatusCode: http.StatusUnauthorized,
		},
		{
			name: "incorrect username and password",
			args: handlerArgs{
				BasicAuthUser: &IncorrectAuthUser,
				BasicAuthPass: &IncorrectAuthPass,
			},
			wantStatusCode: http.StatusUnauthorized,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/", nil).WithContext(ctx)
			setBasicAuth(request, test.args.BasicAuthUser, test.args.BasicAuthPass)
			response := httptest.NewRecorder()

			// Call the middleware with the mock handler
			BasicAuth()(mockHandler)(response, request, nil)

			if !assert.Equal(t, test.wantStatusCode, response.Code) {
				t.Errorf("middleware BasicAuth() Http Status Code want %v, got %v", test.wantStatusCode, response.Code)
			}

		})
	}
}
