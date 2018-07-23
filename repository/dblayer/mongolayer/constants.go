package mongolayer

import (
	"log"

	mgo "gopkg.in/mgo.v2"
)

const (
	DB     = "myevents"
	USERS  = "users"
	EVENTS = "events"
)

// MongoDBLayer - Struct to make use of session
type MongoDBLayer struct {
	session *mgo.Session
}

// NewMongoDBLayer returns handle to database handler interface - mongodb://127.0.0.1
func NewMongoDBLayer(connection string) (*MongoDBLayer, error) {
	s, err := mgo.Dial(connection)
	if err != nil {
		log.Println("Error connecting to MongoDB -", err)
		return nil, err
	}

	return &MongoDBLayer{
		session: s,
	}, err
}

func (m *MongoDBLayer) getFreshSession() *mgo.Session {
	return m.session.Copy()
}


