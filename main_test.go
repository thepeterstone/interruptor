package main_test

import (
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/thepeterstone/interruptor"
)

var _ = Describe("Main", func() {
	It("prints usage when request method is GET", func() {
		w := &spyWriter{}
		Interrupt(w, &http.Request{Method: http.MethodGet})

		Expect(w.headers).To(ConsistOf(200))
		Expect(string(w.body)).To(Equal(Usage()))
	})

	It("rejects unknown request methods", func() {
		w := &spyWriter{}
		Interrupt(w, &http.Request{Method: http.MethodPut})

		Expect(w.headers).To(ConsistOf(405))
		Expect(string(w.body)).To(BeEmpty())
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
	s.body = append(s.body, b...)
	return len(b), nil
}

// spyWriter implements http.ResponseWriter
func (s *spyWriter) WriteHeader(h int) {
	s.headers = append(s.headers, h)
}
