package models

import (
	"context"
	"fmt"
	"net/url"

	"williamfeng323/mooncake-duty/src/utils"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Connection provide the client to connect the database
type Connection struct {
	client   *mongo.Client
	database string
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
	conn := &Connection{client: client, database: config.Database}
	return conn, nil
}
