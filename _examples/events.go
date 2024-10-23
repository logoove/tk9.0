package main

import (
	"fmt"

	. "modernc.org/tk9.0"
)

func click(e *Event) {
	fmt.Printf("click %+v\n", e)
}

func focusIn(e *Event) {
	fmt.Printf("focus in %+v\n", e)
}

func focusOut(e *Event) {
	fmt.Printf("focus out %+v\n", e)
}

func main() {
	b1 := TButton(Txt("Hello"), Command(click))
	b2 := TButton(Txt("World"), Command(click))
	fmt.Printf("Button Hello = %s, button World = %s\n", b1, b2)
	opts := Opts{Padx("1m"), Pady("2m"), Ipadx("1m"), Ipady("1m")}
	Grid(b1, b2, opts)
	Grid(TExit(), Columnspan(2), opts)
	Bind(App, "<FocusIn>", Command(focusIn))
	Bind(App, "<FocusOut>", Command(focusOut))
	App.Wait()
}
