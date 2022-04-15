package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"
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

	commitHash := getCommitHash()
	if commitHash != "" {
		commitHash = " " + commitHash
	}

	if err := b.Emit([]i3bar.BlockGenerator{
		providers.NewPlainText("cdmbar" + commitHash),
	}); err != nil {
		return err
	}

	blocks := []i3bar.BlockGenerator{
		providers.NewIPAddress("wlp0s20f3"),
		providers.NewWiFi("wlp0s20f3", 75),
		providers.NewBattery("BAT0", 80, 30, 20),
		providers.NewDisk("/", 30, 10),
		providers.NewCPU(20, 50),
		providers.NewMemory(7, 5),
		providers.NewPulseaudioVolume(),
		providers.NewDateTime(),
	}

	ticker := time.NewTicker(time.Second * 5)
	sigUpdate := make(chan os.Signal, 1)

	signal.Notify(sigUpdate, syscall.SIGUSR1)

	go func() {
		time.Sleep(time.Millisecond * 1000)
		sigUpdate <- os.Signal(syscall.SIGUSR1)
	}()

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

// TODO: Spotify provider!

func getCommitHash() string {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return ""
	}

	var hash string

	for _, setting := range bi.Settings {
		if setting.Key == "vcs.revision" {
			hash = setting.Value[0:7]
			break
		}
	}

	return hash
}