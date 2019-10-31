package dao

// import (
// 	"testing"
// 	"williamfeng323/mooncake-duty/src/utils"

// 	. "github.com/smartystreets/goconvey/convey"
// )

// func TestNewModel(t *testing.T) {
// 	// Only pass t into top-level Convey calls
// 	Convey("Given a Connection instance and a TestModel", t, func() {
// 		type TestModel struct {
// 			BaseModel
// 			IsTest bool
// 		}
// 		conf := utils.GetConf()
// 		db := &Connection{}
// 		db.InitConnection(nil, conf.Mongo)
// 		db.Register(&TestModel{})
// 		Convey("When trigger New", func() {
// 			// testModel := db.CollectionRegistry["TestModel"].New(&TestModel{})
// 			Convey("should return a instance model instance with connection to collection", func() {
// 				// So(err, ShouldBeNil)
// 				// So(pingError, ShouldBeNil)
// 				// So(db.Client, ShouldNotBeNil)
// 			})
// 		})
// 	})
// }
