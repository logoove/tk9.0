// ctext demo

/*
See https://core.tcl-lang.org/tklib/doc/trunk/embedded/md/tklib/files/modules/ctext/ctext.md

	Using color scheme derived from https://github.com/berni23/berni-dark
*/
package main

import (
	. "modernc.org/tk9.0"
	. "modernc.org/tk9.0/extensions/autoscroll"
	. "modernc.org/tk9.0/extensions/ctext"
	_ "modernc.org/tk9.0/themes/azure"
	"os"
)

func main() {
	ActivateTheme("azure dark")
	InitializeExtension("autoscroll")
	InitializeExtension("ctext")
	var yscroll *Window
	t := Ctext(Foreground("#f8d115"), Background("#263238"), Relief("flat"), Borderwidth(0),
		Font(FixedFont), Setgrid(true), Width(130), Height(40),
		Yscrollcommand(func(e *Event) { e.ScrollSet(yscroll) }))
	t.Linemapfg("#858582")
	yscroll = Autoscroll(TScrollbar(Command(func(e *Event) { e.Yview(t) })).Window)
	t.AddHighlightClassForRegexp("ident", "#34bfc9", `{[_A-Za-z][_A-Za-z0-9]*}`)
	t.AddHighlightClassForRegexp("number", "#34bfc9", `{[0-9]+}`)
	t.AddHighlightClass("keyword", "#cf5fd8", "break", "default", "func", "interface", "select", "case", "defer",
		"go", "map", "struct", "chan", "else", "goto", "package", "switch", "const", "fallthrough", "if", "range",
		"type", "continue", "for", "import", "return", "var")
	t.AddHighlightClassForRegexp("string", "#81a892", `{"[^"]*"}`)
	t.AddHighlightClassForRegexp("string2", "#81a892", "{`.*`}")
	t.AddHighlightClassForRegexp("lineComment", "#6c7782", `{\/\/[^\n\r]*}`)
	t.AddHighlightClassForRegexp("blockComment", "#6c7782", `{/\*(.|\n)*\*/}`)
	b, _ := os.ReadFile("ctext.go")
	t.Insert("1.0", string(b))
	sty := Opts{Padx("1m"), Pady("2m"), Ipadx("1m"), Ipady("1m")}
	Grid(t, Sticky(NEWS), sty)
	Grid(yscroll, Row(0), Column(1), Sticky(NS+E), Pady("2m"), Ipady("1m"))
	GridRowConfigure(App, 0, Weight(1))
	GridColumnConfigure(App, 0, Weight(1))
	Grid(TExit(Style("Accent.TButton")), sty)
	App.Wait()
}
