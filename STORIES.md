# Resume Stories

Polished STAR-format (Situation / Task / Action / Result) answers for resume bullets, built and drilled during this month's "resume-story days" (see `CURRICULUM.md`). Source: `Rakesh_Rakshit_software_engineer_6YOE_15_jan.pdf`.

Background for context: 6 YOE — Kartbites Pvt Ltd (2018-2022, food delivery app: Go, PostgreSQL, PostGIS, GraphQL, Redis, nsqd, Stripe, Firebase/APNS, OAuth), CoffeeBeans Consulting (2022-2024, Hyperledger Fabric for a fintech client: throughput/concurrency optimization, MVCC, Raft/etcd), Tech Prescient (2024-present, IAM platform on distributed microservices, onboarding-workflow automation).

## 1. TPS optimization (2000 → 8000 TPS) — CoffeeBeans / Hyperledger Fabric
_Status: not yet drilled (Week 1)._

- **Situation**:
- **Task**:
- **Action**:
- **Result**:
- **Likely follow-ups**: what was the actual bottleneck (endorsement, ordering, commit phase)? how was payload size reduced? how did that reduce bandwidth specifically? how was 8000 TPS measured/verified?

## 2. MVCC + snapshot management — CoffeeBeans / Hyperledger Fabric
_Status: not yet drilled (Week 2)._

- **Situation**:
- **Task**:
- **Action**:
- **Result**:
- **Likely follow-ups**: what consistency problem did this solve? what would've broken without it? how does Fabric's MVCC differ from Postgres's?

## 3. Driver-order-route matching algorithm (20% cost reduction) — Kartbites
_Status: not yet drilled (Week 2)._

- **Situation**:
- **Task**:
- **Action**:
- **Result**:
- **Likely follow-ups**: what did the matching algorithm actually optimize for? how was "cost" measured? what were the tradeoffs vs. a naive nearest-driver assignment?

## 4. IAM onboarding-workflow automation — Tech Prescient
_Status: not yet drilled (Week 12 — draw the parallel to the CRD/controller reconciliation loop built that week)._

- **Situation**:
- **Task**:
- **Action**:
- **Result**:
- **Likely follow-ups**: what did the approval/rejection flow look like end-to-end? how is this conceptually similar to a Kubernetes controller reconciling desired vs. actual state?

## 5. Mentoring junior developers — Kartbites
_Status: not yet drilled (Week 17 — behavioral-round material for mock interviews)._

- **Situation**:
- **Task**:
- **Action**:
- **Result**:
- **Likely follow-ups**: what specifically did you teach (Golang idioms? PostgreSQL patterns?)? how did you know it landed? a time mentoring didn't go well and what changed after.

## 6. Distributed systems background (Raft/etcd, CAP theorem) — CoffeeBeans
_Status: informal drilling throughout Month 2 Week 9 (leader election week) — capture any new angle here as it surfaces, formal pass in Week 17._

- **Situation**:
- **Task**:
- **Action**:
- **Result**:
- **Likely follow-ups**: how does etcd's Raft implementation compare to the K8s Lease-based leader election used in this project's dispatcher? concrete CAP-theorem tradeoff you navigated in production, not textbook.
