package handler

import (
	"hng-task-two/internal/organization/service"

	"github.com/go-playground/validator"
)

type OrganizationHandler struct {
	organizationService service.OrganizationService
	validate            *validator.Validate
}

func NewOrganizationHandler(orgService service.OrganizationService) *OrganizationHandler {

	return &OrganizationHandler{
		organizationService: orgService,
		validate:            validator.New(),
	}
}
