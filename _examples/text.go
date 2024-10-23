package main

import . "modernc.org/tk9.0"

func main() {
	var scroll *TScrollbarWidget
	t := Text(Font("helvetica", 10), Yscrollcommand(func(e *Event) { e.ScrollSet(scroll) }), Setgrid(true), Wrap("word"), Padx("2m"), Pady("2m"))
	scroll = TScrollbar(Command(func(e *Event) { e.Yview(t) }))
	Grid(t, Sticky("news"), Pady("2m"))
	Grid(scroll, Row(0), Column(1), Sticky("nes"), Pady("2m"))
	GridRowConfigure(App, 0, Weight(1))
	GridColumnConfigure(App, 0, Weight(1))
	Grid(TExit(), Padx("1m"), Pady("2m"), Ipadx("1m"), Ipady("1m"))
	t.TagConfigure("bgstipple", Background(Black), Borderwidth(0), Bgstipple(Gray12))
	t.TagConfigure("big", Font("helvetica", 12, "bold"))
	t.TagConfigure("bold", Font("helvetica", 10, "bold", "italic"))
	t.TagConfigure("center", Justify("center"))
	t.TagConfigure("color1", Background("#a0b7ce"))
	t.TagConfigure("color2", Foreground(Red))
	t.TagConfigure("margins", Lmargin1("12m"), Lmargin2("6m"), Rmargin("10m"))
	t.TagConfigure("overstrike", Overstrike(1))
	t.TagConfigure("raised", Relief("raised"), Borderwidth(1))
	t.TagConfigure("right", Justify("right"))
	t.TagConfigure("spacing", Spacing1("10p"), Spacing2("2p"), Lmargin1("12m"), Lmargin2("6m"), Rmargin("10m"))
	t.TagConfigure("sub", Offset("-2p"), Font("helvetica", 8))
	t.TagConfigure("sunken", Relief("sunken"), Borderwidth(1))
	t.TagConfigure("super", Offset("4p"), Font("helvetica", 8))
	t.TagConfigure("tiny", Font("times", 8, "bold"))
	t.TagConfigure("underline", Underline(1))
	t.TagConfigure("verybig", Font(CourierFont(), 22, "bold"))
	t.InsertML(`Text widgets like this one allow you to display information in a variety of styles. Display styles are controlled
using a mechanism called <bold>tags</bold>. Tags are just textual names that you can apply to one or more ranges of characters within a
text widget. You can configure tags with various display styles. If you do this, then the tagged characters will be displayed with the
styles you chose. The available display styles are:
<br><br><big>1. Font.</big> You can choose any system font, <verybig>large</verybig> or <tiny>small</tiny>.
<br><br><big>2. Color.</big> You can change either the <color1>background</color1> or <color2>foreground</color2> color, or
<color1><color2>both</color2></color1>.
<br><br><big>3. Stippling.</big> You can cause the <bgstipple>background</bgstipple> information to be drawn with a stipple fill instead
of a solid fill.
<br><br><big>4. Underlining.</big> You can <underline>underline</underline> ranges of text.
<br><br><big>5. Overstrikes.</big> You can <overstrike>draw lines through</overstrike> ranges of text.
<br><br><big>6. 3-D effects.</big> You can arrange for the background to be drawn with a border that makes characters appear either
<raised>raised</raised> or <sunken>sunken</sunken>.
<br><br><big>7. Justification.</big> You can arrange for lines to be displayed <br>left-justified <br><right>right-justified, or</right>
<br><center>centered.</center>
<br><br><big>8. Superscripts and subscripts.</big> You can control the vertical position of text to generate superscript effects like
10<super>n</super> or subscript effects like X<sub>i</sub>.
<br><br><big>9. Margins.</big> You can control the amount of extra space left on each side of the text
<br><br><margins>This paragraph is an example of the use of margins. It consists of a single line of text that wraps around on the
screen.  There are two separate left margin values, one for the first display line associated with the text line, and one for the
subsequent display lines, which occur because of wrapping. There is also a separate specification for the right margin, which is used to
choose wrap points for lines.</margins>
<br><br><big>10. Spacing.</big> You can control the spacing of lines with three separate parameters. "Spacing1" tells how much extra
space to leave above a line, "spacing3" tells how much space to leave below a line, and if a text line wraps, "spacing2" tells how much
space to leave between the display lines that make up the text line.
<br><spacing>These indented paragraphs illustrate how spacing can be used. Each paragraph is actually a single line in the text widget,
which is word-wrapped by the widget.</spacing>
<br><spacing>Spacing1 is set to 10 points for this text, which results in relatively large gaps between the paragraphs. Spacing2 is set
to 2 points, which results in just a bit of extra space within a pararaph. Spacing3 isn't used in this example.</spacing>
<br><spacing>To see where the space is, select ranges of text within these paragraphs. The selection highlight will cover the extra
space.</spacing>`)
	App.Center().Wait()
}
