package repoimpl

import (
	"context"
	"sync"
	"williamfeng323/mooncake-duty/src/infrastructure/db"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//ShiftRepo is the implementation of Shift repository
type ShiftRepo struct {
	collection *mongo.Collection
}

func init() {
	db.Register(&ShiftRepo{})
}

//GetName returns the name of Shift collection
func (shifta *ShiftRepo) GetName() string {
	return "Shift"
}

// SetCollection set the collection that communicate with db to the instance
func (shifta *ShiftRepo) SetCollection(coll *mongo.Collection) {
	shifta.collection = coll
}

// InsertOne creates new Shift
func (shifta *ShiftRepo) InsertOne(ctx context.Context, document interface{},
	opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	//TO-DO Default validator
	return shifta.collection.InsertOne(ctx, document, opts...)
}

//Find get one Shift document
func (shifta *ShiftRepo) Find(ctx context.Context, filter interface{},
	opts ...*options.FindOptions) (*mongo.Cursor, error) {
	return shifta.collection.Find(ctx, filter, opts...)
}

// FindOne returns the single document that meets the criteria
func (shifta *ShiftRepo) FindOne(ctx context.Context, filter interface{},
	opts ...*options.FindOneOptions) *mongo.SingleResult {
	return shifta.collection.FindOne(ctx, filter, opts...)
}

// UpdateOne update the document according to the filter.
func (shifta *ShiftRepo) UpdateOne(ctx context.Context, filter interface{}, update interface{},
	opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return shifta.collection.UpdateOne(ctx, filter, update, opts...)
}

// DeleteOne executes a delete command to delete at most one document from the collection.
func (shifta *ShiftRepo) DeleteOne(ctx context.Context, filter interface{},
	opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return shifta.collection.DeleteOne(ctx, filter, opts...)
}

var shiftRepo *ShiftRepo
var shiftLock sync.RWMutex

// GetShiftRepo get the account repository instance,
// Create if not exist.
func GetShiftRepo() *ShiftRepo {
	shiftLock.Lock()
	defer shiftLock.Unlock()
	if shiftRepo == nil {
		shiftRepo = &ShiftRepo{}
		shiftRepo.SetCollection(db.GetConnection().GetCollection("Shift"))
	}
	return shiftRepo
}
