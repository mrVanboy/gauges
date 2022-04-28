package device

import (
	"context"
	"fmt"
	"time"

	"github.com/mrVanboy/gauges/internal"
	cobs "github.com/mrVanboy/go-simple-cobs"
	"github.com/tarm/serial"
)

const baudRate = 115200

const (
	magicMemCPU byte = iota
	magicHandshake
)

type Collector interface {
	Collect(ctx context.Context, usages chan<- float32, rate time.Duration) error
}

type Gauge struct {
	p           *serial.Port
	Log         internal.Logger
	CPU, Memory Collector
}

func (g *Gauge) Run(ctx context.Context, port string, rate time.Duration) error {
	p, err := serial.OpenPort(&serial.Config{Name: port, Baud: baudRate})

	if err != nil {
		return fmt.Errorf("cannot run the device: %w", err)
	}
	g.p = p
	if err := g.write(magicHandshake); err != nil {
		return fmt.Errorf("cannot send handshake: %w", err)
	}

	cpus := make(chan float32, 1)
	mems := make(chan float32, 1)

	if err := g.CPU.Collect(ctx, cpus, rate); err != nil {
		return fmt.Errorf("cannot collect cpu: %w", err)
	}
	if err := g.Memory.Collect(ctx, mems, rate); err != nil {
		return fmt.Errorf("cannot collect memory: %w", err)
	}

	go func() {
		var c, m float32
		for {
			select {
			case <-ctx.Done():
				return
			case c = <-cpus:
				g.send(c, m)
			case m = <-mems:
				g.send(c, m)
			}
		}
	}()
	return nil
}

func (g *Gauge) send(cpu, mem float32) {
	var c byte = uint8(cpu * 100)
	var m byte = uint8(mem * 100)
	err := g.write(magicMemCPU, m, c)
	if err != nil {
		g.Log.Printf("cannot send usages cpu=%d mem=%d: %v", c, m, err)
	}
	return
}

func (g *Gauge) write(b ...byte) error {
	enc, err := cobs.Encode(b)
	if err != nil {
		return fmt.Errorf("cannot encode data %0x: %w", b, err)
	}
	enc = append(enc, 0x00)
	_, err = g.p.Write(enc)
	if err != nil {
		return fmt.Errorf("cannot send data %0x: %w", enc, err)
	}
	return nil
}
