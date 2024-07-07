package repository

import (
	"fmt"
	"hng-task-two/pkg/models"

	"gorm.io/gorm"
)

type OrganizationRepository interface {
	Create(organization *models.Organization) error
	Update(userID, orgID string) error
	AddUserToOrganization(org *models.Organization, user *models.User) error
	FindByID(id, requestID string) (*models.Organization, error)
	FindByUserID(userID string) ([]models.ViewOrganization, error)
}

type organizationRepository struct {
	db *gorm.DB
}

func NewOrganizationRepository(db *gorm.DB) OrganizationRepository {
	return &organizationRepository{db}
}

func (r *organizationRepository) Create(organization *models.Organization) error {
	return r.db.Create(organization).Error
}

func (r *organizationRepository) AddUserToOrganization(org *models.Organization, user *models.User) error {
	return r.db.Model(org).Association("Users").Append(user)
}

func (r *organizationRepository) FindByID(id string, requesterID string) (*models.Organization, error) {
	var organization models.Organization

	query := r.db.Model(&models.Organization{}).
		Joins(`
            INNER JOIN user_organizations uo ON organizations.org_id = uo.organization_org_id
        `).
		Where("organizations.org_id = ?", id).
		Where("uo.user_user_id = ?", requesterID)

	if err := query.First(&organization).Error; err != nil {
		return nil, err
	}

	return &organization, nil
}

func (r *organizationRepository) FindByUserID(userID string) ([]models.ViewOrganization, error) {
	var viewOrganizations []models.ViewOrganization

	err := r.db.
		Table("organizations").
		Select("organizations.org_id, organizations.name, organizations.description").
		Joins("JOIN user_organizations ON user_organizations.organization_org_id = organizations.org_id").
		Where("user_organizations.user_user_id = ?", userID).
		Scan(&viewOrganizations).Error

	return viewOrganizations, err
}

func (r *organizationRepository) Update(userID, orgID string) error {

	var existingLink models.UserOrganization

	if err := r.db.First(&existingLink, "organization_org_id = ? AND user_user_id = ?", orgID, userID).Error; err == nil {
		fmt.Println("user exist")
		return err
	}

	link := models.UserOrganization{
		OrganizationOrgID: orgID,
		UserUserID:        userID,
	}
	if err := r.db.Create(&link).Error; err != nil {
		return err
	}

	return nil
}
