package filereader_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestFileReader(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "File Reader Suite")
}
