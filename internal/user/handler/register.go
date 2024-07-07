package handler

import (
	"hng-task-two/pkg/models"
	"hng-task-two/pkg/responses"
	"hng-task-two/pkg/reuseable"

	"github.com/gofiber/fiber/v2"
)

func (h *UserHandler) Register(c *fiber.Ctx) error {

	var newUser = new(models.User)

	if err := c.BodyParser(&newUser); err != nil {

		return c.Status(fiber.StatusBadRequest).JSON(responses.UserResponse{
			Status:     "Bad request",
			Message:    "Registration unsuccessful",
			StatusCode: 400,
		})
	}

	if err := h.validate.Struct(newUser); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(reuseable.FormatValidationError(err))
	}

	existingUser, err := h.userService.GetUserByEmail(newUser.Email)
	if err == nil && existingUser != nil {
		errorMsg := []map[string]string{
			{
				"field":   "email",
				"message": "email exist",
			},
		}
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"errors": errorMsg,
		})
	}

	user, token, err := h.userService.Register(newUser)
	if err != nil {

		return c.Status(fiber.StatusBadRequest).JSON(responses.UserResponse{
			Status:     "Bad request",
			Message:    "Registration unsuccessful",
			StatusCode: 400,
		})
	}

	userResponse := models.ViewUser{
		UserID:    user.UserID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Phone:     user.Phone,
	}

	return c.Status(fiber.StatusCreated).JSON(responses.LoginResponse{
		Status:  "success",
		Message: "Registration successful",
		Data: &fiber.Map{
			"accessToken": token,
			"user":        userResponse,
		},
	})

}
