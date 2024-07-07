package models

import "github.com/golang-jwt/jwt/v4"

type Authorization struct {
	Authorization string `reqHeader:"Authorization"`
}

type UserDataClaims struct {
	UserId    string
	FirstName string
	LastName  string
	Email     string
	Phone     string
	jwt.RegisteredClaims
}
