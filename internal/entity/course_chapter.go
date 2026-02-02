package entity

import "time"

type CourseChapter struct {
	Id           string `db:"id"`
	InstructorId string `db:"instructor_id"`
	CourseId     string `db:"course_id"`
	Title        string `db:"title"`
	Order        int64  `db:"order"`
	Status       string `db:"status"`

	CreatedAt time.Time  `db:"created_at"`
	CreatedBy string     `db:"created_by"`
	UpdatedAt *time.Time  `db:"updated_at"`
	UpdatedBy *string    `db:"updated_by"`
	DeletedAt *time.Time `db:"deleted_at"`
	DeletedBy *string    `db:"deleted_by"`
}
