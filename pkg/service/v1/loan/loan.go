package loan

import (
	"log"
	"net/http"

	e "aspire-assignment/pkg/errors"

	"github.com/gin-gonic/gin"
)

func (obj *loanService) FuncLoanServiceSample(c *gin.Context) {
	c.JSON(http.StatusOK, &gin.H{"status": "working"})
}

func (obj *loanService) CreateLoan(c *gin.Context) {
	var (
		request  CreateLoanRequest
		response CreateLoanResponse
	)
	if err := c.BindJSON(&request); err != nil {
		log.Printf("unable to marshal request. Error:%s", err.Error())
		response.Errors = append(response.Errors, *e.ErrorInfo[e.BadRequest])
		response.Message = "failed to create loan"
		c.JSON(http.StatusBadRequest, response)
		return
	}

	//make a loan entry in db
	loanId, err := obj.dbObj.CreateLoan(c, request.UserId, request.Amount, request.Installments)
	if err != nil {
		log.Printf("failed to create a loan. Error:%s", err.Error())
		response.Errors = append(response.Errors, *e.ErrorInfo[e.AddDBError])
		response.Message = "failed to create loan"
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	loanDetail := LoanDetails{
		LoanId:       loanId,
		Amount:       request.Amount,
		Installments: request.Installments,
		Status:       LOAN_PENDING,
	}
	response.Status = true
	response.Data = &loanDetail
	response.Message = "successfully created loan"

	c.JSON(http.StatusOK, response)
}

func (obj *loanService) ModifyLoan(c *gin.Context) {
	var (
		request  ModifyLoanRequest
		response ModifyLoanResponse
	)
	if err := c.BindJSON(&request); err != nil {
		log.Printf("unable to marshal request. Error:%s", err.Error())
		response.Errors = append(response.Errors, *e.ErrorInfo[e.BadRequest])
		response.Message = "failed to modify loan"
		c.JSON(http.StatusBadRequest, response)
		return
	}

	//modify the loan if the loan is pending
	loanId, err := obj.dbObj.ModifyLoan(c, request.UserId, request.LoanId, request.Amount, request.Installments)
	if err != nil {
		log.Printf("failed to modify a loan. Error:%s", err.Error())
		response.Errors = append(response.Errors, *e.ErrorInfo[e.AddDBError])
		response.Message = "failed to modify loan"
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	//incorrect loan-user relation or the status is not PENDING
	if loanId == 0 {
		log.Printf("failed to modify a loan. This can be because loan is not pending or user-loan relation is incorrect")
		response.Errors = append(response.Errors, e.ErrorInfo[e.BadRequest].GetErrorDetails("only loans created by user in PENDING status can be modified"))
		response.Message = "failed to modify loan"
		c.JSON(http.StatusBadRequest, response)
		return
	}

	loanDetail := LoanDetails{
		LoanId:       loanId,
		Amount:       request.Amount,
		Installments: request.Installments,
		Status:       LOAN_PENDING,
	}
	response.Status = true
	response.Data = &loanDetail
	response.Message = "successfully modified loan"

	c.JSON(http.StatusOK, response)
}
func (obj *loanService) CancelLoan(c *gin.Context) {
	var (
		request  CancelLoanRequest
		response CancelLoanResponse
	)
	if err := c.BindJSON(&request); err != nil {
		log.Printf("unable to marshal request. Error:%s", err.Error())
		response.Errors = append(response.Errors, *e.ErrorInfo[e.BadRequest])
		response.Message = "failed to cancel loan"
		c.JSON(http.StatusBadRequest, response)
		return
	}

	//cancel the loan if the loan is pending
	loanId, err := obj.dbObj.CancelLoan(c, request.UserId, request.LoanId)
	if err != nil {
		log.Printf("failed to modify a loan. Error:%s", err.Error())
		response.Errors = append(response.Errors, *e.ErrorInfo[e.AddDBError])
		response.Message = "failed to modify loan"
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	//incorrect loan-user relation or the status is not PENDING
	if loanId == 0 {
		log.Printf("failed to cancel a loan. This can be because loan is not pending or user-loan relation is incorrect")
		response.Errors = append(response.Errors, e.ErrorInfo[e.BadRequest].GetErrorDetails("only loans created by user in PENDING status can be cancelled"))
		response.Message = "failed to cancel loan"
		c.JSON(http.StatusBadRequest, response)
		return
	}

	loanDetail := LoanDetails{
		LoanId: loanId,
		Status: LOAN_CANCELLED,
	}
	response.Status = true
	response.Data = &loanDetail
	response.Message = "successfully modified loan"

	c.JSON(http.StatusOK, response)
}

func (obj *loanService) GetLoans(c *gin.Context) {
	var (
		request  GetLoanRequest
		response GetLoanResponse
	)
	if err := c.BindQuery(&request); err != nil {
		log.Printf("unable to marshal request. Error:%s", err.Error())
		response.Errors = append(response.Errors, *e.ErrorInfo[e.BadRequest])
		response.Message = "failed to fetch loans"
		c.JSON(http.StatusBadRequest, response)
		return
	}

	//TODO: add custom status like only loans which are pending or cancelled. add a query scan param

	//fetch loans
	loans, err := obj.dbObj.GetUserLoans(c, request.UserId)
	if err != nil {
		log.Printf("failed to fetch loans. Error:%s", err.Error())
		response.Errors = append(response.Errors, *e.ErrorInfo[e.AddDBError])
		response.Message = "failed to fetch loans"
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response.Status = true
	if len(loans) == 0 {
		response.Message = "no loans available for user"
		c.JSON(http.StatusNotFound, response)
		return
	}

	//init loan slice
	response.Data = make([]LoanDetails, 0)
	for _, loan := range loans {
		response.Data = append(response.Data, LoanDetails{
			LoanId:       loan.LoanId.Int64,
			Amount:       loan.Amount.Float64,
			Installments: loan.Installments.Int64,
			Status:       loan.Status.String,
			CreatedAt:    loan.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		})
	}
	response.Message = "successfully fetched user loans"
	c.JSON(http.StatusOK, response)
}
