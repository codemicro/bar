package providers

import (
	"fmt"
	"io/ioutil"
	"path"
	"strconv"
	"strings"

	"github.com/codemicro/bar/internal/i3bar"
)

const (
	batteryStateFull        = "FULL"
	batteryStateDischarging = "BAT"
	batteryStateCharging    = "CHR"
	batteryStateUnknown     = "UNK"
)

type Battery struct {
	FullThreshold    float32
	OkThreshold      float32
	WarningThreshold float32

	DeviceName         string
	UseDesignMaxEnergy bool

	name                         string
	previousWasBackgroundWarning bool
	isAlert                      bool
}

func NewBattery(deviceName string, fullThreshold, okThreshold, warningThreshold float32) i3bar.BlockGenerator {
	return &Battery{
		DeviceName:       deviceName,
		FullThreshold:    fullThreshold,
		OkThreshold:      okThreshold,
		WarningThreshold: warningThreshold,
		name:             "battery",
	}
}

func (g *Battery) Frequency() uint8 {
	if g.isAlert {
		return 1
	}
	return 5
}

func (g *Battery) infoPath() string {
	return path.Join("/sys/class/power_supply", g.DeviceName)
}

func (g *Battery) getPercentage() (float32, error) {
	// TODO: Cache the read of energy_full
	var maxEnergy, energyNow int32
	{
		maxFile := "energy_full"
		if g.UseDesignMaxEnergy {
			maxFile = "energy_full_design"
		}

		me, err := ioutil.ReadFile(path.Join(g.infoPath(), maxFile))
		if err != nil {
			return 0, err
		}

		en, err := ioutil.ReadFile(path.Join(g.infoPath(), "energy_now"))
		if err != nil {
			return 0, err
		}

		mei, err := strconv.ParseInt(strings.TrimSpace(string(me)), 10, 32)
		if err != nil {
			return 0, err
		}

		eni, err := strconv.ParseInt(strings.TrimSpace(string(en)), 10, 32)
		if err != nil {
			return 0, err
		}

		maxEnergy = int32(mei)
		energyNow = int32(eni)
	}

	return (float32(energyNow) / float32(maxEnergy)) * 100, nil
}

func (g *Battery) getState() (string, error) {

	sa, err := ioutil.ReadFile(path.Join(g.infoPath(), "status"))
	if err != nil {
		return "", err
	}

	var x string

	switch strings.TrimSpace(string(sa)) {
	case "Full":
		x = batteryStateFull
	case "Discharging":
		x = batteryStateDischarging
	case "Charging":
		x = batteryStateCharging
	case "Unknown":
		fallthrough
	default:
		x = batteryStateUnknown
	}

	return x, nil
}

func (g *Battery) Block(colors *i3bar.ColorSet) (*i3bar.Block, error) {
	percentage, err := g.getPercentage()
	if err != nil {
		return nil, err
	}

	state, err := g.getState()
	if err != nil {
		return nil, err
	}

	block := &i3bar.Block{
		Name:      g.name,
		Instance:  g.DeviceName,
		FullText:  fmt.Sprintf("%s %.1f%%", state, percentage),
		ShortText: fmt.Sprintf("%.1f%%", percentage),
	}

	if percentage < g.WarningThreshold && g.WarningThreshold != 0 {

		g.isAlert = true

		if g.previousWasBackgroundWarning || state == batteryStateCharging { // disable flashing when on charge
			block.TextColor = colors.Bad
		} else {
			block.BackgroundColor = colors.Bad
		}

		g.previousWasBackgroundWarning = !g.previousWasBackgroundWarning

	} else if percentage < g.OkThreshold && g.OkThreshold != 0 {
		block.TextColor = colors.Warning
	} else {
		g.isAlert = false
	}

	switch state {
	case batteryStateCharging:
		if percentage > g.FullThreshold && g.FullThreshold != 0 {
			block.TextColor = colors.Good
		} else {
			// Set text/background color to white
			block.BackgroundColor = nil
			block.TextColor = nil
		}
	case batteryStateFull:
		block.BackgroundColor = colors.Warning
		block.TextColor = colors.Background
	case batteryStateUnknown:
		block.TextColor = colors.Warning
	}

	return block, nil
}

func (g *Battery) GetNameAndInstance() (string, string) {
	return g.name, g.DeviceName
}
