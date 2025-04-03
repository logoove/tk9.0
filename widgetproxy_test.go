// Copyright 2025 The tk9.0-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tk9_0 // import "modernc.org/tk9.0"

import (
	"strings"
	"testing"
)

func TestWidgetProxy(t *testing.T) {
	// Widget proxying checks that the eval extension is initialized.
	// However it can't be initialized by calling InitializeExtension() here,
	// as that requires the extension to have been registered. But importing the
	// modernc.org/tk9.0/extensions/eval package would result in an import cycle.
	//
	// Instead let's be a little naughty, and simply pretend that it's initialized.
	Extensions[ExtensionKey{Name: "eval"}] = &extension{initialized: true}

	return
	text := NewTextWidgetProxy(Text())

	assertContent := func(expected string) {
		content := text.Get("1.0", "end -1 chars")[0]
		if expected != content {
			t.Errorf("expected '%s', but is '%s'", expected, content)
		}
	}

	insertCallback := func(args []string) {
		args[2] = strings.ToUpper(args[2]) // upper case the text to be inserted
		text.EvalWrapped(args)
	}

	// no callback registered
	text.Insert(END, "abc ")
	assertContent("abc ")

	// callback registered
	text.Register("insert", insertCallback)
	text.Insert(END, "def ")
	assertContent("abc DEF ")

	// callback unregistered
	text.Unregister("insert")
	text.Insert(END, "ghi ")
	assertContent("abc DEF ghi ")

	// callback registered again
	text.Register("insert", insertCallback)
	text.Insert(END, "jkl ")
	assertContent("abc DEF ghi JKL ")

	// proxying removed
	text.Close()
	text.Insert(END, "mno ")
	assertContent("abc DEF ghi JKL mno ")
}
