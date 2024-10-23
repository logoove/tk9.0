// Copyright 2024 The tk9.0-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tk9_0 // import "modernc.org/tk9.0"

import (
	_ "embed"
	"flag"
	"fmt"
	"os"
	"testing"

	_ "github.com/adrg/xdg"       // generator.go
	_ "github.com/expr-lang/expr" // examples
	_ "golang.org/x/net/html"     // generator.go
	_ "modernc.org/ngrab/lib"     // generator.go
	_ "modernc.org/rec/lib"       // generator.go
)

func TestMain(m *testing.M) {
	if isBuilder {
		os.Exit(0)
	}

	if Error != nil {
		fmt.Fprintln(os.Stderr, Error)
		os.Exit(1)
	}

	flag.Parse()
	rc := m.Run()
	Finalize()
	os.Exit(rc)
}

func TestTokenizer(t *testing.T) {
	for i, test := range []struct {
		s    string
		ids  []int
		toks []string
	}{
		{},
		{"a", []int{0}, []string{"a"}},
		{"\\$", []int{0}, []string{"\\$"}},
		{"\\$\\$", []int{0}, []string{"\\$\\$"}},
		{"\\$\\$\\$", []int{0}, []string{"\\$\\$\\$"}},

		{"\\$\\$\\$\\$", []int{0}, []string{"\\$\\$\\$\\$"}},
		{"a\\$", []int{0}, []string{"a\\$"}},
		{"a\\$\\$", []int{0}, []string{"a\\$\\$"}},
		{"a\\$\\$\\$", []int{0}, []string{"a\\$\\$\\$"}},
		{"a\\$\\$\\$\\$", []int{0}, []string{"a\\$\\$\\$\\$"}},

		{"$a$", []int{1}, []string{"$a$"}},
		{"$$a$", []int{2}, []string{"$$a$"}},
		{"$$a$$", []int{2}, []string{"$$a$$"}},
		{"$a$$", []int{2}, []string{"$a$$"}},
		{"x$a$", []int{0, 1}, []string{"x", "$a$"}},

		{"x$$a$", []int{0, 2}, []string{"x", "$$a$"}},
		{"x$$a$$", []int{0, 2}, []string{"x", "$$a$$"}},
		{"x$a$$", []int{0, 2}, []string{"x", "$a$$"}},
		{"x$a$y", []int{0, 1, 0}, []string{"x", "$a$", "y"}},
		{"x$$a$y", []int{0, 2, 0}, []string{"x", "$$a$", "y"}},

		{"x$$a$$y", []int{0, 2, 0}, []string{"x", "$$a$$", "y"}},
		{"x$a$$y", []int{0, 2, 0}, []string{"x", "$a$$", "y"}},
		{"x\\$0$a\\$1b$$\\$y", []int{0, 2, 0}, []string{"x\\$0", "$a\\$1b$$", "\\$y"}},
	} {
		ids, toks := tokenize(test.s)
		if g, e := fmt.Sprintf("%v %q", ids, toks), fmt.Sprintf("%v %q", test.ids, test.toks); g != e {
			t.Errorf("#%3v: `%s`\ngot %s\nexp %s", i, test.s, g, e)
		}
	}
}

// func capitalize(s string) string {
// 	return strings.ToUpper(string(s[0])) + s[1:]
// }
//
// func TestTmp(t *testing.T) {
// 	ErrorMode = CollectErrors
// 	themes := StyleThemeNames()
// 	slices.Sort(themes)
// 	for _, theme := range themes {
// 		fmt.Printf("\n//\n// # %q theme style guide", theme)
// 		StyleThemeUse(theme)
// 		styles := StyleThemeStyles()
// 		slices.Sort(styles)
// 		elements := StyleElementNames()
// 		slices.Sort(elements)
// 		for _, element := range elements {
// 			fmt.Printf("\n//\n// %q style element options:\n//", element)
// 			options := StyleElementOptions(element)
// 			slices.Sort(options)
// 			for _, option := range options {
// 				fmt.Printf("\n//  - [%s]", capitalize(option[1:]))
// 			}
// 		}
// 		fmt.Printf("\n//\n// %q theme style list", theme)
// 		needsep := false
// 		for _, style := range styles {
// 			if needsep {
// 				fmt.Printf("\n//\n// -\n//\n//")
// 			}
// 			needsep = true
// 			fmt.Printf("\n//\n//  %s", style)
// 			if s := StyleLayout(style); s != "" {
// 				needsep = false
// 				fmt.Printf("\n//\n// Layout: %s", s)
// 			}
// 			if s := strings.TrimSpace(StyleMap(style)); s != "" {
// 				needsep = false
// 				s = strings.Join(strings.Fields(s), " ")
// 				fmt.Printf("\n//\n// Style map: %s", s)
// 			}
// 		}
// 	}
// }

// func TestKeys(t *testing.T) {
// 	Pack(Label(Txt("Hello World!")))
// 	Bind(App, "<Key>", Command(func(e *Event) { fmt.Printf("%q\n", e.Keysym) }))
// 	App.Wait()
// }

// func TestTmp(t *testing.T) {
// }
