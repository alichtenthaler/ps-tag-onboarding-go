package rest

import (
	"github.com/alichtenthaler/ps-tag-onboarding-go/api/src/adapter/in/web"
	"github.com/alichtenthaler/ps-tag-onboarding-go/api/src/middleware"
	"github.com/gorilla/mux"
	"net/http"
)

func newRouter(getUserHandler *web.GetUserHandler, createUserHandler *web.CreateUserHandler) *mux.Router {
	router := mux.NewRouter()

	// ENDPOINTS
	// User
	router.HandleFunc("/find/{userId}", middleware.Logger(getUserHandler.HandleGetUser)).Methods(http.MethodGet)
	router.HandleFunc("/save", middleware.Logger(createUserHandler.HandleCreteUser)).Methods(http.MethodPost)

	return router
}
