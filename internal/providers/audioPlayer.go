package providers

import (
	"fmt"
	"strings"
	"time"

	"github.com/codemicro/bar/internal/i3bar"
)

const (
	musicNoteString     = "♪"
	pausedIconString    = "⏸"
	playerctlExecutable = "playerctl"

	playerStatusStopped = "Stopped"
	playerStatusPlaying = "Playing"
	playerStatusPaused  = "Paused"
	playerStatusUnknown = "Unknown"
)

type AudioPlayer struct {
	ShowTextOnPause bool
	MaxLabelLen     int
	TickerSteps     int

	name string

	lastText     string
	tickerCursor int
}

func NewAudioPlayer(maxLabelLength int) *AudioPlayer {
	return &AudioPlayer{
		MaxLabelLen: maxLabelLength,
		TickerSteps: 3,
		name:        "audioPlayer",
	}
}

func (g *AudioPlayer) Frequency() uint8 {
	return 1
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
		// If there's no player open, an error will be thrown by this command
		// with the below output
		if string(rawMetadataOutput) == "No players found" {
			return &playingAudioInfo{
				Status: playerStatusUnknown,
			}, nil
		}
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

func (g *AudioPlayer) AnimateTicker(x string) string {
	asRunes := []rune(x)

	if len(asRunes) <= g.MaxLabelLen {
		g.lastText = string(asRunes)
		return string(asRunes)
	}
	mod := append(asRunes, []rune("    ")...)

	if sm := string(mod); sm != g.lastText {
		g.tickerCursor = 0
		g.lastText = sm
		return string(mod[:g.MaxLabelLen])
	}

	g.tickerCursor += g.TickerSteps
	if l := len(mod); g.tickerCursor >= l {
		g.tickerCursor -= l
	}

	if g.tickerCursor+g.MaxLabelLen > len(mod) {
		diff := len(mod) - g.tickerCursor
		return string(mod[g.tickerCursor:]) + string(mod[:g.MaxLabelLen-diff])
	}

	return string(mod[g.tickerCursor : g.tickerCursor+g.MaxLabelLen])
}

func (g *AudioPlayer) Block(colors *i3bar.ColorSet) (*i3bar.Block, error) {
	info, err := g.getInfo()
	if err != nil {
		return nil, err
	}

	b := new(i3bar.Block)
	b.Name = g.name

	b.FullText = musicNoteString

	if info.Status == playerStatusPlaying || (info.Status == playerStatusPaused && g.ShowTextOnPause) {

		b.FullText += " "

		if info.Status == playerStatusPaused {
			b.FullText += pausedIconString + " "
		}

		b.FullText += g.AnimateTicker(fmt.Sprintf("%s - %s", info.Track, info.Artist))
	}

	return b, nil
}

func (g *AudioPlayer) GetNameAndInstance() (string, string) {
	return g.name, ""
}

func (g *AudioPlayer) OnClick(event *i3bar.ClickEvent) bool {
	switch event.Button {
	case i3bar.LeftMouseButton:
		_, _ = runCommand(playerctlExecutable, "play-pause")
	case i3bar.MouseWheelScrollUp:
		_, _ = runCommand(playerctlExecutable, "next")
	case i3bar.MouseWheelScrollDown:
		_, _ = runCommand(playerctlExecutable, "previous")
	default:
		return false
	}
	
	time.Sleep(time.Millisecond * 50)
	return true
}
