package v1

import (
	v1 "aspire-assignment/pkg/db/v1"
	"aspire-assignment/pkg/service/v1/loan"
	"aspire-assignment/pkg/service/v1/usermanagement"
)

type serviceObj struct {
	loan.LoanInterface
	usermanagement.UserManagementInterface
}

type ServiceLayer interface {
	loan.LoanInterface
	usermanagement.UserManagementInterface
}

func NewServiceObject(db v1.V1DBLayer) ServiceLayer {
	return &serviceObj{
		loan.NewLoanService(db),
		usermanagement.NewUserManagementService(db),
	}
}
