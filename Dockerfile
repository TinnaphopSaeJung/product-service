# ---------- Build stage ----------
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o product-service ./cmd/api


# ---------- Run stage ----------
FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/product-service .

EXPOSE 8080

CMD ["./product-service"]