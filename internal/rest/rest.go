package rest

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

// Rest represents the response handler
type Rest struct {
	server *http.Server
}

// New creates the response handler
func New(router *mux.Router, options ...func(*Rest)) *Rest {

	server := &Rest{
		&http.Server{
			Handler: router,
		},
	}

	for _, o := range options {
		o(server)
	}

	return server
}

func WithPort(port int) func(*Rest) {
	return func(s *Rest) {
		s.server.Addr = fmt.Sprintf(":%d", port)
	}
}

// Start starts the response server
func (rest *Rest) Start() {

	port := "80"
	if rest.server.Addr != "" {
		port = rest.server.Addr
	}

	log.Info().Msgf("Server listening on %s", port)
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

// Shutdown shuts down the response server gracefully
func (rest *Rest) shutdown(ctx context.Context) error {

	if err := rest.server.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}
