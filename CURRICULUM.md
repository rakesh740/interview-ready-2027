# 18-Week Curriculum — Go + Kubernetes Interview Prep

## How this works (read this first, every session)

- User is a senior Go/backend engineer (6 YOE) — strong Go, near-zero Kubernetes going in. Full background: see `STORIES.md` for resume-derived context.
- Trigger phrases: **"quiz me"** or **"today's task"** — Claude reads this file for the current week/day, picks the right mode below, and runs it.
- **Monday**: skipped entirely. No quiz, no coding. Fixed rest day (in-office, tired) — not something to feel behind about.
- **Tuesday**: light only. A short 3-4 question quiz round. No coding push.
- **Wednesday–Friday**: full weekday pace. One quiz round (5-8 questions) + a small, scoped coding increment (~30-60 min).
- **Saturday–Sunday**: ~3-4 hr build blocks — the week's main project milestone gets built here.
- **Resume-story days** (marked per week below): instead of a topic quiz, grill the user to produce and tighten a STAR answer (situation/task/action/result, real numbers) for that week's resume bullet. Write the polished version into `STORIES.md`. Probe it with realistic interviewer follow-ups.
- After any quiz or story session, append weak spots / topics covered to `QUESTIONS.md`.
- No scheduled/automated nudges — this only runs when the user shows up and asks.
- If a week's stretch goal slips, it rolls into the next week's slack rather than blocking — Month 4 (Weeks 15-18) has deliberate buffer for this.

---

## The Project: a DAG-based Workflow Orchestration Platform on Kubernetes

What Weeks 1-4's simple job scheduler grows into — a mini Argo-Workflows/Temporal: submit a **workflow** (a DAG of dependent jobs, not just a single job), it gets scheduled across a fleet of workers via **Kafka**, runs on Kubernetes with autoscaling and full observability (**Elasticsearch**-backed logs/search), persists workflow/job state in **Postgres** (strong consistency, deep dive in Week 2) and execution history/output in **MongoDB** (flexible schema, deep dive Weeks 6-7) — a real polyglot-persistence story. It's managed through a **custom Kubernetes operator** (`kubectl apply -f workflow.yaml` triggers it), is multi-tenant with RBAC, and ships through a real CI/CD pipeline. Every piece is chosen to map onto a resume strength (Raft/consensus → leader election, MVCC → Postgres internals and workflow-state consistency, IAM background → multi-tenancy/RBAC, nsqd → Kafka) or close a real gap (Kubernetes end to end, plus the explicitly-requested MongoDB/Kafka/Elasticsearch/event-driven-architecture-in-depth/Postgres-in-depth additions).

Repo layout (see `README.md` for the up-to-date version as it fills in):
```
cmd/{api,dispatcher,worker}/
internal/{workflow,job,leaderelection}/
internal/storage/{postgres,mongo}/   internal/search/   internal/events/
proto/   operator/   deploy/{k8s,helm}/   .github/workflows/
```

---

## Month 1 (Weeks 1-4) — Foundation: core scheduler, DAG workflows, Postgres deep dive

**Target by end of Month 1**: a real scheduler is running inside the local `kind` cluster, backed by a Postgres schema you can defend in depth (indexing, isolation levels, query plans), and can execute a 3-step DAG workflow end to end, surviving a pod restart.

### Week 1 — Go senior-level depth + project skeleton
- [ ] **Quiz topics**: goroutine scheduler (GMP model), GC internals (tricolor mark-sweep, write barriers, GOGC/GOMEMLIMIT), memory model & happens-before, context cancellation patterns, generics in practice.
- [ ] **Resume-story day**: the TPS optimization story (2000→8000 TPS, Hyperledger Fabric) — bottleneck, what changed, how measured.
- [ ] **Weekday coding**: gRPC service skeleton (`proto/`, `cmd/api`), job submission endpoint, in-memory queue with a generic worker pool (`internal/job`).
- [ ] **Weekend build**: single-job path working end-to-end — submit via gRPC → queued → worker executes → status queryable. Graceful shutdown, structured logging.
- [ ] **Milestone**: `make run` starts api + worker locally; a CLI or `grpcurl` call submits and completes a job.

### Week 2 — Postgres in depth + reliability + Kubernetes fundamentals
- [ ] **Quiz topics**: **Postgres deep dive** — isolation levels (read committed vs. repeatable read vs. serializable) and how Postgres's own MVCC implementation works (tuple versions, xmin/xmax, vacuum) — explicit compare/contrast with the MVCC you built at Fabric; indexing strategies (B-tree vs. GIN/GIST, partial/covering indexes); reading `EXPLAIN ANALYZE` output and spotting a bad query plan; row-level locking & deadlocks; connection pooling (pgbouncer, why it matters under many short-lived connections). Plus: idempotency keys, outbox pattern, retry/backoff+jitter; K8s basics (pods, deployments, services, configmaps/secrets, namespaces).
- [ ] **Resume-story day**: MVCC/snapshot-management story, and the driver-order-route matching algorithm story (20% cost reduction).
- [ ] **Weekday coding**: `internal/storage/postgres` — schema design for workflow/job state, deliberate index choices (not defaults), retry logic with exponential backoff + dead-letter handling.
- [ ] **Weekend build**: Dockerize api + worker, first `deploy/k8s/` manifests, deploy to local `kind` cluster, verify via `kubectl port-forward`. Run `EXPLAIN ANALYZE` on the real job-lookup query and tune the index for it.
- [ ] **Milestone**: scheduler runs entirely inside `kind`; jobs survive a worker pod restart; can explain the schema's isolation level and index choices, not just that they exist.

### Week 3 — Kubernetes depth + early observability
- [ ] **Quiz topics**: liveness/readiness/startup probes, resource requests/limits, ConfigMap/Secret patterns, rolling updates, basic HPA mechanics.
- [ ] **Weekday coding**: Prometheus metrics via `client_golang` (queue depth, job latency), structured logs.
- [ ] **Weekend build**: probes tuned correctly, resource limits set deliberately (not copy-pasted), first Grafana dashboard (Helm-installed in-cluster) showing live metrics.
- [ ] **Milestone**: can watch a live dashboard while jobs run; a killed worker pod is caught by probes and replaced automatically.

### Week 4 — DAG workflows (the complexity pivot)
- [ ] **Quiz topics**: DAG scheduling / topological sort, workflow state machines, distributed state consistency (tie to MVCC/Postgres background), saga pattern vs. simple retries.
- [ ] **Weekday coding**: extend the job model into `internal/workflow` — a DAG of jobs with dependencies, a state machine (pending → running → succeeded/failed) per node, dependency resolution.
- [ ] **Weekend build**: a 3-step DAG workflow (A → B,C → D style fan-out/fan-in) runs end to end in the cluster; partial-failure behavior defined; verify concurrent node-status updates don't race (row-level locking from Week 2 pays off here).
- [ ] **Milestone**: `Month 1 done` — a real DAG workflow executes reliably inside `kind`, with metrics/dashboards and a Postgres schema you can defend in an interview.

---

## Month 2 (Weeks 5-10, 6 weeks) — Scale-out: event-driven architecture in depth + polyglot persistence + observability

**Target by end of Month 2**: the monolith is split into api / dispatcher / worker-fleet services communicating over **Kafka**, workflow execution history lives in **MongoDB** as a deliberate second datastore, logs/search run through **Elasticsearch**, the dispatcher is leader-elected, and the system autoscales under real load without losing workflows when pods die. This is the widest month — 6 weeks, not 4, to give Kafka, Mongo, and Elasticsearch real depth rather than token integrations.

### Week 5 — Event-driven architecture in depth + Kafka via Strimzi
- [ ] **Quiz topics**: **EDA in depth** — choreography vs. orchestration (which does this project use, and why), event notification vs. event-carried state transfer vs. event sourcing (which is this project doing, could it do event sourcing instead), schema evolution & compatibility (Avro/Protobuf + schema registry, backward vs. forward compatibility), exactly-once vs. at-least-once semantics in Kafka specifically (idempotent producers, consumer groups & offset management, "how would this differ from the nsqd you used at Kartbites").
- [ ] **Weekday coding**: split `cmd/api` (accepts workflows) from `cmd/dispatcher` (resolves DAG deps, emits job-ready events) from `cmd/worker` (consumes and executes) — wired via `internal/events` on top of **Kafka**, deployed in-cluster via **Strimzi** (a free Kafka operator for Kubernetes — also a preview of Month 3's "build your own operator" theme, since you'll use a real production-grade operator before writing one).
- [ ] **Weekend build**: full round-trip through Kafka works in-cluster; api, dispatcher, and workers are separately deployable/scalable; a topic partitioning scheme chosen deliberately (not default).
- [ ] **Milestone**: three independently deployable services, communicating only through Kafka, running in `kind` via Strimzi.

### Week 6 — MongoDB in depth (1/2): document design + indexing
- [ ] **Quiz topics**: document schema design (embedding vs. referencing), when a flexible schema genuinely helps vs. when it's a foot-gun, indexing in MongoDB (compound, multikey, text indexes) vs. Postgres indexing from Week 2 — direct compare/contrast.
- [ ] **Weekday coding**: `internal/storage/mongo` — a workflow **execution-history/output** store (each job's output is naturally document-shaped and varies by job type, unlike the fixed relational workflow/job state schema in Postgres).
- [ ] **Weekend build**: worker execution results write into Mongo; a query for "all outputs of workflow X" is indexed correctly, verified with `.explain()`.
- [ ] **Milestone**: a real second datastore, chosen for a real reason — can explain in an interview exactly why this data is in Mongo and the workflow state is in Postgres, not "because I wanted to use both."

### Week 7 — MongoDB in depth (2/2): consistency, aggregation, sharding concepts
- [ ] **Quiz topics**: write concern & read concern (and how MongoDB's tunable consistency contrasts with Postgres's isolation levels from Week 2 — a strong "how do you choose your datastore" interview answer), the aggregation pipeline, replica sets for HA, sharding concepts (even without running a real sharded cluster locally).
- [ ] **Weekday coding**: an aggregation pipeline building a "workflow history" view (e.g. success rate by job type over time) directly in Mongo.
- [ ] **Weekend build**: a small local replica set (`docker compose` or in-`kind`) to see failover behavior firsthand; document the polyglot-persistence decision (what's in Postgres, what's in Mongo, why) in `README.md`.
- [ ] **Milestone**: polyglot persistence is real and defensible — two datastores, each earning its place.

### Week 8 — Observability depth: Elasticsearch + distributed tracing
- [ ] **Quiz topics**: distributed tracing concepts (spans, context propagation across Kafka — the hard part, since there's no HTTP header to piggyback on), the three pillars (metrics/logs/traces) and when each earns its cost, Elasticsearch fundamentals (inverted index, analyzers/tokenizers, why it's fast at text search), cardinality problems in metrics, SLI vs. SLO.
- [ ] **Weekday coding**: OpenTelemetry tracing across submit → dispatch → execute, propagated through Kafka message headers; ship structured logs to **Elasticsearch** (Helm, in-cluster) with trace-ID correlation.
- [ ] **Weekend build**: Kibana (or a raw ES query) answering "find every failed job for tenant X in the last hour" fast, by index design, not brute force; Tempo/Jaeger showing a full trace of one workflow across all three services.
- [ ] **Milestone**: can debug a synthetic "why is this workflow slow / where did it fail" scenario purely from search + traces, live.

### Week 9 — Resilience + leader election
- [ ] **Quiz topics**: circuit breakers, rate limiting algorithms (token bucket, sliding window), graceful degradation, leader election mechanics (tie explicitly to the Raft/etcd background from CoffeeBeans — how does K8s Lease-based election differ from Raft leader election).
- [ ] **Weekday coding**: leader election for the dispatcher (`internal/leaderelection`, via the K8s Lease API — client-go's `leaderelection` package), rate limiting on the api service, circuit breaker around worker calls.
- [ ] **Weekend build**: chaos test — kill the leader dispatcher pod mid-workflow, verify a standby takes over with no lost jobs; kill a worker mid-job, verify recovery via retry.
- [ ] **Milestone**: the system survives an unannounced pod kill at any tier without losing a workflow.

### Week 10 — Autoscaling + load testing
- [ ] **Quiz topics**: HPA mechanics with custom metrics, Prometheus Adapter, load-testing methodology (open vs. closed load models), capacity planning basics.
- [ ] **Weekday coding**: HPA on the worker deployment driven by Kafka consumer-lag (via Prometheus Adapter, not just CPU), tuning scale-up/down behavior.
- [ ] **Weekend build**: load test with `k6` or `vegeta` (free/local) generating a burst of workflows; watch workers autoscale live via `kubectl get hpa -w` and the Grafana dashboard.
- [ ] **Milestone**: `Month 2 done` — a load spike triggers visible, correct autoscaling; Kafka, Mongo, Elasticsearch, and Postgres are all live, load-bearing, and each individually defensible in an interview.

---

## Month 3 (Weeks 11-14) — Kubernetes-native: operator + multi-tenancy + security

**Target by end of Month 3**: the system is driven by a custom `Workflow` CRD and controller (not a direct API call), is multi-tenant with real RBAC/network isolation, and has baseline security hardening (secrets, mTLS) in place.

### Week 11 — Operator foundations
- [ ] **Quiz topics**: the Kubernetes controller pattern, reconciliation loops (level-triggered vs. edge-triggered), informers/listers/workqueues, CRDs vs. built-in resources — when to build one (contrast with using Strimzi's CRD from Week 5 — now you write one instead of consuming one).
- [ ] **Weekday coding**: `kubebuilder`-scaffold a `Workflow` CRD in `operator/`, basic controller that watches `Workflow` custom resources and creates the corresponding internal workflow record.
- [ ] **Weekend build**: `kubectl apply -f examples/three-step-workflow.yaml` triggers a real run through the existing api/dispatcher/worker pipeline.
- [ ] **Milestone**: the operator is a real, working entry point to the system — not a toy.

### Week 12 — Controller depth
- [ ] **Quiz topics**: status subresources, finalizers, owner references & garbage collection, admission webhooks (validating vs. mutating), controller idempotency.
- [ ] **Resume-story day**: the IAM onboarding-automation story — draw the explicit parallel to the reconciliation loop just built (approval/rejection/environment-generation as a desired-state-vs-actual-state problem).
- [ ] **Weekday coding**: `Workflow.status` subresource reflecting live progress, a finalizer for cleanup, a validating webhook rejecting malformed DAGs (e.g. cycles) before they're accepted.
- [ ] **Weekend build**: end-to-end demo — a bad DAG (with a cycle) is rejected at `kubectl apply` time by the webhook, not at runtime.
- [ ] **Milestone**: the operator behaves like a production-grade K8s extension: validated input, observable status, clean teardown.

### Week 13 — Multi-tenancy & RBAC
- [ ] **Quiz topics**: Kubernetes RBAC model (Roles/ClusterRoles/Bindings), namespace-based isolation strategies, NetworkPolicies, tenant-isolation tradeoffs at scale.
- [ ] **Weekday coding**: namespace-per-tenant (or label-scoped) isolation for workflows, RBAC roles restricting what each tenant's service account can touch, NetworkPolicies restricting cross-tenant traffic.
- [ ] **Weekend build**: two synthetic tenants running workflows in the same cluster, provably isolated — tenant A cannot see or affect tenant B's workflows via RBAC or network policy.
- [ ] **Milestone**: a real multi-tenancy story to tell, backed by working RBAC + NetworkPolicy, not just namespaces as decoration.

### Week 14 — Security hardening
- [ ] **Quiz topics**: Pod Security Standards, secrets management approaches (Sealed Secrets/SOPS vs. raw K8s Secrets), mTLS basics, supply-chain basics (image signing/scanning, at a conceptual level).
- [ ] **Weekday coding**: Pod Security Standards enforced (restricted profile), secrets moved off raw plaintext manifests (Sealed Secrets, free/local), baseline mTLS between services (via a lightweight mesh like Linkerd, free, or manual TLS if the mesh proves too heavy for the time budget).
- [ ] **Weekend build**: a security pass — run through a checklist (least-privilege RBAC, no privileged containers, secrets encrypted at rest in git, mTLS verified with `kubectl exec` + curl from inside a pod).
- [ ] **Milestone**: `Month 3 done` — the system is genuinely Kubernetes-native (CRD-driven), multi-tenant, and passes a basic security review.

---

## Month 4 (Weeks 15-18) — Production hardening + interview readiness

**Target by end of Month 4**: the system ships through a real CI/CD pipeline, has SLOs/alerting and demonstrated chaos-resilience, and — most importantly — every design decision and every resume bullet can be explained cold in a mock interview.

### Week 15 — CI/CD + GitOps
- [ ] **Quiz topics**: CI/CD pipeline design, GitOps principles (declarative, git as source of truth), Helm chart design, image tagging/versioning strategies.
- [ ] **Weekday coding**: GitHub Actions pipeline (build, test, lint, push image to `ghcr.io` — free), a proper Helm chart in `deploy/helm/` replacing raw manifests (now covering api/dispatcher/worker + Strimzi/Mongo/Elasticsearch dependencies).
- [ ] **Weekend build**: install ArgoCD in-cluster (free/local), point it at the repo, verify a change pushed to `main` rolls out to `kind` automatically.
- [ ] **Milestone**: a real GitOps loop — `git push` to deploy, no manual `kubectl apply`.

### Week 16 — SLOs, alerting & chaos
- [ ] **Quiz topics**: SLIs/SLOs/error budgets, alerting design (avoiding alert fatigue), chaos engineering principles, failure-mode-and-effects thinking.
- [ ] **Weekday coding**: Alertmanager rules for real SLOs (e.g. "95% of workflows complete within N seconds"), a couple of Chaos Mesh (free) experiments — pod-kill and network-delay, including killing a Kafka broker and a Mongo node.
- [ ] **Weekend build**: run the chaos experiments live, confirm alerts fire correctly and the system self-heals within its SLO.
- [ ] **Milestone**: can show — with real dashboards and alert history — that the system meets a stated SLO under injected failure across every datastore and service.

### Week 17 — Mock interviews
- [ ] **Quiz topics**: none — full system-design mock sessions using this project as the concrete answer to prompts like "design a distributed workflow orchestrator," "design a multi-tenant SaaS scheduling platform," "design a job queue at scale," "when would you choose Postgres vs. MongoDB vs. Elasticsearch for a given access pattern."
- [ ] **Resume-story day**: final pass over all of `STORIES.md` — every story drilled at least twice, cold, out loud.
- [ ] **Weekday + weekend**: 3-4 full mock interviews (mix of system design and behavioral), reviewing `QUESTIONS.md` for any topic that's come up as a repeated weak spot across the whole 18 weeks.
- [ ] **Milestone**: can walk through the whole system's design, unprompted, in under 10 minutes, then go arbitrarily deep on any component or datastore choice.

### Week 18 — Final polish (buffer week)
- [ ] **Weekday + weekend**: catch up anything that slipped from Weeks 1-17; polish `README.md` into a portfolio piece (architecture diagram, key tradeoffs, what you'd do differently at 10x scale, why each datastore/tool was chosen); final review of `STORIES.md` and `QUESTIONS.md`.
- [ ] **Milestone**: `Month 4 done` / project complete — a portfolio-grade, multi-tenant, event-driven, polyglot-persistence, Kubernetes-native workflow orchestration platform with CI/CD, observability, security, and chaos-tested resilience, plus every resume bullet backed by a tight, rehearsed story.

---

## Day 0 — Environment setup (done)

- [x] Docker Desktop running.
- [x] `kind` installed (`brew install kind`).
- [x] Local cluster created: `kind create cluster --name interview-ready`.
- [x] Verified: `kubectl cluster-info --context kind-interview-ready`.
