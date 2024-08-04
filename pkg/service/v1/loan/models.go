package loan

import (
	e "aspire-assignment/pkg/errors"
)

type CreateLoanRequest struct {
	UserId       int64   `json:"userId" binding:"required"`
	Amount       float64 `json:"amount" binding:"required"`
	Installments int64   `json:"installments" binding:"required"`
}

type CreateLoanResponse struct {
	Data    *LoanDetails `json:"data,omitempty"`
	Status  bool         `json:"success"`
	Errors  []e.Error    `json:"errors,omitempty"`
	Message string       `json:"message,omitempty"`
}

type LoanDetails struct {
	LoanId       int64   `json:"loanId"`
	UserId       int64   `json:"userId,omitempty"`
	UserName     string  `json:"username,omitempty"`
	Amount       float64 `json:"amount,omitempty"`
	Installments int64   `json:"installments,omitempty"`
	Status       string  `json:"status"`
	CreatedAt    string  `json:"createdAt,omitempty"`
}

type ModifyLoanRequest struct {
	UserId       int64   `json:"userId" binding:"required"`
	LoanId       int64   `json:"loanId" binding:"required"`
	Amount       float64 `json:"amount" binding:"required"`
	Installments int64   `json:"installments" binding:"required"`
}

type ModifyLoanResponse struct {
	Data    *LoanDetails `json:"data,omitempty"`
	Status  bool         `json:"success"`
	Errors  []e.Error    `json:"errors,omitempty"`
	Message string       `json:"message,omitempty"`
}

type CancelLoanRequest struct {
	UserId int64 `json:"userId" binding:"required"`
	LoanId int64 `json:"loanId" binding:"required"`
}

type CancelLoanResponse struct {
	Data    *LoanDetails `json:"data,omitempty"`
	Status  bool         `json:"success"`
	Errors  []e.Error    `json:"errors,omitempty"`
	Message string       `json:"message,omitempty"`
}

type GetLoanRequest struct {
	UserId int64 `form:"userId" binding:"required"`
}

type GetLoanResponse struct {
	Data    []LoanDetails `json:"data,omitempty"`
	Status  bool          `json:"success"`
	Errors  []e.Error     `json:"errors,omitempty"`
	Message string        `json:"message,omitempty"`
}

type PendingLoanRequest struct {
	UserId int64 `form:"userId" binding:"required"`
}

type PendingLoanResponse struct {
	Data    []LoanDetails `json:"data,omitempty"`
	Status  bool          `json:"success"`
	Errors  []e.Error     `json:"errors,omitempty"`
	Message string        `json:"message,omitempty"`
}

type ApproveRejectLoanApplicationRequest struct {
	UserId   int64  `json:"userId" binding:"required"`
	LoanId   int64  `json:"loanId" binding:"required"`
	Approval string `json:"approval" binding:"required,oneof=APPROVE REJECT"`
}

type ApproveRejectLoanApplicationResponse struct {
	Data    []LoanDetails `json:"data,omitempty"`
	Status  bool          `json:"success"`
	Errors  []e.Error     `json:"errors,omitempty"`
	Message string        `json:"message,omitempty"`
}
