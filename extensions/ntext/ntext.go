// Copyright 2025 The tk9.0-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package ntext provides the [tklib ntext] package.
//
// To make the extension available in an application:
//
//	import "modernc.org/tk9.0/extension/ntext"
//
// See the [modernc.org/tk9.0.Extension] documentation for information about
// initalizing extensions at runtime.
//
// [modernc.org/tk9.0.Extension]: https://pkg.go.dev/modernc.org/tk9.0#Extension
// [tklib ntext]: https://github.com/tcltk/tklib/tree/master/modules/ntext
package ntext // import "modernc.org/tk9.0/extensions/ntext"

import (
	_ "embed"
	"fmt"
	"strings"

	. "modernc.org/tk9.0"
)

var (
	_ Extension = (*extension)(nil)

	//go:embed embed/ntext.zip
	zip string

	ctx         ExtensionContext
	initialized bool
)

func init() {
	RegisterExtension("ntext", newExtension())
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
	mount := "/extensions/ntext"
	_, err = ctx.Eval(fmt.Sprintf(`
lappend auto_path [zipfs mountdata %s %s]
package require ntext
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
