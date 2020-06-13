package repoimpl

import (
	"context"
	"sync"
	"williamfeng323/mooncake-duty/src/infrastructure/db"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//ProjectRepo is the implementation of Project repository
type ProjectRepo struct {
	collection *mongo.Collection
}

func init() {
	db.Register(&ProjectRepo{})
}

//GetName returns the name of project collection
func (pr *ProjectRepo) GetName() string {
	return "Project"
}

// SetCollection set the collection that communicate with db to the instance
func (pr *ProjectRepo) SetCollection(coll *mongo.Collection) {
	pr.collection = coll
}

// InsertOne creates new project
func (pr *ProjectRepo) InsertOne(ctx context.Context, document interface{},
	opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	//TO-DO Default validator
	return pr.collection.InsertOne(ctx, document, opts...)
}

//Find get one project document
func (pr *ProjectRepo) Find(ctx context.Context, filter interface{},
	opts ...*options.FindOptions) (*mongo.Cursor, error) {
	return pr.collection.Find(ctx, filter, opts...)
}

// FindOne returns the single document that meets the criteria
func (pr *ProjectRepo) FindOne(ctx context.Context, filter interface{},
	opts ...*options.FindOneOptions) *mongo.SingleResult {
	return pr.collection.FindOne(ctx, filter, opts...)
}

// UpdateOne update the document according to the filter.
func (pr *ProjectRepo) UpdateOne(ctx context.Context, filter interface{}, update interface{},
	opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return pr.collection.UpdateOne(ctx, filter, update, opts...)
}

// DeleteOne executes a delete command to delete at most one document from the collection.
func (pr *ProjectRepo) DeleteOne(ctx context.Context, filter interface{},
	opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return pr.collection.DeleteOne(ctx, filter, opts...)
}

var projectRepo *ProjectRepo
var projectLock sync.RWMutex

// GetProjectRepo get the account repository instance,
// Create if not exist.
func GetProjectRepo() *ProjectRepo {
	projectLock.Lock()
	defer projectLock.Unlock()
	if projectRepo == nil {
		projectRepo = &ProjectRepo{}
		projectRepo.SetCollection(db.GetConnection().GetCollection("Project"))
	}
	return projectRepo
}
