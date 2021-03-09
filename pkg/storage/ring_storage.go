package storage

import (
	"container/ring"
	"sync"

	"github.com/artofey/sysmon"
)

type RingStorage struct {
	count int
	size  int

	mu sync.Mutex
	s  *ring.Ring
}

func NewRingStorage(size int) *RingStorage {
	return &RingStorage{
		size: size,
		s:    ring.New(size),
	}
}

func (s *RingStorage) Add(st sysmon.Stats) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.s.Value = st
	s.s = s.s.Next()
	if s.count <= s.size {
		s.count++
	}
	return nil
}

func (s *RingStorage) Len() int {
	return s.count
}

func (s *RingStorage) GetLast(l int) []sysmon.Stats {
	s.mu.Lock()
	defer s.mu.Unlock()

	rLen := s.Len()
	if rLen < l {
		return s.getLast(rLen)
	}
	return s.getLast(l)
}

func (s *RingStorage) getLast(l int) []sysmon.Stats {
	res := make([]sysmon.Stats, 0, l)
	ring := s.s.Prev()
	for i := 0; i < l; i++ {
		res = append(res, ring.Value.(sysmon.Stats))
		ring = ring.Prev()
	}
	return res
}
