package rest

import (
	"context"
	"errors"
	"fmt"
	"github.com/alichtenthaler/ps-tag-onboarding-go/api/internal/user"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Rest represents the http handler
type Rest struct {
	server *http.Server
}

// New creates the http handler
func New(port int, userService *user.Service) *Rest {

	return &Rest{
		&http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: newRouter(userService),
		},
	}
}

// Start starts the http server
func (rest *Rest) Start() {

	log.Info().Msgf("Server listening on %s", rest.server.Addr)
	go func() {

		if err := rest.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Msgf("error starting server: %s", err)
		}
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-done
	log.Info().Msg("Shutting down HTTP server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := rest.shutdown(ctx); err != nil {
		log.Fatal().Msg("Failed to shutdown HTTP server")
	}

	log.Info().Msg("Server was shutdown properly")
}

// Shutdown shuts down the http server gracefully
func (rest *Rest) shutdown(ctx context.Context) error {

	if err := rest.server.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}
