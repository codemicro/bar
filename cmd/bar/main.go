package main

import (
	"os"

	"github.com/codemicro/bar/internal/i3bar"
)

type basicGenerator struct {
	Text string
}

func (b *basicGenerator) Block() (*i3bar.Block, error) {
	return &i3bar.Block{
		FullText: b.Text,
	}, nil
}

func main() {
	b := i3bar.New(os.Stdout)
	_ = b.Initialise()
	_ = b.Emit([]i3bar.BlockGenerator{&basicGenerator{Text: "hello world"}})
	_ = b.Emit([]i3bar.BlockGenerator{&basicGenerator{Text: "something something words"}})
}

// TODO: Accept signals to refresh