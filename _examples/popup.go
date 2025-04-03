package main

import . "modernc.org/tk9.0"
import _ "modernc.org/tk9.0/themes/azure"

func main() {
	Initialize()
	label1 := TLabel(Txt("Left-click me!"))
	label2 := TLabel(Txt("Right-click me!"))
	label3 := TLabel(Txt("No menu selected."))
	menu := Menu()
	ex1 := menu.AddCommand(Lbl("Example 1"), Command(func(e *Event) {
		label3.Configure(Txt("Menu 'Example 1' selected."))
	}))
	ex2 := menu.AddCommand(Lbl("Example 2"), Command(func(e *Event) {
		label3.Configure(Txt("Menu 'Example 2' selected."))
	}))
	Bind(label1, "<Button-1>", Command(func(e *Event) {
		Popup(menu.Window, e.XRoot, e.YRoot, ex1)
	}))
	Bind(label2, "<Button-3>", Command(func(e *Event) {
		Popup(menu.Window, e.XRoot, e.YRoot, ex2)
	}))
	Pack(label1,
		label2,
		label3,
		TExit(),
		Padx("1m"), Pady("2m"), Ipadx("1m"), Ipady("1m"))
	ActivateTheme("azure light")
	App.SetResizable(false, false)
	App.Wait()
}
