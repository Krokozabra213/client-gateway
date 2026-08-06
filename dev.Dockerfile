FROM golang:1.25-alpine AS build

RUN apk add --no-cache git bash curl

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -o client-gateway ./cmd/client-gateway

FROM alpine:3.18

RUN echo "https://mirror.yandex.ru/mirrors/alpine/v3.18/main" > /etc/apk/repositories && \
    echo "https://mirror.yandex.ru/mirrors/alpine/v3.18/community" >> /etc/apk/repositories && \
    apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=build /app/client-gateway .
COPY --from=build /app/configs ./configs

RUN chmod +x /app/client-gateway
CMD ["./client-gateway"]
