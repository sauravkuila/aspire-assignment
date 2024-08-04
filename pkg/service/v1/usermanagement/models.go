package usermanagement

import (
	e "aspire-assignment/pkg/errors"
)

type UserSignupRequest struct {
	UserName      string  `json:"username" binding:"required"`
	Password      string  `json:"password" binding:"required,min=6"`
	UserType      string  `json:"type" binding:"required,oneof=CUSTOMER ADMIN"`
	Email         string  `json:"email" binding:"required"`
	Mobile        string  `json:"mobile" binding:"required"`
	MonthlySalary float64 `json:"salary"`
	BankBalance   float64 `json:"bankBalance"`
}

type UserSignupResponse struct {
	Data    *UserSignup `json:"data,omitempty"`
	Status  bool        `json:"success"`
	Errors  []e.Error   `json:"errors,omitempty"`
	Message string      `json:"message,omitempty"`
}

type UserSignup struct {
	UserName string `json:"username"`
	UserId   int64  `json:"userId"`
}

type UserLoginRequest struct {
	UserName string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

type UserLoginResponse struct {
	Data    *UserLogin `json:"data,omitempty"`
	Status  bool       `json:"success"`
	Errors  []e.Error  `json:"errors,omitempty"`
	Message string     `json:"message,omitempty"`
}

type UserLogin struct {
	Token  string `json:"token"`
	Expiry string `json:"expiry"`
}
