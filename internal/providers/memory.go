package providers

import (
	"errors"
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"

	"github.com/codemicro/bar/internal/i3bar"
)

type Memory struct {
	OkThreshold float32
}

func NewMemory() i3bar.BlockGenerator {
	return new(Memory)
}

var (
	memoryTotalRegexp     = regexp.MustCompile(`MemTotal: +(\d+) kB`)
	memoryAvailableRegexp = regexp.MustCompile(`MemAvailable: +(\d+) kB`)
)

func (g *Memory) getStats() (used float32, total float32, err error) {
	fcont, err := ioutil.ReadFile("/proc/meminfo")
	if err != nil {
		return 0, 0, err
	}

	if x := memoryTotalRegexp.FindSubmatch(fcont); len(x) == 2 {
		totalKB, err := strconv.ParseInt(string(x[1]), 10, 64)
		if err != nil {
			return 0, 0, err
		}
		total = float32(totalKB) / float32(1000*1000)
	} else {
		return 0, 0, errors.New("could not fetch total system memory")
	}

	var available float32
	if x := memoryAvailableRegexp.FindSubmatch(fcont); len(x) == 2 {
		availableKB, err := strconv.ParseInt(string(x[1]), 10, 64)
		if err != nil {
			return 0, 0, err
		}
		available = float32(availableKB) / float32(1000*1000)
	} else {
		return 0, 0, errors.New("could not fetch available system memory")
	}

	used = total - available

	return
}

func (g *Memory) Block(colors *i3bar.ColorSet) (*i3bar.Block, error) {
	used, total, err := g.getStats()
	if err != nil {
		return nil, err
	}

	// TODO: Make the readout change between KB/MB/GB

	block := &i3bar.Block{
		Name:      "memory",
		FullText:  fmt.Sprintf("Mem: %.1fGB of %.1fGB", used, total),
		ShortText: fmt.Sprintf("M: %.1fGB", used),
	}

	if total - used < g.OkThreshold && g.OkThreshold != 0 {
		block.TextColor = colors.Bad
	}

	return block, nil
}