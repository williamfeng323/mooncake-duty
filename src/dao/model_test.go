package dao

import (
	"fmt"
	"reflect"
	"testing"
	"time"
	"williamfeng323/mooncake-duty/src/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDefaultValidate(t *testing.T) {
	Convey("Giving a struct with tag required", t, func() {
		type Test struct {
			BaseModel
			TestWhat string `json:"testing" bson:"testing" required:"true"`
		}
		db := &Connection{}
		db.InitConnection(nil, utils.GetConf().Mongo)
		db.Register(&Test{})
		Convey("Should throw error if the document does not init the required field", func() {
			test := &Test{}
			db.CollectionRegistry["Test"].New(test)
			err := test.DefaultValidator()
			So(err, ShouldNotBeNil)
		})
		Convey("Should not throw error if the document initted the required field", func() {
			test := &Test{}
			db.CollectionRegistry["Test"].New(test)
			test.CreatedAt = time.Now()
			test.ID = primitive.NewObjectID()
			test.TestWhat = "What to test"
			testValue := reflect.ValueOf(test)
			ti := testValue.Interface()
			fmt.Printf("Value of test is %v", ti)
			err := test.DefaultValidator()
			So(err, ShouldBeNil)
		})
	})
}
