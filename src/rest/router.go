package rest

import (
	"github.com/alichtenthaler/ps-tag-onboarding-go/api/src/middleware"
	"github.com/alichtenthaler/ps-tag-onboarding-go/api/src/user"
	"github.com/gorilla/mux"
	"net/http"
)

func newRouter(userProcessor *user.Service) *mux.Router {
	router := mux.NewRouter()
	router.Use(middleware.Logger)

	router.HandleFunc("/user/find/{userId}", userProcessor.FindUserById).Methods(http.MethodGet)
	router.HandleFunc("/user/save", userProcessor.CreateUser).Methods(http.MethodPost)

	return router
}
