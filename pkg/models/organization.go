package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Organization struct {
	OrgID       string `gorm:"primaryKey;type:uuid" json:"orgId"`
	Name        string `gorm:"not null" json:"name" validate:"required"`
	Description string `gorm:"null" json:"description"`
	Users       []User `gorm:"many2many:user_organizations" json:"users"`
}

func (o *Organization) BeforeCreate(tx *gorm.DB) (err error) {

	if o.OrgID == "" {
		o.OrgID = uuid.New().String()
	}
	return
}

type AddToOrg struct {
	UserID string `json:"userId" validate:"required"`
}

type ViewOrganization struct {
	OrgID       string `gorm:"primaryKey;type:uuid" json:"orgId"`
	Name        string `gorm:"not null" json:"name" validate:"required"`
	Description string `gorm:"null" json:"description"`
}
