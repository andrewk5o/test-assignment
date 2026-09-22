package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type stubSource struct {
	names []string
	err   error
}

func (s stubSource) CityNames(context.Context) ([]string, error) { return s.names, s.err }

func doCount(t *testing.T, src CitySource, query string) (*httptest.ResponseRecorder, map[string]any) {
	t.Helper()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/cities/count"+query, nil)
	NewServer(src).ServeHTTP(rec, req)

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("response is not JSON: %v\n%s", err, rec.Body.String())
	}
	return rec, body
}

func TestCountHandler_ReturnsCountAndMatches(t *testing.T) {
	src := stubSource{names: []string{"Shanghai", "Cairo", "Chongqing", "Chengdu", "Rio de Janeiro"}}
	rec, body := doCount(t, src, "?letter=c")

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if body["letter"] != "c" || body["count"] != float64(3) {
		t.Errorf("body = %v", body)
	}
	want := []any{"Cairo", "Chongqing", "Chengdu"}
	if !reflect.DeepEqual(body["cities"], want) {
		t.Errorf("cities = %v, want %v", body["cities"], want)
	}
}

func TestCountHandler_RejectsBadLetter(t *testing.T) {
	src := stubSource{names: []string{"Cairo"}}
	for _, q := range []string{"", "?letter=", "?letter=ab", "?letter=1", "?letter=%20"} {
		rec, body := doCount(t, src, q)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("%q: status = %d, want 400", q, rec.Code)
		}
		if body["error"] == nil {
			t.Errorf("%q: missing error message", q)
		}
	}
}

func TestCountHandler_UpstreamFailureIs502(t *testing.T) {
	src := stubSource{err: errors.New("boom")}
	rec, body := doCount(t, src, "?letter=c")
	if rec.Code != http.StatusBadGateway {
		t.Errorf("status = %d, want 502", rec.Code)
	}
	if body["error"] == nil {
		t.Error("missing error message")
	}
}

func TestCountHandler_MethodNotAllowed(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/cities/count?letter=c", nil)
	NewServer(stubSource{}).ServeHTTP(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("status = %d, want 405", rec.Code)
	}
}
