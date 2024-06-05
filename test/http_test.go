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
	addr string
	mock *httpmock.Server
}

func (s *HTTPSuite) SetupSuite() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8886"
	}
	s.addr = "localhost:" + port
	log := logger.New()
	s.mock = httpmock.NewServer(s.addr, log, time.Now().UnixNano())
	s.mock.ServeAsync()

	go func() {
		err := <-s.mock.Err()
		s.NoError(err)
	}()
}

func (s *HTTPSuite) TearDownSuite() {
	err := s.mock.Shutdown(context.Background())
	s.NoError(err)
}

func (s *HTTPSuite) SetupTest() {
	s.mock.Stats().Reset()
}

func (s *HTTPSuite) Test_Hello() {
	_, err := http.Get("http://" + s.addr + "/hello")
	s.NoError(err)

	stats := s.mock.Stats().Response()
	s.Equal(int64(1), stats.Hello)
}
