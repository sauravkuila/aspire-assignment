package loan

import (
	"aspire-assignment/pkg/config"
	e "aspire-assignment/pkg/errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (obj *loanService) GetPendingLoans(c *gin.Context) {
	var (
		request  PendingLoanRequest
		response PendingLoanResponse
	)
	// if err := c.BindQuery(&request); err != nil {
	// 	log.Printf("unable to marshal request. Error:%s", err.Error())
	// 	response.Errors = append(response.Errors, *e.ErrorInfo[e.BadRequest])
	// 	response.Message = "failed to fetch loans"
	// 	c.JSON(http.StatusBadRequest, response)
	// 	return
	// }
	request.UserId = c.GetInt64(config.USERID)

	//check if the user is an admin and fetch only pending loans
	loans, err := obj.dbObj.GetUnapprovedLoans(c)
	if err != nil {
		log.Printf("failed to fetch loans. Error:%s", err.Error())
		response.Errors = append(response.Errors, *e.ErrorInfo[e.GetDBError])
		response.Message = "failed to fetch loans"
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	if len(loans) == 0 {
		response.Message = "no loans available"
		c.JSON(http.StatusNotFound, response)
		return
	}

	response.Status = true
	//init loan slice
	response.Data = make([]LoanDetails, 0)
	for _, loan := range loans {
		response.Data = append(response.Data, LoanDetails{
			LoanId:    loan.LoanId.Int64,
			UserName:  loan.UserName.String,
			Amount:    loan.Amount.Float64,
			Tenure:    loan.Installments.Int64,
			Status:    loan.Status.String,
			CreatedAt: loan.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		})
	}
	response.Message = "successfully fetched unapproved loans"
	c.JSON(http.StatusOK, response)

}

func (obj *loanService) ApproveRejectLoanApplication(c *gin.Context) {
	var (
		request  ApproveRejectLoanApplicationRequest
		response ApproveRejectLoanApplicationResponse
	)
	if err := c.BindJSON(&request); err != nil {
		log.Printf("unable to marshal request. Error:%s", err.Error())
		response.Errors = append(response.Errors, *e.ErrorInfo[e.BadRequest])
		response.Message = "failed to update loan approval"
		c.JSON(http.StatusBadRequest, response)
		return
	}
	request.UserId = c.GetInt64(config.USERID)

	//check loan details to create transactions
	loanDetail, err := obj.dbObj.FetchLoanDetails(c, request.LoanId)
	if err != nil {
		log.Printf("failed to fetch loan detail. Error:%s", err.Error())
		response.Errors = append(response.Errors, *e.ErrorInfo[e.NoDataFound])
		response.Message = "failed to fetch loan detail"
		c.JSON(http.StatusNotFound, response)
		return
	}

	if loanDetail.Status.String != LOAN_PENDING {
		response.Errors = append(response.Errors, e.ErrorInfo[e.BadRequest].GetErrorDetails("loan not in PENDING state"))
		response.Message = "failed to update loan status"
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if request.Approval == LOAN_REJECT {
		log.Printf("loan is being rejected by admin. LoanId: %d, Status: %s", request.LoanId, request.Approval)
		//update the rejection in db
		err := obj.dbObj.UpdateUnapprovedLoan(c, request.LoanId, false)
		if err != nil {
			log.Printf("failed to update loan status. Error:%s", err.Error())
			response.Errors = append(response.Errors, *e.ErrorInfo[e.AddDBError])
			response.Message = "failed to update loan status"
			c.JSON(http.StatusInternalServerError, response)
			return
		}
		response.Status = true
		response.Message = "successfully updated loan status"
		c.JSON(http.StatusOK, response)
		return
	}

	//finding installment per week but any other logic for installment can be applied here
	equalInstallmentAmount := loanDetail.Amount.Float64 / float64(loanDetail.Tenure.Int64)
	//update and insert transactions
	err = obj.dbObj.UpdateAndInsertInstallments(c, request.LoanId, equalInstallmentAmount, loanDetail.Tenure.Int64)
	if err != nil {
		log.Printf("failed to prepare loan installments. Error:%s", err.Error())
		response.Errors = append(response.Errors, *e.ErrorInfo[e.AddDBError])
		response.Message = "failed to prepare loan installments"
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response.Status = true
	response.Message = "successfully updated loan status"
	c.JSON(http.StatusOK, response)
}
