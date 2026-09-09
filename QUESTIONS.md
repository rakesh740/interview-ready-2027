# Quiz Log

Running log of interview-quiz sessions: topic covered, date, and weak spots to revisit. Appended after each "quiz me" session — see `CURRICULUM.md` for the schedule this follows.

<!-- Format:
## YYYY-MM-DD — <topic>
- Covered: ...
- Weak spots: ...
-->

## 2026-08-30 — Week 1: Go senior-level depth
- Covered: GMP scheduler model (M/P detach on blocking syscalls, sysmon lazy retake), tricolor mark-sweep GC + write barrier invariant, Go memory model / happens-before (data races, atomic Store/Load as fences), `errgroup.WithContext` vs. manual `context.WithCancel`+`WaitGroup` for fan-out/cancel-on-first-error, generics constraints (`comparable` vs. custom arithmetic unions, `~underlying` types).
- Weak spots:
  - GMP: initially thought a *new goroutine* gets scheduled onto M1 during a blocking syscall — actually M1 stays blocked holding G1; **P1** is what gets detached and handed to another M. Revisit the sysmon retake thresholds (~20µs backing off to ~10ms).
  - Tricolor GC: didn't recall the color semantics (white/grey/black) or the black-points-to-white invariant; couldn't construct the write-barrier race scenario unprompted. Redrill next Fabric-adjacent week (ties to the MVCC story).
  - Memory model: got the happens-before question backwards — assumed program order alone guarantees cross-goroutine visibility. Core rule to re-internalize: no happens-before edge (channel/mutex/atomic/WaitGroup/Once) = undefined behavior, regardless of program order.
  - Generics: needed the full explanation for why `comparable` doesn't grant arithmetic ops, and why a custom union constraint doesn't cover named types over a primitive without `~`.

## 2026-08-30 — Backlog: review of 19 past mock-interview transcripts
- Source: `~/Projects/meetings/reviewed*_clean_text.txt` — real practice/mock interview recordings, not a session run through this curriculum. Logged here as a standing backlog to drive future "quiz me" topic selection until each item is drilled and can be struck through.
- Covered (topics the transcripts touched): Go concurrency internals (GMP/GOMAXPROCS, channels, mutexes, race detection, escape analysis), Go language fundamentals (interfaces, nil-interface gotcha, comparable map keys, internal/pkg), five system designs (rate limiter, URL shortener, Stripe payment/idempotency, OAuth "Sign in with Google", product catalog at scale), Hyperledger Fabric/Raft consensus, Postgres/SQL/NoSQL, Kafka vs. message queues, API design (middleware, validation, protobuf compat), design patterns, Kubernetes/Docker/AWS, networking (TCP/TLS/IPv6), git, and explaining own past project work under adversarial questioning.
- Weak spots:
  - **GOMAXPROCS reasoning** — repeatedly can't justify a concurrency number from first principles (why GOMAXPROCS > core count adds no parallelism, why 10 goroutines was chosen for the Fabric batching optimization). Recurring across multiple interviews — needs a rehearsed, quantitative answer.
  - **`go build -race`/`-race` flag scope** — hedged on which commands support it (works with `go run`, `go test`, `go build`).
  - **Escape analysis** — flat "not sure," no attempt. Needs: compiler decides stack vs. heap via `go build -gcflags="-m"`, based on whether a pointer escapes the function's scope.
  - **Goroutine-leak detection tooling** — only named prevention (context/timeouts), never detection tools (`pprof` goroutine profile, `runtime.NumGoroutine()`, leaktest helpers).
  - **Live data-race debugging** — stayed generic ("logging tools"); never named `-race`, `pprof` mutex/block profiles, or `SIGQUIT` goroutine dumps, or how to handle a heisenbug that disappears under logging (observer effect).
  - **Object-pool sizing** and **string-concat hot-path optimization** (`strings.Builder` vs. `+=`) both drew flat "not sure."
  - **Typed-nil-in-interface gotcha** — got the right answer but thin on the "(type, value) pair" mechanic; keep sharp, common trick question.
  - **Interface as map key** — missed that a non-comparable dynamic type inside an interface panics at runtime on `==`.
  - **Isolation levels** — forgot read-uncommitted as the 4th level and the anomaly each level prevents (dirty/non-repeatable/phantom read).
  - **Rate limiter design** — never resolved the check-then-increment race (needs atomic `INCR`+`EXPIRE`/Lua script); replicated-Redis failover discussion broke down without a real answer; never decisively compared token bucket vs. fixed/sliding window before picking one.
  - **URL shortener** — abandoned the "deterministic and unique" requirement for random UUIDs instead of hash+base62+collision-retry; never discussed collision handling.
  - **Payment/Stripe idempotency** — no proactive reconciliation job for stuck webhook/payment-intent state (only reactive once prompted); "webhook arrives before DB commit" ordering problem answered circularly instead of naming idempotency-key-in-transaction or outbox pattern; state-machine transitions "at scale" answered at single-request level only (missing optimistic concurrency/version column/CAS); distributed transactions across DBs drew a flat "not worked with it" instead of naming 2PC/saga as the standard answers even without hands-on experience.
  - **OAuth "Sign in with Google"** — authorization-code exchange step missing entirely from the first whiteboard pass; token-custody model (backend holds it, not client) needed self-correction; didn't know EventBridge Scheduler as a cron alternative for recurring charges.
  - **Catalog/search at scale** — DB choice defaulted to "MongoDB, it's NoSQL" without grounding in access patterns; never proactively proposed the operational-store + synced-search-index (CQRS-style) split; cross-shard query fan-out/merge stayed abstract.
  - **Hyperledger Fabric architecture rationale** — could describe the outcome ("etcd-raft, it plateaus") but not own the reasoning behind alternatives or why horizontal scaling of the orderer doesn't work beyond CAP/CP trade-off; never mentioned sharding the ordering service (multiple Raft groups/channels) as Fabric's actual scaling lever.
  - **Postgres query-plan tuning** — flat "haven't faced this" instead of naming `ANALYZE`/stats freshness, planner hints, or index rewrites.
  - **HTTP redirect codes stated wrong under pressure** — said 302 is permanent/303 is temporary (backwards: 301 = permanent, 302/303/307 = temporary flavors).
  - **NoSQL/BASE model** — self-admittedly incomplete on which ACID property NoSQL drops (isolation) and whether NoSQL supports indexes (it does).

## 2026-09-09 — Week 2: Postgres MVCC deep dive (resumed, Q1 only)
- Covered: `xmin`/`xmax` tuple versioning and snapshot-based visibility rules; write-write conflict mechanics (second writer blocks on first writer's uncommitted `xmax`, via the implicit per-XID lock); read committed (re-fetch + `EvalPlanQual`, silent retry) vs. repeatable read/serializable (`ERROR: could not serialize access due to concurrent update`) divergence after unblocking; `VACUUM`'s two distinct jobs (dead-tuple/bloat reclamation vs. anti-wraparound freezing); Fabric OCC vs. Postgres MVCC compare/contrast.
- Weak spots:
  - **Initial MVCC answer defaulted straight to Fabric's model** (single version + key, reject on stale write) when asked specifically about Postgres — conflated the two before being told to separate them. Didn't know `xmin`/`xmax` existed unprompted.
  - **Write-write conflict mechanics** — first guess had it backwards (thought the second writer queries `xmin`, not `xmax`); needed correction before landing the actual mechanism (wait on the first writer's XID via the transaction's implicit self-lock).
  - **RC vs. RR/serializable divergence after unblocking** — got this one right largely unprompted, with correct reasoning tied to what "repeatable read" promises. Keep as a confirmed-solid area, don't over-drill.
  - **VACUUM's two failure modes** — got dead-tuple bloat immediately, but transaction-ID wraparound needed two guesses ("multiple tuples same xmax") before landing as "no idea, please explain." Needed the full walkthrough (circular 32-bit comparison, freezing, `autovacuum_freeze_max_age`) fully explained rather than derived. **This is a genuine new gap, not previously logged** — drill again from scratch next Postgres-adjacent session, don't assume it stuck after one explanation.
  - **Fabric/Postgres architectural compare-contrast — confirms the flagged recurring pattern** (see backlog entry above: "explaining own project's architecture/rationale under pressure"). First pass was vocabulary-swap-level ("Fabric uses version+reject, Postgres uses xmin/xmax") with no "why" behind the mechanism difference. Needed the full explanation (Postgres retains multiple versions for long-lived-reader snapshot isolation without blocking; Fabric doesn't need that because it's a short optimistic-validate pipeline) fully given rather than derived, even after being told directly this is a known gap.
- Not yet covered this week: isolation levels compare/contrast beyond RC/RR, indexing strategies (B-tree vs. GIN/GIST, partial/covering), `EXPLAIN ANALYZE`, row-level locking/deadlocks beyond what came up here, connection pooling, idempotency keys/outbox/retry-backoff, K8s basics. Resume Week 2 quiz from indexing or isolation-levels-in-general next session.
  - **AWS access control** — proposed IP-allowlisting for S3-from-EC2 access instead of IAM roles/instance profiles + bucket policy conditions.
  - **Kafka vs. message queue distinction** — never landed on durable/replayable partitioned log with consumer-tracked offsets vs. queue removing on ack; unsure whether classic queues support replay.
  - **Middleware for cross-cutting concerns** — didn't connect global panic-recovery/error-handling to middleware unprompted; validation design started per-handler instead of proposing a single struct-tag-driven validation middleware (`go-playground/validator`) upfront.
  - **Design patterns** — factory pattern conflated with "a single constructor"; singleton thread-safety mitigation stalled after one idea (immutability, DI of scoped instances not raised).
  - **Protobuf/gRPC forward compatibility** and the `[deprecated = true]` field option — open gaps (backward compat was fine).
  - **Kubernetes**: DaemonSet unknown outright; Secrets described vaguely ("a separate file") instead of the `Secret` object/base64/mount model; StatefulSet identity criterion confused with "has a volume attached" (it's stable network identity + per-ordinal storage).
  - **Docker `stop` vs `kill`** — described `kill` as removing the image rather than force-terminating without the SIGTERM→SIGKILL grace period; `ENTRYPOINT` purpose unsure.
  - **AWS storage taxonomy** — conflated storage services (S3/EBS/EFS) with managed DB services (RDS/DynamoDB); region vs. AZ left an unanswered follow-up.
  - **Git CLI fluency** — relies on GitHub's UI for squashing instead of knowing `git reset`/`git revert`/`git rebase -i` from the terminal.
  - **TCP vs. TLS conflation** — attributed certificate/signature validation to TCP itself rather than the TLS layer on top of it.
  - **IPv4 vs. IPv6** — hazy on 32-bit vs. 128-bit address space and the address-exhaustion motivation.
  - **Explaining own project work under pressure** — repeatedly couldn't get past "it's a login integration" when asked what an onboarding workflow's actual data model/entity lifecycle was (10+ exchanges without landing it); couldn't articulate architecture rationale for a system built by others beyond "they needed an enterprise blockchain solution" when pushed on design vs. development involvement; scale numbers (onboarding frequency, unit of scale) given inconsistently. Worth rehearsing a fixed "one entity, its lifecycle states, who acts at each state" template plus a consistent scale-numbers story per project.
