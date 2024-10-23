package main

import "time"
import . "modernc.org/tk9.0"

func main() {
	lbl := Label()
	NewTicker(100*time.Millisecond, func() {
		lbl.Configure(Txt(time.Now().Format(time.DateTime)))
	})
	Pack(lbl, TExit(), Padx("1m"), Pady("2m"), Ipadx("1m"), Ipady("1m"))
	App.Wait()
}
