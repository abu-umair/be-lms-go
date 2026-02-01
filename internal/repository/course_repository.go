package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/abu-umair/be-lms-go/internal/entity"
	"github.com/abu-umair/be-lms-go/pkg/database"
	"github.com/jmoiron/sqlx"
)

type ICourseRepository interface {
	WithTransaction(tx *sqlx.Tx) ICourseRepository
	CreateNewCourse(ctx context.Context, course *entity.Course) error
	GetCourseById(ctx context.Context, courseId string) (*entity.Course, error)
	UpdateCourse(ctx context.Context, course *entity.Course) error
	DeleteCourse(ctx context.Context, id string, deletedAt time.Time, deletedBy string) error
}

type courseRepository struct {
	db database.DatabaseQuery
}

func (ss *courseRepository) WithTransaction(tx *sqlx.Tx) ICourseRepository {
	return &courseRepository{
		db: tx,
	}
}

func (sr *courseRepository) CreateNewCourse(ctx context.Context, course *entity.Course) error {
	query := `
        INSERT INTO courses (
            id, name, image_file_name, address, slug, user_id, category_id, course_type, 
            seo_description, duration, timezone, thumbnail, demo_video_storage, 
            demo_video_source, description, capacity, price, discount, certificate, 
            gna, message_for_reviewer, is_approved, course_level_id, 
            course_language_id, created_at, created_by, updated_at, updated_by, deleted_by
        )
        VALUES (
            :id, :name, :image_file_name, :address, :slug, :user_id, :category_id, :course_type, 
            :seo_description, :duration, :timezone, :thumbnail, :demo_video_storage, 
            :demo_video_source, :description, :capacity, :price, :discount, :certificate, 
            :gna, :message_for_reviewer, :is_approved, :course_level_id, 
            :course_language_id, :created_at, :created_by, :updated_at, :updated_by, :deleted_by
        )`

	// NamedExecContext akan otomatis mencocokkan :id dengan field di struct
	// yang memiliki tag db:"id"
	_, err := sr.db.NamedExecContext(ctx, query, course)
	if err != nil {
		return err
	}

	return nil
}

func (sr *courseRepository) GetCourseById(ctx context.Context, courseId string) (*entity.Course, error) {
	var courseEntity entity.Course

	// 1. Tentukan query
	query := `SELECT id, name, address, image_file_name 
	          FROM courses 
	          WHERE id = $1 AND deleted_at IS NULL`

	// 2. Gunakan GetContext untuk mapping otomatis
	// sqlx akan mencocokkan kolom SELECT dengan tag `db` di struct entity.Course
	err := sr.db.GetContext(ctx, &courseEntity, query, courseId)

	if err != nil {
		// 3. Tangani jika data tidak ditemukan
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &courseEntity, nil
}

func (sr *courseRepository) UpdateCourse(ctx context.Context, course *entity.Course) error {
	// Menggunakan Named Query (:field) yang merujuk pada tag db di struct entity
	query := `
		UPDATE courses 
		SET 
			name = :name, 
			address = :address, 
			image_file_name = :image_file_name, 
			updated_at = :updated_at, 
			updated_by = :updated_by 
		WHERE id = :id`

	_, err := sr.db.NamedExecContext(ctx, query, course)
	// Langsung return err jika ada, atau nil jika sukses
	return err
}

func (sr *courseRepository) DeleteCourse(ctx context.Context, id string, deletedAt time.Time, deletedBy string) error {
	query := `UPDATE courses SET deleted_at = :deleted_at, deleted_by = :deleted_by WHERE id = :id`

	// Kita bungkus data ke dalam map agar bisa dibaca oleh NamedExecContext
	data := map[string]any{
		"deleted_at": deletedAt,
		"deleted_by": deletedBy,
		"id":         id,
	}

	_, err := sr.db.NamedExecContext(ctx, query, data)

	if err != nil {
		return err
	}

	return nil
}

func NewCourseRepository(db database.DatabaseQuery) ICourseRepository {
	return &courseRepository{db: db}
}
