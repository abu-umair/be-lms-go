package utils

import (
	"strconv"
	"time"

	"github.com/shopspring/decimal"
)

// StringToPtr mengubah string menjadi pointer string.
// Jika kosong, return nil agar tidak muncul di JSON gRPC (optional).
func StringToPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// PtrStringToPtr menangani field yang sudah *string di struct.
// Berguna untuk meneruskan nilai pointer dari DB langsung ke Response.
func PtrStringToPtr(s *string) *string {
	if s == nil || *s == "" {
		return nil
	}
	return s
}

// Int64ToPtr mengubah int64 menjadi pointer.
func Int64ToPtr(i int64) *int64 {
	if i == 0 {
		return nil
	}
	return &i
}

// Int32ToPtr mengubah int32 menjadi pointer.
func Int32ToPtr(i int32) *int32 {
	if i == 0 {
		return nil
	}
	return &i
}

// PtrInt32ToPtr menangani *int32.
func PtrInt32ToPtr(i *int32) *int32 {
	if i == nil || *i == 0 {
		return nil
	}
	return i
}

// Int32ToStringPtr mengubah int32 menjadi pointer string.
func Int32ToStringPtr(i int32) *string {
	s := strconv.FormatInt(int64(i), 10)
	return &s
}

// PtrInt64ToPtr menangani *int64 (seperti field Capacity).
func PtrInt64ToPtr(i *int64) *int64 {
	return i
}

// Int64ToStringPtr mengubah int64 menjadi pointer string.
func Int64ToStringPtr(i int64) *string {
	s := strconv.FormatInt(i, 10)
	return &s
}

// PtrInt64ToStringPtr mengubah *int64 menjadi pointer string.
func PtrInt64ToStringPtr(i *int64) *string {
	if i == nil {
		return nil
	}
	s := strconv.FormatInt(*i, 10)
	return &s
}

// DecimalToPtr mengubah decimal.Decimal menjadi pointer string (format gRPC).
func DecimalToPtr(d decimal.Decimal) *string {
	if d.IsZero() {
		return nil
	}
	s := d.String()
	return &s
}

// PtrDecimalToPtr menangani *decimal.Decimal (seperti field Price & Discount).
func PtrDecimalToPtr(d *decimal.Decimal) *string {
	if d == nil || d.IsZero() {
		return nil
	}
	s := d.String()
	return &s
}

// TimeToPtr mengubah time.Time menjadi string pointer (format ISO8601).
func TimeToPtr(t time.Time) *string {
	if t.IsZero() {
		return nil
	}
	s := t.Format(time.RFC3339)
	return &s
}

// PtrTimeToPtr menangani *time.Time (seperti field DeletedAt).
func PtrTimeToPtr(t *time.Time) *string {
	if t == nil || t.IsZero() {
		return nil
	}
	s := t.Format(time.RFC3339)
	return &s
}
