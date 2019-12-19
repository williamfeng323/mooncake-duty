package account

import (
	"testing"

	"williamfeng323/mooncake-duty/src/dao"
	"williamfeng323/mooncake-duty/src/utils"

	. "github.com/smartystreets/goconvey/convey"
)

func TestInsertAccount(t *testing.T) {
	Convey("Given an registered account model", t, func() {
		db := &dao.Connection{}
		db.InitConnection(nil, utils.GetConf().Mongo)
		db.Register(&Account{})
		Convey("insert account with empty data", func() {
			acct := &Account{}
			db.CollectionRegistry["Account"].New(acct)
			rst, err := acct.InsertAccount()
			So(err, ShouldNotBeNil)
			So(rst, ShouldBeNil)
		})
		Convey("insert account with valid data", func() {
			acct := &Account{}
			db.CollectionRegistry["Account"].New(acct)
			acct.Email = "test@abc.com"
			acct.Password = "test"
			rst, err := acct.InsertAccount()
			So(err, ShouldBeNil)
			So(rst, ShouldNotBeNil)
			acct.DeleteByID(acct.ID)
		})
	})
}
