// Set the app icon using an .png file.

package main

import _ "embed"
import . "modernc.org/tk9.0"

//go:embed testico.png
var ico []byte

func main() {
	Pack(Button(Txt("Hello"), Command(func() { Destroy(App) })))
	App.IconPhoto(NewPhoto(Data(ico)))
	App.Wait()
}
