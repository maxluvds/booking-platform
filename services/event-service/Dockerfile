FROM golang:1.26-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/bin/event-service ./services/event-service/cmd/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

COPY --from=builder /app/bin/event-service .
COPY --from=builder /app/services/event-service/config ./services/event-service/config
COPY --from=builder /app/migrations/event-service ./migrations/event-service

EXPOSE 50051

CMD ["./event-service"]