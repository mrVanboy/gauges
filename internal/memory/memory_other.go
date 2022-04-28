//go:build !(darwin && cgo)

package memory

import (
	"fmt"
	"runtime"
)

var _ pager = pages

func pages() (avail, total uint64, err error) {
	return 0, 0, fmt.Errorf("memory metrics are not impleneted on %q", runtime.GOARCH)
}
