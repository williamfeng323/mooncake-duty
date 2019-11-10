package account

import (
	"testing"
	"williamfeng323/mooncake-duty/src/dao"
	"williamfeng323/mooncake-duty/src/role"
	"williamfeng323/mooncake-duty/src/utils"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetAccountsByRole(t *testing.T) {
	Convey("Giving a role and a user binds with the role", t, func() {
		db := &dao.Connection{}
		db.InitConnection(nil, utils.GetConf().Mongo)
		db.Register(&Account{})
		db.Register(&role.Role{})
		rol := &role.Role{}
		db.CollectionRegistry["Role"].New(rol)
		rol.Name = "developer"

		acct := &Account{}
		db.CollectionRegistry["Account"].New(acct)
		acct.Email = "test@abc.com"
		acct.Password = "test"

		Convey("Should return error When role does not exists", func() {
			role := role.Role{
				Name: "test role",
			}
			rst, err := getAccountsByRole(role.Name)
			So(rst, ShouldBeNil)
			So(err, ShouldBeError)
		})
		Convey("Should return accounts When role exists", func() {
			role := role.Role{
				Name: "developer",
			}
			rst, err := getAccountsByRole(role.Name)
			So(rst, ShouldNotBeNil)
			So(err, ShouldBeNil)
		})
	})
}
