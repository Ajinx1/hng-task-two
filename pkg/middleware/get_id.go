package middleware

import (
	"hng-task-two/pkg/models"
	"hng-task-two/pkg/reuseable"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

// ////////////////////////////////////////////////////////////
// Validate jwt is valid
// /////////////////////////////////////////////////////////////

func GetUserId(c *fiber.Ctx) string {
	bearer := new(models.Authorization)

	if err := c.ReqHeaderParser(bearer); err != nil {
		return "0"
	}
	var tokenVal string

	if val := strings.Split(bearer.Authorization, " "); len(val) < 2 {
		return "0"
	} else {
		tokenVal = val[1]
	}

	token, _ := jwt.ParseWithClaims(
		tokenVal,
		&models.UserDataClaims{},
		func(*jwt.Token) (interface{}, error) {
			return []byte(reuseable.GetEnvVar("SECRET_KEY")), nil
		})

	claims, _ := token.Claims.(*models.UserDataClaims)

	return claims.UserId

}
