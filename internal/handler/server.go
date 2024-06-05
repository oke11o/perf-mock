package handler

import (
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"sync"

	"github.com/yandex/pandora/lib/str"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/oke11o/perf-mock/internal/stats"
)

const (
	userCount         = 10
	userMultiplicator = 1000
	itemMultiplicator = 100
)

func New(stats *stats.Stats) *Handler {
	keys := make(map[string]int64, userCount)
	for i := int64(1); i <= userCount; i++ {
		keys[str.RandStringRunes(64, "")] = i
	}
	return &Handler{stats: stats, keys: keys}
}

type Handler struct {
	log   *slog.Logger
	stats *stats.Stats
	keys  map[string]int64
	mu    sync.RWMutex
}

func (h *Handler) Hello() {
	h.stats.IncHello()
}

func (h *Handler) Reset() stats.Response {
	h.stats.Reset()
	return h.Stats()
}

func (h *Handler) Stats() stats.Response {
	return h.stats.Response()
}

func (h *Handler) Auth(login string, pass string) (AuthResponse, error) {
	userID, token, err := h.checkLoginPass(login, pass)
	if err != nil {
		h.stats.IncAuth400()
		return AuthResponse{}, status.Error(codes.InvalidArgument, "invalid credentials")
	}
	h.stats.IncAuth200(userID)
	return AuthResponse{
		UserID: userID,
		Token:  token,
	}, nil
}

func (h *Handler) List(token string, userID int64) (*ListResponse, error) {
	h.mu.RLock()
	u := h.keys[token]
	h.mu.RUnlock()
	if userID == 0 {
		h.stats.IncList400()
		return nil, errors.New("invalid token")
	}
	if u != userID {
		h.stats.IncList400()
		return nil, errors.New("invalid user_id")
	}

	// Logic
	result := &ListResponse{}
	userID *= userMultiplicator
	result.Result = make([]*ListItem, itemMultiplicator)
	for i := int64(0); i < itemMultiplicator; i++ {
		result.Result[i] = &ListItem{ItemId: userID + i}
	}
	h.stats.IncList200(userID)
	return result, nil
}

func (h *Handler) Order(token string, userID int64, ItemId int64) (*OrderResponse, error) {
	h.mu.RLock()
	u := h.keys[token]
	h.mu.RUnlock()
	if userID == 0 {
		h.stats.IncList400()
		return nil, errors.New("invalid token")
	}
	if u != userID {
		h.stats.IncList400()
		return nil, errors.New("invalid user_id")
	}

	// Logic
	ranger := userID * userMultiplicator
	if ItemId < ranger || ItemId >= ranger+itemMultiplicator {
		h.stats.IncOrder400()
		return nil, errors.New("invalid item_id")
	}

	result := &OrderResponse{}
	result.OrderId = ItemId + 12345
	h.stats.IncOrder200(userID)
	return result, nil
}

func (h *Handler) ExtStats(login string, pass string) {

}

func (h *Handler) checkLoginPass(login string, pass string) (int64, string, error) {
	if login != pass {
		return 0, "", fmt.Errorf("invalid login %s or pass %s", login, pass)
	}
	userID, err := strconv.ParseInt(login, 10, 64)
	if err != nil {
		return 0, "", fmt.Errorf("invalid login %s", login)
	}
	token := ""
	h.mu.RLock()
	for k, v := range h.keys {
		if v == userID {
			token = k
			break
		}
	}
	h.mu.RUnlock()
	if token == "" {
		return 0, "", fmt.Errorf("invalid login %s and pass %s", login, pass)
	}

	return userID, token, nil
}
