package domainerror

import (
	"fmt"
	"go-rest-boilerplate/constants"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/stoewer/go-strcase"
)

type BadRequestError struct{}

func (e BadRequestError) Code() string {
	return constants.StatusCodeGenericBadRequestError
}
func (e BadRequestError) Error() string {
	return e.GetMessage()
}
func (e BadRequestError) GetMessage() string {
	return "bad request"
}
func (e BadRequestError) GetHttpCode() int {
	return http.StatusBadRequest
}
func (e BadRequestError) Wrap(err error) error {
	return fmt.Errorf("%w: %v", e, err)
}

type RequiredParametersMissing struct {
	BadRequestError
}

func (e RequiredParametersMissing) Code() string {
	return constants.StatusCodeMissingRequiredParameters
}
func (e RequiredParametersMissing) Error() string {
	return e.GetMessage()
}
func (e RequiredParametersMissing) GetMessage() string {
	return "required parameter(s) are missing"
}
func (e RequiredParametersMissing) Wrap(err error) error {
	return fmt.Errorf("%w: %v", e, err)
}

type DuplicateError struct {
	fieldNames []string
	BadRequestError
}

func (e DuplicateError) Code() string {
	return constants.StatusCodeDuplicatedEntry
}
func (e DuplicateError) Error() string {
	return e.GetMessage()
}
func (e DuplicateError) GetMessage() string {
	var errFields []string
	for _, fe := range e.fieldNames {
		errFields = append(errFields, strcase.SnakeCase(fe))
	}
	return fmt.Sprintf("%s already exists", strings.Join(errFields, ", "))
}
func (e DuplicateError) Wrap(err error) error {
	return fmt.Errorf("%w: %v", e, err)
}

type ValidationError struct {
	ValidatorErrors validator.ValidationErrors
	BadRequestError
}

func (e ValidationError) Code() string {
	return constants.StatusCodeInvalidParameters
}
func (e ValidationError) Error() string {
	return e.GetMessage()
}
func (e ValidationError) GetMessage() string {
	var errFields []string
	for _, fe := range e.ValidatorErrors {
		errFields = append(errFields, strcase.LowerCamelCase(fe.Field()))
	}
	return fmt.Sprintf("%s has invalid format", strings.Join(errFields, ", "))
}
func (e ValidationError) Wrap(err error) error {
	return fmt.Errorf("%w: %v", e, err)
}
