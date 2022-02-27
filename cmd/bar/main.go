package main

import (
	"os"

	"github.com/codemicro/bar/internal/i3bar"
	"github.com/codemicro/bar/internal/providers"
)

func main() {
	b := i3bar.New(os.Stdout)
	_ = b.Initialise()
	_ = b.Emit([]i3bar.BlockGenerator{&providers.Memory{}})
}

// TODO: Accept signals to refresh