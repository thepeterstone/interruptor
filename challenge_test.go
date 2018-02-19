package interruptor_test

import (
	"bytes"
	"encoding/json"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/thepeterstone/interruptor"
)

var _ = Describe("Challenge", func() {
	var (
		config *interruptor.Config
	)

	BeforeEach(func() {
		config = &interruptor.Config{
			VerificationToken: "some-token",
		}

	})
	It("echoes a json-encoded challenge", func() {
		body, err := json.Marshal(interruptor.SlackChallenge{
			Challenge:         "some-challenge-string",
			VerificationToken: "some-token",
		})
		if err != nil {
			panic(err)
		}
		r, _ := http.NewRequest(http.MethodGet, "", bytes.NewReader(body))
		w := newSpyWriter()

		handler := interruptor.ChallengeEchoer(config)
		handler(w, r)

		Expect(w.text).To(ContainElement("some-challenge-string"))
	})

	It("returns an authorization error if the token doesn't match", func() {
		body, err := json.Marshal(interruptor.SlackChallenge{
			Challenge:         "some-challenge-string",
			VerificationToken: "invalid-token",
		})
		if err != nil {
			panic(err)
		}
		r, _ := http.NewRequest(http.MethodGet, "", bytes.NewReader(body))
		w := newSpyWriter()

		handler := interruptor.ChallengeEchoer(config)
		handler(w, r)

		Expect(w.headers).To(ContainElement(http.StatusUnauthorized))
		Expect(w.text).To(BeEmpty())
	})
})

type spyWriter struct {
	headers []int
	text    []string
}

func newSpyWriter() *spyWriter {
	return &spyWriter{}
}

// spyWriter implements io.Writer
func (w *spyWriter) Write(p []byte) (int, error) {
	w.text = append(w.text, string(p))
	return len(p), nil
}

// spyWriter implements http.ResponseWriter
func (w *spyWriter) Header() http.Header {
	return nil
}

func (w *spyWriter) WriteHeader(h int) {
	w.headers = append(w.headers, h)
}
