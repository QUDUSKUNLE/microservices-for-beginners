FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN  go mod download


COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server main.go


FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/server .

ENV PORT=8002

EXPOSE 8002

ENTRYPOINT [ "./server" ]