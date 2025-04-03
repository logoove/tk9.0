// Copyright 2025 The tk9.0-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The original Tcl/Tk code in embed/ctext.zip license is
/*

Copyright (c) 2024 The tk9.0-go Authors. All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this
list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice,
this list of conditions and the following disclaimer in the documentation
and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its contributors
may be used to endorse or promote products derived from this software without
specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

*/

// Package ctext provides the [tklib ctext] package.
//
// To make the extension available in an application:
//
//	import "modernc.org/tk9.0/extension/ctext"
//
// See the [modernc.org/tk9.0.Extension] documentation for information about
// initalizing extensions at runtime.
//
// [modernc.org/tk9.0.Extension]: https://pkg.go.dev/modernc.org/tk9.0#Extension
// [tklib ctext]: https://github.com/tcltk/tklib/tree/master/modules/ctext
package ctext // import "modernc.org/tk9.0/extensions/ctext"

import (
	_ "embed"
	"fmt"
	"strings"

	. "modernc.org/tk9.0"
)

var (
	_ Extension = (*extension)(nil)

	//go:embed embed/ctext.zip
	zip string

	ctx         ExtensionContext
	id0         int
	initialized bool
)

func init() {
	RegisterExtension("ctext", newExtension())
}

func tclBinaryString(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		fmt.Fprintf(&b, "\\x%02x", s[i])
	}
	s = b.String()
	return s
}

func setup(context ExtensionContext) (err error) {
	defer func() {
		initialized = true
	}()

	if initialized {
		return nil
	}

	ctx = context // initialize the "global" context
	mount := "/extensions/ctext"
	_, err = ctx.Eval(fmt.Sprintf(`
lappend auto_path [zipfs mountdata %s %s]
package require ctext
`, tclBinaryString(zip), mount))
	return err
}

type extension struct{}

func newExtension() *extension {
	return &extension{}
}

func (e *extension) Initialize(context ExtensionContext) error {
	return setup(context)
}

// Ctextidget represents a tklib ctext. Ctext extends
// [modernc.org/tk9.0.TextWidget] and has all its methods, like
// [modernc.org/tk9.0.TextWidget.Insert] etc.
//
// The documentation is available [here].
//
// [here]: https://core.tcl-lang.org/tklib/doc/trunk/embedded/md/tklib/files/modules/ctext/ctext.md
type CtextWidget struct {
	*Window
	*TextWidget
}

func id() string {
	id0++
	return fmt.Sprintf(".ctext%v", id0)
}

// Ctext returns a newly created CtextWidget. Pass the parent *Window as the
// first argument to make the widget parented.
func Ctext(options ...any) (r *CtextWidget) {
	path := id()
	if len(options) != 0 {
		if x, ok := options[0].(*Window); ok {
			path = x.String() + path
			options = options[1:]
		}
	}
	w := ctx.RegisterWindow(path)
	ctx.EvalErr(fmt.Sprintf("ctext %s %s", path, ctx.Collect(w, options...)))
	return &CtextWidget{Window: w, TextWidget: &TextWidget{Window: w}}
}

// LinemapSelectFg changes the selected line foreground. The default is black.
func (w *CtextWidget) LinemapSelectFg(color string) {
	ctx.EvalErr(fmt.Sprintf("%s configure -linemap_select_fg %s", w, ctx.TclSafeString(color)))
}

// LinemapSelectBg changes the selected line background. The default is yellow.
func (w *CtextWidget) LinemapSelectBg(color string) {
	ctx.EvalErr(fmt.Sprintf("%s configure -linemap_select_bg %s", w, ctx.TclSafeString(color)))
}

// Linemapfg changes the foreground of the linemap. The default is the same
// color as the main text widget.
func (w *CtextWidget) Linemapfg(color string) {
	ctx.EvalErr(fmt.Sprintf("%s configure -linemapfg %s", w, ctx.TclSafeString(color)))
}

// Linemapbg changes the background of the linemap. The default is the same
// color as the main text widget.
func (w *CtextWidget) Linemapbg(color string) {
	ctx.EvalErr(fmt.Sprintf("%s configure -linemapbg %s", w, ctx.TclSafeString(color)))
}

// AddHighlightClass adds a highlighting class class to the ctext widget
// pathName. The highlighting will be done with the color color. All words in
// the keywordlist will be highlighted.
func (w *CtextWidget) AddHighlightClass(class, color string, keywordList ...string) {
	for i, v := range keywordList {
		keywordList[i] = ctx.TclSafeString(v)
	}
	ctx.EvalErr(fmt.Sprintf("::ctext::addHighlightClass %s %s %s {%s}", w, ctx.TclSafeString(class), ctx.TclSafeString(color), strings.Join(keywordList, " ")))
}

// AddHighlightClassForRegexp adds a highlighting class class to the ctext
// widget pathName. The highlighting will be done with the color color. All
// text parts matching the regexp pattern will be highlighted.
//
// 'pattern' is a Tcl regexp and is passed unchanged.
func (w *CtextWidget) AddHighlightClassForRegexp(class, color, pattern string) {
	ctx.EvalErr(fmt.Sprintf("::ctext::addHighlightClassForRegexp %s %s %s %s", w, ctx.TclSafeString(class), ctx.TclSafeString(color), pattern))
}

// AddHighlightClassForSpecialChars adds a highlighting class class to the
// ctext widget pathName. The highlighting will be done with the color color.
// All chars in charstring will be highlighted.
func (w *CtextWidget) AddHighlightClassForSpecialChars(class, color, chars string) {
	ctx.EvalErr(fmt.Sprintf("::ctext::addHighlightClassForSpecialChars %s %s %s %s", w, ctx.TclSafeString(class), ctx.TclSafeString(color), ctx.TclSafeString(chars)))
}

// Fastinsert inserts text without updating the highlighting. Arguments are
// identical to the [modernc.org/TextWidget.Insert] command.
func (w *CtextWidget) Fastinsert(index any, chars string, options ...any) any {
	idx := fmt.Sprint(index)
	ctx.EvalErr(fmt.Sprintf("%s fastinsert %s %s %s", w, ctx.TclSafeString(idx), ctx.TclSafeString(chars), ctx.Collect(w.Window, options...)))
	return index
}

// EnableComments enablec C-style comment highlighting. The class for c-style
// comments is _cComment. The C comment highlighting is disabled by default.
func (w *CtextWidget) EnableComments() {
	ctx.EvalErr(fmt.Sprintf("::ctext::enableComments %s", w))
}

// DisableComments disables C comment highlighting.
func (w *CtextWidget) DisableComments() {
	ctx.EvalErr(fmt.Sprintf("::ctext::disableComments %s", w))
}
