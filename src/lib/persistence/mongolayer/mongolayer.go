package mongolayer

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jalexanderII/MyEventsMicro/src/lib/persistence"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	EVENTS      = "events"
	Performance = 100
)

var (
	// Used to create a singleton object of MongoDB client.
	// Initialized and exposed through  GetMongoClient()
	clientInstance *mongo.Client
	// Used during creation of singleton client object in GetMongoClient()
	clientInstanceError error
	// Used to execute client creation procedure only once
	mongoOnce sync.Once
)

type MongoDBLayer struct {
	session mongo.Database
	ctx     context.Context
}

func NewMongoDBLayer(connection, database string) (persistence.DatabaseHandler, error) {
	// Perform connection creation operation only once.
	mongoOnce.Do(func() {
		// Set client options
		clientOptions := options.Client().ApplyURI(connection)
		ctx, cancel := NewDBContext(10 * time.Second)
		defer cancel()
		// Connect to MongoDB
		client, err := mongo.Connect(ctx, clientOptions)
		if err != nil {
			clientInstanceError = err
		}
		// Check the connection
		err = client.Ping(ctx, nil)
		if err != nil {
			clientInstanceError = err
		}
		clientInstance = client
	})

	return &MongoDBLayer{
		session: *clientInstance.Database(database), ctx: context.Background(),
	}, clientInstanceError
}

func (mgoLayer *MongoDBLayer) AddEvent(e persistence.Event) ([]byte, error) {
	if e.ID.IsZero() {
		e.ID = primitive.NewObjectID()
	}

	if e.Location.ID.IsZero() {
		e.Location.ID = primitive.NewObjectID()
	}

	_, err := mgoLayer.session.Collection(EVENTS).InsertOne(mgoLayer.ctx, e)
	if err != nil {
		fmt.Println("Error inserting new PaymentTask", err)
		return nil, err
	}

	return []byte(e.ID.Hex()), nil
}

func (mgoLayer *MongoDBLayer) FindEvent(id []byte) (persistence.Event, error) {
	filter := bson.D{{Key: "_id", Value: primitive.ObjectIDFromHex(string(id))}}
	e := persistence.Event{}
	err := mgoLayer.session.Collection(EVENTS).FindOne(mgoLayer.ctx, filter).Decode(&e)
	return e, err
}

func (mgoLayer *MongoDBLayer) FindEventByName(name string) (persistence.Event, error) {
	filter := bson.D{{Key: "name", Value: name}}
	e := persistence.Event{}
	err := mgoLayer.session.Collection(EVENTS).FindOne(mgoLayer.ctx, filter).Decode(&e)
	return e, err
}

func (mgoLayer *MongoDBLayer) FindAllAvailableEvents() ([]persistence.Event, error) {
	var results []persistence.Event
	cursor, err := mgoLayer.session.Collection(EVENTS).Find(mgoLayer.ctx, bson.D{})
	if err != nil {
		return results, err
	}
	if err = cursor.All(mgoLayer.ctx, &results); err != nil {
		fmt.Println("[PaymentTaskDB] Error getting all users", "error", err)
		return results, err
	}
	return results, err
}

// NewDBContext returns a new Context according to app performance
func NewDBContext(d time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), d*Performance/100)
}
