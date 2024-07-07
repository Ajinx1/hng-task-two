package reuseable

import (
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

func FormatValidationError(err error) *fiber.Map {
	var errors []map[string]string
	for _, err := range err.(validator.ValidationErrors) {
		errors = append(errors, map[string]string{
			"field":   err.Field(),
			"message": err.Tag() + " validation failed",
		})
	}
	return &fiber.Map{
		"errors": errors,
	}
}
