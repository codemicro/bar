package providers

import (
	"fmt"
	"strings"

	"github.com/codemicro/bar/internal/i3bar"
	"github.com/samber/lo"
)

type IPAddress struct {
	Adapter string

	name string
}

func NewIPAddress(adapter string) i3bar.BlockGenerator {
	return &IPAddress{
		Adapter: adapter,
		name:    "ipAddr",
	}
}

func (g *IPAddress) Frequency() uint8 {
	return 5
}

func (g *IPAddress) getAdapterIPAddress() (string, error) {
	// call ifconfig
	output, err := runCommand("ifconfig")
	if err != nil {
		return "", err
	}

	adapters := lo.Filter(
		strings.Split(string(output), "\n\n"),
		func(x string, _ int) bool {
			return x != ""
		},
	)

	var ipAddr string

	// parse output
	// split by \n\n
	for _, adapter := range adapters {
		fields := strings.Fields(adapter)

		if !strings.EqualFold(
			strings.TrimSuffix(fields[0], ":"), g.Adapter,
		) {
			continue
		}

		for i, field := range fields {
			if field == "inet" {
				ipAddr = fields[i+1]
				break
			}
		}
	}

	return ipAddr, nil
}

func (g *IPAddress) Block(colors *i3bar.ColorSet) (*i3bar.Block, error) {
	ipAddr, err := g.getAdapterIPAddress()
	if err != nil {
		return nil, err
	}

	block := &i3bar.Block{
		Name: g.name,
		Instance: g.Adapter,
	}

	if ipAddr == "" {
		block.TextColor = colors.Bad
		block.FullText = fmt.Sprintf("%s no IP", g.Adapter)
		block.ShortText = "no IP"
	} else {
		block.TextColor = colors.Good
		block.FullText = ipAddr
	}

	return block, nil
}

func (g *IPAddress) GetNameAndInstance() (string, string) {
	return g.name, g.Adapter
}