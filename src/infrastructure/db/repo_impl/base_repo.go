package repoimpl

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// BaseRepo implements the basic repository features
type BaseRepo struct {
	name       string
	collection *mongo.Collection
}

// SetCollection will set the document's own collection instance
func (repo *BaseRepo) SetCollection(coll *mongo.Collection) {
	repo.collection = coll
}

// GetName returns the name of the repository
func (repo *BaseRepo) GetName() string {
	return repo.name
}

//Find return the result cursor base on your query
func (repo *BaseRepo) Find(query ...interface{}) (*mongo.Cursor, error) {
	finalQuery := initQuery(query)
	ctx, cancelFunc := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelFunc()
	//accept zero or one query param
	cursor, err := repo.collection.Find(ctx, finalQuery)
	if err != nil {
		return nil, err
	}
	return cursor, nil
}

//FindByID find document by document _id
func (repo *BaseRepo) FindByID(id primitive.ObjectID) *mongo.SingleResult {
	filter := bson.M{"_id": id}
	ctx, cancelFunc := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelFunc()
	return repo.collection.FindOne(ctx, filter)
}

//DeleteByID find document by document _id
func (repo *BaseRepo) DeleteByID(id primitive.ObjectID) (*mongo.DeleteResult, error) {
	filter := bson.M{"_id": id}
	ctx, cancelFunc := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelFunc()
	return repo.collection.DeleteOne(ctx, filter)
}

//Delete deletes the documents according to the query provided.
func (repo *BaseRepo) Delete(query ...interface{}) (*mongo.DeleteResult, error) {
	finalQuery := initQuery(query)
	ctx, cancelFunc := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelFunc()
	rst, err := repo.collection.DeleteOne(ctx, finalQuery)
	return rst, err
}

func initQuery(query []interface{}) interface{} {
	var finalQuery interface{}
	//accept zero or one query param
	if len(query) == 0 {
		finalQuery = bson.M{}
	} else if len(query) == 1 {
		finalQuery = query[0]
	} else {
		panic("DB: Find method accepts no or maximum one query param.")
	}
	return finalQuery
}
