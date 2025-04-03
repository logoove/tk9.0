// Copyright 2025 The tk9.0-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package eval provides raw Tcl eval.
//
// To make the extension available in an application:
//
//	import "modernc.org/tk9.0/extension/eval"
//
// See the [modernc.org/tk9.0.Extension] documentation for information about
// initalizing extensions at runtime.
//
// [modernc.org/tk9.0.Extension]: https://pkg.go.dev/modernc.org/tk9.0#Extension
package eval // import "modernc.org/tk9.0/extensions/eval"

import (
	. "modernc.org/tk9.0"
)

var (
	_ Extension = (*extension)(nil)

	ctx         ExtensionContext
	initialized bool
)

func init() {
	RegisterExtension("eval", newExtension())
}

type extension struct{}

func newExtension() *extension {
	return &extension{}
}

func (e *extension) Initialize(context ExtensionContext) error {
	defer func() {
		initialized = true
	}()

	if initialized {
		return nil
	}

	ctx = context // initialize the "global" context
	return nil
}

// Eval evaluates raw Tcl string 'tcl'. Use with caution.
func Eval(tcl string) (r string, err error) {
	return ctx.Eval(tcl)
}

// EvalErr evaluates raw Tcl string 'tcl'. Errors are handled accoring to
// [ErrorMode]. Use with caution.
//
// [ErrorMode]: https://pkg.go.dev/modernc.org/tk9.0#ErrorMode
func EvalErr(tcl string) (r string) {
	return ctx.EvalErr(tcl)
}
