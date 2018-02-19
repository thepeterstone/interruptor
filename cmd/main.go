package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"

	envstruct "github.com/cloudfoundry/go-envstruct"
	"github.com/gorilla/mux"
	"github.com/nlopes/slack"
	"github.com/thepeterstone/interruptor"
)

var (
	cfg *interruptor.Config
)

func main() {
	cfg = loadConfig()
	r := mux.NewRouter()
	r.HandleFunc("/", Readme)
	r.HandleFunc("/debug", PostPrinter)
	r.HandleFunc("/challenge", interruptor.ChallengeEchoer(cfg))
	r.HandleFunc("/interrupt", Interrupt)
	r.HandleFunc("/interrupt-channels", InterruptChannels)

	http.Handle("/", r)
	log.Printf("%+v\n", cfg)

	log.Println("Starting interruptor server...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func Interrupt(w http.ResponseWriter, r *http.Request) {
	if !checkMethod(w, r) {
		return
	}

	s, err := slack.SlashCommandParse(r)
	if err != nil {
		log.Printf("error decoding command: %s\n", err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, Usage())
		return
	}

	if !s.ValidateToken(cfg.VerificationToken) {
		log.Printf("token invalid: %s\n", cfg.VerificationToken)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if s.Text == "" {
		fmt.Fprint(w, Usage())
		return
	}

	api := slack.New(cfg.ApiKey)
	setInterrupt(api, s)
}

func InterruptChannels(w http.ResponseWriter, r *http.Request) {
	if !checkMethod(w, r) {
		return
	}

	s, err := slack.SlashCommandParse(r)
	if err != nil {
		log.Printf("error decoding command: %s\n", err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, Usage())
		return
	}

	if !s.ValidateToken(cfg.VerificationToken) {
		log.Printf("token invalid: %s\n", cfg.VerificationToken)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if s.Text == "" {
		fmt.Fprint(w, Usage())
		return
	}

	log.Printf("%+v", s)
	api := slack.New(cfg.ApiKey)
	setInterruptChannels(api, s)

	var names []string
	for _, c := range cfg.Channels {
		ci, err := api.GetChannelInfo(c)
		if err != nil {
			log.Printf("error getting channel info: %s", err)
			continue
		}
		names = append(names, ci.Name)
	}

	fmt.Fprintf(w, "OK, set interrupt channels to %s", names)
}

func PostPrinter(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	log.Printf("Request: %+v\n", r)
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %s", err)
		return
	}

	log.Printf("Body: %s", data)
	w.WriteHeader(http.StatusOK)
}

func Readme(w http.ResponseWriter, _ *http.Request) {
	readme, _ := ioutil.ReadFile("README.md")
	w.Write(readme)
}

func Usage() string {
	return "Use with /interrupt [users]."
}

func checkMethod(w http.ResponseWriter, r *http.Request) bool {
	if r.Method == http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(Usage()))
		return false
	}
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return false
	}
	return true
}

func setInterruptChannels(api *slack.Client, cmd slack.SlashCommand) {
	var c []string
	for _, n := range strings.Split(cmd.Text, " ") {
		re := regexp.MustCompile(`<#(\w+)\|(\w+)>`)
		m := re.FindStringSubmatch(n)
		c = append(c, m[1])
	}
	log.Println(c)
	cfg.Channels = c
}
func setInterrupt(api *slack.Client, cmd slack.SlashCommand) {
	for _, id := range cfg.Channels {
		t, err := api.SetChannelTopic(
			id,
			fmt.Sprintf("%s %s", cfg.MessagePrefix, cmd.Text),
		)
		if err != nil {
			log.Printf("error setting topic for %s: %s", id, err)
			continue
		}
		log.Printf("%s set topic '%s' for %s", cmd.UserName, t, id)
	}
}

func loadConfig() *interruptor.Config {
	var cfg interruptor.Config
	if err := envstruct.Load(&cfg); err != nil {
		log.Fatalf("failed to load config: %s", err)
	}

	return &cfg
}
