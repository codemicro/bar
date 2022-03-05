package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
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
		providers.NewCPU(20, 50),
		providers.NewMemory(7, 5),
		providers.NewPulseaudioVolume(),
		providers.NewDateTime(),
	}

	if err := b.Emit(blocks); err != nil {
		return err
	}

	ticker := time.NewTicker(time.Second * 5)
	sigUpdate := make(chan os.Signal, 1)

	signal.Notify(sigUpdate, syscall.SIGUSR1)

	for {
		select {
		case <-sigUpdate:
			if err := b.Emit(blocks); err != nil {
				return err
			}
		case <-ticker.C:
			if err := b.Emit(blocks); err != nil {
				return err
			}
		}
	}
}

// TODO: Accept signals to refresh
// TODO: Spotify provider!
