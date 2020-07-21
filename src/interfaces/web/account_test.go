package web_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	account "williamfeng323/mooncake-duty/src/domains/account"
	repoimpl "williamfeng323/mooncake-duty/src/infrastructure/db/repo_impl"
	"williamfeng323/mooncake-duty/src/infrastructure/middlewares"
	webInterface "williamfeng323/mooncake-duty/src/interfaces/web"
)

func TestAccountRouter(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Account Router Suite")
}

var _ = Describe("AccountRouter", func() {
	var router *gin.Engine
	var testRecorder *httptest.ResponseRecorder
	var acct *account.Account
	testEmail := "feng.z.2@test.com"
	BeforeEach(func() {
		router = gin.Default()
		router.Use(middlewares.Logger())
		webInterface.RegisterAccountRoute(router)
		testRecorder = httptest.NewRecorder()
		_, acct, _ = account.GetAccountService().Register(testEmail, "123", false)
	})
	AfterEach(func() {
		acctRepo := repoimpl.GetAccountRepo()
		acctRepo.DeleteOne(context.Background(), bson.M{"email": testEmail})
	})
	Describe("GET /accounts", func() {
		It("Should return corresponding account when provide valid email", func() {
			req, _ := http.NewRequest("GET", fmt.Sprintf("/accounts?email=%s", testEmail), nil)
			router.ServeHTTP(testRecorder, req)
			Expect(testRecorder.Code).To(Equal(200))
			rsp := make(map[string]account.Account)
			json.Unmarshal(testRecorder.Body.Bytes(), &rsp)
			Expect(rsp["account"].ID).ToNot(BeEmpty())
			Expect(rsp["account"].Email).To(Equal(testEmail))
		})
	})
	Describe("PUT /accounts", func() {
		It("It should return account with status code 200 Giving email, password as body", func() {
			repoimpl.GetAccountRepo().DeleteOne(context.Background(), bson.M{"email": testEmail})
			req, _ := http.NewRequest("PUT", "/accounts", bytes.NewReader([]byte(fmt.Sprintf(`{"email": "%s", "password": "whatEver"}`, testEmail))))
			router.ServeHTTP(testRecorder, req)
			Expect(testRecorder.Code).To(Equal(200))
			rsp := make(map[string]interface{})
			json.Unmarshal(testRecorder.Body.Bytes(), &rsp)
			Expect(rsp["token"]).ToNot(BeEmpty())
		})
		It("It should return status code 400 and error Giving no email or password in body", func() {
			req, _ := http.NewRequest("PUT", "/accounts", bytes.NewReader([]byte(`{"email": "", "password": "whatEver"}`)))
			router.ServeHTTP(testRecorder, req)
			Expect(testRecorder.Code).To(Equal(400))

			req, _ = http.NewRequest("PUT", "/accounts", bytes.NewReader([]byte(`{"password": "whatEver"}`)))
			router.ServeHTTP(testRecorder, req)
			Expect(testRecorder.Code).To(Equal(400))

			req, _ = http.NewRequest("PUT", "/accounts", bytes.NewReader([]byte(`{"email": "whatEver@test.com"}`)))
			router.ServeHTTP(testRecorder, req)
			Expect(testRecorder.Code).To(Equal(400))

			req, _ = http.NewRequest("PUT", "/accounts", bytes.NewReader([]byte(`{"email": "testing", "password": "whatEver"}`)))
			router.ServeHTTP(testRecorder, req)
			Expect(testRecorder.Code).To(Equal(400))
		})
	})
	Describe("GET /accounts/:id", func() {
		It("It should return account with status code 200 Giving valid account id", func() {
			id := acct.ID.Hex()
			req, _ := http.NewRequest("GET", fmt.Sprintf("/accounts/%s", id), bytes.NewReader([]byte(`{"email": "", "password": "whatEver"}`)))
			router.ServeHTTP(testRecorder, req)
			Expect(testRecorder.Code).To(Equal(200))

			rsp := make(map[string]account.Account)
			json.Unmarshal(testRecorder.Body.Bytes(), &rsp)
			Expect(rsp["account"]).ToNot(BeNil())
			Expect(rsp["account"].Email).To(Equal(testEmail))
		})
		It("It should return error with status code 404 Giving invalid account id", func() {
			id := primitive.NewObjectID().Hex()
			req, _ := http.NewRequest("GET", fmt.Sprintf("/accounts/%s", id), bytes.NewReader([]byte(`{"email": "", "password": "whatEver"}`)))
			router.ServeHTTP(testRecorder, req)
			Expect(testRecorder.Code).To(Equal(404))
		})
		It("It should return error with status code 400 Giving invalid format account id", func() {
			req, _ := http.NewRequest("GET", "/accounts/12314", bytes.NewReader([]byte(`{"email": "", "password": "whatEver"}`)))
			router.ServeHTTP(testRecorder, req)
			Expect(testRecorder.Code).To(Equal(400))
		})
	})
	Describe("POST /accounts/:id", func() {
		Context("Giving valid id", func() {
			It("should 304 when body is empty", func() {
				req, _ := http.NewRequest("POST", fmt.Sprintf("/accounts/%s", acct.ID.Hex()), bytes.NewReader([]byte(`{}`)))
				router.ServeHTTP(testRecorder, req)
				Expect(testRecorder.Code).To(Equal(304))
			})
			It("Should update the mobile when mobile provided in body", func() {
				req, _ := http.NewRequest("POST", fmt.Sprintf("/accounts/%s", acct.ID.Hex()), bytes.NewReader([]byte(`{"mobile": "12345678910"}`)))
				router.ServeHTTP(testRecorder, req)
				Expect(testRecorder.Code).To(Equal(200))
				r := repoimpl.GetAccountRepo().FindOne(context.Background(), bson.M{"_id": acct.ID})
				v := account.Account{}
				r.Decode(&v)
				Expect(v.Mobile).To(Equal("12345678910"))
			})
			It("Should update the contactMethod when contactMethod provided in body", func() {
				cm := account.ContactMethods{SentSMS: true}
				cm.SentHook = account.SentHook{
					URL: "http://fake.com/sofake",
				}
				rawJson, _ := json.Marshal(cm)
				fmt.Printf(`{contactMethods: "%s"}`, rawJson)
				req, _ := http.NewRequest("POST", fmt.Sprintf("/accounts/%s", acct.ID.Hex()), bytes.NewReader([]byte(fmt.Sprintf(`{"mobile": "12345678910", "contactMethods": %s}`, string(rawJson)))))
				router.ServeHTTP(testRecorder, req)
				Expect(testRecorder.Code).To(Equal(200))
				r := repoimpl.GetAccountRepo().FindOne(context.Background(), bson.M{"_id": acct.ID})
				v := account.Account{}
				r.Decode(&v)
				Expect(v.Mobile).To(Equal("12345678910"))
				Expect(v.ContactMethods.SentHook.URL).To(Equal("http://fake.com/sofake"))
			})
		})
		It("Should return 404 when id invalid", func() {
			req, _ := http.NewRequest("POST", "/accounts/12345", bytes.NewReader([]byte(`{"mobile": "12345678910"}`)))
			router.ServeHTTP(testRecorder, req)
			Expect(testRecorder.Code).To(Equal(404))
		})
	})
	Describe("GET /login", func() {
		It("Should return 200 and token when valid email and password provided", func() {
			req, _ := http.NewRequest("POST", "/login", bytes.NewReader([]byte(fmt.Sprintf(`{"email": "%s", "password": "123"}`, testEmail))))
			router.ServeHTTP(testRecorder, req)
			Expect(testRecorder.Code).To(Equal(200))
			rsp := make(map[string]string)
			json.Unmarshal(testRecorder.Body.Bytes(), &rsp)
			Expect(rsp["token"]).ToNot(BeEmpty())
		})
		It("Should return 401 when no matching email and password", func() {
			req, _ := http.NewRequest("POST", "/login", bytes.NewReader([]byte(`{"email": "12131@fjfj.com", "password": "12344"}`)))
			router.ServeHTTP(testRecorder, req)
			Expect(testRecorder.Code).To(Equal(401))
		})
	})
})
