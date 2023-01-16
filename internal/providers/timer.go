package providers

import (
	"fmt"
	"time"

	"github.com/codemicro/bar/internal/i3bar"
)

const (
	timerSymbolPause = "⏸"
	timerSymbolPlay  = "▶"
	timerSymbolClock = "⏰"
)

type Timer struct {
	UseShortLabel bool
	
	times []time.Time

	name string
}

func NewTimer(useShortLabel bool) i3bar.BlockGenerator {
	return &Timer{
		UseShortLabel: useShortLabel,
		name: "timer",
	}
}

func (g *Timer) Frequency() uint8 {
	return 1
}

func (g *Timer) OnClick(event *i3bar.ClickEvent) bool {
	resetButtonPressed := event.Button == i3bar.RightMouseButton
	triggerButtonPressed := event.Button == i3bar.LeftMouseButton

	numStoredTimes := len(g.times)

	if numStoredTimes == 0 {
		// start only if the left mouse button pressed
		if !triggerButtonPressed {
			return false
		}
		g.times = []time.Time{time.Now()}
	} else if resetButtonPressed {
		g.times = nil
	} else if triggerButtonPressed {
		// play/pause
		g.times = append(g.times, time.Now())
	}

	return true
}

func (g *Timer) calculateDuration() time.Duration {
	var sigma time.Duration
	for i := 0; i < len(g.times); i += 2 {
		next := time.Now()
		if i+1 < len(g.times) {
			next = g.times[i+1]
		}
		sigma += next.Sub(g.times[i])
	}
	return sigma.Round(time.Second)
}

func (g *Timer) Block(*i3bar.ColorSet) (*i3bar.Block, error) {
	block := &i3bar.Block{
		Name: g.name,
	}

	numStoredTimes := len(g.times)

	if numStoredTimes == 0 {
		if g.UseShortLabel {
			block.FullText = timerSymbolClock
			block.ShortText = timerSymbolClock
		} else {
			block.FullText = fmt.Sprintf("%s Click to start", timerSymbolClock)
			block.ShortText = fmt.Sprintf("%s Click", timerSymbolClock)
		}
	} else {
		symbol := timerSymbolPlay
		if numStoredTimes%2 == 0 {
			symbol = timerSymbolPause
		}
		block.FullText = fmt.Sprintf("%s %s %s", timerSymbolClock, symbol, g.calculateDuration())
	}

	return block, nil
}

func (g *Timer) GetNameAndInstance() (string, string) {
	return g.name, ""
}
