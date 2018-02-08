package main_test

import (
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/thepeterstone/interruptor"
)

var _ = Describe("Main", func() {
	It("fails when no request is provided", func() {
		w := &spyWriter{}
		Interrupt(w, &http.Request{})

		Expect(w.headers).To(ConsistOf(200))
	})
})

type spyWriter struct {
	body     []byte
	writeErr error
	headers  []int
}

// spyWriter implements http.ResponseWriter
func (s *spyWriter) Header() http.Header {
	return nil
}

// spyWriter implements http.ResponseWriter
func (s *spyWriter) Write(b []byte) (int, error) {
	return 0, nil
}

// spyWriter implements http.ResponseWriter
func (s *spyWriter) WriteHeader(h int) {
	s.headers = append(s.headers, h)
}
