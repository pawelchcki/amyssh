package amyssh

import (
	"log"
	"math/rand"
	"time"
)

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

func timeFuzz(maxFuzz time.Duration) time.Duration {
	return time.Duration(rand.Int63n(int64(maxFuzz)))
}

func IntervalLoop(cfg *Config, fn func(cfg *Config) error) {
	s := sleepSchedule{
		interval:      cfg.MaxPollInterval,
		intervalDelta: 100 * time.Millisecond,
	}

	for {
		// Measure operation time
		start := time.Now()
		err := fn(cfg)
		duration := time.Since(start)

		if err != nil {
			log.Printf("operation returned error: %+v\n", err)
			s.interval = cfg.MaxPollInterval
		} else {
			s.adjustInterval(cfg, duration)
		}

		// sleep
		time.Sleep(s.interval + timeFuzz(100*time.Millisecond))
	}
}
