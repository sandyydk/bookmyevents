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
	log.Println("Starting EVents API Server")
	log.Fatal(servicehandler.ServeAPI(conf.RestfulEndpoint, dbHandler))
}
