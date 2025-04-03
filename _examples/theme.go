package main

import . "modernc.org/tk9.0"

func main() {
	style := Opts{Padx("1m"), Pady("2m"), Ipadx("1m"), Ipady("1m")}
	Pack(TButton(Txt("alt"), Command(func() { ActivateTheme("alt") })), style)
	Pack(TButton(Txt("default"), Command(func() { ActivateTheme("default") })), style)
	Pack(TButton(Txt("clam"), Command(func() { ActivateTheme("clam") })), style)
	Pack(TButton(Txt("classic"), Command(func() { ActivateTheme("classic") })), style)
	Pack(TExit(), style)
	App.SetResizable(false, false)
	App.Wait()
}
