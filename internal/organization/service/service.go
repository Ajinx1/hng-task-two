package service

import (
	"hng-task-two/internal/organization/repository"
	user_repo "hng-task-two/internal/user/repository"
	"hng-task-two/pkg/models"
)

type OrganizationService interface {
	Create(organization *models.Organization, userID string) (*models.Organization, error)
	GetOrganizationByID(id, requesterID string) (*models.Organization, error)
	GetUserOrganizations(userID string) ([]models.ViewOrganization, error)
	AddUserToOrganization(userID, orgID, requesterID string) error
}

type organizationService struct {
	orgRepository  repository.OrganizationRepository
	userRepository user_repo.UserRepository
}

func NewOrganizationService(orgRepo repository.OrganizationRepository, userRepo user_repo.UserRepository) OrganizationService {
	return &organizationService{orgRepo, userRepo}
}

func (s *organizationService) Create(organization *models.Organization, userID string) (*models.Organization, error) {
	if err := s.orgRepository.Create(organization); err != nil {
		return nil, err
	}

	user, err := s.userRepository.FindByID(userID)
	if err != nil {
		return nil, err
	}

	if err := s.orgRepository.AddUserToOrganization(organization, user); err != nil {
		return nil, err
	}

	return organization, nil
}

func (s *organizationService) GetOrganizationByID(id, requesterID string) (*models.Organization, error) {
	return s.orgRepository.FindByID(id, requesterID)
}

func (s *organizationService) GetUserOrganizations(userID string) ([]models.ViewOrganization, error) {
	return s.orgRepository.FindByUserID(userID)
}

func (s *organizationService) AddUserToOrganization(userID, orgID, requesterID string) error {
	_, err := s.userRepository.FindByID(userID)
	if err != nil {
		return err
	}

	return s.orgRepository.Update(userID, orgID)
}
