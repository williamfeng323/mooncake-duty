package validatorimpl

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"williamfeng323/mooncake-duty/src/infrastructure/db"
)

func TestModel(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Model domain Suite")
}

var _ = Describe("DefaultValidator", func() {
	Describe("Giving a struct with tag required", func() {
		type TestModel struct {
			db.BaseModel `json:"inline" bson:"inline"`
			TestWhat     string `json:"testing" bson:"testing" required:"true"`
		}
		var validator db.Validator
		BeforeEach(func() {
			validator = NewDefaultValidator()
		})
		It("Should throw error if the document does not init the required field", func() {
			testModel := &TestModel{}
			err := validator.Verify(testModel)
			Expect(err).NotTo(BeNil())
		})
		It("Should not throw error if the document initted the required field", func() {
			test := &TestModel{}
			test.CreatedAt = time.Now()
			test.ID = primitive.NewObjectID()
			test.TestWhat = "What to test"
			testValue := reflect.ValueOf(test)
			ti := testValue.Interface()
			fmt.Printf("Value of test is %v", ti)
			err := validator.Verify(test)
			Expect(err).To(BeNil())
		})
	})
})
