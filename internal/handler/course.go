package handler

import (
	"context"

	"github.com/abu-umair/be-lms-go/internal/service"
	"github.com/abu-umair/be-lms-go/internal/utils"
	"github.com/abu-umair/be-lms-go/pb/course"
)

type courseHandler struct {
	course.UnimplementedCourseServiceServer

	courseService service.ICourseService //? layer service
}

func (sh *courseHandler) CreateCourse(ctx context.Context, request *course.CreateCourseRequest) (*course.CreateCourseResponse, error) {
	//? validasi request
	validationErrors, err := utils.CheckValidation(request)
	if err != nil {
		return nil, err
	}

	if validationErrors != nil {
		return &course.CreateCourseResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}

	res, err := sh.courseService.CreateCourse(ctx, request)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (sh *courseHandler) DetailCourse(ctx context.Context, request *course.DetailCourseRequest) (*course.DetailCourseResponse, error) {
	//? validasi request
	validationErrors, err := utils.CheckValidation(request)
	if err != nil {
		return nil, err
	}

	if validationErrors != nil {
		return &course.DetailCourseResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}

	res, err := sh.courseService.DetailCourse(ctx, request)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (sh *courseHandler) EditCourse(ctx context.Context, request *course.EditCourseRequest) (*course.EditCourseResponse, error) {
	validationErrors, err := utils.CheckValidation(request)
	if err != nil {
		return nil, err
	}

	if validationErrors != nil {
		return &course.EditCourseResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}

	res, err := sh.courseService.EditCourse(ctx, request)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (sh *courseHandler) DeleteCourse(ctx context.Context, request *course.DeleteCourseRequest) (*course.DeleteCourseResponse, error) {
	validationErrors, err := utils.CheckValidation(request)
	if err != nil {
		return nil, err
	}

	if validationErrors != nil {
		return &course.DeleteCourseResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}

	res, err := sh.courseService.DeleteCourse(ctx, request)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func NewCourseHandler(courseService service.ICourseService) *courseHandler {
	return &courseHandler{
		courseService: courseService,
	}
}
