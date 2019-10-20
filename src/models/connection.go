package models

import (
	"context"
	"fmt"
	"net/url"
	"reflect"

	"williamfeng323/mooncake-duty/src/utils"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Collection is the struct to contain correlated collection connection
//and also the collection operate
type Collection struct {
	*mongo.Collection
}

//Connection provide the client to connect the database
type Connection struct {
	database           string
	client             *mongo.Client
	collectionRegistry map[string]*Collection
	typeRegistry       map[string]reflect.Type
}

// Connect create a connection to the database
func Connect(ctx context.Context, config utils.MongoConfig) (*Connection, error) {
	u := url.URL{
		Scheme:   "mongodb",
		User:     url.UserPassword(config.Username, config.Password),
		Host:     fmt.Sprintf("%s:%s", config.URL, config.Port),
		Path:     config.Database,
		RawQuery: config.ConnectionOptions,
	}
	if ctx == nil {
		ctx = context.Background()
	}
	clientOption := options.Client().ApplyURI(u.String())
	client, err := mongo.Connect(ctx, clientOption)
	if err != nil {
		return nil, err
	}
	conn := &Connection{
		database:           config.Database,
		client:             client,
		collectionRegistry: make(map[string]*Collection),
		typeRegistry:       make(map[string]reflect.Type),
	}
	return conn, nil
}

//Register will register the collection connection to the Connection.
func (conn *Connection) Register(document IDocumentBase, collectionName string) {

	if document == nil {
		panic("document can not be nil")
	}

	reflectType := reflect.TypeOf(document)
	typeName := reflectType.Elem().Name()

	//check if model was already registered
	if _, ok := conn.collectionRegistry[typeName]; !ok {
		collection := &Collection{conn.client.Database(conn.database).Collection(typeName)}
		conn.collectionRegistry[typeName] = collection
		fmt.Printf("Registered collection '%v'", typeName)
	} else {
		fmt.Printf("Tried to register collection '%v' twice", typeName)
	}
}
