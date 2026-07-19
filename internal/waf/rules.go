// Copyright (c) BeduSec. All rights reserved.
package waf

import (
	"net"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/tidwall/gjson"
	"go.uber.org/zap"
)

const (
	ActionBlock = "block"
	ActionLog   = "log"
	ActionAllow = "allow"
)

type Condition struct {
	Field    string      `json:"field"`
	Operator string      `json:"operator"`
	Value    interface{} `json:"value"`
}

type Rule struct {
	ID        string      `json:"id"`
	Priority  int         `json:"priority"`
	Action    string      `json:"action"`
	Conditions []Condition `json:"conditions"`
	MatchType string      `json:"match_type"`
}

type compiledRule struct {
	rule       Rule
	matchers   []matcher
}

type matcher func(r *http.Request) bool

type Engine struct {
	mu        sync.RWMutex
	rules     []compiledRule
	dryRun    bool
	logger    *zap.SugaredLogger
}

func NewEngine(dryRun bool, logger *zap.SugaredLogger) *Engine {
	return &Engine{dryRun: dryRun, logger: logger}
}

func (e *Engine) Load(rules []Rule) error {
	compiled := make([]compiledRule, len(rules))
	for i, rule := range rules {
		matchers := make([]matcher, len(rule.Conditions))
		for j, cond := range rule.Conditions {
			m, err := compileCondition(cond)
			if err != nil {
				return err
			}
			matchers[j] = m
		}
		compiled[i] = compiledRule{rule: rule, matchers: matchers}
	}
	sort.Slice(compiled, func(i, j int) bool {
		return compiled[i].rule.Priority < compiled[j].rule.Priority
	})
	e.mu.Lock()
	e.rules = compiled
	e.mu.Unlock()
	return nil
}

func (e *Engine) Match(r *http.Request) (string, string) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	for _, cr := range e.rules {
		matched := true
		if cr.rule.MatchType == "all" {
			for _, m := range cr.matchers {
				if !m(r) {
					matched = false
					break
				}
			}
		} else {
			matched = false
			for _, m := range cr.matchers {
				if m(r) {
					matched = true
					break
				}
			}
		}
		if matched {
			if e.dryRun {
				e.logger.Infow("WAF rule matched (dry-run)", "rule", cr.rule.ID)
				return ActionLog, cr.rule.ID
			}
			return cr.rule.Action, cr.rule.ID
		}
	}
	return ActionAllow, ""
}

func compileCondition(cond Condition) (matcher, error) {
	switch cond.Field {
	case "ip":
		return ipMatcher(cond), nil
	case "path":
		return pathMatcher(cond)
	case "method":
		return methodMatcher(cond), nil
	case "header":
		parts := strings.SplitN(cond.Field, ".", 2)
		if len(parts) == 2 {
			return headerMatcher(parts[1], cond), nil
		}
		return noMatch, nil
	case "body":
		return bodyMatcher(cond)
	default:
		return noMatch, nil
	}
}

func noMatch(r *http.Request) bool { return false }

func ipMatcher(cond Condition) matcher {
	ipStr, _ := cond.Value.(string)
	ip := net.ParseIP(ipStr)
	return func(r *http.Request) bool {
		remoteIP, _, _ := net.SplitHostPort(r.RemoteAddr)
		if remoteIP == "" {
			remoteIP = r.RemoteAddr
		}
		return net.ParseIP(remoteIP).Equal(ip)
	}
}

func pathMatcher(cond Condition) (matcher, error) {
	pattern, ok := cond.Value.(string)
	if !ok {
		return noMatch, nil
	}
	switch cond.Operator {
	case "regex":
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, err
		}
		return func(r *http.Request) bool {
			return re.MatchString(r.URL.Path)
		}, nil
	case "eq":
		return func(r *http.Request) bool {
			return r.URL.Path == pattern
		}, nil
	case "contains":
		return func(r *http.Request) bool {
			return strings.Contains(r.URL.Path, pattern)
		}, nil
	default:
		return noMatch, nil
	}
}

func methodMatcher(cond Condition) matcher {
	method, _ := cond.Value.(string)
	return func(r *http.Request) bool {
		return r.Method == method
	}
}

func headerMatcher(headerName string, cond Condition) matcher {
	value, _ := cond.Value.(string)
	switch cond.Operator {
	case "eq":
		return func(r *http.Request) bool {
			return r.Header.Get(headerName) == value
		}
	case "contains":
		return func(r *http.Request) bool {
			return strings.Contains(r.Header.Get(headerName), value)
		}
	case "regex":
		re, _ := regexp.Compile(value)
		return func(r *http.Request) bool {
			return re.MatchString(r.Header.Get(headerName))
		}
	default:
		return noMatch
	}
}

func bodyMatcher(cond Condition) (matcher, error) {
	path, _ := cond.Value.(string)
	return func(r *http.Request) bool {
		bodyBytes := make([]byte, r.ContentLength)
		r.Body.Read(bodyBytes)
		defer r.Body.Close()
		result := gjson.GetBytes(bodyBytes, path)
		return result.Exists()
	}, nil
}