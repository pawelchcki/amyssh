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

func DispatchLoop(cfg *Config) {
	interval := cfg.MaxPollInterval
	intervalDelta := 100 * time.Millisecond

	for {
		// Measure operation time
		start := time.Now()
		operation(cfg)
		duration := time.Since(start)

		// adjust interval
		switch {
		case duration > cfg.BackoffThreshold:
			interval = cfg.MaxPollInterval
			intervalDelta = 100 * time.Millisecond

		case duration <= cfg.PerformanceThreshold:
			interval = interval - intervalDelta
			if intervalDelta < 20*time.Second {
				intervalDelta = intervalDelta * 2
			} else {
				intervalDelta = 20 * time.Second
			}

		case duration > cfg.PerformanceThreshold:
			interval = interval + intervalDelta
			if intervalDelta > 100*time.Millisecond {
				intervalDelta = intervalDelta / 2
			} else {
				intervalDelta = 100 * time.Millisecond
			}
		}
		if interval < cfg.MinPollInterval {
			interval = cfg.MinPollInterval
		}
		if interval > cfg.MaxPollInterval {
			interval = cfg.MaxPollInterval
		}

		// sleep
		fuzz := time.Duration(rand.Int63n(100)) * time.Millisecond
		time.Sleep(interval + fuzz)
	}
}
