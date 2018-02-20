package interruptor_test

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/nlopes/slack"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/thepeterstone/interruptor"
)

var _ = Describe("SlackResponder", func() {
	var (
		handler    func(w http.ResponseWriter, r *http.Request)
		challenge  interruptor.SlackChallenge
		appMention interruptor.AppMention
		w          *spyWriter
	)

	BeforeEach(func() {
		config := &interruptor.Config{
			VerificationToken: "some-token",
		}
		baseRequest := interruptor.SlackRequest{
			Type:              "url_verification",
			VerificationToken: "some-token",
		}

		handler = interruptor.SlackResponder(config)
		challenge = interruptor.SlackChallenge{
			Challenge:    "some-challenge-string",
			SlackRequest: baseRequest,
		}
		appMention = interruptor.AppMention{
			SlackRequest: baseRequest,
			Event: slack.Msg{
				Text: "some-text",
			},
		}
		w = newSpyWriter()
	})

	It("echoes a json-encoded challenge", func() {
		r := requestWithChallenge(challenge)
		handler(w, r)

		Expect(w.headers).To(ContainElement(http.StatusOK))
		Expect(w.text).To(ConsistOf("some-challenge-string"))
	})

	It("acknowledges an app mention", func() {
		r := requestWithAppMention(appMention)
		handler(w, r)

		Expect(w.headers).To(ContainElement(http.StatusOK))
		Expect(w.text).To(ConsistOf(""))
	})

	It("returns an authorization error if the token doesn't match", func() {
		challenge.SlackRequest.VerificationToken = "invalid-token"
		r := requestWithChallenge(challenge)
		handler(w, r)

		Expect(w.headers).To(ContainElement(http.StatusUnauthorized))
		Expect(w.text).To(BeEmpty())
	})

	It("returns a bad request error if the type isn't recognized", func() {
		challenge.SlackRequest.Type = "unknown-type"
		r := requestWithChallenge(challenge)
		handler(w, r)

		Expect(w.headers).To(ContainElement(http.StatusBadRequest))
		Expect(w.text).To(BeEmpty())
	})

})

func requestWithAppMention(m interruptor.AppMention) *http.Request {
	body, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	r, err := http.NewRequest(http.MethodGet, "", bytes.NewReader(body))
	if err != nil {
		panic(err)
	}
	return r
}

func requestWithChallenge(c interruptor.SlackChallenge) *http.Request {
	body, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}
	r, err := http.NewRequest(http.MethodGet, "", bytes.NewReader(body))
	if err != nil {
		panic(err)
	}
	return r
}

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
