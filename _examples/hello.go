package main

import . "modernc.org/tk9.0"

func main() {
	Pack(Button(Txt("Hello"), Command(func() { Destroy(App) })))
	App.Wait()
}
