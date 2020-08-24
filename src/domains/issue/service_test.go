package issue_test

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.mongodb.org/mongo-driver/bson"

	"williamfeng323/mooncake-duty/src/domains/issue"
	"williamfeng323/mooncake-duty/src/domains/project"
	repoimpl "williamfeng323/mooncake-duty/src/infrastructure/db/repo_impl"
)

func TestIssueService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Issue Service Suite")
}

var _ = Describe("Issue Service", func() {
	prj := project.NewProject("issueTestProject", "test project for testing issue")
	BeforeEach(func() {
		prj.Create()
	})
	AfterEach(func() {
		repoimpl.GetProjectRepo().DeleteOne(context.Background(), bson.M{"_id": prj.ID})
	})
	Describe("#CreateNewIssue", func() {
		It("Should create new issue in DB when project ID is correct", func() {
			print(prj.ID.Hex())
			i, e := issue.GetIssueService().CreateNewIssue(prj.ID, "testService")
			Expect(e).To(BeNil())
			Expect(i.ProjectID).To(Equal(prj.ID))
			Expect(i.IssueKey).To(Equal("testService"))
			repoimpl.GetIssueRepo().DeleteOne(context.Background(), bson.M{"_id": i.ID})
		})
		It("should return project not found error when project id is invalid", func() {
			i, e := issue.GetIssueService().CreateNewIssue(primitive.NewObjectID(), "testService")
			Expect(e).ToNot(BeNil())
			Expect(e).To(MatchError(project.NotFoundError{}))
			Expect(i).To(BeNil())
		})
	})
	Describe("#GetIssueLists", func() {
		prj2 := project.NewProject("issueTestProject2", "test project for testing issue")

		var i1 *issue.Issue
		var i2 *issue.Issue
		var i3 *issue.Issue
		var i4 *issue.Issue
		var i5 *issue.Issue
		var i6 *issue.Issue
		BeforeEach(func() {
			prj2.Create()
			i1, _ = issue.GetIssueService().CreateNewIssue(prj.ID, "mock1")
			i2, _ = issue.GetIssueService().CreateNewIssue(prj.ID, "mock1")
			i3, _ = issue.GetIssueService().CreateNewIssue(prj.ID, "mock3")
			i4, _ = issue.GetIssueService().CreateNewIssue(prj2.ID, "mock1")
			i5, _ = issue.GetIssueService().CreateNewIssue(prj2.ID, "mock2")
			i6, _ = issue.GetIssueService().CreateNewIssue(prj2.ID, "mock3")
		})
		AfterEach(func() {
			repoimpl.GetProjectRepo().DeleteOne(context.Background(), bson.M{"_id": prj2.ID})
			repoimpl.GetIssueRepo().DeleteOne(context.Background(), bson.M{"_id": i1.ID})
			repoimpl.GetIssueRepo().DeleteOne(context.Background(), bson.M{"_id": i2.ID})
			repoimpl.GetIssueRepo().DeleteOne(context.Background(), bson.M{"_id": i3.ID})
			repoimpl.GetIssueRepo().DeleteOne(context.Background(), bson.M{"_id": i4.ID})
			repoimpl.GetIssueRepo().DeleteOne(context.Background(), bson.M{"_id": i5.ID})
			repoimpl.GetIssueRepo().DeleteOne(context.Background(), bson.M{"_id": i6.ID})
		})
		It("Should return project not found error if project id does not exist", func() {
			issues, err := issue.GetIssueService().GetIssueLists(primitive.NewObjectID(), "", issue.Init)
			Expect(issues).To(BeNil())
			Expect(err).To(MatchError(project.NotFoundError{}))
		})
		It("Should return issue list with specific status, and issue key", func() {
			issues, err := issue.GetIssueService().GetIssueLists(prj.ID, "mock1", issue.Init)
			Expect(err).To(BeNil())
			Expect(len(issues)).To(Equal(2))

			i4.UpdateStatus(issue.Resolved, "test")
			issues, err = issue.GetIssueService().GetIssueLists(prj2.ID, "", issue.Init)
			Expect(err).To(BeNil())
			Expect(len(issues)).To(Equal(2))

			issues, err = issue.GetIssueService().GetIssueLists(prj2.ID, "", -1)
			Expect(err).To(BeNil())
			Expect(len(issues)).To(Equal(3))
		})
	})
	Describe("#GetIssueByID", func() {

		It("should return not found error when issue id doesn't exist", func() {
			i, e := issue.GetIssueService().GetIssueByID(primitive.NewObjectID())
			Expect(e).To(MatchError(issue.NotFoundError{}))
			Expect(i).To(BeNil())
		})
		It("should return the issue instance if issue exist", func() {
			i0, _ := issue.GetIssueService().CreateNewIssue(prj.ID, "testService")
			i, e := issue.GetIssueService().GetIssueByID(i0.ID)
			Expect(e).To(BeNil())
			Expect(i.ID).To(Equal(i0.ID))
			repoimpl.GetIssueRepo().DeleteOne(context.Background(), bson.M{"_id": i0.ID})
		})
	})
})
