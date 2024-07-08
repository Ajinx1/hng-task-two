package handler

import (
	"hng-task-two/pkg/middleware"
	"hng-task-two/pkg/models"
	"hng-task-two/pkg/responses"

	"github.com/gofiber/fiber/v2"
)

func (h *UserHandler) GetAUser(c *fiber.Ctx) error {
	ID := c.Params("id")
	userID := middleware.GetUserId(c)
	user, err := h.userService.GetUserByIDs(ID, userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.UserResponse{
			Status:     "Bad request",
			Message:    "Client error",
			StatusCode: 400,
		})
	}

	userResponse := &models.ViewUser{
		UserID:    user.UserID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Phone:     user.Phone,
	}

	return c.Status(fiber.StatusOK).JSON(responses.UResponse{
		Status:  "success",
		Message: "<message>",
		Data:    userResponse,
	})

}
