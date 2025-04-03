//go:build !windows

package main

import "fmt"
import . "modernc.org/tk9.0"
import _ "modernc.org/tk9.0/themes/azure"

func main() {
	lbl := TLabel(Txt(fmt.Sprintf("type is %q", WmAttributes(App, Type))), Justify("center"))
	Pack(
		TButton(Txt("Type dialog"), Command(func() {
			WmAttributes(App, Type("dialog"))
		})),
		TButton(Txt("Type normal"), Command(func() {
			WmAttributes(App, Type("normal"))
		})),
		TButton(Txt("Query type"), Command(func() {
			lbl.Configure(Txt(fmt.Sprintf("type is %q", WmAttributes(App, Type))))
		})),
		lbl,
		TExit(),
		Padx("1m"), Pady("2m"), Ipadx("1m"), Ipady("1m"))
	ActivateTheme("azure light")
	App.SetResizable(false, false)
	App.Wait()
}
