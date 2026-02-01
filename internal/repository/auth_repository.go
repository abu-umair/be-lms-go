package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/abu-umair/be-lms-go/internal/entity"
	"github.com/jmoiron/sqlx"
)

type IAuthRepository interface {
	GetUserByEmail(ctx context.Context, email string) (*entity.Users, error) //?ctx: mengakses DB nya. User: entitynya (table di DB, terhubung di entity), (User, error) adalah return data atau error
	InsertUser(ctx context.Context, user *entity.Users) error
	UpdateUserPassword(ctx context.Context, userId string, hashedPassword string, updatedBy string) error
	MarkAsVerified(ctx context.Context, userId string) error

	// otp
	UpsertOTP(ctx context.Context, otp *entity.UserOTP) error
	GetOTPByEmail(ctx context.Context, email string) (*entity.UserOTP, error)
	DeleteOTP(ctx context.Context, email string) error
}

type authRepository struct {
	db *sqlx.DB //?Menyimpan koneksi database
}

func (ar *authRepository) GetUserByEmail(ctx context.Context, email string) (*entity.Users, error) {
	row := ar.db.QueryRowContext(ctx, "SELECT id, email, password, full_name, role_code, created_at,verified_at FROM users WHERE email = $1 AND deleted_at IS NULL", email)

	if row.Err() != nil {
		return nil, row.Err()
	}

	var user entity.Users
	err := row.Scan(
		&user.Id,
		&user.Email,
		&user.Password,
		&user.FullName,
		&user.RoleCode,
		&user.CreatedAt,
		&user.VerifiedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (ar *authRepository) InsertUser(ctx context.Context, user *entity.Users) error {
	_, err := ar.db.ExecContext(
		ctx,
		`INSERT INTO users (id, full_name, email, password, role_code, created_at, created_by, updated_at, updated_by, deleted_by)
                VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,

		user.Id,
		user.FullName,
		user.Email,
		user.Password,
		user.RoleCode,
		user.CreatedAt,
		user.CreatedBy,
		user.UpdatedAt,
		user.UpdatedBy,
		user.DeletedBy,
	)

	if err != nil {
		return err
	}

	return nil
}

func (ar *authRepository) UpdateUserPassword(ctx context.Context, userId string, hashedPassword string, updatedBy string) error {
	_, err := ar.db.ExecContext(
		ctx,
		"UPDATE users SET password = $1, updated_at = $2, updated_by = $3 WHERE id = $4",
		hashedPassword,
		time.Now(),
		updatedBy,
		userId,
	)

	if err != nil {
		return err
	}

	return nil
}

func (ar *authRepository) MarkAsVerified(ctx context.Context, userId string) error {
	_, err := ar.db.ExecContext(
		ctx,
		"UPDATE users SET verified_at = $1 WHERE id = $2",
		time.Now(),
		userId,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *authRepository) UpsertOTP(ctx context.Context, otp *entity.UserOTP) error {
	query := `
        INSERT INTO user_otps (email, otp_code, expired_at, created_at)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (email) DO UPDATE SET otp_code = $2, expired_at = $3, created_at = $4`
	_, err := r.db.ExecContext(ctx, query, otp.Email, otp.OTPCode, otp.ExpiredAt, otp.CreatedAt)
	return err
}

func (ar *authRepository) GetOTPByEmail(ctx context.Context, email string) (*entity.UserOTP, error) {
	var otp entity.UserOTP
	query := "SELECT email, otp_code, expired_at,created_at FROM user_otps WHERE email = $1"
	err := ar.db.QueryRowContext(ctx, query, email).Scan(&otp.Email, &otp.OTPCode, &otp.ExpiredAt, &otp.CreatedAt)
	return &otp, err
}

func (ar *authRepository) DeleteOTP(ctx context.Context, email string) error {
	_, err := ar.db.ExecContext(ctx, "DELETE FROM user_otps WHERE email = $1", email)
	return err
}

func NewAuthRepository(db *sqlx.DB) IAuthRepository {
	return &authRepository{db: db}
}
