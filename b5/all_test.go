// Copyright 2024 The tk9.0-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package b5 // import "modernc.org/tk9.0/b5"

import (
	"testing"

	. "modernc.org/tk9.0"
)

func Test1(t *testing.T) {
	width := TkScaling() * 72 * buttonFocusDecoratorCorner
	stroke := TkScaling() * 72 * buttonFocusDecorator
	r := width - stroke/2
	clip := width - stroke
	trc("width=%v clip=%v r=%v stroke=%v", width, clip, r, stroke)
	const k = 10
	getCorners(k*round(width), k*round(clip), k*round(r), k*round(stroke), "#0b5ed7", "#fff", "#fff")
}

func Test2(t *testing.T) {
	width := TkScaling() * 72 * buttonFocusDecoratorCorner
	stroke := TkScaling() * 72 * buttonFocusDecorator
	r := width - stroke/2
	clip := width
	trc("width=%v clip=%v r=%v stroke=%v", width, clip, r, stroke)
	const k = 10
	getCorners(k*round(width), k*round(clip), k*round(r), k*round(stroke), "#0b5ed7", "#97c1fe", "#fff")
}
