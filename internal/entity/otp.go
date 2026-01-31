package entity

import "time"

type UserOTP struct {
	Email     string    `db:"email"`
	OTPCode   string    `db:"otp_code"`
	ExpiredAt time.Time `db:"expired_at"`
    CreatedAt time.Time `db:"created_at"`
}
