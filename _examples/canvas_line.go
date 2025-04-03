// https://www.tutorialspoint.com/tcl-tk/tk_canvas_line.htm

package main

import . "modernc.org/tk9.0"
import _ "modernc.org/tk9.0/themes/azure"

func main() {
	ActivateTheme("azure light")
	c := Canvas(Background(Red), Width(100), Height(100))
	c.CreateLine(10, 10, 50, 50, 30, 100, Linearrow("both"), Fill(Yellow), Smooth(true), Splinesteps(2))
	Pack(c,
		TExit(),
		Padx("1m"), Pady("2m"), Ipadx("1m"), Ipady("1m"))
	App.SetResizable(false, false)
	App.Wait()
}
