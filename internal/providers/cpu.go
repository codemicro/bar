package providers

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/codemicro/bar/internal/i3bar"
)

type CPU struct {
	OkThreshold      float32
	WarningThreshold float32

	idle0, total0 uint64
	idle1, total1 uint64
}

func NewCPU(okThreshold, warningThreshold float32) i3bar.BlockGenerator {
	m := &CPU{
		OkThreshold:      okThreshold,
		WarningThreshold: warningThreshold,
	}
	_ = m.doSample()
	return m
}

func (g *CPU) doSample() error {
	contents, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		return err
	}

	g.idle0 = g.idle1
	g.total0 = g.total1

	g.idle1 = 0
	g.total1 = 0

	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if fields[0] == "cpu" {
			numFields := len(fields)
			for i := 1; i < numFields; i++ {
				val, err := strconv.ParseUint(fields[i], 10, 64)
				if err != nil {
					return err
				}
				g.total1 += val // tally up all the numbers to get total ticks
				if i == 4 {     // idle is the 5th field in the cpu line
					g.idle1 = val
				}
			}
			return nil
		}
	}

	return errors.New("no CPU field")
}

func (g *CPU) getPercentage() float32 {
	idleTicks := float64(g.idle1 - g.idle0)
	totalTicks := float64(g.total1 - g.total0)
	return float32(100 * (totalTicks - idleTicks) / totalTicks)
}

func (g *CPU) Block(colors *i3bar.ColorSet) (*i3bar.Block, error) {
	if err := g.doSample(); err != nil {
		return nil, err
	}
	p := g.getPercentage()

	block := &i3bar.Block{
		Name:      "cpu",
		FullText:  fmt.Sprintf("CPU: %.1f%%", p),
		ShortText: fmt.Sprintf("C: %.1f%%", p),
	}

	if p > g.WarningThreshold && g.WarningThreshold != 0 {
		block.TextColor = colors.Bad
	} else if p > g.OkThreshold && g.OkThreshold != 0 {
		block.TextColor = colors.Warning
	}

	return block, nil
}
