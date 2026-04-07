FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -ldflags="-w -s" -o bin/server ./cmd/server/main.go

# ── Runtime image ──────────────────────────────────────────────────
FROM alpine:3.19

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /app/bin/server .
COPY --from=builder /app/config ./config

RUN mkdir -p logs

EXPOSE 8080

CMD ["./server"]
