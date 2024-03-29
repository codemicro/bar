package i3bar

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
)

type generatorInfo struct {
	Provider         BlockGenerator
	HasClickConsumer bool
	Last      *Block
}

type I3bar struct {
	writer         io.Writer
	reader         io.Reader
	updateSignal   syscall.Signal

	generators []*generatorInfo
	tickNumber uint8
}

func New(writer io.Writer, reader io.Reader, updateSignal syscall.Signal) *I3bar {
	return &I3bar{
		writer:         writer,
		reader:         reader,
		updateSignal:   updateSignal,
	}
}

func (b *I3bar) Initialise() error {
	capabilities, err := json.Marshal(map[string]any{
		"version":      1,
		"click_events": true,
	})
	if err != nil {
		return err
	}

	if _, err := b.writer.Write(append(capabilities, []byte("\n[\n")...)); err != nil {
		return err
	}

	return nil
}

var defaultColorSet = &ColorSet{
	Good:       &Color{0xb8, 0xbb, 0x26},
	Bad:        &Color{251, 73, 52},
	Warning:    &Color{250, 189, 47},
	Background: &Color{0x28, 0x28, 0x28},
}

func (b *I3bar) Emit(blocks []*Block) error {
	jsonData, err := json.Marshal(blocks)
	if err != nil {
		return err
	}

	jsonData = append(jsonData, []byte(",\n")...)

	if _, err := b.writer.Write(jsonData); err != nil {
		return err
	}

	return nil
}

// RegisterBlockGenerator registers a block generator with the status bar. This
// function should not be called after StartLoop is called.
func (b *I3bar) RegisterBlockGenerator(bg ...BlockGenerator) {
	for _, bgx := range bg {
		_, hasClickConsumer := bgx.(ClickEventConsumer)

		metadata := new(generatorInfo)
		metadata.Provider = bgx
		metadata.HasClickConsumer = hasClickConsumer

		b.generators = append([]*generatorInfo{metadata}, b.generators...)
	}
}

func (b *I3bar) StartLoop() error {

	// The ticker will start after the specified duration, not when we
	// instantiate it. Circumventing that here by calling Emit now.
	if err := b.tick(false); err != nil {
		return err
	}

	ticker := time.NewTicker(time.Second)
	sigUpdate := make(chan os.Signal, 1)
	signal.Notify(sigUpdate, os.Signal(b.updateSignal))

	go b.consumerLoop(func() {
		sigUpdate <- b.updateSignal
	})

	for {
		select {
		case <-sigUpdate:
			if err := b.tick(true); err != nil {
				log.Error().Err(err).Msg("could not tick")
			}
		case <-ticker.C:
			if err := b.tick(false); err != nil {
				log.Error().Err(err).Msg("could not tick")
			}
		}
	}
}

func (b *I3bar) tick(override bool) error {
	var hasChanged bool
	for _, gen := range b.generators {
		if override || (gen.Provider.Frequency() == 0 && gen.Last == nil) || (gen.Provider.Frequency() != 0 && b.tickNumber%gen.Provider.Frequency() == 0) {
			block, err := gen.Provider.Block(defaultColorSet)
			if err != nil {
				log.Error().Err(err).Str("generator", fmt.Sprintf("%T", gen.Provider)).Send()
				block = &Block{
					FullText:  "ERROR",
					TextColor: defaultColorSet.Bad,
				}
			}
			if block == nil {
				block = &Block{
					FullText:  "MISSING",
					TextColor: defaultColorSet.Warning,
				}
			}

			if block != gen.Last {
				gen.Last = block
				hasChanged = true
			}
		}
	}

	if hasChanged {
		var blocks []*Block
		for _, gen := range b.generators {
			blocks = append(blocks, gen.Last)
		}
		if err := b.Emit(blocks); err != nil {
			return err
		}
	}

	b.tickNumber += 1

	return nil
}

func (b *I3bar) consumerLoop(requestBarRefresh func()) {
	r := bufio.NewReader(b.reader)
	for {
		inputBytes, err := r.ReadBytes('\n')
		if err != nil {
			log.Error().Err(err).Msg("could not read from input reader")
			continue
		}

		// "ReadBytes reads until the first occurrence of delim in the input,
		// returning a slice containing the data up to and including the
		// delimiter."
		inputBytes = inputBytes[:len(inputBytes)-1]

		// try and parse inputted JSON
		event := new(ClickEvent)
		if err := json.Unmarshal(bytes.Trim(inputBytes, ","), event); err != nil {
			continue // idk what this could be but it's not relevant so BYE!
		}

		for _, consumer := range b.generators {
			if !consumer.HasClickConsumer {
				continue
			}
			consumerName, consumerInstance := consumer.Provider.GetNameAndInstance()
			if consumerName == event.Name && (consumerName == "" || consumerInstance == event.Instance) {
				if consumer.Provider.(ClickEventConsumer).OnClick(event) {
					requestBarRefresh()
				}
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

type ProvidesNameAndInstance interface {
	GetNameAndInstance() (name, instance string)
}

type BlockGenerator interface {
	ProvidesNameAndInstance
	Block(*ColorSet) (*Block, error)
	// Frequency should return the frequency at which the block is updated, as
	// "1 per n seconds".
	//
	// A frequency of zero means that the Block method will only be called
	// once, ever.
	Frequency() uint8
}

type ClickEvent struct {
	Name      string          `json:"name"`
	Instance  string          `json:"instance"`
	Button    MouseButtonType `json:"button"`
	Modifiers []string        `json:"modifiers"`
	X         int             `json:"x"`
	Y         int             `json:"y"`
	RelativeX int             `json:"relative_x"`
	RelativeY int             `json:"relative_y"`
	OutputX   int             `json:"output_x"`
	OutputY   int             `json:"output_y"`
	Width     int             `json:"width"`
	Height    int             `json:"height"`
}

type ClickEventConsumer interface {
	ProvidesNameAndInstance
	// OnClick is called when a new ClickEvent is recieved with the
	// corresponding name and instance is recieved. If OnClick returns true, a
	// refresh of the entire statusbar will be performed.
	//
	// OnClick must not modify the ClickEvent as it may be reused elsewhere.
	OnClick(*ClickEvent) (shouldRefresh bool)
}

type MouseButtonType uint8

const (
	LeftMouseButton MouseButtonType = iota + 1
	MiddleMouseButton
	RightMouseButton
	MouseWheelScrollUp
	MouseWheelScrollDown
)
