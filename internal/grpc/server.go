package grpc

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"sync"

	"github.com/oke11o/perf-mock/internal/stats"

	"github.com/yandex/pandora/lib/str"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	userCount         = 10
	userMultiplicator = 1000
	itemMultiplicator = 100
)

func NewServer(logger *slog.Logger, seed int64) *GRPCServer {
	keys := make(map[string]int64, userCount)
	for i := int64(1); i <= userCount; i++ {
		keys[str.RandStringRunes(64, "")] = i
	}
	logger.Info("New server created", slog.Any("keys", keys))

	return &GRPCServer{Log: logger, keys: keys, stats: stats.NewStats(userCount)}
}

type GRPCServer struct {
	UnimplementedTargetServiceServer
	Log   *slog.Logger
	stats *stats.Stats
	keys  map[string]int64
	mu    sync.RWMutex
}

var _ TargetServiceServer = (*GRPCServer)(nil)

func (s *GRPCServer) Hello(ctx context.Context, request *HelloRequest) (*HelloResponse, error) {
	s.stats.IncHello()
	return &HelloResponse{
		Hello: fmt.Sprintf("Hello %s!", request.Name),
	}, nil
}

func (s *GRPCServer) Auth(ctx context.Context, request *AuthRequest) (*AuthResponse, error) {
	userID, token, err := s.checkLoginPass(request.GetLogin(), request.GetPass())
	if err != nil {
		s.stats.IncAuth400()
		return nil, status.Error(codes.InvalidArgument, "invalid credentials")
	}
	result := &AuthResponse{
		UserId: userID,
		Token:  token,
	}
	s.stats.IncAuth200(userID)
	return result, nil
}

func (s *GRPCServer) List(ctx context.Context, request *ListRequest) (*ListResponse, error) {
	s.mu.RLock()
	userID := s.keys[request.Token]
	s.mu.RUnlock()
	if userID == 0 {
		s.stats.IncList400()
		return nil, status.Error(codes.InvalidArgument, "invalid token")
	}
	if userID != request.UserId {
		s.stats.IncList400()
		return nil, status.Error(codes.InvalidArgument, "invalid user_id")
	}

	// Logic
	result := &ListResponse{}
	userID *= userMultiplicator
	result.Result = make([]*ListItem, itemMultiplicator)
	for i := int64(0); i < itemMultiplicator; i++ {
		result.Result[i] = &ListItem{ItemId: userID + i}
	}
	s.stats.IncList200(request.UserId)
	return result, nil
}

func (s *GRPCServer) Order(ctx context.Context, request *OrderRequest) (*OrderResponse, error) {
	s.mu.RLock()
	userID := s.keys[request.Token]
	s.mu.RUnlock()
	if userID == 0 {
		s.stats.IncOrder400()
		return nil, status.Error(codes.InvalidArgument, "invalid token")
	}
	if userID != request.UserId {
		s.stats.IncOrder400()
		return nil, status.Error(codes.InvalidArgument, "invalid user_id")
	}

	// Logic
	ranger := userID * userMultiplicator
	if request.ItemId < ranger || request.ItemId >= ranger+itemMultiplicator {
		s.stats.IncOrder400()
		return nil, status.Error(codes.InvalidArgument, "invalid item_id")
	}

	result := &OrderResponse{}
	result.OrderId = request.ItemId + 12345
	s.stats.IncOrder200(userID)
	return result, nil
}

func (s *GRPCServer) Stats(ctx context.Context, _ *StatsRequest) (*StatsResponse, error) {
	ss := s.stats.Response()
	result := &StatsResponse{
		Hello: ss.Hello,
		Auth: &StatisticBodyResponse{
			Code200: ss.Auth.Code200,
			Code400: ss.Auth.Code400,
			Code500: ss.Auth.Code500,
		},
		List: &StatisticBodyResponse{
			Code200: ss.List.Code200,
			Code400: ss.List.Code400,
			Code500: ss.List.Code500,
		},
		Order: &StatisticBodyResponse{
			Code200: ss.Order.Code200,
			Code400: ss.Order.Code400,
			Code500: ss.Order.Code500,
		},
	}
	return result, nil
}

func (s *GRPCServer) Reset(ctx context.Context, _ *ResetRequest) (*ResetResponse, error) {
	s.stats.Reset()
	ss := s.stats.Response()

	result := &ResetResponse{
		Auth: &StatisticBodyResponse{
			Code200: ss.Auth.Code200,
			Code400: ss.Auth.Code400,
			Code500: ss.Auth.Code500,
		},
		List: &StatisticBodyResponse{
			Code200: ss.List.Code200,
			Code400: ss.List.Code400,
			Code500: ss.List.Code500,
		},
		Order: &StatisticBodyResponse{
			Code200: ss.Order.Code200,
			Code400: ss.Order.Code400,
			Code500: ss.Order.Code500,
		},
	}
	return result, nil
}

func (s *GRPCServer) checkLoginPass(login string, pass string) (int64, string, error) {
	userID, err := strconv.ParseInt(login, 10, 64)
	if err != nil {
		return 0, "", fmt.Errorf("invalid login %s", login)
	}
	if login != pass {
		return 0, "", fmt.Errorf("invalid login %s or pass %s", login, pass)
	}
	token := ""
	s.mu.RLock()
	for k, v := range s.keys {
		if v == userID {
			token = k
			break
		}
	}
	s.mu.RUnlock()
	if token == "" {
		return 0, "", fmt.Errorf("invalid login %s and pass %s", login, pass)
	}

	return userID, token, nil
}
