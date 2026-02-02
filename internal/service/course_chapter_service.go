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
	DetailCourseChapter(ctx context.Context, request *course_chapter.DetailCourseChapterRequest) (*course_chapter.DetailCourseChapterResponse, error)
	EditCourseChapter(ctx context.Context, request *course_chapter.EditCourseChapterRequest) (*course_chapter.EditCourseChapterResponse, error)
	DeleteCourseChapter(ctx context.Context, request *course_chapter.DeleteCourseChapterRequest) (*course_chapter.DeleteCourseChapterResponse, error)
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
		OrderChapter: request.OrderChapter,
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

func (cs *courseChapterService) DetailCourseChapter(ctx context.Context, request *course_chapter.DetailCourseChapterRequest) (*course_chapter.DetailCourseChapterResponse, error) {

	//* Get data token
	claims, err := jwtentity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	//* apakah role user adl Owner
	if claims.Role != entity.UserRoleOwner {
		return nil, utils.UnauthenticatedResponse()
	}

	// * Get course_chapters by chapter_id
	// Misal request.FieldMask.Paths berisi ["name", "address"]
	paths := []string{"id"} // ID wajib ada untuk mapping
	if request.FieldMask != nil {
		paths = append(paths, request.FieldMask.Paths...)
	}

	courseChapterEntity, err := cs.courseChapterRepository.GetCourseChapterByIdFieldMask(ctx, request.Id, paths)
	if err != nil {
		return nil, err
	}

	//* Apabila null chapter_id, return not found
	if courseChapterEntity == nil {
		return &course_chapter.DetailCourseChapterResponse{
			Base: utils.NotFoundResponse("Course chapter not found"),
		}, nil
	}

	// *success
	res := &course_chapter.DetailCourseChapterResponse{
		Base: utils.SuccessResponse("Course Chapter Detail Success"),
		Id:   courseChapterEntity.Id,
	}

	// Cek Name: Jika kosong (tidak di-select), res.Name tetap nil (tidak muncul di JSON)
	//? Mapping Field String Biasa (Non-Pointer di Struct)
	res.Title = utils.StringToPtr(courseChapterEntity.Title)
	res.InstructorId = utils.StringToPtr(courseChapterEntity.InstructorId)
	res.CourseId = utils.StringToPtr(courseChapterEntity.CourseId)
	res.Status = utils.StringToPtr(courseChapterEntity.Status)
	res.CreatedBy = utils.StringToPtr(courseChapterEntity.CreatedBy)

	//?Mapping Field Pointer String (*string di Struct)
	res.UpdatedBy = utils.PtrStringToPtr(courseChapterEntity.UpdatedBy)
	res.DeletedBy = utils.PtrStringToPtr(courseChapterEntity.DeletedBy)

	//? Mapping Angka int64 biasa
	res.OrderChapter = utils.Int64ToPtr(courseChapterEntity.OrderChapter)

	//? Mapping Waktu (Time)
	res.CreatedAt = utils.TimeToPtr(courseChapterEntity.CreatedAt)
	res.UpdatedAt = utils.TimeToPtr(courseChapterEntity.UpdatedAt)
	res.DeletedAt = utils.PtrTimeToPtr(courseChapterEntity.DeletedAt)

	return res, nil
}

func (cs *courseChapterService) EditCourseChapter(ctx context.Context, request *course_chapter.EditCourseChapterRequest) (*course_chapter.EditCourseChapterResponse, error) {
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
	courseEntity, err := cs.courseChapterRepository.GetCourseChapterById(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	if courseEntity == nil {
		return &course_chapter.EditCourseChapterResponse{
			Base: utils.NotFoundResponse("Course chapter not found"),
		}, nil
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

	// *update ke DB
	newCourse := entity.CourseChapter{
		Id:           request.Id,
		InstructorId: request.InstructorId,
		CourseId:     request.CourseId,
		Title:        request.Title,
		OrderChapter: request.OrderChapter,
		Status:       request.Status,

		UpdatedAt: time.Now(),
		UpdatedBy: &claims.FullName,
	}

	err = courseChapterRepo.UpdateCourseChapter(ctx, &newCourse)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	// *success
	return &course_chapter.EditCourseChapterResponse{
		Base: utils.SuccessResponse("Edit Course Chapter Success"),
		Id:   request.Id,
	}, nil
}

func (cs *courseChapterService) DeleteCourseChapter(ctx context.Context, request *course_chapter.DeleteCourseChapterRequest) (*course_chapter.DeleteCourseChapterResponse, error) {
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
	courseChapterEntity, err := cs.courseChapterRepository.GetCourseChapterById(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	if courseChapterEntity == nil {
		return &course_chapter.DeleteCourseChapterResponse{
			Base: utils.NotFoundResponse("Course Chapter not found"),
		}, nil
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

	// *update delete_at & delete_by ke DB

	err = courseChapterRepo.DeleteCourseChapter(ctx, request.Id, time.Now(), claims.FullName)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	// *success
	return &course_chapter.DeleteCourseChapterResponse{
		Base: utils.SuccessResponse("Delete with SoftDelete Course Success"),
	}, nil
}

func NewCourseChapterService(db *sqlx.DB, courseChapterRepository repository.ICourseChapterRepository) ICourseChapterService {
	return &courseChapterService{
		db:                      db,
		courseChapterRepository: courseChapterRepository,
	}
}
