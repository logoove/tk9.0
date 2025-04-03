package main

import (
	"fmt"
	"os"

	. "modernc.org/tk9.0"
)

func click(e *Event) {
	fmt.Fprintf(os.Stderr, "click %+v\n", e)
}

func focusIn(e *Event) {
	fmt.Fprintf(os.Stderr, "focus in %+v\n", e)
}

func focusOut(e *Event) {
	fmt.Fprintf(os.Stderr, "focus out %+v\n", e)
}

func keyPress(e *Event) {
	fmt.Fprintf(os.Stderr, "key press   '%s' %s\n", e.Keysym, e.State)
}

func keyRelease(e *Event) {
	fmt.Fprintf(os.Stderr, "key release '%s' %s\n", e.Keysym, e.State)
}

func main() {
	b1 := TButton(Txt("Hello"), Command(click))
	b2 := TButton(Txt("World"), Command(click))
	fmt.Fprintf(os.Stderr, "Button Hello = %s, button World = %s\n", b1, b2)
	opts := Opts{Padx("1m"), Pady("2m"), Ipadx("1m"), Ipady("1m")}
	Grid(b1, b2, opts)
	Grid(TExit(), Columnspan(2), opts)
	Bind(App, "<FocusIn>", Command(focusIn))
	Bind(App, "<FocusOut>", Command(focusOut))
	Bind(App, "<KeyPress>", Command(keyPress))
	Bind(App, "<KeyRelease>", Command(keyRelease))
	App.Wait()
}
