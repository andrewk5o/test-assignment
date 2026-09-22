package main

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

const sampleJSON = `{"totalResultsCount":3,"geonames":[
  {"geonameId":1,"name":"Shanghai","population":24874500},
  {"geonameId":2,"name":"Cairo","population":20000000},
  {"geonameId":3,"name":"Rio de Janeiro","population":6000000}]}`

func newUpstream(t *testing.T, status int, body string, hits *int32) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(hits, 1)
		if r.URL.Query().Get("username") != "testuser" {
			t.Errorf("username query = %q, want testuser", r.URL.Query().Get("username"))
		}
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))
}

func TestGeonamesClient_FetchesCityNames(t *testing.T) {
	var hits int32
	up := newUpstream(t, http.StatusOK, sampleJSON, &hits)
	defer up.Close()

	c := NewGeonamesClient(up.URL, "testuser", time.Minute)
	got, err := c.CityNames(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"Shanghai", "Cairo", "Rio de Janeiro"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestGeonamesClient_CachesWithinTTL(t *testing.T) {
	var hits int32
	up := newUpstream(t, http.StatusOK, sampleJSON, &hits)
	defer up.Close()

	c := NewGeonamesClient(up.URL, "testuser", time.Minute)
	for i := 0; i < 3; i++ {
		if _, err := c.CityNames(context.Background()); err != nil {
			t.Fatalf("call %d: %v", i, err)
		}
	}
	if hits != 1 {
		t.Errorf("upstream hit %d times, want 1", hits)
	}
}

func TestGeonamesClient_RefetchesAfterTTL(t *testing.T) {
	var hits int32
	up := newUpstream(t, http.StatusOK, sampleJSON, &hits)
	defer up.Close()

	c := NewGeonamesClient(up.URL, "testuser", 10*time.Millisecond)
	_, _ = c.CityNames(context.Background())
	time.Sleep(20 * time.Millisecond)
	_, _ = c.CityNames(context.Background())
	if hits != 2 {
		t.Errorf("upstream hit %d times, want 2", hits)
	}
}

func TestGeonamesClient_UpstreamErrorIsReturnedNotCached(t *testing.T) {
	var hits int32
	up := newUpstream(t, http.StatusServiceUnavailable, `{"status":{"message":"down"}}`, &hits)
	defer up.Close()

	c := NewGeonamesClient(up.URL, "testuser", time.Minute)
	if _, err := c.CityNames(context.Background()); err == nil {
		t.Fatal("expected error on 503, got nil")
	}
	_, _ = c.CityNames(context.Background())
	if hits != 2 {
		t.Errorf("upstream hit %d times, want 2 (errors must not be cached)", hits)
	}
}

func captureLog(t *testing.T) *bytes.Buffer {
	t.Helper()
	var buf bytes.Buffer
	log.SetOutput(&buf)
	t.Cleanup(func() { log.SetOutput(io.Discard) })
	return &buf
}

func TestCityNames_LogsWhyItFetchedOrUsedCache(t *testing.T) {
	var hits int32
	up := newUpstream(t, http.StatusOK, sampleJSON, &hits)
	defer up.Close()

	buf := captureLog(t)
	c := NewGeonamesClient(up.URL, "testuser", 30*time.Millisecond)

	if _, err := c.CityNames(context.Background()); err != nil {
		t.Fatal(err)
	}
	if got := buf.String(); !strings.Contains(got, "cache empty") || !strings.Contains(got, "fetched 3 cities") {
		t.Errorf("first call should log an empty-cache fetch, got: %q", got)
	}

	buf.Reset()
	if _, err := c.CityNames(context.Background()); err != nil {
		t.Fatal(err)
	}
	if got := buf.String(); !strings.Contains(got, "cache hit") {
		t.Errorf("second call should log a cache hit, got: %q", got)
	}

	buf.Reset()
	time.Sleep(40 * time.Millisecond)
	if _, err := c.CityNames(context.Background()); err != nil {
		t.Fatal(err)
	}
	if got := buf.String(); !strings.Contains(got, "cache expired") || !strings.Contains(got, "fetched 3 cities") {
		t.Errorf("call after TTL should log an expiry-driven fetch, got: %q", got)
	}
}
