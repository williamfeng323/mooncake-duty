package web_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	account "williamfeng323/mooncake-duty/src/domains/account"
	"williamfeng323/mooncake-duty/src/domains/project"
	repoimpl "williamfeng323/mooncake-duty/src/infrastructure/db/repo_impl"
	"williamfeng323/mooncake-duty/src/infrastructure/middlewares"
	webInterface "williamfeng323/mooncake-duty/src/interfaces/web"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestProjectRouter(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Project Router Suite")
}

var _ = Describe("ProjectRouter", func() {
	var router *gin.Engine
	var testRecorder *httptest.ResponseRecorder

	testProjectName := "TestBigProject"
	testProjectDescription := "Test Big Project!!"
	BeforeEach(func() {
		router = gin.Default()
		router.Use(middlewares.Logger())
		webInterface.RegisterProjectRoute(router)
		testRecorder = httptest.NewRecorder()
	})
	AfterEach(func() {
		repoimpl.GetProjectRepo().DeleteOne(context.Background(), bson.M{"name": testProjectName})
	})
	Describe("GET /", func() {
		BeforeEach(func() {
			req, _ := http.NewRequest("PUT", "/projects",
				bytes.NewReader([]byte(fmt.Sprintf(`{"name": "%s", "description": "%s"}`, testProjectName, testProjectDescription))))
			router.ServeHTTP(httptest.NewRecorder(), req)
		})
		It("should return 200 if project name provided properly", func() {
			req, _ := http.NewRequest("GET", fmt.Sprintf(`/projects?name=%s`, testProjectName), nil)
			router.ServeHTTP(testRecorder, req)
			Expect(testRecorder.Code).To(Equal(200))
		})
		It("should return 404 if project name cannot be found", func() {
			req, _ := http.NewRequest("GET", `/projects?name=NotEvenAProject`, nil)
			router.ServeHTTP(testRecorder, req)
			Expect(testRecorder.Code).To(Equal(404))
		})
	})
	Describe("PUT /", func() {
		It("should return 400 if no name and description provided", func() {
			req, _ := http.NewRequest("PUT", "/projects", bytes.NewReader([]byte(`{}`)))
			router.ServeHTTP(testRecorder, req)
			Expect(testRecorder.Code).To(Equal(400))
		})
		It("should return 400 if provided member but member does not exist", func() {
			req, _ := http.NewRequest("PUT", "/projects", bytes.NewReader([]byte(
				fmt.Sprintf(`{"name": "%s", "description": "%s", "members": [{"memberId": "123456789012345678901234"}]}`, testProjectName, testProjectDescription))))
			router.ServeHTTP(testRecorder, req)
			rsp := make(map[string]interface{})
			json.Unmarshal(testRecorder.Body.Bytes(), &rsp)
			Expect(testRecorder.Code).To(Equal(400))
			Expect(rsp["error"]).To(Equal("Invalid member account"))
		})
		It("should return 200 and create project if name and description provided", func() {
			req, _ := http.NewRequest("PUT", "/projects",
				bytes.NewReader([]byte(fmt.Sprintf(`{"name": "%s", "description": "%s"}`, testProjectName, testProjectDescription))))
			router.ServeHTTP(testRecorder, req)
			Expect(testRecorder.Code).To(Equal(200))
			rsp := make(map[string]project.Project)
			json.Unmarshal(testRecorder.Body.Bytes(), &rsp)
			Expect(rsp["project"].Name).To(Equal(testProjectName))
			Expect(rsp["project"].Description).To(Equal(testProjectDescription))
		})
	})
	Describe("GET /:id", func() {
		prj := project.NewProject(testProjectName, testProjectDescription)

		BeforeEach(func() {
			prj.Create()
		})
		AfterEach(func() {
			repoimpl.GetProjectRepo().DeleteOne(context.Background(), bson.M{"_id": prj.ID})
		})
		It("Should return 200 and the corresponding project details if the id is correct", func() {
			req, _ := http.NewRequest("GET", fmt.Sprintf(`/projects/%s`, prj.ID.Hex()), nil)
			router.ServeHTTP(testRecorder, req)
			Expect(testRecorder.Code).To(Equal(200))
			rst := map[string]project.Project{}
			json.Unmarshal(testRecorder.Body.Bytes(), &rst)
			Expect(rst["project"].Name).To(Equal(testProjectName))
			Expect(rst["project"].Description).To(Equal(testProjectDescription))
		})
		It("Should return 404 if the id is invalid", func() {
			req, _ := http.NewRequest("GET", `/projects/123456789012345678901234`, nil)
			router.ServeHTTP(testRecorder, req)
			Expect(testRecorder.Code).To(Equal(404))
		})
		It("Should return 400 if the id is not objectID format", func() {
			req, _ := http.NewRequest("GET", `/projects/12345678`, nil)
			router.ServeHTTP(testRecorder, req)
			Expect(testRecorder.Code).To(Equal(400))
		})
	})
	Describe("POST /:id", func() {
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
			repoimpl.GetAccountRepo().DeleteOne(context.Background(), bson.M{"_id": acct1.ID})
			repoimpl.GetAccountRepo().DeleteOne(context.Background(), bson.M{"_id": acct2.ID})
			repoimpl.GetProjectRepo().DeleteOne(context.Background(), bson.M{"_id": prj.ID})
		})
		It("Should return 404 if project id incorrect", func() {
			req, _ := http.NewRequest("POST", fmt.Sprintf(`/projects/%s`, primitive.NewObjectID().Hex()), bytes.NewReader([]byte(`{"name":"test"}`)))
			router.ServeHTTP(testRecorder, req)
			Expect(testRecorder.Code).To(Equal(404))
		})
		It("Should return 400 if no data in body", func() {
			req, _ := http.NewRequest("POST", fmt.Sprintf(`/projects/%s`, prj.ID.Hex()), bytes.NewReader([]byte(`{}`)))
			router.ServeHTTP(testRecorder, req)
			Expect(testRecorder.Code).To(Equal(400))
		})
		It("Should return 200 if name or description updated", func() {
			req, _ := http.NewRequest("POST", fmt.Sprintf(`/projects/%s`, prj.ID.Hex()), bytes.NewReader([]byte(`{"name": "newTestName"}`)))
			router.ServeHTTP(testRecorder, req)
			Expect(testRecorder.Code).To(Equal(200))
		})
		It("Should return 200 if member updated", func() {
			req, _ := http.NewRequest("POST", fmt.Sprintf(`/projects/%s`, prj.ID.Hex()),
				bytes.NewReader([]byte(fmt.Sprintf(`{"members": [{"memberId": "%s", "isAdmin":true}]}`, acct1.ID.Hex()))))
			router.ServeHTTP(testRecorder, req)
			Expect(testRecorder.Code).To(Equal(200))
		})
	})
})
