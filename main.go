package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	log.Println("Starting")

	wwwPort := os.Getenv("PORT")
	if wwwPort == "" {
		wwwPort = "10080"
	}

	http.HandleFunc("/", HandleRoot)

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	http.HandleFunc("/readiness", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	log.Println("Ready to serve")
	http.ListenAndServe(":"+wwwPort, nil)
}
