package tray

import (
	"context"
	"sync"
	"time"

	"github.com/getlantern/systray"
	"github.com/mrVanboy/gauges/internal"
	"github.com/mrVanboy/gauges/internal/tray/assets"
)

var refreshRates = [...]time.Duration{
	time.Second / 30,
	time.Second / 10,
	time.Second / 5,
	time.Second / 2,
	time.Second,
	time.Second * 2,
	time.Second * 5,
	time.Second * 10,
}

var _ menuItemer = &refreshRateMenu{}

type refreshRateMenu struct {
	log internal.Logger
	mu  sync.Mutex
	checklist[time.Duration]
	rateReporter func(time.Duration)

	*systray.MenuItem
}

func (rr *refreshRateMenu) isValid() bool {
	rr.mu.Lock()
	defer rr.mu.Unlock()
	for _, item := range rr.items {
		if item.Checked() {
			return true
		}
	}
	return false
}

func (rr *refreshRateMenu) onClick() {
}

func (rr *refreshRateMenu) onItemClick(selectedItem *systray.MenuItem) {
	rr.checklist.onItemClick(selectedItem)
	select {
	case rr.ClickedCh <- struct{}{}:
		// notify parent about click too
	default:
	}
}

func (rr *refreshRateMenu) init(ctx context.Context) <-chan struct{} {
	if rr.log == nil {
		rr.log = internal.NoOpLogger{}
	}
	rr.MenuItem = systray.AddMenuItem("Refresh rate", "How often to pull the stats")
	rr.SetTemplateIcon(assets.RefreshRateIcon, assets.RefreshRateIcon)
	rr.ClickedCh = make(chan struct{}, 1)
	rr.refresh()
	rr.run(ctx)

	return rr.ClickedCh
}

func (r *refreshRateMenu) run(ctx context.Context) {
	go func() {
		for {
			for rate, item := range r.items {
				select {
				case <-ctx.Done():
					return
				case <-item.ClickedCh:
					r.log.Printf("clicked on %v", item)
					r.onItemClick(item)
					r.rateReporter(rate)
				default:
				}
			}
			time.Sleep(time.Second / 30)
		}
	}()
}

func (rr *refreshRateMenu) refresh() {
	rr.mu.Lock()
	defer rr.mu.Unlock()

	if rr.items == nil {
		rr.items = make(map[time.Duration]*systray.MenuItem, len(refreshRates))
	}
	rr.addRates()
}

func (rr *refreshRateMenu) addRates() {
	for _, v := range refreshRates {
		if _, ok := rr.items[v]; ok {
			continue
		}
		ri := rr.AddSubMenuItemCheckbox(
			v.Truncate(time.Millisecond*10).String(),
			"",
			false,
		)
		ri.ClickedCh = make(chan struct{}, 1)
		rr.items[v] = ri
	}
}
