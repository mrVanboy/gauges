package cfg

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

var configPath = func() string {
	const configName = ".config.json"

	s, err := os.Executable()
	if err != nil {
		panic(err)
	}

	s, err = filepath.EvalSymlinks(s)
	if err != nil {
		panic(err)
	}

	s = filepath.Dir(s)
	return filepath.Join(s, configName)
}()

type Config struct {
	Port        string
	RefreshRate time.Duration
	Autostart   bool
}

func Default() Config {
	return Config{
		Port:        "",
		RefreshRate: 1 * time.Second,
	}
}

func Load() (Config, error) {
	c := Config{}
	b, err := os.ReadFile(configPath)
	switch {
	case errors.Is(err, os.ErrNotExist):
		_ = c.Save()
		fallthrough
	case err != nil:
		return Default(), fmt.Errorf("cannot load config: %w", err)
	}

	if err = json.Unmarshal(b, &c); err != nil {
		return Default(), fmt.Errorf("cannot parse config: %w", err)
	}
	return c, nil
}

func (c Config) Save() error {
	b, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("cannot serialize config: %w", err)
	}
	err = os.WriteFile(configPath, b, 0766)
	if err != nil {
		return fmt.Errorf("cannot save config file: %w", err)
	}
	return nil
}
