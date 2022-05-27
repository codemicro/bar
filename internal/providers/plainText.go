package providers

import "github.com/codemicro/bar/internal/i3bar"

type PlainText struct {
	Text string

	name string
}

func NewPlainText(text string) i3bar.BlockGenerator {
	return &PlainText{
		Text: text,
		name: "plaintext",
	}
}

func (g *PlainText) Block(*i3bar.ColorSet) (*i3bar.Block, error) {
	return &i3bar.Block{
		Name:     g.name,
		FullText: g.Text,
	}, nil
}

func (g *PlainText) GetNameAndInstance() (string, string) {
	return g.name, ""
}