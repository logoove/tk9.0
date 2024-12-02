package main

import . "modernc.org/tk9.0"

func main() {
	sharedVar := Variable("")
	radio := Radiobutton(Txt("abc"), sharedVar, Value("111"))
	radio2 := Radiobutton(Txt("def"), sharedVar, Value(222))
	display := Label(Background(White), Width(15))
	Pack(display, radio, radio2,
		TButton(Txt("Read Value"), Command(func() { display.Configure(Txt(radio.Variable())) })),
		TExit(),
		Pady("2m"), Ipady("1m"))
	App.Wait()
}
