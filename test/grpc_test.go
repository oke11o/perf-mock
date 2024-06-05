package test

import (
	"context"
	"net"
	"os"
	"testing"

	"github.com/oke11o/perf-mock/internal/handler"
	"github.com/oke11o/perf-mock/internal/stats"

	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	grpcmock "github.com/oke11o/perf-mock/internal/grpc"
	"github.com/oke11o/perf-mock/internal/logger"
)

func TestGRPCSuite(t *testing.T) {
	suite.Run(t, new(GRPCSuite))
}

type GRPCSuite struct {
	suite.Suite

	addr string
	srv  *grpc.Server

	conn   *grpc.ClientConn
	client grpcmock.TargetServiceClient

	handler *handler.Handler
	mock    *grpcmock.GRPCServer
}

func (s *GRPCSuite) SetupSuite() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8886"
	}
	log := logger.New()
	s.handler = handler.New(stats.NewStats(10))

	s.srv = grpc.NewServer()
	s.mock = grpcmock.NewServer(s.handler, log)
	grpcmock.RegisterTargetServiceServer(s.srv, s.mock)
	reflection.Register(s.srv)

	s.addr = ":" + port
	l, err := net.Listen("tcp", s.addr)
	s.NoError(err)

	go func() {
		err = s.srv.Serve(l)
		s.NoError(err)
	}()
}

func (s *GRPCSuite) TearDownSuite() {
	s.srv.GracefulStop()
}

func (s *GRPCSuite) SetupTest() {
	s.handler.Reset()

	conn, err := grpc.NewClient(s.addr, grpc.WithInsecure(), grpc.WithBlock())
	s.NoError(err)
	s.conn = conn
	s.client = grpcmock.NewTargetServiceClient(conn)
}

func (s *GRPCSuite) TearDownTest() {
	s.conn.Close()
}

func (s *GRPCSuite) Test_SuccessScenario() {
	res, err := s.client.Hello(context.Background(), &grpcmock.HelloRequest{Name: "John"})
	s.NoError(err)
	s.Equal("Hello John!", res.Hello)
	s.Equal(int64(1), s.handler.Stats().Hello)
}
