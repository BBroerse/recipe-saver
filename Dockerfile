# syntax=docker/dockerfile:1
FROM golang:tip-alpine3.22 AS builder

WORKDIR /app
COPY . .

RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o dist/api .

FROM alpine:3.22.1

# Create a minimal user
RUN addgroup -S app && adduser -S app -G app && \
  apk --no-cache add ca-certificates

WORKDIR /root
COPY --from=builder /app/dist/api .

USER app
EXPOSE 8080
ENTRYPOINT ["./api"]