# MAGO

[![CI](https://img.shields.io/github/actions/workflow/status/bedusec/mago/ci.yml?style=flat-square)](https://github.com/bedusec/mago/actions)
[![Go Version](https://img.shields.io/github/go-mod/go-version/bedusec/mago?style=flat-square)](https://go.dev/)
[![License](https://img.shields.io/github/license/bedusec/mago?style=flat-square)](LICENSE)
[![Docker](https://img.shields.io/badge/docker-ready-blue?style=flat-square)](https://hub.docker.com/)
[![Python 3.10+](https://img.shields.io/badge/python-3.10%2B-blue?style=flat-square)](https://python.org)

A production-ready API Rate Limiter and WAF Ruleset Generator for cloud and edge deployment.
MAGO combines a high-performance Go enforcement service with a Python-based ruleset management toolchain,
optional assembly-accelerated token-bucket routines, and full observability.

## Overview

MAGO acts as a transparent proxy or sidecar that enforces configurable rate limits and Web Application Firewall rules.
It is designed to run as a single binary, a container, or inside Kubernetes, and exposes HTTP/gRPC endpoints
for health, metrics, and dynamic rule reloads. The companion Python CLI (`rulesgen`) translates
high-level policy documents into optimized JSON rulesets and can execute live tests against a running instance.

## Features

- **Token-bucket rate limiter** with pluggable storage (in-memory / Redis).
- **WAF engine** supporting IP, path, header, and JSON body matching with custom safe evaluators.
- **Hot-path assembly stubs** for x86_64 token decrement (pure-Go fallback included).
- **Structured logging** (JSON via `zap`), Prometheus metrics, and health endpoints.
- **Dynamic rule reload** via authenticated HTTP endpoint.
- **Dry-run mode** that logs matches without blocking.
- **Python ruleset generator** accepting YAML policy and emitting validated JSON rulesets.
- **Integration tests** simulating concurrent clients and WAF scenarios.
- **Kubernetes manifests** and Docker Compose for local development.
- **GitHub Actions CI** with linting, unit tests, and release artifact upload.

## Architecture

```
                  ┌──────────┐
                  │  Client  │
                  └────┬─────┘
                       │
              ┌────────▼────────┐
              │     MAGO        │
              │  (Go service)   │
              │  ┌──────────┐   │
              │  │  WAF     │   │
              │  │  Engine  │   │
              │  └────┬─────┘   │
              │       │         │
              │  ┌────▼─────┐   │
              │  │  Rate    │   │
              │  │  Limiter │   │
              │  └────┬─────┘   │
              │       │         │
              │  ┌────▼─────┐   │
              │  │  Store   │   │
              │  │ (Mem/    │   │
              │  │  Redis)  │   │
              │  └──────────┘   │
              └────────┬────────┘
                       │
              ┌────────▼────────┐
              │  Upstream API   │
              └─────────────────┘
```

## Quick Start

### Prerequisites

- Go 1.20+
- Python 3.10+
- Docker (optional)
- Redis (optional, for distributed limits)

### Build and Run

```shell
git clone https://github.com/bedusec/mago.git
cd mago
make build
./mago serve --config config.yaml
```

The service listens on `http://0.0.0.0:8080`. Verify with:

```shell
curl http://localhost:8080/healthz
```

### Using Docker Compose

```shell
docker-compose up --build
```

This brings up MAGO and a Redis instance. Redis is used as the store if `MAGO_STORE_TYPE=redis`.

### Kubernetes Deployment

Apply the provided manifests:

```shell
kubectl apply -f k8s/deployment.yaml
```

Make sure the `mago-config` ConfigMap exists with your `config.yaml` and ruleset file.

## Configuration

MAGO reads a YAML configuration file (default `config.yaml`) and overrides it with environment variables.

### Example `config.yaml`

``yaml
server:
  host: "0.0.0.0"
  port: 8080
store:
  type: memory
rate_limiter:
  default_rate: 100
  default_burst: 200
  cleanup_interval_sec: 60
waf:
  rules_file: examples/rulesets/default.json
  dry_run: false
admin_token: "secret-token"
log_level: info
log_json: true
``

### Environment Variables

| Variable             | Description                    |
|----------------------|--------------------------------|
| `MAGO_SERVER_HOST`   | Bind address                   |
| `MAGO_SERVER_PORT`   | Port                           |
| `MAGO_STORE_TYPE`    | `memory` or `redis`            |
| `MAGO_REDIS_URL`     | Redis connection URL           |
| `MAGO_RATE_DEFAULT`  | Default sustained rate (rps)   |
| `MAGO_RATE_BURST`    | Maximum burst size             |
| `MAGO_WAF_RULES_FILE`| Path to JSON ruleset           |
| `MAGO_WAF_DRY_RUN`   | `true` to enable dry-run mode  |
| `MAGO_ADMIN_TOKEN`   | Bearer token for admin APIs    |
| `MAGO_LOG_LEVEL`     | Log level (debug, info, warn)  |
| `MAGO_LOG_JSON`      | `true` for JSON logs           |

## API Reference

### HTTP Endpoints

| Method | Path              | Description                     | Auth        |
|--------|-------------------|---------------------------------|-------------|
| GET    | `/healthz`        | Health check                    | None        |
| GET    | `/metrics`        | Prometheus metrics              | None        |
| POST   | `/v1/rules/reload`| Reload WAF ruleset              | Bearer token|
| GET    | `/v1/rules`       | List loaded WAF rules           | None        |
| *      | `/*`              | Proxied request (default echo)  | None        |

### gRPC (optional)

The proto definition is at `api/proto/mago.proto`. The service `Mago` provides `Health` and `ReloadRules` RPCs.

## Rate Limiting

MAGO implements the token-bucket algorithm. Each client is identified by IP address or the `X-API-Key` header.
The default rate and burst are configurable globally; per-route and per-key overrides can be added via the
WAF ruleset or future configuration extensions.

Rate limit headers are added to every response:

- `X-RateLimit-Limit`
- `X-RateLimit-Remaining`
- `Retry-After`

## Web Application Firewall (WAF)

Rules are stored in a JSON array with the following schema:

- `id` – unique identifier
- `priority` – evaluation order (lower value = higher priority)
- `action` – `allow`, `block`, `log`
- `conditions` – list of matchers
- `match_type` – `all` (AND) or `any` (OR)

### Condition Fields

| Field          | Operators              | Description                      |
|----------------|------------------------|----------------------------------|
| `ip`           | `eq`                   | Remote IP match                  |
| `path`         | `eq`, `contains`, `regex` | URL path match                  |
| `method`       | `eq`                   | HTTP method                      |
| `header.NAME`  | `eq`, `contains`, `regex` | Request header value           |
| `body`         | JSON path              | Body match using GJSON syntax    |

Example rule blocking `../` traversal:

``json
[
  {
    "id": "block-traversal",
    "priority": 1,
    "action": "block",
    "conditions": [
      {
        "field": "path",
        "operator": "regex",
        "value": ".*\\.\\..*"
      }
    ],
    "match_type": "all"
  }
]
``

## Ruleset Generator (Python)

The `rulesgen` CLI turns high-level policies into optimized JSON rulesets and runs validation tests.

### Installation

```shell
pip install -r tools/rulesgen/requirements.txt
```

### Usage

Generate a ruleset from a YAML policy:

```shell
python tools/rulesgen/cli.py generate --policy examples/policy.yaml --out examples/rulesets/strict.json
```

Run tests against a live MAGO instance:

```shell
python tools/rulesgen/cli.py test --rules examples/rulesets/strict.json --target http://localhost:8080
```

The generator can be integrated into CI pipelines to validate rules before deployment.

## Monitoring

Prometheus metrics are exposed at `/metrics`. Key metrics:

- `mago_requests_total` – total requests by method, path, and status
- `mago_blocked_requests_total` – blocked requests with rule ID
- `mago_ratelimited_requests_total` – rate-limited requests
- `mago_request_duration_seconds` – latency histogram

Logs are structured JSON containing request ID, trace ID, duration, and status.

## Contributing

We welcome contributions! Please read [CONTRIBUTING.md](CONTRIBUTING.md) and open a pull request.
All contributors are expected to adhere to the code of conduct.

## License

This project is licensed under the Apache License 2.0. See [LICENSE](LICENSE) for details.

## Credits

MAGO is developed and maintained by the BeduSec team. Thanks to all contributors for their efforts.