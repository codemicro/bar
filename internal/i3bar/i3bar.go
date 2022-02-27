package i3bar

import (
	"encoding/json"
	"io"
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

func (b *I3bar) Emit(generators []BlockGenerator) error {
	var blocks []*Block
	for _, generator := range generators {
		b, err := generator.Block(&ColorSet{
			Bad:     &Color{255, 0, 0},
			Warning: &Color{0, 0, 255},
		})
		if err != nil {
			return err
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