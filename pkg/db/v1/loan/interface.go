package loan

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type loanDb struct {
	dbObj *gorm.DB
}

type DbLoanInterface interface {
	FuncLoanSample(*gin.Context) error
}

func NewLoanDbObject(db *gorm.DB) DbLoanInterface {
	return &loanDb{
		dbObj: db,
	}
}
