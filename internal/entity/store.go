package entity

import "time"

type Store struct {
	Id            string
	Name          string
	Address       string
	ImageFileName string
	CreatedAt     time.Time
	CreatedBy     string
	UpdatedAt     time.Time
	UpdatedBy     *string
	DeletedAt     *time.Time
	DeletedBy     *string
	// IsDeleted  bool
}
