package handler

import (
	"hng-task-two/pkg/models"
	"hng-task-two/pkg/responses"
	"hng-task-two/pkg/reuseable"

	"github.com/gofiber/fiber/v2"
)

func (h *UserHandler) Login(c *fiber.Ctx) error {

	var request models.LoginRequest

	if err := c.BodyParser(&request); err != nil {

		return c.Status(fiber.StatusBadRequest).JSON(responses.UserResponse{
			Status:     "Bad request",
			Message:    "Authentication failed",
			StatusCode: 401,
		})

	}

	if err := h.validate.Struct(request); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(reuseable.FormatValidationError(err))
	}

	user, token, err := h.userService.Login(request.Email, request.Password)
	if err != nil {

		return c.Status(fiber.StatusUnauthorized).JSON(responses.UserResponse{
			Status:     "Bad request",
			Message:    "Authentication failed",
			StatusCode: 401,
		})
	}

	userResponse := models.ViewUser{
		UserID:    user.UserID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Phone:     user.Phone,
	}

	return c.Status(fiber.StatusOK).JSON(responses.LoginResponse{
		Status:  "success",
		Message: "Login successful",
		Data: &fiber.Map{
			"accessToken": token,
			"user":        userResponse,
		},
	})
}
