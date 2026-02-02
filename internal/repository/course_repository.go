package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/abu-umair/be-lms-go/internal/entity"
	"github.com/abu-umair/be-lms-go/pkg/database"
	"github.com/jmoiron/sqlx"
)

type ICourseRepository interface {
	WithTransaction(tx *sqlx.Tx) ICourseRepository
	CreateNewCourse(ctx context.Context, course *entity.Course) error
	GetCourseById(ctx context.Context, courseId string) (*entity.Course, error)
	GetCourseByIdFieldMask(ctx context.Context, courseId string, paths []string) (*entity.Course, error)
	UpdateCourse(ctx context.Context, course *entity.Course) error
	DeleteCourse(ctx context.Context, id string, deletedAt time.Time, deletedBy string) error
}

type courseRepository struct {
	db database.DatabaseQuery
	// Kita simpan di sini agar tidak perlu buat map berulang-ulang di setiap request
	whitelist map[string]bool
}

func (ss *courseRepository) WithTransaction(tx *sqlx.Tx) ICourseRepository {
	return &courseRepository{
		db: tx,
	}
}

func (sr *courseRepository) CreateNewCourse(ctx context.Context, course *entity.Course) error {
	query := `
        INSERT INTO courses (
            id, name, image_file_name, address, slug, instructor_id, category_id, course_type, 
            seo_description, duration, timezone, thumbnail, demo_video_storage, 
            demo_video_source, description, capacity, price, discount, certificate, 
            gna, message_for_reviewer, is_approved, course_level_id, 
            course_language_id, created_at, created_by, updated_at, updated_by, deleted_by
        )
        VALUES (
            :id, :name, :image_file_name, :address, :slug, :instructor_id, :category_id, :course_type, 
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
	query := `SELECT id, image_file_name 
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

func (sr *courseRepository) GetCourseByIdFieldMask(ctx context.Context, courseId string, paths []string) (*entity.Course, error) {
	var courseEntity entity.Course

	// 1. Tentukan kolom yang akan di-select
	selectedColumns := "*" // Default jika paths kosong
	if len(paths) > 0 {
		var validColumns []string
		for _, p := range paths {
			// Cek apakah kolom yang diminta ada di whitelist kita
			if sr.whitelist[p] {
				validColumns = append(validColumns, p)
			}
		}

		if len(validColumns) > 0 {
			selectedColumns = strings.Join(validColumns, ", ")
		}
	}

	// 2. Tentukan query dengan kolom dinamis
	query := fmt.Sprintf(`SELECT %s FROM courses WHERE id = $1 AND deleted_at IS NULL`, selectedColumns)

	// 3. Gunakan GetContext (sqlx tetap bisa memetakan meskipun kolomnya cuma sedikit)
	err := sr.db.GetContext(ctx, &courseEntity, query, courseId)

	if err != nil {
		// 3. Tangani jika data tidak ditemukan
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
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
			slug = :slug,
			instructor_id = :instructor_id,
			category_id = :category_id,
			course_type = :course_type,
			seo_description = :seo_description,
			duration = :duration,
			timezone = :timezone,
			thumbnail = :thumbnail,
			demo_video_storage = :demo_video_storage,
			demo_video_source = :demo_video_source,
			description = :description,
			capacity = :capacity,
			price = :price,
			discount = :discount,
			certificate = :certificate,
			gna = :gna,
			message_for_reviewer = :message_for_reviewer,
			is_approved = :is_approved,
			status = :status,
			course_level_id = :course_level_id,
			course_language_id = :course_language_id,

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
	return &courseRepository{
		db: db,
		whitelist: map[string]bool{
			"id": true, "name": true, "address": true, "image_file_name": true,
			"created_at": true, "created_by": true, "updated_at": true,
			"updated_by": true, "deleted_at": true, "deleted_by": true,
			"slug": true, "instructor_id": true, "category_id": true,
			"course_type": true, "seo_description": true, "duration": true,
			"timezone": true, "thumbnail": true, "demo_video_storage": true,
			"demo_video_source": true, "description": true, "capacity": true,
			"price": true, "discount": true, "certificate": true,
			"gna": true, "message_for_reviewer": true, "is_approved": true,
			"status": true, "course_level_id": true, "course_language_id": true,
		},
	}
}
