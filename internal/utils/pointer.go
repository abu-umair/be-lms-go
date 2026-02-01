package utils

import "github.com/shopspring/decimal"

// StringToPtr mengubah string menjadi pointer string.
// Jika string kosong (""), maka mengembalikan nil.
// Ini sangat berguna untuk field 'optional' di gRPC.
func StringToPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// Int64ToPtr mengubah int64 menjadi pointer int64.
// Jika nilainya 0, Anda bisa memilih mengembalikan nil atau tetap pointernya.
// Di sini kita asumsikan jika 0 tetap dikirim, kecuali nilainya nil di DB.
func Int64ToPtr(i int64) *int64 {
	return &i
}

// DecimalToPtr mengubah decimal.Decimal menjadi pointer string (untuk gRPC).
// Karena gRPC tidak punya tipe Decimal, biasanya kita kirim sebagai string.
func DecimalToPtr(d decimal.Decimal) *string {
	if d.IsZero() {
		return nil
	}
	s := d.String()
	return &s
}
