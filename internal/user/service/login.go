package service

import (
	"errors"
	"hng-task-two/pkg/middleware"
	"hng-task-two/pkg/models"

	"golang.org/x/crypto/bcrypt"
)

func (s *userService) Login(email, password string) (*models.User, string, error) {
	user, err := s.userRepository.FindByEmail(email)
	if err != nil {
		return nil, "", errors.New("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	token, err := middleware.GenerateJWT(user)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}
