# Build stage
FROM golang:1.22-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /gateway .

# Runtime
FROM alpine:3.20
ENV GIN_MODE=release
COPY --from=build /gateway /bin/gateway
EXPOSE 8080
ENTRYPOINT ["/bin/gateway"]
