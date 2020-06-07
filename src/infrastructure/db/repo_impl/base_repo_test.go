package repoimpl

import (
	"testing"
	"williamfeng323/mooncake-duty/src/infrastructure/db"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestModel(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Model domain Suite")
}

var _ = Describe("BaseRepo", func() {
	connection := db.GetConnection()
	repo := &BaseRepo{
		name: "Test",
	}
	db.Register(repo)
	repo.SetCollection(connection.GetCollection("Test"))
	Context("#Find", func() {
		It("", func() {})
	})
})
