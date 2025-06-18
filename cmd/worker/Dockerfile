# build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ../../worker-config.yaml /config/worker-config.yaml

RUN go build -o worker ./cmd/worker

FROM alpine:3.22.0

WORKDIR /app

COPY --from=builder /app/worker ./
COPY --from=builder /config/worker-config.yaml /config/worker-config.yaml

CMD ["./worker"]
