package interruptor

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var (
	log = DefaultLogger()
)

type SlackChallenge struct {
	Challenge         string `json:"challenge"`
	VerificationToken string `json:"token"`
}

func ChallengeEchoer(cfg *Config) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var c SlackChallenge
		t, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err := json.Unmarshal(t, &c); err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if c.VerificationToken != cfg.VerificationToken {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		fmt.Fprintf(w, "%s", c.Challenge)
	}
}
