package project

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestModel(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Project domain Suite")
}

var _ = Describe("Project", func() {
	Describe("#NewProject", func() {
		It("should return a project instance", func() {
			project1 := NewProject("Test", "This is a test project", Member{primitive.NewObjectID(), true})
			Expect(project1).NotTo(BeNil())
			project2 := NewProject("Test", "This is a test project")
			Expect(project2).NotTo(BeNil())
		})
	})
})
