package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

type Course struct {
	Id                 string
	Name               string
	Address            string
	ImageFileName      string
	CreatedAt          time.Time
	CreatedBy          string
	UpdatedAt          time.Time
	UpdatedBy          *string
	DeletedAt          *time.Time
	DeletedBy          *string
	Slug               *string
	UserId             *string
	CategoryId         *string
	CourseType         *string
	SeoDescription     *string
	Duration           *string
	Timezone           *string
	Thumbnail          *string
	DemoVideoStorage   *string
	DemoVideoSource    *string
	Description        *string
	Capacity           *int64
	Price              *decimal.Decimal
	Discount           *decimal.Decimal
	Certificate        *string
	Gna                *string
	MessageForReviewer *string
	IsApproved         *string
	Status             *string
	CourseLevelId      *string
	CourseLanguageId   *string
}
