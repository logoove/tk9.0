package main

import . "modernc.org/tk9.0"
import _ "modernc.org/tk9.0/themes/azure"

func main() {
	tex := `$$\int _0 ^\infty {{\sin ax \sin bx}\over{x^2}}\,dx = {\pi a\over 2}$$`
	Pack(Label(Relief("sunken"), Image(NewPhoto(Data(TeX(tex, 2*TkScaling()*72/600))))),
		TExit(),
		Padx("1m"), Pady("2m"), Ipadx("1m"), Ipady("1m"))
	ActivateTheme("azure light")
	App.Center().Wait()
}
