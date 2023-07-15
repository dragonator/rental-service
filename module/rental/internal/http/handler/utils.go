package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/dragonator/rental-service/module/rental/internal/http/contract"
	"github.com/dragonator/rental-service/module/rental/internal/http/service/svc"
)

const (
	_contentTypeHeaderName = "Content-Type"
	_contentTypeJSON       = "application/json"
	_xContentTypeOptions   = "X-Content-Type-Options"
	_noSniff               = "nosniff"
)

type multiError interface {
	Unwrap() []error
}

func errorResponse(w http.ResponseWriter, err error) {
	w.Header().Set(_contentTypeHeaderName, _contentTypeJSON)
	w.Header().Set(_xContentTypeOptions, _noSniff)

	er := toErrorResponse(err)

	var e *svc.Error
	if errors.As(err, &e) {
		w.WriteHeader(e.StatusCode)
		json.NewEncoder(w).Encode(er)
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(er)
}

func toErrorResponse(err error) *contract.ErrorResponse {
	var er *contract.ErrorResponse

	switch e := err.(type) {
	case multiError:
		errs := e.Unwrap()

		er = &contract.ErrorResponse{
			Errors: make([]string, 0, len(errs)),
		}

		for _, err := range errs {
			er.Errors = append(er.Errors, err.Error())
		}

	default:
		er = &contract.ErrorResponse{Error: e.Error()}
	}

	return er
}

func successResponse(w http.ResponseWriter, resp interface{}) {
	w.Header().Set(_contentTypeHeaderName, _contentTypeJSON)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
