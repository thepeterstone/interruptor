package interruptor

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	l "log"
	"net/http"
	"regexp"
	"strings"

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
	api      *slack.Client
	config   *Config
	log      Logger = &l.Logger{}
)

type Logger interface {
	Println(v ...interface{})
	Printf(format string, v ...interface{})
}

func SlackResponder(cfg *Config, l Logger) func(w http.ResponseWriter, r *http.Request) {
	log = l
	api = slack.New(cfg.ApiKey)
	config = cfg

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

func processMessage(w http.ResponseWriter, e []byte) {
	var m AppMention
	if err := json.Unmarshal(e, &m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	users := parseUsers(m.Event.Text)
	if len(users) < 1 {
		w.Write([]byte("Couldn't find a user, try using @?"))
		return
	}

	var u []string
	for _, n := range users {
		ui, err := api.GetUserInfo(n)
		if err != nil {
			log.Println(err)
			continue
		}
		u = append(u, "@"+ui.Name)
	}

	t := fmt.Sprintf(
		"%s %s",
		"Interrupt:",
		strings.Join(u, " "),
	)
	c := channels
	if c == nil {
		c = []string{m.Event.Channel}
	}
	setChannelTopics(api, c, t)
	log.Printf("\n%+v\n\n", m)
}

func parseUsers(message string) []string {
	var u []string
	re := regexp.MustCompile(`<@(\w+)>`)
	m := re.FindAllStringSubmatch(message, -1)
	if len(m) < 2 {
		return nil
	}
	for _, n := range m[1] {
		u = append(u, n)
	}
	log.Println(u)
	return u
}

func parseChannels(message string) []string {
	var c []string
	re := regexp.MustCompile(`<#(\w+)\|(\w+)>`)
	m := re.FindAllStringSubmatch(message, -1)
	for _, n := range m[1] {
		c = append(c, n)
	}
	log.Println(c)
	return c
}

func setChannelTopics(api *slack.Client, channels []string, message string) {
	log.Printf("Setting %s topic: %s\n", channels, message)
	for _, id := range channels {
		t, err := api.SetChannelTopic(
			id,
			message,
		)
		if err != nil {
			log.Printf("error setting topic for %s: %s", id, err)
			continue
		}
		log.Printf("[%s]: %s\n", id, t)
	}
}

func logRequest(r *http.Request, t []byte) {
	log.Printf("Request: %+v\n", r)
	log.Printf("Body: %s\n\n", t)
}
