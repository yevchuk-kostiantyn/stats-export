package main

import (
	"fmt"
	"log"
	"os"

	"github.com/yevchuk-kostiantyn/stats-export/rest/edbo"
	"github.com/yevchuk-kostiantyn/stats-export/service"
)

const (
	defaultPort = "8080"
)

func main() {
	port := getPort()
	log.Printf("{port=%s} Starting HTTP server", port)
	bindAddress := fmt.Sprintf(":%s", port)
	edboClient, err := edbo.NewClient()
	if err != nil {
		log.Fatal(err.Error())
	}
	svc := service.NewService(edboClient)
	svc.Run(bindAddress)
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	return port
}
