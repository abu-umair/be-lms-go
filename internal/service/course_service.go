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
	"github.com/abu-umair/be-lms-go/pb/course"
)

type ICourseService interface {
	CreateCourse(ctx context.Context, request *course.CreateCourseRequest) (*course.CreateCourseResponse, error)
	DetailCourse(ctx context.Context, request *course.DetailCourseRequest) (*course.DetailCourseResponse, error)
	EditCourse(ctx context.Context, request *course.EditCourseRequest) (*course.EditCourseResponse, error)
	DeleteCourse(ctx context.Context, request *course.DeleteCourseRequest) (*course.DeleteCourseResponse, error)
}

type courseService struct {
	db               *sql.DB
	courseRepository repository.ICourseRepository
}

func (ss *courseService) CreateCourse(ctx context.Context, request *course.CreateCourseRequest) (*course.CreateCourseResponse, error) {
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

	courseRepo := ss.courseRepository.WithTransaction(tx)

	// *insert ke DB
	courseEntity := entity.Course{
		Id:            request.Id,
		Name:          request.Name,
		Address:       request.Address,
		ImageFileName: request.ImageFileName,
		CreatedAt:     time.Now(),
		CreatedBy:     claims.FullName,
	}

	err = courseRepo.CreateNewCourse(ctx, &courseEntity)
	if err != nil {
		return nil, err
	}

	// *apakah image ada
	imagePath := filepath.Join("storage", courseEntity.Id, "course", request.ImageFileName)
	_, err = os.Stat(imagePath)
	if err != nil {
		if os.IsNotExist(err) {
			return &course.CreateCourseResponse{
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
	return &course.CreateCourseResponse{
		Base: utils.SuccessResponse("Course successfully created"),
	}, nil
}

func (ss *courseService) DetailCourse(ctx context.Context, request *course.DetailCourseRequest) (*course.DetailCourseResponse, error) {

	//* Get data token
	claims, err := jwtentity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	//* apakah role user adl Owner
	if claims.Role != entity.UserRoleOwner {
		return nil, utils.UnauthenticatedResponse()
	}

	// * Get course by course_id
	courseEntity, err := ss.courseRepository.GetCourseById(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	//* Apabila null course_id, return not found
	if courseEntity == nil {
		return &course.DetailCourseResponse{
			Base: utils.NotFoundResponse("Course not found"),
		}, nil
	}

	// *success
	return &course.DetailCourseResponse{
		Base:          utils.SuccessResponse("Course Detail Success"),
		Id:            courseEntity.Id,
		Name:          courseEntity.Name,
		Address:       courseEntity.Address,
		ImageFileName: fmt.Sprintf("%s/%s/course/%s", os.Getenv("STORAGE_SERVICE_URL"), courseEntity.Id, courseEntity.ImageFileName),
	}, nil
}

func (ss *courseService) EditCourse(ctx context.Context, request *course.EditCourseRequest) (*course.EditCourseResponse, error) {
	//* Get data token
	claims, err := jwtentity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	//* apakah role user adl Owner
	if claims.Role != entity.UserRoleOwner {
		return nil, utils.UnauthenticatedResponse()
	}

	// *Apakah Id course ada di DB
	courseEntity, err := ss.courseRepository.GetCourseById(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	if courseEntity == nil {
		return &course.EditCourseResponse{
			Base: utils.NotFoundResponse("Course not found"),
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

	courseRepo := ss.courseRepository.WithTransaction(tx)

	// *update ke DB
	newCourse := entity.Course{
		Id:            request.Id,
		Name:          request.Name,
		Address:       request.Address,
		ImageFileName: request.ImageFileName,
		UpdatedAt:     time.Now(),
		UpdatedBy:     &claims.FullName,
	}

	err = courseRepo.UpdateCourse(ctx, &newCourse)
	if err != nil {
		return nil, err
	}

	// *jika ada image baru, hapus image lama
	if courseEntity.ImageFileName != request.ImageFileName {
		newImagePath := filepath.Join("storage", request.Id, "course", request.ImageFileName)
		_, err := os.Stat(newImagePath)
		if err != nil {
			if os.IsNotExist(err) {
				return &course.EditCourseResponse{
					Base: utils.BadRequestResponse("Image not found"),
				}, nil
			}
			return nil, err
		}

		oldImagePath := filepath.Join("storage", courseEntity.Id, "course", courseEntity.ImageFileName)
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
	return &course.EditCourseResponse{
		Base: utils.SuccessResponse("Edit Course Success"),
		Id:   request.Id,
	}, nil
}

func (ss *courseService) DeleteCourse(ctx context.Context, request *course.DeleteCourseRequest) (*course.DeleteCourseResponse, error) {
	//* Get data token
	claims, err := jwtentity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	//* apakah role user adl Owner
	if claims.Role != entity.UserRoleOwner {
		return nil, utils.UnauthenticatedResponse()
	}

	// *Apakah Id course ada di DB
	courseEntity, err := ss.courseRepository.GetCourseById(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	if courseEntity == nil {
		return &course.DeleteCourseResponse{
			Base: utils.NotFoundResponse("Course not found"),
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

	courseRepo := ss.courseRepository.WithTransaction(tx)

	// *update delete_at & delete_by ke DB

	err = courseRepo.DeleteCourse(ctx, request.Id, time.Now(), claims.FullName)
	if err != nil {
		return nil, err
	}

	// *jika ada image, hapus image
	if courseEntity.ImageFileName != "" {
		imagePath := filepath.Join("storage", courseEntity.Id, "course", courseEntity.ImageFileName)
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
	return &course.DeleteCourseResponse{
		Base: utils.SuccessResponse("Delete with SoftDelete Course Success"),
	}, nil
}

func NewCourseService(db *sql.DB, courseRepository repository.ICourseRepository) ICourseService {
	return &courseService{
		db:               db,
		courseRepository: courseRepository,
	}
}
