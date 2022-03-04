package main

import (
	"fmt"
	"os"
	"time"

	"github.com/codemicro/bar/internal/i3bar"
	"github.com/codemicro/bar/internal/providers"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "Unhandled error:", err)
		os.Exit(1)
	}
}

func run() error {
	b := i3bar.New(os.Stdout)
	if err := b.Initialise(); err != nil {
		return err
	}

	blocks := []i3bar.BlockGenerator{
		providers.NewBattery("BAT0", 30, 15),
		providers.NewDisk("/", 30, 10),
		providers.NewCPU(10, 20),
		providers.NewMemory(),
		providers.NewPulseaudioVolume(),
		providers.NewDateTime(),
	}

	if err := b.Emit(blocks); err != nil {
		return err
	}

	ticker := time.NewTicker(time.Second * 5)

	for {
		select {
		case <-ticker.C:
			err := b.Emit(blocks)
			if err != nil {
				return err
			}
		}
	}
}

// TODO: Accept signals to refresh
