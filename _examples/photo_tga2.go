package main

import (
	_ "embed"

	. "modernc.org/tk9.0"
	_ "modernc.org/tk9.0/themes/azure"
)

//go:embed img.tga
var picture []byte

func main() {
	Pack(Label(Image(NewPhoto(Data(picture)))),
		TExit(),
		Padx("1m"), Pady("2m"), Ipadx("1m"), Ipady("1m"))
	ActivateTheme("azure light")
	App.Center().Wait()
}
