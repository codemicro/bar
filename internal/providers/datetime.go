package providers

import (
	"time"

	"github.com/codemicro/bar/internal/i3bar"
)

type DateTime struct {
	// TODO: 12 hour mode?
	name string
}

func NewDateTime() i3bar.BlockGenerator {
	return &DateTime{
		name: "datetime",
	}
}

func (g *DateTime) Frequency() uint8 {
	return 1
}

func (g *DateTime) Block(*i3bar.ColorSet) (*i3bar.Block, error) {
	cTime := time.Now().Local()

	return &i3bar.Block{
		Name:      g.name,
		FullText:  cTime.Weekday().String()[:2] + cTime.Format(" 2006-01-02 15:04:05"),
		ShortText: cTime.Format("15:04:05"),
	}, nil
}

func (g *DateTime) GetNameAndInstance() (string, string) {
	return g.name, ""
}