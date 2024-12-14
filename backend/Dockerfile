# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Production stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates gettext
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/config.yaml ./config.yaml.template
COPY docker-entrypoint.sh .
COPY --from=builder /app/.env ./
RUN chmod +x docker-entrypoint.sh

EXPOSE 8080
ENTRYPOINT ["./docker-entrypoint.sh"]
