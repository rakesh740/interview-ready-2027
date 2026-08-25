# interview-ready-2027

A DAG-based **Workflow Orchestration Platform** — event-driven, polyglot-persistence, multi-tenant, Kubernetes-native — built incrementally over an 18-week interview-prep plan (see [`CURRICULUM.md`](./CURRICULUM.md)).

Submit a workflow (a DAG of dependent jobs, not just a single job) and it's scheduled across a fleet of workers via **Kafka**, runs on Kubernetes with metrics-driven autoscaling and full distributed tracing, persists workflow/job state in **Postgres** (strong consistency) and execution history/output in **MongoDB** (flexible schema), searches logs and history through **Elasticsearch**, is managed through a custom Kubernetes operator (`kubectl apply -f workflow.yaml`), is multi-tenant with RBAC and network isolation, and ships through a real CI/CD + GitOps pipeline.

## Architecture (fills in as each phase lands)

```
cmd/
  api/           control-plane: accepts workflow submissions
  dispatcher/    resolves DAG dependencies, emits job-ready events (leader-elected)
  worker/        consumes job-ready events, executes jobs
internal/
  workflow/      DAG model + per-node state machine
  job/           job domain types
  storage/
    postgres/    workflow/job state — strong consistency, deep-dive schema (Week 2)
    mongo/       execution history/output — flexible schema (Weeks 6-7)
  search/        Elasticsearch client — log/history search (Week 8)
  events/        Kafka client, deployed via Strimzi (Week 5)
  leaderelection/  dispatcher leader election (K8s Lease API)
proto/           gRPC service definitions
operator/        kubebuilder-scaffolded Workflow CRD + controller
deploy/
  k8s/           raw manifests (Month 1-2)
  helm/          Helm chart (Month 4)
.github/workflows/  CI pipeline
```

## Why two datastores (polyglot persistence)

- **Postgres**: workflow/job state — needs strong consistency and row-level locking for concurrent DAG-node updates. Deep-dive week (Week 2) covers isolation levels, indexing, and how Postgres's own MVCC compares to the MVCC built at CoffeeBeans/Hyperledger Fabric.
- **MongoDB**: execution history/output — each job's output is naturally document-shaped and varies by job type; read-heavy, schema-flexible. Deep-dive weeks (6-7) cover document design, indexing, aggregation pipelines, and tunable consistency (write/read concern) vs. Postgres's isolation levels.
- **Elasticsearch**: fast search over logs and execution history (Week 8) — answers "find every failed job for tenant X in the last hour" by index design, not brute force.

## Running locally

```
make run
```

(Details filled in as `cmd/api`, `cmd/dispatcher`, and `cmd/worker` are implemented — see Week 1 in `CURRICULUM.md`.)

## Roadmap

- **Month 1 (Weeks 1-4)**: core scheduler → DAG workflows running in `kind`, Postgres schema defensible in depth.
- **Month 2 (Weeks 5-10, 6 weeks)**: event-driven microservices over Kafka (via Strimzi), MongoDB as a second datastore, Elasticsearch-backed observability, leader election, autoscaling under load.
- **Month 3 (Weeks 11-14)**: custom Kubernetes operator/CRD, multi-tenancy + RBAC, security hardening.
- **Month 4 (Weeks 15-18)**: CI/CD + GitOps, SLOs/alerting/chaos testing, interview polish.

## Project docs

- [`CURRICULUM.md`](./CURRICULUM.md) — the 18-week learning plan, checkbox-tracked, week by week.
- [`QUESTIONS.md`](./QUESTIONS.md) — running log of interview quiz topics covered and weak spots to revisit.
- [`STORIES.md`](./STORIES.md) — polished STAR-format answers for resume talking points.
