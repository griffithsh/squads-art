package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
)

type rgb struct {
	R, G, B uint8
}

// MarshalJSON adds custom marshaling to support simplified data entry of
// colors.
func (c rgb) MarshalJSON() ([]byte, error) {
	dst := make([]byte, 6)
	hex.Encode(dst, []byte{c.R, c.G, c.B})

	dst = append([]byte{'"'}, dst...)
	dst = append(dst, '"')
	return dst, nil
}

// UnmarshalJSON adds custom marshaling to support simplified data entry of
// colors.
func (c *rgb) UnmarshalJSON(data []byte) error {
	var sraw string
	if err := json.Unmarshal(data, &sraw); err != nil {
		return fmt.Errorf("unmarshal raw: %v", err)
	}
	if len(sraw) != 6 {
		return fmt.Errorf("cannot Unmarshal \"%s\" to type rgb", sraw)
	}

	raw := []byte(sraw)
	dec := make([]byte, 3)
	if _, err := hex.Decode(dec, raw); err != nil {
		return fmt.Errorf("decode hex: %v", err)
	}

	c.R = dec[0]
	c.G = dec[1]
	c.B = dec[2]
	return nil
}

type variant struct {
	Name   string
	Colors colorSet
}
type colorSet struct {
	XLight rgb `json:"xlt"`
	Light  rgb `json:"lt"`
	Medium rgb `json:"md"`
	Dark   rgb `json:"dk"`
	XDark  rgb `json:"xdk"`
}

// config is the type that unmarshals from config.json in character-appearance/.
type config struct {
	HairTemplate, SkinTemplate colorSet
	HairColors                 []variant
	SkinColors                 []variant
}
