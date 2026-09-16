# Handoff — interview-ready-2027

_Written 2026-09-16. Read this at the start of a new session before doing anything else, then delete or ignore once its content is stale._

## What this project is

An 18-week Go/Kubernetes interview-prep curriculum, tracked in `CURRICULUM.md`, that also incrementally builds a DAG workflow orchestration platform (see `README.md`). Trigger phrases: **"quiz me"** / **"today's task"** — read `CURRICULUM.md` for the current week/day. After any quiz/story session, append weak spots to `QUESTIONS.md`.

## Critical working agreement — read this first

**The user writes all implementation code themselves from here forward.** Claude quizzes, explains, and reviews — it does not write feature code. Full rationale in Claude's persistent memory: `~/.claude/projects/-Users-rakeshrakshit/memory/feedback_interview_ready_user_codes.md` (a fresh agent should read that too). One exception already closed: Week 1 code (commit `78c09d2`) stays as Claude-authored; everything after is user-written.

## State of the repo

- Working tree has one uncommitted change at the time of writing: this handoff update plus a new `QUESTIONS.md` entry from this session — not yet committed. Everything else committed through `c1cee99`.
- Week 1 milestone: commit `78c09d2` — gRPC `JobService`, `internal/job` (Store, generic `Pool[T]`), `cmd/api`. `make run` works, verified via `grpcurl`.

## Curriculum progress

- **Week 1 resume-story day — still NOT done.** TPS-optimization story (2000→8000 TPS, Hyperledger Fabric) in `STORIES.md` is still blank. Tied to a real interview gap (see backlog below).
- **Week 2 quiz — indexing + `EXPLAIN ANALYZE` covered this session** (B-tree fit/leftmost-prefix rule, seq-scan-vs-index-scan cost model, GIN vs. GiST structural difference, `EXPLAIN ANALYZE` row-estimate mismatch → nested-loop blowup). Full detail in `QUESTIONS.md` under `## 2026-09-16 — Week 2: indexing + EXPLAIN ANALYZE` — don't re-derive it here, read that entry. Most important thing to carry forward:
  - **The "surface-level, no why" pattern is now confirmed across three distinct topics** (Fabric/Postgres MVCC compare-contrast on 2026-09-09; random-I/O cost and EXPLAIN ANALYZE row-mismatch on 2026-09-16). Treat this as a standing habit to interrupt on **every** quiz question going forward, not a topic-specific gap — when an answer names the right conclusion but not the mechanism, push at least once for the physical/structural reason before moving on.
  - **New redrill items:** B-tree/seq-scan cost mechanism (heap-page amortization) and GIN-vs-GiST structure (inverted index vs. lossy bounding-predicate tree) both needed full explanations given directly, not derived — redrill from scratch next session, don't assume retention.
  - User explicitly asked to revisit `EXPLAIN ANALYZE` again later even though it went reasonably well this pass — don't treat it as closed.
- **Week 2 topics not yet touched:** isolation levels beyond RC/RR (serializable/predicate locks/serialization anomalies), row-locking/deadlocks, connection pooling, idempotency keys/outbox/retry-backoff, K8s basics.

## Mock-interview transcript backlog

Full extraction (~55 weak spots) lives in `QUESTIONS.md` under `## 2026-08-30 — Backlog: review of 19 past mock-interview transcripts` — do not re-read the raw transcripts (`~/Projects/meetings/reviewed*_clean_text.txt`), that entry is the durable record. Two flagged as highest-priority (recurring across multiple interviews, still undrilled):
1. **GOMAXPROCS reasoning** — can't justify a concurrency number from first principles.
2. **Explaining own project's architecture/data-model rationale under adversarial questioning** — same gap this session's Fabric compare-contrast just re-confirmed.

## Options for next session

No fixed decision pending — pick based on what's fresher:
- Continue Week 2 quiz (isolation levels beyond RC/RR, row-locking/deadlocks, connection pooling are the natural next slice).
- Redrill B-tree/seq-scan cost mechanism and GIN-vs-GiST structure (new gaps from 2026-09-16, didn't stick on first pass).
- Drill the two flagged backlog priorities (GOMAXPROCS, or architecture-under-pressure — the latter pairs naturally with the still-undone TPS resume-story day).

## Suggested skills for the next session

- Plain quiz continuation needs no skill — just `CURRICULUM.md`'s protocol.
- For drilling GOMAXPROCS or rehearsing the TPS/architecture story under adversarial pressure, call the Skill tool for **`mattpocock-skills:grilling`** — built for relentlessly stress-testing an explanation, which fits better than a plain quiz round.
