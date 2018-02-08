package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestInterruptor(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Interruptor Suite")
}
