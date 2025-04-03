package main

import (
	. "modernc.org/tk9.0"
	. "modernc.org/tk9.0/extensions/tablelist"
	_ "modernc.org/tk9.0/themes/azure"
)

func main() {
	InitializeExtension("tablelist")
	style := Opts{Padx("1m"), Pady("2m"), Ipadx("1m"), Ipady("1m")}
	t := Tablelist0(`.t -columns {0 "First Column" 0 "Another column"}`)
	t.Do(`insert end [list "first row" "another value"]`)
	t.Do(`insert end [list "another row" "bla bla"]`)
	t.Do(`insert end [list "more rowz" "ha ha ha"]`)
	Grid(t, Sticky("news"), style)
	GridColumnConfigure(App, 0, Weight(1))
	GridRowConfigure(App, 0, Weight(1))
	Grid(TLabel(Txt("tablelist version: "+Version)), style)
	Grid(TExit(), style)
	ActivateTheme("azure light")
	App.Wait()
}
