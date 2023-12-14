package rest

import (
	"github.com/alichtenthaler/ps-tag-onboarding-go/api/internal/middleware"
	"github.com/gorilla/mux"
	"net/http"
)

// NewRouter creates a new router for the API
func NewRouter(createUserHandler CreateUserHandler, findUserHandler FindUserHandler) *mux.Router {
	router := mux.NewRouter()
	router.Use(middleware.Logger)

	router.HandleFunc("/user/save", createUserHandler.CreateUser).Methods(http.MethodPost)
	router.HandleFunc("/user/find/{userId}", findUserHandler.FindUser).Methods(http.MethodGet)

	return router
}
