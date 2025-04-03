package main

import (
	"fmt"
	"time"

	. "modernc.org/tk9.0"
	_ "modernc.org/tk9.0/themes/azure"
)

func main() {
	lbl := TLabel(Txt("Wide Awake!"))
	Pack(TButton(Txt("Click to Sleep for 2 sec"), Command(func() {
		lbl.Configure(Txt("Sleeping …"))
		Update()
		t := time.Now()
		TclAfterIdle(func() {
			lbl.Configure(Txt(fmt.Sprintf("Started sleeping at %s.\nBecome idle at %s",
				t.Format(time.DateTime), time.Now().Format(time.DateTime))))
		})
		TclAfter(time.Second * 2)
	})),
		lbl,
		TExit(),
		Padx("1m"), Pady("2m"), Ipadx("1m"), Ipady("1m"))
	ActivateTheme("azure light")
	App.SetResizable(false, false)
	App.Wait()
}
