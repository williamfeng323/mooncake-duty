package issue_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestIssueService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Issue Service Suite")
}

var _ = Describe("Issue Service", func() {
	// var testShift *shift.Shift
	// acct1, _ := account.NewAccount("Test1@test.com", "Testaccount1")
	// acct2, _ := account.NewAccount("Test2@test.com", "Testaccount2")
	// acct3, _ := account.NewAccount("Test3@test.com", "Testaccount3")
	// acct4, _ := account.NewAccount("Test4@test.com", "Testaccount4")
	// acct5, _ := account.NewAccount("Test5@test.com", "Testaccount5")
	// acct6, _ := account.NewAccount("Test6@test.com", "Testaccount6")
	// acct7, _ := account.NewAccount("Test7@test.com", "Testaccount7")
	// acctSet := []*account.Account{acct1, acct2, acct3, acct4, acct5, acct6, acct7}
	// for _, v := range acctSet {
	// 	v.Save(false)
	// }
	// prj := project.NewProject("IssueServiceTest", "This is a test project")
	// prj.AlertInterval = 10 * time.Minute
	// prj.CallsPerTier = 2
	// prj.Create()
	// project.GetProjectService().SetMembers(prj.Name,
	// 	project.Member{MemberID: acct1.ID, IsAdmin: false}, project.Member{MemberID: acct2.ID, IsAdmin: false}, project.Member{MemberID: acct3.ID, IsAdmin: false},
	// 	project.Member{MemberID: acct4.ID, IsAdmin: false}, project.Member{MemberID: acct5.ID, IsAdmin: false}, project.Member{MemberID: acct6.ID, IsAdmin: false},
	// 	project.Member{MemberID: acct7.ID, IsAdmin: true})
	// testShift, _ = shift.NewShift(prj.ID, shift.Mon, time.Date(2020, 5, 16, 0, 0, 0, 0, time.Now().Location()), shift.Weekly)
	// testShift.T1Members = []primitive.ObjectID{acct1.ID, acct2.ID, acct3.ID}
	// testShift.T2Members = []primitive.ObjectID{acct4.ID, acct5.ID, acct6.ID}
	// testShift.T3Members = []primitive.ObjectID{acct7.ID}

	// testShift.Create()

	// AfterSuite(func() {
	// 	acctRepo := repoimpl.GetAccountRepo()
	// 	for _, v := range acctSet {
	// 		acctRepo.DeleteOne(context.Background(), bson.M{"_id": v.ID})
	// 	}
	// 	repoimpl.GetProjectRepo().DeleteOne(context.Background(), bson.M{"_id": prj.ID})
	// 	repoimpl.GetShiftRepo().DeleteOne(context.Background(), bson.M{"_id": testShift.ID})
	// })
	// prj.Create()
	// Describe("CreateIssue", func() {

	// })
})
