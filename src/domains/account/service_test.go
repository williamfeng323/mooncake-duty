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
		acct2.Save(true)
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
				Expect(err.Error()).To(Equal("Account Not Found"))
			})
		})
	})
	Describe("#Register", func() {
		Context("call with exist email", func() {
			It("should return error", func() {
				rst, err := acctSrv.Register("test@test.com", "12345test", true)
				Expect(rst).To(Equal(""))
				Expect(err.Error()).To(Equal("Account already exist"))
			})
		})
		Context("call with non-exist email", func() {
			It("should return jwt token", func() {
				rst, err := acctSrv.Register("test1@test.com", "12345test", true)
				Expect(rst).ToNot(BeNil())
				Expect(err).To(BeNil())
				repo.DeleteOne(context.Background(), bson.M{"email": "test1@test.com"})
			})
		})
	})
	Describe("#UpdateContactMethods", func() {
		Context("call with exist account", func() {
			Context("update partial ContactMethods", func() {
				It("should throw error if enable sendSMS when no mobile set", func() {
					cm := ContactMethods{SentSMS: true}
					err := acctSrv.UpdateContactMethods("test@test.com", cm, "", "")
					Expect(err.Error()).To(Equal("Mobile must be set before you active send email notification"))
				})
				It("should update the contactMethods", func() {
					cm := ContactMethods{SentSMS: true}
					cm.SentHook = SentHook{
						URL: "http://fake.com/sofake",
					}
					err := acctSrv.UpdateContactMethods("test@test.com", cm, "", "12345678910")
					Expect(err).To(BeNil())
				})
			})
			Context("when existing account has contactMethods", func() {
				It("Should update the account without error", func() {
					acct, _ := NewAccount("test@test.com", "12345test")
					acct.ContactMethods.SentEmail = true
					acct.Mobile = "12345678910"
					acct.Save(true)
					cm := ContactMethods{SentSMS: true, SentEmail: false}
					cm.SentHook = SentHook{
						URL: "http://fake.com/sofake",
					}
					err := acctSrv.UpdateContactMethods("test@test.com", cm, "", "")
					Expect(err).To(BeNil())
				})
			})
		})
		Context("call with non-exist account", func() {
			It("should return error", func() {
				err := acctSrv.UpdateContactMethods("test2@test.com", ContactMethods{}, "", "")
				Expect(err).ToNot(BeNil())
			})
		})
	})
})
