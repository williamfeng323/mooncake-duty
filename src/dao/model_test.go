package dao

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewModel(t *testing.T) {
	// Only pass t into top-level Convey calls
	Convey("Given an sample model", t, func() {
		Convey("When trigger New", func() {
			// db := &Connection{}
			// conf := utils.MongoConfig{
			// 	URL:               "localhost",
			// 	Port:              "27017",
			// 	Username:          "root",
			// 	Password:          "example",
			// 	Database:          "mooncake",
			// 	ConnectionOptions: "authSource=admin",
			// }
			// err := db.InitConnection(nil, conf)
			// pingError := db.Client.Ping(nil, nil)
			Convey("should return a instance model instance with connection to collection", func() {
				// So(err, ShouldBeNil)
				// So(pingError, ShouldBeNil)
				// So(db.Client, ShouldNotBeNil)
			})
		})
	})
}
