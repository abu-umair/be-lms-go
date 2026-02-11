package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

type Course struct {
	Id                 string           `db:"id"`
	Name               string           `db:"name"`
	Address            *string          `db:"address"`
	ImageFileName      string           `db:"image_file_name"`
	CreatedAt          time.Time        `db:"created_at"`
	CreatedBy          string           `db:"created_by"`
	UpdatedAt          time.Time        `db:"updated_at"`
	UpdatedBy          *string          `db:"updated_by"`
	DeletedAt          *time.Time       `db:"deleted_at"`
	DeletedBy          *string          `db:"deleted_by"`
	Slug               *string          `db:"slug"`
	InstructorId       *string          `db:"instructor_id"`
	CategoryId         *string          `db:"category_id"`
	CourseType         *string          `db:"course_type"`
	SeoDescription     *string          `db:"seo_description"`
	Duration           *string          `db:"duration"`
	Timezone           *string          `db:"timezone"`
	Thumbnail          *string          `db:"thumbnail"`
	DemoVideoStorage   *string          `db:"demo_video_storage"`
	DemoVideoSource    *string          `db:"demo_video_source"`
	Description        *string          `db:"description"`
	Capacity           *int32           `db:"capacity"`
	Price              *decimal.Decimal `db:"price"`
	Discount           *decimal.Decimal `db:"discount"`
	Certificate        *string          `db:"certificate"`
	Gna                *string          `db:"gna"`
	MessageForReviewer *string          `db:"message_for_reviewer"`
	IsApproved         *string          `db:"is_approved"`
	Status             *string          `db:"status"`
	CourseLevelId      *string          `db:"course_level_id"`
	CourseLanguageId   *string          `db:"course_language_id"`
}
