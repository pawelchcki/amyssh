package amyssh

import (
	"fmt"
	"math/rand"
	"time"
)

var _ = fmt.Println

func operation(cfg *Config) {
	time.Sleep(time.Duration(rand.Intn(130)) * time.Millisecond)
	// fmt.Printf("%d\n", time.Duration(rand.Intn(110))*time.Millisecond)
}

type sleepSchedule struct {
	interval      time.Duration
	intervalDelta time.Duration
}

func (s *sleepSchedule) adjustInterval(cfg *Config, duration time.Duration) {
	// adjust interval
	switch {
	case duration > cfg.BackoffThreshold:
		s.interval = cfg.MaxPollInterval
		s.intervalDelta = 100 * time.Millisecond

	case duration <= cfg.PerformanceThreshold:
		s.interval = s.interval - s.intervalDelta
		if s.intervalDelta < 20*time.Second {
			s.intervalDelta = s.intervalDelta * 2
		} else {
			s.intervalDelta = 20 * time.Second
		}

	case duration > cfg.PerformanceThreshold:
		s.interval = s.interval + s.intervalDelta
		if s.intervalDelta > 100*time.Millisecond {
			s.intervalDelta = s.intervalDelta / 2
		} else {
			s.intervalDelta = 100 * time.Millisecond
		}
	}
	if s.interval < cfg.MinPollInterval {
		s.interval = cfg.MinPollInterval
	}
	if s.interval > cfg.MaxPollInterval {
		s.interval = cfg.MaxPollInterval
	}
}

func DispatchLoop(cfg *Config, fn func(cfg *Config)) {
	s := sleepSchedule{cfg.MaxPollInterval, 100 * time.Millisecond}

	for {
		// Measure operation time
		start := time.Now()
		operation(cfg)
		duration := time.Since(start)
		s.adjustInterval(cfg, duration)

		// sleep
		fuzz := time.Duration(rand.Int63n(100)) * time.Millisecond
		time.Sleep(s.interval + fuzz)
	}
}
