.PHONY: generate
# generate
generate:
	go mod tidy
	go get github.com/google/wire/cmd/wire@latest
	go generate ./...
	go mod tidy

.PHONY: all
# generate all
all:
	make generate;

.DEFAULT_GOAL := generate

all_tests:
	go test -v ./internal/biz/... ./internal/service/handler/... -bench=. -cover -coverprofile=coverage.out -benchmem -cpu=1 -timeout=500ms -json


bench_tests:
	go test -v ./internal/service/handler/... ./internal/biz/... -bench=. -benchmem -cpu=1,2,3,4 -timeout=500ms

unit_tests:
	go test -v ./internal/service/handler/... ./internal/biz/...

coverage_tests:
	go test -v ./internal/service/handler/... ./internal/biz/... -cover -coverprofile=coverage.out



devtools:
	@echo "Installing devtools"
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install mvdan.cc/gofumpt@latest
	go install go.uber.org/mock/mockgen@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/daixiang0/gci@v0.11.2



swagger:
	swag fmt
	swag i  -d internal/service,internal/media,internal/competition,internal/user,internal/draw -g server.go  -pd



check:
	golangci-lint run  \
		--build-tags "${BUILD_TAG}" \
		--timeout=20m0s \
		--enable=gofmt \
		--enable=unconvert \
		--enable=unparam \
		--enable=asciicheck \
		--enable=misspell \
		--enable=revive \
		--enable=decorder \
		--enable=reassign \
		--enable=usestdlibvars \
		--enable=nilerr \
		--enable=gosec \
		--enable=exportloopref \
		--enable=whitespace \
		--enable=gocyclo \
		--enable=nestif \
		--enable=gochecknoinits \
		--enable=gocognit \
		--enable=funlen \
		--enable=forbidigo \
		--enable=godox \
		--enable=gocritic \
		--enable=gci \
		--enable=lll \
		--config=issues.exclude.yaml




build:
	docker build -t buildf .


.PHONY: e2e
# Run the end-to-end suite against a LIVE compose stack through the gateway.
# Bring the stack up first: docker compose up -d --build && make seed
E2E_BASE_URL ?= http://localhost:8080
e2e:
	E2E_BASE_URL=$(E2E_BASE_URL) go test -tags e2e -count=1 -v ./test/e2e/...


.PHONY: seed
# Seed Postgres + MinIO with the canonical sample inventory. Requires
# APP_DATASOURCE_POSTGRES_DSN and the SEED_MINIO_* env vars (see cmd/seed).
# Idempotent: IDs are deterministic UUIDv5, so re-running upserts, never dupes.
seed:
	go run ./cmd/seed

