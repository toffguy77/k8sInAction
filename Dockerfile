FROM golang:1.21-alpine AS build_base

RUN apk add --no-cache git

# Set the Current Working Directory inside the container
WORKDIR /tmp/http_server

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.* .

RUN go mod download

COPY . .

# Build the Go app
RUN go build -o ./out/http_server ./cmd/main.go

# Start fresh from a smaller image
FROM alpine:latest
RUN apk add ca-certificates

COPY --from=build_base /tmp/http_server/out/http_server /app/http_server

# This container exposes port 8080 to the outside world
EXPOSE 8080

# Run the binary program produced by `go install`
CMD ["/app/http_server"]