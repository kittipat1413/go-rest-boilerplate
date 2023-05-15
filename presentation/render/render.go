package render

import (
	"bytes"
	"encoding/json"
	"errors"
	"go-rest-boilerplate/config"
	"go-rest-boilerplate/constants"
	"go-rest-boilerplate/domain/domainerror"
	"go-rest-boilerplate/presentation/view"
	"net/http"
)

func Text(resp http.ResponseWriter, r *http.Request, text string) {
	resp.Header().Set("Content-Type", "text/plain")
	resp.WriteHeader(200)
	if _, err := resp.Write([]byte(text)); err != nil {
		Error(resp, r, new(domainerror.InternalError))
	}
}

func JSON(resp http.ResponseWriter, r *http.Request, obj interface{}) {
	successResponse := &struct {
		Code string      `json:"code"`
		Data interface{} `json:"data"`
	}{
		Code: constants.StatusCodeSuccess,
		Data: obj,
	}

	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(200)

	if err := json.NewEncoder(resp).Encode(successResponse); err != nil {
		Error(resp, r, new(domainerror.InternalError))
	}
}

func JSONWithCode(resp http.ResponseWriter, r *http.Request, code string, obj interface{}) {
	successResponse := &struct {
		Code string      `json:"code"`
		Data interface{} `json:"data"`
	}{
		Code: code,
		Data: obj,
	}

	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(200)

	if err := json.NewEncoder(resp).Encode(successResponse); err != nil {
		Error(resp, r, new(domainerror.InternalError))
	}
}

func Error(resp http.ResponseWriter, r *http.Request, err error) {
	resp.Header().Set("Content-Type", "application/json")
	errObj := unwrapError(err)
	resp.WriteHeader(errObj.HttpCode)

	if err_ := json.NewEncoder(resp).Encode(errObj); err_ != nil {
		config.FromRequest(r).Printf("%s %s %s - %s\n",
			r.RemoteAddr, r.Method, r.RequestURI, err_.Error())
	}
	cfg := config.FromRequest(r)
	cfg.Printf("%s %s %s %s\n", r.RemoteAddr, r.Method, r.RequestURI, err.Error())
}

func FileTransfer(resp http.ResponseWriter, r *http.Request, filename string, bytes *bytes.Buffer) {
	resp.Header().Set("Content-Description", "File Transfer")
	resp.Header().Set("Content-Transfer-Encoding", "binary")
	resp.Header().Set("Content-Disposition", "attachment; filename="+filename)
	resp.Header().Set("Content-Type", "application/octet-stream")
	resp.WriteHeader(200)
	if _, err := resp.Write(bytes.Bytes()); err != nil {
		Error(resp, r, new(domainerror.InternalError))
	}
}

func unwrapError(err error) view.Error {
	errObj := view.Error{
		Code:     constants.StatusCodeGenericInternalError,
		Message:  err.Error(),
		HttpCode: http.StatusInternalServerError,
	}

	if code, ok := err.(domainerror.Interface); ok {
		if data, ok := code.(interface{ GetData() interface{} }); ok {
			errObj.Data = data.GetData()
		}
		errObj.Code = code.Code()
		errObj.Message = code.GetMessage()
		errObj.HttpCode = code.GetHttpCode()
		return errObj
	}

	unwrapErr := errors.Unwrap(err)
	for unwrapErr != nil {
		if domainErr, ok := unwrapErr.(domainerror.Interface); ok {
			if data, ok := domainErr.(interface{ GetData() interface{} }); ok {
				errObj.Data = data.GetData()
			}
			errObj.Code = domainErr.Code()
			errObj.Message = domainErr.GetMessage()
			errObj.HttpCode = domainErr.GetHttpCode()
			break
		} else {
			unwrapErr = errors.Unwrap(unwrapErr)
		}
	}

	return errObj
}
