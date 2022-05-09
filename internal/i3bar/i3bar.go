package i3bar

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/rs/zerolog/log"
)

type I3bar struct {
	writer io.Writer

	hasSentFirstLine bool
}

func New(writer io.Writer) *I3bar {
	return &I3bar{
		writer: writer,
	}
}

func (b *I3bar) Initialise() error {
	_, err := b.writer.Write([]byte(
		[]byte("{\"version\":1}\n"), // This means that versions of i3 prior to
		// 4.3 can still work with this bar. We do not use touch features, nor
		//do we use any special stop/start handling. That's handled by the OS.
	))
	return err
}

var defaultColorSet = &ColorSet{
	Good:       &Color{0xb8, 0xbb, 0x26},
	Bad:        &Color{251, 73, 52},
	Warning:    &Color{250, 189, 47},
	Background: &Color{0x28, 0x28, 0x28},
}

func (b *I3bar) Emit(generators []BlockGenerator) error {
	var blocks []*Block
	for _, generator := range generators {
		b, err := generator.Block(defaultColorSet)
		if err != nil {
			log.Error().Err(err).Str("generator", fmt.Sprintf("%T", generator)).Send()
			b = &Block{
				FullText:  "ERROR",
				TextColor: defaultColorSet.Bad,
			}
		}
		blocks = append(blocks, b)
	}

	jsonData, err := json.Marshal(blocks)
	if err != nil {
		return err
	}

	if !b.hasSentFirstLine {
		jsonData = append([]byte("[\n"), jsonData...)
		b.hasSentFirstLine = true
	} else {
		jsonData = append([]byte{','}, jsonData...)
	}

	jsonData = append(jsonData, '\n')

	if _, err := b.writer.Write(jsonData); err != nil {
		return err
	}

	return nil
}

type Block struct {
	FullText            string `json:"full_text"`
	ShortText           string `json:"short_text,omitempty"`
	TextColor           *Color `json:"color,omitempty"`
	BackgroundColor     *Color `json:"background,omitempty"`
	BorderColor         *Color `json:"border,omitempty"`
	BorderTop           int    `json:"border_top,omitempty"`
	BorderRight         int    `json:"border_right,omitempty"`
	BorderBottom        int    `json:"border_bottom,omitempty"`
	BorderLeft          int    `json:"border_left,omitempty"`
	MinWidth            string `json:"min_width,omitempty"`
	Align               string `json:"align,omitempty"`
	Urgent              bool   `json:"urgent,omitempty"`
	Name                string `json:"name,omitempty"`
	Instance            string `json:"instance,omitempty"`
	Separator           bool   `json:"separator,omitempty"`
	SeparatorBlockWidth int    `json:"separator_block_width,omitempty"`
	Markup              string `json:"markup,omitempty"`
}

type BlockGenerator interface {
	Block(*ColorSet) (*Block, error)
}
