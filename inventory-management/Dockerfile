FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/inventory-api ./cmd/api

FROM alpine:3.20

WORKDIR /app
COPY --from=builder /app/bin/inventory-api /app/inventory-api

EXPOSE 8080
CMD ["/app/inventory-api"]
