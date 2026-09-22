package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"unicode"
	"unicode/utf8"
)

type CitySource interface {
	CityNames(ctx context.Context) ([]string, error)
}

type countResponse struct {
	Letter string   `json:"letter"`
	Count  int      `json:"count"`
	Cities []string `json:"cities"`
}

func NewServer(src CitySource) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/cities/count", func(w http.ResponseWriter, r *http.Request) {
		handleCount(w, r, src)
	})
	return mux
}

func handleCount(w http.ResponseWriter, r *http.Request, src CitySource) {
	letter := r.URL.Query().Get("letter")
	if !isSingleLetter(letter) {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "letter must be exactly one alphabetic character"})
		return
	}

	names, err := src.CityNames(r.Context())
	if err != nil {
		log.Printf("upstream error: %v", err)
		writeJSON(w, http.StatusBadGateway, map[string]string{"error": "could not fetch cities from geonames"})
		return
	}

	matches := CitiesStartingWith(names, letter)
	writeJSON(w, http.StatusOK, countResponse{Letter: letter, Count: len(matches), Cities: matches})
}

func isSingleLetter(s string) bool {
	r, size := utf8.DecodeRuneInString(s)
	return size > 0 && size == len(s) && unicode.IsLetter(r)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("write response: %v", err)
	}
}
