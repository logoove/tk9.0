// Copyright 2024 The tk9.0-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tk9_0 // import "modernc.org/tk9.0"

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"
)

var (
	// Themes register Tk themes. User code must not directly mutate Themes.
	Themes = map[ThemeKey]Theme{}

	AlreadyActivated  = errors.New("Already activated")
	AlreadyRegistered = errors.New("Already registered")
	Finalized         = errors.New("Finalized")
	NotActivated      = errors.New("Not activated")
	NotFound          = errors.New("Not found")

	currentTheme    Theme
	currentThemeKey ThemeKey

	_ Theme        = (*theme)(nil)
	_ Theme        = (*builtinTheme)(nil)
	_ ThemeContext = themeContext{}
)

// https://tkdocs.com/tutorial/styles.html
//
// Besides the built-in themes (alt, default, clam, and classic), macOS
// includes a theme named aqua to match the system-wide style, while Windows
// includes themes named vista, winxpnative, and winnative.
func init() {
	RegisterTheme("alt", &builtinTheme{"alt"})
	RegisterTheme("default", &builtinTheme{"default"})
	RegisterTheme("clam", &builtinTheme{"clam"})
	RegisterTheme("classic", &builtinTheme{"classic"})
	switch goos {
	case "darwin":
		RegisterTheme("aqua", &builtinTheme{"aqua"})
	case "windows":
		RegisterTheme("vista", &builtinTheme{"vista"})
		RegisterTheme("winxpnative", &builtinTheme{"winxpnative"})
		RegisterTheme("winnative", &builtinTheme{"winnative"})
	}
}

type builtinTheme struct {
	name string
}

func (t *builtinTheme) Activate(context ThemeContext) error {
	StyleThemeUse(t.name)
	return nil
}

func (t *builtinTheme) Deactivate(context ThemeContext) error {
	StyleThemeUse("default")
	return nil
}

func (t *builtinTheme) Finalize(context ThemeContext) error {
	return nil
}

func (t *builtinTheme) Initialize(context ThemeContext) error {
	return nil
}

// CurrentTheme returns the currently activated theme, if any.
func CurrentTheme() Theme {
	return currentTheme
}

// CurrentThemeName returns the name of the currently activated theme, if any.
func CurrentThemeName() (r string) {
	return currentThemeKey.Name
}

// ActivateTheme searches [Themes] to find first theme named like 'name' and
// call its Activate method. The search is case insensitive, using
// strings.ToLower, and white space is normalized. If there's no match,
// ActivateTheme returns [NotFound].
//
// Any package can register themes but only the main package can activate a
// theme.
func ActivateTheme(name string) (err error) {
	if !isCalledFromMain() {
		return NotActivated
	}

	var keys []ThemeKey
	for k := range Themes {
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
			return Themes[k].Activate(nil)
		}
	}
	return NotFound
}

func matchName(s string) string {
	return strings.Join(strings.Fields(strings.ToLower(s)), " ")
}

// RegisterTheme registers t.
func RegisterTheme(name string, t Theme) (r ThemeKey, err error) {
	k := ThemeKey{Type: typeName(t), Name: name}
	if _, ok := Themes[k]; ok {
		return r, AlreadyRegistered
	}

	Themes[k] = &theme{inner: t, k: k}
	return k, nil
}

// ThemeKey indexes Themes
type ThemeKey struct {
	Type string
	Name string
}

// ThemeContext provides context to Theme methods.
type ThemeContext interface {
	Eval(tcl string) (r string, err error)
}

type themeContext struct{}

func newThemeContext() (r themeContext) {
	evalFunc = eval
	return r
}

// Eval evaluates 'tcl' and returns a result value and an error, if any.
func (themeContext) Eval(tcl string) (r string, err error) {
	return evalFunc(tcl)
}

var evalFunc func(string) (string, error)

// Theme provides handling of a Tk theme. When calling Theme methods registered
// in Theme, the context argument is ignored and an instance is created
// automatically.
type Theme interface {
	// Activate makes the theme active/in use. The Activate method of themes in
	// Themes automatically call Initialize if it was not called before.
	Activate(context ThemeContext) error
	// Deactivate makes the theme not active. Deactivate cannot be called before
	// Activate().
	Deactivate(context ThemeContext) error
	// Finalize is called to perform any cleanup. After Finalize returns, the theme
	// cannot be used. The Finalize method of themes in Themes automatically remove
	// the theme from Themes after Finalize completes.
	Finalize(context ThemeContext) error
	// Initialize is called to perform any one-time initialization of a theme. The
	// Initialize method of themes in Themes can be called multiple times but only
	// the first successful call to Initialize will have any effect.
	Initialize(context ThemeContext) error
}

type theme struct {
	inner Theme
	k     ThemeKey

	activated   bool
	finalized   bool
	initialized bool
}

func (t *theme) Activate(context ThemeContext) (err error) {
	if currentTheme != nil {
		currentTheme.Deactivate(nil)
		currentTheme = nil
		currentThemeKey = ThemeKey{}
	}
	context = newThemeContext()

	defer func() {
		evalFunc = nil
		if err == nil {
			t.activated = true
			currentTheme = t
			currentThemeKey = t.k
		}
	}()

	if t.finalized {
		return Finalized
	}

	if t.activated {
		return AlreadyActivated
	}

	if !t.initialized {
		if err = t.inner.Initialize(context); err != nil {
			return err
		}
	}

	evalFunc = eval
	return t.inner.Activate(context)
}

func (t *theme) Deactivate(context ThemeContext) (err error) {
	context = newThemeContext()

	defer func() {
		evalFunc = nil
		t.activated = false
		currentTheme = nil
		currentThemeKey = ThemeKey{}
	}()

	if t.finalized {
		return Finalized
	}

	if !t.activated {
		return NotActivated
	}

	return t.inner.Deactivate(context)
}

func (t *theme) Finalize(context ThemeContext) (err error) {
	context = newThemeContext()

	defer func() {
		evalFunc = nil
		t.finalized = true
		delete(Themes, t.k)
		currentTheme = nil
		currentThemeKey = ThemeKey{}
	}()

	if !t.finalized {
		err = t.inner.Finalize(context)
	}
	return err
}

func (t *theme) Initialize(context ThemeContext) (err error) {
	context = newThemeContext()

	defer func() {
		evalFunc = nil
		if err == nil {
			t.initialized = true
		}
	}()

	if t.finalized {
		return Finalized
	}

	if !t.initialized {
		err = t.inner.Initialize(context)
	}
	return err
}

func typeName(th any) string {
	t := reflect.TypeOf(th)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	return fmt.Sprintf("%s.%s", t.PkgPath(), t.Name())
}
