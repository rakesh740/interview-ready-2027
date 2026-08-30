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
