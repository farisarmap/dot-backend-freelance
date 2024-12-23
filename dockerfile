# Tahap Build
FROM golang:1.23.1-alpine AS build

# Mengatur environment variable untuk modul Go
ENV GO111MODULE=on

# Menginstal git dan netcat
RUN apk update && apk add --no-cache git

# Menetapkan direktori kerja
WORKDIR /app

# Menyalin file go.mod dan go.sum
COPY go.mod go.sum ./

# Mengunduh dependensi
RUN go mod download

# Menyalin seluruh kode sumber ke dalam kontainer
COPY . .

# Menyalin config.json ke direktori kerja
COPY config.json .

# Membangun aplikasi utama
RUN go build -o /app/main cmd/main.go

# Membangun binary migrate
RUN go build -o /app/migrate cmd/migrate/main.go

# Tahap Run
FROM alpine:latest

# Menginstal netcat untuk skrip penunggu
RUN apk update && apk add --no-cache netcat-openbsd

# Menetapkan direktori kerja
WORKDIR /app

# Menyalin aplikasi yang telah dibangun dari tahap build
COPY --from=build /app/main .
COPY --from=build /app/migrate .
COPY --from=build /app/config.json .

COPY --from=build /app/migration ./migration

# Mengekspos port aplikasi
EXPOSE 8080

# Menyalin dan menetapkan skrip entrypoint
COPY entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh

# Menetapkan entrypoint
ENTRYPOINT ["/app/entrypoint.sh"]
