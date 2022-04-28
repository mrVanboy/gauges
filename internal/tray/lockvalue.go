package tray

import "sync"

type lock[T comparable] struct {
	sync.Mutex
	v T
}

func (l *lock[T]) get() T {
	l.Lock()
	defer l.Unlock()
	return l.v
}

func (l *lock[T]) set(v T) {
	l.Lock()
	l.v = v
	l.Unlock()
}
