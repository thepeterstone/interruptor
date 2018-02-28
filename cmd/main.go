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
	r.HandleFunc("/api", interruptor.SlackResponder(cfg, &log.Logger{}))

	http.Handle("/", r)
	log.Printf("%+v\n", cfg)

	log.Println("Starting interruptor server...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func Readme(w http.ResponseWriter, _ *http.Request) {
	readme, _ := ioutil.ReadFile("README.md")
	w.Write(readme)
}

func loadConfig() *interruptor.Config {
	var cfg interruptor.Config
	if err := envstruct.Load(&cfg); err != nil {
		log.Fatalf("failed to load config: %s", err)
	}

	return &cfg
}
