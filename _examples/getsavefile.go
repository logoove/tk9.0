package main

import "fmt"
import . "modernc.org/tk9.0"

func main() {
	Pack(
		TButton(Txt("Save As..."), Command(func() {
			fmt.Printf("%q\n", GetSaveFile(
				Title("Save File"),
				Confirmoverwrite(true),
				Filetypes([]FileType{{"Go files", []string{".go"}, ""}})),
			)
		})),
		TExit(),
		Padx("1m"), Pady("2m"), Ipadx("1m"), Ipady("1m"),
	)
	App.Wait()
}
