package auth

import (
	e "aspire-assignment/pkg/errors"
	"time"

	"github.com/golang-jwt/jwt"
)

type AuthResponse struct {
	Data    interface{} `json:"data,omitempty"`
	Status  bool        `json:"success"`
	Errors  []e.Error   `json:"errors,omitempty"`
	Message string      `json:"message,omitempty"`
}

type Token struct {
	UserName string    `json:"username"`
	UserId   int64     `json:"userId"`
	UserType string    `json:"userType"`
	Exp      time.Time `json:"expiry"`
}

type Claims struct {
	Payload Token `json:"token"`
	jwt.StandardClaims
}
