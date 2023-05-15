package domainerror

import (
	"fmt"
	"go-rest-boilerplate/constants"
	"net/http"
)

type InternalError struct{}

func (e InternalError) Code() string {
	return constants.StatusCodeGenericInternalError
}
func (e InternalError) Error() string {
	return e.GetMessage()
}
func (e InternalError) GetMessage() string {
	return "internal server error"
}
func (e InternalError) GetHttpCode() int {
	return http.StatusInternalServerError
}
func (e InternalError) Wrap(err error) error {
	return fmt.Errorf("%w: %v", e, err)
}

type DatabaseError struct {
	InternalError
}

func (e DatabaseError) Code() string {
	return constants.StatusCodeDatabaseError
}
func (e DatabaseError) Error() string {
	return e.GetMessage()
}
func (e DatabaseError) GetMessage() string {
	return "database error(s)"
}
func (e DatabaseError) Wrap(err error) error {
	return fmt.Errorf("%w: %v", e, err)
}
