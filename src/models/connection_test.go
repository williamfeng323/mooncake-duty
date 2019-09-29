package models_test

import (
	"testing"

	"williamfeng323/mooncake-duty/src/models"
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
			conn, err := models.Connect(nil, conf)
			Convey("The client should be initialized", func() {
				So(err, ShouldBeNil)
				So(conn, ShouldNotBeNil)
			})
		})
		Convey("When trigger Connect without conf", func() {
			conf := utils.Config{}
			conn, err := models.Connect(nil, conf.Mongo)
			Convey("The error should be return", func() {
				So(err, ShouldNotBeNil)
				So(conn, ShouldBeNil)
			})
		})
	})
}
