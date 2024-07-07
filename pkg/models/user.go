package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	UserID    string `gorm:"primaryKey;type:uuid" json:"userId"`
	FirstName string `gorm:"not null" json:"firstName" validate:"required"`
	LastName  string `gorm:"not null" json:"lastName" validate:"required" `
	Email     string `gorm:"unique;not null" json:"email" validate:"required,email"`
	Password  string `gorm:"not null" json:"password" validate:"required"`
	Phone     string `gorm:"null" json:"phone"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {

	if u.UserID == "" {
		u.UserID = uuid.New().String()
	}
	return
}

type ViewUser struct {
	UserID    string `json:"userId"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}

type UserOrganization struct {
	OrganizationOrgID string `gorm:"primaryKey" json:"organization_org_id"`
	UserUserID        string `gorm:"primaryKey" json:"user_user_id"`
}
