package main

import (
	"fmt"
	"time"

	. "modernc.org/tk9.0"
	_ "modernc.org/tk9.0/themes/azure"
)

func main() {
	btnFire := TButton(Txt("Fire the event"))
	lblBound := TLabel()
	lblFired := TLabel()
	handler := Command(func() {
		lblFired.Configure(Txt(fmt.Sprintf("Event fired at %s", time.Now().Format(time.DateTime))))
	})
	btnBind := TButton(Txt("Bind events"), Command(func() {
		Bind(btnFire, "<Button-1>", handler)
		lblBound.Configure(Txt("Fire button is bound: true"))
	}))
	btnUnbind := TButton(Txt("Unbind events"), Command(func() {
		Bind(btnFire, "<Button-1>", "")
		lblBound.Configure(Txt("Fire button is bound: false"))
	}))
	btnBind.Invoke()
	Pack(btnFire, lblBound, lblFired, btnBind, btnUnbind, TExit(),
		Padx("1m"), Pady("2m"), Ipadx("1m"), Ipady("1m"))
	ActivateTheme("azure light")
	App.SetResizable(false, false)
	App.Wait()
}
