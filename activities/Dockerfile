# Build stage
FROM golang:latest AS build
RUN apt-get update && apt-get install -y --no-install-recommends git ca-certificates && rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /api ./cmd/api

# Runtime
FROM alpine:latest
ENV GIN_MODE=release
COPY --from=build /api /bin/api
EXPOSE 8080
ENTRYPOINT ["/bin/api"]
