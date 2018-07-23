package repository

import (
	"bookmyevents/repository/dblayer/mongolayer"
)

type DBTYPE string

const (
	MONGODB DBTYPE = "mongodb"
)

func NewDBLayer(options DBTYPE, connectionString string) (mongolayer.DatabaseHandler, error) {
	switch options {
	case MONGODB:
		return mongolayer.NewMongoDBLayer(connectionString)
	}

	return nil, nil
}
