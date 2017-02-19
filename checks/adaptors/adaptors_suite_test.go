package adaptors_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestAdaptors(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Adaptors Suite")
}
