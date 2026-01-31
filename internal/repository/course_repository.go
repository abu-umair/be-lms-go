package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/abu-umair/be-lms-go/internal/entity"
	"github.com/abu-umair/be-lms-go/pkg/database"
)

type ICourseRepository interface {
	WithTransaction(tx *sql.Tx) ICourseRepository
	CreateNewCourse(ctx context.Context, course *entity.Course) error
	GetCourseById(ctx context.Context, courseId string) (*entity.Course, error)
	UpdateCourse(ctx context.Context, course *entity.Course) error
	DeleteCourse(ctx context.Context, id string, deletedAt time.Time, deletedBy string) error
}

type courseRepository struct {
	db database.DatabaseQuery
}

func (ss *courseRepository) WithTransaction(tx *sql.Tx) ICourseRepository {
	return &courseRepository{
		db: tx,
	}
}

func (sr *courseRepository) CreateNewCourse(ctx context.Context, course *entity.Course) error {
	_, err := sr.db.ExecContext(
		ctx,
		`INSERT INTO courses (id, name, image_file_name, address, slug, user_id, category_id, course_type, seo_description, duration, timezone, thumbnail, demo_video_storage, demo_video_source, description, capacity, price, discount, certificate, gna, message_for_reviewer, is_approved, course_level_id, course_language_id, created_at, created_by, updated_at, updated_by, deleted_by)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12,$13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29)`,

		course.Id,
		course.Name,
		course.ImageFileName,
		course.Address,
		course.Slug,
		course.UserId,
		course.CategoryId,
		course.CourseType,
		course.SeoDescription,
		course.Duration,
		course.Timezone,
		course.Thumbnail,
		course.DemoVideoStorage,
		course.DemoVideoSource,
		course.Description,
		course.Capacity,
		course.Price,
		course.Discount,
		course.Certificate,
		course.Gna,
		course.MessageForReviewer,
		course.IsApproved,
		course.CourseLevelId,
		course.CourseLanguageId,
		course.CreatedAt,
		course.CreatedBy,
		course.UpdatedAt,
		course.UpdatedBy,
		course.DeletedBy,
	)

	if err != nil {
		return err
	}

	return nil
}

func (sr *courseRepository) GetCourseById(ctx context.Context, courseId string) (*entity.Course, error) {
	var courseEntity entity.Course

	row := sr.db.QueryRowContext(
		ctx,
		"SELECT id, name, address, image_file_name FROM courses WHERE id = $1 AND deleted_at IS NULL",
		courseId,
	)

	if row.Err() != nil {
		return nil, row.Err()
	}

	err := row.Scan(
		&courseEntity.Id,
		&courseEntity.Name,
		&courseEntity.Address,
		&courseEntity.ImageFileName,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &courseEntity, nil
}

func (sr *courseRepository) UpdateCourse(ctx context.Context, course *entity.Course) error {
	_, err := sr.db.ExecContext(
		ctx, "UPDATE courses SET name = $1, address = $2, image_file_name = $3, updated_at = $4, updated_by = $5 WHERE id = $6",
		course.Name,
		course.Address,
		course.ImageFileName,
		course.UpdatedAt,
		course.UpdatedBy,
		course.Id,
	)

	if err != nil {
		return nil
	}

	return nil
}

func (sr *courseRepository) DeleteCourse(ctx context.Context, id string, deletedAt time.Time, deletedBy string) error {
	_, err := sr.db.ExecContext(
		ctx, "UPDATE courses SET deleted_at = $1, deleted_by = $2 WHERE id = $3",
		deletedAt,
		deletedBy,
		id,
	)

	if err != nil {
		return nil
	}

	return nil
}

func NewCourseRepository(db database.DatabaseQuery) ICourseRepository {
	return &courseRepository{db: db}
}
