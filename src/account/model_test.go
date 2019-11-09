package account

import (
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"williamfeng323/mooncake-duty/src/dao"
	"williamfeng323/mooncake-duty/src/utils"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFindByRoleID(t *testing.T) {
	Convey("Given an registered account model", t, func() {
		db := &dao.Connection{}
		db.InitConnection(nil, utils.GetConf().Mongo)
		db.Register(&Account{})
		Convey("trigger FindByRoleID with ramdom data", func() {
			acct := &Account{}
			db.CollectionRegistry["Account"].New(acct)
			testObjectID := primitive.NewObjectID()
			rst, err := acct.FindByRoleID(testObjectID)
			So(len(rst), ShouldEqual, 0)
			So(err, ShouldBeNil)
		})
	})
}

func TestInsertAccount(t *testing.T) {
	Convey("Given an registered account model", t, func() {
		db := &dao.Connection{}
		db.InitConnection(nil, utils.GetConf().Mongo)
		db.Register(&Account{})
		Convey("trigger FindByRoleID with ramdom data", func() {
			acct := &Account{}
			db.CollectionRegistry["Account"].New(acct)
			rst, err := acct.InsertAccount()
			So(rst, ShouldNotBeNil)
			So(err, ShouldBeNil)
		})
	})
}
