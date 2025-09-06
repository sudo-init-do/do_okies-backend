# -------- build stage --------
# Use a Go version that matches your go.mod / toolchain requirement
FROM --platform=$BUILDPLATFORM golang:1.24.4 AS build
WORKDIR /src

# Make the container auto-download matching minor toolchains if needed
ENV GOTOOLCHAIN=auto

# Leverage cache for deps
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest and build
COPY . .
ARG TARGETOS=linux
ARG TARGETARCH=amd64
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -trimpath -ldflags="-s -w" -o /out/api ./apps/api

# -------- run stage --------
FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=build /out/api /app/api
EXPOSE 8081
USER 65532:65532
ENTRYPOINT ["/app/api"]
