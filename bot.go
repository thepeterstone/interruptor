package interruptor

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"

	"github.com/nlopes/slack"
)

type SlackRequest struct {
	VerificationToken string `json:"token"`
	Type              string `json:"type"`
}

func (r *SlackRequest) validate(token string) bool {
	return r.VerificationToken == token
}

type SlackChallenge struct {
	Challenge string `json:"challenge"`
	SlackRequest
}

type AppMention struct {
	Event slack.Msg `json:"event"`
	SlackRequest
}

var (
	channels []string
	users    []string
)

func SlackResponder(cfg *Config) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		t, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var q SlackRequest
		if err := json.Unmarshal(t, &q); err != nil {
			return
		}

		if !q.validate(cfg.VerificationToken) {
			log.Printf("bad token: %s\n", q.VerificationToken)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		switch q.Type {
		case "url_verification":
			echoChallenge(w, t)
			return
		case "event_callback":
			log.Println("event callback:")
			processMessage(w, t)
			return
		}

		log.Println("Unknown request type.")
		logRequest(r, t)
		w.WriteHeader(http.StatusBadRequest)
	}
}

func echoChallenge(w http.ResponseWriter, t []byte) {
	var c SlackChallenge
	if err := json.Unmarshal(t, &c); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	log.Printf("%+v\n", c)
	fmt.Fprintf(w, "%s", c.Challenge)
}

func processMessage(w http.ResponseWriter, t []byte) {
	var m AppMention
	if err := json.Unmarshal(t, &m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Printf("\n%+v\n\n", m.Event)
}

func setInterruptChannels(message string) {
	var c []string
	re := regexp.MustCompile(`<#(\w+)\|(\w+)>`)
	m := re.FindAllStringSubmatch(message, -1)
	for _, n := range m[1] {
		c = append(c, n)
	}
	log.Println(c)
	channels = c
}

func setChannelTopics(api *slack.Client, channels []string, message string) {
	for _, id := range channels {
		_, err := api.SetChannelTopic(
			id,
			message,
		)
		if err != nil {
			log.Printf("error setting topic for %s: %s", id, err)
			continue
		}
	}
}

func PostPrinter(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	t, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	logRequest(r, t)
}

func logRequest(r *http.Request, t []byte) {
	log.Printf("Request: %+v\n", r)
	log.Printf("Body: %s\n\n", t)
}
