package mongolayer

import (
	models "bookmyevents/lib"

	"gopkg.in/mgo.v2/bson"
)

func (m *MongoDBLayer) AddEvent(e models.Event) ([]byte, error) {
	session := m.getFreshSession()
	defer session.Close()

	if !e.ID.Valid() {
		e.ID = bson.NewObjectId()
	}

	if !e.Location.ID.Valid() {
		e.Location.ID = bson.NewObjectId()
	}

	return []byte(e.ID), session.DB(DB).C(EVENTS).Insert(e)
}

func (m *MongoDBLayer) FindEvent(id []byte) (models.Event, error) {
	session := m.getFreshSession()
	defer session.Close()

	var event models.Event

	err := session.DB(DB).C(EVENTS).FindId(bson.ObjectId(id)).One(&event)
	return event, err
}

func (m *MongoDBLayer) FindEventByName(name string) (models.Event, error) {
	session := m.getFreshSession()
	defer session.Close()

	var event models.Event

	err := session.DB(DB).C(EVENTS).Find(bson.M{"name": name}).One(&event)
	return event, err
}

func (m *MongoDBLayer) FindAllAvailableEvents() ([]models.Event, error) {
	session := m.getFreshSession()
	defer session.Close()

	var events []models.Event

	err := session.DB(DB).C(EVENTS).Find(nil).All(&events)
	return events, err
}
