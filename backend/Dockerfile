# syntax=docker/dockerfile:1
#
# Multi-stage build producing ONE service binary, selected by the SERVICE
# build-arg (gateway | competition | user | draw | media). Each compose service
# builds this same Dockerfile with `args: { SERVICE: <name> }`.
#
# Base images use docker.io only (the original gcr.aban.io / gcr.io bases are
# not reachable from CI). CGO is disabled so the binary runs on a bare alpine.

ARG GO_IMAGE=golang:1.25-alpine
ARG RUNTIME_IMAGE=alpine:3.20

# --- build stage ------------------------------------------------------------
FROM ${GO_IMAGE} AS builder

# Which cmd/<SERVICE> to compile. Required — the build fails fast if unset.
ARG SERVICE

WORKDIR /src
ENV CGO_ENABLED=0 GOOS=linux

COPY . .

RUN test -n "${SERVICE}" || (echo "ERROR: SERVICE build-arg is required (gateway|competition|user|draw|media)" >&2; exit 1)

# Prefer a vendored build (hermetic, no network — required where the Go module
# proxy is unreachable). Fall back to downloading modules when vendor/ is absent.
RUN if [ -d vendor ]; then \
        go build -trimpath -ldflags="-s -w" -mod=vendor -o /out/service ./cmd/${SERVICE}; \
    else \
        go mod download && go build -trimpath -ldflags="-s -w" -o /out/service ./cmd/${SERVICE}; \
    fi

# --- runtime stage ----------------------------------------------------------
FROM ${RUNTIME_IMAGE} AS runtime

RUN apk add --no-cache ca-certificates \
    && adduser -D -H -u 10001 app

COPY --from=builder /out/service /app/bin/service
# Baked-in defaults; every value is overridable via APP_* env (see app/config.go).
COPY ./config.example.yaml /app/config.yaml

USER app
EXPOSE 8080

ENTRYPOINT ["/app/bin/service", "--config", "/app/config.yaml"]
