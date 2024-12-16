FROM golang:1.20-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o pubsub-publisher .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/pubsub-publisher .
COPY order_message.yaml .
CMD ["./pubsub-publisher"]