package router

import (
	"net/http"

	"book-management/internal/api/middleware"
	"book-management/internal/api/v1/book/handler"

	"github.com/gorilla/mux"
)

func NewRouter(bookHandler *handler.Handler) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) { w.Write([]byte("OK")) }).Methods(http.MethodGet)

	v1 := r.PathPrefix("/v1").Subrouter()

	// public
	v1.HandleFunc("/books", bookHandler.List).Methods(http.MethodGet)
	v1.HandleFunc("/books/{id}", bookHandler.Get).Methods(http.MethodGet)

	// secured
	secured := v1.NewRoute().Subrouter()
	secured.Use(middleware.AuthJWT)
	secured.HandleFunc("/books", bookHandler.Create).Methods(http.MethodPost)
	secured.HandleFunc("/books/{id}", bookHandler.Update).Methods(http.MethodPut)
	secured.HandleFunc("/books/{id}", bookHandler.Delete).Methods(http.MethodDelete)

	return r
}
