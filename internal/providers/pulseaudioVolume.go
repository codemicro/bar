package providers

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/codemicro/bar/internal/i3bar"
)

type PulseaudioVolume struct {
	// Sink is the target sink name to look for in Pulseaudio. Leave blank
	// to use the default sink.
	Sink string
}

func NewPulseaudioVolume() i3bar.BlockGenerator {
	return new(PulseaudioVolume)
}

func (g *PulseaudioVolume) getInfo() (string, error) {
	x, err := runCommand("pacmd", "info")
	return string(x), err
}

var (
	pulseaudioInfoParseRegexp = regexp.MustCompile(`(\d+ sink\(s\) available\.|\d+ source\(s\) available\.)`)
	pulseaudioVolumeRegexp    = regexp.MustCompile(`volume: front-left: \d+ \/ +(\d{1,3})% \/ (?:\d|.)+ dB, +front-right: \d+ \/ +(\d{1,3})% \/ (?:\d|.)+ dB`)
	pulseaudioMutedRegexp     = regexp.MustCompile(`muted: (.+)`)
	pulseaudioNameRegexp      = regexp.MustCompile(`name: <(.+)>`)
)

type volumeInfo struct {
	Left  int
	Right int
	Muted bool
}

func (g *PulseaudioVolume) getVolume(info string) (*volumeInfo, error) {
	x := pulseaudioInfoParseRegexp.Split(info, 5)

	sinkInfoStrings := strings.Split(
		strings.TrimLeft(x[1], "\n\t* "),
		"index: ",
	)[1:] // 1 onwards to remove the blank string at the beginning

	for _, sinkText := range sinkInfoStrings {
		v := new(volumeInfo)
		name := pulseaudioNameRegexp.FindStringSubmatch(sinkText)
		if len(name) == 2 && (strings.EqualFold(name[1], g.Sink) || g.Sink == "") {
			volumes := pulseaudioVolumeRegexp.FindStringSubmatch(sinkText)
			muted := pulseaudioMutedRegexp.FindStringSubmatch(sinkText)[1]

			if len(volumes) != 3 {
				return nil, fmt.Errorf("could not parse volumes from sink %s", name[1])
			}

			vl, err := strconv.ParseInt(volumes[1], 10, 32)
			if err != nil {
				return nil, err
			}
			v.Left = int(vl)

			vr, err := strconv.ParseInt(volumes[2], 10, 32)
			if err != nil {
				return nil, err
			}
			v.Right = int(vr)

			switch strings.ToLower(muted) {
			case "yes":
				v.Muted = true
			case "no":
				v.Muted = false
			default:
				return nil, fmt.Errorf("unknown error state %#v", muted)
			}

			return v, nil
		}
	}

	return nil, errors.New("no sink with the specified name found")
}

func (g *PulseaudioVolume) Block(colors *i3bar.ColorSet) (*i3bar.Block, error) {
	info, err := g.getInfo()
	if err != nil {
		return nil, err
	}

	v, err := g.getVolume(info)
	if err != nil {
		return nil, err
	}

	block := new(i3bar.Block)
	block.Name = "pulseaudioVolume"
	block.Instance = g.Sink

	if v.Muted {
		block.FullText = "Vol: muted"
		block.ShortText = "V: mute"
		block.TextColor = colors.Warning
	} else if v.Left == v.Right {
		block.FullText = fmt.Sprintf("Vol: %d%%", v.Left)
		block.ShortText = fmt.Sprintf("V: %d%%", v.Left)
	} else {
		block.FullText = fmt.Sprintf("Vol: L%d%% R%d%%", v.Left, v.Right)
		block.ShortText = fmt.Sprintf("V: %d%%", v.Left)
	}

	return block, nil
}
