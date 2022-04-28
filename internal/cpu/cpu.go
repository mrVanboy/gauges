package cpu

import (
	"context"
	"time"

	"github.com/mrVanboy/gauges/internal"
)

// ticker returns current state of the CPU ticks. Calculate the difference between to ticks and divide by difference of totals to get relative value.
//   oldIdle, oldTotal, _ := ticks()
//   newIdle, newTotal, _ := ticks()
//   idlePercentage := float64(newIdle - oldIdle)/float64(newTotal-oldTotal)
type ticker func() (idle uint64, total uint64, err error)

type CPU struct {
	ticker
	log internal.Logger
}

func (c *CPU) init() {
	if c.ticker == nil {
		c.ticker = ticks
	}
	if c.log == nil {
		c.log = internal.NoOpLogger{}
	}
}

func (c *CPU) Collect(ctx context.Context, usages chan<- float32, interval time.Duration) error {
	c.init()
	oldIdle, oldTotal, err := c.ticker()
	if err != nil {
		return err
	}

	ticker := time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
			}

			newIdle, newTotal, err := c.ticker()
			if err != nil {
				continue
			}
			idle, total := float32(newIdle-oldIdle), float32(newTotal-oldTotal)
			oldIdle, oldTotal = newIdle, newTotal

			var usage float32
			if total > 0 && idle > 0 {
				usage = 1 - idle/total
			}

			select {
			case <-ctx.Done():
				return
			case usages <- usage:
			}
		}
	}()
	return nil
}
