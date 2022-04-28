package tray

import "github.com/getlantern/systray"

type checklist[K comparable] struct {
	items map[K]*systray.MenuItem
}

func (l *checklist[K]) onItemClick(selectedItem *systray.MenuItem) {
	for _, item := range l.items {
		switch item {
		case selectedItem:
			item.Check()
		default:
			item.Uncheck()
		}
	}
}

func (l *checklist[K]) selectDefault(defaultValue K) {
	for key, item := range l.items {
		switch key {
		case defaultValue:
			item.Check()
		default:
			item.Uncheck()
		}
	}
}
