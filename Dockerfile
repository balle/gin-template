# Build
FROM golang:1.25-alpine AS builder
WORKDIR /code
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o app .

# Run
FROM alpine:latest
WORKDIR /app
COPY --from=builder /code/app .
COPY --from=builder /code/templates ./templates

EXPOSE 8000
CMD ["./app"]