package tests

import (
	"bytes"
	"encoding/json"
	"hng-task-two/pkg/middleware"
	"hng-task-two/pkg/models"
	"hng-task-two/pkg/reuseable"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestRegisterUser_Success(t *testing.T) {
	app := setupTestApp()

	t.Run("RegisterUser_Success", func(t *testing.T) {
		requestBody := models.User{
			UserID:    "742aee44-1d78-4ef5-b1c7-716657550065",
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
			Password:  "password123",
			Phone:     "1234567890",
		}

		requestJSON, _ := json.Marshal(requestBody)
		req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(requestJSON))
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req)
		assert.NoError(t, err, "Request failed")
		assert.Equal(t, fiber.StatusCreated, res.StatusCode, "Expected status code 201")

		var response struct {
			Status  string `json:"status"`
			Message string `json:"message"`
			Data    struct {
				AccessToken string `json:"accessToken"`
				User        struct {
					UserID    string `json:"userId"`
					FirstName string `json:"firstName"`
					LastName  string `json:"lastName"`
					Email     string `json:"email"`
					Phone     string `json:"phone"`
				} `json:"user"`
			} `json:"data"`
		}
		err = json.NewDecoder(res.Body).Decode(&response)
		assert.NoError(t, err, "Failed to decode response")

		assert.Equal(t, "success", response.Status, "Expected success status")
		assert.Equal(t, "Registration successful", response.Message, "Expected 'Registration successful' message")
		assert.NotNil(t, response.Data.AccessToken, "Expected accessToken in response")
		assert.Equal(t, requestBody.UserID, response.Data.User.UserID, "Expected UserID in response")
		assert.Equal(t, requestBody.FirstName, response.Data.User.FirstName, "Expected FirstName in response")
		assert.Equal(t, requestBody.LastName, response.Data.User.LastName, "Expected LastName in response")
		assert.Equal(t, requestBody.Email, response.Data.User.Email, "Expected Email in response")
		assert.Equal(t, requestBody.Phone, response.Data.User.Phone, "Expected Phone in response")
	})
}

func TestDefaultOrganizationNameGeneration(t *testing.T) {
	app := setupTestApp()

	t.Run("DefaultOrganizationNameGeneration", func(t *testing.T) {
		requestBody := models.User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
			Password:  "password123",
			Phone:     "1234567890",
		}

		requestJSON, _ := json.Marshal(requestBody)
		req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(requestJSON))
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req)
		assert.NoError(t, err, "Request failed")
		assert.Equal(t, fiber.StatusCreated, res.StatusCode, "Expected status code 201")

		var response struct {
			Status  string `json:"status"`
			Message string `json:"message"`
			Data    struct {
				User struct {
					OrganisationName string `json:"organisationName"`
				} `json:"user"`
			} `json:"data"`
		}
		err = json.NewDecoder(res.Body).Decode(&response)
		assert.NoError(t, err, "Failed to decode response")

		expectedOrgName := requestBody.FirstName + "'s Organisation"
		assert.Equal(t, expectedOrgName, response.Data.User.OrganisationName, "Expected default organisation name")
	})
}

func TestLoginUser_Success(t *testing.T) {
	app := setupTestApp()

	t.Run("LoginUser_Success", func(t *testing.T) {
		loginData := models.LoginRequest{
			Email:    "john.doe@example.com",
			Password: "password123",
		}

		requestJSON, _ := json.Marshal(loginData)
		req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(requestJSON))
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req)
		assert.NoError(t, err, "Request failed")
		assert.Equal(t, fiber.StatusOK, res.StatusCode, "Expected status code 200")

		var response struct {
			Status  string `json:"status"`
			Message string `json:"message"`
			Data    struct {
				AccessToken string `json:"accessToken"`
				User        struct {
					UserID    string `json:"userId"`
					FirstName string `json:"firstName"`
					LastName  string `json:"lastName"`
					Email     string `json:"email"`
					Phone     string `json:"phone"`
				} `json:"user"`
			} `json:"data"`
		}
		err = json.NewDecoder(res.Body).Decode(&response)
		assert.NoError(t, err, "Failed to decode response")

		assert.Equal(t, "success", response.Status, "Expected success status")
		assert.Equal(t, "Login successful", response.Message, "Expected 'Login successful' message")
		assert.NotNil(t, response.Data.AccessToken, "Expected accessToken in response")
		assert.Equal(t, "John", response.Data.User.FirstName, "Expected FirstName in response")
		assert.Equal(t, "Doe", response.Data.User.LastName, "Expected LastName in response")
		assert.Equal(t, "john.doe@example.com", response.Data.User.Email, "Expected Email in response")
		assert.Equal(t, "1234567890", response.Data.User.Phone, "Expected Phone in response")
	})
}

func TestMissingRequiredFields(t *testing.T) {
	app := setupTestApp()

	t.Run("MissingFirstName", func(t *testing.T) {
		requestBody := models.User{
			// Missing FirstName
			LastName: "Doe",
			Phone:    "1234567890",
			Email:    "john@example.com",
			Password: "password123",
		}

		testValidation(t, app, requestBody, "FirstName")
	})

	t.Run("MissingLastName", func(t *testing.T) {
		requestBody := models.User{
			FirstName: "John",
			Phone:     "1234567890",
			// Missing LastName
			Email:    "john@example.com",
			Password: "password123",
		}

		testValidation(t, app, requestBody, "LastName")
	})

	t.Run("MissingEmail", func(t *testing.T) {
		requestBody := models.User{
			FirstName: "John",
			LastName:  "Doe",
			Phone:     "1234567890",
			// Missing Email
			Password: "password123",
		}

		testValidation(t, app, requestBody, "Email")
	})

	t.Run("MissingPassword", func(t *testing.T) {
		requestBody := models.User{
			FirstName: "John",
			LastName:  "Doe",
			Phone:     "1234567890",
			Email:     "john@example.com",
			// Missing Password
		}

		testValidation(t, app, requestBody, "Password")
	})
}

func testValidation(t *testing.T, app *fiber.App, requestBody models.User, expectedField string) {
	requestJSON, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(requestJSON))
	req.Header.Set("Content-Type", "application/json")

	res, err := app.Test(req)
	assert.NoError(t, err, "Request failed")
	assert.Equal(t, fiber.StatusUnprocessableEntity, res.StatusCode, "Expected status code 422")

	var response struct {
		Errors []struct {
			Field   string `json:"field"`
			Message string `json:"message"`
		} `json:"errors"`
	}
	err = json.NewDecoder(res.Body).Decode(&response)
	assert.NoError(t, err, "Failed to decode response")

	assert.Len(t, response.Errors, 1, "Expected 1 error message")
	assert.Equal(t, expectedField, response.Errors[0].Field, "Expected error for "+expectedField+" field")
}

func TestDuplicateEmail(t *testing.T) {
	app := setupTestApp()

	t.Run("DuplicateEmail", func(t *testing.T) {
		requestBody := models.User{
			FirstName: "Jane",
			LastName:  "Doe",
			Email:     "jane.doe@example.com", // Already registered email
			Password:  "password456",
			Phone:     "0987654321",
		}

		requestJSON, _ := json.Marshal(requestBody)
		req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(requestJSON))
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req)
		assert.NoError(t, err, "Request failed")
		assert.Equal(t, fiber.StatusUnprocessableEntity, res.StatusCode, "Expected status code 422")

		var response struct {
			Errors []struct {
				Field   string `json:"field"`
				Message string `json:"message"`
			} `json:"errors"`
		}
		err = json.NewDecoder(res.Body).Decode(&response)
		assert.NoError(t, err, "Failed to decode response")

		assert.Len(t, response.Errors, 1, "Expected 1 error message")
		assert.Equal(t, "email", response.Errors[0].Field, "Expected error for Email field")
		assert.Equal(t, "email exist", response.Errors[0].Message, "Expected error message 'email exist'")
	})
}

func TestOrganizationAccess(t *testing.T) {
	app := setupTestApp()

	var tokenUser2 string
	var tokenUser1 string

	// Register User1 and generate token
	t.Run("RegisterUser1", func(t *testing.T) {
		requestBody := models.User{
			UserID:    "user1_unique_id",
			FirstName: "Janey",
			LastName:  "Doe",
			Email:     "janey.doe@example.com",
			Password:  "password123",
			Phone:     "1234567890",
		}

		requestJSON, _ := json.Marshal(requestBody)
		req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(requestJSON))
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req)
		assert.NoError(t, err, "Request failed")
		assert.Equal(t, fiber.StatusCreated, res.StatusCode, "Expected status code 201")

		// Generate a mock token for User1
		user1 := models.User{
			UserID:    "user1_unique_id",
			FirstName: "Jane",
			LastName:  "Doe",
			Email:     "jane.doe@example.com",
		}
		tokenUser1, err = middleware.MockGenerateJWT(&user1)
		assert.NoError(t, err, "Failed to generate mock token for User1")
	})

	// Register User2 and generate token
	t.Run("RegisterUser2", func(t *testing.T) {
		requestBody := models.User{
			UserID:    "user2_unique_id",
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
			Password:  "password123",
			Phone:     "1234567890",
		}

		requestJSON, _ := json.Marshal(requestBody)
		req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(requestJSON))
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req)
		assert.NoError(t, err, "Request failed")
		assert.Equal(t, fiber.StatusCreated, res.StatusCode, "Expected status code 201")

		// Generate a mock token for User2
		user2 := models.User{
			UserID:    "user2_unique_id",
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
		}
		tokenUser2, err = middleware.MockGenerateJWT(&user2)
		assert.NoError(t, err, "Failed to generate mock token for User2")
	})

	// User2 tries to access User1's organization
	t.Run("User2CannotAccessUser1Organization", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/organisations/bec0c2bf-dfa4-45a6-88d2-9dad2814e949", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+tokenUser2)

		res, err := app.Test(req)
		assert.NoError(t, err, "Request failed")
		assert.Equal(t, fiber.StatusForbidden, res.StatusCode, "Expected status code 403")

		var response struct {
			Status     string `json:"status"`
			Message    string `json:"message"`
			StatusCode int    `json:"statusCode"`
		}
		err = json.NewDecoder(res.Body).Decode(&response)
		assert.NoError(t, err, "Failed to decode response")

		assert.Equal(t, "Forbidden", response.Status, "Expected status Forbidden")
		assert.Equal(t, "You do not have access to this organization", response.Message, "Expected forbidden message")
		assert.Equal(t, 403, response.StatusCode, "Expected status code 403")
	})

	// User1 tries to access their own organization
	t.Run("User1CanAccessTheirOwnOrganization", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/organisations/bec0c2bf-dfa4-45a6-88d2-9dad2814e949", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+tokenUser1)

		res, err := app.Test(req)
		assert.NoError(t, err, "Request failed")
		assert.Equal(t, fiber.StatusOK, res.StatusCode, "Expected status code 200")

		var response struct {
			Status     string `json:"status"`
			Message    string `json:"message"`
			StatusCode int    `json:"statusCode"`
		}
		err = json.NewDecoder(res.Body).Decode(&response)
		assert.NoError(t, err, "Failed to decode response")

		assert.Equal(t, "success", response.Status, "Expected status success")
		assert.Equal(t, "Organization details fetched successfully", response.Message, "Expected success message")
		assert.Equal(t, 200, res.StatusCode, "Expected status code 200")
	})
}

func setupTestApp() *fiber.App {
	app := fiber.New()

	h := &handler{
		userService: &mockUserService{},
		orgService:  &mockOrganizationService{},
	}

	app.Post("/auth/register", h.registerUser)
	app.Post("/auth/login", h.loginUser)
	app.Get("/api/organisations/:orgID", h.getOrganization)

	return app
}

type handler struct {
	userService userService
	orgService  organizationService
}

type userService interface {
	DupUserByEmail(email string) (*models.User, error)
}

type organizationService interface {
	GetOrganizationByID(orgID, userID string) (*models.Organization, error)
}

type mockUserService struct{}

func (m *mockUserService) DupUserByEmail(email string) (*models.User, error) {
	if email == "jane.doe@example.com" {
		return &models.User{Email: email}, nil
	}
	return nil, nil
}

type mockOrganizationService struct{}

func (m *mockOrganizationService) GetOrganizationByID(orgID, userID string) (*models.Organization, error) {
	if orgID == "bec0c2bf-dfa4-45a6-88d2-9dad2814e949" && userID == "user1_unique_id" {
		return &models.Organization{
			OrgID:       orgID,
			Name:        "Janey's Organization",
			Description: "A test organization",
		}, nil
	}
	return nil, fiber.ErrForbidden
}

func (h *handler) registerUser(c *fiber.Ctx) error {
	var validate = validator.New()
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":     "Bad request",
			"message":    "Registration unsuccessful",
			"statusCode": 400,
		})
	}

	if err := validate.Struct(user); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(reuseable.FormatValidationError(err))
	}

	existingUser, err := h.userService.DupUserByEmail(user.Email)
	if err == nil && existingUser != nil {
		errors := []map[string]string{
			{
				"field":   "email",
				"message": "email exist",
			},
		}
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"errors": errors,
		})
	}

	userToken, _ := middleware.MockGenerateJWT(&user)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Registration successful",
		"data": fiber.Map{
			"accessToken": userToken,
			"user": fiber.Map{
				"userId":           user.UserID,
				"firstName":        user.FirstName,
				"lastName":         user.LastName,
				"email":            user.Email,
				"phone":            user.Phone,
				"organisationName": user.FirstName + "'s Organisation",
			},
		},
	})
}

func (h *handler) loginUser(c *fiber.Ctx) error {
	var validate = validator.New()
	var loginData models.LoginRequest
	if err := c.BodyParser(&loginData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":     "Bad request",
			"message":    "Authentication failed",
			"statusCode": 400,
		})
	}

	if err := validate.Struct(loginData); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(reuseable.FormatValidationError(err))
	}

	if loginData.Email == "john.doe@example.com" && loginData.Password == "password123" {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"message": "Login successful",
			"data": fiber.Map{
				"accessToken": "mock_access_token",
				"user": fiber.Map{
					"userId":    "user1_unique_id",
					"firstName": "John",
					"lastName":  "Doe",
					"email":     "john.doe@example.com",
					"phone":     "1234567890",
				},
			},
		})
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":     "Unauthorized",
			"message":    "Invalid credentials",
			"statusCode": 401,
		})
	}
}

func (h *handler) getOrganization(c *fiber.Ctx) error {
	orgID := c.Params("orgID")
	authHeader := c.Get("Authorization")
	tokenString := authHeader[len("Bearer "):]

	claims, err := middleware.MockParseJWT(tokenString)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":     "Unauthorized",
			"message":    "Invalid token",
			"statusCode": 401,
		})
	}

	org, err := h.orgService.GetOrganizationByID(orgID, claims.UserId)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":     "Forbidden",
			"message":    "You do not have access to this organization",
			"statusCode": 403,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Organization details fetched successfully",
		"data":    org,
	})
}
