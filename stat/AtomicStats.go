package stat

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"
)

type AtomicStats struct {
	mu    sync.RWMutex
	stats map[string]*AtomicStat
}

func NewAtomicStats() *AtomicStats {
	return &AtomicStats{}
}

func (s *AtomicStats) AddTime(key string, start time.Time) time.Time {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.stats == nil {
		s.stats = make(map[string]*AtomicStat)
	}
	stat := s.stats[key]

	now := time.Now()

	duration := now.Sub(start)
	stat.AddDuration(duration)

	return now
}

func (s *AtomicStats) AddDuration(key string, duration time.Duration) *AtomicStats {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.stats == nil {
		s.stats = make(map[string]*AtomicStat)
	}
	stat := s.stats[key]

	if stat == nil {
		stat = NewAtomicStat()
		s.stats[key] = stat
	}

	stat.AddDuration(duration)

	return s
}

func (s *AtomicStats) GetCount() int32 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.stats == nil {
		return 0
	}
	var count int32
	for _, stat := range s.stats {
		count += stat.count
	}
	return count
}

func (s *AtomicStats) GetMap() map[string]*AtomicStat {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.stats
}

func (s *AtomicStats) String(indent string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.stats == nil {
		return "0"
	}

	var keys []string
	for key := range s.stats {
		keys = append(keys, key)
	}

	sort.StringSlice(keys).Sort()

	var sb strings.Builder
	var totalCount int32
	var totalDuration time.Duration

	for _, key := range keys {
		stat := s.stats[key]

		totalCount += stat.count
		totalDuration += stat.duration

		if sb.Len() > 0 {
			sb.WriteString(indent)
		}
		sb.WriteString(fmt.Sprintf("%20s: %5d (%s)\n", key, stat.count, stat.GetAverageString()))
	}

	avgDuration := totalDuration / time.Duration(totalCount)
	sb.WriteString(indent)
	sb.WriteString(fmt.Sprintf("%20s: %5d (%s)\n", "TOTAL", totalCount, avgDuration.String()))

	return sb.String()
}

func (s *AtomicStats) GetTotalDuration() time.Duration {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.stats == nil {
		return 0
	}

	var totalDuration time.Duration

	for _, stat := range s.stats {
		totalDuration += stat.duration
	}

	return totalDuration
}

func (s *AtomicStats) GetTotalDurationAndCount() (time.Duration, int32) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.stats == nil {
		return 0, 0
	}

	var totalDuration time.Duration
	var totalCount int32

	for _, stat := range s.stats {
		totalDuration += stat.duration
		totalCount += stat.count
	}

	if totalCount == 0 {
		return 0, 0
	}

	return totalDuration, totalCount
}

func (s *AtomicStats) GetAverageDuration() time.Duration {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.stats == nil {
		return 0
	}

	var totalDuration time.Duration
	var totalCount int32

	for _, stat := range s.stats {
		totalDuration += stat.duration
		totalCount += stat.count
	}

	if totalCount == 0 {
		return 0
	}

	return totalDuration / time.Duration(totalCount)
}
