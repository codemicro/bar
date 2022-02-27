package providers

import (
	"time"

	"github.com/codemicro/bar/internal/i3bar"
)

type DateTime struct {
	// TODO: 12 hour mode?
}

func NewDateTime() i3bar.BlockGenerator {
	return new(DateTime)
}

func (g *DateTime) Block(*i3bar.ColorSet) (*i3bar.Block, error) {
	cTime := time.Now().Local()

	return &i3bar.Block{
		Name: "datetime",
		FullText: cTime.Format("2006-01-02 15:04:05"),
		ShortText: cTime.Format("15:04:05"),
	}, nil
}
