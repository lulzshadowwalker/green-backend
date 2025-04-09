# Pin specific version for stability
# Use separate stage for building image
# Use debian for easier build utilities
FROM golang:1.23.8-bullseye AS build-base

WORKDIR /app 

# Copy only files required to install dependencies (better layer caching)
COPY go.mod go.sum ./

# Use cache mount to speed up install of existing dependencies
RUN --mount=type=cache,target=/go/pkg/mod \
  --mount=type=cache,target=/root/.cache/go-build \
  go mod download

FROM build-base AS dev

# Install air for hot reload & delve for debugging
RUN go install github.com/cosmtrek/air@latest && \
  go install github.com/go-delve/delve/cmd/dlv@latest

COPY . .

CMD ["air", "-c", ".air.toml"]

FROM build-base AS build-production

# Add non root user
RUN useradd -u 1001 nonroot

COPY . .

# Compile application during build rather than at runtime
# Add flags to statically link binary
RUN go build \
  -ldflags="-linkmode external -extldflags -static" \
  -tags netgo \
  -o http \
  ./cmd/http/main.go #  NOTE: For now, we only build the http server, but other binaries can be added if needed e.g. `./cmd/cli/main.go`

# Use separate stage for deployable image
FROM scratch

WORKDIR /

# Copy the passwd file
COPY --from=build-production /etc/passwd /etc/passwd

# Copy the app binary from the build stage
COPY --from=build-production /app/http http

# Use nonroot user
USER nonroot

# Indicate expected port
EXPOSE 8080

CMD ["/http"]

