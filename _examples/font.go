package main

import (
	"fmt"
	"slices"
	"strings"

	. "modernc.org/tk9.0"
	_ "modernc.org/tk9.0/themes/azure"
)

func main() {
	App.WmTitle("Standard and System Fonts")
	var scroll *TScrollbarWidget
	t := Text(Wrap("none"), Setgrid(true), Yscrollcommand(
		func(e *Event) { e.ScrollSet(scroll) }))
	scroll = TScrollbar(Command(func(e *Event) { e.Yview(t) }))
	Grid(t, Sticky("news"), Pady("2m"), Ipadx("1m"), Ipady("1m"))
	Grid(scroll, Row(0), Column(1), Sticky("nes"), Pady("2m"))
	GridRowConfigure(App, 0, Weight(1))
	GridColumnConfigure(App, 0, Weight(1))
	Grid(TExit(), Padx("1m"), Pady("2m"), Ipadx("1m"), Ipady("1m"))
	insertFontFamilies(t)
	ActivateTheme("azure light")
	Bind(App, "<Escape>", Command(func() { Destroy(App) }))
	App.Center().Wait()
}

func insertFontFamilies(t *TextWidget) {
	m := map[string]bool{}
	std := true
	for i, family := range getFontFamilies() {
		if m[family] {
			continue
		}
		if std && !strings.HasPrefix(family, "Tk") {
			std = false
			t.Insert(END, "\n")
		}
		m[family] = true
		tag := fmt.Sprintf("t%v", i)
		t.TagConfigure(tag, Font(NewFont(Family(family))))
		t.Insert(END, family+": ", "",
			"The quick brown fox jumped over the lazy dogs. 0O1lZ2UVuv\n",
			tag)
	}
}

func getFontFamilies() []string {
	families := FontFamilies()
	for _, family := range []string{
		DefaultFont, TextFont, FixedFont, MenuFont, HeadingFont,
		CaptionFont, SmallCaptionFont, IconFont, TooltipFont,
	} {
		families = append(families, family)
	}
	slices.SortFunc(families, func(a, b string) int { // put std fonts first
		if strings.HasPrefix(a, "Tk") {
			a = " " + a
		}
		if strings.HasPrefix(b, "Tk") {
			b = " " + b
		}
		return strings.Compare(a, b)
	})
	return families
}
