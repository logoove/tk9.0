package main

import (
	_ "embed"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	tk "modernc.org/tk9.0"
)

const APPNAME = "Ring"

//go:embed bell.svg
var ICON_SVG string

type App struct {
	button  *tk.TButtonWidget
	when    time.Time
	message string
}

func main() {
	if len(os.Args) < 2 || os.Args[1] == "-h" || os.Args[1] == "--help" {
		usage() // doesn't return
	}
	when := getRingTime(os.Args[1])
	message := titleCase(strings.Join(os.Args[2:], " "))
	now := time.Now()
	if when.After(now) {
		time.Sleep(when.Sub(now))
	}
	NewApp(when, message).Run()
}

func usage() {
	fmt.Println("usage: ring <H[:M]|+M> [message]")
	os.Exit(1)
}

func getRingTime(arg string) time.Time {
	when := time.Now()
	rx := regexp.MustCompile(`(\d\d?)(?::(\d\d?))?|(\+\d\d?)`)
	if matches := rx.FindStringSubmatch(os.Args[1]); len(matches) == 4 {
		if matches[3] != "" {
			if mins, err := strconv.Atoi(matches[3][1:]); err != nil {
				usage()
			} else {
				when = when.Add(time.Minute * time.Duration(mins))
			}
		} else {
			if hours, err := strconv.Atoi(matches[1]); err != nil {
				usage()
			} else {
				when = time.Date(when.Year(), when.Month(), when.Day(),
					hours, 0, 0, 0, when.Location())
			}
			if matches[2] != "" {
				if mins, err := strconv.Atoi(matches[2]); err != nil {
					usage()
				} else {
					when = time.Date(when.Year(), when.Month(), when.Day(),
						when.Hour(), mins, 0, 0, when.Location())
				}
			}
		}
	}
	return when
}

func titleCase(line string) string {
	rx := regexp.MustCompile(`\b(\pL)`)
	return rx.ReplaceAllStringFunc(line, func(x string) string {
		return strings.ToUpper(x)
	})
}

func NewApp(when time.Time, message string) *App {
	app := &App{when: when, message: message}
	tk.StyleThemeUse("clam")
	tk.WmWithdraw(tk.App)
	tk.WmAttributes(tk.App, tk.Topmost(true))
	tk.App.IconPhoto(tk.NewPhoto(tk.Data(ICON_SVG)))
	tk.App.WmTitle(APPNAME)
	tk.App.Configure(tk.Background(tk.LightYellow), tk.Pady(0), tk.Padx(0))
	tk.WmProtocol(tk.App, tk.WM_DELETE_WINDOW, app.onQuit)
	for _, key := range []string{"<Escape>", "<q>", "<Return>"} {
		tk.Bind(tk.App, key, tk.Command(app.onQuit))
	}
	tk.StyleConfigure("TButton", tk.Font(tk.HELVETICA, 36, tk.BOLD),
		tk.Background(tk.LightYellow), tk.Foreground(tk.Red))
	app.button = tk.TButton(tk.Txt(app.getMesage()), tk.Command(app.onQuit),
		tk.Justify(tk.CENTER))
	app.update()
	tk.Pack(app.button, tk.Fill(tk.FILL_BOTH), tk.Expand(true),
		tk.Ipadx(15), tk.Ipady(15))
	return app
}

func (me *App) Run() {
	tk.App.SetResizable(false, false)
	tk.App.Center()
	tk.WmDeiconify(tk.App)
	tk.Focus(me.button)
	tk.App.Wait()
}

func (me *App) getMesage() string {
	now := time.Now()
	if me.message == "" {
		return now.Format("15:04")
	} else {
		return now.Format("15:04") + "\n" + me.message
	}
}

func (me *App) update() {
	me.button.Configure(tk.Txt(me.getMesage()))
	tk.TclAfter(time.Second*5, me.update)
}

func (me *App) onQuit() { tk.Destroy(tk.App) }
