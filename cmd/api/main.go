package main

import (
	"fmt"
	"log"
	"net/http"
)

const webPort = "80"

type Config struct{}

func main() {
	app := Config{}

	// Defines the server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	log.Printf("Starting broker service on port %s", webPort)

	// Start the server
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
