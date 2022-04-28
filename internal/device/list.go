package device

import (
	"fmt"
	"path/filepath"
)

const ttyPathPattern = "/dev/tty.*"

func List() ([]string, error) {
	matches, err := filepath.Glob(ttyPathPattern)
	if err != nil {
		return nil, fmt.Errorf("cannot list the devices: %w", err)
	}
	return matches, nil
}
