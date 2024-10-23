package main

import (
	. "modernc.org/tk9.0"
)

func main() {
	style := Opts{Padx("1m"), Pady("2m"), Ipadx("1m"), Ipady("1m"), Sticky("e")}
	a := TEntry(Textvariable("foo"))
	b := TEntry(Textvariable("bar"))
	Grid(Label(Txt("A:")), a, TButton(Txt("Focus A"), Command(func() { Focus(a) })), style)
	Grid(Label(Txt("B:")), b, TButton(Txt("Focus B"), Command(func() { Focus(b) })), style)
	Grid(TExit(), Columnspan(3), style)
	App.Wait()
}
