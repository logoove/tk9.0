package main

import (
	"fmt"
	. "modernc.org/tk9.0"
)

func main() {
	Pack(
		TButton(Txt("Select font..."), Command(func() {
			Fontchooser(Command(func() { fmt.Printf("%q\n", FontchooserFont()) }))
			FontchooserShow()
		})),
		TExit(),
		Padx("1m"), Pady("2m"), Ipadx("1m"), Ipady("1m"),
	)
	App.Wait()
}
