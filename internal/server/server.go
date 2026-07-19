// Copyright (c) BeduSec. All rights reserved.
package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/bedusec/mago/internal/enforcer"
	"github.com/bedusec/mago/internal/store"
	"github.com/bedusec/mago/pkg/config"
	"github.com/bedusec/mago/pkg/logging"
	"github.com/bedusec/mago/pkg/metrics"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

type Server struct {
	httpServer *http.Server
	enforcer   *enforcer.Enforcer
	logger     *zap.SugaredLogger
	cfg        *config.Config
}

func New(cfg *config.Config, logger *zap.SugaredLogger) (*Server, error) {
	st, err := store.New(cfg.Store)
	if err != nil {
		return nil, fmt.Errorf("store init: %w", err)
	}
	enf := enforcer.New(cfg, st, logger)

	if err := enf.LoadWAFRules(cfg.WAF.RulesFile); err != nil {
		logger.Warnw("failed to load initial WAF rules", "error", err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/healthz", healthHandler).Methods("GET")
	router.Handle("/metrics", promhttp.Handler())
	router.HandleFunc("/v1/rules/reload", enf.ReloadRulesHandler).Methods("POST")
	router.HandleFunc("/v1/rules", enf.ListRulesHandler).Methods("GET")

	router.Use(loggingMiddleware(logger))
	router.Use(authMiddleware(cfg.AdminToken, "/v1/rules/reload"))
	router.Use(enf.WAFMiddleware)
	router.Use(enf.RateLimitMiddleware)
	router.Use(metricsMiddleware)

	router.PathPrefix("/").HandlerFunc(defaultHandler)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	return &Server{
		httpServer: &http.Server{
			Addr:    addr,
			Handler: router,
		},
		enforcer: enf,
		logger:   logger,
		cfg:      cfg,
	}, nil
}

func (s *Server) ListenAndServe() error {
	s.logger.Infow("starting server", "addr", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	s.httpServer.Shutdown(ctx)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("MAGO"))
}

func loggingMiddleware(logger *zap.SugaredLogger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			logger.Infow("request",
				"method", r.Method,
				"path", r.URL.Path,
				"remote", r.RemoteAddr,
				"duration", time.Since(start).String(),
			)
		})
	}
}

func authMiddleware(token string, protectedPath string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == protectedPath {
				if token == "" {
					http.Error(w, "unauthorized", http.StatusUnauthorized)
					return
				}
				if r.Header.Get("Authorization") != "Bearer "+token {
					http.Error(w, "unauthorized", http.StatusUnauthorized)
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

func metricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timer := prometheus.NewTimer(metrics.RequestLatency)
		defer timer.ObserveDuration()
		sw := &statusWriter{ResponseWriter: w, status: 200}
		next.ServeHTTP(sw, r)
		metrics.RequestsTotal.WithLabelValues(r.Method, r.URL.Path, http.StatusText(sw.status)).Inc()
	})
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (s *statusWriter) WriteHeader(code int) {
	s.status = code
	s.ResponseWriter.WriteHeader(code)
}