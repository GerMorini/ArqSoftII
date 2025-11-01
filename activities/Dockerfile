# Build stage
FROM golang:latest AS build
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /api ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /reindex ./cmd/reindex

# Runtime
FROM alpine:latest
ENV GIN_MODE=release
COPY --from=build /api /bin/api
COPY --from=build /reindex /bin/reindex
EXPOSE 8080
ENTRYPOINT ["/bin/api"]
