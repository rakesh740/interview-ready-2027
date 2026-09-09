# Handoff — interview-ready-2027

_Written 2026-08-30. Read this at the start of a new session before doing anything else, then delete or ignore once its content is stale._

## What this project is

An 18-week Go/Kubernetes interview-prep curriculum, tracked in `CURRICULUM.md`, that also incrementally builds a DAG workflow orchestration platform (see `README.md` for architecture). Trigger phrases from the curriculum: **"quiz me"** or **"today's task"** — read `CURRICULUM.md` for the current week/day and run the right mode (Monday off, Tuesday light quiz, Wed–Fri quiz+small coding increment, Sat–Sun weekend build block). After any quiz/story session, append weak spots to `QUESTIONS.md`.

## Critical working agreement — read this first

**The user writes all implementation code themselves from here forward.** Claude's role is to quiz, explain, and review — not to write feature code. This was established explicitly after Claude built the entire Week 1 milestone unprompted and the user pushed back ("I did not code a single line? how can I learn like this"). Full detail and rationale is saved in Claude's persistent memory at `~/.claude/projects/-Users-rakeshrakshit/memory/feedback_interview_ready_user_codes.md` — a fresh agent should read that file too, since it governs how every future coding session in this repo should run.

One exception was already made and is closed: the Week 1 code (commit `78c09d2`, see below) stays as-is; the user chose not to revert/redo it. Every coding task from Week 2 onward should have the user write the code, with Claude reviewing/interrogating it afterward like a real interviewer would.

## State of the repo

- Commit `78c09d2` — Week 1 milestone: `proto/job.proto` (`JobService` with `SubmitJob`/`GetJobStatus`), generated stubs in `proto/jobpb/`, `internal/job/` (domain types, thread-safe `Store`, generic `Pool[T any]` worker pool), `cmd/api/` (gRPC server, `slog` structured logging, graceful shutdown on SIGINT/SIGTERM). Verified end-to-end via `grpcurl` (submit → running → succeeded, `NotFound` on unknown job, clean shutdown on `SIGTERM`). `make run` starts it.
- **Uncommitted right now:** `QUESTIONS.md` has a large backlog entry appended (see below) that has not been committed. The user was asked whether to commit it and whether to continue Week 2's quiz or tackle backlog items first — that question was pending when this handoff was written; check the live conversation/user reply before assuming either.

## Curriculum progress

- **Week 1 quiz topics** — all covered and logged in `QUESTIONS.md` (`## 2026-08-30 — Week 1: Go senior-level depth`): GMP scheduler/syscall handoff, tricolor GC + write barrier, Go memory model/happens-before, `errgroup.WithContext` vs. manual context+WaitGroup, generics constraints.
- **Week 1 resume-story day — still NOT done.** The TPS-optimization story (2000→8000 TPS, Hyperledger Fabric) in `STORIES.md` is still blank. This is explicitly flagged as connected to a real gap found in the transcript backlog (see below) — the user can state the outcome of the Fabric architecture but struggles to own the design reasoning when pushed. Worth doing the resume-story drill and the backlog item together rather than separately.
- **Week 2 quiz** — started, paused after Q1. First question was Postgres MVCC (`xmin`/`xmax`, tuple visibility, why `VACUUM` is needed) — the user had not answered yet when they paused ("stop for now, will restart later"). Resume from that exact question if continuing Week 2 quiz next.

## Mock-interview transcript backlog (new this session)

The user has 19 real mock-interview transcripts at `~/Projects/meetings/reviewed*_clean_text.txt` (rough ASR transcriptions, not part of this repo). A subagent reviewed all 19 and extracted ~55 concrete weak spots across ~30 topic areas (Go concurrency, 5 system-design problems, Postgres/SQL/NoSQL, Kafka, K8s/Docker/AWS, networking, git, and explaining-own-work-under-pressure). Full list is now in `QUESTIONS.md` under `## 2026-08-30 — Backlog: review of 19 past mock-interview transcripts` — do not re-read the raw transcripts again, that entry is the durable extraction.

Two things flagged as the highest-priority items to drill first (recurring across *multiple separate interviews*, not one-off slips):
1. **GOMAXPROCS reasoning** — can't justify a concurrency number from first principles when pushed.
2. **Explaining own project's data model/architecture rationale under adversarial questioning** — repeatedly collapses to a surface-level description ("it's a login integration") instead of a concrete entity/lifecycle/decision answer. This is the same gap underlying the still-pending TPS resume-story day.

A secondary pattern worth naming to the user if it comes up again: "I'm not sure" used as a terminal answer rather than a reasoned attempt (showed up for escape analysis, DaemonSets, NoSQL indexing, distributed transactions, etc.).

## Suggested skills for the next session

- No skill is strictly required to continue quizzing/reviewing — that's plain conversation following `CURRICULUM.md`'s protocol.
- If the next session focuses on drilling the flagged weak spots (GOMAXPROCS reasoning, or rehearsing the project-explanation/TPS story under adversarial pressure), consider calling the Skill tool for **`mattpocock-skills:grilling`** — it's built for relentlessly stress-testing a plan/answer/explanation, which matches "make the user defend a design decision like a real interviewer would" better than a plain quiz round.
