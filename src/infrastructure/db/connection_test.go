package db_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"williamfeng323/mooncake-duty/src/infrastructure/db"
	repoimpl "williamfeng323/mooncake-duty/src/infrastructure/db/repo_impl"
	"williamfeng323/mooncake-duty/src/utils"
)

func TestModel(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Connection Suite")
}

type TestRepo struct {
	repoimpl.BaseRepo
	name string
}

func (tr *TestRepo) GetName() string {
	return "Test"
}

var _ = Describe("Connection", func() {
	// Only pass t into top-level Convey calls
	Context("Given an empty context", func() {
		Describe("When init Connect with proper conf", func() {
			db := &db.Connection{}
			conf := utils.MongoConfig{
				URL:               "localhost",
				Port:              "27017",
				Username:          "root",
				Password:          "example",
				Database:          "mooncake",
				ConnectionOptions: "authSource=admin",
			}
			err := db.InitConnection(nil, conf)
			pingError := db.Client.Ping(nil, nil)
			It("a client should be initialized and no error when ping", func() {
				Expect(err).To(BeNil())
				Expect(pingError).To(BeNil())
				Expect(db.Client).ToNot(BeNil())
			})
		})
		Describe("When trigger Connect without conf", func() {
			db := &db.Connection{}
			conf := utils.Config{}
			err := db.InitConnection(nil, conf.Mongo)
			It("should return error", func() {
				Expect(err).ToNot(BeNil())
				Expect(db.Client).To(BeNil())
			})
		})
	})
	Describe("#Register", func() {
		Context("Giving a test struct", func() {
			Describe("when register to a Connection instance", func() {
				testRepo := &TestRepo{}
				db.Register(testRepo)
				It("should register the collection to connection", func() {
					conn := db.GetConnection()
					Expect(conn.GetCollection("Test")).ToNot(BeNil())
				})
			})
		})
	})
})
