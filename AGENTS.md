# Shipline

A personal status board: URLs you actually care about, checked by one Go
process, shown on one page. Real use — you open it to see if something is
down. The DevOps path is the other half of the project: we take that same
process from this laptop to a public URL on AWS.

Chaos Gym stays parked. This repo does not reuse that cluster, those
manifests, or that instance.

## What it is

You give Shipline a list of URLs (hardcoded in phase 1). It probes them
every 30 seconds, remembers the last answers, and renders a status page.
`/healthz` answers for Shipline itself — the page can show a target down
and still be a healthy process. That distinction is the whole product.

## Phase 1 (current)

A Go process Josh can run and test on this machine. He is learning Go as we
build — explain each new language construct the first time it appears, then
keep moving. No Docker, no AWS, no CI file until the process exists and he
can say what it is doing.

**Done looks like:** a real `*_test.go` fails if `/healthz` is missing;
`go test` and `go run .` both work; `/` shows probe results for the
configured targets; Josh can explain `package main`, `func main`, a TCP
port, and why `/healthz` is not the status page.

## Later phases (not started)

2. Docker — same process, now in an image
3. GitHub Actions CI — tests on every push, not just this Mac
4. AWS budget + VPC + IAM — nothing billable before a budget alarm.
   Deploy auth is GitHub OIDC (no long-lived AWS keys in the repo).
5. Deploy — process runs on AWS
6. CD — a merge to main updates the live URL
7. Persist targets (SQLite, then Postgres) and check in the background
8. Signals — logs, then metrics
9. Rollback — a bad deploy undone on purpose

## Stack

Go 1.26, html/template (one binary, no separate frontend), Docker, GitHub
Actions with OIDC to AWS, Terraform. No Kubernetes until this path exists —
k8s is a scheduler for processes you already know how to ship.

## Never

- Never write a complete file for a piece Josh hasn't reached yet.
- Never create any billable AWS resource before a budget alarm exists.
- Never reuse the Chaos Gym instance, VPC, or manifests.
- Never add Kubernetes, a service mesh, or a second cloud account to skip
  the boring path.
- Never add a separate Node/React app in phase 1 — that is a second
  deployable and hides the pipeline.
- Never commit secrets or real credentials.
- Never run `terraform destroy` or terminate instances without asking.
- Never write the phase post. That reflection is Josh's.

## How to work with me

- I'm learning Go and DevOps in the same repo. Explain the why before the
  how. Simple and concrete — one idea at a time, no jargon without a
  plain sentence next to it. When a Go construct shows up for the first
  time, say what it is in one or two sentences, then use it.
- Step by step, gated. State the step, why it exists, and what files
  change. Wait for Josh to approve. Do not write or run the step until
  he says go.
- Smallest piece that moves the phase forward, then stop.
- Name tradeoffs when there's a real choice, pick one, keep moving.
- Don't quiz me after every piece. Questions go in QUESTIONS.md and get
  reviewed at the end of a phase. Ask inline only when my answer changes
  what gets built next.
- Don't skip networking, IAM, DNS, or state files when we reach them.
- Never write the phase post for me.
- After every approved step, update `docs/walkthrough.md` and `README.md`
  in plain language so Josh can reread the project without the chat.

## Do not

- Commit secrets or write real credentials to any file
- Run destructive commands without asking
- Add dependencies without telling me what they cost
