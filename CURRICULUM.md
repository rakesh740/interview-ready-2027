# 4-Week Curriculum — Go + Kubernetes Interview Prep

## How this works (read this first, every session)

- User is a senior Go/backend engineer (6 YOE) — strong Go, near-zero Kubernetes going in. Full background: see `STORIES.md` for resume-derived context.
- Trigger phrases: **"quiz me"** or **"today's task"** — Claude reads this file for the current week/day, picks the right mode below, and runs it.
- **Monday**: skipped entirely. No quiz, no coding. Fixed rest day (in-office, tired) — not something to feel behind about.
- **Tuesday**: light only. A short 3-4 question quiz round. No coding push.
- **Wednesday–Friday**: full weekday pace. One quiz round (5-8 questions) + a small, scoped coding increment (~30-60 min).
- **Saturday–Sunday**: ~3-4 hr build blocks — the week's main project milestone gets built here.
- **Resume-story days** (marked per week below): instead of a topic quiz, grill the user to produce and tighten a STAR answer (situation/task/action/result, real numbers) for that week's resume bullet. Write the polished version into `STORIES.md`. Probe it with realistic interviewer follow-ups (e.g. "why did payload size affect bandwidth specifically", "what would've happened without MVCC there").
- After any quiz or story session, append weak spots / topics covered to `QUESTIONS.md`.
- No scheduled/automated nudges — this only runs when the user shows up and asks.

---

## Week 1 — Go senior-level depth + project skeleton

- [ ] **Quiz topics** (Tue light / Wed-Fri full): goroutine scheduler (GMP model), GC internals (tricolor mark-sweep, write barriers, GOGC/GOMEMLIMIT), memory model & happens-before, context cancellation patterns, generics in practice.
- [ ] **Resume-story day** (any Wed-Fri): the TPS optimization story — 2000→8000 TPS via concurrency tuning + payload reduction on Hyperledger Fabric. What was the bottleneck, what specifically changed, how was the improvement measured.
- [ ] **Weekday coding** (Wed-Fri): gRPC service skeleton (`proto/`, `cmd/server`), job submission endpoint, in-memory queue with a generic worker pool (`internal/job`, `internal/scheduler`).
- [ ] **Weekend build**: working end-to-end scheduler — submit job via gRPC → queued → worker picks up via channel-based pool → executes → status queryable. Graceful shutdown, structured logging.
- [ ] **Milestone**: `make run` starts server + workers locally; a CLI or `grpcurl` call submits and completes a job.

## Week 2 — Persistence, reliability + Kubernetes fundamentals

- [ ] **Quiz topics**: idempotency keys, outbox pattern, retry/backoff+jitter, at-least-once vs exactly-once delivery (include "how would this differ with Kafka vs. nsqd" framing); K8s basics — pods, deployments, services, configmaps/secrets, namespaces.
- [ ] **Resume-story day**: the MVCC/snapshot-management story (what consistency problem it solved, what durability guarantees it gave) and the driver-order-route matching algorithm story (the 20% cost reduction — what the matching logic actually did).
- [ ] **Weekday coding**: swap in-memory storage for SQLite/Postgres (`internal/storage`), add retry logic with exponential backoff + dead-letter handling.
- [ ] **Weekend build**: Dockerize server + worker, write first `deploy/k8s/` manifests (Deployment, Service, ConfigMap), `kind load docker-image`, deploy to local cluster, verify via `kubectl port-forward`.
- [ ] **Milestone**: scheduler runs entirely inside the local `kind` cluster; jobs survive a worker pod restart.

## Week 3 — Kubernetes depth + observability

- [ ] **Quiz topics**: liveness/readiness/startup probes, resource requests/limits, HPA mechanics, StatefulSet vs Deployment, rolling updates, ConfigMap/Secret patterns, Helm basics.
- [ ] **Weekday coding**: Prometheus metrics via `client_golang` (queue depth, job latency, worker throughput), structured logs, basic OpenTelemetry tracing across submit → dispatch → execute.
- [ ] **Weekend build**: add HPA scaling workers on queue-depth/CPU, probes tuned correctly, load-test with many concurrent jobs, observe scaling live via `kubectl get hpa -w`.
- [ ] **Milestone**: can demo autoscaling under load with metrics to show it, live in the cluster.

## Week 4 — Advanced Kubernetes + interview polish

- [ ] **Quiz topics**: CRDs & the controller/operator pattern, client-go informers, reconciliation loops, consensus trade-offs at depth (Raft leader election/log replication edge cases — depth quizzing, not re-teaching), rate limiting algorithms (token bucket).
- [ ] **Resume-story day**: the IAM onboarding-automation story (approval/rejection/environment-generation workflow) — explicitly draw the parallel to the CRD/controller reconciliation loop just built.
- [ ] **Weekday coding** (stretch): a minimal CRD + controller (e.g. `JobPolicy` CRD) using `kubebuilder` or raw `client-go`, or a solid write-up of the pattern if time-constrained.
- [ ] **Weekend**: 1-2 full mock system-design sessions ("design a distributed job scheduler at scale" — using this project as the concrete answer), review `QUESTIONS.md` for recurring weak spots and `STORIES.md` for delivery, polish README with architecture diagram and tradeoff explanations.
- [ ] **Milestone**: project + README function as a portfolio piece; can walk through every design decision and every resume bullet out loud, cold.

---

## Day 0 — Environment setup (do before Week 1 build)

- [ ] Start Docker Desktop (`open -a Docker`, wait for daemon).
- [ ] Install `kind` (`brew install kind`) — lightest, fastest local cluster for active iteration, fully free/local.
- [ ] `kind create cluster --name interview-ready`
- [ ] Verify: `kubectl cluster-info --context kind-interview-ready`
