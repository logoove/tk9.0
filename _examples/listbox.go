// Adapted from https://www.tutorialspoint.com/tcl-tk/tk_listbox_widget.htm
package main

import (
	. "modernc.org/tk9.0"
	_ "modernc.org/tk9.0/themes/azure"
)

func main() {
	myList := Listbox()
	label := TLabel(Txt("No choice selected"))
	Bind(myList, "<<ListboxSelect>>", Command(func() {
		label.Configure(Txt(myList.Get(myList.Curselection()[0])))
	}))
	style := Opts{Padx("1m"), Pady("2m"), Ipadx("1m"), Ipady("1m")}
	Grid(myList, Row(0), Column(0), Sticky(NEWS), style)
	Grid(label, Row(1), Column(0), Columnspan(2), style)
	Grid(TExit(), Row(2), Column(0), Columnspan(2), style)
	myList.Insert(0, "Choice1", "Choice2", "Choice3")
	ActivateTheme("azure light")
	App.SetResizable(false, false)
	App.Wait()
}
