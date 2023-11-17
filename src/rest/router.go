package rest

import (
	"github.com/alichtenthaler/ps-tag-onboarding-go/api/src/adapter/in/web"
	"github.com/alichtenthaler/ps-tag-onboarding-go/api/src/middleware"
	"github.com/gorilla/mux"
	"net/http"
)

func newRouter(getUserHandler *web.GetUserHandler, createUserHandler *web.CreateUserHandler) *mux.Router {
	router := mux.NewRouter()
	router.Use(middleware.Logger)

	router.HandleFunc("/user/find/{userId}", getUserHandler.HandleGetUser).Methods(http.MethodGet)
	router.HandleFunc("/user/save", createUserHandler.HandleCreteUser).Methods(http.MethodPost)

	return router
}
