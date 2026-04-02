FROM golang:1.25.5 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .
WORKDIR /app/
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o trip-service ./services/trip-service/cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/trip-service .
CMD ["./trip-service"] 
