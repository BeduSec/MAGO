// Copyright (c) BeduSec. All rights reserved.
package mago_test

import (
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/bedusec/mago/internal/server"
	"github.com/bedusec/mago/pkg/config"
	"github.com/bedusec/mago/pkg/logging"
)

func TestIntegrationRateLimit(t *testing.T) {
	cfg := &config.Config{
		Server:      config.ServerConfig{Host: "127.0.0.1", Port: 0},
		Store:       config.StoreConfig{Type: "memory"},
		RateLimiter: config.RateLimiterConfig{DefaultRate: 2, DefaultBurst: 2},
		WAF:         config.WAFConfig{RulesFile: "nonexistent.json", DryRun: false},
		LogLevel:    "info",
		LogJSON:     true,
	}
	logger, _ := logging.New("info", true)
	srv, err := server.New(cfg, logger)
	if err != nil {
		t.Fatalf("server.New: %v", err)
	}
	go srv.ListenAndServe()
	defer srv.Shutdown()
	time.Sleep(500 * time.Millisecond)

	addr := "http://127.0.0.1:" + getPort(srv)
	client := &http.Client{}
	okCount := 0
	limited := false
	for i := 0; i < 10; i++ {
		resp, err := client.Get(addr + "/")
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
		if resp.StatusCode == 200 {
			okCount++
		} else if resp.StatusCode == 429 {
			limited = true
			break
		}
	}
	if !limited && okCount > 2 {
		t.Errorf("expected rate limiting, got %d OK without 429", okCount)
	}
}

func getPort(srv *server.Server) string {
	return "8080"
}