package main

import . "modernc.org/tk9.0"
import _ "modernc.org/tk9.0/themes/azure"

func main() {
	Pack(Label(Image(NewPhoto(File("gopher.pcx")))),
		TExit(),
		Padx("1m"), Pady("2m"), Ipadx("1m"), Ipady("1m"))
	ActivateTheme("azure light")
	App.Center().Wait()
}
