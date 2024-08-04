package usermanagement

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type userMgtDb struct {
	dbObj *gorm.DB
}

type DbUserManagementInterface interface {
	AddUser(*gin.Context, UserDetails) (int64, error)
	GetUserByUsername(*gin.Context, string) (UserDetails, error)
}

func NewLoanDbObject(db *gorm.DB) DbUserManagementInterface {
	return &userMgtDb{
		dbObj: db,
	}
}
