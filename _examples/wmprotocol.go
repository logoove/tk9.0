package main

import (
	"fmt"
	"time"

	. "modernc.org/tk9.0"
	_ "modernc.org/tk9.0/themes/azure"
)

func main() {
	lbl := TLabel(Justify("center"))
	WmProtocol(App, "WM_TAKE_FOCUS", func() {
		lbl.Configure(Txt(fmt.Sprintf(`WM_TAKE_FOCUS handler func()
invoked at %v`, time.Now().Format(time.DateTime))))
	})
	WmProtocol(App, "WM_DELETE_WINDOW", Command(func() {
		lbl.Configure(Txt(fmt.Sprintf(`WM_DELETE_WINDOW handler Command()
invoked at %v`, time.Now().Format(time.DateTime))))
	}))
	Pack(TLabel(Txt(`Try changing app focus
or closing the app window using
the Window manager close button. `), Justify("center")),
		lbl,
		TExit(),
		Padx("1m"), Pady("2m"), Ipadx("1m"), Ipady("1m"))
	ActivateTheme("azure light")
	App.SetResizable(false, false)
	App.Wait()
}
