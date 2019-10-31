package utils

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetConf(t *testing.T) {
	Convey("When config file provided", t, func() {
		Convey("should load data correctly", func() {
			rst := GetConf()
			So(rst, ShouldNotBeNil)
			So(rst.Mongo.URL, ShouldNotBeEmpty)
		})
	})
}
