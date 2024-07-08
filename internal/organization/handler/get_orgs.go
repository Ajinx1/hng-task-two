package handler

import (
	"hng-task-two/pkg/middleware"
	"hng-task-two/pkg/models"
	"hng-task-two/pkg/responses"

	"github.com/gofiber/fiber/v2"
)

func (h *OrganizationHandler) GetAllOrganizations(c *fiber.Ctx) error {
	userID := middleware.GetUserId(c)
	organizations, err := h.organizationService.GetUserOrganizations(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.UserResponse{
			Status:     "Bad request",
			Message:    "Client error",
			StatusCode: 400,
		})
	}

	return c.Status(fiber.StatusOK).JSON(responses.AllOrgResponse{
		Status:  "success",
		Message: "<message>",
		Data: &fiber.Map{
			"organisations": &organizations,
		},
	})
}

func (h *OrganizationHandler) GetOrganization(c *fiber.Ctx) error {
	orgID := c.Params("orgId")
	userID := middleware.GetUserId(c)
	organization, err := h.organizationService.GetOrganizationByID(orgID, userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.UserResponse{
			Status:     "Bad request",
			Message:    "Client error",
			StatusCode: 400,
		})
	}

	newOrganization := &models.ViewOrganization{
		OrgID:       organization.OrgID,
		Name:        organization.Name,
		Description: organization.Description,
	}

	return c.Status(fiber.StatusOK).JSON(responses.OrgResponse{
		Status:  "success",
		Message: "<message>",
		Data:    newOrganization,
	})

}
