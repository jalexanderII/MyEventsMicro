package dblayer

import (
	"github.com/jalexanderII/MyEventsMicro/src/lib/persistence"
	"github.com/jalexanderII/MyEventsMicro/src/lib/persistence/mongolayer"
)

type DBTYPE string

const (
	MONGODB  DBTYPE = "mongodb"
	DYNAMODB DBTYPE = "dynamodb"
)

func NewPersistenceLayer(options DBTYPE, connection, databasename string) (persistence.DatabaseHandler, error) {

	switch options {
	case MONGODB:
		return mongolayer.NewMongoDBLayer(connection, databasename)
	}
	return nil, nil
}
