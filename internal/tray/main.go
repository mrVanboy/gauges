package tray

import (
	"context"
	"time"

	"github.com/getlantern/systray"
	"github.com/mrVanboy/gauges/internal"
	"github.com/mrVanboy/gauges/internal/cfg"
	"github.com/mrVanboy/gauges/internal/tray/assets"
)

type Runner interface {
	Run(ctx context.Context, port string, rate time.Duration) error
}

type menuItemer interface {
	init(context.Context) <-chan struct{}
	isValid() bool
	refresh()
	onClick()
}

type MainMenu struct {
	Log       internal.Logger
	Device    Runner
	items     map[<-chan struct{}]menuItemer
	port      lock[string]
	rate      lock[time.Duration]
	autostart lock[bool]
}

func (m *MainMenu) init(ctx context.Context) {
	if m.Log == nil {
		m.Log = internal.NoOpLogger{}
	}
	systray.SetTemplateIcon(assets.TrayIcon, assets.TrayIcon)
	m.items = make(map[<-chan struct{}]menuItemer)

	ss := &startStopItem{
		collector:       m.Device,
		portGetter:      m.port.get,
		rateGetter:      m.rate.get,
		validator:       m.isValid,
		autostartSetter: m.autostart.set,
	}
	defer func() {
		m.refresh()
		if m.autostart.get() && m.isValid() {
			ss.onStartClick()
			ss.refresh()
		}
	}()
	m.add(ctx, ss)

	systray.AddSeparator()

	rr := &refreshRateMenu{log: m.Log, rateReporter: m.rate.set}
	m.add(ctx, rr)
	preSelect(&m.rate, &rr.checklist)

	dm := &devicesMenu{log: m.Log, portReporter: m.port.set}
	m.add(ctx, dm)
	preSelect(&m.port, &dm.checklist)

	systray.AddSeparator()

	m.add(ctx, &exitItem{})

	m.Log.Printf("Menu initialized with %d items", len(m.items))
}

func (m *MainMenu) add(ctx context.Context, item menuItemer) {
	clicks := item.init(ctx)
	m.items[clicks] = item
}

func (m *MainMenu) Run(ctx context.Context, config cfg.Config) {
	go func() {
		for {
			for ch, item := range m.items {
				select {
				case <-ctx.Done():
					return
				case <-ch:
					m.Log.Printf("Clicked on %v", item)
					item.onClick()
					m.refresh()
				default:
				}
			}
			time.Sleep(time.Second / 30)
		}
	}()

	// load value from config, need to be injeted in submenus later during m.init
	m.port.set(config.Port)
	m.rate.set(config.RefreshRate)
	m.autostart.set(config.Autostart)

	systray.Run(
		func() {
			m.init(ctx)
		},
		func() {
			config.Port = m.port.get()
			config.RefreshRate = m.rate.get()
			config.Autostart = m.autostart.get()

			if err := config.Save(); err != nil {
				m.Log.Print(err.Error())
			}
		})
}
func (m *MainMenu) Stop() {
	systray.Quit()
}

func (m *MainMenu) isValid() bool {
	for _, item := range m.items {
		if !item.isValid() {
			m.Log.Printf("Item %v is not valid", item)
			return false
		}
	}
	return true
}

func (m *MainMenu) refresh() {
	for _, item := range m.items {
		item.refresh()
	}
}

func preSelect[K comparable](v *lock[K], k *checklist[K]) {
	k.selectDefault(v.get())
}
