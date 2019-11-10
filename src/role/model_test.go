package role

import (
	"testing"

	"williamfeng323/mooncake-duty/src/dao"
	"williamfeng323/mooncake-duty/src/utils"

	. "github.com/smartystreets/goconvey/convey"
)

func TestInsertRole(t *testing.T) {
	Convey("Given an registered account model", t, func() {
		db := &dao.Connection{}
		db.InitConnection(nil, utils.GetConf().Mongo)
		db.Register(&Role{})
		Convey("insert account with empty data", func() {
			rol := &Role{}
			db.CollectionRegistry["Role"].New(rol)
			rst, err := rol.InsertRole()
			So(err, ShouldNotBeNil)
			So(rst, ShouldBeNil)
		})
		Convey("insert account with valid data", func() {
			rol := &Role{}
			db.CollectionRegistry["Role"].New(rol)
			rol.Name = "Admin"
			rol.Auth = Create
			rst, err := rol.InsertRole()
			So(err, ShouldBeNil)
			So(rst, ShouldNotBeNil)
			rol.DeleteByID(rol.ID)
		})
	})
}
