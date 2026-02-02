package handler

import (
	"context"

	"github.com/abu-umair/be-lms-go/internal/service"
	"github.com/abu-umair/be-lms-go/internal/utils"
	"github.com/abu-umair/be-lms-go/pb/course_chapter"
)

type courseChapterHandler struct {
	course_chapter.UnimplementedCourseChapterServiceServer

	courseChapterService service.ICourseChapterService //? layer service
}

func (ch *courseChapterHandler) CreateCourseChapter(ctx context.Context, request *course_chapter.CreateCourseChapterRequest) (*course_chapter.CreateCourseChapterResponse, error) {
	//? validasi request
	validationErrors, err := utils.CheckValidation(request)
	if err != nil {
		return nil, err
	}

	if validationErrors != nil {
		return &course_chapter.CreateCourseChapterResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}

	res, err := ch.courseChapterService.CreateCourseChapter(ctx, request)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (ch *courseChapterHandler) DetailCourseChapter(ctx context.Context, request *course_chapter.DetailCourseChapterRequest) (*course_chapter.DetailCourseChapterResponse, error) {
	//? validasi request
	validationErrors, err := utils.CheckValidation(request)
	if err != nil {
		return nil, err
	}

	if validationErrors != nil {
		return &course_chapter.DetailCourseChapterResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}

	res, err := ch.courseChapterService.DetailCourseChapter(ctx, request)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (ch *courseChapterHandler) EditCourseChapter(ctx context.Context, request *course_chapter.EditCourseChapterRequest) (*course_chapter.EditCourseChapterResponse, error) {
	validationErrors, err := utils.CheckValidation(request)
	if err != nil {
		return nil, err
	}

	if validationErrors != nil {
		return &course_chapter.EditCourseChapterResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}

	res, err := ch.courseChapterService.EditCourseChapter(ctx, request)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (ch *courseChapterHandler) DeleteCourseChapter(ctx context.Context, request *course_chapter.DeleteCourseChapterRequest) (*course_chapter.DeleteCourseChapterResponse, error) {
	validationErrors, err := utils.CheckValidation(request)
	if err != nil {
		return nil, err
	}

	if validationErrors != nil {
		return &course_chapter.DeleteCourseChapterResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}

	res, err := ch.courseChapterService.DeleteCourseChapter(ctx, request)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func NewCourseChapterHandler(courseChapterService service.ICourseChapterService) *courseChapterHandler {
	return &courseChapterHandler{
		courseChapterService: courseChapterService,
	}
}
