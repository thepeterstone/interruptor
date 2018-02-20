package main

import (
	"io/ioutil"
	"log"
	"net/http"

	envstruct "github.com/cloudfoundry/go-envstruct"
	"github.com/gorilla/mux"
	"github.com/thepeterstone/interruptor"
)

var (
	cfg *interruptor.Config
)

func main() {
	cfg = loadConfig()
	r := mux.NewRouter()
	r.HandleFunc("/", Readme)
	r.HandleFunc("/debug", interruptor.PostPrinter)
	r.HandleFunc("/api", interruptor.SlackResponder(cfg))

	http.Handle("/", r)
	log.Printf("%+v\n", cfg)

	log.Println("Starting interruptor server...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func Readme(w http.ResponseWriter, _ *http.Request) {
	readme, _ := ioutil.ReadFile("README.md")
	w.Write(readme)
}

func Usage() string {
	return "Use with /interrupt [users]."
}

func loadConfig() *interruptor.Config {
	var cfg interruptor.Config
	if err := envstruct.Load(&cfg); err != nil {
		log.Fatalf("failed to load config: %s", err)
	}

	return &cfg
}
