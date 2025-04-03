package main

import . "modernc.org/tk9.0"
import _ "modernc.org/tk9.0/themes/azure"

var cm = int(TkScaling()*72/2.54 + 0.5)

func main() {
	Pack(Label(Image(NewPhoto(Width(20*cm), Height(15*cm)).Graph("set grid; splot x**2+y**2, x**2-y**2"))),
		TExit(),
		Padx("1m"), Pady("2m"), Ipadx("1m"), Ipady("1m"))
	ActivateTheme("azure light")
	App.Center().Wait()
}
