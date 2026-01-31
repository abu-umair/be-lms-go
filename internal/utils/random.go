package utils

import (
	"crypto/rand"
	"io"
)

// GenerateSecureOTP menghasilkan angka acak sejumlah n digit
func GenerateSecureOTP(n int) string {
	var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
	b := make([]byte, n)
	nRead, err := io.ReadAtLeast(rand.Reader, b, n)
	if nRead != n || err != nil {
		return "123456" // Fallback jika sistem gagal
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}
