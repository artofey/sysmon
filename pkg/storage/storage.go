package storage

import (
	"fmt"
	"sync"

	"github.com/artofey/sysmon"
)

type Storage struct {
	count int

	mu sync.Mutex
	s  []sysmon.Stats
}

func NewStorage(count int) *Storage {
	return &Storage{
		count: count,
		s:     make([]sysmon.Stats, 0, count),
	}
}

func (s *Storage) Add(st sysmon.Stats) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.s) > s.count {
		return fmt.Errorf("Error")
	}
	s.s = append(s.s, st)
	return nil
}

func (s *Storage) Len() int {
	return len(s.s)
}

func (s *Storage) GetLast(l int) []sysmon.Stats {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.s) < l {
		return s.s[:]
	}
	return s.s[len(s.s)-l:]
}
