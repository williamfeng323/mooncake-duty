package middlewares_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"williamfeng323/mooncake-duty/src/domains/account"
	repoimpl "williamfeng323/mooncake-duty/src/infrastructure/db/repo_impl"
	"williamfeng323/mooncake-duty/src/infrastructure/middlewares"
	"williamfeng323/mooncake-duty/src/utils"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/mgo.v2/bson"
)

func TestAuthMiddleware(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Auth middleware Suite")
}

var _ = Describe("Authenticate", func() {
	var router *gin.Engine
	var testRecorder *httptest.ResponseRecorder
	BeforeEach(func() {
		router = gin.Default()
		router.Use(middlewares.Authenticate())
		router.GET("/ping", func(c *gin.Context) {
			c.String(200, "pong")
		})
		testRecorder = httptest.NewRecorder()
	})
	It("should return 403 if header not provided", func() {
		req, _ := http.NewRequest("GET", "/ping", nil)
		router.ServeHTTP(testRecorder, req)
		Expect(testRecorder.Code).To(Equal(403))
	})
	It("should return 403 if header not valid", func() {
		req, _ := http.NewRequest("GET", "/ping", nil)
		req.Header.Set(middlewares.AuthHeader, "1312314124124")
		router.ServeHTTP(testRecorder, req)
		Expect(testRecorder.Code).To(Equal(403))
	})
	It("should return 403 if header not valid", func() {
		tk, _ := utils.SignToken("123@mock.com")
		req, _ := http.NewRequest("GET", "/ping", nil)
		req.Header.Set(middlewares.AuthHeader, tk)
		router.ServeHTTP(testRecorder, req)
		Expect(testRecorder.Code).To(Equal(403))
	})
	It("should return 200 if token and account valid", func() {
		token, acct, _ := account.GetAccountService().Register("test@mock.com", "123123", false)
		req, _ := http.NewRequest("GET", "/ping", nil)
		req.Header.Set(middlewares.AuthHeader, token)
		router.ServeHTTP(testRecorder, req)
		Expect(testRecorder.Code).To(Equal(200))
		Expect(testRecorder.Body.Bytes()).To(Equal([]byte("pong")))
		repoimpl.GetAccountRepo().DeleteOne(context.Background(), bson.M{"_id": acct.ID})
	})
})
