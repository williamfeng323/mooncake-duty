package shift_test

import (
	"context"
	"testing"
	"time"
	"williamfeng323/mooncake-duty/src/domains/project"
	"williamfeng323/mooncake-duty/src/domains/shift"
	repoimpl "williamfeng323/mooncake-duty/src/infrastructure/db/repo_impl"

	"go.mongodb.org/mongo-driver/bson/primitive"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.mongodb.org/mongo-driver/bson"
)

func TestShift(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Shift domain Suite")
}

var _ = Describe("Shift", func() {
	Describe("#NewShift", func() {
		Context("If project does not exist", func() {
			It("should return project not found error and nil shift", func() {
				shift, err := shift.NewShift(primitive.NewObjectID(), shift.Mon, time.Now(), shift.Weekly)
				Expect(shift).To(BeNil())
				Expect(err).To(MatchError(project.NotFoundError{}))
			})
		})
		Context("If the project exists", func() {
			project := project.NewProject("TestShift", "Project for testing the shift")
			project.Create()

			AfterSuite(func() {
				repoimpl.GetProjectRepo().DeleteOne(context.Background(), bson.M{"_id": project.ID})
			})

			Context("If the start date is not the first date of the shift", func() {
				It("should replace it to the first date of the shift base on weekStart", func() {

				})
			})
		})
	})
	Describe("#DefaultShifts", func() {
		Context("If no member in the domain", func() {
			It("should return nil instantly", func() {
			})
		})
	})
})
