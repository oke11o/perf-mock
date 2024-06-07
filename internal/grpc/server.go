package grpc

import (
	"context"
	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/oke11o/perf-mock/internal/handler"
)

func NewServer(h *handler.Handler, logger *slog.Logger) *GRPCServer {
	return &GRPCServer{handler: h, log: logger}
}

type GRPCServer struct {
	UnimplementedTargetServiceServer

	handler *handler.Handler
	log     *slog.Logger
}

var _ TargetServiceServer = (*GRPCServer)(nil)

func (s *GRPCServer) Hello(ctx context.Context, request *HelloRequest) (*HelloResponse, error) {
	res := s.handler.Hello(request.Name, request.SkipStats, request.Sleep.AsDuration())
	return &HelloResponse{
		Hello: res,
	}, nil
}

func (s *GRPCServer) Auth(_ context.Context, request *AuthRequest) (*AuthResponse, error) {
	res, err := s.handler.Auth(request.GetLogin(), request.GetPass())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid credentials")
	}
	return &AuthResponse{
		UserId: res.UserID,
		Token:  res.Token,
	}, nil
}

func (s *GRPCServer) List(_ context.Context, request *ListRequest) (*ListResponse, error) {
	res, err := s.handler.List(request.Token, request.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	result := &ListResponse{}
	result.Result = make([]*ListItem, len(res.Result))
	for i, r := range res.Result {
		result.Result[i] = &ListItem{ItemId: r.ItemId}
	}
	return result, nil
}

func (s *GRPCServer) Order(ctx context.Context, request *OrderRequest) (*OrderResponse, error) {
	res, err := s.handler.Order(request.Token, request.UserId, request.ItemId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &OrderResponse{OrderId: res.OrderId}, nil
}

func (s *GRPCServer) Stats(ctx context.Context, _ *StatsRequest) (*StatsResponse, error) {
	ss := s.handler.Stats()
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

func (s *GRPCServer) Reset(ctx context.Context, _ *ResetRequest) (*StatsResponse, error) {
	ss := s.handler.Reset()

	result := &StatsResponse{
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
