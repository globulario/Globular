# ---- Project settings ----
# We'll auto-detect PKG below; BIN defaults to "globular" (lowercase for Docker tags)
BIN ?= globular
IMG := ghcr.io/globulario/$(BIN):$(shell git rev-parse --short HEAD)

# --- Auto-detect entrypoint (root or cmd/*) ---
# Prefer root main.go, else cmd/Globular, else cmd/globular.
ifeq ("$(wildcard main.go)","")
  ifneq ("$(wildcard cmd/Globular/main.go)","")
    PKG ?= ./cmd/Globular
    BIN ?= Globular
  else ifneq ("$(wildcard cmd/globular/main.go)","")
    PKG ?= ./cmd/globular
    BIN ?= globular
  else
    # Fallback (still allows override: make build PKG=./some/other/cmd)
    PKG ?= .
  endif
else
  PKG ?= .
endif

.PHONY: tidy fmt vet lint test bench run dev-up dev-down build image vuln sec sbom clean print-vars

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

build:
	mkdir -p ./.bin
	CGO_ENABLED=$(CGO) go build -trimpath -ldflags "-s -w" -o ./.bin/$(BIN) $(PKG)

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

