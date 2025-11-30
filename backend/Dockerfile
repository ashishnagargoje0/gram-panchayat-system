# backend/Dockerfile - safe two-stage Go build
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server

FROM alpine:latest
RUN apk add --no-cache ca-certificates
COPY --from=builder /server /usr/local/bin/server
EXPOSE 8080
USER 1000:1000
ENTRYPOINT ["/usr/local/bin/server"]
