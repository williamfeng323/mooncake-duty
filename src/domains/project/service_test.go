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
		prj := project.NewProject("Test", "This is a test project", project.Member{MemberID: primitive.NewObjectID(), IsAdmin: true})
		prj.Create()

		prjService := &project.Service{}
		prjService.SetRepo(repoimpl.GetProjectRepo())
		AfterSuite(func() {
			repo := repoimpl.GetProjectRepo()
			repo.DeleteOne(context.Background(), bson.M{"name": prj.Name})
		})

		acct1, _ := account.NewAccount("Test1@test.com", "Testaccount1")
		acct2, _ := account.NewAccount("Test2@test.com", "Testaccount1")
		BeforeEach(func() {
			acct1.Save(false)
			acct2.Save(false)
		})
		AfterEach(func() {
			acctRepo := repoimpl.GetAccountRepo()
			acctRepo.DeleteOne(context.Background(), bson.M{"_id": acct1.ID})
			acctRepo.DeleteOne(context.Background(), bson.M{"_id": acct2.ID})
		})

		Context("Provides a valid project name", func() {
			Context("Provides an invalid member id", func() {
				It("should return account not found error and empty succeeded member id", func() {
					succeededIDs, failedIDs, err := prjService.SetMembers(prj.Name, project.Member{MemberID: primitive.NewObjectID(), IsAdmin: false})
					Expect(succeededIDs).To(BeNil())
					Expect(failedIDs).ToNot(Equal([]project.Member{}))
					Expect(err.Error()).To(Equal("Account Not Found"))
				})
				It("should return the account not found error and return succeeded member id if only few members do not exist", func() {
					succeededIDs, failedIDs, err := prjService.SetMembers(prj.Name, project.Member{MemberID: primitive.NewObjectID(), IsAdmin: false}, project.Member{MemberID: acct1.ID, IsAdmin: true})
					Expect(err).To(BeNil())
					Expect(succeededIDs[0].MemberID).To(Equal(acct1.ID))
					Expect(failedIDs).ToNot(Equal([]project.Member{}))
				})
			})
			Context("Provides valid member id", func() {
				It("should return the succeeded member id and nil error if only provides 1 member id", func() {
					succeededIDs, failedIDs, err := prjService.SetMembers(prj.Name, project.Member{MemberID: acct1.ID, IsAdmin: true})
					Expect(succeededIDs[0].MemberID).To(Equal(acct1.ID))
					Expect(failedIDs).To(Equal([]project.Member{}))
					Expect(err).To(BeNil())
				})
				It("should be able to return succeeded member ids and nil error if provides multiple member id ", func() {
					succeededIDs, failedIDs, err := prjService.SetMembers(prj.Name, project.Member{MemberID: acct2.ID, IsAdmin: false}, project.Member{MemberID: acct1.ID, IsAdmin: true})
					Expect(succeededIDs[0].MemberID).To(Equal(acct2.ID))
					Expect(succeededIDs[1].MemberID).To(Equal(acct1.ID))
					Expect(failedIDs).To(Equal([]project.Member{}))
					Expect(err).To(BeNil())
				})
			})
		})
		Context("Provides an invalid project name", func() {
			It("Should return project not found error", func() {
				succeededIDs, failedIDs, err := prjService.SetMembers("FakeProject", project.Member{MemberID: acct1.ID, IsAdmin: true})
				Expect(succeededIDs).To(BeNil())
				Expect(failedIDs).ToNot(Equal([]project.Member{}))
				Expect(err).ToNot(BeNil())
			})
		})
	})
})
