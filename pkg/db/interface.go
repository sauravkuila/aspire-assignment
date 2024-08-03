package db

import (
	v1 "aspire-assignment/pkg/db/v1"

	"gorm.io/gorm"
)

type dbService struct {
	V1 v1.V1DBLayer //v1
}

func NewDBObject(psqlDB *gorm.DB) DBLayer {
	temp := &dbService{
		v1.NewV1DbLayer(psqlDB),
	}
	return temp
}

type DBLayer interface {
	GetV1DBLayer() v1.V1DBLayer //v1 db layer
}

func (obj *dbService) GetV1DBLayer() v1.V1DBLayer {
	return obj.V1
}
