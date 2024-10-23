package main

import (
	"os"

	. "modernc.org/tk9.0"
)

func main() {
	home, _ := os.UserHomeDir()
	lbl := Label(Txt(home))
	Pack(
		lbl,
		TButton(Txt("Choose directory"), Command(func() {
			lbl.Configure(Txt(ChooseDirectory(Initialdir(lbl.Txt()))))
		})),
		TExit(),
		Padx("1m"), Pady("2m"), Ipadx("1m"), Ipady("1m"),
	)
	App.Wait()
}
