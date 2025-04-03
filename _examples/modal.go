package main

import (
	"fmt"
	"strconv"
	"strings"

	tk "modernc.org/tk9.0"
)

const APPNAME = "Modal"

type App struct {
	percent      float64
	label        *tk.TLabelWidget
	entry        *tk.TEntryWidget
	buttonFrame  *tk.TFrameWidget
	configButton *tk.TButtonWidget
	quitButton   *tk.TButtonWidget
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
	me.label = tk.TLabel(tk.Txt("This is modal.go; see also modeless.go"),
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
	tk.Grid(me.configButton, tk.Row(0), tk.Column(0), tk.Sticky(tk.W), opts)
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
	data := ConfigDialogData{Percent: me.percent}
	dlg := NewConfigDialog(&data)
	dlg.ShowModal()
	if data.Ok {
		me.percent = data.Percent
		me.entry.Configure(tk.Textvariable(fmt.Sprintf("%.0f%%",
			me.percent)))
	}
}

func (me *App) onQuit() { tk.Destroy(tk.App) }

type ConfigDialogData struct {
	Ok      bool
	Percent float64
}

type ConfigDialog struct {
	data           *ConfigDialogData
	win            *tk.ToplevelWidget
	percentLabel   *tk.TLabelWidget
	percentSpinbox *tk.TSpinboxWidget
	buttonFrame    *tk.TFrameWidget
	okButton       *tk.TButtonWidget
	cancelButton   *tk.TButtonWidget
}

func NewConfigDialog(data *ConfigDialogData) *ConfigDialog {
	dlg := &ConfigDialog{data: data}
	dlg.win = tk.App.Toplevel()
	dlg.win.WmTitle("Modal — Config")
	// tk.WmAttributes(dlg.win, tk.Type("dialog")) // TODO
	tk.WmProtocol(dlg.win.Window, tk.WM_DELETE_WINDOW, dlg.onCancel)
	dlg.percentLabel = dlg.win.TLabel(tk.Txt("Percent"))
	dlg.percentSpinbox = dlg.win.TSpinbox(tk.Format("%.0f%%"),
		tk.Increment(1), tk.From(0), tk.To(100),
		tk.Textvariable(fmt.Sprintf("%.0f%%", data.Percent)))
	dlg.buttonFrame = dlg.win.TFrame()
	dlg.okButton = dlg.buttonFrame.TButton(tk.Txt("OK"),
		tk.Command(dlg.onOk))
	dlg.cancelButton = dlg.buttonFrame.TButton(tk.Txt("Cancel"),
		tk.Command(dlg.onCancel))
	opts := tk.Opts{tk.Padx(3), tk.Pady(3)}
	tk.Grid(dlg.percentLabel, tk.Row(0), tk.Column(0), tk.Sticky(tk.W),
		opts)
	tk.Grid(dlg.percentSpinbox, tk.Row(0), tk.Column(1), tk.Sticky(tk.WE),
		opts)
	tk.Grid(dlg.buttonFrame, tk.Row(1), tk.Column(0), tk.Columnspan(2),
		opts)
	tk.Grid(dlg.okButton, tk.Row(0), tk.Column(0), tk.Sticky(tk.E), opts)
	tk.Grid(dlg.cancelButton, tk.Row(0), tk.Column(1), tk.Sticky(tk.E),
		opts)
	tk.GridColumnConfigure(dlg.win, 1, tk.Weight(1))
	return dlg
}

func (me *ConfigDialog) onOk() {
	text := strings.TrimSuffix(me.percentSpinbox.Textvariable(), "%")
	if percent, err := strconv.ParseFloat(text, 64); err == nil {
		me.data.Percent = percent
		me.data.Ok = true
	}
	tk.Destroy(me.win)
}

func (me *ConfigDialog) onCancel() { tk.Destroy(me.win) }

func (me *ConfigDialog) ShowModal() {
	me.win.Raise(tk.App)
	tk.Focus(me.win)
	tk.Focus(me.percentSpinbox)
	tk.GrabSet(me.win)
	me.win.Center().Wait()
}
