package dao

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

//IDocumentBase Interface which each collection document (model) has to implement
type IDocumentBase interface {
	SetCollection(*mongo.Collection)
	SetDocument(IDocumentBase)
	GetCollection() *mongo.Collection
	FindByID(primitive.ObjectID) *mongo.SingleResult
	Find(...interface{}) (*mongo.Cursor, error)
}

// BaseModel is the base model that other models should embedded.
// Please use tag `json:",inline" bson:",inline"` to make the exported
// fields inline.
type BaseModel struct {
	document   IDocumentBase
	collection *mongo.Collection
	typeValue  *reflect.Value
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	CreatedAt  time.Time          `json:"createdAt" bson:"createdAt,omitempty"`
	UpdatedAt  time.Time          `json:"updatedAt" bson:"updatedAt,omitempty"`
	Deleted    bool               `json:"-" bson:"deleted,omitempty"`
}

// SetCollection will set the document's own collection instance
func (model *BaseModel) SetCollection(coll *mongo.Collection) {
	model.collection = coll
}

// SetDocument will set the document's itself into the IDocumentBase.
// It guaranteens the generic function DefaultValidator can get the
// run time data of the document.
func (model *BaseModel) SetDocument(doc IDocumentBase) {
	model.document = doc
}

// GetCollection will set the document's own collection instance
func (model *BaseModel) GetCollection() *mongo.Collection {
	return model.collection
}

//FindByID find document by document _id
func (model *BaseModel) FindByID(id primitive.ObjectID) *mongo.SingleResult {
	filter := bson.M{"_id": id}
	ctx, cancelFunc := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelFunc()
	return model.collection.FindOne(ctx, filter)
}

//Find return the result cursor base on your query
func (model *BaseModel) Find(query ...interface{}) (*mongo.Cursor, error) {
	finalQuery := initQuery(query)
	ctx, cancelFunc := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelFunc()
	//accept zero or one query param
	cursor, err := model.collection.Find(ctx, finalQuery)
	if err != nil {
		return nil, err
	}
	return cursor, nil
}

//DeleteByID find document by document _id
func (model *BaseModel) DeleteByID(id primitive.ObjectID) (*mongo.DeleteResult, error) {
	filter := bson.M{"_id": id}
	ctx, cancelFunc := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelFunc()
	return model.collection.DeleteOne(ctx, filter)
}

//Delete deletes the documents according to the query provided.
func (model *BaseModel) Delete(query ...interface{}) (*mongo.DeleteResult, error) {
	finalQuery := initQuery(query)
	ctx, cancelFunc := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelFunc()
	rst, err := model.collection.DeleteOne(ctx, finalQuery)
	return rst, err
}

//DefaultValidator will validate below tag:
//    required: boolean <-- The field could not be empty.
//    TBD: More default validators
func (model *BaseModel) DefaultValidator() []error {

	docElem := reflect.ValueOf(model.document).Elem()
	fieldType := docElem.Type()
	validationErrors := []error{}

	for fieldIndex := 0; fieldIndex < docElem.NumField(); fieldIndex++ {
		var required bool
		var err error
		var fieldValue reflect.Value
		field := fieldType.Field(fieldIndex)
		fieldTag := field.Tag
		requiredTag := fieldTag.Get("required")
		fieldElem := docElem.Field(fieldIndex)
		fieldName := field.Name
		if fieldElem.Kind() == reflect.Ptr || fieldElem.Kind() == reflect.Interface {
			fieldValue = fieldElem.Elem()
		} else {
			fieldValue = fieldElem
		}

		if len(requiredTag) > 0 {
			required, err = strconv.ParseBool(requiredTag)
			if err != nil {
				panic("Check your required tag - must be boolean")
			}
			if required {
				if err := validateRequired(fieldValue, fieldName); err != nil {
					validationErrors = append(validationErrors, err)
				}
			}
		}
		if len(validationErrors) > 0 {
			return validationErrors
		}
	}
	return nil
}

func validateRequired(fieldValue reflect.Value, fieldName string) error {
	isSet := false
	if !fieldValue.IsValid() {
		isSet = false
	} else if fieldValue.Kind() == reflect.Ptr && !fieldValue.IsNil() {
		isSet = true
	} else if fieldValue.Kind() == reflect.Slice || fieldValue.Kind() == reflect.Map {
		isSet = fieldValue.Len() > 0
	} else if fieldValue.Kind() == reflect.Interface {
		isSet = fieldValue.Interface() != nil
	} else {
		va := fieldValue.Interface()
		zeroValue := reflect.Zero(reflect.TypeOf(fieldValue.Interface())).Interface()
		isSet = !reflect.DeepEqual(va, zeroValue)
	}

	if !isSet {
		return fmt.Errorf("Field - %s cannot be null", fieldName)
	}
	return nil
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
