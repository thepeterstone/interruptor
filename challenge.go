package interruptor

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type SlackRequest struct {
	VerificationToken string `json:"token"`
}

func (r *SlackRequest) validate(token string) bool {
	return r.VerificationToken == token
}

type SlackChallenge struct {
	Challenge string `json:"challenge"`
	SlackRequest
}

func ChallengeEchoer(cfg *Config) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		t, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		logRequest(r, t)

		var q SlackRequest
		if err := json.Unmarshal(t, &q); err != nil {
			return
		}

		if !q.validate(cfg.VerificationToken) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		var c SlackChallenge
		if err := json.Unmarshal(t, &c); err != nil {
			PostPrinter(w, r)
			return
		}

		fmt.Fprintf(w, "%s", c.Challenge)
	}
}

func PostPrinter(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func logRequest(r *http.Request, t []byte) {
	log.Printf("Request: %+v\n", r)
	log.Printf("Body: %s\n", t)
}
