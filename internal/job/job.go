// Package job holds the domain types and execution primitives for the
// single-job path: a job's state machine, a thread-safe in-memory store, and
// a generic worker pool. It has no knowledge of gRPC or the proto layer —
// cmd/api maps between this and proto/jobpb.
package job

// State is a job's position in its pending -> running -> succeeded|failed
// state machine.
type State int

const (
	StatePending State = iota
	StateRunning
	StateSucceeded
	StateFailed
)

func (s State) String() string {
	switch s {
	case StatePending:
		return "pending"
	case StateRunning:
		return "running"
	case StateSucceeded:
		return "succeeded"
	case StateFailed:
		return "failed"
	default:
		return "unspecified"
	}
}

// Job is one unit of work submitted through the API.
type Job struct {
	ID      string
	Payload string
	State   State
	Result  string
	Err     string
}
