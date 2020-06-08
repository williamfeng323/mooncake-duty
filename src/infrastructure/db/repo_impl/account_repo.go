package repoimpl

import (
	"context"
	"sync"

	"williamfeng323/mooncake-duty/src/infrastructure/db"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//AccountRepo is the implementation of Project repository
type AccountRepo struct {
	collection *mongo.Collection
}

func init() {
	db.Register(&AccountRepo{})
}

//GetName returns the name of project collection
func (pr *AccountRepo) GetName() string {
	return "Account"
}

// SetCollection set the collection that communicate with db to the instance
func (pr *AccountRepo) SetCollection(coll *mongo.Collection) {
	pr.collection = coll
}

// InsertOne creates new project
func (pr *AccountRepo) InsertOne(ctx context.Context, document interface{},
	opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return pr.collection.InsertOne(ctx, document, opts...)
}

//Find get project documents that meet the criteria
func (pr *AccountRepo) Find(ctx context.Context, filter interface{},
	opts ...*options.FindOptions) (*mongo.Cursor, error) {
	return pr.collection.Find(ctx, filter, opts...)
}

// FindOne returns the single document that meets the criteria
func (pr *AccountRepo) FindOne(ctx context.Context, filter interface{},
	opts ...*options.FindOneOptions) *mongo.SingleResult {
	return pr.collection.FindOne(ctx, filter, opts...)
}

// UpdateOne update the document according to the filter.
func (pr *AccountRepo) UpdateOne(ctx context.Context, filter interface{}, update interface{},
	opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return pr.collection.UpdateOne(ctx, filter, update, opts...)
}

// DeleteOne executes a delete command to delete at most one document from the collection.
func (pr *AccountRepo) DeleteOne(ctx context.Context, filter interface{},
	opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return pr.collection.DeleteOne(ctx, filter, opts...)
}

var accountRepo *AccountRepo
var lock sync.RWMutex

// GetAccountRepo get the account repository instance,
// Create if not exist.
func GetAccountRepo() *AccountRepo {
	lock.Lock()
	defer lock.Unlock()
	if accountRepo == nil {
		accountRepo = &AccountRepo{}
		accountRepo.SetCollection(db.GetConnection().GetCollection("Account"))
	}
	return accountRepo
}
