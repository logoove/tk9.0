// Copyright 2025 The tk9.0-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package tablelist provides the [tklib tablelist].
//
// To make the extension available in an application:
//
//	import "modernc.org/tk9.0/extension/tablelist"
//
// See the [modernc.org/tk9.0.Extension] documentation for information about
// initalizing extensions at runtime.
//
// [modernc.org/tk9.0.Extension]: https://pkg.go.dev/modernc.org/tk9.0#Extension
// [tklib tablelist]: https://github.com/tcltk/tklib/tree/master/modules/tablelist
package tablelist // import "modernc.org/tk9.0/extensions/tablelist"

import (
	_ "embed"
	"fmt"
	"strings"

	. "modernc.org/tk9.0"
)

var (
	_ Extension = (*extension)(nil)

	//go:embed embed/tablelist.zip
	zip string

	ctx         ExtensionContext
	initialized bool

	// Version reports the version of the tablelist package. Valid after successful
	// initialization.
	Version string
)

func init() {
	RegisterExtension("tablelist", newExtension())
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
	mount := "/extensions/tablelist"
	Version, err = ctx.Eval(fmt.Sprintf(`
lappend auto_path [zipfs mountdata %s %s]
package require tablelist_tile
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

// TablelistWidget represents a tklib tablelist.
//
// The extensive documentation is available [here].
//
// [here]: https://www.nemethi.de/tablelist/index.html
type TablelistWidget struct {
	*Window
}

// Tablelist0 returns a newly created tablelist using raw Tcl string 'args'.
// Example:
//
//	Tablelist0(`.t -columns {0 "First Column" 0 "Another column"}`)
//
// This is a bootstrap function enabling to use the megawidget without first
// implementing the overwhelming amount of its options and methods in Go.
//
// Use with caution.
func Tablelist0(args string) (r *TablelistWidget) {
	return &TablelistWidget{Window: ctx.RegisterWindow(ctx.EvalErr(fmt.Sprintf("tablelist::tablelist %s", args)))}
}

// Do executes the raw Tcl string 'args'. Example:
//
//	t.Do(`insert end [list "first row" "another value"]`)
//
// This is a bootstrap function enabling to use the megawidget without first
// implementing the overwhelming amount of its options and methods in Go.
//
// Use with caution.
func (t *TablelistWidget) Do(args string) string {
	return ctx.EvalErr(fmt.Sprintf("%s %s", t, args))
}
