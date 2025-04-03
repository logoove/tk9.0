package main

import . "modernc.org/tk9.0"
import _ "modernc.org/tk9.0/themes/azure"

func main() {
	pw := TPanedwindow(Orient("horizontal"))
	left := pw.TButton(Txt("Left"))
	right := pw.TButton(Txt("Right"))
	pw.Add(left.Window)
	pw.Add(right.Window)
	Pack(pw,
		TExit(),
		Padx("1m"), Pady("2m"), Ipadx("1m"), Ipady("1m"), Expand(1))
	ActivateTheme("azure light")
	App.Configure(Width("10c"))
	App.SetResizable(false, false)
	App.Wait()
}
