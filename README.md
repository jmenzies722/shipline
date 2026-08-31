# Shipline

A status board on your laptop. It checks a list of websites and shows
whether each one is up.

```
go test ./...
go run .
```

Then open http://127.0.0.1:8080

- `/` — the board (GitHub and example.com right now)
- `/healthz` — Shipline itself. Prints `ok` if this process is running.

The story of what we built, in order, is in [docs/walkthrough.md](docs/walkthrough.md).
Read that when you come back to this repo.
