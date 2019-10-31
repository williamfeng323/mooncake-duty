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

func TestRegister(t *testing.T) {
	Convey("Giving a test struct", t, func() {
		type TestModel struct {
			BaseModel
			IsTest bool
		}
		db := &Connection{}
		db.InitConnection(nil, utils.GetConf().Mongo)
		Convey("when register to a Connection instance", func() {
			db.Register(&TestModel{})
			Convey("Should register the collection to connection", func() {
				So(db.CollectionRegistry["TestModel"], ShouldNotBeNil)
			})
		})
	})
}

func TestNewModel(t *testing.T) {
	// Only pass t into top-level Convey calls
	Convey("Given a Connection instance and a TestModel", t, func() {
		type TestModel struct {
			BaseModel
			IsTest bool
		}
		conf := utils.GetConf()
		db := &Connection{}
		db.InitConnection(nil, conf.Mongo)
		db.Register(&TestModel{})
		Convey("When trigger New", func() {
			testModel := &TestModel{}
			err := db.CollectionRegistry["TestModel"].New(testModel)
			Convey("should return a instance model instance with connection to collection", func() {
				So(err, ShouldBeNil)
				So(testModel, ShouldNotBeNil)
				So(testModel.collection, ShouldNotBeNil)
			})
		})
	})
}
