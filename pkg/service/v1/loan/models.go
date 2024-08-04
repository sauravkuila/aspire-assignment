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
	LoanId       int64                `json:"loanId"`
	UserId       int64                `json:"userId,omitempty"`
	UserName     string               `json:"username,omitempty"`
	Amount       float64              `json:"amount,omitempty"`
	Installments int64                `json:"installments,omitempty"`
	Status       string               `json:"status"`
	Details      []InstallmentDetails `json:"details,omitempty"`
	CreatedAt    string               `json:"createdAt,omitempty"`
}

type InstallmentDetails struct {
	LoanId            int64   `json:"loanId,omitempty"`
	AmoundDue         float64 `json:"amountDue,omitempty"`
	AmountPaid        float64 `json:"amountPaid,omitempty"`
	Status            string  `json:"status,omitempty"`
	InstallmentNumber int64   `json:"installmentNumber,omitempty"`
	TransactionId     string  `json:"transactionId,omitempty"`
	DueDate           string  `json:"dueDate,omitempty"`
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

type GetLoanDetailRequest struct {
	UserId int64 `form:"userId" binding:"required"`
	LoanId int64 `form:"loanId" binding:"required"`
}

type GetLoanDetailResponse struct {
	Data    *GetLoanDetail `json:"data,omitempty"`
	Status  bool           `json:"success"`
	Errors  []e.Error      `json:"errors,omitempty"`
	Message string         `json:"message,omitempty"`
}

type GetLoanDetail struct {
	LoanId            int64                `json:"loanId"`
	LoanAmount        float64              `json:"loanAmount,omitempty"`
	OutstandingAmount float64              `json:"outstandingAmount,omitempty"`
	TotalInstallments int                  `json:"totalInstallments,omitempty"`
	Status            string               `json:"status"`
	Installments      []InstallmentDetails `json:"installments,omitempty"`
}

type ProcessLoanPaymentRequest struct {
	UserId        int64   `json:"userId" binding:"required"`
	LoanId        int64   `json:"loanId" binding:"required"`
	Amount        float64 `json:"amount" binding:"required"`
	TransactionId string  `json:"transactionId" binding:"required"`
}

type ProcessLoanPaymentResponse struct {
	Data    interface{} `json:"data,omitempty"`
	Status  bool        `json:"success"`
	Errors  []e.Error   `json:"errors,omitempty"`
	Message string      `json:"message,omitempty"`
}
