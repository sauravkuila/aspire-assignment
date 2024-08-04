package loan

import (
	e "aspire-assignment/pkg/errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (obj *loanService) GetInstallments(c *gin.Context) {
	var (
		request  GetLoanDetailRequest
		response GetLoanDetailResponse
	)
	if err := c.BindQuery(&request); err != nil {
		log.Printf("unable to marshal request. Error:%s", err.Error())
		response.Errors = append(response.Errors, *e.ErrorInfo[e.BadRequest])
		response.Message = "failed to fetch installments"
		c.JSON(http.StatusBadRequest, response)
		return
	}

	installments, err := obj.dbObj.GetUserLoanInstallments(c, request.UserId, request.LoanId)
	if err != nil {
		log.Printf("failed to fetch loan installments. Error:%s", err.Error())
		response.Errors = append(response.Errors, e.ErrorInfo[e.GetDBError].GetErrorDetails(err.Error()))
		response.Message = "failed to fetch installments"
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response.Status = true
	if len(installments) == 0 {
		response.Message = "no installments against loan available"
		c.JSON(http.StatusNotFound, response)
		return
	}
	response.Data = &GetLoanDetail{
		LoanId:            request.LoanId,
		TotalInstallments: len(installments),
		LoanAmount:        installments[0].LoanAmount.Float64,
		OutstandingAmount: installments[0].LoanAmount.Float64,
		Status:            installments[0].LoanStatus.String,
		Installments:      make([]InstallmentDetails, 0),
	}
	for _, installment := range installments {
		response.Data.Installments = append(response.Data.Installments, InstallmentDetails{
			AmoundDue:         installment.AmountDue.Float64,
			AmountPaid:        installment.AmountPaid.Float64,
			Status:            installment.Status.String,
			InstallmentNumber: installment.InstallmentSeq.Int64,
			TransactionId:     installment.TransactionId.String,
			DueDate:           installment.DueDate.Time.Format("2006-01-02"),
		})
		if installment.Status.String == TXN_PAID {
			response.Data.OutstandingAmount -= installment.AmountPaid.Float64
		}
	}
	response.Message = "successfully fetched installments"
	c.JSON(http.StatusOK, response)
}

func (obj *loanService) ProcessLoanPayment(c *gin.Context) {
	var (
		request  ProcessLoanPaymentRequest
		response ProcessLoanPaymentResponse
	)

	if err := c.BindJSON(&request); err != nil {
		log.Printf("unable to marshal request. Error:%s", err.Error())
		response.Errors = append(response.Errors, *e.ErrorInfo[e.BadRequest])
		response.Message = "failed to process payment"
		c.JSON(http.StatusBadRequest, response)
		return
	}

	//scope: validate transaction id with any service if available
	//get existing installments and check if payment for an installment is valid
	installments, err := obj.dbObj.GetUserLoanInstallments(c, request.UserId, request.LoanId)
	if err != nil {
		log.Printf("failed to fetch loan installments. Error:%s", err.Error())
		response.Errors = append(response.Errors, e.ErrorInfo[e.GetDBError].GetErrorDetails(err.Error()))
		response.Message = "failed to fetch installments to process payment"
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	if len(installments) == 0 {
		response.Message = "no installments against loan available"
		c.JSON(http.StatusBadRequest, response)
		return
	}

	loanPaid := 0.0
	txn := 0
	//find next available installment
	for i, installment := range installments {
		loanPaid += installment.AmountPaid.Float64
		if installment.Status.String == TXN_PENDING {
			txn = i
			break
		}
	}
	if request.Amount < installments[txn].AmountDue.Float64 {
		log.Println("amount payable is less than installment amount")
		response.Errors = append(response.Errors, e.ErrorInfo[e.BadRequest].GetErrorDetails("amount payable is less than installment amount"))
		response.Message = "failed to  process payment"
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	//mark the current txn as paid
	installments[txn].AmountPaid.Float64 = request.Amount
	installments[txn].Status.String = TXN_PAID
	installments[txn].TransactionId.String = request.TransactionId

	loanDue := installments[0].LoanAmount.Float64 - loanPaid - request.Amount
	//if amount paid in installment is so big that it covers more than the entire loan amount, reject the transactions
	if loanDue < 0 {
		log.Println("transaction covers more than loan amount")
		response.Errors = append(response.Errors, e.ErrorInfo[e.BadRequest].GetErrorDetails("transaction repays more than loan amount. transaction not allowed"))
		response.Message = "failed to  process payment"
		c.JSON(http.StatusNotAcceptable, response)
		return
	}

	//if loanDue is greater than 0, make changes to amount due in recurring installments. if not, mark recurring installments as CANCELLED and the loan needs to be marked as PAID
	duePerInstallment := installments[txn].AmountDue.Float64
	loanClosed := true
	if loanDue > 0 {
		duePerInstallment = loanDue / float64(len(installments)-(txn+1))
		loanClosed = false
	}

	//update installment if repayment amount is exactly as due
	if installments[txn].AmountDue.Float64 == request.Amount {
		//update only this installment
		err := obj.dbObj.UpdateSingleInstallmentPayment(c, request.LoanId, installments[txn], loanClosed)
		if err != nil {
			log.Printf("failed to update payment. Error: %s", err.Error())
			response.Errors = append(response.Errors, e.ErrorInfo[e.AddDBError].GetErrorDetails(err.Error()))
			response.Message = "failed to  process payment"
			c.JSON(http.StatusInternalServerError, response)
			return
		}
		response.Status = true
		response.Message = "successfully processed payment"
		c.JSON(http.StatusOK, response)
		return
	}

	//update following transactions with adjusted due amount
	for i := txn + 1; i < len(installments); i++ {
		installments[i].AmountDue.Float64 = duePerInstallment
		if loanDue < 0 {
			installments[i].Status.String = TXN_CANCELLED
		}
	}

	//update these transactions in DB
	err = obj.dbObj.UpdateInstallment(c, request.LoanId, installments[txn:], loanClosed)
	if err != nil {
		log.Printf("failed to update payment. Error: %s", err.Error())
		response.Errors = append(response.Errors, e.ErrorInfo[e.AddDBError].GetErrorDetails(err.Error()))
		response.Message = "failed to  process payment"
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	response.Status = true
	response.Message = "successfully processed payment"
	c.JSON(http.StatusOK, response)
}
