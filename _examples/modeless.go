package main

import (
	"fmt"
	"strconv"
	"strings"

	tk "modernc.org/tk9.0"
)

const APPNAME = "Modeless"

type App struct {
	percent      float64
	label        *tk.TLabelWidget
	entry        *tk.TEntryWidget
	buttonFrame  *tk.TFrameWidget
	configButton *tk.TButtonWidget
	quitButton   *tk.TButtonWidget
	configDialog *ConfigDialog
}

func main() {
	NewApp().Run()
}

func NewApp() *App {
	app := &App{percent: 50}
	tk.StyleThemeUse("clam")
	tk.WmWithdraw(tk.App)
	tk.App.WmTitle(APPNAME)
	tk.WmProtocol(tk.App, tk.WM_DELETE_WINDOW, app.onQuit)
	app.makeWidgets()
	app.makeLayout()
	app.makeBindings()
	return app
}

func (me *App) makeWidgets() {
	me.label = tk.TLabel(tk.Txt("This is modeless.go; see also modal.go"),
		tk.Relief(tk.GROOVE), tk.Background(tk.LightYellow))
	me.entry = tk.TEntry(tk.Textvariable(fmt.Sprintf("%.0f%%", me.percent)))
	me.buttonFrame = tk.TFrame()
	me.configButton = me.buttonFrame.TButton(tk.Txt("Config…"),
		tk.Command(me.onConfig))
	me.quitButton = me.buttonFrame.TButton(tk.Txt("Quit"),
		tk.Command(me.onQuit))
}

func (me *App) makeLayout() {
	opts := tk.Opts{tk.Padx(3), tk.Pady(3)}
	tk.Grid(me.label, tk.Row(0), tk.Column(0), tk.Sticky(tk.WE), opts)
	tk.Grid(me.entry, tk.Row(1), tk.Column(0), tk.Sticky(tk.WE), opts)
	tk.Grid(me.buttonFrame, tk.Row(2), tk.Column(0), tk.Columnspan(2),
		tk.Sticky(tk.WE), opts)
	tk.GridColumnConfigure(tk.App, 0, tk.Weight(1))
	tk.Grid(me.configButton, tk.Row(0), tk.Column(0), tk.Sticky(tk.W),
		opts)
	tk.Grid(me.quitButton, tk.Row(0), tk.Column(1), tk.Sticky(tk.E), opts)
	tk.GridColumnConfigure(me.buttonFrame, 1, tk.Weight(1))
}

func (me *App) makeBindings() {
	tk.Bind(tk.App, "<Escape>", tk.Command(me.onQuit))
}

func (me *App) Run() {
	tk.App.Center()
	tk.WmDeiconify(tk.App)
	tk.App.Wait()
}

func (me *App) onConfig() {
	if me.configDialog == nil {
		me.configDialog = NewConfigDialog(me.entry, &me.percent)
	}
	me.configDialog.Show()
}

func (me *App) onQuit() { tk.Destroy(tk.App) }

type ConfigDialog struct {
	notFirstUse    bool
	percent        *float64
	entry          *tk.TEntryWidget
	win            *tk.ToplevelWidget
	percentLabel   *tk.TLabelWidget
	percentSpinbox *tk.TSpinboxWidget
	closeButton    *tk.TButtonWidget
}

func NewConfigDialog(entry *tk.TEntryWidget,
	percent *float64) *ConfigDialog {
	dlg := &ConfigDialog{entry: entry, percent: percent}
	dlg.win = tk.App.Toplevel()
	dlg.win.WmTitle("Modeless — Config")
	// tk.WmAttributes(dlg.win, tk.Type("dialog")) // TODO
	tk.WmProtocol(dlg.win.Window, tk.WM_DELETE_WINDOW, dlg.onHide)
	dlg.percentLabel = dlg.win.TLabel(tk.Txt("Percent"))
	dlg.percentSpinbox = dlg.win.TSpinbox(tk.Format("%.0f%%"),
		tk.Increment(1), tk.From(0), tk.To(100),
		tk.Textvariable(fmt.Sprintf("%.0f%%", *percent)),
		tk.Command(dlg.onPercentChange))
	dlg.closeButton = dlg.win.TButton(tk.Txt("Close"),
		tk.Command(dlg.onHide))
	tk.Grid(dlg.percentLabel, tk.Row(0), tk.Column(0), tk.Sticky(tk.W))
	tk.Grid(dlg.percentSpinbox, tk.Row(0), tk.Column(1), tk.Sticky(tk.WE))
	tk.Grid(dlg.closeButton, tk.Row(1), tk.Column(0), tk.Columnspan(2))
	return dlg
}

func (me *ConfigDialog) onPercentChange() {
	text := strings.TrimSuffix(me.percentSpinbox.Textvariable(), "%")
	if percent, err := strconv.ParseFloat(text, 64); err == nil {
		*me.percent = percent
		me.entry.Configure(tk.Textvariable(fmt.Sprintf("%.0f%%", percent)))
	}
}

func (me *ConfigDialog) onHide() {
	tk.WmWithdraw(me.win.Window)
	tk.GrabRelease(me.win)
}

func (me *ConfigDialog) Show() {
	tk.WmDeiconify(me.win.Window)
	me.win.Raise(tk.App)
	if me.notFirstUse {
		me.notFirstUse = true
		me.win.Center()
	}
	tk.Focus(me.win)
	tk.Focus(me.percentSpinbox)
}
