// Copyright 2025 The tk9.0-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package autoscroll provides the [tklib autoscroll].
//
// This package allows scrollbars to be mapped and unmapped as needed depending
// on the size and content of the scrollbars scrolled widget. The scrollbar
// must be managed by either pack or grid, other geometry managers are not
// supported.
//
// When managed by pack, any geometry changes made in the scrollbars parent
// between the time a scrollbar is unmapped, and when it is mapped will be
// lost. It is an error to destroy any of the scrollbars siblings while the
// scrollbar is unmapped. When managed by grid, if anything becomes gridded in
// the same row and column the scrollbar occupied it will be replaced by the
// scrollbar when remapped.
//
// This package may be used on any scrollbar-like widget as long as it supports
// the set subcommand in the same style as scrollbar. If the set subcommand is
// not used then this package will have no effect.
//
// To make the extension available in an application:
//
//	import "modernc.org/tk9.0/extension/autoscroll"
//
// See the [modernc.org/tk9.0.Extension] documentation for information about
// initalizing extensions at runtime.
//
// [modernc.org/tk9.0.Extension]: https://pkg.go.dev/modernc.org/tk9.0#Extension
// [tklib autoscroll]: https://github.com/tcltk/tklib/tree/master/modules/autoscroll
package autoscroll // import "modernc.org/tk9.0/extensions/autoscroll"

import (
	_ "embed"
	"fmt"
	"strings"

	. "modernc.org/tk9.0"
)

var (
	_ Extension = (*extension)(nil)

	//go:embed embed/autoscroll.zip
	zip string

	ctx         ExtensionContext
	initialized bool

	// Version reports the version of the autscroll package.
	Version = "0.7"
)

func init() {
	RegisterExtension("autoscroll", newExtension())
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
	mount := "/extensions/autoscroll"
	_, err = ctx.Eval(fmt.Sprintf(`
lappend auto_path [zipfs mountdata %s %s]
package require autoscroll
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

// Autoscroll arranges for the already existing scrollbar 'scrollbar' to be
// mapped and unmapped as needed. The function returns its 'scrollbar'
// argument.
func Autoscroll(scrollbar *Window) (r *Window) {
	ctx.EvalErr(fmt.Sprintf("::autoscroll::autoscroll %s", scrollbar))
	return scrollbar
}

// Unautoscroll returns the named scrollbar to its original static state. The
// function returns its 'scrollbar' argument.
func Unautoscroll(scrollbar *Window) (r *Window) {
	ctx.EvalErr(fmt.Sprintf("::autoscroll::unautoscroll %s", scrollbar))
	return scrollbar
}
