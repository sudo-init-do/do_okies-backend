# -------- build stage --------
FROM --platform=$BUILDPLATFORM golang:1.24.4 AS build
WORKDIR /src

# Ensure toolchain auto-updates to match go.mod
ENV GOTOOLCHAIN=auto

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build
COPY . .
ARG TARGETOS=linux
ARG TARGETARCH=amd64
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -trimpath -ldflags="-s -w" -o /out/api ./apps/api

# -------- run stage --------
FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=build /out/api /app/api
EXPOSE 8080                
USER 65532:65532          
ENTRYPOINT ["/app/api"]
