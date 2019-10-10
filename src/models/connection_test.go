package models

import (
	"testing"

	"williamfeng323/mooncake-duty/src/utils"

	. "github.com/smartystreets/goconvey/convey"
)

func TestConnect(t *testing.T) {
	// Only pass t into top-level Convey calls
	Convey("Given an empty context", t, func() {
		Convey("When trigger Connect with conf", func() {
			conf := utils.MongoConfig{
				URL:      "localhost",
				Port:     "27017",
				Username: "root",
				Password: "example",
				Database: "mooncake",
			}
			conn, err := Connect(nil, conf)
			Convey("The client should be initialized", func() {
				So(err, ShouldBeNil)
				So(conn, ShouldNotBeNil)
			})
		})
		Convey("When trigger Connect without conf", func() {
			conf := utils.Config{}
			conn, err := Connect(nil, conf.Mongo)
			Convey("The error should be return", func() {
				So(err, ShouldNotBeNil)
				So(conn, ShouldBeNil)
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
