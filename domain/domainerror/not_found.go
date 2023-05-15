package domainerror

import (
	"fmt"
	"go-rest-boilerplate/constants"
	"net/http"
)

type NotFoundError struct{}

func (e NotFoundError) Code() string {
	return constants.StatusCodeGenericNotFoundError
}
func (e NotFoundError) Error() string {
	return e.GetMessage()
}
func (e NotFoundError) GetMessage() string {
	return "not found"
}
func (e NotFoundError) GetHttpCode() int {
	return http.StatusNotFound
}
func (e NotFoundError) Wrap(err error) error {
	return fmt.Errorf("%w: %v", e, err)
}
