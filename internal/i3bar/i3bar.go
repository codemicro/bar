package i3bar

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
)

type I3bar struct {
	writer               io.Writer
	updateInterval       time.Duration
	updateSignal         syscall.Signal
	registeredGenerators []BlockGenerator

	hasInitialised   bool
	hasSentFirstLine bool
}

func New(writer io.Writer, updateInterval time.Duration, updateSignal syscall.Signal) *I3bar {
	return &I3bar{
		writer:         writer,
		updateInterval: updateInterval,
		updateSignal:   updateSignal,
	}
}

func (b *I3bar) Initialise() error {
	capabilities, err := json.Marshal(map[string]any{
		"version": 1,
		"click_events": true,
	})
	if err != nil {
		return err
	}
	
	if _, err := b.writer.Write(append(capabilities, '\n')); err != nil {
		return err
	}

	b.hasInitialised = true
	return nil
}

var defaultColorSet = &ColorSet{
	Good:       &Color{0xb8, 0xbb, 0x26},
	Bad:        &Color{251, 73, 52},
	Warning:    &Color{250, 189, 47},
	Background: &Color{0x28, 0x28, 0x28},
}

func (b *I3bar) Emit(generators []BlockGenerator) error {
	if !b.hasInitialised {
		if err := b.Initialise(); err != nil {
			return err
		}
	}

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

func (b *I3bar) RegisterBlockGenerator(bg ...BlockGenerator) {
	b.registeredGenerators = append(b.registeredGenerators, bg...)
}

func (b *I3bar) StartLoop() error {
	// The ticker will start after the specified duration, not when we
	// instantiate it. Circumventing that here by calling Emit now.
	if err := b.Emit(b.registeredGenerators); err != nil {
		return err
	}

	ticker := time.NewTicker(b.updateInterval)
	sigUpdate := make(chan os.Signal, 1)
	signal.Notify(sigUpdate, os.Signal(b.updateSignal))

	for {
		select {
		case <-sigUpdate:
			if err := b.Emit(b.registeredGenerators); err != nil {
				return err
			}
		case <-ticker.C:
			if err := b.Emit(b.registeredGenerators); err != nil {
				return err
			}
		}
	}
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
	GetNameAndInstance() (name, instance string)
}
