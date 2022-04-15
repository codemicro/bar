package providers

import "github.com/codemicro/bar/internal/i3bar"

type PlainText struct {
	Text string
}

func NewPlainText(text string) i3bar.BlockGenerator {
	return &PlainText{
		Text: text,
	}
}

func (g *PlainText) Block(*i3bar.ColorSet) (*i3bar.Block, error) {
	return &i3bar.Block{
		Name: "plaintext",
		FullText: g.Text,
	}, nil
}