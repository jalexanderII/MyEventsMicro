package mongolayer

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	persistence2 "github.com/jalexanderII/MyEventsMicro/src/lib/persistence"
)

const (
	DB     = "MyEvents"
	USERS  = "users"
	EVENTS = "events"
)

type MongoDBLayer struct {
	session *mgo.Session
}

func NewMongoDBLayer(connection string) (persistence2.DatabaseHandler, error) {
	s, err := mgo.Dial(connection)
	return &MongoDBLayer{
		session: s,
	}, err
}

func (mgoLayer *MongoDBLayer) AddEvent(e persistence2.Event) ([]byte, error) {
	s := mgoLayer.getFreshSession()
	defer s.Close()

	if !e.ID.Valid() {
		e.ID = bson.NewObjectId()
	}

	if !e.Location.ID.Valid() {
		e.Location.ID = bson.NewObjectId()
	}

	return []byte(e.ID), s.DB(DB).C(EVENTS).Insert(e)
}

func (mgoLayer *MongoDBLayer) FindEvent(id []byte) (persistence2.Event, error) {
	s := mgoLayer.getFreshSession()
	defer s.Close()
	e := persistence2.Event{}

	err := s.DB(DB).C(EVENTS).FindId(bson.ObjectId(id)).One(&e)
	return e, err
}

func (mgoLayer *MongoDBLayer) FindEventByName(name string) (persistence2.Event, error) {
	s := mgoLayer.getFreshSession()
	defer s.Close()
	e := persistence2.Event{}
	err := s.DB(DB).C(EVENTS).Find(bson.M{"name": name}).One(&e)
	return e, err
}

func (mgoLayer *MongoDBLayer) FindAllAvailableEvents() ([]persistence2.Event, error) {
	s := mgoLayer.getFreshSession()
	defer s.Close()
	var events []persistence2.Event
	err := s.DB(DB).C(EVENTS).Find(nil).All(&events)
	return events, err
}

func (mgoLayer *MongoDBLayer) getFreshSession() *mgo.Session {
	return mgoLayer.session.Copy()
}
