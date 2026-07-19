// Copyright (c) BeduSec. All rights reserved.
package enforcer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/bedusec/mago/internal/limiter"
	"github.com/bedusec/mago/internal/store"
	"github.com/bedusec/mago/internal/waf"
	"github.com/bedusec/mago/pkg/config"
	"github.com/bedusec/mago/pkg/metrics"
	"go.uber.org/zap"
)

type Enforcer struct {
	cfg        *config.Config
	store      store.Store
	logger     *zap.SugaredLogger
	rateLimit  *limiter.RateLimiter
	wafEngine  *waf.Engine
	wafRules   []waf.Rule
}

func New(cfg *config.Config, st store.Store, logger *zap.SugaredLogger) *Enforcer {
	return &Enforcer{
		cfg:       cfg,
		store:     st,
		logger:    logger,
		rateLimit: limiter.NewRateLimiter(st, cfg.RateLimiter.DefaultRate, cfg.RateLimiter.DefaultBurst),
		wafEngine: waf.NewEngine(cfg.WAF.DryRun, logger),
	}
}

func (e *Enforcer) LoadWAFRules(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("reading rules file: %w", err)
	}
	var rules []waf.Rule
	if err := json.Unmarshal(data, &rules); err != nil {
		return fmt.Errorf("unmarshalling rules: %w", err)
	}
	if err := e.wafEngine.Load(rules); err != nil {
		return err
	}
	e.wafRules = rules
	e.logger.Infow("WAF rules loaded", "count", len(rules))
	return nil
}

func (e *Enforcer) ReloadRulesHandler(w http.ResponseWriter, r *http.Request) {
	path := e.cfg.WAF.RulesFile
	if err := e.LoadWAFRules(path); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("rules reloaded"))
}

func (e *Enforcer) ListRulesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(e.wafRules)
}

func (e *Enforcer) WAFMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		action, ruleID := e.wafEngine.Match(r)
		switch action {
		case waf.ActionBlock:
			metrics.BlockedRequests.WithLabelValues(ruleID).Inc()
			http.Error(w, "blocked by WAF", http.StatusForbidden)
			return
		case waf.ActionLog:
			e.logger.Infow("WAF match (log)", "rule_id", ruleID, "path", r.URL.Path)
		}
		next.ServeHTTP(w, r)
	})
}

func (e *Enforcer) RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.RemoteAddr
		if apiKey := r.Header.Get("X-API-Key"); apiKey != "" {
			key = apiKey
		}
		allowed, remaining, resetAfter, err := e.rateLimit.Allow(key)
		if err != nil {
			e.logger.Errorw("rate limiter error", "error", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("X-RateLimit-Limit", strconv.Itoa(e.cfg.RateLimiter.DefaultBurst))
		w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(remaining))
		w.Header().Set("Retry-After", strconv.Itoa(int(resetAfter.Seconds())))
		if !allowed {
			metrics.RateLimitedRequests.Inc()
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}