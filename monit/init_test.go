package monit_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

var (
	fixturesPath string
)

func TestMonit(t *testing.T) {
	fixturesPath = "files"

	RegisterFailHandler(Fail)
	RunSpecs(t, "monit")
}
