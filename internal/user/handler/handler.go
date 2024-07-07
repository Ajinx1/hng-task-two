package handler

import (
	"hng-task-two/internal/user/service"

	"github.com/go-playground/validator"
)

type UserHandler struct {
	userService service.UserService
	validate    *validator.Validate
}

func NewUserHandler(userService service.UserService) *UserHandler {

	return &UserHandler{
		userService: userService,
		validate:    validator.New(),
	}
}
