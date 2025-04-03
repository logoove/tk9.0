//https://www.tutorialspoint.com/tcl-tk/tk_canvas_arc.htm/1000

package main

import . "modernc.org/tk9.0"
import _ "modernc.org/tk9.0/themes/azure"

func main() {
	ActivateTheme("azure light")
	c := Canvas(Background(Red), Width(200), Height(200))
	c.CreateArc(20, 20, 180, 180, Fill(Yellow), Dash(3, 3), Width(3))
	Pack(c,
		TExit(),
		Padx("1m"), Pady("2m"), Ipadx("1m"), Ipady("1m"))
	App.SetResizable(false, false)
	App.Wait()
}
