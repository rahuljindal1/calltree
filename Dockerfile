# -----------------------
# Build stage
# -----------------------
FROM golang:1.22-alpine AS builder

WORKDIR /app

ENV CGO_ENABLED=0

COPY go.mod ./
RUN go mod download

COPY . .

RUN go build \
    -ldflags="-s -w" \
    -o calltree \
    ./cmd/calltree


# -----------------------
# Runtime stage
# -----------------------
FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/calltree /usr/local/bin/calltree

ENTRYPOINT ["calltree"]
