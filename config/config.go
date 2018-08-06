package config

import (
	"bookmyevents/repository"
	"encoding/json"
	"log"
	"os"
)

var (
	DBTypeDefault       = repository.DBTYPE("mongodb")
	DBConnectionDefault = "mongodb://127.0.0.1"
	RestfulEPDefault    = "localhost:8181"
	RestfulTLSEndpoint  = "localhost:8182"
	AMQPMessageBroker   = ""
)

// ServiceConfig struct for service configuration to be read from a json file
type ServiceConfig struct {
	Databasetype       repository.DBTYPE `json:"databasetype"`
	DBConnection       string            `json:"dbconnection"`
	RestfulEndpoint    string            `json:restfulapi_endpoint`
	RestfulTLSEndpoint string            `json:"restfultls_endpoint"`
	AMQPMessageBroker  string            `json:"amqp_message_broker"`
}

// ExtractConfig extracts config from a json file or returns the default
func ExtractConfig(filename string) (ServiceConfig, error) {
	conf := ServiceConfig{
		Databasetype:       DBTypeDefault,
		DBConnection:       DBConnectionDefault,
		RestfulEndpoint:    RestfulEPDefault,
		RestfulTLSEndpoint: RestfulTLSEndpoint,
		AMQPMessageBroker:  AMQPMessageBroker,
	}

	file, err := os.Open(filename)
	if err != nil {
		log.Println("Error opening config file", err)
		return conf, err
	}

	err = json.NewDecoder(file).Decode(&conf)
	if broker := os.Getenv("AMQP_URL"); broker != "" {
		conf.AMQPMessageBroker = broker
	}
	return conf, err
}
