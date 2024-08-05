package loan

import (
	"aspire-assignment/pkg/config"
	v1 "aspire-assignment/pkg/db/v1"
	"aspire-assignment/pkg/db/v1/loan"
	dbmock "aspire-assignment/pkg/db/v1/mock"
	e "aspire-assignment/pkg/errors"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
)

func Test_loanService_GetPendingLoans(t *testing.T) {
	var (
		dbObj  v1.V1DBLayer
		userId int64 = 1
	)

	//init error to be used in function
	e.ErrorInit()

	tests := []struct {
		name           string
		httpMethod     string
		httpStatus     int
		input          PendingLoanRequest
		setup          func(*gin.Context, PendingLoanRequest)
		expectedOutput PendingLoanResponse
		actualOutput   PendingLoanResponse
	}{
		{
			name:  "FailToGetLoan",
			input: PendingLoanRequest{},
			setup: func(c *gin.Context, data PendingLoanRequest) {
				ctrl := gomock.NewController(t)
				repo := dbmock.NewMockV1DBLayer(ctrl)
				dbObj = repo
				repo.EXPECT().GetUnapprovedLoans(c).Return(nil, fmt.Errorf("db error")).Times(1)
			},
			expectedOutput: PendingLoanResponse{
				Status: false,
				Errors: []e.Error{{
					ErrName:     e.ErrorInfo[e.GetDBError].ErrName,
					Description: e.ErrorInfo[e.GetDBError].Description,
					Code:        e.ErrorInfo[e.GetDBError].Code,
				}},
				Message: "failed to fetch loans",
			},
			httpStatus: http.StatusInternalServerError,
			httpMethod: http.MethodGet,
		},
		{
			name:  "NoLoansAvailable",
			input: PendingLoanRequest{},
			setup: func(c *gin.Context, data PendingLoanRequest) {
				ctrl := gomock.NewController(t)
				repo := dbmock.NewMockV1DBLayer(ctrl)
				dbObj = repo
				repo.EXPECT().GetUnapprovedLoans(c).Return(nil, nil).Times(1)
			},
			expectedOutput: PendingLoanResponse{
				Status:  false,
				Message: "no loans available",
			},
			httpStatus: http.StatusNotFound,
			httpMethod: http.MethodGet,
		},
		{
			name:  "SuccessGetLoans",
			input: PendingLoanRequest{},
			setup: func(c *gin.Context, data PendingLoanRequest) {
				ctrl := gomock.NewController(t)
				repo := dbmock.NewMockV1DBLayer(ctrl)
				dbObj = repo
				loans := make([]loan.UnApprovedLoan, 0)
				t1, _ := time.Parse("2006-01-02 15:04:05", "2024-08-08 15:00:00")
				loans = append(loans, loan.UnApprovedLoan{
					LoanId:       sql.NullInt64{Int64: 3, Valid: true},
					UserName:     sql.NullString{String: "testuser", Valid: true},
					Amount:       sql.NullFloat64{Float64: 34000, Valid: true},
					Installments: sql.NullInt64{Int64: 3, Valid: true},
					Status:       sql.NullString{String: LOAN_PENDING, Valid: true},
					CreatedAt:    sql.NullTime{Time: t1, Valid: true},
				})
				repo.EXPECT().GetUnapprovedLoans(c).Return(loans, nil).Times(1)
			},
			expectedOutput: PendingLoanResponse{
				Status: true,
				Data: []LoanDetails{{
					LoanId:    3,
					UserName:  "testuser",
					Amount:    34000,
					Tenure:    3,
					Status:    LOAN_PENDING,
					CreatedAt: "2024-08-08 15:00:00",
				}},
				Message: "uccessfully fetched unapproved loans",
			},
			httpStatus: http.StatusOK,
			httpMethod: http.MethodGet,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			fmt.Println("Starting Create Loan TestCase: ", tt.name)
			w, ctx := getContext(tt.httpMethod, tt.input, nil, nil)
			ctx.Set(config.USERID, userId)

			//setup test
			tt.setup(ctx, tt.input)
			servObj := NewLoanService(dbObj)

			//calling the function
			servObj.GetPendingLoans(ctx)

			//check for result status
			assert.Equal(t, tt.httpStatus, w.Code)

			//create a copy of the output structure
			err := json.Unmarshal(w.Body.Bytes(), &tt.actualOutput)
			if err != nil {
				t.Error("unable to unmarshal response")
			}

			//compare expected vs actual output
			assert.Equal(t, tt.expectedOutput.Status, tt.actualOutput.Status)
			if len(tt.expectedOutput.Errors) != 0 {
				assert.Equal(t, tt.expectedOutput.Errors[0].Code, tt.actualOutput.Errors[0].Code)
			}

			fmt.Println("Ending Create Loan TestCase: ", tt.name)
		})
	}
}

func Test_loanService_ApproveRejectLoanApplication(t *testing.T) {
	var (
		dbObj  v1.V1DBLayer
		userId int64 = 1
	)

	//init error to be used in function
	e.ErrorInit()

	tests := []struct {
		name           string
		httpMethod     string
		httpStatus     int
		input          ApproveRejectLoanApplicationRequest
		setup          func(*gin.Context, ApproveRejectLoanApplicationRequest)
		expectedOutput ApproveRejectLoanApplicationResponse
		actualOutput   ApproveRejectLoanApplicationResponse
	}{
		{
			name: "MissingInputLoanId",
			input: ApproveRejectLoanApplicationRequest{
				Approval: LOAN_APPROVE,
			},
			setup: func(c *gin.Context, data ApproveRejectLoanApplicationRequest) {
				ctrl := gomock.NewController(t)
				repo := dbmock.NewMockV1DBLayer(ctrl)
				dbObj = repo
			},
			expectedOutput: ApproveRejectLoanApplicationResponse{
				Status: false,
				Errors: []e.Error{{
					ErrName:     e.ErrorInfo[e.BadRequest].ErrName,
					Description: e.ErrorInfo[e.BadRequest].Description,
					Code:        e.ErrorInfo[e.BadRequest].Code,
				}},
				Message: "failed to update loan approval",
			},
			httpStatus: http.StatusBadRequest,
			httpMethod: http.MethodPost,
		},
		{
			name: "MissingInputApproval",
			input: ApproveRejectLoanApplicationRequest{
				LoanId: 3,
			},
			setup: func(c *gin.Context, data ApproveRejectLoanApplicationRequest) {
				ctrl := gomock.NewController(t)
				repo := dbmock.NewMockV1DBLayer(ctrl)
				dbObj = repo
			},
			expectedOutput: ApproveRejectLoanApplicationResponse{
				Status: false,
				Errors: []e.Error{{
					ErrName:     e.ErrorInfo[e.BadRequest].ErrName,
					Description: e.ErrorInfo[e.BadRequest].Description,
					Code:        e.ErrorInfo[e.BadRequest].Code,
				}},
				Message: "failed to update loan approval",
			},
			httpStatus: http.StatusBadRequest,
			httpMethod: http.MethodPost,
		},
		{
			name: "ErrorFetchingLoans",
			input: ApproveRejectLoanApplicationRequest{
				LoanId:   3,
				Approval: LOAN_APPROVE,
			},
			setup: func(c *gin.Context, data ApproveRejectLoanApplicationRequest) {
				ctrl := gomock.NewController(t)
				repo := dbmock.NewMockV1DBLayer(ctrl)
				dbObj = repo
				repo.EXPECT().FetchLoanDetails(c, data.LoanId).Return(loan.LoanDetails{}, fmt.Errorf("db error")).Times(1)
			},
			expectedOutput: ApproveRejectLoanApplicationResponse{
				Status: false,
				Errors: []e.Error{{
					ErrName:     e.ErrorInfo[e.NoDataFound].ErrName,
					Description: e.ErrorInfo[e.NoDataFound].Description,
					Code:        e.ErrorInfo[e.NoDataFound].Code,
				}},
				Message: "failed to fetch loan detail",
			},
			httpStatus: http.StatusNotFound,
			httpMethod: http.MethodPost,
		},
		{
			name: "NotPendingLoan",
			input: ApproveRejectLoanApplicationRequest{
				LoanId:   3,
				Approval: LOAN_APPROVE,
			},
			setup: func(c *gin.Context, data ApproveRejectLoanApplicationRequest) {
				ctrl := gomock.NewController(t)
				repo := dbmock.NewMockV1DBLayer(ctrl)
				dbObj = repo
				loanDetail := loan.LoanDetails{
					LoanId: sql.NullInt64{Int64: data.LoanId, Valid: true},
					Status: sql.NullString{String: LOAN_APPROVE, Valid: true},
				}
				repo.EXPECT().FetchLoanDetails(c, data.LoanId).Return(loanDetail, nil).Times(1)
			},
			expectedOutput: ApproveRejectLoanApplicationResponse{
				Status: false,
				Errors: []e.Error{{
					ErrName:     e.ErrorInfo[e.BadRequest].ErrName,
					Description: e.ErrorInfo[e.BadRequest].Description,
					Code:        e.ErrorInfo[e.BadRequest].Code,
				}},
				Message: "failed to update loan status",
			},
			httpStatus: http.StatusBadRequest,
			httpMethod: http.MethodPost,
		},
		{
			name: "RejectLoanError",
			input: ApproveRejectLoanApplicationRequest{
				LoanId:   3,
				Approval: LOAN_REJECT,
			},
			setup: func(c *gin.Context, data ApproveRejectLoanApplicationRequest) {
				ctrl := gomock.NewController(t)
				repo := dbmock.NewMockV1DBLayer(ctrl)
				dbObj = repo
				loanDetail := loan.LoanDetails{
					LoanId: sql.NullInt64{Int64: data.LoanId, Valid: true},
					Status: sql.NullString{String: LOAN_PENDING, Valid: true},
				}
				repo.EXPECT().FetchLoanDetails(c, data.LoanId).Return(loanDetail, nil).Times(1)
				repo.EXPECT().UpdateUnapprovedLoan(c, data.LoanId, false).Return(fmt.Errorf("db error")).Times(1)
			},
			expectedOutput: ApproveRejectLoanApplicationResponse{
				Status: false,
				Errors: []e.Error{{
					ErrName:     e.ErrorInfo[e.AddDBError].ErrName,
					Description: e.ErrorInfo[e.AddDBError].Description,
					Code:        e.ErrorInfo[e.AddDBError].Code,
				}},
				Message: "failed to update loan status",
			},
			httpStatus: http.StatusInternalServerError,
			httpMethod: http.MethodPost,
		},
		{
			name: "RejectLoanSuccess",
			input: ApproveRejectLoanApplicationRequest{
				LoanId:   3,
				Approval: LOAN_REJECT,
			},
			setup: func(c *gin.Context, data ApproveRejectLoanApplicationRequest) {
				ctrl := gomock.NewController(t)
				repo := dbmock.NewMockV1DBLayer(ctrl)
				dbObj = repo
				loanDetail := loan.LoanDetails{
					LoanId: sql.NullInt64{Int64: data.LoanId, Valid: true},
					Status: sql.NullString{String: LOAN_PENDING, Valid: true},
				}
				repo.EXPECT().FetchLoanDetails(c, data.LoanId).Return(loanDetail, nil).Times(1)
				repo.EXPECT().UpdateUnapprovedLoan(c, data.LoanId, false).Return(nil).Times(1)
			},
			expectedOutput: ApproveRejectLoanApplicationResponse{
				Status:  true,
				Message: "successfully updated loan status",
			},
			httpStatus: http.StatusOK,
			httpMethod: http.MethodPost,
		},
		{
			name: "CreateInstallmentsError",
			input: ApproveRejectLoanApplicationRequest{
				LoanId:   3,
				Approval: LOAN_APPROVE,
			},
			setup: func(c *gin.Context, data ApproveRejectLoanApplicationRequest) {
				ctrl := gomock.NewController(t)
				repo := dbmock.NewMockV1DBLayer(ctrl)
				dbObj = repo
				loanDetail := loan.LoanDetails{
					LoanId: sql.NullInt64{Int64: data.LoanId, Valid: true},
					Status: sql.NullString{String: LOAN_PENDING, Valid: true},
					Amount: sql.NullFloat64{Float64: 30000, Valid: true},
					Tenure: sql.NullInt64{Int64: 10, Valid: true},
				}
				installment := loanDetail.Amount.Float64 / float64(loanDetail.Tenure.Int64)
				repo.EXPECT().FetchLoanDetails(c, data.LoanId).Return(loanDetail, nil).Times(1)
				repo.EXPECT().UpdateAndInsertInstallments(c, data.LoanId, installment, loanDetail.Tenure.Int64).Return(fmt.Errorf("db error")).Times(1)
			},
			expectedOutput: ApproveRejectLoanApplicationResponse{
				Status: false,
				Errors: []e.Error{{
					ErrName:     e.ErrorInfo[e.AddDBError].ErrName,
					Description: e.ErrorInfo[e.AddDBError].Description,
					Code:        e.ErrorInfo[e.AddDBError].Code,
				}},
				Message: "failed to prepare loan installments",
			},
			httpStatus: http.StatusInternalServerError,
			httpMethod: http.MethodPost,
		},
		{
			name: "ApproveLoanSuccess",
			input: ApproveRejectLoanApplicationRequest{
				LoanId:   3,
				Approval: LOAN_APPROVE,
			},
			setup: func(c *gin.Context, data ApproveRejectLoanApplicationRequest) {
				ctrl := gomock.NewController(t)
				repo := dbmock.NewMockV1DBLayer(ctrl)
				dbObj = repo
				loanDetail := loan.LoanDetails{
					LoanId: sql.NullInt64{Int64: data.LoanId, Valid: true},
					Status: sql.NullString{String: LOAN_PENDING, Valid: true},
					Amount: sql.NullFloat64{Float64: 30000, Valid: true},
					Tenure: sql.NullInt64{Int64: 10, Valid: true},
				}
				installment := loanDetail.Amount.Float64 / float64(loanDetail.Tenure.Int64)
				repo.EXPECT().FetchLoanDetails(c, data.LoanId).Return(loanDetail, nil).Times(1)
				repo.EXPECT().UpdateAndInsertInstallments(c, data.LoanId, installment, loanDetail.Tenure.Int64).Return(nil).Times(1)
			},
			expectedOutput: ApproveRejectLoanApplicationResponse{
				Status:  true,
				Message: "successfully updated loan status",
			},
			httpStatus: http.StatusOK,
			httpMethod: http.MethodPost,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			fmt.Println("Starting Create Loan TestCase: ", tt.name)
			w, ctx := getContext(tt.httpMethod, tt.input, nil, nil)
			ctx.Set(config.USERID, userId)

			//setup test
			tt.setup(ctx, tt.input)
			servObj := NewLoanService(dbObj)

			//calling the function
			servObj.ApproveRejectLoanApplication(ctx)

			//check for result status
			assert.Equal(t, tt.httpStatus, w.Code)

			//create a copy of the output structure
			err := json.Unmarshal(w.Body.Bytes(), &tt.actualOutput)
			if err != nil {
				t.Error("unable to unmarshal response")
			}

			//compare expected vs actual output
			assert.Equal(t, tt.expectedOutput.Status, tt.actualOutput.Status)
			if len(tt.expectedOutput.Errors) != 0 {
				assert.Equal(t, tt.expectedOutput.Errors[0].Code, tt.actualOutput.Errors[0].Code)
			}

			fmt.Println("Ending Create Loan TestCase: ", tt.name)
		})
	}
}
