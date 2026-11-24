FROM golang:alpine AS builder
WORKDIR /app
COPY server.go .
# Build พร้อม Optimized flags
RUN go build -ldflags="-s -w" -o server server.go

FROM alpine:latest
WORKDIR /root
COPY --from=builder /app/server .
CMD ["./server"]
