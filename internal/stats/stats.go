package stats

import (
	"sync"
	"sync/atomic"
)

func NewStats(capacity int) *Stats {
	stats := Stats{
		auth200:       make(map[int64]uint64, capacity),
		auth200Mutex:  sync.Mutex{},
		auth400:       atomic.Uint64{},
		auth500:       atomic.Uint64{},
		list200:       make(map[int64]uint64, capacity),
		list200Mutex:  sync.Mutex{},
		list400:       atomic.Uint64{},
		list500:       atomic.Uint64{},
		order200:      make(map[int64]uint64, capacity),
		order200Mutex: sync.Mutex{},
		order400:      atomic.Uint64{},
		order500:      atomic.Uint64{},
	}
	return &stats
}

type Stats struct {
	hello         atomic.Uint64
	auth200       map[int64]uint64
	auth200Mutex  sync.Mutex
	auth400       atomic.Uint64
	auth500       atomic.Uint64
	list200       map[int64]uint64
	list200Mutex  sync.Mutex
	list400       atomic.Uint64
	list500       atomic.Uint64
	order200      map[int64]uint64
	order200Mutex sync.Mutex
	order400      atomic.Uint64
	order500      atomic.Uint64
}

func (s *Stats) IncHello() {
	s.hello.Add(1)
}

func (s *Stats) IncAuth400() {
	s.auth400.Add(1)
}

func (s *Stats) IncAuth500() {
	s.auth500.Add(1)
}

func (s *Stats) IncAuth200(userID int64) {
	s.auth200Mutex.Lock()
	s.auth200[userID]++
	s.auth200Mutex.Unlock()
}

func (s *Stats) IncList400() {
	s.list400.Add(1)
}

func (s *Stats) IncList500() {
	s.list500.Add(1)
}

func (s *Stats) IncList200(userID int64) {
	s.list200Mutex.Lock()
	s.list200[userID]++
	s.list200Mutex.Unlock()
}

func (s *Stats) IncOrder400() {
	s.order400.Add(1)
}

func (s *Stats) IncOrder500() {
	s.order500.Add(1)
}

func (s *Stats) IncOrder200(userID int64) {
	s.order200Mutex.Lock()
	s.order200[userID]++
	s.order200Mutex.Unlock()
}

func (s *Stats) Reset() {
	s.hello.Store(0)

	s.auth200Mutex.Lock()
	s.auth200 = map[int64]uint64{}
	s.auth200Mutex.Unlock()
	s.auth400.Store(0)
	s.auth500.Store(0)

	s.list200Mutex.Lock()
	s.list200 = map[int64]uint64{}
	s.list200Mutex.Unlock()
	s.list400.Store(0)
	s.list500.Store(0)

	s.order200Mutex.Lock()
	s.order200 = map[int64]uint64{}
	s.order200Mutex.Unlock()
	s.order400.Store(0)
	s.order500.Store(0)
}

type Response struct {
	Auth  Body  `json:"auth,omitempty"`
	List  Body  `json:"list,omitempty"`
	Order Body  `json:"order,omitempty"`
	Hello int64 `json:"hello,omitempty"`
}

type Body struct {
	Code200 map[int64]uint64 `json:"code200,omitempty"`
	Code400 uint64           `json:"code400,omitempty"`
	Code500 uint64           `json:"code500,omitempty"`
}

func (s *Stats) Response() Response {
	return Response{
		Hello: int64(s.hello.Load()),
		Auth: Body{
			Code200: s.auth200,
			Code400: s.auth400.Load(),
			Code500: s.auth500.Load(),
		},
		List: Body{
			Code200: s.list200,
			Code400: s.list400.Load(),
			Code500: s.list500.Load(),
		},
		Order: Body{
			Code200: s.order200,
			Code400: s.order400.Load(),
			Code500: s.order500.Load(),
		},
	}
}
