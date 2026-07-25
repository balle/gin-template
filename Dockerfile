# Build
FROM golang:1.25-alpine AS builder
WORKDIR /code
COPY . .
RUN go build -o app main.go

# Run
FROM alpine:latest
WORKDIR /app
COPY --from=builder /code/app .
COPY --from=builder /code/templates ./templates

EXPOSE 8000
CMD ["./app"]