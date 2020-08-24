package web_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	account "williamfeng323/mooncake-duty/src/domains/account"
	"williamfeng323/mooncake-duty/src/domains/issue"
	project "williamfeng323/mooncake-duty/src/domains/project"
	"williamfeng323/mooncake-duty/src/domains/shift"
	repoimpl "williamfeng323/mooncake-duty/src/infrastructure/db/repo_impl"
	"williamfeng323/mooncake-duty/src/infrastructure/middlewares"
	webInterface "williamfeng323/mooncake-duty/src/interfaces/web"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestIssueRouter(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Issue Router Suite")
}

var _ = Describe("IssueRouter", func() {
	var router *gin.Engine
	var testRecorder *httptest.ResponseRecorder

	testProjectName := "TestIssueBigProject"
	testProjectDescription := "Test Issue Big Project!!"

	var testShift *shift.Shift
	acct1, _ := account.NewAccount("Test1@test.com", "Testaccount1")
	acct2, _ := account.NewAccount("Test2@test.com", "Testaccount2")
	acct3, _ := account.NewAccount("Test3@test.com", "Testaccount3")
	acct4, _ := account.NewAccount("Test4@test.com", "Testaccount4")
	acct5, _ := account.NewAccount("Test5@test.com", "Testaccount5")
	acct6, _ := account.NewAccount("Test6@test.com", "Testaccount6")
	acct7, _ := account.NewAccount("Test7@test.com", "Testaccount7")
	acctSet := []*account.Account{acct1, acct2, acct3, acct4, acct5, acct6, acct7}

	prj := project.NewProject(testProjectName, testProjectDescription)
	prj.AlertInterval = 10 * time.Minute
	prj.CallsPerTier = 2

	project.GetProjectService().SetMembers(prj.ID,
		project.Member{MemberID: acct1.ID, IsAdmin: false}, project.Member{MemberID: acct2.ID, IsAdmin: false}, project.Member{MemberID: acct3.ID, IsAdmin: false},
		project.Member{MemberID: acct4.ID, IsAdmin: false}, project.Member{MemberID: acct5.ID, IsAdmin: false}, project.Member{MemberID: acct6.ID, IsAdmin: false},
		project.Member{MemberID: acct7.ID, IsAdmin: true})

	BeforeEach(func() {
		router = gin.Default()
		router.Use(middlewares.Logger())
		webInterface.RegisterIssueRoute(router)
		testRecorder = httptest.NewRecorder()

		for _, v := range acctSet {
			v.Save(false)
		}
		prj.Create()
		testShift, _ = shift.NewShift(prj.ID, shift.Mon, time.Date(2020, 5, 16, 0, 0, 0, 0, time.Now().Location()), shift.Weekly)
		testShift.T1Members = []primitive.ObjectID{acct1.ID, acct2.ID, acct3.ID}
		testShift.T2Members = []primitive.ObjectID{acct4.ID, acct5.ID, acct6.ID}
		testShift.T3Members = []primitive.ObjectID{acct7.ID}
		testShift.Create()
	})
	AfterEach(func() {
		acctRepo := repoimpl.GetAccountRepo()
		for _, v := range acctSet {
			acctRepo.DeleteOne(context.Background(), bson.M{"_id": v.ID})
		}
		repoimpl.GetProjectRepo().DeleteOne(context.Background(), bson.M{"_id": prj.ID})
		repoimpl.GetShiftRepo().DeleteOne(context.Background(), bson.M{"_id": testShift.ID})
	})

	Describe("GET /", func() {
		var i1 *issue.Issue
		var i2 *issue.Issue

		BeforeEach(func() {
			i1, _ = issue.GetIssueService().CreateNewIssue(prj.ID, "testService1")
			i2, _ = issue.GetIssueService().CreateNewIssue(prj.ID, "testService2")
		})
		AfterEach(func() {
			repoimpl.GetIssueRepo().DeleteOne(context.Background(), bson.M{"_id": i1.ID})
			repoimpl.GetIssueRepo().DeleteOne(context.Background(), bson.M{"_id": i2.ID})
		})
		It("should get status 400 when no project specified", func() {
			req, _ := http.NewRequest("GET", "/issues?issueKey=whatever", nil)
			router.ServeHTTP(testRecorder, req)
			Expect(testRecorder.Code).To(Equal(400))
		})
		It("should return issue list if param correct", func() {
			req, _ := http.NewRequest("GET", fmt.Sprintf("/issues?projectId=%s", prj.ID.Hex()), nil)
			router.ServeHTTP(testRecorder, req)
			Expect(testRecorder.Code).To(Equal(200))
			rsp := make(map[string][]issue.Issue)
			json.Unmarshal(testRecorder.Body.Bytes(), &rsp)
			Expect(len(rsp["issues"])).To(Equal(2))
		})
	})
})
