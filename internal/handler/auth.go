package handler

import (
	"context"

	"github.com/abu-umair/be-lms-go/internal/service"
	"github.com/abu-umair/be-lms-go/internal/utils"
	"github.com/abu-umair/be-lms-go/pb/auth"
)

type authHandler struct {
	auth.UnimplementedAuthServiceServer

	authService service.IAuthService //? layer service
}

func (sh *authHandler) Register(ctx context.Context, request *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	//? validasi request
	validationErrors, err := utils.CheckValidation(request)
	if err != nil {
		return nil, err
	}

	if validationErrors != nil {
		return &auth.RegisterResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}

	//?biasanya ada proses register (bisnis proses), di buat di layer service
	res, err := sh.authService.Register(ctx, request)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (sh *authHandler) Login(ctx context.Context, request *auth.LoginRequest) (*auth.LoginResponse, error) {
	validationErrors, err := utils.CheckValidation(request)
	if err != nil {
		return nil, err
	}

	if validationErrors != nil {
		return &auth.LoginResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}

	res, err := sh.authService.Login(ctx, request)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (sh *authHandler) Logout(ctx context.Context, request *auth.LogoutRequest) (*auth.LogoutResponse, error) {
	validationErrors, err := utils.CheckValidation(request)
	if err != nil {
		return nil, err
	}

	if validationErrors != nil {
		return &auth.LogoutResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}

	res, err := sh.authService.Logout(ctx, request)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (sh *authHandler) ChangePassword(ctx context.Context, request *auth.ChangePasswordRequest) (*auth.ChangePasswordResponse, error) {
	validationErrors, err := utils.CheckValidation(request)
	if err != nil {
		return nil, err
	}

	if validationErrors != nil {
		return &auth.ChangePasswordResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}

	res, err := sh.authService.ChangePassword(ctx, request)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (sh *authHandler) GetProfile(ctx context.Context, request *auth.GetProfileRequest) (*auth.GetProfileResponse, error) {
	//? langsung proses (tidak ada inputan)
	res, err := sh.authService.GetProfile(ctx, request)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (sh *authHandler) RequestOTP(ctx context.Context, request *auth.RequestOTPRequest) (*auth.RequestOTPResponse, error) {
	//? langsung proses (tidak ada inputan)
	res, err := sh.authService.RequestOTP(ctx, request)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (sh *authHandler) Verify(ctx context.Context, request *auth.VerifyRequest) (*auth.VerifyResponse, error) {
	validationErrors, err := utils.CheckValidation(request)
	if err != nil {
		return nil, err
	}

	if validationErrors != nil {
		return &auth.VerifyResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}

	res, err := sh.authService.Verify(ctx, request)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func NewAuthHandler(authService service.IAuthService) *authHandler {
	return &authHandler{
		authService: authService,
	}
}
