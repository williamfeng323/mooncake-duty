package account

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetAccountsByRole(t *testing.T) {
	Convey("When call getAccountsByRole with a role", t, func() {
		Convey("Should return error When role does not exists", func() {
			role := Role{
				Name: "test role",
			}
			rst, err := getAccountsByRole(role.Name)
			So(rst, ShouldBeNil)
			So(err, ShouldBeError)
		})
		Convey("Should return accounts When role exists", func() {
			role := Role{
				Name: "developer",
			}
			rst, err := getAccountsByRole(role.Name)
			So(rst, ShouldNotBeNil)
			So(err, ShouldBeNil)
		})
	})
}
