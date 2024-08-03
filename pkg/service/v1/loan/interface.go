package loan

import (
	v1 "aspire-assignment/pkg/db/v1"

	"github.com/gin-gonic/gin"
)

type loanService struct {
	dbObj v1.V1DBLayer
}

type LoanInterface interface {
	CreateLoan(*gin.Context)
	ModifyLoan(*gin.Context)
	CancelLoan(*gin.Context)
	GetLoans(*gin.Context)
}

func NewLoanService(db v1.V1DBLayer) LoanInterface {
	return &loanService{
		dbObj: db,
	}
}
