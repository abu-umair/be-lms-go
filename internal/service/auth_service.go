package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/abu-umair/be-lms-go/internal/entity"
	jwtentity "github.com/abu-umair/be-lms-go/internal/entity/jwt"
	"github.com/abu-umair/be-lms-go/internal/repository"
	"github.com/abu-umair/be-lms-go/internal/utils"
	"github.com/abu-umair/be-lms-go/pb/auth"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	gocache "github.com/patrickmn/go-cache"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type IAuthService interface {
	Register(ctx context.Context, request *auth.RegisterRequest) (*auth.RegisterResponse, error)
	Login(ctx context.Context, request *auth.LoginRequest) (*auth.LoginResponse, error)
	Logout(ctx context.Context, request *auth.LogoutRequest) (*auth.LogoutResponse, error)
	ChangePassword(ctx context.Context, request *auth.ChangePasswordRequest) (*auth.ChangePasswordResponse, error)
	GetProfile(ctx context.Context, request *auth.GetProfileRequest) (*auth.GetProfileResponse, error)
	RequestOTP(ctx context.Context, request *auth.RequestOTPRequest) (*auth.RequestOTPResponse, error)
	Verify(ctx context.Context, request *auth.VerifyRequest) (*auth.VerifyResponse, error)
}

type authService struct {
	authRepository repository.IAuthRepository
	cacheService   *gocache.Cache // Bisa digunakan nanti saat pindah ke Redis
	messageSender  IMessageSender // Untuk fleksibilitas Email/WA
}

func (as *authService) Register(ctx context.Context, request *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	//? apakah password sama dengan confirm password
	if request.Password != request.PasswordConfirmation {
		return &auth.RegisterResponse{
			Base: utils.BadRequestResponse("Password is not matched"),
		}, nil
	}

	//? ngecek email ke DB
	//* layer repository, utk akses DB (clean arsitektur)
	user, err := as.authRepository.GetUserByEmail(ctx, request.Email)
	if err != nil {
		return nil, err
	}

	//* jika emal sudah terdaftar/ada, di error in
	if user != nil {
		return &auth.RegisterResponse{
			Base: utils.BadRequestResponse("User already exist"),
		}, nil
	}

	//? Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), 10)
	if err != nil {
		return nil, err
	}

	//? Insert ke DB
	newUser := entity.Users{
		Id:        uuid.NewString(),
		FullName:  request.FullName,
		Email:     request.Email,
		Password:  string(hashedPassword),
		RoleCode:  entity.UserRoleUser,
		CreatedAt: time.Now(),
		CreatedBy: &request.FullName,
	}

	err = as.authRepository.InsertUser(ctx, &newUser)
	if err != nil {
		return nil, err
	}

	return &auth.RegisterResponse{
		Base: utils.SuccessResponse("User is registered"),
	}, nil
}

func (as *authService) Login(ctx context.Context, request *auth.LoginRequest) (*auth.LoginResponse, error) {
	//* check apakah email ada di database
	user, err := as.authRepository.GetUserByEmail(ctx, request.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return &auth.LoginResponse{
			Base: utils.BadRequestResponse("User is not registered"),
		}, nil
	}

	//* check apakah password sama dengan password di database
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)) //? mengecek password yang hash dengan password yang diinput/request
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) { //? jika password tidak sama
			return nil, status.Errorf(codes.Unauthenticated, "Unauthenticated") //?kembalikan error unauthenticated dari middleware (gRPC)
		}
		return nil, err
	}

	//* generate jwt
	now := time.Now()
	var verifiedStr string
	if user.VerifiedAt != nil {
		verifiedStr = user.VerifiedAt.Format(time.RFC3339)
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtentity.JwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.Id,
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
		Email:      user.Email,
		FullName:   user.FullName,
		Role:       user.RoleCode,
		VerifiedAt: verifiedStr,
	})
	secretKey := os.Getenv("JWT_SECRET")
	accessToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return nil, err
	}

	//* kirim response
	return &auth.LoginResponse{
		Base:        utils.SuccessResponse("Login successful"),
		AccessToken: accessToken,
	}, nil

}

func (as *authService) Logout(ctx context.Context, request *auth.LogoutRequest) (*auth.LogoutResponse, error) {
	//? dapatkan token dari metadata
	jwtToken, err := jwtentity.ParseTokenFromContext(ctx)
	if err != nil {
		return nil, err
	}

	tokenClaims, err := jwtentity.GetClaimsFromContext(ctx)

	if err != nil {
		return nil, err
	}

	//? kembalikan token tadi hingga menjadi entity jwt

	//? kita masukkan token ke dalam memory db / cache
	as.cacheService.Set(jwtToken, "", time.Duration(tokenClaims.ExpiresAt.Time.Unix()-time.Now().Unix())*time.Second)

	//? kirim response
	return &auth.LogoutResponse{
		Base: utils.SuccessResponse("Logout success"),
	}, nil
}

func (as *authService) ChangePassword(ctx context.Context, request *auth.ChangePasswordRequest) (*auth.ChangePasswordResponse, error) {
	//*Cek apakah new pass confirmation matched
	if request.NewPassword != request.NewPasswordConfirmation {
		return &auth.ChangePasswordResponse{
			Base: utils.BadRequestResponse("New password is not matched"),
		}, nil
	}

	//* Cek apakah old password sama
	jwtToken, err := jwtentity.ParseTokenFromContext(ctx)
	if err != nil {
		return nil, err
	}

	claims, err := jwtentity.GetClaimsFromToken(jwtToken)
	if err != nil {
		return nil, err
	}

	user, err := as.authRepository.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return &auth.ChangePasswordResponse{
			Base: utils.BadRequestResponse("User does not exist"),
		}, nil
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.OldPassword))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return &auth.ChangePasswordResponse{
				Base: utils.BadRequestResponse("Old password is not matched"),
			}, nil
		}
		return nil, err
	}

	//* Update new password ke DB
	hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), 10)
	if err != nil {
		return nil, err
	}

	err = as.authRepository.UpdateUserPassword(ctx, user.Id, string(hashedNewPassword), user.FullName)
	if err != nil {
		return nil, err
	}

	//* Kirim response
	return &auth.ChangePasswordResponse{
		Base: utils.SuccessResponse("Change password success"),
	}, nil
}

func (as *authService) GetProfile(ctx context.Context, request *auth.GetProfileRequest) (*auth.GetProfileResponse, error) {
	//* Get data token
	claims, err := jwtentity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	//* Ambil data dari DB
	user, err := as.authRepository.GetUserByEmail(ctx, claims.Email)

	if err != nil {
		return nil, err
	}

	if user == nil {
		return &auth.GetProfileResponse{
			Base: utils.BadRequestResponse("User doesn't exist"),
		}, nil
	}

	//* Buat/Kirim Response
	// ? apakah user sudah terverifikasi atau belum
	var verifiedAt *timestamppb.Timestamp //?default-nya nil
	if user.VerifiedAt != nil {
		verifiedAt = timestamppb.New(*user.VerifiedAt)
	}

	return &auth.GetProfileResponse{
		Base:        utils.SuccessResponse("Get Profile success"),
		UserId:      claims.Subject,
		FullName:    claims.FullName,
		Email:       claims.Email,
		RoleCode:    claims.Role,
		VerifiedAt:  verifiedAt, //?tidak akan tampil jika nil
		MemberSince: timestamppb.New(user.CreatedAt),
	}, nil
}

// ! Method tambahan untuk Request/Resend OTP
func (as *authService) RequestOTP(ctx context.Context, request *auth.RequestOTPRequest) (*auth.RequestOTPResponse, error) {
	//* Get data token
	claims, err := jwtentity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	//* --- LOGIKA RATE LIMITING ---
	lastOTP, err := as.authRepository.GetOTPByEmail(ctx, claims.Email)
	if err == nil { //? Jika data ditemukan
		//? Cek apakah permintaan terakhir kurang dari 60 detik yang lalu
		if time.Since(lastOTP.CreatedAt).Seconds() < 60 {
			remaining := 60 - int(time.Since(lastOTP.CreatedAt).Seconds())

			return &auth.RequestOTPResponse{
				Base: utils.BadRequestResponse(fmt.Sprintf("Silakan tunggu %d detik lagi untuk meminta kode baru", remaining)),
			}, nil
		}
	}

	//* Generate 6 angka random
	code := utils.GenerateSecureOTP(6)

	//* Simpan ke DB (Postgres)
	otp := &entity.UserOTP{
		Email:     claims.Email,
		OTPCode:   code,
		ExpiredAt: time.Now().Add(5 * time.Minute),
		CreatedAt: time.Now(),
	}
	as.authRepository.UpsertOTP(ctx, otp)

	//* Kirim via Email (lewat interface)
	subject := "Kode Verifikasi OTP Anda"
	htmlBody := utils.GetOTPEmailTemplate(code)
	go func() {
		errSend := as.messageSender.Send(claims.Email, subject, htmlBody)
		if errSend != nil {
			// Log error kirim email di sini (misal pakai logrus/zap)
			fmt.Printf("Error sending email to %s: %v\n", claims.Email, errSend)
		}
	}()

	return &auth.RequestOTPResponse{
		Base: utils.SuccessResponse("Send or Resend OTP success"),
	}, nil
}

func (as *authService) Verify(ctx context.Context, request *auth.VerifyRequest) (*auth.VerifyResponse, error) {
	//* Get data token
	claims, err := jwtentity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// * Get OTP From repo
	otpData, err := as.authRepository.GetOTPByEmail(ctx, claims.Email)
	if err != nil {
		return &auth.VerifyResponse{
			Base: utils.BadRequestResponse("OTP tidak ditemukan atau kadaluarsa"),
		}, nil
	}

	//* Cek Expiry
	if time.Now().After(otpData.ExpiredAt) {
		return &auth.VerifyResponse{
			Base: utils.BadRequestResponse("OTP sudah kadaluarsa"),
		}, nil
	}

	//* Cek Kecocokan Kode
	if otpData.OTPCode != request.CodeOtp {
		return &auth.VerifyResponse{
			Base: utils.BadRequestResponse("kode OTP salah"),
		}, nil
	}

	//* Update User jadi Verified
	user, err := as.authRepository.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return &auth.VerifyResponse{
			Base: utils.BadRequestResponse("User doesn't exist"),
		}, nil
	}

	//* VerifyEmail
	//? jika sudah verifikasi, kembalikan error
	if user.VerifiedAt != nil {
		return &auth.VerifyResponse{
			Base: utils.BadRequestResponse("Email already verified"),
		}, nil
	}

	//? jika belum verifikasi, update verified_at
	err = as.authRepository.MarkAsVerified(ctx, user.Id)
	if err != nil {
		return nil, err
	}

	//* Hapus OTP setelah berhasil
	as.authRepository.DeleteOTP(ctx, claims.Email)

	//* Buat/Kirim Response
	return &auth.VerifyResponse{
		Base: utils.SuccessResponse("Verify success"),
	}, nil
}

func NewAuthService(authRepository repository.IAuthRepository, cacheService *gocache.Cache, sender IMessageSender) IAuthService {
	return &authService{
		authRepository: authRepository,
		cacheService:   cacheService,
		messageSender:  sender,
	}
}
