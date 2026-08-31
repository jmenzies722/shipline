package main

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func testServer(targets []string) *server {
	return &server{
		targets: targets,
		client:  &http.Client{Timeout: time.Second},
		page:    template.Must(template.New("index").Parse(indexHTML)),
	}
}

func TestHealthz(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	testServer(nil).mux().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status %d", rec.Code)
	}
	if body := rec.Body.String(); body != "ok\n" {
		t.Fatalf("body %q", body)
	}
}

func TestIndexReportsUpstream(t *testing.T) {
	up := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(up.Close)

	s := testServer([]string{up.URL})
	s.refresh()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	s.mux().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status %d", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, up.URL) {
		t.Fatalf("missing target url in %s", body)
	}
	if !strings.Contains(body, "Up 200") {
		t.Fatalf("missing up status in %s", body)
	}
}

func TestIndexUsesLastCheck(t *testing.T) {
	s := testServer([]string{"http://127.0.0.1:1"})
	s.results = []result{{URL: "http://cached.example", Up: true, Status: 200}}

	rec := httptest.NewRecorder()
	s.mux().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	body := rec.Body.String()
	if !strings.Contains(body, "http://cached.example") {
		t.Fatalf("page should show the last check, got %s", body)
	}
	if strings.Contains(body, "127.0.0.1") {
		t.Fatalf("page probed again instead of using memory: %s", body)
	}
}

func TestUnknownPathIsNotFound(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/favicon.ico", nil)
	rec := httptest.NewRecorder()
	testServer([]string{"http://127.0.0.1:1"}).mux().ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status %d", rec.Code)
	}
}

func TestIndexStillOKWhenTargetDown(t *testing.T) {
	down := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	t.Cleanup(down.Close)

	s := testServer([]string{down.URL})
	s.refresh()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	s.mux().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status page should stay 200 when a target is down, got %d", rec.Code)
	}
	if body := rec.Body.String(); !strings.Contains(body, "Down 503") {
		t.Fatalf("missing down status in %s", body)
	}
}
