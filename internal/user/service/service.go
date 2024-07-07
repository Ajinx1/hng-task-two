package service

import (
	org_repo "hng-task-two/internal/organization/repository"
	"hng-task-two/internal/user/repository"
	"hng-task-two/pkg/models"
)

type UserService interface {
	Register(user *models.User) (*models.User, string, error)
	Login(email, password string) (*models.User, string, error)
	GetUserByID(id string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUserByIDs(id, requestId string) (*models.User, error)
}

type userService struct {
	userRepository repository.UserRepository
	orgRepository  org_repo.OrganizationRepository
}

func NewUserService(userRepo repository.UserRepository, orgRepo org_repo.OrganizationRepository) UserService {
	return &userService{userRepo, orgRepo}
}

func (s *userService) GetUserByEmail(email string) (*models.User, error) {
	return s.userRepository.FindByEmail(email)
}

func (s *userService) GetUserByID(id string) (*models.User, error) {
	return s.userRepository.FindByID(id)
}

func (s *userService) GetUserByIDs(id, requestId string) (*models.User, error) {
	return s.userRepository.FindByIDs(id, requestId)
}
