package handler

import (
	"hng-task-two/pkg/middleware"
	"hng-task-two/pkg/models"
	"hng-task-two/pkg/responses"
	"hng-task-two/pkg/reuseable"

	"github.com/gofiber/fiber/v2"
)

func (h *OrganizationHandler) AddUserToOrganization(c *fiber.Ctx) error {
	var request models.AddToOrg
	userID := middleware.GetUserId(c)

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.UserResponse{
			Status:     "Bad request",
			Message:    "Client error",
			StatusCode: 400,
		})
	}

	if err := h.validate.Struct(request); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(reuseable.FormatValidationError(err))
	}

	orgID := c.Params("orgId")

	if err := h.organizationService.AddUserToOrganization(request.UserID, orgID, userID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.UserResponse{
			Status:     "Bad request",
			Message:    "Client error",
			StatusCode: 400,
		})
	}

	return c.Status(fiber.StatusOK).JSON(responses.OtherResponse{
		Status:  "success",
		Message: "User added to organisation successfully",
	})
}
