package account

import (
	"context"
	"testing"
	repoimpl "williamfeng323/mooncake-duty/src/infrastructure/db/repo_impl"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.mongodb.org/mongo-driver/bson"
)

func TestAccountService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Account Service Suite")
}

var _ = Describe("AccountService", func() {
	repo := repoimpl.GetAccountRepo()
	acctSrv := &Service{}
	acctSrv.SetRepo(repoimpl.GetAccountRepo())
	BeforeEach(func() {
		acct2, _ := NewAccount("test@test.com", "12345test")
		acct2.IsAdmin = true
		acct2.Save()
	})
	AfterEach(func() {
		repo.DeleteOne(context.Background(), bson.M{"email": "test@test.com"})
	})
	Describe("#SignIn", func() {
		Context("call with incorrect email/password", func() {
			It("should return error", func() {
				rst, err := acctSrv.SignIn("test@test.com", "12345")
				Expect(rst).To(Equal(""))
				Expect(err).ToNot(BeNil())
			})
		})
		Context("call with correct email/password", func() {
			It("should jwt token", func() {
				rst, err := acctSrv.SignIn("test@test.com", "12345test")
				Expect(rst).ToNot(BeNil())
				Expect(err).To(BeNil())
			})
		})
		Context("call with non-exist email", func() {
			It("should return error", func() {
				rst, err := acctSrv.SignIn("test1@test.com", "12345test")
				Expect(rst).To(Equal(""))
				Expect(err.Error()).To(Equal("Account does not exist"))
			})
		})
	})
})
