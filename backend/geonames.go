package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type GeonamesClient struct {
	searchURL string
	username  string
	ttl       time.Duration
	http      *http.Client

	mu        sync.Mutex
	names     []string
	fetchedAt time.Time
}

func NewGeonamesClient(searchURL, username string, ttl time.Duration) *GeonamesClient {
	return &GeonamesClient{
		searchURL: searchURL,
		username:  username,
		ttl:       ttl,
		http:      &http.Client{Timeout: 10 * time.Second},
	}
}

type searchResponse struct {
	Geonames []struct {
		Name string `json:"name"`
	} `json:"geonames"`
	Status *struct {
		Message string `json:"message"`
	} `json:"status,omitempty"`
}

func (c *GeonamesClient) CityNames(ctx context.Context) ([]string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if age := time.Since(c.fetchedAt); c.names != nil && age < c.ttl {
		log.Printf("cities: cache hit (%d cities, age %s, expires in %s)",
			len(c.names), age.Round(time.Second), (c.ttl - age).Round(time.Second))
		return c.names, nil
	}

	reason := "cache empty"
	if c.names != nil {
		reason = "cache expired (ttl " + c.ttl.String() + ")"
	}

	start := time.Now()
	names, err := c.fetch(ctx)
	if err != nil {
		return nil, err
	}
	c.names, c.fetchedAt = names, time.Now()
	log.Printf("cities: %s, fetched %d cities from geonames in %s",
		reason, len(names), time.Since(start).Round(time.Millisecond))
	return names, nil
}

func (c *GeonamesClient) fetch(ctx context.Context) ([]string, error) {
	q := url.Values{
		"featureClass": {"P"},
		"maxRows":      {"50"},
		"orderby":      {"population"},
		"username":     {c.username},
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.searchURL+"?"+q.Encode(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("geonames request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, fmt.Errorf("geonames read: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("geonames status %d: %s", resp.StatusCode, body)
	}

	var sr searchResponse
	if err := json.Unmarshal(body, &sr); err != nil {
		return nil, fmt.Errorf("geonames decode: %w", err)
	}
	if sr.Status != nil && sr.Status.Message != "" {
		return nil, fmt.Errorf("geonames error: %s", sr.Status.Message)
	}

	names := make([]string, 0, len(sr.Geonames))
	for _, g := range sr.Geonames {
		names = append(names, g.Name)
	}
	return names, nil
}
