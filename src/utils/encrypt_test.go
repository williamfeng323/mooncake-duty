package utils

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestModel(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Encrypt Suite")
}

var _ = Describe("#Encrypt", func() {
	Context("With predefined string", func() {
		rawString := "teststring"
		It("should be able to encrypt and decrpt to raw string", func() {
			encryptedString, err := Encrypt(rawString)
			Expect(err).To(BeNil())
			decryptedString, err := Decrypt(string(encryptedString))
			Expect(string(decryptedString)).To(Equal(rawString))
		})
	})
})
