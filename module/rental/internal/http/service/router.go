package service

import (
	"net/http"

	"github.com/go-chi/chi"
)

type RentalHandler interface {
	GetRentalByID(method, path string) func(w http.ResponseWriter, r *http.Request)
}

// NewRouter is a construction function for router that handles operations for rentals.
func NewRouter(rh RentalHandler) http.Handler {
	router := chi.NewRouter()

	api := []struct {
		MethodFunc func(pattern string, handlerFn http.HandlerFunc)
		Method     string
		Path       string
		HandleFunc func(string, string) func(w http.ResponseWriter, r *http.Request)
	}{
		{router.Get, "GET", "/rental/{id}", rh.GetRentalByID},
	}

	for _, endpoint := range api {
		endpoint.MethodFunc(endpoint.Path, endpoint.HandleFunc(endpoint.Method, endpoint.Path))
	}

	return router
}
