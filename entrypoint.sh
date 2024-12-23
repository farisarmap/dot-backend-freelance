#!/bin/sh

# Menunggu Postgres siap
echo "Menunggu Postgres..."

while ! nc -z $DB_HOST $DB_PORT; do
  sleep 1
done

echo "Postgres siap - menjalankan migrasi"

# Menjalankan migrasi
./migrate -config=config.json -dir=migration up

echo "Migrasi selesai - memulai aplikasi"

# Menjalankan aplikasi utama
./main
