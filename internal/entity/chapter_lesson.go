package entity

import "time"

type ChapterLesson struct {
	Id            string  `db:"id"`
	InstructorId  *string `db:"instructor_id"`
	CourseId      *string `db:"course_id"`
	Title         string  `db:"title"`
	OrderLesson   int64   `db:"order_lesson"`
	ChapterId     *string `db:"chapter_id"`
	Slug          *string `db:"slug"`
	Description   *string `db:"description"`
	FilePath      *string `db:"file_path"`
	StorageLesson *string `db:"storage_lesson"`
	LessonType    *string `db:"lesson_type"`
	Volume        *string `db:"volume"`
	Duration      *string `db:"duration"`
	FileType      *string `db:"file_type"`
	Downloadable  *string `db:"downloadable"`
	IsPreview     *int64  `db:"is_preview"`
	Status        *string `db:"status"`

	CreatedAt time.Time  `db:"created_at"`
	CreatedBy string     `db:"created_by"`
	UpdatedAt time.Time  `db:"updated_at"`
	UpdatedBy *string    `db:"updated_by"`
	DeletedAt *time.Time `db:"deleted_at"`
	DeletedBy *string    `db:"deleted_by"`
}
