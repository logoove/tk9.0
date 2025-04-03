package main

import "fmt"
import . "modernc.org/tk9.0"
import _ "modernc.org/tk9.0/themes/azure"

func main() {
	lbl := TLabel(Txt(fmt.Sprintf("topmost is %v", WmAttributes(App, Topmost))), Justify("center"))
	Pack(
		TButton(Txt("Topmost 1"), Command(func() {
			WmAttributes(App, Topmost(true))
		})),
		TButton(Txt("Topmost 0"), Command(func() {
			WmAttributes(App, Topmost(false))
		})),
		TButton(Txt("Query topmost"), Command(func() {
			lbl.Configure(Txt(fmt.Sprintf("topmost is %v", WmAttributes(App, Topmost))))
		})),
		lbl,
		TExit(),
		Padx("1m"), Pady("2m"), Ipadx("1m"), Ipady("1m"))
	ActivateTheme("azure light")
	App.SetResizable(false, false)
	App.Wait()
}
