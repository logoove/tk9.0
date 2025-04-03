package main

import (
	"fmt"

	. "modernc.org/tk9.0"
	. "modernc.org/tk9.0/extensions/eval"
	_ "modernc.org/tk9.0/themes/azure"
)

func main() {
	InitializeExtension("eval")
	style := Opts{Padx("1m"), Pady("2m"), Ipadx("1m"), Ipady("1m")}
	in := TEntry(Textvariable("expr {3 * 5}"))
	out := TLabel()
	Pack(in,
		out,
		TButton(Txt("Eval"), Command(func() {
			r, err := Eval(in.Textvariable())
			out.Configure(Textvariable(fmt.Sprintf("result: %s\nerr: %v", r, err)))
		})),
		TExit(), style)
	ActivateTheme("azure light")
	App.Wait()
}
