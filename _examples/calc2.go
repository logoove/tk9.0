// The same calculator as in calc.go with additional handling of keyboard input.
package main

import "github.com/expr-lang/expr"
import . "modernc.org/tk9.0"

var (
	out *LabelWidget
	m   = map[rune][]string{
		'*': {"<KeyPress-KP_Multiply>"},
		'+': {"<KeyPress-KP_Add>"},
		'-': {"<KeyPress-KP_Subtract>"},
		'.': {"<KeyPress-KP_Decimal>"},
		'/': {"<KeyPress-KP_Divide>"},
		'0': {"<KeyPress-KP_0>"},
		'1': {"<KeyPress-KP_1>"},
		'2': {"<KeyPress-KP_2>"},
		'3': {"<KeyPress-KP_3>"},
		'4': {"<KeyPress-KP_4>"},
		'5': {"<KeyPress-KP_5>"},
		'6': {"<KeyPress-KP_6>"},
		'7': {"<KeyPress-KP_7>"},
		'8': {"<KeyPress-KP_8>"},
		'9': {"<KeyPress-KP_9>"},
		'=': {"<KeyPress-KP_Enter>", "<KeyPress-Return>"},
		'C': {"<KeyPress-Escape>"},
	}
)

func main() {
	out = Label(Height(2), Anchor("e"), Txt("(123+232)/(123-10)"))
	Grid(out, Columnspan(4), Sticky("e"))
	var b *ButtonWidget
	for i, c := range "C()/789*456-123+0.=" {
		h := Command(func() {
			c := c
			switch c {
			case 'C':
				out.Configure(Txt(""))
			case '=':
				x, err := expr.Eval(out.Txt(), nil)
				if err != nil {
					MessageBox(Icon("error"), Msg(err.Error()), Title("Error"))
					x = ""
				}
				out.Configure(Txt(x))
			default:
				out.Configure(Txt(out.Txt() + string(c)))
			}
		})
		b = Button(Txt(string(c)), Width(-4), h)
		Grid(b, Row(i/4+1), Column(i%4), Sticky("news"), Ipadx("1.5m"), Ipady("2.6m"))
		Bind(App, string(c), h)
		for _, v := range m[c] {
			Bind(App, v, h)
		}
	}
	Grid(b, Columnspan(2))
	App.Configure(Padx(0), Pady(0)).Wait()
}
