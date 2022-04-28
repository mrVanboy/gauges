package memory

import (
	"context"
	"time"

	"github.com/mrVanboy/gauges/internal"
)

// pager returns available and total number of memory pages. Use unix.Getpagesize() from golang.org/x/sys/unix to convert to the bytes.
// Values are not 100% accurate, but for relative calculationas it's enough.
type pager func() (avail, total uint64, err error)

type Memory struct {
	pager
	log internal.Logger
}

func (m *Memory) init() {
	if m.pager == nil {
		m.pager = pages
	}
	if m.log == nil {
		m.log = internal.NoOpLogger{}
	}
}

func (m *Memory) Collect(ctx context.Context, usages chan<- float32, interval time.Duration) error {
	m.init()
	// to return immediately if we are not able to collect memory pages.
	_, _, err := m.pager()
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

			avail, total, err := m.pager()
			if err != nil {
				m.log.Print("cannot collect pages:", err.Error())
				continue
			}

			var usage float32
			if avail > total {
				continue
			}
			if total > 0 && avail > 0 {
				usage = 1 - float32(avail)/float32(total)
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
