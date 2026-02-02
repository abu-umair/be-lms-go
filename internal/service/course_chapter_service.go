package service

import (
	"context"
	"runtime/debug"
	"time"

	"github.com/abu-umair/be-lms-go/internal/entity"
	jwtentity "github.com/abu-umair/be-lms-go/internal/entity/jwt"
	"github.com/abu-umair/be-lms-go/internal/repository"
	"github.com/abu-umair/be-lms-go/internal/utils"
	"github.com/abu-umair/be-lms-go/pb/course_chapter"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ICourseChapterService interface {
	CreateCourseChapter(ctx context.Context, request *course_chapter.CreateCourseChapterRequest) (*course_chapter.CreateCourseChapterResponse, error)
	// DetailCourse(ctx context.Context, request *course_chapter.DetailCourseChapterRequest) (*course_chapter.DetailCourseChapterResponse, error)
	// EditCourse(ctx context.Context, request *course_chapter.EditCourseChapterRequest) (*course_chapter.EditCourseChapterResponse, error)
	// DeleteCourse(ctx context.Context, request *course_chapter.DeleteCourseChapterRequest) (*course_chapter.DeleteCourseChapterResponse, error)
}

type courseChapterService struct {
	db                      *sqlx.DB
	courseChapterRepository repository.ICourseChapterRepository
}

func (cs *courseChapterService) CreateCourseChapter(ctx context.Context, request *course_chapter.CreateCourseChapterRequest) (*course_chapter.CreateCourseChapterResponse, error) {
	//* Get data token
	claims, err := jwtentity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	//* apakah role user adl Owner
	if claims.Role != entity.UserRoleOwner {
		return nil, utils.UnauthenticatedResponse()
	}

	tx, err := cs.db.BeginTxx(ctx, nil)
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

	courseChapterRepo := cs.courseChapterRepository.WithTransaction(tx)

	// *insert ke DB
	courseChapterEntity := entity.CourseChapter{
		Id:           uuid.NewString(),
		InstructorId: request.InstructorId,
		CourseId:     request.CourseId,
		Title:        request.Title,
		Order:        request.Order,
		Status:       request.Status,

		CreatedAt: time.Now(),
		CreatedBy: claims.FullName,
	}

	err = courseChapterRepo.CreateNewCourseChapter(ctx, &courseChapterEntity)
	if err != nil {
		return nil, err
	}

	err = tx.Commit() //?harus dicommit agar data tersimpan
	if err != nil {
		return nil, err
	}

	// *success
	return &course_chapter.CreateCourseChapterResponse{
		Base: utils.SuccessResponse("Course chapter successfully created"),
	}, nil
}

// func (ss *courseChapterService) DetailCourse(ctx context.Context, request *course_chapter.DetailCourseChapterRequest) (*course_chapter.DetailCourseChapterResponse, error) {

// 	//* Get data token
// 	claims, err := jwtentity.GetClaimsFromContext(ctx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	//* apakah role user adl Owner
// 	if claims.Role != entity.UserRoleOwner {
// 		return nil, utils.UnauthenticatedResponse()
// 	}

// 	// * Get course by course_id
// 	// Misal request.FieldMask.Paths berisi ["name", "address"]
// 	paths := []string{"id"} // ID wajib ada untuk mapping
// 	if request.FieldMask != nil {
// 		paths = append(paths, request.FieldMask.Paths...)
// 	}

// 	courseEntity, err := cs.courseRepository.GetCourseByIdFieldMask(ctx, request.Id, paths)
// 	if err != nil {
// 		return nil, err
// 	}

// 	//* Apabila null course_id, return not found
// 	if courseEntity == nil {
// 		return &course.DetailCourseChapterResponse{
// 			Base: utils.NotFoundResponse("Course not found"),
// 		}, nil
// 	}

// 	// *success
// 	res := &course.DetailCourseChapterResponse{
// 		Base: utils.SuccessResponse("Course Detail Success"),
// 		Id:   courseEntity.Id,
// 	}

// 	// Cek Name: Jika kosong (tidak di-select), res.Name tetap nil (tidak muncul di JSON)
// 	//? Mapping Field String Biasa (Non-Pointer di Struct)
// 	res.Name = utils.StringToPtr(courseEntity.Name)
// 	res.Address = utils.StringToPtr(courseEntity.Address)
// 	res.CreatedBy = utils.StringToPtr(courseEntity.CreatedBy)

// 	//?Mapping Field Pointer String (*string di Struct)
// 	res.Slug = utils.PtrStringToPtr(courseEntity.Slug)
// 	res.InstructorId = utils.PtrStringToPtr(courseEntity.InstructorId)
// 	res.CategoryId = utils.PtrStringToPtr(courseEntity.CategoryId)
// 	res.CourseType = utils.PtrStringToPtr(courseEntity.CourseType)
// 	res.SeoDescription = utils.PtrStringToPtr(courseEntity.SeoDescription)
// 	res.Duration = utils.PtrStringToPtr(courseEntity.Duration)
// 	res.Timezone = utils.PtrStringToPtr(courseEntity.Timezone)
// 	res.Thumbnail = utils.PtrStringToPtr(courseEntity.Thumbnail)
// 	res.DemoVideoStorage = utils.PtrStringToPtr(courseEntity.DemoVideoStorage)
// 	res.DemoVideoSource = utils.PtrStringToPtr(courseEntity.DemoVideoSource)
// 	res.Description = utils.PtrStringToPtr(courseEntity.Description)
// 	res.Certificate = utils.PtrStringToPtr(courseEntity.Certificate)
// 	res.Gna = utils.PtrStringToPtr(courseEntity.Gna)
// 	res.MessageForReviewer = utils.PtrStringToPtr(courseEntity.MessageForReviewer)
// 	res.IsApproved = utils.PtrStringToPtr(courseEntity.IsApproved)
// 	res.Status = utils.PtrStringToPtr(courseEntity.Status)
// 	res.CourseLevelId = utils.PtrStringToPtr(courseEntity.CourseLevelId)
// 	res.CourseLanguageId = utils.PtrStringToPtr(courseEntity.CourseLanguageId)
// 	res.UpdatedBy = utils.PtrStringToPtr(courseEntity.UpdatedBy)
// 	res.DeletedBy = utils.PtrStringToPtr(courseEntity.DeletedBy)

// 	//? Mapping Angka dan Harga (Int64 & Decimal)
// 	res.Capacity = utils.PtrInt64ToPtr(courseEntity.Capacity)
// 	res.Price = utils.PtrDecimalToPtr(courseEntity.Price)
// 	res.Discount = utils.PtrDecimalToPtr(courseEntity.Discount)

// 	//? Mapping Waktu (Time)
// 	res.CreatedAt = utils.TimeToPtr(courseEntity.CreatedAt)
// 	res.UpdatedAt = utils.TimeToPtr(courseEntity.UpdatedAt)
// 	res.DeletedAt = utils.PtrTimeToPtr(courseEntity.DeletedAt)

// 	//? khusus image: Cek dulu apakah ImageFileName ada di database
// 	if courseEntity.ImageFileName != "" {
// 		fullUrl := fmt.Sprintf("%s/%s/course/%s", os.Getenv("STORAGE_SERVICE_URL"), courseEntity.Id, courseEntity.ImageFileName)
// 		res.ImageFileName = &fullUrl
// 	}

// 	return res, nil
// }

// func (ss *courseChapterService) EditCourse(ctx context.Context, request *course_chapter.EditCourseChapterRequest) (*course_chapter.EditCourseChapterResponse, error) {
// 	//* Get data token
// 	claims, err := jwtentity.GetClaimsFromContext(ctx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	//* apakah role user adl Owner
// 	if claims.Role != entity.UserRoleOwner {
// 		return nil, utils.UnauthenticatedResponse()
// 	}

// 	// *Apakah Id course ada di DB
// 	courseEntity, err := ss.courseRepository.GetCourseById(ctx, request.Id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if courseEntity == nil {
// 		return &course.EditCourseChapterResponse{
// 			Base: utils.NotFoundResponse("Course not found"),
// 		}, nil
// 	}

// 	tx, err := ss.db.BeginTxx(ctx, nil)
// 	if err != nil {
// 		return nil, err
// 	}

// 	defer func() {
// 		if e := recover(); e != nil {
// 			if tx != nil {
// 				tx.Rollback() //?rollback jika ada error saan runtime
// 			}

// 			debug.PrintStack() //?agar ada stock tracenya yang digunakan utk debug
// 			panic(e)           //?agar bisa nyampai ke Middleware
// 		}
// 	}()

// 	defer func() {
// 		if err != nil && tx != nil {
// 			tx.Rollback() //?rollback jika ada error
// 		}
// 	}()

// 	courseRepo := ss.courseRepository.WithTransaction(tx)

// 	// *update ke DB
// 	var priceDecimal *decimal.Decimal
// 	if request.Price != nil {
// 		// Konversi string ke decimal
// 		d, err := decimal.NewFromString(*request.Price)
// 		if err != nil {
// 			return &course.EditCourseChapterResponse{
// 				Base: utils.BadRequestResponse("Invalid price format"),
// 			}, nil
// 		}
// 		priceDecimal = &d
// 	}

// 	var discountDecimal *decimal.Decimal
// 	if request.Price != nil {
// 		d, err := decimal.NewFromString(*request.Discount)
// 		if err != nil {
// 			return &course.EditCourseChapterResponse{
// 				Base: utils.BadRequestResponse("Invalid discount format"),
// 			}, nil
// 		}
// 		discountDecimal = &d
// 	}

// 	newCourse := entity.Course{
// 		Id:                 request.Id,
// 		Name:               request.Name,
// 		Address:            request.Address,
// 		ImageFileName:      request.ImageFileName,
// 		Slug:               request.Slug,
// 		InstructorId:       request.InstructorId,
// 		CategoryId:         request.CategoryId,
// 		CourseType:         request.CourseType,
// 		SeoDescription:     request.SeoDescription,
// 		Duration:           request.Duration,
// 		Timezone:           request.Timezone,
// 		Thumbnail:          request.Thumbnail,
// 		DemoVideoStorage:   request.DemoVideoStorage,
// 		DemoVideoSource:    request.DemoVideoSource,
// 		Description:        request.Description,
// 		Capacity:           request.Capacity,
// 		Price:              priceDecimal,
// 		Discount:           discountDecimal,
// 		Certificate:        request.Certificate,
// 		Gna:                request.Gna,
// 		MessageForReviewer: request.MessageForReviewer,
// 		IsApproved:         request.IsApproved,
// 		Status:             request.Status,
// 		CourseLevelId:      request.CourseLevelId,
// 		CourseLanguageId:   request.CourseLanguageId,

// 		UpdatedAt: time.Now(),
// 		UpdatedBy: &claims.FullName,
// 	}

// 	err = courseRepo.UpdateCourse(ctx, &newCourse)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// *jika ada image baru, hapus image lama
// 	if courseEntity.ImageFileName != request.ImageFileName {
// 		newImagePath := filepath.Join("storage", request.Id, "course", request.ImageFileName)
// 		_, err := os.Stat(newImagePath)
// 		if err != nil {
// 			if os.IsNotExist(err) {
// 				return &course.EditCourseChapterResponse{
// 					Base: utils.BadRequestResponse("Image not found"),
// 				}, nil
// 			}
// 			return nil, err
// 		}

// 		oldImagePath := filepath.Join("storage", courseEntity.Id, "course", courseEntity.ImageFileName)
// 		err = os.Remove(oldImagePath)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}

// 	err = tx.Commit()
// 	if err != nil {
// 		return nil, err
// 	}

// 	// *success
// 	return &course.EditCourseChapterResponse{
// 		Base: utils.SuccessResponse("Edit Course Success"),
// 		Id:   request.Id,
// 	}, nil
// }

// func (ss *courseChapterService) DeleteCourse(ctx context.Context, request *course_chapter.DeleteCourseChapterRequest) (*course_chapter.DeleteCourseChapterResponse, error) {
// 	//* Get data token
// 	claims, err := jwtentity.GetClaimsFromContext(ctx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	//* apakah role user adl Owner
// 	if claims.Role != entity.UserRoleOwner {
// 		return nil, utils.UnauthenticatedResponse()
// 	}

// 	// *Apakah Id course ada di DB
// 	courseEntity, err := ss.courseRepository.GetCourseById(ctx, request.Id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if courseEntity == nil {
// 		return &course.DeleteCourseChapterResponse{
// 			Base: utils.NotFoundResponse("Course not found"),
// 		}, nil
// 	}

// 	tx, err := ss.db.BeginTxx(ctx, nil)
// 	if err != nil {
// 		return nil, err
// 	}

// 	defer func() {
// 		if e := recover(); e != nil {
// 			if tx != nil {
// 				tx.Rollback() //?rollback jika ada error saan runtime
// 			}

// 			debug.PrintStack() //?agar ada stock tracenya yang digunakan utk debug
// 			panic(e)           //?agar bisa nyampai ke Middleware
// 		}
// 	}()

// 	defer func() {
// 		if err != nil && tx != nil {
// 			tx.Rollback() //?rollback jika ada error
// 		}
// 	}()

// 	courseRepo := ss.courseRepository.WithTransaction(tx)

// 	// *update delete_at & delete_by ke DB

// 	err = courseRepo.DeleteCourse(ctx, request.Id, time.Now(), claims.FullName)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// *jika ada image, hapus image
// 	if courseEntity.ImageFileName != "" {
// 		imagePath := filepath.Join("storage", courseEntity.Id, "course", courseEntity.ImageFileName)
// 		err = os.Remove(imagePath)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}

// 	err = tx.Commit()
// 	if err != nil {
// 		return nil, err
// 	}

// 	// *success
// 	return &course.DeleteCourseChapterResponse{
// 		Base: utils.SuccessResponse("Delete with SoftDelete Course Success"),
// 	}, nil
// }

func NewCourseChapterService(db *sqlx.DB, courseChapterRepository repository.ICourseChapterRepository) ICourseChapterService {
	return &courseChapterService{
		db:                      db,
		courseChapterRepository: courseChapterRepository,
	}
}
