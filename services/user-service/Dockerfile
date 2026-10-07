FROM golang:1.26-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/bin/user-service ./services/user-service/cmd/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

COPY --from=builder /app/bin/user-service .
COPY --from=builder /app/services/user-service/config ./services/user-service/config
COPY --from=builder /app/migrations/user-service ./migrations/user-service

EXPOSE 50053

CMD ["./user-service"]