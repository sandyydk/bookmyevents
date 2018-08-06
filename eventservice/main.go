package main

import (
	"bookmyevents/config"
	servicehandler "bookmyevents/eventservice/handlers"
	msgqueue "bookmyevents/lib/msgqueue/amqp"
	"bookmyevents/repository"
	"flag"
	"log"

	"github.com/streadway/amqp"
)

func main() {
	configPath := flag.String("config", `.\config\config.json`, "flag to set the path to the configuration json file")
	flag.Parse()

	conf, _ := config.ExtractConfig(*configPath)
	log.Println("Configuration parsed")

	conn, err := amqp.Dial(conf.AMQPMessageBroker)
	if err != nil {
		log.Println("Error dialing amqp -", err.Error())
		panic(err)
	}

	emitter, err := msgqueue.NewAMQPEventEmitter(conn)
	if err != nil {
		log.Println("Error connecting to amqp -", err.Error())
		panic(err)
	}

	dbHandler, _ := repository.NewDBLayer(conf.Databasetype, conf.DBConnection)

	// Start REST API
	log.Println("Starting Events API Server")

	httpErrorChan, httpTLSErrorChan := servicehandler.ServeAPI(conf.RestfulTLSEndpoint, conf.RestfulEndpoint, dbHandler, emitter)

	select {
	case err := <-httpErrorChan:
		log.Fatal("HTTP Error found:", err)
	case err := <-httpTLSErrorChan:
		log.Fatal("HTTPS Error found:", err)
	}
}
