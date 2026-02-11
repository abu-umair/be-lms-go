package service

import (
	"context"
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
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
)

type ICourseService interface {
	CreateCourse(ctx context.Context, request *course.CreateCourseRequest) (*course.CreateCourseResponse, error)
	DetailCourse(ctx context.Context, request *course.DetailCourseRequest) (*course.DetailCourseResponse, error)
	EditCourse(ctx context.Context, request *course.EditCourseRequest) (*course.EditCourseResponse, error)
	DeleteCourse(ctx context.Context, request *course.DeleteCourseRequest) (*course.DeleteCourseResponse, error)
}

type courseService struct {
	db               *sqlx.DB
	courseRepository repository.ICourseRepository
}

func (ss *courseService) CreateCourse(ctx context.Context, request *course.CreateCourseRequest) (*course.CreateCourseResponse, error) {
	//* Get data token
	claims, err := jwtentity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	//* apakah role user adl Instructor
	if claims.Role != entity.UserRoleInstructor {
		return nil, utils.UnauthenticatedResponse()
	}

	tx, err := ss.db.BeginTxx(ctx, nil)
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
	var priceDecimal *decimal.Decimal
	if request.Price != nil {
		// Konversi string ke decimal
		d, err := decimal.NewFromString(*request.Price)
		if err != nil {
			return &course.CreateCourseResponse{
				Base: utils.BadRequestResponse("Invalid price format"),
			}, nil
		}
		priceDecimal = &d
	}

	var discountDecimal *decimal.Decimal
	if request.Price != nil {
		d, err := decimal.NewFromString(*request.Discount)
		if err != nil {
			return &course.CreateCourseResponse{
				Base: utils.BadRequestResponse("Invalid discount format"),
			}, nil
		}
		discountDecimal = &d
	}

	courseEntity := entity.Course{
		Id:                 request.Id,
		Name:               request.Name,
		Address:            request.Address,
		ImageFileName:      request.ImageFileName,
		Slug:               request.Slug,
		InstructorId:       request.InstructorId,
		CategoryId:         request.CategoryId,
		CourseType:         request.CourseType,
		SeoDescription:     request.SeoDescription,
		Duration:           request.Duration,
		Timezone:           request.Timezone,
		Thumbnail:          request.Thumbnail,
		DemoVideoStorage:   request.DemoVideoStorage,
		DemoVideoSource:    request.DemoVideoSource,
		Description:        request.Description,
		Capacity:           request.Capacity,
		Price:              priceDecimal,
		Discount:           discountDecimal,
		Certificate:        request.Certificate,
		Gna:                request.Gna,
		MessageForReviewer: request.MessageForReviewer,
		IsApproved:         request.IsApproved,
		CourseLevelId:      request.CourseLevelId,
		CourseLanguageId:   request.CourseLanguageId,

		CreatedAt: time.Now(),
		CreatedBy: claims.FullName,
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

	//* apakah role user adl Instructor
	if claims.Role != entity.UserRoleInstructor {
		return nil, utils.UnauthenticatedResponse()
	}

	// * Get course by course_id
	// Misal request.FieldMask.Paths berisi ["name", "address"]
	paths := []string{"id"} // ID wajib ada untuk mapping
	if request.FieldMask != nil {
		paths = append(paths, request.FieldMask.Paths...)
	}

	courseEntity, err := ss.courseRepository.GetCourseByIdFieldMask(ctx, request.Id, paths)
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
	res := &course.DetailCourseResponse{
		Base: utils.SuccessResponse("Course Detail Success"),
		Id:   courseEntity.Id,
	}

	// Cek Name: Jika kosong (tidak di-select), res.Name tetap nil (tidak muncul di JSON)
	//? Mapping Field String Biasa (Non-Pointer di Struct)
	res.Name = utils.StringToPtr(courseEntity.Name)
	res.Address = utils.PtrStringToPtr(courseEntity.Address)
	res.CreatedBy = utils.StringToPtr(courseEntity.CreatedBy)

	//?Mapping Field Pointer String (*string di Struct)
	res.Slug = utils.PtrStringToPtr(courseEntity.Slug)
	res.InstructorId = utils.PtrStringToPtr(courseEntity.InstructorId)
	res.CategoryId = utils.PtrStringToPtr(courseEntity.CategoryId)
	res.CourseType = utils.PtrStringToPtr(courseEntity.CourseType)
	res.SeoDescription = utils.PtrStringToPtr(courseEntity.SeoDescription)
	res.Duration = utils.PtrStringToPtr(courseEntity.Duration)
	res.Timezone = utils.PtrStringToPtr(courseEntity.Timezone)
	res.Thumbnail = utils.PtrStringToPtr(courseEntity.Thumbnail)
	res.DemoVideoStorage = utils.PtrStringToPtr(courseEntity.DemoVideoStorage)
	res.DemoVideoSource = utils.PtrStringToPtr(courseEntity.DemoVideoSource)
	res.Description = utils.PtrStringToPtr(courseEntity.Description)
	res.Certificate = utils.PtrStringToPtr(courseEntity.Certificate)
	res.Gna = utils.PtrStringToPtr(courseEntity.Gna)
	res.MessageForReviewer = utils.PtrStringToPtr(courseEntity.MessageForReviewer)
	res.IsApproved = utils.PtrStringToPtr(courseEntity.IsApproved)
	res.Status = utils.PtrStringToPtr(courseEntity.Status)
	res.CourseLevelId = utils.PtrStringToPtr(courseEntity.CourseLevelId)
	res.CourseLanguageId = utils.PtrStringToPtr(courseEntity.CourseLanguageId)
	res.UpdatedBy = utils.PtrStringToPtr(courseEntity.UpdatedBy)
	res.DeletedBy = utils.PtrStringToPtr(courseEntity.DeletedBy)

	//? Mapping Angka dan Harga (Int64 & Decimal)
	res.Capacity = utils.PtrInt32ToPtr(courseEntity.Capacity)
	res.Price = utils.PtrDecimalToPtr(courseEntity.Price)
	res.Discount = utils.PtrDecimalToPtr(courseEntity.Discount)

	//? Mapping Waktu (Time)
	res.CreatedAt = utils.TimeToPtr(courseEntity.CreatedAt)
	res.UpdatedAt = utils.TimeToPtr(courseEntity.UpdatedAt)
	res.DeletedAt = utils.PtrTimeToPtr(courseEntity.DeletedAt)

	//? khusus image: Cek dulu apakah ImageFileName ada di database
	if courseEntity.ImageFileName != "" {
		fullUrl := fmt.Sprintf("%s/%s/course/%s", os.Getenv("STORAGE_SERVICE_URL"), courseEntity.Id, courseEntity.ImageFileName)
		res.ImageFileName = &fullUrl
	}

	return res, nil
}

func (ss *courseService) EditCourse(ctx context.Context, request *course.EditCourseRequest) (*course.EditCourseResponse, error) {
	//* Get data token
	claims, err := jwtentity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	//* apakah role user adl Instructor
	if claims.Role != entity.UserRoleInstructor {
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

	tx, err := ss.db.BeginTxx(ctx, nil)
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
	var priceDecimal *decimal.Decimal
	if request.Price != nil {
		// Konversi string ke decimal
		d, err := decimal.NewFromString(*request.Price)
		if err != nil {
			return &course.EditCourseResponse{
				Base: utils.BadRequestResponse("Invalid price format"),
			}, nil
		}
		priceDecimal = &d
	}

	var discountDecimal *decimal.Decimal
	if request.Price != nil {
		d, err := decimal.NewFromString(*request.Discount)
		if err != nil {
			return &course.EditCourseResponse{
				Base: utils.BadRequestResponse("Invalid discount format"),
			}, nil
		}
		discountDecimal = &d
	}

	newCourse := entity.Course{
		Id:                 request.Id,
		Name:               request.Name,
		Address:            request.Address,
		ImageFileName:      request.ImageFileName,
		Slug:               request.Slug,
		InstructorId:       request.InstructorId,
		CategoryId:         request.CategoryId,
		CourseType:         request.CourseType,
		SeoDescription:     request.SeoDescription,
		Duration:           request.Duration,
		Timezone:           request.Timezone,
		Thumbnail:          request.Thumbnail,
		DemoVideoStorage:   request.DemoVideoStorage,
		DemoVideoSource:    request.DemoVideoSource,
		Description:        request.Description,
		Capacity:           request.Capacity,
		Price:              priceDecimal,
		Discount:           discountDecimal,
		Certificate:        request.Certificate,
		Gna:                request.Gna,
		MessageForReviewer: request.MessageForReviewer,
		IsApproved:         request.IsApproved,
		Status:             request.Status,
		CourseLevelId:      request.CourseLevelId,
		CourseLanguageId:   request.CourseLanguageId,

		UpdatedAt: time.Now(),
		UpdatedBy: &claims.FullName,
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

	//* apakah role user adl Instructor
	if claims.Role != entity.UserRoleInstructor {
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

	tx, err := ss.db.BeginTxx(ctx, nil)
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

func NewCourseService(db *sqlx.DB, courseRepository repository.ICourseRepository) ICourseService {
	return &courseService{
		db:               db,
		courseRepository: courseRepository,
	}
}
