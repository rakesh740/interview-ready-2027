package job

import "sync"

// Store is a thread-safe in-memory job registry. Reads use RLock so
// concurrent status lookups don't block each other; writes take the full
// lock since they mutate the job in place.
type Store struct {
	mu   sync.RWMutex
	jobs map[string]*Job
}

func NewStore() *Store {
	return &Store{jobs: make(map[string]*Job)}
}

func (s *Store) Put(j *Job) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.jobs[j.ID] = j
}

// Get returns a copy of the job so callers can't mutate store state
// without going through UpdateState.
func (s *Store) Get(id string) (Job, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	j, ok := s.jobs[id]
	if !ok {
		return Job{}, false
	}
	return *j, true
}

func (s *Store) UpdateState(id string, state State, result, errMsg string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if j, ok := s.jobs[id]; ok {
		j.State = state
		j.Result = result
		j.Err = errMsg
	}
}
