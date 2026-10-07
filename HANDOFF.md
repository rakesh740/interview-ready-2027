# Handoff — interview-ready-2027

_Written 2026-09-29. Read this at the start of a new session before doing anything else, then delete or ignore once its content is stale._

## What this project is

An 18-week Go/Kubernetes interview-prep curriculum, tracked in `CURRICULUM.md`, that also incrementally builds a DAG workflow orchestration platform (see `README.md`). Trigger phrases: **"quiz me"** / **"today's task"** — read `CURRICULUM.md` for the current week/day. After any quiz/story session, append weak spots to `QUESTIONS.md`.

## Critical working agreement — read this first

**The user writes all implementation code themselves from here forward.** Claude quizzes, explains, and reviews — it does not write feature code. Full rationale in Claude's persistent memory: `~/.claude/projects/-Users-rakeshrakshit/memory/feedback_interview_ready_user_codes.md` (a fresh agent should read that too). One exception already closed: Week 1 code (commit `78c09d2`) stays as Claude-authored; everything after is user-written.

**Standing quiz rule (see Claude's own memory `feedback_quiz_push_for_mechanism.md`):** on every question, every topic, push at least once for the physical/structural mechanism before accepting an answer that only names the right conclusion. If a second prompt still doesn't produce it, give the full explanation directly and log the item as a redrill — don't keep nudging.

## State of the repo

- Working tree clean as of this writing, through commit `4fe94e2`. This handoff update plus the new `QUESTIONS.md` entry from this session are the only pending (uncommitted) changes.
- Week 1 milestone: commit `78c09d2` — gRPC `JobService`, `internal/job` (Store, generic `Pool[T]`), `cmd/api`. `make run` works, verified via `grpcurl`.

## Curriculum progress

- **Week 1 resume-story day — still NOT done.** TPS-optimization story (2000→8000 TPS, Hyperledger Fabric) in `STORIES.md` is still blank. Tied to a real interview gap (see backlog below).
- **Week 2 quiz topics — now fully covered at least once.** As of this session (2026-09-29), every topic on the Week 2 quiz-topics list in `CURRICULUM.md` has been touched except idempotency keys/outbox/retry-backoff and K8s basics. Full detail lives in `QUESTIONS.md` — read the entries dated 2026-09-09, 2026-09-16, and 2026-09-29 rather than re-deriving them here. Headline items to carry forward:
  - **"Zero-attempt don't-know" is a new pattern, distinct from the older "right conclusion, no mechanism" habit.** In the 2026-09-29 session, the last three questions in a row (write skew/SSI, deadlock detection, connection pooling) drew an immediate "don't know" with no attempted guess at all — not even a wrong one. Next session, nudge for a guess before accepting "don't know," since a wrong guess is more diagnostic and gives something to correct.
  - **Three genuinely new gaps from 2026-09-29, none of which stuck on a single explanation — redrill from scratch:** write skew / Serializable Snapshot Isolation (predicate/SIREAD locks, dangerous-structure detection); row-lock/deadlock mechanics (optimistic blocking + `deadlock_timeout`-triggered wait-for-graph cycle detection, not proactive checking); connection pooling (process-per-connection fork/memory cost as the resource pgbouncer economizes on, session vs. transaction pooling mode tradeoffs).
  - **Still-open redrill items from 2026-09-16, not yet re-tested:** B-tree/seq-scan cost mechanism (heap-page amortization) and GIN-vs-GiST structure (inverted index vs. lossy bounding-predicate tree).
  - User explicitly asked to revisit `EXPLAIN ANALYZE` depth again later (from 2026-09-16) — still not re-tested, don't treat it as closed.
- **Week 2 topics genuinely not yet touched:** idempotency keys/outbox pattern/retry-backoff+jitter, K8s basics (pods, deployments, services, configmaps/secrets, namespaces).

## Mock-interview transcript backlog

Full extraction (~55 weak spots) lives in `QUESTIONS.md` under `## 2026-08-30 — Backlog: review of 19 past mock-interview transcripts` — do not re-read the raw transcripts (`~/Projects/meetings/reviewed*_clean_text.txt`), that entry is the durable record. Two flagged as highest-priority (recurring across multiple interviews, still undrilled):
1. **GOMAXPROCS reasoning** — can't justify a concurrency number from first principles.
2. **Explaining own project's architecture/data-model rationale under adversarial questioning** — same gap the 2026-09-09 Fabric compare-contrast re-confirmed.

## Options for next session

**K8s basics was started but not finished** (2026-10-07, different device/session than the one that logged the 2026-09-29 entry above): asked how a Deployment's controller detects and replaces a deleted Pod (ReplicaSet controller watching pod status via the API server, reconciling desired vs. actual replica count). Session was paused before the user answered — no weak spot to log yet, just restart K8s basics from that same question rather than assuming it's covered.

No fixed decision pending otherwise — pick based on what's fresher:
- Finish the last two Week 2 quiz topics: idempotency keys/outbox/retry-backoff, K8s basics (resume from the Deployment/ReplicaSet question above) — this closes out Week 2's quiz-topics list entirely.
- Redrill the four accumulated "didn't stick" items in one pass: B-tree/seq-scan cost, GIN-vs-GiST structure, write skew/SSI, deadlock wait-for-graph mechanics, connection pooling resource/pooling-mode tradeoff. This is getting like a real backlog — worth a dedicated redrill session rather than letting it grow further.
- Drill the two flagged backlog priorities (GOMAXPROCS, or architecture-under-pressure — the latter pairs naturally with the still-undone TPS resume-story day).

## Suggested skills for the next session

- Plain quiz continuation needs no skill — just `CURRICULUM.md`'s protocol.
- For drilling GOMAXPROCS or rehearsing the TPS/architecture story under adversarial pressure, call the Skill tool for **`mattpocock-skills:grilling`** — built for relentlessly stress-testing an explanation, which fits better than a plain quiz round.
