package main

import "fmt"
import "slices"
import . "modernc.org/tk9.0"

func main() {
	var scroll *TScrollbarWidget
	t := Text(Wrap("none"), Setgrid(true), Yscrollcommand(func(e *Event) { e.ScrollSet(scroll) }))
	scroll = TScrollbar(Command(func(e *Event) { e.Yview(t) }))
	fonts := FontFamilies()
	slices.Sort(fonts)
	Grid(t, Sticky("news"), Pady("2m"), Ipadx("1m"), Ipady("1m"))
	Grid(scroll, Row(0), Column(1), Sticky("nes"), Pady("2m"))
	GridRowConfigure(App, 0, Weight(1))
	GridColumnConfigure(App, 0, Weight(1))
	Grid(TExit(), Padx("1m"), Pady("2m"), Ipadx("1m"), Ipady("1m"))
	m := map[string]bool{}
	for i, font := range fonts {
		if m[font] {
			continue
		}
		m[font] = true
		tag := fmt.Sprintf("t%v", i)
		t.TagConfigure(tag, Font(NewFont(Family(font))))
		t.Insert("end", font+": ", "", "Lorem ipsum dolor sit amet, consectetur adipiscing elit...\n", tag)
	}
	App.Center().Wait()
}
