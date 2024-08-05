package v1

import (
	"aspire-assignment/pkg/db/v1/loan"
	"aspire-assignment/pkg/db/v1/usermanagement"

	"gorm.io/gorm"
)

type dbV1LayerObj struct {
	loan.DbLoanInterface
	usermanagement.DbUserManagementInterface
}

//go:generate mockgen -destination=mock/mock.go -package=mock aspire-assignment/pkg/db/v1  V1DBLayer
type V1DBLayer interface {
	loan.DbLoanInterface
	usermanagement.DbUserManagementInterface
}

func NewV1DbLayer(db *gorm.DB) V1DBLayer {
	return dbV1LayerObj{
		loan.NewLoanDbObject(db),
		usermanagement.NewLoanDbObject(db),
	}
}
