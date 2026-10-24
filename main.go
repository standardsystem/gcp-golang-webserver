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
		wwwPort = "8080"
	}

	http.HandleFunc("/", AuthMiddleware(HandleRoot))

	health := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	}
	http.HandleFunc("/health", health)
	http.HandleFunc("/healthz", health)
	http.HandleFunc("/readiness", health)

	log.Println("Ready to serve")
	http.ListenAndServe(":"+wwwPort, nil)
}
