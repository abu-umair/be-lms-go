package entity

import "time"

const ( //? kebutuhan di layer service (insert ke DB)
	UserRoleUser         = "user"
	UserRoleKitchenStaff = "kitchenstaff"
	UserRoleCashier      = "cashier"
	UserRoleOwner        = "owner"
	UserRoleAdmin        = "admin"
)

type UserRoles struct {
	Id        string
	Name      string
	Code      string
	CreatedAt time.Time
	CreatedBy *string
	UpdatedAt time.Time
	UpdatedBy *string
	DeletedAt time.Time
	DeletedBy *string
	IsDeleted bool
}

type Users struct {
	Id         string
	FullName   string
	Email      string
	Password   string
	RoleCode   string
	VerifiedAt *time.Time
	CreatedAt  time.Time
	CreatedBy  *string
	UpdatedAt  time.Time
	UpdatedBy  *string
	DeletedAt  *time.Time
	DeletedBy  *string
	// IsDeleted  bool
}
