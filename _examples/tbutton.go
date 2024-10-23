package main

import _ "embed"
import . "modernc.org/tk9.0"

//go:embed red_corner.png
var red []byte

//go:embed green_corner.png
var green []byte

func main() {
	StyleThemeUse("default")
	StyleElementCreate("Red.Corner.TButton.indicator", "image", NewPhoto(Data(red)))
	StyleElementCreate("Green.Corner.TButton.indicator", "image", NewPhoto(Data(green)))
	StyleLayout("Red.Corner.TButton",
		"Button.border", Sticky("nswe"), Border(1), Children(
			"Button.focus", Sticky("nswe"), Children(
				"Button.padding", Sticky("nswe"), Children(
					"Button.label", Sticky("nswe"),
					"Red.Corner.TButton.indicator", Side("right"), Sticky("ne")))))
	StyleLayout("Green.Corner.TButton",
		"Button.border", Sticky("nswe"), Border(1), Children(
			"Button.focus", Sticky("nswe"), Children(
				"Button.padding", Sticky("nswe"), Children(
					"Button.label", Sticky("nswe"),
					"Green.Corner.TButton.indicator", Side("right"), Sticky("ne")))))
	opts := Opts{Padx("1m"), Pady("2m"), Ipadx("1m"), Ipady("1m")}
	rb := TButton(Txt("Red"))
	gb := TButton(Txt("Green"))
	Grid(rb, gb, opts)
	Grid(TButton(Txt("Use style"), Command(func() {
		rb.Configure(Style("Red.Corner.TButton"))
		gb.Configure(Style("Green.Corner.TButton"))
	})), TExit(), opts)
	App.Wait()
}
