package main

import (
	_ "embed"

	tk "modernc.org/tk9.0"
)

//go:embed document-new.svg
var NEW_SVG string

//go:embed document-open.svg
var OPEN_SVG string

//go:embed document-save.svg
var SAVE_SVG string

const APPNAME = "Layouts"

type App struct {
	toolbarFrame   *tk.TFrameWidget
	newToolButton  *tk.TButtonWidget
	openToolButton *tk.TButtonWidget
	saveToolButton *tk.TButtonWidget
	editFrame      *tk.TFrameWidget
	editText       *tk.TextWidget
	editVScrollbar *tk.TScrollbarWidget
	statusFrame    *tk.TFrameWidget
	statusLabel    *tk.TLabelWidget
	statusModLabel *tk.TLabelWidget
}

func main() {
	NewApp().Run()
}

func NewApp() *App {
	app := &App{}
	tk.StyleThemeUse("clam")
	tk.WmWithdraw(tk.App)
	tk.WmMinSize(tk.App, 240, 320)
	tk.App.WmTitle(APPNAME)
	tk.WmProtocol(tk.App, tk.WM_DELETE_WINDOW, app.onQuit)
	app.makeWidgets()
	app.makeLayout()
	app.makeBindings()
	return app
}

func (me *App) makeWidgets() {
	me.makeToolbar()
	me.makeCentralArea()
	me.makeStatusBar()
}

func (me *App) makeToolbar() {
	me.toolbarFrame = tk.TFrame(tk.Relief(tk.RAISED), tk.Borderwidth(5))
	me.newToolButton = me.toolbarFrame.TButton(tk.Image(tk.NewPhoto(
		tk.Data(NEW_SVG))))
	me.openToolButton = me.toolbarFrame.TButton(tk.Image(tk.NewPhoto(
		tk.Data(OPEN_SVG))))
	me.saveToolButton = me.toolbarFrame.TButton(tk.Image(tk.NewPhoto(
		tk.Data(SAVE_SVG))))
}

func (me *App) makeCentralArea() {
	me.editFrame = tk.TFrame()
	me.editText = me.editFrame.Text(tk.Font(tk.HELVETICA, 15),
		tk.Wrap(tk.WORD), tk.Yscrollcommand(func(event *tk.Event) {
			event.ScrollSet(me.editVScrollbar)
		}))
	me.editText.InsertML(MESSAGE)
	me.editVScrollbar = me.editFrame.TScrollbar(tk.Command(
		func(event *tk.Event) { event.Yview(me.editText) }))
}

func (me *App) makeStatusBar() {
	me.statusFrame = tk.TFrame(tk.Relief(tk.SUNKEN), tk.Borderwidth(5))
	me.statusLabel = me.statusFrame.TLabel(tk.Txt("Ready"))
	me.statusModLabel = me.statusFrame.TLabel(tk.Txt("MOD"),
		tk.Relief(tk.GROOVE))
}

func (me *App) makeLayout() {
	me.layoutToolbar()
	tk.Grid(me.toolbarFrame, tk.Row(0), tk.Column(0), tk.Sticky(tk.WE))
	me.layoutCentralArea()
	tk.Grid(me.editFrame, tk.Row(1), tk.Column(0), tk.Sticky(tk.NEWS))
	me.layoutStatusbar()
	tk.Grid(me.statusFrame, tk.Row(2), tk.Column(0), tk.Sticky(tk.WE))
	tk.GridColumnConfigure(tk.App, 0, tk.Weight(1))
	tk.GridRowConfigure(tk.App, 1, tk.Weight(1))
	tk.App.Configure(tk.Padx(5), tk.Pady(5))
}

func (me *App) layoutToolbar() {
	opts := tk.Opts{tk.Sticky(tk.W), tk.Padx(2.5)}
	tk.Grid(me.newToolButton, tk.Row(0), tk.Column(1), opts)
	tk.Grid(me.openToolButton, tk.Row(0), tk.Column(2), opts)
	tk.Grid(me.saveToolButton, tk.Row(0), tk.Column(3), opts)
}

func (me *App) layoutCentralArea() {
	tk.Grid(me.editText, tk.Row(0), tk.Column(0), tk.Sticky(tk.NEWS))
	tk.Grid(me.editVScrollbar, tk.Row(0), tk.Column(1), tk.Sticky(tk.NS))
	tk.GridRowConfigure(me.editFrame, 0, tk.Weight(1))
	tk.GridColumnConfigure(me.editFrame, 0, tk.Weight(1))
}

func (me *App) layoutStatusbar() {
	tk.Grid(me.statusLabel, tk.Row(0), tk.Column(0), tk.Sticky(tk.WE))
	tk.Grid(me.statusModLabel, tk.Row(0), tk.Column(2), tk.Sticky(tk.E))
	tk.GridColumnConfigure(me.statusFrame, 0, tk.Weight(1))
}

func (me *App) makeBindings() {
	tk.Bind(tk.App, "<Escape>", tk.Command(me.onQuit))
}

func (me *App) Run() {
	tk.WmGeometry(tk.App, "480x320")
	tk.App.Center()
	tk.WmDeiconify(tk.App)
	tk.App.Wait()
}

func (me *App) onQuit() { tk.Destroy(tk.App) }

const MESSAGE = `This example shows how to do nested layouts.
<br><br>The toolbar and status bar are laid out as one line grids inside a
frame.
<br><br>The central area is this Text widget and a vertical scrollbar, again
laid out as a one line grid.
<br><br>The toolbar frames and central area's frame are laid out as a grid
in the main window. So the overall grid has three grids nested inside it and
each of these nested grids has widgets nested inside.
<br><br>The key to using nested layouts is to create the nested widgets as
children of the widgets they are to be nested inside.
<br><br>Try resizing the window to see how the layout adapts.`
