package models

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//https://github.com/entropyx/mango/blob/master/mango.go
//Connection provide the client to connect the database
type Connection struct {
	Client *mongo.Client
}

// Connect create a connection to the database
func (connection *Connection) Connect(config *options.ClientOptions) error {
	return nil
}
