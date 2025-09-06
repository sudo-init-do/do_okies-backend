# syntax=docker/dockerfile:1.7

########################
# Build stage
########################
ARG GO_VERSION=1.22
FROM golang:${GO_VERSION} AS build

WORKDIR /src

# Cache Go modules
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Copy source
COPY . .

# Build (static, trimmed)
ARG TARGETOS=linux
ARG TARGETARCH=amd64
ENV CGO_ENABLED=0
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -trimpath -ldflags="-s -w" -o /out/api ./apps/api

########################
# Runtime stage
########################
FROM gcr.io/distroless/base-debian12:nonroot

WORKDIR /srv
COPY --from=build /out/api /srv/api

# Documented port (actual port controlled by env/compose)
EXPOSE 8081

# Non-root user provided by distroless:nonroot
ENTRYPOINT ["/srv/api"]
