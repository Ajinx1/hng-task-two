package tests

import (
	"bytes"
	"encoding/json"
	"hng-task-two/pkg/middleware"
	"hng-task-two/pkg/models"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestGenerateJWT(t *testing.T) {
	app := setupTestApp()

	t.Run("TokenGeneration and Expiry", func(t *testing.T) {

		userID := "unique_id"
		email := "jane.doe1@example.com"
		firstName := "John"
		lastName := "Doe"
		password := "password123"
		phone := "1234567890"

		requestBody := models.User{
			UserID:    userID,
			FirstName: firstName,
			LastName:  lastName,
			Email:     email,
			Password:  password,
			Phone:     phone,
		}

		requestJSON, _ := json.Marshal(requestBody)
		req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(requestJSON))
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req)

		assert.NoError(t, err, "Request failed")
		assert.Equal(t, fiber.StatusCreated, res.StatusCode, "Expected status code 201")

		var response struct {
			Data struct {
				AccessToken string `json:"accessToken"`
			} `json:"data"`
		}
		err = json.NewDecoder(res.Body).Decode(&response)
		assert.NoError(t, err, "Failed to decode response")

		token := response.Data.AccessToken

		claims, err := middleware.MockParseJWT(token)
		assert.NoError(t, err, "Failed to parse JWT token")

		assert.Equal(t, userID, claims.UserId, "Expected user ID in token claims")
		assert.Equal(t, firstName, claims.FirstName, "Expected firstName in token claims")
		assert.Equal(t, lastName, claims.LastName, "Expected lastName in token claims")
		assert.Equal(t, email, claims.Email, "Expected email in token claims")
		assert.Equal(t, phone, claims.Phone, "Expected phone in token claims")

		expirationTime := claims.ExpiresAt.Time
		expectedExpirationTime := time.Now().Add(1 * time.Hour)
		assert.WithinDuration(t, expectedExpirationTime, expirationTime, 5*time.Second, "Token expiration time mismatch")
	})
}
