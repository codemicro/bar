package i3bar

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"
)

type ColorSet struct {
	Bad     *Color
	Warning *Color
	Good    *Color
}

type Color struct {
	R, G, B uint8
}

func NewColorFromHexString(hexString string) (*Color, error) {
	hexString = strings.TrimPrefix(hexString, "#")

	if !(len(hexString) == 3 || len(hexString) == 6) {
		return nil, errors.New("invalid color length")
	}

	if len(hexString) == 3 {
		var newHexString string
		for _, char := range hexString {
			newHexString += string(char) + string(char)
		}
		hexString = newHexString
	}

	colorBytes, err := hex.DecodeString(hexString)
	if err != nil {
		return nil, err
	}

	return &Color{
		R: colorBytes[0], G: colorBytes[1], B: colorBytes[2],
	}, nil
}

func (c *Color) String() string {
	return "#" + hex.EncodeToString([]byte{c.R, c.G, c.B})
}

func (c *Color) MarshalJSON() ([]byte, error) {
	return []byte(`"` + c.String() + `"`), nil
}

func (c *Color) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, []byte("null")) {
		return nil
	}

	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	nc, err := NewColorFromHexString(s)
	if err != nil {
		return err
	}

	*c = *nc
	return nil
}
