# ---- Project settings ----
BIN ?= globular
IMG := ghcr.io/globulario/$(BIN):$(shell git rev-parse --short HEAD)
PKG ?= .
LEGACY_CMD := ./cmd/globular
LEGACY_BIN ?= globular-dev

.PHONY: tidy fmt vet lint test bench run dev-up dev-down build image vuln sec sbom clean print-vars build-legacy build-gateway build-xds build-all check-legacy-guard

tidy:
	go mod tidy

fmt:
	go fmt ./...

vet:
	go vet ./...

lint:
	golangci-lint run ./...

test:
	go test ./... -race -count=1

bench:
	go test ./... -run=^$$ -bench=. -benchmem


# allow override: make build CGO=1
CGO ?= 1   # set 1 by default since you use chai2010/webp

build: build-all

build-gateway: check-gateway-no-exec
	mkdir -p ./.bin
	CGO_ENABLED=$(CGO) go build -trimpath -ldflags "-s -w" -o ./.bin/gateway ./cmd/globular-gateway

build-xds:
	mkdir -p ./.bin
	CGO_ENABLED=$(CGO) go build -trimpath -ldflags "-s -w" -o ./.bin/xds ./cmd/globular-xds

build-all: check-legacy-guard build-gateway build-xds

check-legacy-guard:
	@pattern=$$(printf "%s($$|[^-])" "$(LEGACY_CMD)"); \
	count=$$(rg -c "$$pattern" Makefile 2>/dev/null || true); \
	if [ "$$count" -ne 1 ]; then \
		echo "expected exactly one legacy reference to $(LEGACY_CMD) in the Makefile, found $$count"; \
		exit 1; \
	fi

check-gateway-no-exec:
	@if rg -n "os/exec" cmd/globular-gateway internal/gateway >/dev/null 2>&1; then \
		echo "gateway must not import os/exec"; \
		exit 1; \
	else \
		echo "gateway exec check: OK"; \
	fi

run:
	GOLOG_LOG_LEVEL=info CGO_ENABLED=$(CGO) go run $(PKG)

image:
	docker build -t $(IMG) .

vuln:
	govulncheck ./...

sec:
	gosec ./...

sbom:
	syft packages . -o cyclonedx-json > sbom.json

clean:
	rm -rf ./.bin sbom.json

print-vars:
	@echo "BIN=$(BIN)"
	@echo "PKG=$(PKG)"
