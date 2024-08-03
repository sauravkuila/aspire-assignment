package db

import (
	"aspire-assignment/pkg/config"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type dbLogger struct{}

func PsqlConnect() (*gorm.DB, error) {

	c := config.GetConfig()

	// Init DB connection string params from config
	dsn := fmt.Sprintf("host=%s port=%d user=%s password='%s' dbname=%s sslmode=%s connect_timeout=%d",
		c.GetString("databases.postgres.host"), c.GetInt("databases.postgres.port"), c.GetString("databases.postgres.user"), c.GetString("databases.postgres.password"), c.GetString("databases.postgres.db"), c.GetString("databases.postgres.sslmode"), c.GetInt("databases.postgres.connect_timeout"))

	// Connect to database
	dbc, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		SkipDefaultTransaction: true, /// Removed defuault transaction used by GORM to faster the query and execution
	})
	if err != nil {
		log.Printf("Failed to connect to database. Error: %s, conn: %s", err.Error(), dsn)
		return nil, err
	}

	log.Println("Postgres Database Connected")
	return dbc, nil
}
