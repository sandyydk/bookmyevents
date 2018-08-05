package main

import (
	"bookmyevents/config"
	"bookmyevents/repository"
	servicehandler "bookmyevents/servicehandlers/event"
	"flag"
	"log"
)

func main() {
	configPath := flag.String("config", `.\config\config.json`, "flag to set the path to the configuration json file")
	flag.Parse()

	conf, _ := config.ExtractConfig(*configPath)
	log.Println("Configuration parsed")

	dbHandler, _ := repository.NewDBLayer(conf.Databasetype, conf.DBConnection)

	// Start REST API
	log.Println("Starting Events API Server")

	httpErrorChan, httpTLSErrorChan := servicehandler.ServeAPI(conf.RestfulTLSEndpoint, conf.RestfulEndpoint, dbHandler)

	select {
	case err := <-httpErrorChan:
		log.Fatal("HTTP Error found:", err)
	case err := <-httpTLSErrorChan:
		log.Fatal("HTTPS Error found:", err)
	}
}
