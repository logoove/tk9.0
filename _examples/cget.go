// hello.go with some padding added.

package main

import . "modernc.org/tk9.0"
import _ "modernc.org/tk9.0/themes/azure"

func main() {
	btn := TButton(Txt("Hello"))
	Pack(btn,
		TLabel(Txt(btn.Txt())), // Tcl: .btn cget -text
		TExit(),
		Ipadx(10), Ipady(5), Padx(20), Pady(10))
	ActivateTheme("azure light")
	App.SetResizable(false, false)
	App.Wait()
}
