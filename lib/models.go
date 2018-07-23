package lib

import (
	"reflect"

	"gopkg.in/mgo.v2/bson"
)

type Event struct {
	ID        bson.ObjectId `bson:"_id"`
	Name      string
	Duration  int
	StartDate int64
	EndDate   int64
	Location  Location
}

type Location struct {
	ID        bson.ObjectId `bson:"_id"`
	Name      string
	Address   string
	Country   string
	OpenTime  int
	CloseTime int
	Halls     []Hall
}

type Hall struct {
	ID       bson.ObjectId `bson:"_id"`
	Name     string        `json:"name"`
	Location string        `json:"location,omitempty"`
	Capacity int           `json:"capacity"`
}

// Valid returns valididty of an Event
func (e Event) Valid() bool {
	if e.Duration == 0 || e.EndDate == 0 || e.StartDate == 0 || e.StartDate > e.EndDate || e.StartDate == e.EndDate || e.Name == "" || reflect.DeepEqual(e.Location, Location{}) {
		return false
	}
	return true
}
