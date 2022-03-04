package providers

import (
	"fmt"
	"io/ioutil"
	"path"
	"strconv"
	"strings"

	"github.com/codemicro/bar/internal/i3bar"
)

type Battery struct {
	OkThreshold      float32
	WarningThreshold float32

	Name               string
	UseDesignMaxEnergy bool
}

func NewBattery(name string, okThreshold, warningThreshold float32) i3bar.BlockGenerator {
	return &Battery{
		Name:             name,
		OkThreshold:      okThreshold,
		WarningThreshold: warningThreshold,
	}
}

func (g *Battery) infoPath() string {
	return path.Join("/sys/class/power_supply", g.Name)
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
		x = "FULL"
	case "Discharging":
		x = "BAT"
	case "Charging":
		x = "CHR"
	case "Unknown":
		fallthrough
	default:
		x = "UNK"
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
		Name:      "battery",
		Instance:  g.Name,
		FullText:  fmt.Sprintf("%s %.1f%%", state, percentage),
		ShortText: fmt.Sprintf("%.1f%%", percentage),
	}

	if percentage < g.WarningThreshold && g.WarningThreshold != 0 {
		block.TextColor = colors.Bad
	} else if percentage < g.OkThreshold && g.OkThreshold != 0 {
		block.TextColor = colors.Warning
	}

	return block, nil
}
