package api

import (
	org_handler "hng-task-two/internal/organization/handler"
	org_repo "hng-task-two/internal/organization/repository"
	org_service "hng-task-two/internal/organization/service"
	"hng-task-two/internal/user/handler"
	"hng-task-two/internal/user/repository"
	"hng-task-two/internal/user/service"
	"hng-task-two/pkg/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterAllRoutes(app fiber.Router, db *gorm.DB) {

	userRepo := repository.NewUserRepository(db)
	orgRepo := org_repo.NewOrganizationRepository(db)

	userService := service.NewUserService(userRepo, orgRepo)
	orgService := org_service.NewOrganizationService(orgRepo, userRepo)

	userHandler := handler.NewUserHandler(userService)
	orgHandler := org_handler.NewOrganizationHandler(orgService)

	//Users routes
	authRoutes := app.Group("/auth")
	authRoutes.Post("/register", userHandler.Register)
	authRoutes.Post("/login", userHandler.Login)
	app.Get("/api/users/:id", middleware.ValidateJWT, userHandler.GetAUser)

	//Organizations routes
	orgRoutes := app.Group("/api/organisations", middleware.ValidateJWT)
	orgRoutes.Post("/", orgHandler.Create)
	orgRoutes.Get("/", orgHandler.GetAllOrganizations)
	orgRoutes.Get("/:orgId", orgHandler.GetOrganization)
	app.Post("/api/organisations/:orgId/users", orgHandler.AddUserToOrganization)
}
