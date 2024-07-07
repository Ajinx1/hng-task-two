package middleware

import (
	"hng-task-two/pkg/models"
	"hng-task-two/pkg/reuseable"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// ////////////////////////////////////////////////////////////
// Generate JWT
// /////////////////////////////////////////////////////////////

func GenerateJWT(user *models.User) (tokenString string, err error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(10 * time.Hour)

	the_claims := models.UserDataClaims{
		UserId:    user.UserID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Phone:     user.Phone,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    user.Email,
			ExpiresAt: jwt.NewNumericDate(expireTime),
		},
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS512, the_claims)
	token, _ := claims.SignedString([]byte(reuseable.GetEnvVar("SECRET_KEY")))
	return token, err

}

// ////////////////////////////////////////////////////////////
// Mock Generate JWT
// /////////////////////////////////////////////////////////////

func MockGenerateJWT(user *models.User) (tokenString string, err error) {

	expireTime := time.Now().Add(1 * time.Hour)

	mockClaims := models.UserDataClaims{
		UserId:    user.UserID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Phone:     user.Phone,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    user.Email,
			ExpiresAt: jwt.NewNumericDate(expireTime),
		},
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS512, &mockClaims)
	token, _ := claims.SignedString([]byte("mock_secret_key"))

	return token, nil
}
