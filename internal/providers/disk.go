package providers

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/codemicro/bar/internal/i3bar"
)

type Disk struct {
	OkThreshold      float32
	WarningThreshold float32

	MountPath string
}

func NewDisk(mountPath string, okThreshold, warningThreshold float32) i3bar.BlockGenerator {
	return &Disk{
		OkThreshold:      okThreshold,
		WarningThreshold: warningThreshold,
		MountPath:        mountPath,
	}
}

func (g *Disk) getAvailable() (float32, error) {
	cmdout, err := runCommand("df")
	if err != nil {
		return 0, err
	}
	for _, line := range strings.Split(string(cmdout), "\n") {
		fields := strings.Fields(line)
		if fields[5] == g.MountPath || (g.MountPath == "" && fields[5] == "/") {
			y, _ := strconv.ParseFloat(fields[3], 64)
			return float32(y / 1000 / 1000), nil // to GB
		}
	}
	return 0, errors.New("could not find specified mounted drive")
}

func (g *Disk) Block(colors *i3bar.ColorSet) (*i3bar.Block, error) {
	da, err := g.getAvailable()
	if err != nil {
		return nil, err
	}

	block := &i3bar.Block{
		Name:      "disk",
		FullText:  fmt.Sprintf("Disk avail: %.1fGB", da),
		ShortText: fmt.Sprintf("D: %.1fGB", da),
	}

	if da < g.WarningThreshold && g.WarningThreshold != 0 {
		block.TextColor = colors.Bad
	} else if da < g.OkThreshold && g.OkThreshold != 0 {
		block.TextColor = colors.Warning
	}

	return block, nil
}
