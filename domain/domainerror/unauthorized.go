package domainerror

import (
	"fmt"
	"go-rest-boilerplate/constants"
	"net/http"
	"strings"

	"github.com/stoewer/go-strcase"
)

type UnauthorizedError struct{}

func (e UnauthorizedError) Code() string {
	return constants.StatusCodeUnauthorized
}
func (e UnauthorizedError) Error() string {
	return e.GetMessage()
}
func (e UnauthorizedError) GetMessage() string {
	return "unauthorized"
}
func (e UnauthorizedError) GetHttpCode() int {
	return http.StatusUnauthorized
}
func (e UnauthorizedError) Wrap(err error) error {
	return fmt.Errorf("%w: %v", e, err)
}

type UnauthorizedTokenExpired struct {
	TokenType []string
	UnauthorizedError
}

func (e UnauthorizedTokenExpired) Code() string {
	return constants.StatusCodeTokenExpired
}
func (e UnauthorizedTokenExpired) Error() string {
	return e.GetMessage()
}
func (e UnauthorizedTokenExpired) GetMessage() string {
	var tokenType []string
	for _, tkType := range e.TokenType {
		tokenType = append(tokenType, strcase.LowerCamelCase(tkType))
	}
	return fmt.Sprintf("%s token has expired", strings.Join(tokenType, ", "))
}
func (e UnauthorizedTokenExpired) Wrap(err error) error {
	return fmt.Errorf("%w: %v", e, err)
}
