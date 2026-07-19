// Copyright (c) BeduSec. All rights reserved.
package mago_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bedusec/mago/internal/enforcer"
	"github.com/bedusec/mago/internal/store"
	"github.com/bedusec/mago/internal/waf"
	"github.com/bedusec/mago/pkg/config"
	"github.com/bedusec/mago/pkg/logging"
)

func TestWAFBlock(t *testing.T) {
	logger, _ := logging.New("info", true)
	cfg := &config.Config{
		WAF: config.WAFConfig{DryRun: false},
		RateLimiter: config.RateLimiterConfig{DefaultRate: 100, DefaultBurst: 200},
		Store:       config.StoreConfig{Type: "memory"},
	}
	st, _ := store.New(cfg.Store)
	enf := enforcer.New(cfg, st, logger)
	rules := []waf.Rule{
		{
			ID:       "test1",
			Priority: 1,
			Action:   "block",
			Conditions: []waf.Condition{
				{Field: "path", Operator: "regex", Value: "/admin"},
			},
			MatchType: "all",
		},
	}
	data, _ := json.Marshal(rules)
	tmpFile, _ := os.CreateTemp("", "rules-*.json")
	tmpFile.Write(data)
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())
	enf.LoadWAFRules(tmpFile.Name())

	handler := enf.WAFMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	req := httptest.NewRequest("GET", "/admin", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != 403 {
		t.Errorf("expected 403, got %d", w.Code)
	}
}