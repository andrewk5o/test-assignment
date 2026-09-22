package main

import (
	"log"
	"net/http"
	"os"
	"time"
)

const geonamesSearchURL = "http://api.geonames.org/searchJSON"

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func main() {
	port := env("PORT", "8080")
	username := env("GEONAMES_USERNAME", "hsample")
	ttl, err := time.ParseDuration(env("CACHE_TTL", "10m"))
	if err != nil {
		log.Fatalf("invalid CACHE_TTL: %v", err)
	}

	src := NewGeonamesClient(geonamesSearchURL, username, ttl)
	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           NewServer(src),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("listening on :%s (geonames user=%s, cache ttl=%s)", port, username, ttl)
	log.Fatal(srv.ListenAndServe())
}
