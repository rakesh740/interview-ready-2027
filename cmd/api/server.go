package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/rakeshrakshit/interview-ready-2027/internal/job"
	jobpb "github.com/rakeshrakshit/interview-ready-2027/proto/jobpb"
)

type jobServer struct {
	jobpb.UnimplementedJobServiceServer
	store *job.Store
	pool  *job.Pool[string]
	log   *slog.Logger
}

func newJobServer(log *slog.Logger, workers, queueSize int) *jobServer {
	s := &jobServer{
		store: job.NewStore(),
		log:   log,
	}
	s.pool = job.NewPool(workers, queueSize, s.process)
	return s
}

func (s *jobServer) SubmitJob(ctx context.Context, req *jobpb.SubmitJobRequest) (*jobpb.SubmitJobResponse, error) {
	id := genID()
	s.store.Put(&job.Job{ID: id, Payload: req.GetPayload(), State: job.StatePending})
	s.log.Info("job submitted", "job_id", id)

	if err := s.pool.Submit(ctx, id); err != nil {
		return nil, status.Errorf(codes.Unavailable, "submit: %v", err)
	}
	return &jobpb.SubmitJobResponse{JobId: id}, nil
}

func (s *jobServer) GetJobStatus(_ context.Context, req *jobpb.GetJobStatusRequest) (*jobpb.GetJobStatusResponse, error) {
	j, ok := s.store.Get(req.GetJobId())
	if !ok {
		return nil, status.Errorf(codes.NotFound, "job %q not found", req.GetJobId())
	}
	return &jobpb.GetJobStatusResponse{
		JobId:  j.ID,
		State:  toProtoState(j.State),
		Result: j.Result,
		Error:  j.Err,
	}, nil
}

// process is the worker-pool handler: it moves a job from pending through
// running to a terminal state. Runs on a pool worker goroutine, not the
// gRPC request goroutine.
func (s *jobServer) process(ctx context.Context, id string) {
	j, ok := s.store.Get(id)
	if !ok {
		return
	}
	s.store.UpdateState(id, job.StateRunning, "", "")
	s.log.Info("job running", "job_id", id)

	select {
	case <-time.After(500 * time.Millisecond): // stand-in for real work
	case <-ctx.Done():
		return
	}

	result := fmt.Sprintf("processed:%s", j.Payload)
	s.store.UpdateState(id, job.StateSucceeded, result, "")
	s.log.Info("job succeeded", "job_id", id)
}

func toProtoState(st job.State) jobpb.JobState {
	switch st {
	case job.StatePending:
		return jobpb.JobState_JOB_STATE_PENDING
	case job.StateRunning:
		return jobpb.JobState_JOB_STATE_RUNNING
	case job.StateSucceeded:
		return jobpb.JobState_JOB_STATE_SUCCEEDED
	case job.StateFailed:
		return jobpb.JobState_JOB_STATE_FAILED
	default:
		return jobpb.JobState_JOB_STATE_UNSPECIFIED
	}
}

func genID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
