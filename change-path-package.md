# Berikut adalah langkah-langkah untuk mengubahnya:

## 1. Ubah File go.mod
Buka file go.mod di root project Anda. Ubah baris pertama (baris module):

Dari:

module github.com/abu-umair/test-be-microservice
Menjadi:

module github.com/abu-umair/be-microservice


## 2. Update Semua Import di File .go
Ini bagian yang paling krusial. Semua file yang meng-import package internal Anda (seperti internal/entity) harus diperbarui manual atau otomatis.

Jika Anda menggunakan VS Code, Anda bisa melakukan Global Search & Replace:

Find: github.com/abu-umair/test-be-microservice

Replace: github.com/abu-umair/be-microservice

## 3. Jalankan Tidy Modul
Setelah mengubah semua import, jalankan perintah ini di terminal untuk memastikan dependensi Anda sinkron kembali:

```bash

go mod tidy
```

## 4. Sesuaikan di GitHub (Jika sudah dipush)
Jika project ini sudah ada di GitHub, jangan lupa untuk mengubah nama repositorinya di bagian Settings > Repository Name agar sesuai dengan nama modul yang baru (be-microservice).