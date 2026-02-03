package service

import (
	"context"
	"runtime/debug"
	"time"

	"github.com/abu-umair/be-lms-go/internal/entity"
	jwtentity "github.com/abu-umair/be-lms-go/internal/entity/jwt"
	"github.com/abu-umair/be-lms-go/internal/repository"
	"github.com/abu-umair/be-lms-go/internal/utils"
	"github.com/abu-umair/be-lms-go/pb/chapter_lesson"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type IChapterLessonService interface {
	CreateChapterLesson(ctx context.Context, request *chapter_lesson.CreateChapterLessonRequest) (*chapter_lesson.CreateChapterLessonResponse, error)
	// DetailChapterLesson(ctx context.Context, request *chapter_lesson.DetailChapterLessonRequest) (*chapter_lesson.DetailChapterLessonResponse, error)
	// EditChapterLesson(ctx context.Context, request *chapter_lesson.EditChapterLessonRequest) (*chapter_lesson.EditChapterLessonResponse, error)
	// DeleteChapterLesson(ctx context.Context, request *chapter_lesson.DeleteChapterLessonRequest) (*chapter_lesson.DeleteChapterLessonResponse, error)
}

type chapterLessonService struct {
	db                      *sqlx.DB
	chapterLessonRepository repository.IChapterLessonRepository
}

func (ls *chapterLessonService) CreateChapterLesson(ctx context.Context, request *chapter_lesson.CreateChapterLessonRequest) (*chapter_lesson.CreateChapterLessonResponse, error) {
	//* Get data token
	claims, err := jwtentity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	//* apakah role user adl Owner
	if claims.Role != entity.UserRoleOwner {
		return nil, utils.UnauthenticatedResponse()
	}

	tx, err := ls.db.BeginTxx(ctx, nil)
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

	chapterLessonRepo := ls.chapterLessonRepository.WithTransaction(tx)

	// *insert ke DB
	chapterLessonEntity := entity.ChapterLesson{
		Id:            uuid.NewString(),
		InstructorId:  request.InstructorId,
		CourseId:      request.CourseId,
		Title:         request.Title,
		OrderLesson:   request.OrderLesson,
		ChapterId:     request.ChapterId,
		Slug:          request.Slug,
		Description:   request.Description,
		FilePath:      request.FilePath,
		StorageLesson: request.StorageLesson,
		LessonType:    request.LessonType,
		Volume:        request.Volume,
		Duration:      request.Duration,
		FileType:      request.FileType,
		Downloadable:  request.Downloadable,
		IsPreview:     request.IsPreview,
		Status:        request.Status,

		CreatedAt: time.Now(),
		CreatedBy: claims.FullName,
	}

	err = chapterLessonRepo.CreateNewChapterLesson(ctx, &chapterLessonEntity)
	if err != nil {
		return nil, err
	}

	err = tx.Commit() //?harus dicommit agar data tersimpan
	if err != nil {
		return nil, err
	}

	// *success
	return &chapter_lesson.CreateChapterLessonResponse{
		Base: utils.SuccessResponse("Course chapter lesson successfully created"),
	}, nil
}

// func (cs *chapterLessonService) DetailChapterLesson(ctx context.Context, request *chapter_lesson.DetailChapterLessonRequest) (*chapter_lesson.DetailChapterLessonResponse, error) {

// 	//* Get data token
// 	claims, err := jwtentity.GetClaimsFromContext(ctx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	//* apakah role user adl Owner
// 	if claims.Role != entity.UserRoleOwner {
// 		return nil, utils.UnauthenticatedResponse()
// 	}

// 	// * Get chapter_lessons by chapter_id
// 	// Misal request.FieldMask.Paths berisi ["name", "address"]
// 	paths := []string{"id"} // ID wajib ada untuk mapping
// 	if request.FieldMask != nil {
// 		paths = append(paths, request.FieldMask.Paths...)
// 	}

// 	chapterLessonEntity, err := cs.chapterLessonRepository.GetChapterLessonByIdFieldMask(ctx, request.Id, paths)
// 	if err != nil {
// 		return nil, err
// 	}

// 	//* Apabila null chapter_id, return not found
// 	if chapterLessonEntity == nil {
// 		return &chapter_lesson.DetailChapterLessonResponse{
// 			Base: utils.NotFoundResponse("Course chapter not found"),
// 		}, nil
// 	}

// 	// *success
// 	res := &chapter_lesson.DetailChapterLessonResponse{
// 		Base: utils.SuccessResponse("Course Chapter Detail Success"),
// 		Id:   chapterLessonEntity.Id,
// 	}

// 	// Cek Name: Jika kosong (tidak di-select), res.Name tetap nil (tidak muncul di JSON)
// 	//? Mapping Field String Biasa (Non-Pointer di Struct)
// 	res.Title = utils.StringToPtr(chapterLessonEntity.Title)
// 	res.InstructorId = utils.StringToPtr(chapterLessonEntity.InstructorId)
// 	res.CourseId = utils.StringToPtr(chapterLessonEntity.CourseId)
// 	res.Status = utils.StringToPtr(chapterLessonEntity.Status)
// 	res.CreatedBy = utils.StringToPtr(chapterLessonEntity.CreatedBy)

// 	//?Mapping Field Pointer String (*string di Struct)
// 	res.UpdatedBy = utils.PtrStringToPtr(chapterLessonEntity.UpdatedBy)
// 	res.DeletedBy = utils.PtrStringToPtr(chapterLessonEntity.DeletedBy)

// 	//? Mapping Angka int64 biasa
// 	res.OrderChapter = utils.Int64ToPtr(chapterLessonEntity.OrderChapter)

// 	//? Mapping Waktu (Time)
// 	res.CreatedAt = utils.TimeToPtr(chapterLessonEntity.CreatedAt)
// 	res.UpdatedAt = utils.TimeToPtr(chapterLessonEntity.UpdatedAt)
// 	res.DeletedAt = utils.PtrTimeToPtr(chapterLessonEntity.DeletedAt)

// 	return res, nil
// }

// func (cs *chapterLessonService) EditChapterLesson(ctx context.Context, request *chapter_lesson.EditChapterLessonRequest) (*chapter_lesson.EditChapterLessonResponse, error) {
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
// 	courseEntity, err := cs.chapterLessonRepository.GetChapterLessonById(ctx, request.Id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if courseEntity == nil {
// 		return &chapter_lesson.EditChapterLessonResponse{
// 			Base: utils.NotFoundResponse("Course chapter not found"),
// 		}, nil
// 	}

// 	tx, err := cs.db.BeginTxx(ctx, nil)
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

// 	chapterLessonRepo := cs.chapterLessonRepository.WithTransaction(tx)

// 	// *update ke DB
// 	newCourse := entity.ChapterLesson{
// 		Id:           request.Id,
// 		InstructorId: request.InstructorId,
// 		CourseId:     request.CourseId,
// 		Title:        request.Title,
// 		OrderChapter: request.OrderChapter,
// 		Status:       request.Status,

// 		UpdatedAt: time.Now(),
// 		UpdatedBy: &claims.FullName,
// 	}

// 	err = chapterLessonRepo.UpdateChapterLesson(ctx, &newCourse)
// 	if err != nil {
// 		return nil, err
// 	}

// 	err = tx.Commit()
// 	if err != nil {
// 		return nil, err
// 	}

// 	// *success
// 	return &chapter_lesson.EditChapterLessonResponse{
// 		Base: utils.SuccessResponse("Edit Course Chapter Success"),
// 		Id:   request.Id,
// 	}, nil
// }

// func (cs *chapterLessonService) DeleteChapterLesson(ctx context.Context, request *chapter_lesson.DeleteChapterLessonRequest) (*chapter_lesson.DeleteChapterLessonResponse, error) {
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
// 	chapterLessonEntity, err := cs.chapterLessonRepository.GetChapterLessonById(ctx, request.Id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if chapterLessonEntity == nil {
// 		return &chapter_lesson.DeleteChapterLessonResponse{
// 			Base: utils.NotFoundResponse("Course Chapter not found"),
// 		}, nil
// 	}

// 	tx, err := cs.db.BeginTxx(ctx, nil)
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

// 	chapterLessonRepo := cs.chapterLessonRepository.WithTransaction(tx)

// 	// *update delete_at & delete_by ke DB

// 	err = chapterLessonRepo.DeleteChapterLesson(ctx, request.Id, time.Now(), claims.FullName)
// 	if err != nil {
// 		return nil, err
// 	}

// 	err = tx.Commit()
// 	if err != nil {
// 		return nil, err
// 	}

// 	// *success
// 	return &chapter_lesson.DeleteChapterLessonResponse{
// 		Base: utils.SuccessResponse("Delete with SoftDelete Course Success"),
// 	}, nil
// }

func NewChapterLessonService(db *sqlx.DB, chapterLessonRepository repository.IChapterLessonRepository) IChapterLessonService {
	return &chapterLessonService{
		db:                      db,
		chapterLessonRepository: chapterLessonRepository,
	}
}
