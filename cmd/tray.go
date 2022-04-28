package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"log"

	"github.com/mrVanboy/gauges/internal"
	"github.com/mrVanboy/gauges/internal/cfg"
	"github.com/mrVanboy/gauges/internal/cpu"
	"github.com/mrVanboy/gauges/internal/device"
	"github.com/mrVanboy/gauges/internal/memory"
	"github.com/mrVanboy/gauges/internal/tray"
	"golang.org/x/sys/unix"
)

func main() {

	var logger internal.Logger
	switch _, err := unix.IoctlGetWinsize(int(os.Stdout.Fd()), unix.TIOCGWINSZ); err {
	case nil:
		logger = log.Default()
	default:
		logger = internal.NoOpLogger{}
	}

	config, err := cfg.Load()
	if err != nil {
		logger.Print(err.Error())
	}

	menu := &tray.MainMenu{
		Log: logger,
		Device: &device.Gauge{
			Log:    logger,
			CPU:    &cpu.CPU{},
			Memory: &memory.Memory{},
		},
	}
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	go func() {
		<-ctx.Done()
		menu.Stop()
	}()
	menu.Run(ctx, config)
}
