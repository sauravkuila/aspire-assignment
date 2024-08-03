package loan

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type loanDb struct {
	dbObj *gorm.DB
}

type DbLoanInterface interface {
	CreateLoan(*gin.Context, int64, float64, int64) (int64, error)
	ModifyLoan(*gin.Context, int64, int64, float64, int64) (int64, error)
	CancelLoan(*gin.Context, int64, int64) (int64, error)
	GetAllLoansForAgainstUser(*gin.Context, int64) ([]LoanDetails, error)
}

func NewLoanDbObject(db *gorm.DB) DbLoanInterface {
	return &loanDb{
		dbObj: db,
	}
}
