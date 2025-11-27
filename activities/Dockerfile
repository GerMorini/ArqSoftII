# Build stage
FROM golang:1.24-alpine AS build
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN go test ./internal/services/... -v
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /api ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /reindex ./cmd/reindex

# Runtime
FROM alpine:3.20
ENV GIN_MODE=release
COPY --from=build /api /bin/api
COPY --from=build /reindex /bin/reindex
EXPOSE 8080
ENTRYPOINT ["/bin/api"]
