package syncservice_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSyncService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "syncservice Suite")
}
