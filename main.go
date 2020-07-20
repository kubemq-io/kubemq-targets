package main

import (
	"context"
	"github.com/kubemq-hub/kubemq-targets/binding"
	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/kubemq-hub/kubemq-targets/pkg/logger"

	"os"
	"os/signal"
	"syscall"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var (
	log *logger.Logger
)

func run() error {
	var gracefulShutdown = make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, syscall.SIGTERM)
	signal.Notify(gracefulShutdown, syscall.SIGINT)
	signal.Notify(gracefulShutdown, syscall.SIGQUIT)

	cfg, err := config.Load()
	if err != nil {
		return err
	}
	err = cfg.Validate()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	bindingsService := binding.New()
	err = bindingsService.Start(ctx, cfg)
	if err != nil {
		return err
	}
	<-gracefulShutdown
	bindingsService.Stop()
	return nil
}
func main() {
	log = logger.NewLogger("main")
	log.Infof("starting kubemq targets connector version: %s, commit: %s, date %s", version, commit, date)
	if err := run(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}
