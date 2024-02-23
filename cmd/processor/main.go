package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/broganross/salad-interview/internal/config"
	"github.com/broganross/salad-interview/internal/domain"
	"github.com/broganross/salad-interview/internal/handler"
	"github.com/broganross/salad-interview/internal/handler/middleware"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Info().Msg("setting up")
	conf := config.Config{}
	if err := envconfig.Process("salad", &conf); err != nil {
		log.Error().Err(err).Msg("parsing environment variables")
		os.Exit(1)
	}
	zerolog.SetGlobalLevel(zerolog.Level(conf.Log.Level))

	// set up signal handling
	ctx := context.Background()
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	// build the listener, and handlers
	listener := handler.NewListener(ctx, &conf)
	logHandler := middleware.Log{}
	listener.AddHandler(&logHandler)
	domain := domain.Example{}
	handler := handler.PlaneStatusHandler{
		Domain: &domain,
	}
	listener.AddHandler(&handler)

	// start listening for messages from the server
	go func() {
		if err := listener.Listen(); err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}
			log.Error().Err(err).Msg("listening")
			os.Exit(2)
		}
	}()

	<-ctx.Done()
	stop()
	log.Info().Msg("shutting down")
}
