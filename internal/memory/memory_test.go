package memory

import (
	"context"
	"math"
	"testing"
	"time"
)

// testTicker returns always ticks with desiredLoad.
func testPager(desiredLoad float32) pager {
	total := uint64(math.MaxUint32) // we need huge number to cover rounding
	return func() (uint64, uint64, error) {
		total--
		avail := uint64(float32(total) * desiredLoad)
		return avail, total, nil
	}
}

func assertUsage(t *testing.T, want, got float32) {
	t.Helper()
	switch {
	case got == want:
	default:
		t.Errorf("invalid load: want %v, got: %v", want, got)
	}
}

func TestMemory_Collect(t *testing.T) {
	tests := map[string]struct {
		interval    time.Duration
		fixtureLoad float32 // load to return from testTicker
		wantLoad    float32
	}{
		"happy path": {
			interval:    time.Millisecond,
			fixtureLoad: 0.5,
			wantLoad:    0.5,
		},
		"negative - no panic": {
			interval:    time.Millisecond,
			fixtureLoad: -1.1,
		},
		"same - no panic": {
			interval:    time.Millisecond,
			fixtureLoad: math.MaxFloat32,
		},
	}

	for name := range tests {
		tt := tests[name]
		t.Run(name, func(t *testing.T) {
			// number of samples collected
			const samples = 5

			c := Memory{
				pager: testPager(tt.wantLoad),
			}
			usages := make(chan float32)
			ctx, cancel := context.WithTimeout(context.Background(), (1+samples)*tt.interval)
			err := c.Collect(ctx, usages, tt.interval)
			if err != nil {
				t.Fatal("cannot start cpu collection:", err)
			}
			for i := 0; i < samples; i++ {
				time.Sleep(tt.interval)
				select {
				case <-ctx.Done():
					t.Fatalf("only %d samples collected instead of %d", i, samples)
				case gotLoad := <-usages:
					assertUsage(t, tt.wantLoad, gotLoad)
				}
			}
			cancel()
			close(usages)
			// nothing should panic, like writing to closed channel etc.
			time.Sleep(tt.interval)
		})
	}
}
