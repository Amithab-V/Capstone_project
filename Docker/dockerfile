# docker/Dockerfile
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o selfhealer ./cmd/selfhealer

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/selfhealer .
CMD ["./selfhealer"]
