package main

import "fmt"
import . "modernc.org/tk9.0"
import _ "modernc.org/tk9.0/themes/azure"

func main() {
	label := TLabel()
	Pack(label,
		TExit(),
		Padx("1m"), Pady("2m"), Ipadx("1m"), Ipady("1m"), Expand(1))
	ActivateTheme("azure light")
	App.Configure(Width("10c"))
	App.SetResizable(false, false)
	label.Configure(Txt(fmt.Sprintf("App children windows:\n%v", WinfoChildren(App))))
	App.Wait()
}
