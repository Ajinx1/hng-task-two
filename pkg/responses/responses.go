package responses

import (
	"hng-task-two/pkg/models"

	"github.com/gofiber/fiber/v2"
)

type ValidationResponse struct {
	Errors interface{} `json:"errors"`
}

type UserResponse struct {
	Status     string `json:"status"`
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
}

type OtherResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type JwtResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

type LoginResponse struct {
	Status  string     `json:"status"`
	Message string     `json:"message"`
	Data    *fiber.Map `json:"data"`
}

type OrgResponse struct {
	Status  string                   `json:"status"`
	Message string                   `json:"message"`
	Data    *models.ViewOrganization `json:"data"`
}

type AllOrgResponse struct {
	Status  string     `json:"status"`
	Message string     `json:"message"`
	Data    *fiber.Map `json:"data"`
}

type UResponse struct {
	Status  string           `json:"status"`
	Message string           `json:"message"`
	Data    *models.ViewUser `json:"data"`
}
