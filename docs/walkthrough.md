# Walkthrough

Plain notes so you can reopen this repo and remember what we did.
We add a section after each step.

## Step 1 — A process on this laptop

Shipline is a Go program. `package main` means “this is a program.”
`func main` is where it starts. If `main` finishes, the program quits.

It listens on port 8080 — a numbered door on this Mac. Your browser
talks to Shipline, not to GitHub.

Two doors, two jobs:

- `/healthz` → `ok` means Shipline is alive
- `/` → the board. A site can be down and Shipline is still fine.

The first version checked GitHub and example.com **every time you
refreshed**. That worked. It also meant two tabs hit those sites twice.

## Step 2 — Check on a clock, show from memory

Shipline now checks the list every 30 seconds and remembers the last
answers. Opening `/` only reads that memory and draws the table.

On startup it checks once right away, so the first page load is not
empty. Then a background loop waits 30 seconds and checks again.

Two pieces of Go that showed up here:

- `go s.loop()` — start the loop **and keep going**. The server can
  listen while the loop checks sites. That is a goroutine: extra work
  that runs at the same time.
- `sync.Mutex` — a lock. The loop writes the results. The page reads
  them. The lock means they take turns so they do not trip over the
  same list.

You must restart `go run .` to pick this up. An old process is still
the old code.
