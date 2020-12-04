package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", handleHome())

	log.Fatal(http.ListenAndServe(":5000", router))
}

func handleHome() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("Hello world"))
	}
}
