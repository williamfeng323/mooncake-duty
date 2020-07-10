package issue_test

import (
	"context"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"williamfeng323/mooncake-duty/src/domains/account"
	"williamfeng323/mooncake-duty/src/domains/issue"
	"williamfeng323/mooncake-duty/src/domains/project"
	"williamfeng323/mooncake-duty/src/domains/shift"
	repoimpl "williamfeng323/mooncake-duty/src/infrastructure/db/repo_impl"
)

func TestIssue(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Issue Domain Suite")
}

var _ = Describe("Issue", func() {
	var testShift *shift.Shift
	acct1, _ := account.NewAccount("Test1@test.com", "Testaccount1")
	acct2, _ := account.NewAccount("Test2@test.com", "Testaccount2")
	acct3, _ := account.NewAccount("Test3@test.com", "Testaccount3")
	acct4, _ := account.NewAccount("Test4@test.com", "Testaccount4")
	acct5, _ := account.NewAccount("Test5@test.com", "Testaccount5")
	acct6, _ := account.NewAccount("Test6@test.com", "Testaccount6")
	acct7, _ := account.NewAccount("Test7@test.com", "Testaccount7")
	acctSet := []*account.Account{acct1, acct2, acct3, acct4, acct5, acct6, acct7}
	for _, v := range acctSet {
		v.Save(false)
	}
	prj := project.NewProject("IssueTest", "This is a test project")
	prj.AlertInterval = 10 * time.Minute
	prj.CallsPerTier = 2
	prj.Create()
	project.GetProjectService().SetMembers(prj.Name,
		project.Member{MemberID: acct1.ID, IsAdmin: false}, project.Member{MemberID: acct2.ID, IsAdmin: false}, project.Member{MemberID: acct3.ID, IsAdmin: false},
		project.Member{MemberID: acct4.ID, IsAdmin: false}, project.Member{MemberID: acct5.ID, IsAdmin: false}, project.Member{MemberID: acct6.ID, IsAdmin: false},
		project.Member{MemberID: acct7.ID, IsAdmin: true})
	testShift, _ = shift.NewShift(prj.ID, shift.Mon, time.Date(2020, 5, 16, 0, 0, 0, 0, time.Now().Location()), shift.Weekly)
	testShift.T1Members = []primitive.ObjectID{acct1.ID, acct2.ID, acct3.ID}
	testShift.T2Members = []primitive.ObjectID{acct4.ID, acct5.ID, acct6.ID}
	testShift.T3Members = []primitive.ObjectID{acct7.ID}

	testShift.Create()

	AfterSuite(func() {
		acctRepo := repoimpl.GetAccountRepo()
		for _, v := range acctSet {
			acctRepo.DeleteOne(context.Background(), bson.M{"_id": v.ID})
		}
		repoimpl.GetProjectRepo().DeleteOne(context.Background(), bson.M{"_id": prj.ID})
		repoimpl.GetShiftRepo().DeleteOne(context.Background(), bson.M{"_id": testShift.ID})
	})
	prj.Create()
	Describe("#GetNotificationTier", func() {
		It("should return project not found error and 0 when project doesn't exist", func() {
			iss := issue.Issue{ProjectID: primitive.NewObjectID()}
			rst, err := iss.GetNotificationTier()
			Expect(rst).To(Equal(issue.Tier(0)))
			Expect(err).To(MatchError(project.NotFoundError{}))
		})
		It("Should return the T1 members when T1 notification do not over number configure in project.CallsPerTier", func() {
			iss := issue.Issue{ProjectID: prj.ID, Service: "Test", CreatedAt: time.Now()}
			rst, err := iss.GetNotificationTier()
			Expect(rst).To(Equal(issue.T1))
			Expect(err).To(BeNil())
		})
		It("Should return the T2 members when T1 notification count over number configure in project.CallsPerTier", func() {
			iss := issue.Issue{ProjectID: prj.ID, Service: "Test", CreatedAt: time.Now(), T1NotificationCount: prj.CallsPerTier}
			rst, err := iss.GetNotificationTier()
			Expect(rst).To(Equal(issue.T2))
			Expect(err).To(BeNil())
		})
		It("Should return the T3 members when T2 notification count over number configure in project.CallsPerTier", func() {
			iss := issue.Issue{ProjectID: prj.ID, Service: "Test", CreatedAt: time.Now(), T1NotificationCount: prj.CallsPerTier, T2NotificationCount: prj.CallsPerTier}
			rst, err := iss.GetNotificationTier()
			Expect(rst).To(Equal(issue.T3))
			Expect(err).To(BeNil())
		})
	})
	Describe("#NewIssue", func() {
		Context("If project does not exist", func() {
			It("should return project not found error and nil issue", func() {
				iss, err := issue.NewIssue(primitive.NewObjectID(), "testService")
				Expect(iss).To(BeNil())
				Expect(err).To(MatchError(project.NotFoundError{}))
			})
		})
		Context("If the project exists", func() {
			It("should return issue instance and nil error", func() {
				iss, err := issue.NewIssue(prj.ID, "testService")
				Expect(iss).ToNot(BeNil())
				Expect(err).To(BeNil())
			})
		})
	})
})
