// Original code comes from:
//
//	https://github.com/rdbende/Azure-ttk-theme/blob/main/example.py
//
// License terms:
//
// """
// Example script for testing the Azure ttk theme
// Author: rdbende
// License: MIT license
// Source: https://github.com/rdbende/ttk-widget-factory
// """

// Modification are copyright 2024 The tk9.0-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"

	. "modernc.org/tk9.0"
	_ "modernc.org/tk9.0/themes/azure"
	// This is here only to enable reusing the code for pictures in vnc/README.md.
	_ "modernc.org/tk9.0/vnc"
)

func main() {
	dark := flag.Bool("dark", false, "use dark mode")
	flag.Parse()
	switch {
	case *dark:
		ActivateTheme("azure dark")
	default:
		ActivateTheme("azure light")
	}
	Pack(app(), Fill("both"), Expand(true))
	Update()
	WmMinSize(App, WinfoWidth(App), WinfoHeight(App))
	App.Wait()
}

func app() (r *TFrameWidget) {
	r = TFrame()
	// Make the app responsive
	for i := 0; i < 3; i++ {
		GridColumnConfigure(r, i, Weight(1))
		GridRowConfigure(r, i, Weight(1))
	}

	// Create widgets

	// Create a Frame for the Checkbuttons
	checkFrame := r.TLabelframe(Txt(" Checkbuttons "), Padding("20 10"))
	Grid(checkFrame, Row(0), Column(0), Padx("20 10"), Pady("20 10"), Sticky("nsew"))

	// Checkbuttons
	check1 := checkFrame.TCheckbutton(Txt("Unchecked"), Variable(0))
	Grid(check1, Row(0), Column(0), Padx(5), Pady(10), Sticky("nsew"))
	check2 := checkFrame.TCheckbutton(Txt("Checked"), Variable(1))
	Grid(check2, Row(1), Column(0), Padx(5), Pady(10), Sticky("nsew"))
	check3 := checkFrame.TCheckbutton(Txt("Third state"), Variable(0))
	check3.WidgetState("alternate")
	Grid(check3, Row(2), Column(0), Padx(5), Pady(10), Sticky("nsew"))
	check4 := checkFrame.TCheckbutton(Txt("Disabled"), State("disabled"))
	check4.WidgetState("disabled !alternate")
	Grid(check4, Row(3), Column(0), Padx(5), Pady(10), Sticky("nsew"))

	// Separator
	separator := r.TSeparator()
	Grid(separator, Row(1), Column(0), Padx("20 10"), Pady(10), Sticky("ew"))

	// Create a Frame for the Radiobuttons
	radioFrame := r.TLabelframe(Txt(" Radiobuttons "), Padding("20 10"))
	Grid(radioFrame, Row(2), Column(0), Padx("20 10"), Pady(10), Sticky("nsew"))

	// Radiobuttons
	rv := Variable(2)
	radio1 := radioFrame.TRadiobutton(Txt("Unselected"), rv, Value(1))
	Grid(radio1, Row(0), Column(0), Padx(5), Pady(10), Sticky("nsew"))
	radio2 := radioFrame.TRadiobutton(Txt("Selected"), rv, Value(2))
	Grid(radio2, Row(1), Column(0), Padx(5), Pady(10), Sticky("nsew"))
	radio3 := radioFrame.TRadiobutton(Txt("Disabled"), State("disabled"))
	Grid(radio3, Row(2), Column(0), Padx(5), Pady(10), Sticky("nsew"))

	// Create a Frame for input widgets
	widgetsFrame := r.TFrame(Padding("0 0 0 10"))
	Grid(widgetsFrame, Row(0), Column(1), Padx(10), Pady("30 10"), Sticky("nsew"), Rowspan(3))
	GridColumnConfigure(widgetsFrame, 0, Weight(1))

	// Entry
	entry := widgetsFrame.TEntry(Textvariable("Entry"))
	Grid(entry, Row(0), Column(0), Padx(5), Pady("0 10"), Sticky("ew"))

	// Spinbox
	spinbox := widgetsFrame.TSpinbox(Textvariable("Spinbox"), From(0), To(100), Increment(0.1))
	Grid(spinbox, Row(1), Column(0), Padx(5), Pady(10), Sticky("ew"))

	// Combobox
	combobox := widgetsFrame.TCombobox(Values("Combobox {Editable item 1} {Editable item 2}"))
	combobox.Current(0)
	Grid(combobox, Row(2), Column(0), Padx(5), Pady(10), Sticky("ew"))

	// Read-only combobox
	readonlyCombo := widgetsFrame.TCombobox(State("readonly"), Values("{Readonly combobox} {Item 1} {Item 2}"))
	readonlyCombo.Current(0)
	Grid(readonlyCombo, Row(3), Column(0), Padx(5), Pady(10), Sticky("ew"))

	// Menubutton
	menubutton := widgetsFrame.Menubutton(Txt("Menubutton"), Direction("below"))
	Grid(menubutton, Row(4), Column(0), Padx(5), Pady(10), Sticky("nsew"))

	// Menu for the Menubutton
	menu := menubutton.Menu()
	menu.AddCommand(Lbl("Menu item 1"))
	menu.AddCommand(Lbl("Menu item 2"))
	menu.AddSeparator()
	menu.AddCommand(Lbl("Menu item 3"))
	menu.AddCommand(Lbl("Menu item 4"))
	menubutton.Configure(Mnu(menu))

	// OptionMenu
	optionmenu := widgetsFrame.OptionMenu(Variable(nil), "", "OptionMenu", "Option 1", "Option 2")
	Grid(optionmenu, Row(5), Column(0), Padx(5), Pady(10), Sticky("nsew"))

	// Button
	button := widgetsFrame.TButton(Txt("Button"))
	Grid(button, Row(6), Column(0), Padx(5), Pady(10), Sticky("nsew"))

	// Accentbutton
	accentButton := widgetsFrame.TButton(Txt("Accent button"), Style("Accent.TButton"))
	Grid(accentButton, Row(7), Column(0), Padx(5), Pady(10), Sticky("nsew"))

	// Togglebutton
	toggleButton := widgetsFrame.TCheckbutton(Txt("Toggle button"), Style("Toggle.TButton"))
	Grid(toggleButton, Row(8), Column(0), Padx(5), Pady(10), Sticky("nsew"))

	// Switch
	swtch := widgetsFrame.TCheckbutton(Txt("Switch"), Style("Switch.TCheckbutton"))
	Grid(swtch, Row(9), Column(0), Padx(5), Pady(10), Sticky("nsew"))

	// Exit button
	Grid(widgetsFrame.TExit(), Row(10), Column(0), Padx(5), Pady(10), Sticky("nsew"))

	// Panedwindow
	paned := r.TPanedwindow()
	Grid(paned, Row(0), Column(2), Pady("25 5"), Sticky("nsew"), Rowspan(3))

	// Pane #1
	pane1 := paned.TFrame(Padding(5))
	paned.Add(pane1.Window, Weight(1))

	// Scrollbar
	scrollbar := pane1.TScrollbar()
	Pack(scrollbar, Side("right"), Fill("y"))

	// Treeview
	treeview := pane1.TTreeview(Selectmode("browse"), Columns("1 2"), Height(10),
		Yscrollcommand(func(e *Event) { e.ScrollSet(scrollbar) }))
	Pack(treeview, Expand(true), Fill("both"))
	scrollbar.Configure(Command(func(e *Event) { e.Yview(treeview) }))

	// Treeview columns
	treeview.Column("#0", Anchor("w"), Width(120))
	treeview.Column(1, Anchor("w"), Width(120))
	treeview.Column(2, Anchor("w"), Width(120))

	// Treeview headings
	treeview.Heading("#0", Txt("Column 1"), Anchor("center"))
	treeview.Heading(1, Txt("Column 2"), Anchor("center"))
	treeview.Heading(2, Txt("Column 3"), Anchor("center"))

	// Define treeview data
	treeviewData := [...][4]any{
		{"", 1, "Parent", "{Item 1} {Value 1}"},
		{1, 2, "Child", "{Subitem 1.1} {Value 1.1}"},
		{1, 3, "Child", "{Subitem 1.2} {Value 1.2}"},
		{1, 4, "Child", "{Subitem 1.3} {Value 1.3}"},
		{1, 5, "Child", "{Subitem 1.4} {Value 1.4}"},
		{"", 6, "Parent", "{Item 2} {Value 2}"},
		{6, 7, "Child", "{Subitem 2.1} {Value 2.1}"},
		{6, 8, "Sub-parent", "{Subitem 2.2} {Value 2.2}"},
		{8, 9, "Child", "{Subitem 2.2.1} {Value 2.2.1}"},
		{8, 10, "Child", "{Subitem 2.2.2} {Value 2.2.2}"},
		{8, 11, "Child", "{Subitem 2.2.3} {Value 2.2.3}"},
		{6, 12, "Child", "{Subitem 2.3} {Value 2.3}"},
		{6, 13, "Child", "{Subitem 2.4} {Value 2.4}"},
		{"", 14, "Parent", "{Item 3} {Value 3}"},
		{14, 15, "Child", "{Subitem 3.1} {Value 3.1}"},
		{14, 16, "Child", "{Subitem 3.2} {Value 3.2}"},
		{14, 17, "Child", "{Subitem 3.3} {Value 3.3}"},
		{14, 18, "Child", "{Subitem 3.4} {Value 3.4}"},
		{"", 19, "Parent", "{Item 4} {Value 4}"},
		{19, 20, "Child", "{Subitem 4.1} {Value 4.1}"},
		{19, 21, "Sub-parent", "{Subitem 4.2} {Value 4.2}"},
		{21, 22, "Child", "{Subitem 4.2.1} {Value 4.2.1}"},
		{21, 23, "Child", "{Subitem 4.2.2} {Value 4.2.2}"},
		{21, 24, "Child", "{Subitem 4.2.3} {Value 4.2.3}"},
		{19, 25, "Child", "{Subitem 4.3} {Value 4.3}"},
	}

	// Insert treeview data
	for _, item := range treeviewData {
		treeview.Insert(item[0], "end", Id(item[1]), Txt(item[2]), Values(item[3]))
		if item[0] == "" || item[1] == 8 || item[1] == 21 {
			treeview.Item(item[1], Open(true))
		}
	}

	// Select and scroll
	treeview.Selection("set", 10)
	treeview.See(7)

	// Notebook, pane #2
	pane2 := paned.TFrame(Padding(5))
	paned.Add(pane2.Window, Weight(3))

	// Notebook, pane #2
	notebook := pane2.TNotebook()
	Pack(notebook, Fill("both"), Expand(true))

	// Tab #1
	tab1 := notebook.TFrame()
	for index := 0; index < 2; index++ {
		GridColumnConfigure(tab1, index, Weight(1))
		GridRowConfigure(tab1, index, Weight(1))
	}
	notebook.Add(tab1, Txt("Tab 1"))

	// Scale
	sv := Variable(nil)
	scale := tab1.TScale(From(100), To(0), sv)
	scale.Configure(Command(func() { sv.Set(scale.Get()) }))
	Grid(scale, Row(0), Column(0), Padx("20  10"), Pady("20 0"), Sticky("ew"))

	// Progressbar
	progress := tab1.TProgressbar(Value(0), sv, Mode("determinate"))
	Grid(progress, Row(0), Column(1), Padx("10 20"), Pady("20 0"), Sticky("ew"))

	// Label
	label := tab1.TLabel(Txt(fmt.Sprintf("%q theme for tk9.0", CurrentThemeName())), Justify("center"), Font("helvetica", 15, "bold"))
	Grid(label, Row(1), Column(0), Pady(10), Columnspan(2))

	// Tab #2
	tab2 := notebook.TFrame()
	notebook.Add(tab2, Txt("Tab 2"))

	// Tab #3
	tab3 := notebook.TFrame()
	notebook.Add(tab3, Txt("Tab 3"))

	// Sizegrip
	sizegrip := r.TSizegrip()
	Grid(sizegrip, Row(100), Column(100), Padx("0 5"), Pady("0 5"))

	return r
}
