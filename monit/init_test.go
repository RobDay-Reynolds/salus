package monit_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

var (
	fixturesPath string
)

func TestSky(t *testing.T) {
	fixturesPath = "files"

	RegisterFailHandler(Fail)
	RunSpecs(t, "monit")
}
