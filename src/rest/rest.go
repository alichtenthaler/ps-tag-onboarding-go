package rest

import (
	"context"
	"errors"
	"fmt"
	"github.com/alichtenthaler/ps-tag-onboarding-go/api/src/user"
	"github.com/rs/zerolog/log"
	"net/http"
)

// Rest represents the http handler
type Rest struct {
	server *http.Server
}

// New creates the http handler
func New(port int, userService *user.Service) *Rest {

	return &Rest{
		&http.Server{
			Addr: fmt.Sprintf(":%d", port),
			Handler: newRouter(userService),
		},
	}
}

// Start starts the http server
func (rest *Rest) Start() {

	log.Info().Msgf("Server listening on %s\n", rest.server.Addr)
	if err := rest.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal().Msgf("error starting server: %s\n", err)
	}
}

// Shutdown shuts down the http server gracefully
func (rest *Rest) Shutdown(ctx context.Context) error {

	if err := rest.server.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}
