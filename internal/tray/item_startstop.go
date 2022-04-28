package tray

import (
	"context"
	"time"

	"github.com/getlantern/systray"
	"github.com/mrVanboy/gauges/internal/tray/assets"
)

const (
	stateStopped = iota
	stateRunning
	stateError
)

type startStopItem struct {
	*systray.MenuItem
	collector Runner
	stop      func()

	state     uint8
	lastError string

	rateGetter func() time.Duration
	portGetter func() string

	autostartSetter func(bool)

	validator func() bool
}

func (ss *startStopItem) init(_ context.Context) <-chan struct{} {
	ss.MenuItem = systray.AddMenuItem("Start/Stop", "")
	ss.ClickedCh = make(chan struct{}, 1)

	ss.refresh()

	return ss.ClickedCh
}

func (ss startStopItem) isValid() bool {
	return true
}

func (ss *startStopItem) refresh() {
	switch ss.state {
	case stateStopped:
		ss.SetTitle("Start")
		ss.SetTemplateIcon(assets.PlayIcon, assets.PlayIcon)
	case stateRunning:
		ss.SetTitle("Stop")
		ss.SetTemplateIcon(assets.PauseIcon, assets.PauseIcon)
	case stateError:
		ss.SetTitle("Error")
		ss.SetTemplateIcon(assets.ErrorIcon, assets.ErrorIcon)
	}
	if !ss.validator() {
		ss.Disable()
	} else {
		ss.Enable()
	}
}

func (ss *startStopItem) onClick() {
	defer ss.refresh()

	switch ss.state {
	case stateStopped:
		ss.onStartClick()
	case stateRunning:
		ss.onStopClick()
	case stateError:
		return
	}
}

func (ss *startStopItem) onStartClick() {
	ctx := context.Background()
	ctx, ss.stop = context.WithCancel(ctx)
	err := ss.collector.Run(ctx, ss.portGetter(), ss.rateGetter())
	if err != nil {
		ss.state = stateError
		ss.SetTooltip(ss.lastError)
		return
	}

	ss.state = stateRunning
	ss.autostartSetter(true)
}
func (ss *startStopItem) onStopClick() {
	ss.stop()
	ss.state = stateStopped
	ss.autostartSetter(false)
}
