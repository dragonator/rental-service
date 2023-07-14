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

func errorResponse(w http.ResponseWriter, err error) {
	er := &contract.ErrorResponse{Message: err.Error()}
	w.Header().Set(_contentTypeHeaderName, _contentTypeJSON)
	w.Header().Set(_xContentTypeOptions, _noSniff)

	var e *svc.Error
	if errors.As(err, &e) {
		w.WriteHeader(e.StatusCode)
		json.NewEncoder(w).Encode(er)
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(er)
}

func successResponse(w http.ResponseWriter, resp interface{}) {
	w.Header().Set(_contentTypeHeaderName, _contentTypeJSON)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
