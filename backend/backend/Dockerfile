FROM golang:1.25.4-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main .

FROM alpine:3.21.5

WORKDIR /root/

COPY --from=build /app/main .
COPY --from=build /app/config ./config
# COPY --from=build /app/migrations ./migrations

EXPOSE 3000

ENV SERVER_PORT=3000

CMD ["./main"]
