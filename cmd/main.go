package main

import (
	"io/ioutil"
	"log"
	"net/http"

	envstruct "github.com/cloudfoundry/go-envstruct"
	"github.com/gorilla/mux"
	"github.com/thepeterstone/interruptor"
)

func main() {
	cfg := loadConfig()
	r := mux.NewRouter()
	r.HandleFunc("/", Readme)
	r.HandleFunc("/api", interruptor.SlackResponder(cfg, &log.Logger{}))

	http.Handle("/", r)

	log.Println("Starting interruptor server...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Readme writes the readme file to the provided ResponseWriter
func Readme(w http.ResponseWriter, _ *http.Request) {
	readme, _ := ioutil.ReadFile("README.md")
	_, _ = w.Write(readme)
}

func loadConfig() *interruptor.Config {
	var cfg interruptor.Config
	if err := envstruct.Load(&cfg); err != nil {
		log.Fatalf("failed to load config: %s", err)
	}

	return &cfg
}
