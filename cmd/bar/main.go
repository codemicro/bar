package main

import (
	"os"
	"os/signal"
	"path"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/codemicro/bar/internal/i3bar"
	"github.com/codemicro/bar/internal/providers"
)

func main() {
	logFileName := "cdmbar.log"
	if userHomeDir, err := os.UserHomeDir(); err == nil {
		logFileName = path.Join(userHomeDir, logFileName)
	}

	log.Logger = log.Logger.Output(zerolog.MultiLevelWriter(os.Stderr, &lumberjack.Logger{
		Filename: logFileName,
		MaxSize:  1,  // MB
		MaxAge:   14, // days
	})).Level(zerolog.ErrorLevel)

	if err := run(); err != nil {
		log.Fatal().Err(err).Msg("unhandled error")
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
		providers.NewAudioPlayer(32),
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
		// hacky thing to force an update in a second so we can clear the "cdmbar <VERSION>" readout
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
