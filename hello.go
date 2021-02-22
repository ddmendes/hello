package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const nameKey string = "name"

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", hello)
	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	log.Fatal(server.ListenAndServe())
}

func hello(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue(nameKey)
	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing name parameter"))
	} else {
		w.Write([]byte(fmt.Sprintf("Hello, %s!", name)))
	}
}
