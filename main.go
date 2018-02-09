package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", Readme)
	r.HandleFunc("/interrupt", Interrupt)
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func Interrupt(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(Usage()))
		return
	}
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

func Readme(w http.ResponseWriter, _ *http.Request) {
	readme, _ := ioutil.ReadFile("README.md")
	w.Write(readme)
}

func Usage() string {
	return "Use in slack with <code>/interrupt [users]</code>"
}
