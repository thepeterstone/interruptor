package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
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
}

func Usage() string {
	return "Use in slack with <code>/interrupt [users]</code>"
}
