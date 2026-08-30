package main

import (
	"context"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"

	jobpb "github.com/rakeshrakshit/interview-ready-2027/proto/jobpb"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		logger.Error("failed to listen", "err", err)
		os.Exit(1)
	}

	srv := newJobServer(logger, 4, 100)
	grpcServer := grpc.NewServer()
	jobpb.RegisterJobServiceServer(grpcServer, srv)

	go func() {
		logger.Info("api server listening", "addr", lis.Addr().String())
		if err := grpcServer.Serve(lis); err != nil {
			logger.Error("grpc server stopped", "err", err)
		}
	}()

	<-ctx.Done()
	logger.Info("shutdown signal received")

	grpcServer.GracefulStop()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.pool.Shutdown(shutdownCtx); err != nil {
		logger.Error("pool shutdown incomplete", "err", err)
	}

	logger.Info("shutdown complete")
}
