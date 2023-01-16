package providers

import (
	"github.com/codemicro/bar/internal/i3bar"
	"os"
	"github.com/rs/zerolog/log"
)

type LaunchProgram struct {
	Text       string
	Executable string

	name string
}

func NewLaunchProgram(text string, executable string) i3bar.BlockGenerator {
	return &LaunchProgram{
		Text: text,
		Executable: executable,
		name: "launchProgram",
	}
}

func (g *LaunchProgram) Frequency() uint8 {
	return 0
}

func (g *LaunchProgram) Block(*i3bar.ColorSet) (*i3bar.Block, error) {
	return &i3bar.Block{
		Name:     g.name,
		FullText: g.Text,
	}, nil
}

func (g *LaunchProgram) GetNameAndInstance() (string, string) {
	return g.name, ""
}

func (g *LaunchProgram) OnClick(event *i3bar.ClickEvent) bool {
	if event.Button != i3bar.LeftMouseButton {
		return false
	}

	process, err := os.StartProcess(g.Executable, []string{g.Executable}, &os.ProcAttr{Files: []*os.File{os.Stdin, os.Stdout, os.Stderr}})
	if err != nil {
		log.Error().Err(err).Str("location", "launchProgram_onClick").Msg("Could not start process")
		return false
	}
	_ = process.Release()
	return false
}
