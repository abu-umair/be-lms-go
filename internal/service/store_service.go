package service

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"time"

	"github.com/abu-umair/be-lms-go/internal/entity"
	jwtentity "github.com/abu-umair/be-lms-go/internal/entity/jwt"
	"github.com/abu-umair/be-lms-go/internal/repository"
	"github.com/abu-umair/be-lms-go/internal/utils"
	"github.com/abu-umair/be-lms-go/pb/store"
)

type IStoreService interface {
	CreateStore(ctx context.Context, request *store.CreateStoreRequest) (*store.CreateStoreResponse, error)
	DetailStore(ctx context.Context, request *store.DetailStoreRequest) (*store.DetailStoreResponse, error)
	EditStore(ctx context.Context, request *store.EditStoreRequest) (*store.EditStoreResponse, error)
	DeleteStore(ctx context.Context, request *store.DeleteStoreRequest) (*store.DeleteStoreResponse, error)
}

type storeService struct {
	db              *sql.DB
	storeRepository repository.IStoreRepository
}

func (ss *storeService) CreateStore(ctx context.Context, request *store.CreateStoreRequest) (*store.CreateStoreResponse, error) {
	//* Get data token
	claims, err := jwtentity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	//* apakah role user adl Owner
	if claims.Role != entity.UserRoleOwner {
		return nil, utils.UnauthenticatedResponse()
	}

	tx, err := ss.db.Begin()
	if err != nil {
		return nil, err
	}

	defer func() {
		if e := recover(); e != nil {
			if tx != nil {
				tx.Rollback() //?rollback jika ada error saan runtime
			}

			debug.PrintStack() //?agar ada stock tracenya yang digunakan utk debug
			panic(e)           //?agar bisa nyampai ke Middleware
		}
	}()

	defer func() {
		if err != nil && tx != nil {
			tx.Rollback() //?rollback jika ada error
		}
	}()

	storeRepo := ss.storeRepository.WithTransaction(tx)

	// *insert ke DB
	storeEntity := entity.Store{
		Id:            request.Id,
		Name:          request.Name,
		Address:       request.Address,
		ImageFileName: request.ImageFileName,
		CreatedAt:     time.Now(),
		CreatedBy:     claims.FullName,
	}

	err = storeRepo.CreateNewStore(ctx, &storeEntity)
	if err != nil {
		return nil, err
	}

	// *apakah image ada
	imagePath := filepath.Join("storage", storeEntity.Id, "store", request.ImageFileName)
	_, err = os.Stat(imagePath)
	if err != nil {
		if os.IsNotExist(err) {
			return &store.CreateStoreResponse{
				Base: utils.BadRequestResponse("File not found"),
			}, nil
		}
		return nil, err
	}

	err = tx.Commit() //?harus dicommit agar data tersimpan
	if err != nil {
		return nil, err
	}

	// *success
	return &store.CreateStoreResponse{
		Base: utils.SuccessResponse("Store successfully created"),
	}, nil
}

func (ss *storeService) DetailStore(ctx context.Context, request *store.DetailStoreRequest) (*store.DetailStoreResponse, error) {

	//* Get data token
	claims, err := jwtentity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	//* apakah role user adl Owner
	if claims.Role != entity.UserRoleOwner {
		return nil, utils.UnauthenticatedResponse()
	}

	// * Get store by store_id
	storeEntity, err := ss.storeRepository.GetStoreById(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	//* Apabila null store_id, return not found
	if storeEntity == nil {
		return &store.DetailStoreResponse{
			Base: utils.NotFoundResponse("Store not found"),
		}, nil
	}

	// *success
	return &store.DetailStoreResponse{
		Base:          utils.SuccessResponse("Store Detail Success"),
		Id:            storeEntity.Id,
		Name:          storeEntity.Name,
		Address:       storeEntity.Address,
		ImageFileName: fmt.Sprintf("%s/%s/store/%s", os.Getenv("STORAGE_SERVICE_URL"), storeEntity.Id, storeEntity.ImageFileName),
	}, nil
}

func (ss *storeService) EditStore(ctx context.Context, request *store.EditStoreRequest) (*store.EditStoreResponse, error) {
	//* Get data token
	claims, err := jwtentity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	//* apakah role user adl Owner
	if claims.Role != entity.UserRoleOwner {
		return nil, utils.UnauthenticatedResponse()
	}

	// *Apakah Id store ada di DB
	storeEntity, err := ss.storeRepository.GetStoreById(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	if storeEntity == nil {
		return &store.EditStoreResponse{
			Base: utils.NotFoundResponse("Store not found"),
		}, nil
	}

	tx, err := ss.db.Begin()
	if err != nil {
		return nil, err
	}

	defer func() {
		if e := recover(); e != nil {
			if tx != nil {
				tx.Rollback() //?rollback jika ada error saan runtime
			}

			debug.PrintStack() //?agar ada stock tracenya yang digunakan utk debug
			panic(e)           //?agar bisa nyampai ke Middleware
		}
	}()

	defer func() {
		if err != nil && tx != nil {
			tx.Rollback() //?rollback jika ada error
		}
	}()

	storeRepo := ss.storeRepository.WithTransaction(tx)

	// *update ke DB
	newStore := entity.Store{
		Id:            request.Id,
		Name:          request.Name,
		Address:       request.Address,
		ImageFileName: request.ImageFileName,
		UpdatedAt:     time.Now(),
		UpdatedBy:     &claims.FullName,
	}

	err = storeRepo.UpdateStore(ctx, &newStore)
	if err != nil {
		return nil, err
	}

	// *jika ada image baru, hapus image lama
	if storeEntity.ImageFileName != request.ImageFileName {
		newImagePath := filepath.Join("storage", request.Id, "store", request.ImageFileName)
		_, err := os.Stat(newImagePath)
		if err != nil {
			if os.IsNotExist(err) {
				return &store.EditStoreResponse{
					Base: utils.BadRequestResponse("Image not found"),
				}, nil
			}
			return nil, err
		}

		oldImagePath := filepath.Join("storage", storeEntity.Id, "store", storeEntity.ImageFileName)
		err = os.Remove(oldImagePath)
		if err != nil {
			return nil, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	// *success
	return &store.EditStoreResponse{
		Base: utils.SuccessResponse("Edit Store Success"),
		Id:   request.Id,
	}, nil
}

func (ss *storeService) DeleteStore(ctx context.Context, request *store.DeleteStoreRequest) (*store.DeleteStoreResponse, error) {
	//* Get data token
	claims, err := jwtentity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	//* apakah role user adl Owner
	if claims.Role != entity.UserRoleOwner {
		return nil, utils.UnauthenticatedResponse()
	}

	// *Apakah Id store ada di DB
	storeEntity, err := ss.storeRepository.GetStoreById(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	if storeEntity == nil {
		return &store.DeleteStoreResponse{
			Base: utils.NotFoundResponse("Store not found"),
		}, nil
	}

	tx, err := ss.db.Begin()
	if err != nil {
		return nil, err
	}

	defer func() {
		if e := recover(); e != nil {
			if tx != nil {
				tx.Rollback() //?rollback jika ada error saan runtime
			}

			debug.PrintStack() //?agar ada stock tracenya yang digunakan utk debug
			panic(e)           //?agar bisa nyampai ke Middleware
		}
	}()

	defer func() {
		if err != nil && tx != nil {
			tx.Rollback() //?rollback jika ada error
		}
	}()

	storeRepo := ss.storeRepository.WithTransaction(tx)

	// *update delete_at & delete_by ke DB

	err = storeRepo.DeleteStore(ctx, request.Id, time.Now(), claims.FullName)
	if err != nil {
		return nil, err
	}

	// *jika ada image, hapus image
	if storeEntity.ImageFileName != "" {
		imagePath := filepath.Join("storage", storeEntity.Id, "store", storeEntity.ImageFileName)
		err = os.Remove(imagePath)
		if err != nil {
			return nil, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	// *success
	return &store.DeleteStoreResponse{
		Base: utils.SuccessResponse("Delete with SoftDelete Store Success"),
	}, nil
}

func NewStoreService(db *sql.DB, storeRepository repository.IStoreRepository) IStoreService {
	return &storeService{
		db:              db,
		storeRepository: storeRepository,
	}
}
