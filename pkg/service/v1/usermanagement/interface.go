package usermanagement

import (
	v1 "aspire-assignment/pkg/db/v1"

	"github.com/gin-gonic/gin"
)

type userMgtService struct {
	dbObj v1.V1DBLayer
}

type UserManagementInterface interface {
	UserSignup(*gin.Context)
	UserLogin(*gin.Context)
}

func NewUserManagementService(db v1.V1DBLayer) UserManagementInterface {
	return &userMgtService{
		dbObj: db,
	}
}
