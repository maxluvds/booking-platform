FROM golang:1.26-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/bin/booking-service ./services/booking-service/cmd/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

COPY --from=builder /app/bin/booking-service .
COPY --from=builder /app/services/booking-service/config ./services/booking-service/config
COPY --from=builder /app/migrations/booking-service ./migrations/booking-service

EXPOSE 50054

CMD ["./booking-service"]