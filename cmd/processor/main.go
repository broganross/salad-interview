package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/broganross/salad-interview/internal/config"
	"github.com/broganross/salad-interview/internal/listener"
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

	listener := listener.New(ctx, &conf)
	if err := listener.Start(); err != nil {
		log.Error().Err(err).Msg("starting up listener")
		os.Exit(2)
	}

	select {
	// case <- neverReady:
	// 	fmt.Println("ready")
	case <-ctx.Done():
		log.Info().Msg("context canceled")
		stop()
	}
	log.Info().Msg("shutting down")
}
