package handler

import (
	"context"

	"github.com/abu-umair/be-lms-go/internal/service"
	"github.com/abu-umair/be-lms-go/internal/utils"
	"github.com/abu-umair/be-lms-go/pb/chapter_lesson"
)

type chapterLessonHandler struct {
	chapter_lesson.UnimplementedChapterLessonServiceServer

	chapterLessonService service.IChapterLessonService //? layer service
}

func (lh *chapterLessonHandler) CreateChapterLesson(ctx context.Context, request *chapter_lesson.CreateChapterLessonRequest) (*chapter_lesson.CreateChapterLessonResponse, error) {
	//? validasi request
	validationErrors, err := utils.CheckValidation(request)
	if err != nil {
		return nil, err
	}

	if validationErrors != nil {
		return &chapter_lesson.CreateChapterLessonResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}

	res, err := lh.chapterLessonService.CreateChapterLesson(ctx, request)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (lh *chapterLessonHandler) DetailChapterLesson(ctx context.Context, request *chapter_lesson.DetailChapterLessonRequest) (*chapter_lesson.DetailChapterLessonResponse, error) {
	//? validasi request
	validationErrors, err := utils.CheckValidation(request)
	if err != nil {
		return nil, err
	}

	if validationErrors != nil {
		return &chapter_lesson.DetailChapterLessonResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}

	res, err := lh.chapterLessonService.DetailChapterLesson(ctx, request)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (ch *chapterLessonHandler) EditChapterLesson(ctx context.Context, request *chapter_lesson.EditChapterLessonRequest) (*chapter_lesson.EditChapterLessonResponse, error) {
	validationErrors, err := utils.CheckValidation(request)
	if err != nil {
		return nil, err
	}

	if validationErrors != nil {
		return &chapter_lesson.EditChapterLessonResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}

	res, err := ch.chapterLessonService.EditChapterLesson(ctx, request)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (ch *chapterLessonHandler) DeleteChapterLesson(ctx context.Context, request *chapter_lesson.DeleteChapterLessonRequest) (*chapter_lesson.DeleteChapterLessonResponse, error) {
	validationErrors, err := utils.CheckValidation(request)
	if err != nil {
		return nil, err
	}

	if validationErrors != nil {
		return &chapter_lesson.DeleteChapterLessonResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}

	res, err := ch.chapterLessonService.DeleteChapterLesson(ctx, request)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func NewChapterLessonHandler(chapterLessonService service.IChapterLessonService) *chapterLessonHandler {
	return &chapterLessonHandler{
		chapterLessonService: chapterLessonService,
	}
}
