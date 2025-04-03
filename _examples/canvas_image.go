// https://wiki.tcl-lang.org/page/bitmap

package main

import _ "embed"
import . "modernc.org/tk9.0"
import _ "modernc.org/tk9.0/themes/azure"

//go:embed gotk.png
var png []byte

func main() {
	ActivateTheme("azure light")
	c := Canvas(Background(Red), Width(200), Height(200))
	img := NewPhoto(Data(png))
	c.CreateImage(100, 100, Image(img))
	Pack(c,
		TExit(),
		Padx("1m"), Pady("2m"), Ipadx("1m"), Ipady("1m"))
	App.SetResizable(false, false)
	App.Wait()
}
