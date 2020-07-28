package project_test

import (
	"context"
	"testing"
	"williamfeng323/mooncake-duty/src/domains/account"
	"williamfeng323/mooncake-duty/src/domains/project"
	repoimpl "williamfeng323/mooncake-duty/src/infrastructure/db/repo_impl"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestProjectService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Project service Suite")
}

var _ = Describe("project.Service", func() {
	Describe("#SetMembers", func() {
		prj := project.NewProject("ProjectTestService", "This is a test project", project.Member{MemberID: primitive.NewObjectID(), IsAdmin: true})

		prjService := &project.Service{}
		prjService.SetRepo(repoimpl.GetProjectRepo())

		acct1, _ := account.NewAccount("Test1@test.com", "Testaccount1")
		acct2, _ := account.NewAccount("Test2@test.com", "Testaccount1")
		BeforeEach(func() {
			prj.Create()
			acct1.Save(false)
			acct2.Save(false)
		})
		AfterEach(func() {
			acctRepo := repoimpl.GetAccountRepo()
			acctRepo.DeleteOne(context.Background(), bson.M{"_id": acct1.ID})
			acctRepo.DeleteOne(context.Background(), bson.M{"_id": acct2.ID})
			repo := repoimpl.GetProjectRepo()
			repo.DeleteOne(context.Background(), bson.M{"name": prj.Name})
		})

		Context("Provides a valid project name", func() {
			Context("Provides an invalid member id", func() {
				It("should return account not found error and empty succeeded member id", func() {
					succeededIDs, failedIDs, err := prjService.SetMembers(prj.ID, project.Member{MemberID: primitive.NewObjectID(), IsAdmin: false})
					Expect(succeededIDs).To(BeNil())
					Expect(failedIDs).ToNot(Equal([]project.Member{}))
					Expect(err.Error()).To(Equal("Account Not Found"))
				})
				It("should return the account not found error and return succeeded member id if only few members do not exist", func() {
					succeededIDs, failedIDs, err := prjService.SetMembers(prj.ID, project.Member{MemberID: primitive.NewObjectID(), IsAdmin: false}, project.Member{MemberID: acct1.ID, IsAdmin: true})
					Expect(err).To(BeNil())
					Expect(succeededIDs[0].MemberID).To(Equal(acct1.ID))
					Expect(failedIDs).ToNot(Equal([]project.Member{}))
				})
			})
			Context("Provides valid member id", func() {
				It("should return the succeeded member id and nil error if only provides 1 member id", func() {
					succeededIDs, failedIDs, err := prjService.SetMembers(prj.ID, project.Member{MemberID: acct1.ID, IsAdmin: true})
					Expect(err).To(BeNil())
					Expect(succeededIDs[0].MemberID).To(Equal(acct1.ID))
					Expect(failedIDs).To(Equal([]project.Member{}))
				})
				It("should be able to return succeeded member ids and nil error if provides multiple member id ", func() {
					succeededIDs, failedIDs, err := prjService.SetMembers(prj.ID, project.Member{MemberID: acct2.ID, IsAdmin: false}, project.Member{MemberID: acct1.ID, IsAdmin: true})
					Expect(err).To(BeNil())
					Expect(succeededIDs[0].MemberID).To(Equal(acct2.ID))
					Expect(succeededIDs[1].MemberID).To(Equal(acct1.ID))
					Expect(failedIDs).To(Equal([]project.Member{}))
				})
			})
		})
		Context("Provides an invalid project name", func() {
			It("Should return project not found error", func() {
				succeededIDs, failedIDs, err := prjService.SetMembers(primitive.NewObjectID(), project.Member{MemberID: acct1.ID, IsAdmin: true})
				Expect(succeededIDs).To(BeNil())
				Expect(failedIDs).ToNot(Equal([]project.Member{}))
				Expect(err).To(MatchError(project.NotFoundError{}))
			})
		})
	})
	Describe("#SetNameOrDescription", func() {
		prj := project.NewProject("ProjectTestService", "This is a test project", project.Member{MemberID: primitive.NewObjectID(), IsAdmin: true})

		prjService := &project.Service{}
		prjService.SetRepo(repoimpl.GetProjectRepo())

		BeforeEach(func() {
			prj.Create()
		})
		AfterEach(func() {
			repoimpl.GetProjectRepo().DeleteOne(context.Background(), bson.M{"name": prj.Name})
		})
		It("should update the name and description if any of it provided", func() {
			prjService.SetNameOrDescription(prj.ID, "newName", "new description")
			p := repoimpl.GetProjectRepo().FindOne(context.Background(), bson.M{"_id": prj.ID})
			rstPrj := project.Project{}
			p.Decode(&rstPrj)
			Expect(rstPrj.Name).To(Equal("newName"))
			Expect(rstPrj.Description).To(Equal("new description"))
		})
		It("should return project not found error if id invalid", func() {
			err := prjService.SetNameOrDescription(primitive.NewObjectID(), "newName", "")
			Expect(err).To(MatchError(project.NotFoundError{}))
		})
	})
	Describe("#CreateProject", func() {
		acct1, _ := account.NewAccount("Test1@test.com", "Testaccount1")
		acct2, _ := account.NewAccount("Test2@test.com", "Testaccount1")
		BeforeEach(func() {
			acct1.Save(false)
			acct2.Save(false)
		})
		AfterEach(func() {
			repoimpl.GetAccountRepo().DeleteOne(context.Background(), bson.M{"_id": acct1.ID})
			repoimpl.GetAccountRepo().DeleteOne(context.Background(), bson.M{"_id": acct2.ID})
		})
		It("Should return account not found when provided member that is not exist", func() {
			proj, err := project.GetProjectService().CreateProject("Test", "Test project we crated", project.Member{MemberID: primitive.NewObjectID(), IsAdmin: false})
			Expect(proj).To(BeNil())
			Expect(err).To(MatchError(account.NotFoundError{}))
		})
		It("should return account when all fields are correct", func() {
			projMember1 := project.Member{MemberID: acct1.ID, IsAdmin: false}
			projMember2 := project.Member{MemberID: acct2.ID, IsAdmin: false}
			proj, err := project.GetProjectService().CreateProject("TestService", "Test Project we created", projMember1, projMember2)
			Expect(proj).ToNot(BeNil())
			Expect(err).To(BeNil())
			repoimpl.GetProjectRepo().DeleteOne(context.Background(), bson.M{"_id": proj.ID})
		})
	})
	Describe("#GetProjectByName", func() {
		prj := project.NewProject("ProjectTestService", "This is a test project")

		prjService := &project.Service{}
		prjService.SetRepo(repoimpl.GetProjectRepo())

		BeforeEach(func() {
			prj.Create()
		})
		AfterEach(func() {
			repoimpl.GetProjectRepo().DeleteOne(context.Background(), bson.M{"name": prj.Name})
		})
		It("should return project not found error if no project found", func() {
			prj, err := prjService.GetProjectByName("ProjectWhat")
			Expect(prj).To(BeNil())
			Expect(err).To(MatchError(project.NotFoundError{}))
		})
		It("should return project if project found", func() {
			prj, err := prjService.GetProjectByName("ProjectTestService")
			Expect(prj).ToNot(BeNil())
			Expect(err).To(BeNil())
		})
	})
	Describe("#GetProjectByID", func() {
		prj := project.NewProject("ProjectTestService", "This is a test project")

		prjService := &project.Service{}
		prjService.SetRepo(repoimpl.GetProjectRepo())

		BeforeEach(func() {
			prj.Create()
		})
		AfterEach(func() {
			repoimpl.GetProjectRepo().DeleteOne(context.Background(), bson.M{"name": prj.Name})
		})
		It("should return project not found error if no project found", func() {
			prj, err := prjService.GetProjectByID(primitive.NewObjectID())
			Expect(prj).To(BeNil())
			Expect(err).To(MatchError(project.NotFoundError{}))
		})
		It("should return project if project found", func() {
			prj, err := prjService.GetProjectByID(prj.ID)
			Expect(prj).ToNot(BeNil())
			Expect(err).To(BeNil())
		})
	})
})
