package providers

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/codemicro/bar/internal/i3bar"
	"github.com/samber/lo"
)

type WiFi struct {
	Adapter     string
	OkThreshold float32

	name string
}

func NewWiFi(adapter string, okThreshold float32) i3bar.BlockGenerator {
	return &WiFi{
		Adapter:     adapter,
		OkThreshold: okThreshold,

		name: "wifi",
	}
}

var (
	// For use with iwconfig
	essidRegexp       = regexp.MustCompile(`ESSID:(?:"(.+)"|off/any)`)
	frequencyRegexp   = regexp.MustCompile(`Frequency:(\d(?:\.\d+)? [a-zA-Z]Hz)`)
	linkQualityRegexp = regexp.MustCompile(`Link Quality=(\d+\/\d+)`)
)

func (g *WiFi) getAdapterIPAddress() (string, error) {
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

func (g *WiFi) getConnectionInfo() (ssid, frequency string, linkQuality float32, err error) {
	output, err := runCommand("iwconfig")
	if err != nil {
		return "", "", 0, err
	}

	adapters := lo.Filter(strings.Split(string(output), "\n\n"), func(x string, _ int) bool {
		return x != ""
	})

	for _, adapterInfo := range adapters {
		if !strings.HasPrefix(adapterInfo, g.Adapter) {
			continue
		}

		if essidRegexp.MatchString(adapterInfo) {
			ssid = essidRegexp.FindStringSubmatch(adapterInfo)[1]
		}

		if frequencyRegexp.MatchString(adapterInfo) {
			frequency = frequencyRegexp.FindStringSubmatch(adapterInfo)[1]
		}

		if linkQualityRegexp.MatchString(adapterInfo) {
			eqn := linkQualityRegexp.FindStringSubmatch(adapterInfo)[1]
			sp := strings.Split(eqn, "/")
			num, _ := strconv.Atoi(sp[0])
			denom, _ := strconv.Atoi(sp[1])
			linkQuality = (float32(num) / float32(denom)) * 100
		}
	}

	return
}

func (g *WiFi) Block(colors *i3bar.ColorSet) (*i3bar.Block, error) {
	ssid, frequency, linkQuality, err := g.getConnectionInfo()
	if err != nil {
		return nil, err
	}

	block := &i3bar.Block{
		Name: g.name,
		Instance: g.Adapter,
	}

	if ssid == "" {
		block.TextColor = colors.Bad
		block.FullText = fmt.Sprintf("%s not connected", g.Adapter)
		block.ShortText = "not connected"
	} else {

		if linkQuality < g.OkThreshold && g.OkThreshold != 0 {
			block.TextColor = colors.Warning
		}

		block.TextColor = colors.Good
		block.FullText = fmt.Sprintf("%s (%s) %.0f%%", ssid, strings.ReplaceAll(frequency, " ", ""), linkQuality)
	}

	return block, nil
}

func (g *WiFi) GetNameAndInstance() (string, string) {
	return g.name, g.Adapter
}