package main

import . "modernc.org/tk9.0"
import _ "modernc.org/tk9.0/themes/azure"

func main() {
	ActivateTheme("azure light")
	c := Canvas(Background(Red), Width(400), Height(400))
	text := c.Text(Width(20), Height(5))
	text.Insert("end", "Hello\nWorld!")
	c.CreateWindow(200, 200, ItemWindow(text.Window), Anchor("center"))
	Pack(c,
		TExit(),
		Padx("1m"), Pady("2m"), Ipadx("1m"), Ipady("1m"))
	App.SetResizable(false, false)
	App.Wait()
}
