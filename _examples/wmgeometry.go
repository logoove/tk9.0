package main

import . "modernc.org/tk9.0"
import _ "modernc.org/tk9.0/themes/azure"

func main() {
	style := Opts{Padx("1m"), Pady("2m"), Ipadx("1m"), Ipady("1m")}
	in := TEntry()
	Pack(in, style)
	Pack(TButton(Txt("Get"), Command(func() {
		in.Configure(Textvariable(WmGeometry(App)))
	})), style)
	Pack(TButton(Txt("Set"), Command(func() {
		WmGeometry(App, in.Textvariable())
	})), style)
	Pack(TButton(Txt("Default"), Command(func() {
		in.Configure(Textvariable(WmGeometry(App, "")))
	})), style)
	Pack(TExit(), style)
	ActivateTheme("azure light")
	App.Wait()
}
