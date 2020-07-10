package repoimpl

import (
	"context"
	"sync"
	"williamfeng323/mooncake-duty/src/infrastructure/db"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//IssueRepo is the implementation of Shift repository
type IssueRepo struct {
	collection *mongo.Collection
}

func init() {
	db.Register(&IssueRepo{})
}

//GetName returns the name of Shift collection
func (i *IssueRepo) GetName() string {
	return "Shift"
}

// SetCollection set the collection that communicate with db to the instance
func (i *IssueRepo) SetCollection(coll *mongo.Collection) {
	i.collection = coll
}

// InsertOne creates new Shift
func (i *IssueRepo) InsertOne(ctx context.Context, document interface{},
	opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	//TO-DO Default validator
	return i.collection.InsertOne(ctx, document, opts...)
}

//Find get one Shift document
func (i *IssueRepo) Find(ctx context.Context, filter interface{},
	opts ...*options.FindOptions) (*mongo.Cursor, error) {
	return i.collection.Find(ctx, filter, opts...)
}

// FindOne returns the single document that meets the criteria
func (i *IssueRepo) FindOne(ctx context.Context, filter interface{},
	opts ...*options.FindOneOptions) *mongo.SingleResult {
	return i.collection.FindOne(ctx, filter, opts...)
}

// UpdateOne update the document according to the filter.
func (i *IssueRepo) UpdateOne(ctx context.Context, filter interface{}, update interface{},
	opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return i.collection.UpdateOne(ctx, filter, update, opts...)
}

// DeleteOne executes a delete command to delete at most one document from the collection.
func (i *IssueRepo) DeleteOne(ctx context.Context, filter interface{},
	opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return i.collection.DeleteOne(ctx, filter, opts...)
}

var issueRepo *IssueRepo
var issueLock sync.RWMutex

// GetIssueRepo get the account repository instance,
// Create if not exist.
func GetIssueRepo() *IssueRepo {
	issueLock.Lock()
	defer issueLock.Unlock()
	if issueRepo == nil {
		issueRepo = &IssueRepo{}
		issueRepo.SetCollection(db.GetConnection().GetCollection("Shift"))
	}
	return issueRepo
}
