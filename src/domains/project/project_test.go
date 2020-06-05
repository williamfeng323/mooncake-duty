package project

import (
	"fmt"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// func TestGetRolesByRole(t *testing.T) {
// 	Convey("Giving a role and a user binds with the role", t, func() {
// 		// db := &dao.Connection{}
// 		// db.InitConnection(nil, utils.GetConf().Mongo)
// 		// db.Register(&Role{})
// 		// db.Register(&Role{})
// 		Convey("Should return error When role does not exists", func() {
// 			// role := Role{
// 			// 	Name: "test role",
// 			// }
// 			// rst, err := getRolesByRole(role.Name)
// 			// So(rst, ShouldBeNil)
// 			// So(err, ShouldBeError)
// 		})
// 		Convey("Should return accounts When role exists", func() {
// 			// role := Role{
// 			// 	Name: "developer",
// 			// }
// 			// rst, err := getRolesByRole(role.Name)
// 			// So(rst, ShouldNotBeNil)
// 			// So(err, ShouldBeNil)
// 		})
// 	})
// }
func TestProject(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Project domain Suite")
}

var _ = Describe("Project Domain", func() {
	Describe("#createProject", func() {
		Context("giving empty project instance", func() {
			project := Project{}
			It("should return error", func() {
				rst, err := createProject(project)
				Expect(rst).To(BeNil())
				Expect(err).ToNot(BeEquivalentTo(BeAssignableToTypeOf(fmt.Errorf(""))))
			})
		})
		Context("giving project instance with proper data", func() {
			It("should return ", func() {})
		})
	})
})
