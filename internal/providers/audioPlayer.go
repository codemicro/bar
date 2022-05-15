package providers

import (
	"fmt"
	"strings"

	"github.com/codemicro/bar/internal/i3bar"
)

const (
	musicNoteString     = "♪"
	playerctlExecutable = "playerctl"

	playerStatusStopped = "Stopped"
	playerStatusPlaying = "Playing"
	playerStatusPaused  = "Paused"
	playerStatusUnknown = "Unknown"
)

type AudioPlayer struct {
	ShowTextOnPause bool
	MaxLabelLen     int
}

func NewAudioPlayer(maxLabelLength int) *AudioPlayer {
	return &AudioPlayer{
		MaxLabelLen: maxLabelLength,
	}
}

type playingAudioInfo struct {
	Track  string
	Artist string
	Album  string
	Status string
}

func (g *AudioPlayer) getInfo() (*playingAudioInfo, error) {
	rawMetadataOutput, err := runCommand(playerctlExecutable, "metadata")
	if err != nil {
		return nil, err
	}

	info := new(playingAudioInfo)

	lines := strings.Split(string(rawMetadataOutput), "\n")
	for _, line := range lines {
		splitLine := strings.Fields(line)

		if len(splitLine) < 3 {
			continue
		}

		var (
			// application = splitLine[0]
			fieldName = splitLine[1]
			data      = strings.Join(splitLine[2:], " ")
		)

		switch strings.ToLower(fieldName) {
		case "xesam:artist":
			info.Artist = data
		case "xesam:title":
			info.Track = data
		case "xesam:album":
			info.Album = data
		}
	}

	rawStatusOutput, err := runCommand(playerctlExecutable, "status")
	if err != nil {
		return nil, err
	}

	if x := string(rawStatusOutput); !(x == playerStatusStopped || x == playerStatusPlaying || x == playerStatusPaused) {
		info.Status = playerStatusUnknown
	} else {
		info.Status = x
	}

	return info, nil
}

func (g *AudioPlayer) Block(colors *i3bar.ColorSet) (*i3bar.Block, error) {
	info, err := g.getInfo()
	if err != nil {
		return nil, err
	}

	b := new(i3bar.Block)
	b.Name = "audioPlayer"

	b.FullText = musicNoteString

	if info.Status == playerStatusPlaying || (info.Status == playerStatusPaused && g.ShowTextOnPause) {
		x := fmt.Sprintf("%s - %s", info.Track, info.Artist)
		if len(x) > g.MaxLabelLen {
			x = x[:g.MaxLabelLen] + "..."
		}
		b.FullText += " " + x
	}

	return b, nil
}
