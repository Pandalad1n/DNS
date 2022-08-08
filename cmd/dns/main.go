package main

import (
	"context"
	"errors"
	"github.com/Pandalad1n/DNS/cmd/dns/handler"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	config := Config{}.WithEnv()
	if err := config.Validate(); err != nil {
		log.Fatal().Err(err).Msg("Invalid config.")
	}
	level, err := zerolog.ParseLevel(config.LogLevel)
	if err != nil {
		log.Fatal().Msgf("Unknown log level -- %s.", config.LogLevel)
	}
	ctx := context.Background()
	ctx = log.Logger.Level(level).WithContext(ctx)
	webServer := &http.Server{
		Addr:         config.ListenWebAddress,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      handler.NewHandler(*config.SectorID),
		BaseContext: func(net.Listener) context.Context {
			return ctx
		},
	}

	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Ctx(ctx).Info().Msgf("Web listening on %s.", config.ListenWebAddress)
		err = webServer.ListenAndServe()
		if err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return
			}
			log.Fatal().Err(err).Msg("Web server failed.")
		}
	}()
	<-ctx.Done()
	log.Ctx(ctx).Info().Msg("Shutting down.")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	_ = webServer.Shutdown(shutdownCtx)

	log.Ctx(ctx).Info().Msg("Shutdown complete.")
}
