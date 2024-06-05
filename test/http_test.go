package test

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	httpmock "github.com/oke11o/perf-mock/internal/http"
	"github.com/oke11o/perf-mock/internal/logger"
)

func TestHTTPSuite(t *testing.T) {
	suite.Run(t, new(HTTPSuite))
}

type HTTPSuite struct {
	suite.Suite
	server *httpmock.Server
	addr   string
}

func (s *HTTPSuite) SetupSuite() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8886"
	}
	s.addr = "localhost:" + port
	log := logger.New()
	s.server = httpmock.NewServer(s.addr, log, time.Now().UnixNano())
	s.server.ServeAsync()

	go func() {
		err := <-s.server.Err()
		s.NoError(err)
	}()
}

func (s *HTTPSuite) TearDownSuite() {
	err := s.server.Shutdown(context.Background())
	s.NoError(err)
}

func (s *HTTPSuite) SetupTest() {
	s.server.Stats().Reset()
}

func (s *HTTPSuite) Test_SuccessScenario() {
	_, err := http.Get("http://" + s.addr + "/hello")
	s.NoError(err)

	stats := s.server.Stats().Response()
	s.Equal(map[int64]uint64{1: 3, 2: 3, 3: 3}, stats.Order)
}
