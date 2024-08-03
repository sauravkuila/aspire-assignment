package usermanagement

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type userMgtDb struct {
	dbObj *gorm.DB
}

type DbUserManagementInterface interface {
	FuncUserMgtSample(*gin.Context) error
}

func NewLoanDbObject(db *gorm.DB) DbUserManagementInterface {
	return &userMgtDb{
		dbObj: db,
	}
}
