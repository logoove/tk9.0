package main

import "github.com/expr-lang/expr"
import . "modernc.org/tk9.0"

func main() {
	out := Label(Height(2), Anchor("e"), Txt("(123+232)/(123-10)"))
	Grid(out, Columnspan(4), Sticky("e"))
	var b *ButtonWidget
	for i, c := range "C()/789*456-123+0.=" {
		b = Button(Txt(string(c)),
			Command(
				func() {
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
				},
			),
			Width(-4))
		Grid(b, Row(i/4+1), Column(i%4), Sticky("news"), Ipadx("1.5m"), Ipady("2.6m"))
	}
	Grid(b, Columnspan(2))
	App.Configure(Padx(0), Pady(0)).Wait()
}
