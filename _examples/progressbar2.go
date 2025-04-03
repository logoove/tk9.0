package main

import "time"
import . "modernc.org/tk9.0"
import _ "modernc.org/tk9.0/themes/azure"

func main() {
	ActivateTheme("azure light")
	pb := TProgressbar()
	Pack(pb,
		TButton(Txt("Start"), Command(func() { pb.Start(300 / 10 * time.Millisecond) })),
		TButton(Txt("Stop"), Command(func() { pb.Stop() })),
		TExit(),
		Padx("1m"), Pady("2m"), Ipadx("1m"), Ipady("1m"))
	App.SetResizable(false, false)
	App.Wait()
}
