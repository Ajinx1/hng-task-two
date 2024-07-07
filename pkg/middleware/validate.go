package middleware

import (
	"errors"
	"hng-task-two/pkg/models"
	"hng-task-two/pkg/responses"
	"hng-task-two/pkg/reuseable"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

// ////////////////////////////////////////////////////////////
// Validate jwt is valid
// /////////////////////////////////////////////////////////////

func ValidateJWT(c *fiber.Ctx) error {
	bearer := new(models.Authorization)

	if err := c.ReqHeaderParser(bearer); err != nil {
		return c.Status(403).JSON(responses.JwtResponse{Status: false, Message: noMsg})
	}
	var tokenVal string

	if val := strings.Split(bearer.Authorization, " "); len(val) < 2 {
		return c.Status(403).JSON(responses.JwtResponse{Status: false, Message: invalidMsg})
	} else {
		tokenVal = val[1]
	}

	token, err := jwt.ParseWithClaims(
		tokenVal,
		&models.UserDataClaims{},
		func(*jwt.Token) (interface{}, error) {
			return []byte(reuseable.GetEnvVar("SECRET_KEY")), nil
		})

	claims, ok := token.Claims.(*models.UserDataClaims)
	if claims.UserId == "" || !ok {
		return c.Status(403).JSON(responses.JwtResponse{Status: false, Message: invalidUser})
	}

	if err != nil {
		return c.Status(403).JSON(responses.JwtResponse{Status: false, Message: expireMsg})
	}

	return c.Next()

}

// ////////////////////////////////////////////////////////////
// Parse jwt
// /////////////////////////////////////////////////////////////

func ParseJWT(tokenVal string) (*models.UserDataClaims, error) {

	token, err := jwt.ParseWithClaims(
		tokenVal,
		&models.UserDataClaims{},
		func(*jwt.Token) (interface{}, error) {
			return []byte(reuseable.GetEnvVar("SECRET_KEY")), nil
		})

	claims, ok := token.Claims.(*models.UserDataClaims)
	if claims.UserId == "" || !ok {
		return claims, err
	}

	if err != nil {
		return claims, err
	}

	return claims, nil

}

// MockParseJWT simulates parsing of JWT tokens for testing purposes.
func MockParseJWT(tokenString string) (*models.UserDataClaims, error) {

	secretKey := []byte("mock_secret_key")

	token, err := jwt.ParseWithClaims(tokenString, &models.UserDataClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*models.UserDataClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
