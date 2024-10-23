package main

import . "modernc.org/tk9.0"

func main() {
	f := Frame(Width("50m"), Height("10m"))
	a := f.TEntry(Textvariable("AAA"))
	b := f.TEntry(Textvariable("BBB"))
	Place(a, X("1m"), Y("1m"))
	Place(b, X("4m"), Y("4m"))
	Pack(
		f,
		TButton(Txt("Raise A"), Command(func() { a.Raise(b) })),
		TButton(Txt("Raise B"), Command(func() { b.Raise(a) })),
		TExit(),
		Padx("1m"), Pady("2m"), Ipadx("1m"), Ipady("1m"),
	)
	App.Wait()
}
