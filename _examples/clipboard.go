package main

import (
	. "modernc.org/tk9.0"
)

func main() {
	style := Opts{Padx("1m"), Pady("2m"), Ipadx("1m"), Ipady("1m"), Sticky("e")}
	in := TEntry(Textvariable("abc"))
	out := TEntry()
	Grid(Label(Txt("Append:")), in, TButton(Txt("Append"), Command(func() { ClipboardAppend(in.Textvariable()) })), style)
	Grid(Label(Txt("Clipboard:")), out, TButton(Txt("Get"), Command(func() {
		r := ClipboardGet()
		if Error != nil {
			r = "<empty>"
		}
		out.Configure(Textvariable(r))
	})), style,
	)
	Grid(TExit(), Columnspan(3), style)
	App.Wait()
}
