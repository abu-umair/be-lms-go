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

type IChapterLessonRepository interface {
	WithTransaction(tx *sqlx.Tx) IChapterLessonRepository
	CreateNewChapterLesson(ctx context.Context, chapterLesson *entity.ChapterLesson) error
	GetChapterLessonById(ctx context.Context, chapterLessonId string) (*entity.ChapterLesson, error)
	GetChapterLessonByIdFieldMask(ctx context.Context, chapterLessonId string, paths []string) (*entity.ChapterLesson, error)
	UpdateChapterLesson(ctx context.Context, chapterLesson *entity.ChapterLesson) error
	DeleteChapterLesson(ctx context.Context, id string, deletedAt time.Time, deletedBy string) error
}

type chapterLessonRepository struct {
	db database.DatabaseQuery
	// Kita simpan di sini agar tidak perlu buat map berulang-ulang di setiap request
	whitelist map[string]bool
}

func (cs *chapterLessonRepository) WithTransaction(tx *sqlx.Tx) IChapterLessonRepository {
	return &chapterLessonRepository{
		db: tx,
	}
}

func (cr *chapterLessonRepository) CreateNewChapterLesson(ctx context.Context, chapterLesson *entity.ChapterLesson) error {
	query := `
        INSERT INTO course_chapter_lessons (
			instructor_id,course_id,title,order_lesson,chapter_id,slug,description,file_path,storage_lesson,lesson_type,volume,duration,file_type,downloadable,is_preview,status, created_at, created_by, updated_at, updated_by, deleted_by
        )
        VALUES (
            :instructor_id,:course_id,:title,:order_lesson,:chapter_id,:slug,:description,:file_path,:storage_lesson,:lesson_type,:volume,:duration,:file_type,:downloadable,:is_preview,:status, :created_at, :created_by, :updated_at, :updated_by, :deleted_by
        )`

	// NamedExecContext akan otomatis mencocokkan :id dengan field di struct
	// yang memiliki tag db:"id"
	_, err := cr.db.NamedExecContext(ctx, query, chapterLesson)
	if err != nil {
		return err
	}

	return nil
}

func (cr *chapterLessonRepository) GetChapterLessonById(ctx context.Context, chapterLessonId string) (*entity.ChapterLesson, error) {
	var chapterLessonEntity entity.ChapterLesson

	// 1. Tentukan query
	query := `SELECT id
	          FROM course_chapter_lessons
	          WHERE id = $1 AND deleted_at IS NULL`

	// 2. Gunakan GetContext untuk mapping otomatis
	// sqlx akan mencocokkan kolom SELECT dengan tag `db` di struct entity.Course
	err := cr.db.GetContext(ctx, &chapterLessonEntity, query, chapterLessonId)

	if err != nil {
		// 3. Tangani jika data tidak ditemukan
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &chapterLessonEntity, nil
}

func (cr *chapterLessonRepository) GetChapterLessonByIdFieldMask(ctx context.Context, chapterLessonId string, paths []string) (*entity.ChapterLesson, error) {
	var chapterLessonEntity entity.ChapterLesson

	// 1. Tentukan kolom yang akan di-select
	selectedColumns := "*" // Default jika paths kosong
	if len(paths) > 0 {
		var validColumns []string
		for _, p := range paths {
			// Cek apakah kolom yang diminta ada di whitelist kita
			if cr.whitelist[p] {
				validColumns = append(validColumns, p)
			}
		}

		if len(validColumns) > 0 {
			selectedColumns = strings.Join(validColumns, ", ")
		}
	}

	// 2. Tentukan query dengan kolom dinamis
	query := fmt.Sprintf(`SELECT %s FROM course_chapter_lessons WHERE id = $1 AND deleted_at IS NULL`, selectedColumns)

	// 3. Gunakan GetContext (sqlx tetap bisa memetakan meskipun kolomnya cuma sedikit)
	err := cr.db.GetContext(ctx, &chapterLessonEntity, query, chapterLessonId)

	if err != nil {
		// 3. Tangani jika data tidak ditemukan
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, err
	}

	return &chapterLessonEntity, nil
}

func (sr *chapterLessonRepository) UpdateChapterLesson(ctx context.Context, chapterLesson *entity.ChapterLesson) error {
	// Menggunakan Named Query (:field) yang merujuk pada tag db di struct entity
	query := `
		UPDATE course_chapter_lessons
		SET
			instructor_id = :instructor_id,
			course_id = :course_id,
			title = :title,
			order_chapter = :order_chapter,
			status = :status,
			
			updated_at = :updated_at,
			updated_by = :updated_by
		WHERE id = :id`

	_, err := sr.db.NamedExecContext(ctx, query, chapterLesson)
	return err
}

func (sr *chapterLessonRepository) DeleteChapterLesson(ctx context.Context, id string, deletedAt time.Time, deletedBy string) error {
	query := `UPDATE course_chapter_lessons SET deleted_at = :deleted_at, deleted_by = :deleted_by WHERE id = :id`

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

func NewChapterLessonRepository(db database.DatabaseQuery) IChapterLessonRepository {
	return &chapterLessonRepository{
		db: db,
		whitelist: map[string]bool{
			"id":             true,
			"title":          true,
			"slug":           true,
			"order_lesson":   true,
			"description":    true,
			"file_path":      true,
			"storage_lesson": true,
			"lesson_type":    true,
			"volume":         true,
			"duration":       true,
			"file_type":      true,
			"downloadable":   true,
			"is_preview":     true,
			"status":         true,
			"instructor_id":  true,
			"course_id":      true,
			
			"created_at":     true, "created_by": true, "updated_at": true,
			"updated_by": true, "deleted_at": true, "deleted_by": true,
		},
	}
}
