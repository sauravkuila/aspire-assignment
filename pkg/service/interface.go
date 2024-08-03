package service

import (
	"aspire-assignment/pkg/db"
	v1 "aspire-assignment/pkg/service/v1"

	"github.com/gin-gonic/gin"
)

type service struct {
	V1 v1.ServiceLayer //v1
}

type ServiceGroupLayer interface {
	GetV1Service() v1.ServiceLayer //v1
	Health(*gin.Context)
}

func NewServiceGroupObject(db db.DBLayer) ServiceGroupLayer {
	return &service{
		v1.NewServiceObject(db.GetV1DBLayer()),
	}
}

func (s *service) GetV1Service() v1.ServiceLayer {
	return s.V1
}
