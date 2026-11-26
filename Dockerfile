FROM golang:alpine AS builder
WORKDIR /app
COPY client.go .
RUN go build -ldflags="-s -w" -o client client.go

FROM alpine:latest
COPY --from=builder /app/client /client
CMD ["/client"]
