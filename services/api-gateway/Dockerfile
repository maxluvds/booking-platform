FROM golang:1.26-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/bin/api-gateway ./services/api-gateway/cmd/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

COPY --from=builder /app/bin/api-gateway .

COPY --from=builder /app/services/api-gateway/config ./services/api-gateway/config

EXPOSE 8080

CMD ["./api-gateway"]