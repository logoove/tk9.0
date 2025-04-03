// https://wiki.tcl-lang.org/page/bitmap

package main

import . "modernc.org/tk9.0"
import _ "modernc.org/tk9.0/themes/azure"

func main() {
	ActivateTheme("azure light")
	c := Canvas(Background(Linen), Width(200), Height(200))
	c.CreateBitmap(100, 100, Bitmap("hourglass"))
	Pack(c,
		TExit(),
		Padx("1m"), Pady("2m"), Ipadx("1m"), Ipady("1m"))
	App.SetResizable(false, false)
	App.Wait()
}
