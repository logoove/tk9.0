// Set the app icon using an .ico file.

package main

import _ "embed"
import . "modernc.org/tk9.0"

//go:embed testico.ico
var ico []byte

func main() {
	Pack(Button(Txt("Hello"), Command(func() { Destroy(App) })))
	App.IconPhoto(NewPhoto(Data(ico)))
	App.Wait()
}
