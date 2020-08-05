package main

import (
	"fmt"
	"image/color"
)

type remapper struct {
	remaps map[rgb]rgb
}

func (r *remapper) load(tpl, values colorSet) {
	if r.remaps == nil {
		r.remaps = make(map[rgb]rgb)
	}
	r.remaps[tpl.XLight] = values.XLight
	r.remaps[tpl.Light] = values.Light
	r.remaps[tpl.Medium] = values.Medium
	r.remaps[tpl.Dark] = values.Dark
	r.remaps[tpl.XDark] = values.XDark
}

func (r *remapper) remap(in color.Color) (shouldRemap bool, to color.Color) {
	switch v := in.(type) {
	case color.NRGBA:
		if c, ok := r.remaps[rgb{v.R, v.G, v.B}]; ok {
			return true, color.RGBA{c.R, c.G, c.B, 255}
		}
	case color.RGBA:
		if c, ok := r.remaps[rgb{v.R, v.G, v.B}]; ok {
			return true, color.RGBA{c.R, c.G, c.B, 255}
		}
	default:
		panic(fmt.Sprintf("unhandled pixel of color type %T", v))
	}

	return false, in
}
