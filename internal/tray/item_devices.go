package tray

import (
	"context"
	"sync"
	"time"

	"github.com/getlantern/systray"
	"github.com/mrVanboy/gauges/internal"
	"github.com/mrVanboy/gauges/internal/device"
	"github.com/mrVanboy/gauges/internal/tray/assets"
)

var _ menuItemer = &devicesMenu{}

type devicesMenu struct {
	deviceLister func() ([]string, error)
	log          internal.Logger
	mu           sync.Mutex
	portReporter func(string)

	*systray.MenuItem
	checklist[string]
}

func (dm *devicesMenu) init(ctx context.Context) <-chan struct{} {
	if dm.log == nil {
		dm.log = internal.NoOpLogger{}
	}
	if dm.deviceLister == nil {
		dm.deviceLister = device.List
	}

	dm.MenuItem = systray.AddMenuItem("Device", "Click to refresh the list")
	dm.SetTemplateIcon(assets.DevicesIcon, assets.DevicesIcon)
	dm.ClickedCh = make(chan struct{}, 1)

	dm.refresh()
	dm.run(ctx)
	return dm.ClickedCh
}

func (dm *devicesMenu) onClick() {
	dm.refresh()
}

func (dm *devicesMenu) onItemClick(selectedItem *systray.MenuItem) {
	dm.checklist.onItemClick(selectedItem)
	select {
	case dm.ClickedCh <- struct{}{}:
		// notify parent about click too
	default:
	}
}

func (d *devicesMenu) run(ctx context.Context) {
	go func() {
		for {
			d.mu.Lock()
			for port, item := range d.items {
				select {
				case <-ctx.Done():
					d.mu.Unlock()
					return
				case <-item.ClickedCh:
					d.log.Printf("clicked on %v", item)
					d.onItemClick(item)
					d.portReporter(port)
				default:
				}
			}
			d.mu.Unlock()
			time.Sleep(time.Second / 30)
		}
	}()
}

func (d *devicesMenu) refresh() {
	d.mu.Lock()
	defer d.mu.Unlock()
	paths, err := d.deviceLister()
	if err != nil {
		d.log.Printf("cannot build devices submenu: %w")
	}

	if d.items == nil {
		d.items = make(map[string]*systray.MenuItem, len(paths))
	}
	for _, path := range paths {
		if i, ok := d.items[path]; ok {
			i.Enable()
			continue
		}
		mi := d.MenuItem.AddSubMenuItemCheckbox(path, "", false)
		mi.ClickedCh = make(chan struct{}, 1)
		d.items[path] = mi
	}
items:
	for oldPath := range d.items {
		for _, newPath := range paths {
			if newPath == oldPath {
				continue items
			}
		}
		d.items[oldPath].Disable()
	}
}

func (d *devicesMenu) isValid() bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	for _, item := range d.items {
		if item.Checked() && !item.Disabled() {
			return true
		}
	}
	return false
}
