# Build stage
FROM golang:1.26 AS builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o tiny-url ./cmd

# Run stage
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/tiny-url .
CMD ["./tiny-url"]