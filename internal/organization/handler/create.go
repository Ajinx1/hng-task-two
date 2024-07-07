package handler

import (
	"hng-task-two/pkg/middleware"
	"hng-task-two/pkg/models"
	"hng-task-two/pkg/responses"
	"hng-task-two/pkg/reuseable"

	"github.com/gofiber/fiber/v2"
)

func (h *OrganizationHandler) Create(c *fiber.Ctx) error {
	var newOrg = new(models.Organization)
	userID := middleware.GetUserId(c)

	if err := c.BodyParser(&newOrg); err != nil {

		return c.Status(fiber.StatusBadRequest).JSON(responses.UserResponse{
			Status:     "Bad request",
			Message:    "Client error",
			StatusCode: 400,
		})
	}

	if err := h.validate.Struct(newOrg); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(reuseable.FormatValidationError(err))
	}

	organization, err := h.organizationService.Create(newOrg, userID)
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

	return c.Status(fiber.StatusCreated).JSON(responses.OrgResponse{
		Status:  "success",
		Message: "Organisation created successfully",
		Data:    newOrganization,
	})

}
