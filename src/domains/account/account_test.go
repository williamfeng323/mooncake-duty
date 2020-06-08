package account

import (
	"context"
	"testing"
	repoimpl "williamfeng323/mooncake-duty/src/infrastructure/db/repo_impl"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.mongodb.org/mongo-driver/bson"
)

func TestAccount(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Account Domain Suite")
}

var _ = Describe("Account", func() {
	repo := repoimpl.GetAccountRepo()
	AfterEach(func() {
		repo.DeleteOne(context.Background(), bson.M{"email": "test@test.com"})
	})
	Describe("#Save", func() {
		Context("When database does not have record in same email", func() {
			It("should insert the account to DB directly", func() {
				acct, err := NewAccount("test@test.com", "12345test")
				Expect(err).To(BeNil())
				rst, err := acct.Save()
				Expect(err).To(BeNil())
				Expect(rst).To(Equal(1))
			})
		})
		Context("When database has record in same email", func() {
			It("should update the account to DB directly", func() {
				acct, err := NewAccount("test@test.com", "12345test")
				acct.Save()
				acct2, _ := NewAccount("test@test.com", "12345test")
				acct2.IsAdmin = true
				rst, err := acct2.Save()
				dbAcct := &Account{}
				rst2 := repo.FindOne(context.Background(), bson.M{"email": "test@test.com"})
				rst2.Decode(dbAcct)
				Expect(dbAcct.IsAdmin).To(BeTrue())
				Expect(err).To(BeNil())
				Expect(rst).To(Equal(1))
			})
		})
	})
})
