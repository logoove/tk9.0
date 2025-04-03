package main

import _ "embed"
import . "modernc.org/tk9.0"
import _ "modernc.org/tk9.0/themes/azure"

//go:embed gopher.png
var gopher []byte

func main() {
	Pack(Label(Image(NewPhoto(Data(gopher)))),
		TExit(),
		Padx("1m"), Pady("2m"), Ipadx("1m"), Ipady("1m"))
	ActivateTheme("azure light")
	App.Center().Wait()
}
