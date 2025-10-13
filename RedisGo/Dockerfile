FROM golang:1.24.3-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o redisgo

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/redisgo .
EXPOSE 6379
CMD ["./redisgo"]
