package main

import . "modernc.org/tk9.0"

func main() {
	busy := false
	b := TButton(Txt("Ready"))
	Pack(b,
		TButton(Txt("Change"), Command(func() {
			switch {
			case busy:
				b.BusyForget()
				b.Configure(Txt("Ready"))
			default:
				b.Busy(Watch)
				b.Configure(Txt("Busy"))
			}
			busy = !busy
		})),
		TExit(),
		Padx("1m"), Pady("2m"), Ipadx("1m"), Ipady("1m"))
	App.Wait()
}
