package dao

import (
	"context"
	"fmt"
	"net/url"
	"time"

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
	Database string
	Client   *mongo.Client
}

//Database is the object you get the MongoDB client/connection
var (
	Database Connection
)

// InitConnection initial the connection to the database
func (conn *Connection) InitConnection(ctx context.Context, config utils.MongoConfig) error {
	if ctx == nil {
		ctx, _ = context.WithTimeout(context.Background(), 30*time.Second)
	}
	u := url.URL{
		Scheme:   "mongodb",
		User:     url.UserPassword(config.Username, config.Password),
		Host:     fmt.Sprintf("%s:%s", config.URL, config.Port),
		Path:     config.Database,
		RawQuery: config.ConnectionOptions,
	}
	clientOption := options.Client().ApplyURI(u.String())
	client, err := mongo.Connect(ctx, clientOption)
	if err != nil {
		return err
	}
	conn.Database = config.Database
	conn.Client = client
	return nil
}

//Register will register the collection connection to the Connection.
// func (conn *Connection) Register(document IDocumentBase, collectionName string) {

// 	if document == nil {
// 		panic("document can not be nil")
// 	}

// 	reflectType := reflect.TypeOf(document)
// 	typeName := reflectType.Elem().Name()

// 	//check if model was already registered
// 	if _, ok := conn.collectionRegistry[typeName]; !ok {
// 		collection := &Collection{conn.client.Database(conn.database).Collection(typeName)}
// 		conn.collectionRegistry[typeName] = collection
// 		fmt.Printf("Registered collection '%v'", typeName)
// 	} else {
// 		fmt.Printf("Tried to register collection '%v' twice", typeName)
// 	}
// }
