FROM golang:1.23.1-alpine AS build

ENV GO111MODULE=on

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

COPY config.json .

RUN go build -o /app/main cmd/main.go

RUN go build -o /app/migrate cmd/migrate/main.go

FROM alpine:latest

RUN apk update && apk add --no-cache netcat-openbsd

WORKDIR /app

COPY --from=build /app/main .
COPY --from=build /app/migrate .
COPY --from=build /app/config.json .

COPY --from=build /app/migration ./migration
COPY --from=build /app/test ./test

EXPOSE 8080

COPY entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh

ENTRYPOINT ["/app/entrypoint.sh"]
