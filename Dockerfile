FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod .
COPY main.go .

RUN <<EOF
go mod tidy 
go build
EOF

FROM debian:12

RUN apt-get update && apt-get install -y \
curl \
man

WORKDIR /app
COPY --from=builder /app/mcpsh .

ENTRYPOINT ["./mcpsh"]

