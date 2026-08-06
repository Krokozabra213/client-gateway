FROM golang:1.25-alpine AS build

RUN apk add --no-cache git

WORKDIR /build

COPY go.mod go.sum ./

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /client-gateway ./cmd/client-gateway

FROM alpine:3.23

WORKDIR /app
COPY --from=build /client-gateway .

CMD ["./client-gateway"]
