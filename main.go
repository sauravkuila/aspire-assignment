package main

import (
	"aspire-assignment/api"
	"aspire-assignment/pkg/config"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var (
		environment string
	)
	if len(os.Args) == 2 {
		environment = os.Args[1] // developer custom file
	} else {
		environment = "local"
	}

	config.Load(environment)

	if err := api.Start(); err != nil {
		log.Fatal("Failed to start server, err:", err)
		os.Exit(1)
	}
	addShutdownHook()
}

func addShutdownHook() {
	// when receive interruption from system shutdown server and scheduler
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-quit
	log.Println("Quit/Interrupt signal detected. Gracefully closing connections")
	//shutdown server
	api.ShutdownRouter()
	api.CloseDatabase()

	log.Printf("All done! Wrapping up here for PID: %d", os.Getpid())
}
