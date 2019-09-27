package models_test

import (
	"testing"

	"williamfeng323/mooncake-duty/src/models"

	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestConnect(t *testing.T) {
	// Only pass t into top-level Convey calls
	Convey("Given an empty connection and client options", t, func() {
		conn := &models.Connection{}
		conf := &options.ClientOptions{}
		Convey("When trigger Connect", func() {
			conn.Connect(conf)
			Convey("The client should be initialized", func() {
				So(conn.Client, ShouldNotBeNil)
			})
		})
		Convey("When client option is empty", func() {
			error := conn.Connect(conf)
			Convey("The error should be return", func() {
				So(error, ShouldNotBeNil)
			})
		})
	})
}
