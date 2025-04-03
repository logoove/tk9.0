// Copyright 2025 The tk9.0-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tk9_0 // import "modernc.org/tk9.0"

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

var (
	_ Extension        = (*extension)(nil)
	_ ExtensionContext = extensionContext{}

	// Extensions register Tk extensions or extra packages. User code must
	// not directly mutate Extensions.
	Extensions = map[ExtensionKey]Extension{}

	AlreadyInitialized = errors.New("Already initialized")
	NotInitialized     = errors.New("Not initialized")
)

// RegisterExtension registers e.
func RegisterExtension(name string, e Extension) (r ExtensionKey, err error) {
	k := ExtensionKey{Type: typeName(e), Name: name}
	if _, ok := Extensions[k]; ok {
		return r, AlreadyRegistered
	}

	Extensions[k] = &extension{inner: e}
	return k, nil
}

type extension struct {
	inner Extension

	initialized bool
}

// ExtensionKey indexes Extensions
type ExtensionKey struct {
	Type string
	Name string
}

// ExtensionContext provides context to Extension methods.
type ExtensionContext interface {
	// Eval evaluates the tcl script.
	Eval(tcl string) (r string, err error)
	// EvalErr is like Eval, but handles errors according to the current value of
	// [ErrorMode].
	EvalErr(tcl string) (r string)
	RegisterWindow(path string) *Window
	Collect(w *Window, options ...any) string
	// Returns a single Tcl string, no braces, except "{}" is returned for s == "".
	TclSafeString(string) string
}

type extensionContext struct{}

func newExtensionContext() (r extensionContext) {
	return r
}

func (extensionContext) TclSafeString(s string) (r string) {
	return tclSafeString(s)
}

func (extensionContext) Eval(tcl string) (r string, err error) {
	return eval(tcl)
}

func (extensionContext) EvalErr(tcl string) (r string) {
	return evalErr(tcl)
}

func (extensionContext) RegisterWindow(path string) (w *Window) {
	w = &Window{path}
	windowIndex[path] = w
	return w
}

func (extensionContext) Collect(w *Window, options ...any) string {
	var a []string
	for _, v := range options {
		switch x := v.(type) {
		case Opt:
			a = append(a, x.optionString(w))
		default:
			a = append(a, tclSafeString(fmt.Sprint(x)))
		}
	}
	return strings.Join(a, " ")
}

// Extension handles Tk extensions. When calling Extension methods registered
// in Extension, the context argument is ignored and an instance is created
// automatically.
type Extension interface {
	// Initialize is called to perform any one-time initialization of an extension.
	// The Initialize method of an extension in Extensions can be called multiple
	// times but only the first successful call to Initialize will have any effect.
	Initialize(context ExtensionContext) error
}

func (e *extension) Initialize(context ExtensionContext) (err error) {
	context = newExtensionContext()

	defer func() {
		if err == nil {
			e.initialized = true
		}
	}()

	if e.initialized {
		return AlreadyInitialized
	}

	return e.inner.Initialize(context)
}

// InitializeExtension searches [Extensions] to find first extension named 'name'
// and call its Initialize method. The search is case sensitive, white space is
// normalized. If there's no match, InitializeExtension returns [NotFound].
//
// Any package can register extensions but only the main package can initialize
// an extension.
func InitializeExtension(name string) (err error) {
	if !isCalledFromMain() {
		return NotInitialized
	}

	var keys []ExtensionKey
	for k := range Extensions {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(a, b int) bool {
		c, d := keys[a], keys[b]
		e, f := matchName(c.Name), matchName(d.Name)
		if e < f {
			return true
		}
		if e > f {
			return false
		}

		return c.Type < d.Type
	})
	name = matchName(name)
	for _, k := range keys {
		if matchName(k.Name) == name {
			return Extensions[k].Initialize(nil)
		}
	}
	return NotFound
}

func extensionInitialized(name string) bool {
	for extName, ext := range Extensions {
		if extName.Name == name && ext.(*extension).initialized {
			return true
		}
	}
	return false
}
