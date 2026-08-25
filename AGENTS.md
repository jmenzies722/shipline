# Shipline

A tiny Go HTTP service we take from `go run` on a laptop to a public URL on
AWS, through Docker, GitHub Actions, and Terraform. Solo learning project —
the point is that Josh can defend every layer in an interview, not that the
service is a product.

Chaos Gym stays parked. This repo does not reuse that cluster, those
manifests, or that instance.

## Phase 1 (current)

A Go process Josh can run and test on this machine. He is learning Go as we
build — explain each new language construct the first time it appears, then
keep moving. No Docker, no AWS, no CI file until the process exists and he
can say what it is doing.

**Done looks like:** `go test` and `go run .` both work; the process listens
on a port and answers `/healthz`; Josh can explain what `package main`,
`func main`, and a TCP port are in his own words.

## Later phases (not started)

2. Docker — same process, now in an image
3. GitHub Actions CI — the tests run on every push, not just on this Mac
4. AWS budget + VPC + IAM — nothing billable before a budget alarm
5. Deploy — process runs on AWS
6. CD — a merge to main updates the live URL
7. Signals — logs, then metrics, so a failure is visible
8. Rollback — a bad deploy can be undone on purpose

## Stack

Go 1.26, then Docker, GitHub Actions, Terraform, AWS. No Kubernetes until
the path above exists — k8s is a scheduler for processes you already know
how to ship.

## Never

- Never write a complete file for a piece Josh hasn't reached yet.
- Never create any billable AWS resource before a budget alarm exists.
- Never reuse the Chaos Gym instance, VPC, or manifests.
- Never add Kubernetes, a service mesh, or a second cloud account to skip
  the boring path.
- Never commit secrets or real credentials.
- Never run `terraform destroy` or terminate instances without asking.
- Never write the phase post. That reflection is Josh's.

## How to work with me

- I'm learning Go and DevOps in the same repo. Explain the why before the
  how. When a Go construct shows up for the first time, say what it is in
  one or two sentences, then use it.
- Smallest piece that moves the phase forward, then stop.
- Name tradeoffs when there's a real choice, pick one, keep moving.
- Don't quiz me after every piece. Questions go in QUESTIONS.md and get
  reviewed at the end of a phase. Ask inline only when my answer changes
  what gets built next.
- Don't skip networking, IAM, DNS, or state files when we reach them.
- Never write the phase post for me.

## Do not

- Commit secrets or write real credentials to any file
- Run destructive commands without asking
- Add dependencies without telling me what they cost
