package handler

import (
	"context"

	"github.com/abu-umair/be-lms-go/internal/service"
	"github.com/abu-umair/be-lms-go/internal/utils"
	"github.com/abu-umair/be-lms-go/pb/store"
)

type storeHandler struct {
	store.UnimplementedStoreServiceServer

	storeService service.IStoreService //? layer service
}

func (sh *storeHandler) CreateStore(ctx context.Context, request *store.CreateStoreRequest) (*store.CreateStoreResponse, error) {
	//? validasi request
	validationErrors, err := utils.CheckValidation(request)
	if err != nil {
		return nil, err
	}

	if validationErrors != nil {
		return &store.CreateStoreResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}

	res, err := sh.storeService.CreateStore(ctx, request)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (sh *storeHandler) DetailStore(ctx context.Context, request *store.DetailStoreRequest) (*store.DetailStoreResponse, error) {
	//? validasi request
	validationErrors, err := utils.CheckValidation(request)
	if err != nil {
		return nil, err
	}

	if validationErrors != nil {
		return &store.DetailStoreResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}

	res, err := sh.storeService.DetailStore(ctx, request)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (sh *storeHandler) EditStore(ctx context.Context, request *store.EditStoreRequest) (*store.EditStoreResponse, error) {
	validationErrors, err := utils.CheckValidation(request)
	if err != nil {
		return nil, err
	}

	if validationErrors != nil {
		return &store.EditStoreResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}

	res, err := sh.storeService.EditStore(ctx, request)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (sh *storeHandler) DeleteStore(ctx context.Context, request *store.DeleteStoreRequest) (*store.DeleteStoreResponse, error) {
	validationErrors, err := utils.CheckValidation(request)
	if err != nil {
		return nil, err
	}

	if validationErrors != nil {
		return &store.DeleteStoreResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}

	res, err := sh.storeService.DeleteStore(ctx, request)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func NewStoreHandler(storeService service.IStoreService) *storeHandler {
	return &storeHandler{
		storeService: storeService,
	}
}
