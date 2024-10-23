package main

import . "modernc.org/tk9.0"
import _ "embed"

//go:embed gotk.png
var icon []byte

func main() {
	fontSize := int(10*TkScaling()/NativeScaling + 0.5)
	font := Font("helvetica", fontSize)
	var scroll *TScrollbarWidget
	t := Text(font, Height(20), Yscrollcommand(func(e *Event) { e.ScrollSet(scroll) }), Setgrid(true), Wrap("word"),
		Padx("4p"), Pady("12p"))
	scroll = TScrollbar(Command(func(e *Event) { e.Yview(t) }))
	Grid(t, Sticky("news"), Pady("2m"))
	Grid(scroll, Row(0), Column(1), Sticky("nes"), Pady("2m"))
	GridRowConfigure(App, 0, Weight(1))
	GridColumnConfigure(App, 0, Weight(1))
	Grid(Exit().Configure(font), Padx("1m"), Pady("2m"), Ipadx("1m"), Ipady("1m"))
	t.TagConfigure("c", Justify("center"))
	t.TagConfigure("e", Offset("-2p"))
	t.TagConfigure("t", Font("times", fontSize))
	sym := " <t>T<e>E</e>X</t> "
	tex := `$Q(\xi) = \lambda_1 y_1^2 \sum_{i=2}^n \sum_{j=2}^n y_i b_{ij} y_j$`
	t.InsertML(`<c>Hello Go + Tk`, NewPhoto(Data(icon)), Padx("4p"), `users!
<br><br>Hello Go + Tk +`, sym, tex, ` users! (\$inline math\$)
<br><br>Hello Go + Tk +`, sym, `$`+tex+`$`, ` users! (\$\$display math\$\$)</c>
<br><br>The above exemplifies embeding pictures and`, sym, `scripts. A text widget can also embed other widgets. For example,
when a`, Button(font, Txt("<Tbutton>")), Padx("4p"), Pady("2p"), Align("center"), `and
a`, Entry(font, Textvariable("<TEntry>"), Background(White), Width(8)), Padx("4p"), Pady("2p"), Align("center"), `are part of
the markup, they will reflow when their containing text widget is resized.`)
	App.Center().Wait()
}
