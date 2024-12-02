package main

import (
	"fmt"
	"os"

	. "modernc.org/tk9.0"
	_ "modernc.org/tk9.0/vnc"
)

func main() {
	App.Configure(Padx(0), Pady(0))
	f := TFrame()
	Grid(f.TLabel(Txt(fmt.Sprintf("DISPLAY=%s PID=%d", os.Getenv("DISPLAY"), os.Getpid()))))
	Grid(f.TLabel())
	Grid(f.TExit(), Ipadx("1m"), Ipady("2m"))
	Grid(f, Row(1), Column(1), Sticky("n"), Pady("3m"))
	Grid(TLabel(Txt("↖")), Row(0), Column(0), Sticky("nw"))
	Grid(TLabel(Txt("↗")), Row(0), Column(2), Sticky("ne"))
	Grid(TLabel(Txt("↙")), Row(2), Column(0), Sticky("sw"))
	Grid(TLabel(Txt("↘")), Row(2), Column(2), Sticky("se"))
	GridRowConfigure(App, 1, Weight(1))
	GridColumnConfigure(App, 1, Weight(1))
	App.Wait()
}
