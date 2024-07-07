package main

import (
	v1 "hng-task-two/api"
	"hng-task-two/pkg/database"
	"hng-task-two/pkg/models"
	"log"
	"os"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	/////////////////////////////////////////////
	//Initiate Fiber App
	/////////////////////////////////////////////
	app := fiber.New(fiber.Config{
		Prefork:      false,
		AppName:      "HNG Task Two Backend",
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
		ServerHeader: "Golang Fiber",
	})

	// Create a new middleware handler "c.NEXT()"
	app.Use(recover.New())

	// Loggers
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${time} - ${status} - ${latency} - ${method} ${path}\n",
	}))

	// Define the CSP policy
	cspPolicy := "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'"

	// Middleware to set CSP headers
	app.Use(func(c *fiber.Ctx) error {
		c.Set("Content-Security-Policy", cspPolicy)
		return c.Next()
	})

	// Migrate the schemas
	migrate_err := database.DB.AutoMigrate(&models.User{}, &models.Organization{})
	if migrate_err != nil {
		log.Fatal("failed to auto-migrate database:", migrate_err)
	}

	/////////////////////////////////////////////
	//All Routes Registered
	/////////////////////////////////////////////
	v1.RegisterAllRoutes(app, database.DB)

	/////////////////////////////////////////////
	//Host and port config
	/////////////////////////////////////////////
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")

	if port == "" {
		port = "9090"
	}

	err := app.Listen(host + ":" + port)
	if err != nil {
		panic(err)
	}

}
