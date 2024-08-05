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

func Test_loanService_GetInstallments(t *testing.T) {
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
		input          GetLoanDetailRequest
		queryParams    map[string]string
		setup          func(*gin.Context, GetLoanDetailRequest)
		expectedOutput GetLoanDetailResponse
		actualOutput   GetLoanDetailResponse
	}{
		{
			name:        "MissingInputLoanId",
			input:       GetLoanDetailRequest{},
			queryParams: make(map[string]string),
			setup: func(c *gin.Context, data GetLoanDetailRequest) {
				ctrl := gomock.NewController(t)
				repo := dbmock.NewMockV1DBLayer(ctrl)
				dbObj = repo
			},
			expectedOutput: GetLoanDetailResponse{
				Status: false,
				Errors: []e.Error{{
					ErrName:     e.ErrorInfo[e.BadRequest].ErrName,
					Description: e.ErrorInfo[e.BadRequest].Description,
					Code:        e.ErrorInfo[e.BadRequest].Code,
				}},
				Message: "failed to fetch installments",
			},
			httpStatus: http.StatusBadRequest,
			httpMethod: http.MethodGet,
		},
		{
			name: "ErrorFetchingInstallments",
			input: GetLoanDetailRequest{
				LoanId: 3,
			},
			queryParams: make(map[string]string),
			setup: func(c *gin.Context, data GetLoanDetailRequest) {
				ctrl := gomock.NewController(t)
				repo := dbmock.NewMockV1DBLayer(ctrl)
				dbObj = repo
				repo.EXPECT().GetUserLoanInstallments(c, userId, data.LoanId).Return(nil, fmt.Errorf("db error")).Times(1)
			},
			expectedOutput: GetLoanDetailResponse{
				Status: false,
				Errors: []e.Error{{
					ErrName:     e.ErrorInfo[e.GetDBError].ErrName,
					Description: e.ErrorInfo[e.GetDBError].Description,
					Code:        e.ErrorInfo[e.GetDBError].Code,
				}},
				Message: "failed to fetch installments",
			},
			httpStatus: http.StatusInternalServerError,
			httpMethod: http.MethodGet,
		},
		{
			name: "NoInstallmentsForLoan",
			input: GetLoanDetailRequest{
				LoanId: 3,
			},
			queryParams: make(map[string]string),
			setup: func(c *gin.Context, data GetLoanDetailRequest) {
				ctrl := gomock.NewController(t)
				repo := dbmock.NewMockV1DBLayer(ctrl)
				dbObj = repo
				installments := make([]loan.InstallmentDetails, 0)
				repo.EXPECT().GetUserLoanInstallments(c, userId, data.LoanId).Return(installments, nil).Times(1)
			},
			expectedOutput: GetLoanDetailResponse{
				Status:  false,
				Message: "no installments against loan available",
			},
			httpStatus: http.StatusNotFound,
			httpMethod: http.MethodGet,
		},
		{
			name: "getInstallmentSuccess",
			input: GetLoanDetailRequest{
				LoanId: 3,
			},
			queryParams: make(map[string]string),
			setup: func(c *gin.Context, data GetLoanDetailRequest) {
				ctrl := gomock.NewController(t)
				repo := dbmock.NewMockV1DBLayer(ctrl)
				dbObj = repo
				t1, _ := time.Parse("2006-01-02", "2024-08-08")
				installments := make([]loan.InstallmentDetails, 2)
				installments[0] = loan.InstallmentDetails{
					AmountDue:      sql.NullFloat64{Float64: 5000, Valid: true},
					AmountPaid:     sql.NullFloat64{Float64: 5000, Valid: true},
					Status:         sql.NullString{String: TXN_PAID, Valid: true},
					InstallmentSeq: sql.NullInt64{Int64: 1, Valid: true},
					TransactionId:  sql.NullString{String: "txn1", Valid: true},
					DueDate:        sql.NullTime{Time: t1, Valid: true},
					LoanAmount:     sql.NullFloat64{Float64: 10000, Valid: true},
					LoanStatus:     sql.NullString{String: LOAN_APPROVED, Valid: true},
				}
				installments[1] = loan.InstallmentDetails{
					AmountDue:      sql.NullFloat64{Float64: 5000, Valid: true},
					AmountPaid:     sql.NullFloat64{Float64: 0, Valid: true},
					Status:         sql.NullString{String: TXN_PENDING, Valid: true},
					InstallmentSeq: sql.NullInt64{Int64: 2, Valid: true},
					TransactionId:  sql.NullString{String: "txn2", Valid: true},
					DueDate:        sql.NullTime{Time: t1.Add(24 * time.Hour), Valid: true},
					LoanAmount:     sql.NullFloat64{Float64: 10000, Valid: true},
					LoanStatus:     sql.NullString{String: LOAN_APPROVED, Valid: true},
				}
				repo.EXPECT().GetUserLoanInstallments(c, userId, data.LoanId).Return(installments, nil).Times(1)
			},
			expectedOutput: GetLoanDetailResponse{
				Status: true,
				Data: &GetLoanDetail{
					LoanId:            3,
					LoanAmount:        10000,
					OutstandingAmount: 5000,
					Tenure:            2,
					Status:            LOAN_APPROVED,
					Installments: []InstallmentDetails{{
						AmoundDue:         5000,
						AmountPaid:        5000,
						Status:            TXN_PAID,
						InstallmentNumber: 1,
						TransactionId:     "txn1",
						DueDate:           "2024-08-08",
					}, {
						AmoundDue:         5000,
						AmountPaid:        0,
						Status:            TXN_PENDING,
						InstallmentNumber: 2,
						TransactionId:     "txn2",
						DueDate:           "2024-08-15",
					}},
				},
				Message: "successfully fetched installments",
			},
			httpStatus: http.StatusOK,
			httpMethod: http.MethodGet,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			fmt.Println("Starting Create Loan TestCase: ", tt.name)
			tt.queryParams["loanId"] = fmt.Sprintf("%d", tt.input.LoanId)
			w, ctx := getContext(tt.httpMethod, nil, tt.queryParams, nil)
			ctx.Set(config.USERID, userId)

			//setup test
			tt.setup(ctx, tt.input)
			servObj := NewLoanService(dbObj)

			//calling the function
			servObj.GetInstallments(ctx)

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

func Test_loanService_ProcessLoanPayment(t *testing.T) {
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
		input          ProcessLoanPaymentRequest
		setup          func(*gin.Context, ProcessLoanPaymentRequest)
		expectedOutput ProcessLoanPaymentResponse
		actualOutput   ProcessLoanPaymentResponse
	}{
		{
			name: "MissingInputLoanId",
			input: ProcessLoanPaymentRequest{
				Amount:        5000,
				TransactionId: "txn1",
			},
			setup: func(c *gin.Context, data ProcessLoanPaymentRequest) {
				ctrl := gomock.NewController(t)
				repo := dbmock.NewMockV1DBLayer(ctrl)
				dbObj = repo
			},
			expectedOutput: ProcessLoanPaymentResponse{
				Status: false,
				Errors: []e.Error{{
					ErrName:     e.ErrorInfo[e.BadRequest].ErrName,
					Description: e.ErrorInfo[e.BadRequest].Description,
					Code:        e.ErrorInfo[e.BadRequest].Code,
				}},
				Message: "failed to process payment",
			},
			httpStatus: http.StatusBadRequest,
			httpMethod: http.MethodPost,
		},
		{
			name: "MissingInputAmount",
			input: ProcessLoanPaymentRequest{
				LoanId:        3,
				TransactionId: "txn1",
			},
			setup: func(c *gin.Context, data ProcessLoanPaymentRequest) {
				ctrl := gomock.NewController(t)
				repo := dbmock.NewMockV1DBLayer(ctrl)
				dbObj = repo
			},
			expectedOutput: ProcessLoanPaymentResponse{
				Status: false,
				Errors: []e.Error{{
					ErrName:     e.ErrorInfo[e.BadRequest].ErrName,
					Description: e.ErrorInfo[e.BadRequest].Description,
					Code:        e.ErrorInfo[e.BadRequest].Code,
				}},
				Message: "failed to process payment",
			},
			httpStatus: http.StatusBadRequest,
			httpMethod: http.MethodPost,
		},
		{
			name: "MissingInputTransactionId",
			input: ProcessLoanPaymentRequest{
				LoanId: 3,
				Amount: 5000,
			},
			setup: func(c *gin.Context, data ProcessLoanPaymentRequest) {
				ctrl := gomock.NewController(t)
				repo := dbmock.NewMockV1DBLayer(ctrl)
				dbObj = repo
			},
			expectedOutput: ProcessLoanPaymentResponse{
				Status: false,
				Errors: []e.Error{{
					ErrName:     e.ErrorInfo[e.BadRequest].ErrName,
					Description: e.ErrorInfo[e.BadRequest].Description,
					Code:        e.ErrorInfo[e.BadRequest].Code,
				}},
				Message: "failed to process payment",
			},
			httpStatus: http.StatusBadRequest,
			httpMethod: http.MethodPost,
		},
		{
			name: "ErrorFetchingInstallments",
			input: ProcessLoanPaymentRequest{
				LoanId:        3,
				Amount:        5000,
				TransactionId: "txn1",
			},
			setup: func(c *gin.Context, data ProcessLoanPaymentRequest) {
				ctrl := gomock.NewController(t)
				repo := dbmock.NewMockV1DBLayer(ctrl)
				dbObj = repo
				repo.EXPECT().GetUserLoanInstallments(c, userId, data.LoanId).Return(nil, fmt.Errorf("db error")).Times(1)
			},
			expectedOutput: ProcessLoanPaymentResponse{
				Status: false,
				Errors: []e.Error{{
					ErrName:     e.ErrorInfo[e.GetDBError].ErrName,
					Description: e.ErrorInfo[e.GetDBError].Description,
					Code:        e.ErrorInfo[e.GetDBError].Code,
				}},
				Message: "failed to fetch installments to process payment",
			},
			httpStatus: http.StatusInternalServerError,
			httpMethod: http.MethodPost,
		},
		{
			name: "NoInstallments",
			input: ProcessLoanPaymentRequest{
				LoanId:        3,
				Amount:        5000,
				TransactionId: "txn1",
			},
			setup: func(c *gin.Context, data ProcessLoanPaymentRequest) {
				ctrl := gomock.NewController(t)
				repo := dbmock.NewMockV1DBLayer(ctrl)
				dbObj = repo
				repo.EXPECT().GetUserLoanInstallments(c, userId, data.LoanId).Return(nil, nil).Times(1)
			},
			expectedOutput: ProcessLoanPaymentResponse{
				Status: false,
				// Errors: []e.Error{{
				// 	ErrName:     e.ErrorInfo[e.GetDBError].ErrName,
				// 	Description: e.ErrorInfo[e.GetDBError].Description,
				// 	Code:        e.ErrorInfo[e.GetDBError].Code,
				// }},
				Message: "no installments against loan available",
			},
			httpStatus: http.StatusBadRequest,
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
			servObj.ProcessLoanPayment(ctx)

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
