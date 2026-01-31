### Generate proto service
```bash
protoc --go_out=./pb --go-grpc_out=./pb --proto_path=./proto --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative service/service.proto
```
### Run main.go
```bash
go run cmd/grpc/main.go
```

### install go-cache
```bash
go get github.com/patrickmn/go-cache
```

### yang belum ada
- cmd\rest\main.go
- internal\dto\webhook_recieve_invoice.go
- pkg\database\database_query.go
- storage/product


### install timestamp
```bash
go get google.golang.org/protobuf/types/known/timestamppb
```

### install gomail
```bash
go get gopkg.in/gomail.v2
```

### Untuk OTP pada tabler user_otps:
1. harus menggunakan cronjob untuk menghapus sisa data yang expired (meski sudah ada ada perintah hapus, tetapi untuk kasus jika user tidak melanjutkan OTP otomatis expired dan tidak melanjutkan OTP juga dan tersimpan di DB) 

### Run Server gofiber (utk upload)
```bash
go run cmd/rest/main.go
```

## install hot reload air
```bash
go install github.com/air-verse/air@latest

run: air
https://www.youtube.com/watch?v=V7fY776aGgM
```

## install shopspring/decimal untuk price/discount
```bash
go install github.com/air-verse/air@latest
```