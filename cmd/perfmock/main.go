package main

import (
	"fmt"
	"log/slog"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	grpcmock "github.com/oke11o/perf-mock/internal/grpc"
	"github.com/oke11o/perf-mock/internal/handler"
	httpmock "github.com/oke11o/perf-mock/internal/http"
	"github.com/oke11o/perf-mock/internal/logger"
	"github.com/oke11o/perf-mock/internal/stats"
)

func main() {
	log := logger.New()
	newStats := stats.NewStats(10)
	h := handler.New(newStats)
	errCh := make(chan error, 2)
	runGrpc(errCh, log, h)
	runHTTP(errCh, log, newStats, h)
	err := <-errCh
	if err != nil {
		panic(err)
	}
	fmt.Println("FINISH")
}

func runGrpc(errCh chan error, log *slog.Logger, h *handler.Handler) {
	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "8091"
	}

	srv := grpc.NewServer()
	mock := grpcmock.NewServer(h, log)
	grpcmock.RegisterTargetServiceServer(srv, mock)
	reflection.Register(srv)

	addr := ":" + port
	l, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	go func() {
		err = srv.Serve(l)
		errCh <- err
	}()
}

func runHTTP(errCh chan error, log *slog.Logger, newStats *stats.Stats, h *handler.Handler) {
	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "8092"
	}
	addr := ":" + port

	mock := httpmock.NewServer(addr, log, newStats, h)
	mock.ServeAsync()

	go func() {
		err := <-mock.Err()
		errCh <- err
	}()
}
