package db

import (
	"context"
	"fmt"
	"net/url"
	"reflect"
	"sync"
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
	Database           string
	Client             *mongo.Client
	CollectionRegistry map[string]*Collection
}

// InitConnection initial the connection to the database
func (conn *Connection) InitConnection(ctx context.Context, config utils.MongoConfig) error {
	var contextCancel context.CancelFunc
	if ctx == nil {
		ctx, contextCancel = context.WithTimeout(context.Background(), 30*time.Second)
	}
	defer contextCancel()
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
	conn.CollectionRegistry = map[string]*Collection{}
	return nil
}

//Register will register the collection connection to the Connection.
// Before you can operate your model, you must register it.
func (conn *Connection) Register(document IDocumentBase) {

	if document == nil {
		panic("document can not be nil")
	}

	reflectType := reflect.TypeOf(document)
	typeName := reflectType.Elem().Name()

	// check if model was already registered, if not, register the model
	// into CollectionRegistry[modelName]
	if _, ok := conn.CollectionRegistry[typeName]; !ok {
		collection := &Collection{
			conn.Client.Database(conn.Database).Collection(typeName)}
		conn.CollectionRegistry[typeName] = collection
		fmt.Printf("Registered collection '%v'", typeName)
	} else {
		fmt.Printf("Tried to register collection '%v' twice", typeName)
	}
}

// GetCollection returns the registered collection
func (conn *Connection) GetCollection(collectionName string) *Collection {
	if collection, ok := conn.CollectionRegistry[collectionName]; ok {
		return collection
	} else {
		panic("Account does not registered to the DB")
	}
}

// New return a document instance
// To new a document, you should follow below steps:
// connection.Register(&User{})
// user := &User{}
// connection.CollectionRegistry["User"].New(user)
func (coll *Collection) New(doc IDocumentBase) error {
	doc.SetCollection(coll.Collection)
	doc.SetDocument(doc)
	return nil
}

var connection Connection
var lock sync.RWMutex

//GetConnection return the initted DB connection struct
func GetConnection() *Connection {
	lock.Lock()
	defer lock.Unlock()
	if reflect.DeepEqual(connection, reflect.Zero(reflect.TypeOf(connection)).Interface()) {
		connection = Connection{}
		connection.InitConnection(nil, utils.GetConf().Mongo)
		return &connection
	}
	return &connection
}
