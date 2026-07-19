// Copyright (c) BeduSec. All rights reserved.
package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bedusec/mago/internal/server"
	"github.com/bedusec/mago/pkg/config"
	"github.com/bedusec/mago/pkg/logging"
)

func main() {
	cfgPath := flag.String("config", "config.yaml", "path to configuration file")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	logger, err := logging.New(cfg.LogLevel, cfg.LogJSON)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	srv, err := server.New(cfg, logger)
	if err != nil {
		logger.Fatalw("failed to create server", "error", err)
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-shutdown
		logger.Infow("received signal, shutting down", "signal", sig)
		srv.Shutdown()
	}()

	if err := srv.ListenAndServe(); err != nil {
		logger.Fatalw("server exited with error", "error", err)
	}
}