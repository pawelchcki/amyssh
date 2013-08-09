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
			s.intervalDelta = cfg.MaxPollInterval / 2
		}

	case duration > cfg.PerformanceThreshold:
		s.interval = s.interval + s.intervalDelta
		if s.intervalDelta > 100*time.Millisecond {
			s.intervalDelta = s.intervalDelta / 2
		} else {
			s.intervalDelta = cfg.MinPollInterval / 2
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
		intervalDelta: cfg.MaxPollInterval / 2,
	}

	for {
		// Measure operation time
		start := time.Now()
		err := fn(cfg)
		duration := time.Since(start)
		if err != nil {
			log.Printf("at interval %s operation took %s and returned error: <%+v>\n",
				s.interval.String(), duration.String(), err)
			s.interval = cfg.MaxPollInterval
		} else {
			s.adjustInterval(cfg, duration)
		}

		// sleep
		time.Sleep(s.interval + timeFuzz(s.interval/10))
	}
}
