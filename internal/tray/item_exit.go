package tray

import (
	"context"

	"github.com/getlantern/systray"
	"github.com/mrVanboy/gauges/internal/tray/assets"
)

var _ menuItemer = &exitItem{}

type exitItem struct {
	*systray.MenuItem
}

func (i *exitItem) init(_ context.Context) <-chan struct{} {
	i.MenuItem = systray.AddMenuItem("Quit", "Quit the application")
	i.SetTemplateIcon(assets.ExitIcon, assets.ExitIcon)
	i.ClickedCh = make(chan struct{}, 1)
	return i.ClickedCh
}

func (exitItem) isValid() bool {
	return true
}

func (exitItem) refresh() {}

func (exitItem) onClick() {
	systray.Quit()
}
