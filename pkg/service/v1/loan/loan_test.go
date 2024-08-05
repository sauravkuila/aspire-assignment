package loan

import (
	"aspire-assignment/pkg/config"
	v1 "aspire-assignment/pkg/db/v1"
	"aspire-assignment/pkg/db/v1/loan"
	dbmock "aspire-assignment/pkg/db/v1/mock"
	e "aspire-assignment/pkg/errors"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
)

func Test_loanService_CreateLoan(t *testing.T) {
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
		input          CreateLoanRequest
		setup          func(*gin.Context, CreateLoanRequest)
		expectedOutput CreateLoanResponse
		actualOutput   CreateLoanResponse
	}{
		{
			name: "MissingInputAmount",
			input: CreateLoanRequest{
				Tenure: 3,
			},
			setup: func(c *gin.Context, data CreateLoanRequest) {
				ctrl := gomock.NewController(t)
				repo := dbmock.NewMockV1DBLayer(ctrl)
				dbObj = repo
			},
			expectedOutput: CreateLoanResponse{
				Status: false,
				Errors: []e.Error{{
					ErrName:     e.ErrorInfo[e.BadRequest].ErrName,
					Description: e.ErrorInfo[e.BadRequest].Description,
					Code:        e.ErrorInfo[e.BadRequest].Code,
				}},
				Message: "failed to create loan",
			},
			httpStatus: http.StatusBadRequest,
			httpMethod: http.MethodPost,
		},
		{
			name: "MissingInputTenure",
			input: CreateLoanRequest{
				Amount: 34000,
			},
			setup: func(c *gin.Context, data CreateLoanRequest) {
				ctrl := gomock.NewController(t)
				repo := dbmock.NewMockV1DBLayer(ctrl)
				dbObj = repo
			},
			expectedOutput: CreateLoanResponse{
				Status: false,
				Errors: []e.Error{{
					ErrName:     e.ErrorInfo[e.BadRequest].ErrName,
					Description: e.ErrorInfo[e.BadRequest].Description,
					Code:        e.ErrorInfo[e.BadRequest].Code,
				}},
				Message: "failed to create loan",
			},
			httpStatus: http.StatusBadRequest,
			httpMethod: http.MethodPost,
		},
		{
			name: "FailToCreateLoan",
			input: CreateLoanRequest{
				Amount: 34000,
				Tenure: 3,
			},
			setup: func(c *gin.Context, data CreateLoanRequest) {
				ctrl := gomock.NewController(t)
				repo := dbmock.NewMockV1DBLayer(ctrl)
				dbObj = repo
				repo.EXPECT().CreateLoan(c, userId, data.Amount, data.Tenure).Return(int64(0), fmt.Errorf("failed to create loan")).Times(1)
			},
			expectedOutput: CreateLoanResponse{
				Status: false,
				Errors: []e.Error{{
					ErrName:     e.ErrorInfo[e.AddDBError].ErrName,
					Description: e.ErrorInfo[e.AddDBError].Description,
					Code:        e.ErrorInfo[e.AddDBError].Code,
				}},
				Message: "failed to create loan",
			},
			httpStatus: http.StatusInternalServerError,
			httpMethod: http.MethodPost,
		},
		{
			name: "SuccessCreateLoan",
			input: CreateLoanRequest{
				Amount: 34000,
				Tenure: 3,
			},
			setup: func(c *gin.Context, data CreateLoanRequest) {
				ctrl := gomock.NewController(t)
				repo := dbmock.NewMockV1DBLayer(ctrl)
				dbObj = repo
				repo.EXPECT().CreateLoan(c, userId, data.Amount, data.Tenure).Return(int64(1), nil).Times(1)
			},
			expectedOutput: CreateLoanResponse{
				Status: true,
				Data: &LoanDetails{
					LoanId: 1,
					Amount: 34000,
					Tenure: 3,
					Status: LOAN_PENDING,
				},
				Message: "successfully created loan",
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
			servObj.CreateLoan(ctx)

			//check for result status
			assert.Equal(t, tt.httpStatus, w.Code)

			//create a copy of the output structure
			err := json.Unmarshal(w.Body.Bytes(), &tt.actualOutput)
			if err != nil {
				t.Error("unable to unmarshal response")
			}

			//compare expected vs actual output
			assert.Equal(t, tt.expectedOutput, tt.actualOutput)

			fmt.Println("Ending Create Loan TestCase: ", tt.name)
		})
	}
}

func Test_loanService_ModifyLoan(t *testing.T) {
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
		input          ModifyLoanRequest
		setup          func(*gin.Context, ModifyLoanRequest)
		expectedOutput ModifyLoanResponse
		actualOutput   ModifyLoanResponse
	}{
		{
			name: "MissingInputAmount",
			input: ModifyLoanRequest{
				LoanId: 3,
				Tenure: 3,
			},
			setup: func(c *gin.Context, data ModifyLoanRequest) {
				ctrl := gomock.NewController(t)
				repo := dbmock.NewMockV1DBLayer(ctrl)
				dbObj = repo
			},
			expectedOutput: ModifyLoanResponse{
				Status: false,
				Errors: []e.Error{{
					ErrName:     e.ErrorInfo[e.BadRequest].ErrName,
					Description: e.ErrorInfo[e.BadRequest].Description,
					Code:        e.ErrorInfo[e.BadRequest].Code,
				}},
				Message: "failed to modify loan",
			},
			httpStatus: http.StatusBadRequest,
			httpMethod: http.MethodPost,
		},
		{
			name: "MissingInputTenure",
			input: ModifyLoanRequest{
				LoanId: 3,
				Amount: 34000,
			},
			setup: func(c *gin.Context, data ModifyLoanRequest) {
				ctrl := gomock.NewController(t)
				repo := dbmock.NewMockV1DBLayer(ctrl)
				dbObj = repo
			},
			expectedOutput: ModifyLoanResponse{
				Status: false,
				Errors: []e.Error{{
					ErrName:     e.ErrorInfo[e.BadRequest].ErrName,
					Description: e.ErrorInfo[e.BadRequest].Description,
					Code:        e.ErrorInfo[e.BadRequest].Code,
				}},
				Message: "failed to modify loan",
			},
			httpStatus: http.StatusBadRequest,
			httpMethod: http.MethodPost,
		},
		{
			name: "MissingInputLoanId",
			input: ModifyLoanRequest{
				Amount: 34000,
				Tenure: 3,
			},
			setup: func(c *gin.Context, data ModifyLoanRequest) {
				ctrl := gomock.NewController(t)
				repo := dbmock.NewMockV1DBLayer(ctrl)
				dbObj = repo
			},
			expectedOutput: ModifyLoanResponse{
				Status: false,
				Errors: []e.Error{{
					ErrName:     e.ErrorInfo[e.BadRequest].ErrName,
					Description: e.ErrorInfo[e.BadRequest].Description,
					Code:        e.ErrorInfo[e.BadRequest].Code,
				}},
				Message: "failed to modify loan",
			},
			httpStatus: http.StatusBadRequest,
			httpMethod: http.MethodPost,
		},
		{
			name: "FailToModifyLoan",
			input: ModifyLoanRequest{
				LoanId: 3,
				Amount: 34000,
				Tenure: 3,
			},
			setup: func(c *gin.Context, data ModifyLoanRequest) {
				ctrl := gomock.NewController(t)
				repo := dbmock.NewMockV1DBLayer(ctrl)
				dbObj = repo
				repo.EXPECT().ModifyLoan(c, userId, data.LoanId, data.Amount, data.Tenure).Return(int64(0), fmt.Errorf("failed to modify loan")).Times(1)
			},
			expectedOutput: ModifyLoanResponse{
				Status: false,
				Errors: []e.Error{{
					ErrName:     e.ErrorInfo[e.AddDBError].ErrName,
					Description: e.ErrorInfo[e.AddDBError].Description,
					Code:        e.ErrorInfo[e.AddDBError].Code,
				}},
				Message: "failed to modify loan",
			},
			httpStatus: http.StatusInternalServerError,
			httpMethod: http.MethodPost,
		},
		{
			name: "FailToModifyLoanForNonPendingLoan",
			input: ModifyLoanRequest{
				LoanId: 3,
				Amount: 34000,
				Tenure: 3,
			},
			setup: func(c *gin.Context, data ModifyLoanRequest) {
				ctrl := gomock.NewController(t)
				repo := dbmock.NewMockV1DBLayer(ctrl)
				dbObj = repo
				repo.EXPECT().ModifyLoan(c, userId, data.LoanId, data.Amount, data.Tenure).Return(int64(0), nil).Times(1)
			},
			expectedOutput: ModifyLoanResponse{
				Status: false,
				Errors: []e.Error{{
					ErrName:     e.ErrorInfo[e.BadRequest].ErrName,
					Description: e.ErrorInfo[e.BadRequest].Description + " | only loans created by user in PENDING status can be modified",
					Code:        e.ErrorInfo[e.BadRequest].Code,
				}},
				Message: "failed to modify loan",
			},
			httpStatus: http.StatusBadRequest,
			httpMethod: http.MethodPost,
		},
		{
			name: "SuccessModifyLoan",
			input: ModifyLoanRequest{
				LoanId: 3,
				Amount: 34000,
				Tenure: 3,
			},
			setup: func(c *gin.Context, data ModifyLoanRequest) {
				ctrl := gomock.NewController(t)
				repo := dbmock.NewMockV1DBLayer(ctrl)
				dbObj = repo
				repo.EXPECT().ModifyLoan(c, userId, data.LoanId, data.Amount, data.Tenure).Return(int64(1), nil).Times(1)
			},
			expectedOutput: ModifyLoanResponse{
				Status: true,
				Data: &LoanDetails{
					LoanId: 3,
					Amount: 34000,
					Tenure: 3,
					Status: LOAN_PENDING,
				},
				Message: "successfully modified loan",
			},
			httpStatus: http.StatusOK,
			httpMethod: http.MethodPost,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			fmt.Println("Starting Modify Loan TestCase: ", tt.name)
			w, ctx := getContext(tt.httpMethod, tt.input, nil, nil)
			ctx.Set(config.USERID, userId)

			//setup test
			tt.setup(ctx, tt.input)
			servObj := NewLoanService(dbObj)

			//calling the function
			servObj.ModifyLoan(ctx)

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

			fmt.Println("Ending Modify Loan TestCase: ", tt.name)
		})
	}
}

func Test_loanService_CancelLoan(t *testing.T) {
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
		input          CancelLoanRequest
		setup          func(*gin.Context, CancelLoanRequest)
		expectedOutput CancelLoanResponse
		actualOutput   CancelLoanResponse
	}{
		{
			name:  "MissingInputLoanId",
			input: CancelLoanRequest{},
			setup: func(c *gin.Context, data CancelLoanRequest) {
				ctrl := gomock.NewController(t)
				repo := dbmock.NewMockV1DBLayer(ctrl)
				dbObj = repo
			},
			expectedOutput: CancelLoanResponse{
				Status: false,
				Errors: []e.Error{{
					ErrName:     e.ErrorInfo[e.BadRequest].ErrName,
					Description: e.ErrorInfo[e.BadRequest].Description,
					Code:        e.ErrorInfo[e.BadRequest].Code,
				}},
				Message: "failed to cancel loan",
			},
			httpStatus: http.StatusBadRequest,
			httpMethod: http.MethodPost,
		},
		{
			name: "FailToCancelLoan",
			input: CancelLoanRequest{
				LoanId: 3,
			},
			setup: func(c *gin.Context, data CancelLoanRequest) {
				ctrl := gomock.NewController(t)
				repo := dbmock.NewMockV1DBLayer(ctrl)
				dbObj = repo
				repo.EXPECT().CancelLoan(c, userId, data.LoanId).Return(int64(0), fmt.Errorf("failed to cancel loan")).Times(1)
			},
			expectedOutput: CancelLoanResponse{
				Status: false,
				Errors: []e.Error{{
					ErrName:     e.ErrorInfo[e.AddDBError].ErrName,
					Description: e.ErrorInfo[e.AddDBError].Description,
					Code:        e.ErrorInfo[e.AddDBError].Code,
				}},
				Message: "failed to cancel loan",
			},
			httpStatus: http.StatusInternalServerError,
			httpMethod: http.MethodPost,
		},
		{
			name: "FailToCancelLoanForNonPendingLoan",
			input: CancelLoanRequest{
				LoanId: 3,
			},
			setup: func(c *gin.Context, data CancelLoanRequest) {
				ctrl := gomock.NewController(t)
				repo := dbmock.NewMockV1DBLayer(ctrl)
				dbObj = repo
				repo.EXPECT().CancelLoan(c, userId, data.LoanId).Return(int64(0), nil).Times(1)
			},
			expectedOutput: CancelLoanResponse{
				Status: false,
				Errors: []e.Error{{
					ErrName:     e.ErrorInfo[e.BadRequest].ErrName,
					Description: e.ErrorInfo[e.BadRequest].Description + " | only loans created by user in PENDING status can be cancelled",
					Code:        e.ErrorInfo[e.BadRequest].Code,
				}},
				Message: "failed to modify loan",
			},
			httpStatus: http.StatusBadRequest,
			httpMethod: http.MethodPost,
		},
		{
			name: "SuccessCancelLoan",
			input: CancelLoanRequest{
				LoanId: 3,
			},
			setup: func(c *gin.Context, data CancelLoanRequest) {
				ctrl := gomock.NewController(t)
				repo := dbmock.NewMockV1DBLayer(ctrl)
				dbObj = repo
				repo.EXPECT().CancelLoan(c, userId, data.LoanId).Return(int64(3), nil).Times(1)
			},
			expectedOutput: CancelLoanResponse{
				Status: true,
				Data: &LoanDetails{
					LoanId: 3,
					Status: LOAN_CANCELLED,
				},
				Message: "successfully cancelled loan",
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
			servObj.CancelLoan(ctx)

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

func Test_loanService_GetLoans(t *testing.T) {
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
		input          GetLoanRequest
		setup          func(*gin.Context, GetLoanRequest)
		expectedOutput GetLoanResponse
		actualOutput   GetLoanResponse
	}{
		{
			name:  "FailToGetLoan",
			input: GetLoanRequest{},
			setup: func(c *gin.Context, data GetLoanRequest) {
				ctrl := gomock.NewController(t)
				repo := dbmock.NewMockV1DBLayer(ctrl)
				dbObj = repo
				repo.EXPECT().GetUserLoans(c, userId).Return(nil, fmt.Errorf("db error")).Times(1)
			},
			expectedOutput: GetLoanResponse{
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
			name:  "InvalidLoanId",
			input: GetLoanRequest{},
			setup: func(c *gin.Context, data GetLoanRequest) {
				ctrl := gomock.NewController(t)
				repo := dbmock.NewMockV1DBLayer(ctrl)
				dbObj = repo
				// loans := make([]loan.LoanDetails, 0)
				repo.EXPECT().GetUserLoans(c, userId).Return(nil, nil).Times(1)
			},
			expectedOutput: GetLoanResponse{
				Status:  false,
				Message: "no loans available for user",
			},
			httpStatus: http.StatusNotFound,
			httpMethod: http.MethodGet,
		},
		{
			name:  "SuccessGetLoans",
			input: GetLoanRequest{},
			setup: func(c *gin.Context, data GetLoanRequest) {
				ctrl := gomock.NewController(t)
				repo := dbmock.NewMockV1DBLayer(ctrl)
				dbObj = repo
				loans := make([]loan.LoanDetails, 0)
				t1, _ := time.Parse("2006-01-02 15:04:05", "2024-08-08 15:00:00")
				loans = append(loans, loan.LoanDetails{
					LoanId:    sql.NullInt64{Int64: 3, Valid: true},
					Amount:    sql.NullFloat64{Float64: 34000, Valid: true},
					Tenure:    sql.NullInt64{Int64: 3, Valid: true},
					Status:    sql.NullString{String: LOAN_PENDING, Valid: true},
					CreatedAt: sql.NullTime{Time: t1, Valid: true},
				})
				repo.EXPECT().GetUserLoans(c, userId).Return(loans, nil).Times(1)
			},
			expectedOutput: GetLoanResponse{
				Status: true,
				Data: []LoanDetails{{
					LoanId:    3,
					Amount:    34000,
					Tenure:    3,
					Status:    LOAN_PENDING,
					CreatedAt: "2024-08-08 15:00:00",
				}},
				Message: "successfully fetched user loans",
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
			servObj.GetLoans(ctx)

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

func getContext(method string, data interface{}, queries map[string]string, params map[string]string) (w *httptest.ResponseRecorder, c *gin.Context) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	temp, _ := gin.CreateTestContext(recorder)
	byteData, err := json.Marshal(data)
	if err != nil {
		log.Fatalln(err)
	}
	temp.Request, err = http.NewRequest(method, "/", bytes.NewBuffer(byteData))
	if err != nil {
		log.Fatalln(err)
	}

	//add headers
	temp.Request.Header = http.Header{}
	temp.Request.Header.Set("Content-Type", "application/json")

	//add query params
	if queries != nil {
		q := temp.Request.URL.Query()
		for k, v := range queries {
			q.Add(k, v)
		}
		temp.Request.URL.RawQuery = q.Encode()
	}

	//add path params
	for k, v := range params {
		temp.Params = append(temp.Params, gin.Param{
			Key:   k,
			Value: v,
		})
	}

	return recorder, temp
}
