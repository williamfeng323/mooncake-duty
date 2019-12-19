package team

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetRolesByRole(t *testing.T) {
	Convey("Giving a role and a user binds with the role", t, func() {
		// db := &dao.Connection{}
		// db.InitConnection(nil, utils.GetConf().Mongo)
		// db.Register(&Role{})
		// db.Register(&Role{})
		Convey("Should return error When role does not exists", func() {
			// role := Role{
			// 	Name: "test role",
			// }
			// rst, err := getRolesByRole(role.Name)
			// So(rst, ShouldBeNil)
			// So(err, ShouldBeError)
		})
		Convey("Should return accounts When role exists", func() {
			// role := Role{
			// 	Name: "developer",
			// }
			// rst, err := getRolesByRole(role.Name)
			// So(rst, ShouldNotBeNil)
			// So(err, ShouldBeNil)
		})
	})
}
