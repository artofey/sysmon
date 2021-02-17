package storage

import (
	"container/ring"
	"sync"

	"github.com/artofey/sysmon"
)

type ringStorage struct {
	count int

	mu sync.Mutex
	s  *ring.Ring
}

func NewRingStorage(count int) *ringStorage {
	return &ringStorage{
		count: count,
		s:     ring.New(count),
	}
}

func (s *ringStorage) Add(st sysmon.Stats) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.s.Next()
	s.s.Value = st
	return nil
}

func (s *ringStorage) Len() int {
	return s.s.Len()
}

func (s *ringStorage) GetLast(l int) []sysmon.Stats {
	rLen := s.Len()
	if rLen < l {
		return s.get(rLen)
	}
	return s.get(l)
}

func (s *ringStorage) get(l int) []sysmon.Stats {
	s.mu.Lock()
	defer s.mu.Unlock()

	res := make([]sysmon.Stats, 0, l)
	ring := s.s.Prev()
	for i := 0; i < l; i++ {
		res = append(res, ring.Value.(sysmon.Stats))
		ring = ring.Prev()
	}
	return res
}
