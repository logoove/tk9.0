package main

import . "modernc.org/tk9.0"
import _ "modernc.org/tk9.0/themes/azure"

func main() {
	Pack(TButton(Txt("Light"), Command(func() { ActivateTheme("azure light") })),
		TButton(Txt("Dark"), Command(func() { ActivateTheme("azure dark") })),
		TExit(),
		Pady("2m"), Ipady("1m"))
	App.Wait()
}
