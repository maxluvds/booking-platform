FROM golang:1.26-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/bin/notification-service ./services/notification-service/cmd/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

COPY --from=builder /app/bin/notification-service .
COPY --from=builder /app/services/notification-service/config ./services/notification-service/config
COPY --from=builder /app/migrations/notification-service ./migrations/notification-service

CMD ["./notification-service"]