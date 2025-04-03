package main

import "fmt"
import . "modernc.org/tk9.0"
import _ "modernc.org/tk9.0/themes/azure"

func main() {
	fr := TFrame()

	// Scrollbar
	sb := fr.TScrollbar()
	Pack(sb, Side("right"), Fill("y"))

	// Treeview
	tv := fr.TTreeview(Selectmode("extended"), Columns("1 2"), Height(10),
		Yscrollcommand(func(e *Event) { e.ScrollSet(sb) }))
	Pack(tv, Expand(true), Fill("both"))
	sb.Configure(Command(func(e *Event) { e.Yview(tv) }))

	// Treeview columns
	tv.Column("#0", Anchor("w"), Width(120))
	tv.Column(1, Anchor("w"), Width(120))
	tv.Column(2, Anchor("w"), Width(120))

	// Treeview headings
	tv.Heading("#0", Txt("Column 1"), Anchor("center"))
	tv.Heading(1, Txt("Column 2"), Anchor("center"))
	tv.Heading(2, Txt("Column 3"), Anchor("center"))

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
		tv.Insert(item[0], "end", Id(item[1]), Txt(item[2]), Values(item[3]))
		if item[0] == "" || item[1] == 8 || item[1] == 21 {
			tv.Item(item[1], Open(true))
		}
	}

	tv.TagAdd("blue", "7")
	tv.TagConfigure("blue", Foreground(Blue))
	lbl := TLabel(Txt("Select a treeview item"))
	var children []string
	var del *TButtonWidget
	del = TButton(Txt("Delete []"), Command(func() {
		tv.Delete(children)
		children = nil
		del.Configure(Txt("Delete []"))
	}))
	Bind(tv, "<<TreeviewSelect>>", Command(func() {
		list := tv.Selection("")
		var s string
		if len(list) != 0 {
			sel := list[0]
			children = tv.Children(sel)
			s = fmt.Sprintf(`selected=%q parent=%q index=%v
children=%v
text=%q values=%q
tag("blue", Foreground)=%s focus=%v`,
				sel, tv.Parent(sel), tv.Index(sel),
				children,
				tv.Item(sel, Txt), tv.Item(sel, Values),
				tv.TagConfigure("blue", Foreground), tv.Focus())
		}
		lbl.Configure(Txt(s))
		del.Configure(Txt(fmt.Sprintf("Delete %v", children)))
	}))
	menu := Menu()
	var deleteID string
	menu.AddCommand(Lbl("Delete"), Command(func() {
		defer func() {
			deleteID = ""
		}()

		if deleteID != "" {
			tv.Delete(deleteID)
			lbl.Configure(Txt(fmt.Sprintf("deleted item %q", deleteID)))
		}
	}))
	Bind(tv, "<Button-3>", Command(func(e *Event) {
		deleteID = tv.IdentifyItem(e.X, e.Y)
		Popup(menu.Window, e.XRoot, e.YRoot, nil)
	}))
	Pack(fr,
		lbl,
		del,
		TButton(Txt("Select All"), Command(func() { tv.Selection("set", all(tv, "")) })),
		TButton(Txt("Clear"), Command(func() { tv.Delete(tv.Children("")) })),
		TButton(Txt("Focus 2"), Command(func() {
			tv.Focus(2)
			Focus(tv)
		})),
		TExit(),
		Padx("1m"), Pady("2m"), Ipadx("1m"), Ipady("1m"))
	ActivateTheme("azure light")
	App.SetResizable(false, false)
	App.Wait()
}

func all(t *TTreeviewWidget, root any) (r []string) {
	for _, v := range t.Children(root) {
		if s := fmt.Sprint(v); s != "" {
			r = append(r, s)
			r = append(r, all(t, s)...)
		}
	}
	return r
}
