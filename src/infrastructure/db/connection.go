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

type collectionRegistry struct {
	registries map[string]*mongo.Collection
	lock       sync.Mutex
}

//Connection provide the client to connect the database
type Connection struct {
	Database string
	Client   *mongo.Client
	*collectionRegistry
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
	conn.collectionRegistry = &collectionRegistry{registries: map[string]*mongo.Collection{}}
	return nil
}

//Register will create the collection connection registered it to the connection.
// Before you can operate your repo, you must register it.
func (conn *Connection) register(repo Repository) {

	if repo == nil {
		panic("repo can not be nil")
	}

	repoName := repo.GetName()

	conn.collectionRegistry.lock.Lock()
	defer conn.collectionRegistry.lock.Unlock()
	// check if model was already registered, if not, register the model
	// into CollectionRegistry[modelName]
	if _, ok := conn.collectionRegistry.registries[repoName]; !ok {
		collection := conn.Client.Database(conn.Database).Collection(repoName)
		conn.collectionRegistry.registries[repoName] = collection
		fmt.Printf("Registered collection '%v'", repoName)
	} else {
		fmt.Printf("Tried to register collection '%v' twice", repoName)
	}
}

// GetCollection returns the registered repo's collection instance
func (conn *Connection) GetCollection(repositoryName string) *mongo.Collection {
	if collection, ok := conn.collectionRegistry.registries[repositoryName]; ok {
		return collection
	}
	panic("repository does not registered, please call GetConnection().Register(repo) to register")
}

var connection Connection
var lock sync.RWMutex

func init() {
	GetConnection()
}

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

// Register creates and registers the repo's collection to the global connection
func Register(repo Repository) {
	lock.Lock()
	defer lock.Unlock()
	connection.register(repo)
}
