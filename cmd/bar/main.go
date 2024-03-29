package main

import (
	"os"
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
	})).Level(zerolog.DebugLevel)

	if err := run(); err != nil {
		log.Fatal().Err(err).Msg("unhandled error")
	}
}

func run() error {
	b := i3bar.New(os.Stdout, os.Stdin, syscall.SIGUSR1)
	if err := b.Initialise(); err != nil {
		return err
	}

	commitHash := getCommitHash()
	if commitHash != "" {
		commitHash = " " + commitHash
	}

	// Blocks registered first will be the rightmost in the status bar.
	b.RegisterBlockGenerator(
		providers.NewLaunchProgram("MINI", "/home/akp/.local/bin/minisettings"),
		providers.NewDateTime(),
		providers.NewPulseaudioVolume(),
		providers.NewMemory(7, 5),
		providers.NewCPU(20, 50),
		providers.NewDisk("/", 30, 10),
		providers.NewBattery("BAT0", 80, 30, 20),
		providers.NewWiFi("wlp0s20f3", 75),
		providers.NewIPAddress("wlp0s20f3"),
		providers.NewAudioPlayer(32),
	)

	if err := b.Emit([]*i3bar.Block{
		{FullText: "cdmbar" + commitHash},
	}); err != nil {
		return err
	}

	time.Sleep(time.Second) // show "cdmbar" for one second

	return b.StartLoop()
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
