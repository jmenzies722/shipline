package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync"
	"time"
)

type server struct {
	targets []string
	client  *http.Client
	page    *template.Template
	every   time.Duration

	mu      sync.Mutex
	results []result
}

type result struct {
	URL    string
	Up     bool
	Status int
	Took   time.Duration
	Err    string
}

func main() {
	s := &server{
		targets: []string{
			"https://github.com",
			"https://example.com",
		},
		client: &http.Client{Timeout: 5 * time.Second},
		page:   template.Must(template.New("index").Parse(indexHTML)),
		every:  30 * time.Second,
	}
	s.refresh()
	go s.loop()

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      s.mux(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 20 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	log.Printf("listening on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func (s *server) mux() *http.ServeMux {
	m := http.NewServeMux()
	m.HandleFunc("/healthz", s.healthz)
	m.HandleFunc("/", s.index)
	return m
}

func (s *server) healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintln(w, "ok")
}

func (s *server) index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := s.page.Execute(w, s.snapshot()); err != nil {
		log.Printf("template: %v", err)
	}
}

func (s *server) loop() {
	ticker := time.NewTicker(s.every)
	defer ticker.Stop()
	for range ticker.C {
		s.refresh()
	}
}

func (s *server) refresh() {
	results := make([]result, 0, len(s.targets))
	for _, url := range s.targets {
		results = append(results, s.probe(url))
	}
	s.mu.Lock()
	s.results = results
	s.mu.Unlock()
}

func (s *server) snapshot() []result {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]result, len(s.results))
	copy(out, s.results)
	return out
}

func (s *server) probe(url string) result {
	start := time.Now()
	resp, err := s.client.Get(url)
	took := time.Since(start)
	if err != nil {
		return result{URL: url, Took: took, Err: err.Error()}
	}
	defer resp.Body.Close()

	return result{
		URL:    url,
		Up:     resp.StatusCode < 400,
		Status: resp.StatusCode,
		Took:   took,
	}
}

const indexHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Shipline</title>
<style>
  :root { color-scheme: dark; }
  body { font: 16px/1.4 ui-sans-serif, system-ui, sans-serif; margin: 2rem; background: #111; color: #eee; }
  h1 { font-size: 1.25rem; font-weight: 600; }
  p { color: #9a9a9a; }
  table { width: 100%; max-width: 40rem; border-collapse: collapse; }
  th, td { text-align: left; padding: 0.6rem 0.4rem; border-bottom: 1px solid #2a2a2a; }
  .up { color: #3dd68c; }
  .down { color: #ff6b6b; }
</style>
</head>
<body>
<h1>Shipline</h1>
<p>Things you ship, and whether they are up.</p>
<table>
  <tr><th>Target</th><th>Status</th><th>Time</th></tr>
  {{range .}}
  <tr>
    <td>{{.URL}}</td>
    <td class="{{if .Up}}up{{else}}down{{end}}">
      {{if .Up}}Up {{.Status}}{{else if .Err}}Down — {{.Err}}{{else}}Down {{.Status}}{{end}}
    </td>
    <td>{{.Took}}</td>
  </tr>
  {{end}}
</table>
</body>
</html>
`
