// https://www.tutorialspoint.com/tcl-tk/tk_canvas_text.htm

package main

import . "modernc.org/tk9.0"
import _ "modernc.org/tk9.0/themes/azure"

func main() {
	ActivateTheme("azure light")
	c := Canvas(Background(Red), Width(400), Height(400))
	c.CreateText(200, 200, Fill(Yellow), Justify("center"),
		Txt("Hello World.\nHow are you?"), Font(HELVETICA, 18, "bold"))
	Pack(c,
		TExit(),
		Padx("1m"), Pady("2m"), Ipadx("1m"), Ipady("1m"))
	App.SetResizable(false, false)
	App.Wait()
}
