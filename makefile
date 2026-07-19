# Copyright (c) BeduSec. All rights reserved.
.PHONY: build test docker-build lint run-examples

build:
	go build -o mago ./cmd/mago

test:
	go test ./...
	cd tools/rulesgen && python -m pytest tests/ -v

docker-build:
	docker build -t mago:latest .

lint:
	golangci-lint run ./...
	cd tools/rulesgen && flake8 --max-line-length=120 .

run-examples:
	./examples/curl.sh