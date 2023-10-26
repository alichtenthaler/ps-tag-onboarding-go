package rest

import (
	"github.com/alichtenthaler/ps-tag-onboarding-go/api/src/middleware"
	"github.com/alichtenthaler/ps-tag-onboarding-go/api/src/user"
	"github.com/gorilla/mux"
	"net/http"
)

func newRouter(userProcessor *user.Service) *mux.Router {
	router := mux.NewRouter()

	// ENDPOINTS
	// User
	router.HandleFunc("/find/{userId}", middleware.Logger(userProcessor.FindUserById)).Methods(http.MethodGet)
	router.HandleFunc("/save", middleware.Logger(userProcessor.CreateUser)).Methods(http.MethodPost)

	return router
}
