package domainerror

import (
	"fmt"
	"go-rest-boilerplate/constants"
	"net/http"
)

type UnprocessableEntityError struct {
	Message string `json:"message"`
}

func (e UnprocessableEntityError) Code() string {
	return constants.StatusCodeUnprocessableEntity
}
func (e UnprocessableEntityError) Error() string {
	return e.GetMessage()
}
func (e UnprocessableEntityError) GetMessage() string {
	return e.Message
}
func (e UnprocessableEntityError) GetHttpCode() int {
	return http.StatusUnprocessableEntity
}
func (e UnprocessableEntityError) Wrap(err error) error {
	return fmt.Errorf("%w: %v", e, err)
}
