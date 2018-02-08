package main

import (
	"fmt"
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
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("%+v", r)))
}
