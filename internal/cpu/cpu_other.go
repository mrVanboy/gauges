//go:build !(darwin && cgo)

package cpu

import (
	"fmt"
	"runtime"
)

var _ ticker = ticks

func ticks() (idle, total uint64, err error) {
	return 0, 0, fmt.Errorf("cpu metrics are not impleneted on %q", runtime.GOARCH)
}
