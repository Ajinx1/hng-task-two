package service

import (
	"hng-task-two/pkg/middleware"
	"hng-task-two/pkg/models"

	"golang.org/x/crypto/bcrypt"
)

func (s *userService) Register(user *models.User) (*models.User, string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}
	user.Password = string(hashedPassword)

	if err := s.userRepository.Create(user); err != nil {
		return nil, "", err
	}

	org := &models.Organization{
		Name:  user.FirstName + "'s Organisation",
		Users: []models.User{*user},
	}

	if err := s.orgRepository.Create(org); err != nil {
		return nil, "", err
	}

	if err := s.orgRepository.AddUserToOrganization(org, user); err != nil {
		return nil, "", err
	}

	token, err := middleware.GenerateJWT(user)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}
