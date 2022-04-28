//go:build darwin && cgo

package memory

import (
	"testing"
)

func TestCollect(t *testing.T) {
	free, total, err := pages()
	switch {
	case err != nil:
		t.Fatalf("unexpected error: %q", err.Error())
	case total <= 0:
		t.Error("total memory should be more than zero")
	case total == free:
		t.Error("memory cannot be completely unused")
	case free > total:
		t.Error("free memory cannot be larger than total")
	}
}
