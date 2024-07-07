package repository

import (
	"hng-task-two/pkg/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error
	FindByEmail(email string) (*models.User, error)
	FindByID(id string) (*models.User, error)
	FindByIDs(id string, requesterID string) (*models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *userRepository) FindByID(id string) (*models.User, error) {
	var user models.User
	err := r.db.Where("user_id = ?", id).First(&user).Error
	return &user, err
}

func (r *userRepository) FindByIDs(id string, requesterID string) (*models.User, error) {
	var user models.User

	// Fetch user details only if it belongs to the same organization(s) as the requester
	query := r.db.Model(&models.User{}).
		Joins(`
            INNER JOIN user_organizations uo ON users.user_id = uo.user_user_id
        `).
		Where("users.user_id = ?", id).
		Where("uo.organization_org_id IN (SELECT organization_org_id FROM user_organizations WHERE user_user_id = ?)", requesterID)

	if err := query.First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
