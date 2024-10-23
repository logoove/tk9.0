// Copyright 2024 The tk9.0-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The b5 package is a work in progress with no stable API yet. It will eventually become
// a full theme package.
package b5 // import "modernc.org/tk9.0/b5"

import (
	"fmt"
	"math"

	. "modernc.org/tk9.0"
)

const (
	buttonFocusDecoratorCorner = 9 / 96.  // The rounded corner is 9px on a 96 DPI display.
	buttonFocusDecorator       = 4 / 96.  // 4px on a 96 DPI display.
	buttonTileHeight           = 27 / 96. // 27px on a 96 DPI display
)

var (
	corners = map[cornerKey][4]*Img{}
	tiles   = map[tileKey]*Img{}
)

type Color int

const (
	_ Color = iota
	ButtonFace
	ButtonFocus
	ButtonText
)

type Colors map[Color]string

type tileKey struct {
	width  int
	height int
	color  string
}

type cornerKey struct {
	width       int
	clip        int
	r           int
	strokeWidth int
	fill        string
	stroke      string
	background  string
}

func round(n float64) int {
	return int(math.Round(n))
}

// All sizes in px
func getCorners(width, clip, r, strokeWidth int, fill, stroke, background string) (re [4]*Img) {
	k := cornerKey{width, clip, r, strokeWidth, fill, stroke, background}
	if ex, ok := corners[k]; ok {
		return ex
	}

	svg := fmt.Sprintf(`<svg>
	<rect width="%[7]d" height="%[7]d" fill=%[6]q />
	<circle r="%[2]d" cx="%[1]d" cy="%[1]d" stroke-width="%[3]d" fill=%q stroke=%q />
</svg>`,
		width, r, strokeWidth, fill, stroke, background, 2*width)
	img := NewPhoto(Data(svg))
	re[0] = NewPhoto(Width(clip), Height(clip))
	re[0].Copy(img, From(width, width-clip, width+clip, width))
	re[1] = NewPhoto(Width(clip), Height(clip))
	re[1].Copy(img, From(width-clip, width-clip, width, width))
	re[2] = NewPhoto(Width(clip), Height(clip))
	re[2].Copy(img, From(width-clip, width, width, width+clip))
	re[3] = NewPhoto(Width(clip), Height(clip))
	re[3].Copy(img, From(width, width, width+clip, width+clip))
	corners[k] = re
	return re
}

// All sizes in px
func getTile(width, height int, color string) (r *Img) {
	k := tileKey{width, height, color}
	if ex, ok := tiles[k]; ok {
		return ex
	}

	r = NewPhoto(Width(width), Height(height),
		Data(fmt.Sprintf(`<svg width="%d" height="%d" fill=%q><rect width="%[1]d" height="%d" fill=%q/></svg>`, width, height, color)))
	tiles[k] = r
	return r
}

// ButtonStyle defines a button style. ATM only when using the "default" theme.
//
// This function is intended for prototyping and will be most probably unexported at some time.
func ButtonStyle(style string, colors Colors, background string, focused bool) string {
	width := TkScaling() * 72 * buttonFocusDecoratorCorner
	stroke := TkScaling() * 72 * buttonFocusDecorator
	th := TkScaling() * 72 * buttonTileHeight
	r := width - stroke/2
	clip := width - stroke
	focus := background
	if focused {
		focus = colors[ButtonFocus]
	}
	ocorners := getCorners(round(width), round(width), round(r), round(stroke), colors[ButtonFace], focus, background)
	oq1 := style + ".p1"
	oq2 := style + ".p2"
	oq3 := style + ".p3"
	oq4 := style + ".p4"
	StyleElementCreate(oq1, "image", ocorners[0])
	StyleElementCreate(oq2, "image", ocorners[1])
	StyleElementCreate(oq3, "image", ocorners[2])
	StyleElementCreate(oq4, "image", ocorners[3])
	icorners := getCorners(round(width), round(clip), round(r), round(stroke), colors[ButtonFace], focus, background)
	iq1 := style + ".iq1"
	iq2 := style + ".iq2"
	iq3 := style + ".iq3"
	iq4 := style + ".iq4"
	StyleElementCreate(iq1, "image", icorners[0])
	StyleElementCreate(iq2, "image", icorners[1])
	StyleElementCreate(iq3, "image", icorners[2])
	StyleElementCreate(iq4, "image", icorners[3])
	tile := "Tile." + style + ".tile"
	t := getTile(8, round(th), colors[ButtonFace])
	StyleElementCreate(tile, "image", t)
	StyleLayout(style,
		"Button.border", Sticky("nswe"), Children(
			"Button.focus", Sticky("nswe"), Children(
				oq1, Sticky("ne"),
				oq2, Sticky("nw"),
				oq3, Sticky("sw"),
				oq4, Sticky("se"),
				"Button.padding", Sticky("nswe"), Children(
					tile,
					iq1, Sticky("ne"),
					iq2, Sticky("nw"),
					iq3, Sticky("sw"),
					iq4, Sticky("se"),
					"Button.label", Sticky("nswe")))))
	StyleConfigure(style, Background(focus), Borderwidth(0), Compound(true), Focuscolor(focus), Focussolid(false),
		Focusthickness(0), Foreground(colors[ButtonText]), Padding(round(stroke)), Relief("flat"), Shiftrelief(0))
	StyleMap(style, Background, "disabled", "#edeceb")
	return style
}
