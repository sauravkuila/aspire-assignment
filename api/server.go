package api

import (
	"aspire-assignment/pkg/config"
	"aspire-assignment/pkg/db"
	"aspire-assignment/pkg/service"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"gorm.io/gorm"
)

var srv *http.Server
var ctx context.Context
var databases []*gorm.DB

func Start() error {
	ctx = context.Background()

	databases = make([]*gorm.DB, 0)
	postgresConn, err := db.PsqlConnect()
	if err != nil {
		log.Printf("Failed to connect psql database", err.Error())
		return err
	}

	databases = append(databases, postgresConn)

	dbObj := db.NewDBObject(postgresConn)

	serviceObj := service.NewServiceGroupObject(dbObj)

	startRouter(serviceObj)
	return nil
}

func startRouter(obj service.ServiceGroupLayer) {
	srv = &http.Server{
		Addr:    fmt.Sprintf(":%d", config.GetConfig().GetInt("server.port")),
		Handler: getRouter(obj), //getRouter set the api specs for version-1 routes
	}
	// run api router
	log.Println("starting router")
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server %s", err.Error())
		}
	}()
}

func ShutdownRouter() {
	timeoutCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	log.Println("Shutting down router START")
	defer log.Println("Shutting down router END")
	if err := srv.Shutdown(timeoutCtx); err != nil {
		log.Fatalf("Server forced to shutdown. Error: %s", err.Error())
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-timeoutCtx.Done():
		log.Println("timeout of 2 seconds.")
	}
}

func CloseDatabase() {
	log.Println("disconnecting databases START")
	defer log.Println("disconnecting databases END")
	for _, database := range databases {
		db, _ := database.DB()
		if db != nil {
			err := db.Close()
			if err != nil {
				log.Printf("unable to close db. Error: %s", err.Error())
			}
		} else {
			log.Println("unable to close db as connection is nil")

		}
	}
}
