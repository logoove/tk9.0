// https://www.tutorialspoint.com/tcl-tk/tk_canvas_rectangle.htm

package main

import . "modernc.org/tk9.0"
import _ "modernc.org/tk9.0/themes/azure"

func main() {
	ActivateTheme("azure light")
	c := Canvas(Background(Red), Width(200), Height(200))
	c.CreateRectangle(50, 50, 100, 80, Fill(Yellow))
	Pack(c,
		TExit(),
		Padx("1m"), Pady("2m"), Ipadx("1m"), Ipady("1m"))
	App.SetResizable(false, false)
	App.Wait()
}
