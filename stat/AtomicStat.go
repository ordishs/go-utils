package stat

import (
	"fmt"
	"sync"
	"time"
)

type AtomicStat struct {
	mu       sync.RWMutex
	count    int32
	duration time.Duration
}

func NewAtomicStat() *AtomicStat {
	return &AtomicStat{}
}

func (s *AtomicStat) AddDuration(duration time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.count += 1
	s.duration += duration
}

func (s *AtomicStat) AddTime(start time.Time) time.Time {
	now := time.Now()
	duration := now.Sub(start)

	s.AddDuration(duration)

	return now
}

func (s *AtomicStat) GetCount() int32 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.count
}

func (s *AtomicStat) GetDuration() time.Duration {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.duration
}

func (s *AtomicStat) GetAverage() float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.count == 0 {
		return 0
	}
	return float64(s.duration) / float64(s.count)
}

func (s *AtomicStat) GetAverageString() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.count == 0 {
		return "0"
	}
	avg := s.duration / time.Duration(s.count)
	return avg.String()
}

func (s *AtomicStat) String() string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.count == 0 {
		return "0"
	}

	return fmt.Sprintf("%d (%s)", s.count, s.GetAverageString())
}
