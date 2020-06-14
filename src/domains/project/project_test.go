package project_test

import (
	"context"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"williamfeng323/mooncake-duty/src/domains/project"
	repoimpl "williamfeng323/mooncake-duty/src/infrastructure/db/repo_impl"
)

func TestProject(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Project domain Suite")
}

var _ = Describe("Project", func() {
	Describe("#NewProject", func() {
		It("should return a project instance", func() {
			project1 := project.NewProject("Test", "This is a test project", project.Member{MemberID: primitive.NewObjectID(), IsAdmin: true})
			Expect(project1).NotTo(BeNil())
			project2 := project.NewProject("Test", "This is a test project")
			Expect(project2).NotTo(BeNil())
		})
	})
	Describe("#Create", func() {
		prj := project.NewProject("TestProjectDomain", "This is a test project")
		Context("call by an empty domain instance", func() {
			It("should return error", func() {
				prjEmpty := &project.Project{}
				err := prjEmpty.Create()
				Expect(err.Error()).To(Equal("Project does not initialized"))
			})
		})
		Context("Call by valid domain instance", func() {
			AfterEach(func() {
				repo := repoimpl.GetProjectRepo()
				repo.DeleteOne(context.Background(), bson.M{"name": prj.Name})
			})
			Context("No existing project name in the database", func() {
				It("Should be able to insert into the database", func() {
					err := prj.Create()
					Expect(err).To(BeNil())
				})
			})
			Context("Existing project name in the database", func() {
				It("Should return error to state cannot override existing project", func() {
					prj.Create()
					prj2 := project.NewProject("Test", "This is the 2nd test project")
					err := prj2.Create()
					Expect(err.Error()).To(Equal("Project already exist"))
				})
			})
		})
	})
})
