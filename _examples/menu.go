package main

import (
	"fmt"
	. "modernc.org/tk9.0"
	"runtime"
)

func main() {
	menubar := Menu()

	fileMenu := menubar.Menu()
	fileMenu.AddCommand(Lbl("New"), Underline(0), Accelerator("Ctrl+N"))
	fileMenu.AddCommand(Lbl("Open..."), Underline(0), Accelerator("Ctrl+O"), Command(func() { GetOpenFile() }))
	Bind(App, "<Control-o>", Command(func() { fileMenu.Invoke(1) }))
	fileMenu.AddCommand(Lbl("Save"), Underline(0), Accelerator("Ctrl+S"))
	fileMenu.AddCommand(Lbl("Save As..."), Underline(5))
	fileMenu.AddCommand(Lbl("Close"), Underline(0), Accelerator("Crtl+W"))
	fileMenu.AddSeparator()
	fileMenu.AddCommand(Lbl("Exit"), Underline(1), Accelerator("Ctrl+Q"), ExitHandler())
	Bind(App, "<Control-q>", Command(func() { fileMenu.Invoke(6) }))
	menubar.AddCascade(Lbl("File"), Underline(0), Mnu(fileMenu))

	editMenu := menubar.Menu()
	editMenu.AddCommand(Lbl("Undo"))
	editMenu.AddSeparator()
	editMenu.AddCommand(Lbl("Cut"))
	editMenu.AddCommand(Lbl("Copy"))
	editMenu.AddCommand(Lbl("Paste"))
	editMenu.AddCommand(Lbl("Delete"))
	editMenu.AddCommand(Lbl("Select All"))
	menubar.AddCascade(Lbl("Edit"), Underline(0), Mnu(editMenu))

	helpMenu := menubar.Menu()
	helpMenu.AddCommand(Lbl("Help Index"))
	helpMenu.AddCommand(Lbl("About..."))
	menubar.AddCascade(Lbl("Help"), Underline(0), Mnu(helpMenu))

	App.WmTitle(fmt.Sprintf("%s on %s", App.WmTitle(""), runtime.GOOS))
	App.Configure(Mnu(menubar), Width("8c"), Height("6c")).Wait()
}
