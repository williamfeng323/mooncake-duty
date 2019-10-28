package dao

import (
	"testing"

	"williamfeng323/mooncake-duty/src/utils"

	. "github.com/smartystreets/goconvey/convey"
)

func TestConnect(t *testing.T) {
	// Only pass t into top-level Convey calls
	Convey("Given an empty context", t, func() {
		Convey("When trigger Connect with conf", func() {
			db := &Connection{}
			conf := utils.MongoConfig{
				URL:               "localhost",
				Port:              "27017",
				Username:          "root",
				Password:          "example",
				Database:          "mooncake",
				ConnectionOptions: "authSource=admin",
			}
			err := db.InitConnection(nil, conf)
			pingError := db.Client.Ping(nil, nil)
			Convey("The client should be initialized and no error when ping", func() {
				So(err, ShouldBeNil)
				So(pingError, ShouldBeNil)
				So(db.Client, ShouldNotBeNil)
			})
		})
		Convey("When trigger Connect without conf", func() {
			db := &Connection{}
			conf := utils.Config{}
			err := db.InitConnection(nil, conf.Mongo)
			Convey("The error should be return", func() {
				So(err, ShouldNotBeNil)
				So(db.Client, ShouldBeNil)
			})
		})
	})
}

// func TestGetCollection(t *testing.T) {
// 	// Only pass t into top-level Convey calls
// 	Convey("Given an test struct", t, func() {
// 		type TestStruct struct{}
// 		Convey("When get the model collection", func() {
// 			test := &TestStruct{}
// 			collection := getCollection(test)
// 			Convey("The client should be initialized", func() {
// 				So(collection, ShouldEqual, "testStruct")
// 			})
// 		})
// 	})
// }
