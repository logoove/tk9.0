package main

import . "modernc.org/tk9.0"
import _ "modernc.org/tk9.0/themes/azure"

func main() {
	entry := TEntry(Width(20), Textvariable("Initial value"))
	label := TLabel(Background(LightBlue))
	Pack(TLabel(Txt(`Clicking the "Move" button copies the content
of the entry widget to the label bellow it
and clears the entry.`), Justify("center")),
		entry,
		label,
		TButton(Txt("Move"), Command(func() {
			label.Configure(Txt(entry.Textvariable()))
			entry.Configure(Textvariable(""))
		})),
		TExit(),
		Ipadx(10), Ipady(5), Padx(20), Pady(10))
	ActivateTheme("azure light")
	App.Wait()
}
