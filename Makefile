SVC := globular
IMG := ghcr.io/globulario/$(SVC):$(shell git rev-parse --short HEAD)

.PHONY: tidy lint test bench run dev-up dev-down build image vuln sec sbom

tidy:
	go mod tidy

lint:
	golangci-lint run ./...

test:
	go test ./... -race -count=1

bench:
	go test ./... -run=^$$ -bench=. -benchmem

run:
	GOLOG_LOG_LEVEL=info go run ./cmd/$(SVC)

dev-up:
	docker compose up -d --build

dev-down:
	docker compose down -v

build:
	CGO_ENABLED=0 go build -ldflags "-s -w" -o ./.bin/$(SVC) ./cmd/$(SVC)

image:
	docker build -t $(IMG) .

vuln:
	govulncheck ./...

sec:
	gosec ./...

sbom:
	syft packages . -o cyclonedx-json > sbom.json

