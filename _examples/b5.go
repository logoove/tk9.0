package main

import (
	. "modernc.org/tk9.0"
	"modernc.org/tk9.0/b5"
)

func main() {
	background := White
	primary := b5.Colors{b5.ButtonText: "#fff", b5.ButtonFace: "#0d6efd", b5.ButtonFocus: "#98c1fe"}
	secondary := b5.Colors{b5.ButtonText: "#fff", b5.ButtonFace: "#6c757d", b5.ButtonFocus: "#c0c4c8"}
	success := b5.Colors{b5.ButtonText: "#fff", b5.ButtonFace: "#198754", b5.ButtonFocus: "#9dccb6"}
	danger := b5.Colors{b5.ButtonText: "#fff", b5.ButtonFace: "#dc3545", b5.ButtonFocus: "#f0a9b0"}
	warning := b5.Colors{b5.ButtonText: "#000", b5.ButtonFace: "#ffc107", b5.ButtonFocus: "#ecd182"}
	info := b5.Colors{b5.ButtonText: "#000", b5.ButtonFace: "#0dcaf0", b5.ButtonFocus: "#85d5e5"}
	light := b5.Colors{b5.ButtonText: "#000", b5.ButtonFace: "#f8f9fa", b5.ButtonFocus: "#e9e9ea"}
	dark := b5.Colors{b5.ButtonText: "#fff", b5.ButtonFace: "#212529", b5.ButtonFocus: "#a0a2a4"}
	link := b5.Colors{b5.ButtonText: "#1774fd", b5.ButtonFace: "#fff", b5.ButtonFocus: "#c2dbfe"}
	StyleThemeUse("default")
	opts := Opts{Padx("1m"), Pady("2m"), Ipadx("1m"), Ipady("1m")}
	Grid(TButton(Txt("Primary"), Style(b5.ButtonStyle("primary.TButton", primary, background, false))),
		TButton(Txt("Secondary"), Style(b5.ButtonStyle("secondary.TButton", secondary, background, false))),
		TButton(Txt("Success"), Style(b5.ButtonStyle("success.TButton", success, background, false))),
		opts)
	Grid(TButton(Txt("Danger"), Style(b5.ButtonStyle("danger.TButton", danger, background, false))),
		TButton(Txt("Warning"), Style(b5.ButtonStyle("warning.TButton", warning, background, false))),
		TButton(Txt("Info"), Style(b5.ButtonStyle("info.TButton", info, background, false))),
		opts)
	Grid(TButton(Txt("Light"), Style(b5.ButtonStyle("light.TButton", light, background, false))),
		TButton(Txt("Dark"), Style(b5.ButtonStyle("dark.TButton", dark, background, false))),
		TButton(Txt("Link"), Style(b5.ButtonStyle("link.TButton", link, background, false))),
		opts)
	Grid(TButton(Txt("Primary"), Style(b5.ButtonStyle("focused.primary.TButton", primary, background, true))),
		TButton(Txt("Secondary"), Style(b5.ButtonStyle("focused.secondary.TButton", secondary, background, true))),
		TButton(Txt("Success"), Style(b5.ButtonStyle("focused.success.TButton", success, background, true))),
		opts)
	Grid(TButton(Txt("Danger"), Style(b5.ButtonStyle("focused.danger.TButton", danger, background, true))),
		TButton(Txt("Warning"), Style(b5.ButtonStyle("focused.warning.TButton", warning, background, true))),
		TButton(Txt("Info"), Style(b5.ButtonStyle("focused.info.TButton", info, background, true))),
		opts)
	Grid(TButton(Txt("Light"), Style(b5.ButtonStyle("focused.light.TButton", light, background, true))),
		TButton(Txt("Dark"), Style(b5.ButtonStyle("focused.dark.TButton", dark, background, true))),
		TButton(Txt("Link"), Style(b5.ButtonStyle("focused.link.TButton", link, background, true))),
		opts)
	Grid(TExit(), Columnspan(3), opts)
	App.Configure(Background(background)).Wait()
}
