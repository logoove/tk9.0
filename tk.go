// Copyright 2024 The tk9.0-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tk9_0 // import "modernc.org/tk9.0"

import (
	"bytes"
	"crypto/sha256"
	_ "embed"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/png"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
	"unsafe"

	"golang.org/x/net/html"
	"modernc.org/libtk9.0"
)

const (
	// ThemeEnvVar, if non-blank and containing a valid built-in theme
	// name, is used to set the default application theme. Calling
	// StyleThemeUse() will override the default.
	ThemeEnvVar = "TK9_THEME"

	// ScaleEnvVar, if a valid (floating point) number, sets the TkScaling
	// value at package initialization to NativeScaling*TK9_SCALE.
	ScaleEnvVar = "TK9_SCALE"

	gnuplotTimeout = time.Minute //TODO do not let the UI freeze
	goarch         = runtime.GOARCH
	goos           = runtime.GOOS
	libVersion     = libtk9_0.Version

	tcl_eval_direct = 0x40000
	tcl_ok          = 0
	tcl_error       = 1
	tcl_return      = 2
	tcl_break       = 3
	tcl_continue    = 4

	disconnectButtonTooltip = "Disconnect the application"
	exitButtonTooltip       = "Quit the application"
)

// NativeScaling is the value returned by TKScaling in package initialization before it is possibly
// changed using the [ScaleEnvVar] value.
var NativeScaling float64

// App is the main/root application window.
var App = &Window{}

var target = fmt.Sprintf("%s/%s", goos, goarch)

//TODO? ErrorMsg

// Error modes
const (
	// Errors will panic with a stack trace.
	PanicOnError = iota
	// Errors will be recorded into the Error variable using errors.Join
	CollectErrors
)

const (
	testHookWaitVar = "TK9_TEST_HOOK_WAIT"
)

type dllInfo struct {
	dll      string
	initProc string
}

// ErrorMode selects the action taken on errors.
var ErrorMode int

// Error records errors when [CollectErrors] is true.
var Error error

var (
	_ Opt    = (*MenuItem)(nil)
	_ Widget = (*Window)(nil)

	//go:embed embed/gotk.png
	icon []byte
	//go:embed embed/tklib/tooltip/tooltip.tcl
	tooltip []byte

	appDeiconified     bool
	appIconified       bool
	appWithdrawn       bool
	autocenterDisabled bool
	cleanupDirs        []string
	exitHandler        Opt
	finished           atomic.Int32
	forcedX            = -1
	forcedY            = -1
	handlers           = map[int32]*eventHandler{}
	id                 atomic.Int32
	initialized        bool
	isBuilder          = os.Getenv("MODERNC_BUILDER") != ""
	isVNC              = os.Getenv("TK9_VNC") == "1"
	testHookWait       = os.Getenv(testHookWaitVar)
	wmTitle            string

	// https://pdos.csail.mit.edu/archive/rover/RoverDoc/escape_shell_table.html
	//
	// The following characters are dissallowed or have special meanings in Tcl and
	// so are escaped:
	//
	//	&;`'"|*?~<>^()[]{}$\
	badChars = [...]bool{
		' ':  true,
		'"':  true,
		'$':  true,
		'&':  true,
		'(':  true,
		')':  true,
		'*':  true,
		';':  true,
		'<':  true,
		'>':  true,
		'?':  true,
		'[':  true,
		'\'': true,
		'\\': true,
		'\n': true,
		'\r': true,
		'\t': true,
		']':  true,
		'^':  true,
		'`':  true,
		'{':  true,
		'|':  true,
		'}':  true,
		'~':  true,
	}

	//TODO remove the associated tcl var on window destroy event both from the
	// interp and this map.
	textVariables = map[*Window]string{} // : tclName
	variables     = map[*Window]*VariableOpt{}
	windowIndex   = map[string]*Window{}
)

func commonLazyInit() {
	// Convert DOS line endings to Unix before evaluating.
	// https://gitlab.com/cznic/tk9.0/-/issues/67
	eval(string(bytes.ReplaceAll(tooltip, []byte{'\r', '\n'}, []byte{'\n'})))
}

func checkSig(dir string, sig map[string]string) (r bool) {
	if dmesgs {
		dmesg("checkSig(%q, %q)", dir, sig)
	}
	if err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		base := filepath.Base(path)
		sum := sig[base]
		if sum == "" {
			return nil
		}

		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		if g, e := fmt.Sprintf("%0x", sha256.Sum256(b)), sum; g != e {
			return fmt.Errorf("check failed: %s %s != %s", path, g, e)
		}

		delete(sig, base)
		return nil
	}); err != nil || len(sig) != 0 {
		if dmesgs {
			dmesg("checkSig(%q) failed: %v", dir, err)
		}
		return false
	}

	return true
}

// Returns a single Tcl string, no braces, except "{}" is returned for s == "".
func tclSafeString(s string) string {
	if s == "" {
		return "{}"
	}

	const badString = "&;`'\"|*?~<>^()[]{}$\\\n\r\t "
	if strings.ContainsAny(s, badString) {
		var b strings.Builder
		for _, c := range s {
			switch {
			case int(c) < len(badChars) && badChars[c]:
				fmt.Fprintf(&b, "\\x%02x", c)
			default:
				b.WriteRune(c)
			}
		}
		s = b.String()
	}
	return s
}

// Same as tclSafeStrings but does not escape <>.
func tclSafeStringBind(s string) string {
	if s == "" {
		return "{}"
	}

	const badString = "&;`'\"|*?~^()[]{}$\\\n\r\t "
	if strings.ContainsAny(s, badString) {
		var b strings.Builder
		for _, c := range s {
			switch {
			case int(c) < len(badChars) && badChars[c]:
				fmt.Fprintf(&b, "\\x%02x", c)
			default:
				b.WriteRune(c)
			}
		}
		s = b.String()
	}
	return s
}

// Returns a space separated list of safe Tcl strings.
func tclSafeList(list ...any) string {
	var a []string
	for _, v := range list {
		a = append(a, tclSafeString(fmt.Sprint(v)))
	}
	return strings.Join(a, " ")
}

// Returns a space separated list of safe Tcl strings.
func tclSafeStrings(s ...string) string {
	var a []string
	for _, v := range s {
		a = append(a, tclSafeString(v))
	}
	return strings.Join(a, " ")
}

// Returns a Tcl string that is safe inside {...}
func tclSafeInBraces(s string) string {
	const badString = "}\\\n\r\t "
	if strings.ContainsAny(s, badString) {
		var b strings.Builder
		for _, c := range s {
			switch {
			case int(c) < len(badMLChars) && badMLChars[c]:
				fmt.Fprintf(&b, "\\x%02x", c)
			default:
				b.WriteRune(c)
			}
		}
		s = b.String()
	}
	return s
}

func setDefaults() {
	windowIndex[""] = App
	windowIndex["."] = App
	if goos == "windows" {
		wmWithdraw(App)
	}
	exitHandler = Command(func() {
		Destroy(App)
		for _, v := range Themes {
			v.Deactivate(nil)
			v.Finalize(nil)
		}
	})
	evalErr("option add *tearOff 0") // https://tkdocs.com/tutorial/menus.html
	NativeScaling = TkScaling()
	if s := os.Getenv(ScaleEnvVar); s != "" {
		if k, err := strconv.ParseFloat(s, 64); err == nil {
			TkScaling(min(max(k, 0.5), 5) * NativeScaling)
		}
	}
	if nm := os.Getenv(ThemeEnvVar); nm != "" {
		StyleThemeUse(nm)
	}
	App.IconPhoto(NewPhoto(Data(icon)))
	wmTitle = filepath.Base(os.Args[0])
	wmTitle = strings.TrimSuffix(wmTitle, ".exe")
	App.WmTitle(wmTitle)
	x, y := -1, -1
	if os.Getenv("TK9_DEMO") == "1" {
		for i := 1; i < len(os.Args); i++ {
			s := os.Args[i]
			if !strings.HasPrefix(s, "+") {
				continue
			}

			a := strings.Split(s[1:], "+")
			if len(a) != 2 {
				continue
			}

			var err error
			if x, err = strconv.Atoi(a[0]); err != nil || x < 0 {
				x = -1
				break
			}

			if y, err = strconv.Atoi(a[1]); err != nil || y < 0 {
				y = -1
			}
			break
		}
	}
	App.Configure(Padx("4m"), Pady("3m"))
	if x >= 0 && y >= 0 {
		forcedX, forcedY = x, y
	}
}

// Window represents a Tk window/widget. It implements common widget methods.
//
// Window implements Opt. When a Window instance is used as an Opt, it provides
// its path name.
type Window struct {
	fpath string
}

func (w *Window) isWidget() {}

// Widget is implemented by every *Window
//
// Widget implements Opt. When a Widget instance is used as an Opt, it provides
// its path name.
type Widget interface {
	isWidget()
	optionString(*Window) string
}

// String implements fmt.Stringer.
func (w *Window) String() (r string) {
	if r = w.fpath; r == "" {
		r = "."
	}
	return r
}

func (w *Window) optionString(_ *Window) string {
	return w.String()
}

func (w *Window) split(options []Opt) (opts []Opt, tvs []textVarOpt, vs []*VariableOpt) {
	for _, v := range options {
		switch x := v.(type) {
		case textVarOpt:
			tvs = append(tvs, x)
		case *VariableOpt:
			vs = append(vs, x)
		default:
			opts = append(opts, x)
		}
	}
	return opts, tvs, vs
}

func (w *Window) newChild0(nm string) (rw *Window) {
	rw = &Window{fpath: fmt.Sprintf("%s.%s%v", w, nm, id.Add(1))}
	windowIndex[rw.fpath] = rw
	return rw
}

func (w *Window) newChild(nm string, options ...Opt) (rw *Window) {
	class := strings.Replace(nm, "ttk_", "ttk::", 1)
	nm = strings.Replace(nm, "ttk::", "t", 1)
	if c := nm[len(nm)-1]; c >= '0' && c <= '9' {
		nm += "_"
	}
	path := fmt.Sprintf("%s.%s%v", w, nm, id.Add(1))
	options, tvs, vs := w.split(options)
	rw = &Window{}
	code := fmt.Sprintf("%s %s %s", class, path, winCollect(rw, options...))
	var err error
	if rw.fpath, err = eval(code); err != nil {
		fail(fmt.Errorf("code=%s -> r=%s err=%v", code, rw.fpath, err))
	}
	if len(tvs) != 0 {
		rw.Configure(tvs[len(tvs)-1])
	}
	if len(vs) != 0 {
		rw.Configure(vs[len(vs)-1])
	}
	windowIndex[rw.fpath] = rw
	return rw
}

func evalErr(code string) (r string) {
	r, err := eval(code)
	if err != nil {
		fail(fmt.Errorf("code=%s -> r=%s err=%v", code, r, err))
	}
	return r
}

func fail(err error) {
	switch ErrorMode {
	default:
		fallthrough
	case PanicOnError:
		if dmesgs {
			dmesg("PANIC %v", err)
		}
		panic(err)
	case CollectErrors:
		Error = errors.Join(Error, err)
	}
}

func winCollect(w *Window, options ...Opt) string {
	var a []string
	for _, v := range options {
		a = append(a, v.optionString(w))
	}
	return strings.Join(a, " ")
}

func collectAny(options ...any) string {
	var a []string
	for _, v := range options {
		switch x := v.(type) {
		case Opt:
			a = append(a, x.optionString(nil))
		default:
			a = append(a, tclSafeString(fmt.Sprint(x)))
		}
	}
	return strings.Join(a, " ")
}

func collect(options ...Opt) string {
	var a []string
	for _, v := range options {
		a = append(a, v.optionString(nil))
	}
	return strings.Join(a, " ")
}

func collectOne(name string, options ...Opt) string {
	for _, v := range options {
		opt := v.optionString(nil)
		if strings.HasPrefix(opt, name) {
			return strings.TrimSpace(opt[len(name):])
		}
	}
	return ""
}

// Opts is a list of options. It implements Opt.
type Opts []Opt

func (o Opts) optionString(w *Window) string {
	return winCollect(w, []Opt(o)...)
}

// Opt represents an optional argument.
type Opt interface {
	optionString(w *Window) string
}

type rawOption string

func (s rawOption) optionString(w *Window) string {
	return string(s)
}

type stringOption string

func (s stringOption) optionString(w *Window) string {
	return tclSafeString(string(s))
}

// EventHandler is the type used by call backs.
type EventHandler func(*Event)

// Event communicates information with an event handler. All handlers can use the Err field. Simple
// handlers, like in
//
//	Button(..., Command(func(e *Event) {...}))
//
// can use the 'W' field, if applicable.  All other fields are valid only in
// handlers bound using [Bind].
type Event struct {
	// Event handlers should set Err on failure.
	Err error
	// Result can be optionally set by the handler. The field is returned by
	// certain methods, for example [TCheckbuttonWidget.Invoke].
	Result string
	// Event source, if any. This field is set when the event handler was
	// created.
	W *Window

	returnCode int // One of tcl_ok .. tcl_continue

	// The window to which the event was reported (the window field from
	// the event). Valid for all event types.  This field is set when the
	// event is handled.
	EventWindow *Window
	// The keysym corresponding to the event, substituted as a textual string.
	// Valid only for Key and KeyRelease events.
	Keysym string
	// The number of the last client request processed by the server (the serial
	// field from the event). Valid for all event types.
	Serial int64
	// The width field from the event. Indicates the new or requested width of the
	// window. Valid only for Configure, ConfigureRequest, Create, ResizeRequest,
	// and Expose events.
	Width string
	// The height field from the event. Valid for the Configure, ConfigureRequest,
	// Create, ResizeRequest, and Expose events. Indicates the new or requested
	// height of the window.
	Height string
	// The x and y fields from the event. For Button, ButtonRelease, Motion, Key,
	// KeyRelease, and MouseWheel events, X and Y indicate the position of the
	// mouse pointer relative to the receiving window. For key events on the
	// Macintosh these are the coordinates of the mouse at the moment when an X11
	// KeyEvent is sent to Tk, which could be slightly later than the time of the
	// physical press or release. For Enter and Leave events, the position where
	// the mouse pointer crossed the window, relative to the receiving window. For
	// Configure and Create requests, the x and y coordinates of the window
	// relative to its parent window.
	X, Y int
	// The x_root and y_root fields from the event. If a virtual-root window
	// manager is being used then the substituted values are the corresponding
	// x-coordinate and y-coordinate in the virtual root. Valid only for Button,
	// ButtonRelease, Enter, Key, KeyRelease, Leave and Motion events. Same meaning
	// as X and Y, except relative to the (virtual) root window.
	XRoot, YRoot int
	// The delta value of a MouseWheel event. The delta value represents
	// the rotation units the mouse wheel has been moved. The sign of the
	// value represents the direction the mouse wheel was scrolled.
	Delta int
	// The state field from the event. For KeyPress, KeyRelease, ButtonPress,
	// ButtonRelease, Enter, Leave, and Motion events, it is a bit field.
	// Visibility events are not currently supported, and the value will be 0.
	State Modifier

	args []string
}

// Called from eventDispatcher. Arg1 is handler id, optionally followed by a
// list of Bind substitution values.
func newEvent(arg1 string) (id int, e *Event, err error) {
	e = &Event{returnCode: tcl_ok}
	a := strings.Fields(arg1)
	if len(a) == 0 {
		return -1, e, fmt.Errorf("internal error: missing handler ID")
	}

	if id, err = strconv.Atoi(a[0]); err != nil {
		return id, e, fmt.Errorf("newEvent: parsing event ID %q: %v", a[0], err)
	}

	for i, v := range a[1:] {
		switch i {
		case 0: // %#
			if e.Serial, err = strconv.ParseInt(v, 10, 64); err != nil {
				return id, e, fmt.Errorf("newEvent: parsing event serial %q: %v", v, err)
			}
		case 1: // %W
			e.EventWindow = windowIndex[v]
		case 2: // %K
			e.Keysym = v
		case 3: // %w
			e.Width = v
		case 4: // %h
			e.Height = v
		case 5: // %x
			e.X = atoi(v)
		case 6: // %y
			e.Y = atoi(v)
		case 7: // %X
			e.XRoot = atoi(v)
		case 8: // %Y
			e.YRoot = atoi(v)
		case 9: // %D
			e.Delta = atoi(v)
		case 10: // %s
			e.State = Modifier(atoi(v))
		}
	}
	return id, e, nil
}

func atoi(s string) int {
	if n, err := strconv.Atoi(s); err == nil {
		return n
	}

	return 0
}

// SetReturnCodeOK sets return code of 'e' to TCL_OK.
func (e *Event) SetReturnCodeOK() {
	e.returnCode = tcl_ok
}

// SetReturnCodeError sets return code of 'e' to TCL_ERROR.
func (e *Event) SetReturnCodeError() {
	e.returnCode = tcl_error
}

// SetReturnCodeReturn sets return code of 'e' to TCL_RETURN.
func (e *Event) SetReturnCodeReturn() {
	e.returnCode = tcl_return
}

// SetReturnCodeBreak sets return code of 'e' to TCL_BREAK.
func (e *Event) SetReturnCodeBreak() {
	e.returnCode = tcl_break
}

// SetReturnCodeContinue sets return code of 'e' to TCL_CONTINUE.
func (e *Event) SetReturnCodeContinue() {
	e.returnCode = tcl_continue
}

// ScrollSet communicates events to scrollbars. Example:
//
//	var scroll *TScrollbarWidget
//	// tcl: text .text -yscrollcommand ".scroll set"
//	t := Text(..., Yscrollcommand(func(e *Event) { e.ScrollSet(scroll) }))
func (e *Event) ScrollSet(w Widget) {
	if len(e.args) > 1 {
		evalErr(fmt.Sprintf("%s set %s %s", w, e.args[0], e.args[1]))
	}
}

// Xview communicates events to views. Example:
//
//	var scroll *TScrollbarWidget
//	t := Text(...)
//	// tcl: ttk::scrollbar .scroll -command ".text xview"
//	scroll = TScrollbar(Command(func(e *Event) { e.Xview(t)}))
func (e *Event) Xview(w Widget) {
	if len(e.args) > 1 {
		evalErr(fmt.Sprintf("%s xview %s", w, strings.Join(e.args, " ")))
	}
}

// Yview communicates events to views. Example:
//
//	var scroll *TScrollbarWidget
//	t := Text(...)
//	// tcl: ttk::scrollbar .scroll -command ".text yview"
//	scroll = TScrollbar(Command(func(e *Event) { e.Yview(t)}))
func (e *Event) Yview(w Widget) {
	if len(e.args) > 1 {
		evalErr(fmt.Sprintf("%s yview %s", w, strings.Join(e.args, " ")))
	}
}

type eventHandler struct {
	callback func(*Event)
	id       int32
	tcl      string
	w        *Window

	lateBind bool
}

func newEventHandler(option string, handler any) (r *eventHandler) {
	var callback func(*Event)
	switch x := handler.(type) {
	case EventHandler:
		callback = x
	case func(*Event):
		callback = x
	case func():
		callback = func(*Event) { x() }
	default:
		fail(fmt.Errorf("registering event handler: unsupported handler type: %T", handler))
		return nil
	}

	r = &eventHandler{
		callback: callback,
		id:       id.Add(1),
		tcl:      option,
	}
	handlers[r.id] = r
	return r
}

func (e *eventHandler) optionString(w *Window) string {
	if e == nil {
		return ""
	}

	e.w = w
	switch {
	case e.lateBind:
		return fmt.Sprintf("%s {eventDispatcher {%v %%# %%W %%K %%w %%h %%x %%y %%X %%Y %%D %%s}}", e.tcl, e.id)
	default:
		return fmt.Sprintf("%s {eventDispatcher %v}", e.tcl, e.id)
	}
}

func hexDigit(b byte) byte {
	if b <= 9 {
		return '0' + b
	}

	return 'a' + b - 10
}

func tclBinaryBytes(s []byte) string {
	var b strings.Builder
	for _, v := range s {
		b.WriteByte('\\')
		b.WriteByte('x')
		b.WriteByte(hexDigit(v >> 4))
		b.WriteByte(hexDigit(v & 15))
	}
	return b.String()
}

var xpmSig = []byte("/* XPM */") // https://en.wikipedia.org/wiki/X_PixMap

func optionString(v any) (r string) {
	switch x := v.(type) {
	case time.Duration:
		return fmt.Sprint(int64((x + time.Millisecond/2) / time.Millisecond))
	case []byte:
		switch {
		case bytes.HasPrefix(x, xpmSig):
			return tclSafeString(fmt.Sprintf("%s", v))
		default:
			return tclBinaryBytes(x)
		}
	case image.Image:
		var buf bytes.Buffer
		err := png.Encode(&buf, x)
		if err != nil {
			fail(err)
			return ""
		}

		return base64.StdEncoding.EncodeToString(buf.Bytes())
	case []FileType:
		var a []string
		for _, v := range x {
			a = append(a, fmt.Sprintf("{%s {%s} %s}", tclSafeString(v.TypeName), tclSafeStrings(v.Extensions...), v.MacType))
		}
		return fmt.Sprintf("{%s}", strings.Join(a, " "))
	default:
		return tclSafeString(fmt.Sprint(v))
	}
}

// bind — Arrange for X events to invoke functions
//
// # Description
//
// Bind tag options...
//
// The bind command associates commands with X events. If all three
// arguments are specified, bind will arrange for script (a Tcl script called
// the “binding script”) to be evaluated whenever the event(s) given by
// sequence occur in the window(s) identified by tag. If script is prefixed
// with a “+”, then it is appended to any existing binding for sequence;
// otherwise script replaces any existing binding. If script is an empty string
// then the current binding for sequence is destroyed, leaving sequence
// unbound. In all of the cases where a script argument is provided, bind
// returns an empty string.
//
// If sequence is specified without a script, then the script currently bound
// to sequence is returned, or an empty string is returned if there is no
// binding for sequence. If neither sequence nor script is specified, then the
// return value is a list whose elements are all the sequences for which there
// exist bindings for tag.
//
// The tag argument determines which window(s) the binding applies to. If tag
// begins with a dot, as in .a.b.c, then it must be the path name for a window;
// otherwise it may be an arbitrary string. Each window has an associated list
// of tags, and a binding applies to a particular window if its tag is among
// those specified for the window. Although the [Bindtags] command may be used to
// assign an arbitrary set of binding tags to a window, the default binding
// tags provide the following behavior:
//
//   - If a tag is the name of an internal window the binding applies to that
//     window.
//   - If the tag is the name of a class of widgets, such as Button, the
//     binding applies to all widgets in that class.
//   - If the tag is the name of a toplevel window the binding applies to the
//     toplevel window and all its internal windows.
//   - If tag has the value all, the binding applies to all windows in the
//     application.
//
// Example usage in _examples/events.go.
//
// Additional information might be available at the [Tcl/Tk bind] page.
//
// [Tcl/Tk bind]: https://www.tcl.tk/man/tcl9.0/TkCmd/bind.html
func Bind(options ...any) {
	a := []string{"bind"}
	var w *Window
	for _, v := range options {
		switch x := v.(type) {
		case *Window:
			if w == nil {
				w = x
			}
			a = append(a, x.String())
		case *eventHandler:
			x.tcl = ""
			x.lateBind = true
			a = append(a, x.optionString(w))
		default:
			a = append(a, tclSafeStringBind(fmt.Sprint(x)))
		}
	}
	evalErr(strings.Join(a, " "))
}

// bindtags — Determine which bindings apply to a window, and order of evaluation
//
// # Description
//
// When a binding is created with the bind command, it is associated either
// with a particular window such as .a.b.c, a class name such as Button, the
// keyword all, or any other string. All of these forms are called binding
// tags. Each window contains a list of binding tags that determine how events
// are processed for the window. When an event occurs in a window, it is
// applied to each of the window's tags in order: for each tag, the most
// specific binding that matches the given tag and event is executed. See the
// bind command for more information on the matching process.
//
// By default, each window has four binding tags consisting of the name of the
// window, the window's class name, the name of the window's nearest toplevel
// ancestor, and all, in that order. Toplevel windows have only three tags by
// default, since the toplevel name is the same as that of the window. The
// bindtags command allows the binding tags for a window to be read and
// modified.
//
// If bindtags is invoked with only one argument, then the current set of
// binding tags for window is returned as a list. If the tagList argument is
// specified to bindtags, then it must be a proper list; the tags for window
// are changed to the elements of the list. The elements of tagList may be
// arbitrary strings; however, any tag starting with a dot is treated as the
// name of a window; if no window by that name exists at the time an event is
// processed, then the tag is ignored for that event. The order of the elements
// in tagList determines the order in which binding scripts are executed in
// response to events. For example, the command
//
//	Bindtags(myBtn, "all", App, "Button", myBtn)
//
// reverses the order in which binding scripts will be evaluated for a button
// myBtn so that all bindings are invoked first, following by bindings for
// buttons (App), followed by class bindings, followed by bindings for
// myBtn. If tagList is an empty list then the binding tags for window are
// returned to the default state described above.
//
// The bindtags command may be used to introduce arbitrary additional binding
// tags for a window, or to remove standard tags. For example, the command
//
//	Bindtags(myBtn, myBtn, "TrickyButton", App, "all")
//
// replaces the Button tag for myBtn with TrickyButton. This means that the
// default widget bindings for buttons, which are associated with the Button
// tag, will no longer apply to myBtn, but any bindings associated with
// TrickyButton (perhaps some new button behavior) will apply.
//
// Additional information might be available at the [Tcl/Tk bindtags] page.
//
// [Tcl/Tk bindtags]: https://www.tcl.tk/man/tcl9.0/TkCmd/bindtags.html
func Bindtags(w *Window, tagList ...any) (r []string) {
	switch len(tagList) {
	case 0:
		return parseList(evalErr(fmt.Sprintf("bindtags %s", w)))
	default:
		return parseList(evalErr(fmt.Sprintf("bindtags %s {%s}", w, tclSafeList(flat(tagList...)...))))
	}
}

// Img represents a Tk image.
type Img struct {
	name string
}

// image — Create and manipulate images
//
// # Description
//
// Deletes 'm'. If there are instances of the image displayed in widgets, the
// image will not actually be deleted until all of the instances are released.
// However, the association between the instances and the image manager will be
// dropped. Existing instances will retain their sizes but redisplay as empty
// areas. If a deleted image is recreated with another call to image create,
// the existing instances will use the new image.
func (m *Img) Delete() {
	evalErr(fmt.Sprintf("image delete %s", m))
}

// String implements fmt.Stringer.
func (m *Img) String() string {
	return m.optionString(nil)
}

func (m *Img) optionString(_ *Window) string {
	if m != nil {
		return m.name
	}

	return "img0" // does not exist
}

// Bitmap — Images that display two colors
//
// # Description
//
// A bitmap is an image whose pixels can display either of two colors or be
// transparent. A bitmap image is defined by four things: a background color, a
// foreground color, and two bitmaps, called the source and the mask. Each of
// the bitmaps specifies 0/1 values for a rectangular array of pixels, and the
// two bitmaps must have the same dimensions. For pixels where the mask is
// zero, the image displays nothing, producing a transparent effect. For other
// pixels, the image displays the foreground color if the source data is one
// and the background color if the source data is zero.
//
// Additional information might be available at the [Tcl/Tk bitmap] page.
//
//   - [Background] color
//
// Specifies a background color for the image in any of the standard ways
// accepted by Tk. If this option is set to an empty string then the background
// pixels will be transparent. This effect is achieved by using the source
// bitmap as the mask bitmap, ignoring any -maskdata or -maskfile options.
//
//   - [Data] string
//
// Specifies the contents of the source bitmap as a string. The string must
// adhere to X11 bitmap format (e.g., as generated by the bitmap program). If
// both the -data and -file options are specified, the -data option takes
// precedence.
//
//   - [File] name
//
// name gives the name of a file whose contents define the source bitmap. The
// file must adhere to X11 bitmap format (e.g., as generated by the bitmap
// program).
//
//   - [Foreground] color
//
// Specifies a foreground color for the image in any of the standard ways
// accepted by Tk.
//
//   - [Maskdata] string
//
// Specifies the contents of the mask as a string. The string must adhere to
// X11 bitmap format (e.g., as generated by the bitmap program). If both the
// -maskdata and -maskfile options are specified, the -maskdata option takes
// precedence.
//
//   - [Maskfile] name
//
// name gives the name of a file whose contents define the mask. The file must
// adhere to X11 bitmap format (e.g., as generated by the bitmap program).
//
// [Tcl/Tk bitmap]: https://www.tcl.tk/man/tcl9.0/TkCmd/bitmap.html
func NewBitmap(options ...Opt) *Img {
	nm := fmt.Sprintf("bmp%v", id.Add(1))
	code := fmt.Sprintf("image create bitmap %s %s", nm, collect(options...))
	r, err := eval(code)
	if err != nil {
		fail(fmt.Errorf("code=%s -> r=%s err=%v", code, r, err))
		return nil
	}

	return &Img{name: nm}
}

// Photo — Full-color images
//
// A photo is an image whose pixels can display any color with a varying degree
// of transparency (the alpha channel). A photo image is stored internally in
// full color (32 bits per pixel), and is displayed using dithering if
// necessary. Image data for a photo image can be obtained from a file or a
// string, or it can be supplied from C code through a procedural interface. A
// photo image is (semi)transparent if the image data it was obtained from had
// transparency information. In regions where no image data has been supplied,
// it is fully transparent. Transparency may also be modified with the
// transparency set subcommand.
//
// # Supported image formats
//
//   - BMP
//   - GIF
//   - ICO
//   - JPEG
//   - PCX
//   - PNG
//   - PPM
//   - SVG
//   - TGA
//   - TIFF
//   - XBM
//   - XMP
//
// # Photo Options
//
//   - [Data] string
//
// Specifies the contents of the image as a string. The string should contain
// binary data or, for some formats, base64-encoded data (this is currently
// guaranteed to be supported for PNG and GIF images). The format of the string
// must be one of those for which there is an image file format handler that
// will accept string data. If both the -data and -file options are specified,
// the -file option takes precedence.
//
//   - [Format] format-name
//
// Specifies the name of the file format for the data specified with the -data
// or -file option.
//
//   - [File] name
//
// name gives the name of a file that is to be read to supply data for the
// photo image. The file format must be one of those for which there is an
// image file format handler that can read data.
//
//   - [Gamma] value
//
// Specifies that the colors allocated for displaying this image in a window
// should be corrected for a non-linear display with the specified gamma
// exponent value. (The intensity produced by most CRT displays is a power
// function of the input value, to a good approximation; gamma is the exponent
// and is typically around 2). The value specified must be greater than zero.
// The default value is one (no correction). In general, values greater than
// one will make the image lighter, and values less than one will make it
// darker.
//
//   - [Height] number
//
// Specifies the height of the image, in pixels. This option is useful
// primarily in situations where the user wishes to build up the contents of
// the image piece by piece. A value of zero (the default) allows the image to
// expand or shrink vertically to fit the data stored in it.
//
//   - [Palette] palette-spec
//
// Specifies the resolution of the color cube to be allocated for displaying
// this image, and thus the number of colors used from the colormaps of the
// windows where it is displayed. The palette-spec string may be either a
// single decimal number, specifying the number of shades of gray to use, or
// three decimal numbers separated by slashes (/), specifying the number of
// shades of red, green and blue to use, respectively. If the first form (a
// single number) is used, the image will be displayed in monochrome (i.e.,
// grayscale).
//
//   - [Width] number
//
// Specifies the width of the image, in pixels. This option is useful primarily
// in situations where the user wishes to build up the contents of the image
// piece by piece. A value of zero (the default) allows the image to expand or
// shrink horizontally to fit the data stored in it.
//
// Additional information might be available at the [Tcl/Tk photo] page.
//
// # Format("bmp") additional options
//
// In addition the value of option [Format] is treated as a list and may contain
// any of the special options listed below.
//
//   - -verbose bool
//
//     This option is supported for reading and writing. Available since
//     version 2.0.
//
//     If set to true, additional information about the read or written image
//     is printed to stdout. Default is false.
//
//   - -resolution xres ?yres?
//
//     This option is supported for writing only. Available since version 2.0.
//     An incompatible version of this option was introduced in version 1.4.1.
//
//     Set the resolution values of the written image file. If yres is not
//     specified, it is set to the value of xres.
//
//     If option is not specified, the DPI and aspect values of the metadata
//     dictionary are written. If no metadata values are available, no
//     resolution values are written.
//
//   - -xresolution xres
//
//     This option is supported for writing only. Available since version 2.0.
//
//     Set the horizontal resolution value of the written image file.
//
//   - -yresolution yres
//
//     This option is supported for writing only. Available since version 2.0.
//
//     Set the vertical resolution value of the written image file.
//
// Additional information might be available at the [tkImg-bmp] page.
//
// # Format("ico") additional options
//
// In addition the value of option [Format] is treated as a list and may contain
// any of the special options listed below.
//
//   - -verbose bool
//
//     This option is supported for reading and writing. Available since version
//     1.3.
//
//     If set to true, additional information about the read or written image is
//     printed to stdout. Default is false.
//
//   - -index integer
//
//     This option is supported for reading only. Available since version 1.3.
//
//     Read the page at specified index. The first page is at index 0. Default is
//     0.
//
// Additional information might be available at the [tkImg-ico] page.
//
// # Format("jpeg") additional options
//
// In addition the value of option [Format] is treated as a list and may contain
// any of the special options listed below.
//
//   - -verbose bool
//
//     This option is supported for reading and writing. Available since version
//     2.0.
//
//     If set to true, additional information about the read or written image is
//     printed to stdout. Default is false.
//
//   - -fast
//
//     This option is supported for reading only. Available since version
//     1.2.4.
//
//     If specified, it activates a processing mode which is fast, but
//     provides only low-quality information.
//
//   - -grayscale
//
//     This option is supported for reading and writing. Available since
//     version 1.2.4.
//
//     Usage of this option forces incoming images to grayscale and written
//     images will be monochrome.
//
//   - -optimize
//
//     This option is supported for writing only. Available since version
//     1.2.4.
//
//     If specified, it causes the writer to optimize the Huffman table used
//     to encode the JPEG coefficients.
//
//   - -progressive
//
//     This option is supported for writing only. Available since version
//     1.2.4.
//
//     If specified, it causes the creation of a progressive JPEG file.
//
//   - -quality n
//
//     This option is supported for writing only. Available since version
//     1.2.4.
//
//     It specifies the compression level as a quality percentage. The higher the
//     quality, the less the compression. The nominal range for n is 0...100.
//     Useful values are in the range 5...95. The default value is 75.
//
//   - -smooth n
//
//     This option is supported for writing only. Available since version
//     1.2.4.
//
//     When used the writer will smooth the image before performing the
//     compression. Values in the 10...30 are usually enough. The default is
//     0, i.e no smoothing.
//
//   - -resolution xres ?yres?
//
//     This option is supported for writing only. Available since version 2.0.
//
//     Set the resolution values of the written image file. If yres is not
//     specified, it is set to the value of xres.
//
//     If option is not specified, the DPI and aspect values of the metadata
//     dictionary are written. If no metadata values are available, no
//     resolution values are written.
//
//   - -xresolution xres
//
//     This option is supported for writing only. Available since version 2.0.
//
//     Set the horizontal resolution value of the written image file.
//
//   - -yresolution yres
//
//     This option is supported for writing only. Available since version 2.0.
//
//     Set the vertical resolution value of the written image file.
//
// Additional information might be available at the [tkImg-jpeg] page.
//
// # Format("pcx") additional options
//
// In addition the value of option [Format] is treated as a list and may contain
// any of the special options listed below.
//
//   - -verbose bool
//
//     This option is supported for reading and writing. Available since
//     version 1.3.
//
//     If set to true, additional information about the read or written image
//     is printed to stdout. Default is false.
//
//   - -compression string
//
//     This option is supported for writing only. Available since version 1.3.
//
//     Set the compression mode to either none or rle. Default is rle.
//
//   - -resolution xres ?yres?
//
//     This option is supported for writing only. Available since version 2.0.
//
//     Set the resolution values of the written image file. If yres is not
//     specified, it is set to the value of xres.
//
//     If option is not specified, the DPI and aspect values of the metadata
//     dictionary are written. If no metadata values are available, no
//     resolution values are written.
//
//   - -xresolution xres
//
//     This option is supported for writing only. Available since version 2.0.
//
//     Set the horizontal resolution value of the written image file.
//
//   - -yresolution yres
//
//     This option is supported for writing only. Available since version 2.0.
//
//     Set the vertical resolution value of the written image file.
//
// Additional information might be available at the [tkImg-pcx] page.
//
// # Format("png") additional options
//
// In addition the value of option [Format] is treated as a list and may contain
// any of the special options listed below.
//
//   - -verbose bool
//
//     This option is supported for reading and writing. Available since
//     version 1.4.6.
//
//     If set to true, additional information about the read or written image
//     is printed to stdout. Default is false.
//
//   - -alpha double
//
//     This option is supported for reading only. Available since version
//     1.4.2.
//
//     An additional alpha filtering value for the overall image, which allows
//     the background on which the image is displayed to show through. This
//     usually also has the effect of desaturating the image. The alpha value
//     must be between 0.0 and 1.0. Specifying an alpha value, overrides the
//     setting of the -withalpha flag, i.e. reading a file which has no alpha
//     channel (Grayscale, RGB) will add an alpha channel to the image
//     independent of the -withalpha flag setting.
//
//   - -gamma double
//
//     This option is supported for reading only. Available since version
//     1.4.6.
//
//     Use the specified gamma value when reading an image. This option
//     overwrites gamma values specified in the file. If this option is not
//     specified and no gamma value is in the file, a default value of 1.0 is
//     used.
//
//   - -withalpha bool
//
//     This option is supported for reading and writing. Available since
//     version 1.4.1.
//
//     If set to false, an alpha channel is ignored during reading or writing.
//     Default is true.
//
//     Note: This option was named -matte in previous versions and is still
//     recognized.
//
//   - -resolution xres ?yres?
//
//     This option is supported for writing only. Available since version 2.0.
//
//     Set the resolution values of the written image file. If yres is not
//     specified, it is set to the value of xres.
//
//     If option is not specified, the DPI and aspect values of the metadata
//     dictionary are written. If no metadata values are available, no
//     resolution values are written.
//
//   - -xresolution xres
//
//     This option is supported for writing only. Available since version 2.0.
//
//     Set the horizontal resolution value of the written image file.
//
//   - -yresolution yres
//
//     This option is supported for writing only. Available since version 2.0.
//
//     Set the vertical resolution value of the written image file.
//
//   - -tag key value
//
//     This option is supported for writing only. Available since version 2.0.
//
//     Each key-value pair will be written as a named text chunk where the key
//     provides the name of the chunk and the value its contents. Currently the
//     maximum number of -tag specifications are 10.
//
// Additional information might be available at the [tkImg-png] page.
//
// # Format("ppm") additional options
//
// In addition the value of option [Format] is treated as a list and may contain
// any of the special options listed below.
//
//   - -verbose bool
//
//     This option is supported for reading and writing. Available since
//     version 1.4.0.
//
//     If set to true, additional information about the read or written image
//     is printed to stdout. Default is false.
//
//   - -scanorder string
//
//     This option is supported for reading only. Available since version
//     1.4.0.
//
//     Specify the scanline order of the input image. Possible values: TopDown
//     or BottomUp. Default is TopDown.
//
//   - -min double
//
//     This option is supported for reading only. Available since version
//     1.4.0.
//
//     Specify the minimum pixel value to be used for mapping 16-bit input data
//     to 8-bit image values. If not specified or negative, the minimum value
//     found in the image data.
//
//   - -max float
//
//     This option is supported for reading only. Available since version
//     1.4.0.
//
//     Specify the maximum pixel value to be used for mapping 16-bit input data
//     to 8-bit image values. If not specified or negative, the maximum value
//     found in the image data.
//
//   - -gamma double
//
//     This option is supported for reading only. Available since version
//     1.4.0.
//
//     Specify a gamma correction to be applied when mapping 16-bit input data
//     to 8-bit image values. Default is 1.0.
//
//   - -ascii bool
//
//     This option is supported for writing only. Available since version
//     1.4.0.
//
//     If set to true, the file is written in PPM 8-bit ASCII format (P3).
//     Default is false, i.e. write in PPM 8-bit binary format (P6).
//
// Additional information might be available at the [tkImg-ppm] page.
//
// # Format("tga") additional options
//
// In addition the value of option [Format] is treated as a list and may contain
// any of the special options listed below.
//
//   - -verbose bool
//
//     This option is supported for reading and writing. Available since
//     version 1.3.
//
//     If set to true, additional information about the read or written image
//     is printed to stdout. Default is false.
//
//   - -withalpha bool
//
//     This option is supported for reading and writing. Available since
//     version 1.3.
//
//     If set to false, an alpha channel is ignored during reading or writing.
//     Default is true.
//
//     Note: This option was named -matte in previous versions and is still
//     recognized.
//
//   - -compression string
//
//     This option is supported for writing only. Available since version 1.3.
//
//     Set the compression mode to either none or rle. Default is rle.
//
// Additional information might be available at the [tkImg-tga] page.
//
// # Format("tiff") additional options
//
// In addition the value of option [Format] is treated as a list and may contain
// any of the special options listed below.
//
//   - -verbose bool
//
//     This option is supported for reading and writing. Available since
//     version 2.0.
//
//     If set to true, additional information about the read or written image
//     is printed to stdout. Default is false.
//
//   - -index integer
//
//     This option is supported for reading only. Available since version
//     1.4.0.
//
//     Read the page at specified index. The first page is at index 0. Default
//     is 0.
//
//   - -compression string
//
//     This option is supported for writing only. Available since version
//     1.2.4.
//
//     Set the compression mode to either none, jpeg, packbits, or deflate.
//     Default is none.
//
//   - -byteorder string
//
//     This option is supported for writing only. Available since version
//     1.2.4.
//
//     Set the byteorder to either none, bigendian, littleendian, network or
//     smallendian. Default is none.
//
//     The values bigendian and network are aliases of each other, as are
//     littleendian and smallendian.
//
//   - -resolution xres ?yres?
//
//     This option is supported for writing only. Available since version 2.0.
//
//     Set the resolution values of the written image file. If yres is not
//     specified, it is set to the value of xres.
//
//     If option is not specified, the DPI and aspect values of the metadata
//     dictionary are written. If no metadata values are available, no
//     resolution values are written.
//
//   - -xresolution xres
//
//     This option is supported for writing only. Available since version 2.0.
//
//     Set the horizontal resolution value of the written image file.
//
//   - -yresolution yres
//
//     This option is supported for writing only. Available since version 2.0.
//
//     Set the vertical resolution value of the written image file.
//
// Additional information might be available at the [tkImg-tiff] page.
//
// # Format("xbm") additional options
//
// In addition the value of option [Format] is treated as a list and may contain
// any of the special options listed below.
//
//   - -verbose bool
//
//     This option is supported for reading and writing. Available since
//     version 2.0.
//
//     If set to true, additional information about the read or written image
//     is printed to stdout. Default is false.
//
//   - -foreground string
//
//     This option is supported for reading only. Available since version
//     1.4.15.
//
//     Set the foreground color of the bitmap. Default value is black. The
//     color string may be given in a format as accepted by Tk_GetColor.
//
//   - -background string
//
//     This option is supported for reading only. Available since version
//     1.4.15.
//
//     Set the background color of the bitmap. Default value is transparent.
//     The color string may be given in a format as accepted by Tk_GetColor.
//
// Additional information might be available at the [tkImg-xbm] page.
//
// # Format("xpm") additional options
//
// In addition the value of option [Format] is treated as a list and may contain
// any of the special options listed below.
//
//   - -verbose bool
//
//     This option is supported for reading and writing. Available since
//     version 2.0.
//
//     If set to true, additional information about the read or written image
//     is printed to stdout. Default is false.
//
// Additional information might be available at the [tkImg-xpm] page.
//
// [Tcl/Tk photo]: https://www.tcl.tk/man/tcl9.0/TkCmd/photo.html
// [tkImg-bmp]: https://tkimg.sourceforge.net/RefMan/files/img-bmp.html
// [tkImg-ico]: https://tkimg.sourceforge.net/RefMan/files/img-ico.html
// [tkImg-jpeg]: https://tkimg.sourceforge.net/RefMan/files/img-jpeg.html
// [tkImg-pcx]: https://tkimg.sourceforge.net/RefMan/files/img-pcx.html
// [tkImg-png]: https://tkimg.sourceforge.net/RefMan/files/img-png.html
// [tkImg-ppm]: https://tkimg.sourceforge.net/RefMan/files/img-ppm.html
// [tkImg-tga]: https://tkimg.sourceforge.net/RefMan/files/img-tga.html
// [tkImg-tiff]: https://tkimg.sourceforge.net/RefMan/files/img-tiff.html
// [tkImg-xbm]: https://tkimg.sourceforge.net/RefMan/files/img-xbm.html
// [tkImg-xpm]: https://tkimg.sourceforge.net/RefMan/files/img-xmp.html
func NewPhoto(options ...Opt) *Img {
	nm := fmt.Sprintf("img%v", id.Add(1))
	code := fmt.Sprintf("image create photo %s %s", nm, collect(options...))
	r, err := eval(code)
	if err != nil {
		fail(fmt.Errorf("code=%s -> r=%s err=%v", code, r, err))
		return nil
	}

	return &Img{name: nm}
}

// Width — Get the configured option value.
//
// Additional information might be available at the [Tcl/Tk photo] page.
//
// [Tcl/Tk photo]: https://www.tcl.tk/man/tcl9.0/TkCmd/photo.html
func (m *Img) Width() string {
	return evalErr(fmt.Sprintf(`image width %s`, m))
}

// Height — Get the configured option value.
//
// Additional information might be available at the [Tcl/Tk photo] page.
//
// [Tcl/Tk photo]: https://www.tcl.tk/man/tcl9.0/TkCmd/photo.html
func (m *Img) Height() string {
	return evalErr(fmt.Sprintf(`image height %s`, m))
}

// // Returns photo data.
// //
// // Additional information might be available at the [Tcl/Tk photo] page.
// //
// // [Tcl/Tk photo]: https://www.tcl.tk/man/tcl9.0/TkCmd/photo.html
// func (m *Img) Data(options ...Opt) []byte {
// 	s := evalErr(fmt.Sprintf("%s data %s", collect(options...)))
// 	panic(todo("%q", s))
// }

// photo — Full-color images
//
// Copies a region from the image called sourceImage (which must be a photo
// image) to the image called imageName, possibly with pixel zooming and/or
// subsampling. If no options are specified, this command copies the whole of
// sourceImage into imageName, starting at coordinates (0,0) in imageName.
//
// The following options may be specified:
//
//   - [From] x1 y1 x2 y2
//
// Specifies a rectangular sub-region of the source image to be copied. (x1,y1)
// and (x2,y2) specify diagonally opposite corners of the rectangle. If x2 and
// y2 are not specified, the default value is the bottom-right corner of the
// source image. The pixels copied will include the left and top edges of the
// specified rectangle but not the bottom or right edges. If the -from option
// is not given, the default is the whole source image.
//
//   - [To] x1 y1 x2 y2
//
// Specifies a rectangular sub-region of the destination image to be affected.
// (x1,y1) and (x2,y2) specify diagonally opposite corners of the rectangle. If
// x2 and y2 are not specified, the default value is (x1,y1) plus the size of
// the source region (after subsampling and zooming, if specified). If x2 and
// y2 are specified, the source region will be replicated if necessary to fill
// the destination region in a tiled fashion.
//
// The function returns 'm'.
//
// Additional information might be available at the [Tcl/Tk photo] page.
//
// [Tcl/Tk photo]: https://www.tcl.tk/man/tcl9.0/TkCmd/photo.html
func (m *Img) Copy(src *Img, options ...Opt) (r *Img) {
	evalErr(fmt.Sprintf("%s copy %s %s", m, src, collect(options...)))
	return m
}

// From option.
//
// Known uses:
//   - [Img.Copy]
//   - [Scale] (widget specific)
//   - [Spinbox] (widget specific)
//   - [TScale] (widget specific)
//   - [TSpinbox] (widget specific)
func From(val ...any) Opt {
	return rawOption(fmt.Sprintf(`-from %s`, collectAny(val...)))
}

// To option.
//
// Known uses:
//   - [Img.Copy]
//   - [Scale] (widget specific)
//   - [Spinbox] (widget specific)
//   - [TScale] (widget specific)
//   - [TSpinbox] (widget specific)
func To(val ...any) Opt {
	return rawOption(fmt.Sprintf(`-to %s`, collectAny(val...)))
}

// Graph — use gnuplot to draw on a photo. Graph returns 'm'
//
// The 'script' argument is passed to a gnuplot executable, which must be
// installed on the machine.  See the [gnuplot site] for documentation about
// producing graphs. The script must not use the 'set term <device>' command.
//
// The content of 'm' is replaced, including its internal name.
//
// [gnuplot site]: http://www.gnuplot.info/
func (m *Img) Graph(script string) *Img {
	switch {
	case strings.HasPrefix(m.name, "img"):
		w, h := m.Width(), m.Height()
		script = fmt.Sprintf("set terminal pngcairo size %s, %s\n%s", w, h, script)
		out, err := gnuplot(script)
		if err != nil {
			fail(fmt.Errorf("plot: executing script: %s", err))
			break
		}

		*m = *NewPhoto(Width(w), Height(h), Data(out))
	default:
		fail(fmt.Errorf("plot: %s is not a photo", m))
	}
	return m
}

// Destroy — Destroy one or more windows
//
// # Description
//
// This command deletes the windows given by the window arguments, plus all of
// their descendants. If a window “.” (App) is deleted then all windows will be
// destroyed and the application will (normally) exit. The windows are
// destroyed in order, and if an error occurs in destroying a window the
// command aborts without destroying the remaining windows. No error is
// returned if window does not exist.
func Destroy(options ...Opt) {
	evalErr(fmt.Sprintf("destroy %s", collect(options...)))
}

// Pack — Geometry manager that packs around edges of cavity
//
// # Description
//
// The options consist of one or more content windows followed
// by options that specify how to manage the content. See THE PACKER
// ALGORITHM for details on how the options are used by the packer.
//
// The first argument must be a *Window.
//
// The following options are supported:
//
//   - [After] other
//
// Other must the name of another window. Use its container as the container
// for the content, and insert the content just after other in the packing
// order.
//
//   - [Anchor] anchor
//
// Anchor must be a valid anchor position such as n or sw; it specifies where
// to position each content in its parcel. Defaults to center.
//
//   - [Before] other
//
// Other must the name of another window. Use its container as the container
// for the content, and insert the content just before other in the packing
// order.
//
//   - [Expand] boolean
//
// Specifies whether the content should be expanded to consume extra space in
// their container. Boolean may have any proper boolean value, such as 1 or no.
// Defaults to 0.
//
//   - [Fill] style
//
// If a content's parcel is larger than its requested dimensions, this option
// may be used to stretch the content. Style must have one of the following
// values:
//
//   - "none" - Give the content its requested dimensions plus any internal
//     padding requested with -ipadx or -ipady. This is the default.
//   - "x" - Stretch the content horizontally to fill the entire width of its
//     parcel (except leave external padding as specified by -padx).
//   - "y" - Stretch the content vertically to fill the entire height of its parcel
//     (except leave external padding as specified by -pady).
//   - "both": Stretch the content both horizontally and vertically.
//
// .
//
//   - [In] container
//
// Insert the window at the end of the packing order for the container window
// given by container.
//
//   - [Ipadx] amount
//
// Amount specifies how much horizontal internal padding to leave on each side
// of the content. Amount must be a valid screen distance, such as 2 or .5c. It
// defaults to 0.
//
//   - [Ipady] amount
//
// Amount specifies how much vertical internal padding to leave on each side of
// the content. Amount defaults to 0.
//
//   - [Padx] amount
//
// Amount specifies how much horizontal external padding to leave on each side
// of the content. Amount may be a list of two values to specify padding for
// left and right separately. Amount defaults to 0.
//
//   - [Pady] amount
//
// Amount specifies how much vertical external padding to leave on each side of
// the content. Amount may be a list of two values to specify padding for top
// and bottom separately. Amount defaults to 0.
//
//   - [Side] side
//
// Specifies which side of the container the content will be packed against.
// Must be "left", "right", "top", or "bottom". Defaults to top.
//
// If no -in, -after or -before option is specified then each of the content
// will be inserted at the end of the packing list for its parent unless it is
// already managed by the packer (in which case it will be left where it is).
// If one of these options is specified then all the content will be inserted
// at the specified point. If any of the content are already managed by the
// geometry manager then any unspecified options for them retain their previous
// values rather than receiving default values.
//
// Additional information might be available at the [Tcl/Tk pack] page.
//
// [Tcl/Tk pack]: https://www.tcl.tk/man/tcl9.0/TkCmd/pack.html
func Pack(options ...Opt) {
	evalErr(fmt.Sprintf("pack %s", collect(options...)))
}

// SetResizable — Enable/disable window resizing
//
// # Description
//
// This command controls whether or not the user may interactively resize a
// top-level window. If resizing is disabled, then the window's size will be
// the size from the most recent interactive resize or wm geometry command. If
// there has been no such operation then the window's natural size will be
// used.
//
// More information might be available at the [Tcl/Tk wm] page.
//
// [Tcl/Tk wm]: https://www.tcl.tk/man/tcl9.0/TkCmd/wm.html
func (w *Window) SetResizable(width, height bool) {
	evalErr(fmt.Sprintf("wm resizable %s %v %v", w, width, height))
}

// Wait — Wait for a window to be destroyed
//
// # Description
//
// Wait command waits for 'w' to be destroyed. This is typically used to wait
// for a user to finish interacting with a dialog box before using the result
// of that interaction.
//
// While the Wwait command is waiting it processes events in the normal
// fashion, so the application will continue to respond to user interactions.
// If an event handler invokes Wait again, the nested call to Wait must
// complete before the outer call can complete.
func (w *Window) Wait() {
	if isBuilder && testHookWait != "" {
		TclAfterIdle(Command(func() {
			fmt.Println(testHookWait)
		}))
	}
	if w == App {
		switch {
		case os.Getenv("TK9_VNC") == "1":
			autocenterDisabled = true
			WmGeometry(App, fmt.Sprintf("%sx%s+0+0", os.Getenv("TK9_VNC_WIDTH"), os.Getenv("TK9_VNC_HEIGHT")))
		case forcedX >= 0 && forcedY >= 0: // Behind TK9_DEMO=1.
			evalErr(fmt.Sprintf("wm geometry . +%v+%v", forcedX, forcedY)) //TODO add API func
			forcedX, forcedY = -1, -1                                      // Apply only the first time.
		default:
			if goos == "windows" {
				if !appWithdrawn && !appIconified {
					WmDeiconify(App)
				}
			}
			if !autocenterDisabled {
				autocenterDisabled = true
				w.Center()
			}
		}
	}
	evalErr(fmt.Sprintf("tkwait window %s", w))
}

// WaitVisibility — Wait for a window to change visibility
//
// # Description
//
// WaitVisibility command waits for a change in w's visibility state (as
// indicated by the arrival of a VisibilityNotify event). This form is
// typically used to wait for a newly-created window to appear on the screen
// before taking some action.
//
// While the Wwait command is waiting it processes events in the normal
// fashion, so the application will continue to respond to user interactions.
// If an event handler invokes Wait again, the nested call to Wait must
// complete before the outer call can complete.
func (w *Window) WaitVisibility() {
	evalErr(fmt.Sprintf("tkwait visibility %s", w))
}

// IconPhoto — change window icon
//
// # Description
//
// IconPhoto sets the titlebar icon for window based on the named photo images.
// If -default is specified, this is applied to all future created toplevels as
// well. The data in the images is taken as a snapshot at the time of
// invocation. If the images are later changed, this is not reflected to the
// titlebar icons. Multiple images are accepted to allow different images sizes
// (e.g., 16x16 and 32x32) to be provided. The window manager may scale
// provided icons to an appropriate size.
//
// On Windows, the images are packed into a Windows icon structure. This will
// override an ico specified to wm iconbitmap, and vice versa. This command
// sets the taskbar icon for the window.
//
// On X, the images are arranged into the _NET_WM_ICON X property, which most
// modern window managers support. A wm iconbitmap may exist simultaneously. It is
// recommended to use not more than 2 icons, placing the larger icon first. This
// command also sets the panel icon for the application if the window manager or
// desktop environment supports it.
//
// On Macintosh, the first image called is loaded into an OSX-native icon
// format, and becomes the application icon in dialogs, the Dock, and other
// contexts. At the script level the command will accept only the first image
// passed in the parameters as support for multiple sizes/resolutions on macOS
// is outside Tk's scope. Developers should use the largest icon they can
// support (preferably 512 pixels) to ensure smooth rendering on the Mac.
//
// More information might be available at the [Tcl/Tk wm] page.
//
// [Tcl/Tk wm]: https://www.tcl.tk/man/tcl9.0/TkCmd/wm.html
func (w *Window) IconPhoto(options ...Opt) {
	evalErr(fmt.Sprintf("wm iconphoto %s %s", w, collect(options...)))
}

// DefaultIcon option.
//
// Known uses:
//   - [IconPhoto] (command specific)
func DefaultIcon(val ...any) Opt {
	return rawOption("-default")
}

// WmTitle — change the window manager title
//
// # Description
//
// If string is specified, then it will be passed to the window manager for use
// as the title for window (the window manager should display this string in
// window's title bar). In this case the command returns an empty string. If
// string is not specified then the command returns the current title for the
// window. The title for a window defaults to its name.
//
// More information might be available at the [Tcl/Tk wm] page.
//
// [Tcl/Tk wm]: https://www.tcl.tk/man/tcl9.0/TkCmd/wm.html
func (w *Window) WmTitle(s string) string {
	if s != "" {
		s = tclSafeString(s)
	}
	return evalErr(fmt.Sprintf("wm title %s %s", w, s))
}

// Center centers 'w' and returns 'w'.
func (w *Window) Center() *Window {
	if w == App {
		autocenterDisabled = true
	}
	evalErr(fmt.Sprintf("tk::PlaceWindow %s center", w))
	return w
}

// Grid — Geometry manager that arranges widgets in a grid
//
// # Description
//
// Removes each of the windows from grid for its container and unmaps their
// windows. The content will no longer be managed by the grid geometry manager.
// The configuration options for that window are forgotten, so that if the
// window is managed once more by the grid geometry manager, the initial
// default settings are used.
//
// If the last content window of the container becomes unmanaged, this will
// also send the virtual event <<NoManagedChild>> to the container; the
// container may choose to resize itself (or otherwise respond) to such a
// change.
//
// More information might be available at the [Tcl/Tk grid] page.
//
// [Tcl/Tk grid]: https://www.tcl.tk/man/tcl9.0/TkCmd/grid.html
func GridForget(w ...*Window) {
	if len(w) == 0 {
		return
	}

	var a []string
	for _, v := range w {
		if v != nil {
			a = append(a, v.String())
		}
	}
	evalErr(fmt.Sprintf("grid forget %s", strings.Join(a, " ")))
}

// Grid — Geometry manager that arranges widgets in a grid
//
// # Description
//
// Removes each of the windows from grid for its container and unmaps their
// windows. The content will no longer be managed by the grid geometry manager.
// However, the configuration options for that window are remembered, so that
// if the content window is managed once more by the grid geometry manager, the
// previous values are retained.
//
// If the last content window of the container becomes unmanaged, this will
// also send the virtual event <<NoManagedChild>> to the container; the
// container may choose to resize itself (or otherwise respond) to such a
// change.
//
// More information might be available at the [Tcl/Tk grid] page.
//
// [Tcl/Tk grid]: https://www.tcl.tk/man/tcl9.0/TkCmd/grid.html
func GridRemove(w ...*Window) {
	if len(w) == 0 {
		return
	}

	var a []string
	for _, v := range w {
		if v != nil {
			a = append(a, v.String())
		}
	}
	evalErr(fmt.Sprintf("grid remove %s", strings.Join(a, " ")))
}

// Grid — Geometry manager that arranges widgets in a grid
//
// # Description
//
// The arguments consist of the names of one or more content windows followed
// by pairs of arguments that specify how to manage the content. The characters
// -, x and ^, can be specified instead of a window name to alter the default
// location of a window, as described in the RELATIVE PLACEMENT section, below.
//
// The following options are supported:
//
//   - [Column] n
//
// Insert the window so that it occupies the nth column in the grid. Column
// numbers start with 0. If this option is not supplied, then the window is
// arranged just to the right of previous window specified on this call to
// grid, or column “0” if it is the first window. For each x that immediately
// precedes the window, the column position is incremented by one. Thus the x
// represents a blank column for this row in the grid.
//
//   - [Columnspan] n
//
// Insert the window so that it occupies n columns in the grid. The default is
// one column, unless the window name is followed by a -, in which case the
// columnspan is incremented once for each immediately following -.
//
//   - [In] container
//
// Insert the window(s) in the container window given by container. The default
// is the first window's parent window.
//
//   - [Ipadx] amount
//
// The amount specifies how much horizontal internal padding to leave on each
// side of the content. This is space is added inside the content border. The
// amount must be a valid screen distance, such as 2 or .5c. It defaults to 0.
//
//   - [Ipady] amount
//
// The amount specifies how much vertical internal padding to leave on the top
// and bottom of the content. This space is added inside the content border.
// The amount defaults to 0.
//
//   - [Padx] amount
//
// The amount specifies how much horizontal external padding to leave on each
// side of the content, in screen units. Amount may be a list of two values to
// specify padding for left and right separately. The amount defaults to 0.
// This space is added outside the content border.
//
//   - [Pady] amount
//
// The amount specifies how much vertical external padding to leave on the top
// and bottom of the content, in screen units. Amount may be a list of two
// values to specify padding for top and bottom separately. The amount defaults
// to 0. This space is added outside the content border.
//
//   - [Row] n
//
// Insert the content so that it occupies the nth row in the grid. Row numbers
// start with 0. If this option is not supplied, then the content is arranged
// on the same row as the previous content specified on this call to grid, or
// the next row after the highest occupied row if this is the first content.
//
//   - [Rowspan] n
//
// Insert the content so that it occupies n rows in the grid. The default is
// one row. If the next grid command contains ^ characters instead of content
// that line up with the columns of this content, then the rowspan of this
// content is extended by one.
//
//   - [Sticky] style
//
// If a content's cell is larger than its requested dimensions, this option may
// be used to position (or stretch) the content within its cell. Style is a
// string that contains zero or more of the characters n, s, e or w. The string
// can optionally contain spaces or commas, but they are ignored. Each letter
// refers to a side (north, south, east, or west) that the content will “stick”
// to. If both n and s (or e and w) are specified, the content will be
// stretched to fill the entire height (or width) of its cavity. The -sticky
// option subsumes the combination of -anchor and -fill that is used by pack.
// The default is “”, which causes the content to be centered in its cavity, at
// its requested size.
//
// If any of the content is already managed by the geometry manager then any
// unspecified options for them retain their previous values rather than
// receiving default values.
//
// More information might be available at the [Tcl/Tk grid] page.
//
// [Tcl/Tk grid]: https://www.tcl.tk/man/tcl9.0/TkCmd/grid.html
func Grid(w Widget, options ...Opt) {
	evalErr(fmt.Sprintf("grid configure %s %s", w, collect(options...)))
}

// Grid — Geometry manager that arranges widgets in a grid
//
// # Description
//
// The anchor value controls how to place the grid within the container window
// when no row/column has any weight. See THE GRID ALGORITHM below for further
// details. The default anchor is nw.
//
// More information might be available at the [Tcl/Tk grid] page.
//
// [Tcl/Tk grid]: https://www.tcl.tk/man/tcl9.0/TkCmd/grid.html
func GridAnchor(w *Window, anchor string) string {
	return evalErr(fmt.Sprintf("grid anchor %s %s", w, tclSafeString(anchor)))
}

// Grid — Geometry manager that arranges widgets in a grid
//
// # Description
//
// Query or set the row properties of the index row of the geometry container,
// window. The valid options are -minsize, -weight, -uniform and -pad. If one
// or more options are provided, then index may be given as a list of row
// indices to which the configuration options will operate on. Indices may be
// integers, window names or the keyword all. For all the options apply to all
// rows currently occupied by content windows. For a window name, that window
// must be a content window of this container and the options apply to all rows
// currently occupied by the container window. The -minsize option sets the
// minimum size, in screen units, that will be permitted for this row. The
// -weight option (an integer value) sets the relative weight for apportioning
// any extra spaces among rows. A weight of zero (0) indicates the row will not
// deviate from its requested size. A row whose weight is two will grow at
// twice the rate as a row of weight one when extra space is allocated to the
// layout. The -uniform option, when a non-empty value is supplied, places the
// row in a uniform group with other rows that have the same value for
// -uniform. The space for rows belonging to a uniform group is allocated so
// that their sizes are always in strict proportion to their -weight values.
// See THE GRID ALGORITHM below for further details. The -pad option specifies
// the number of screen units that will be added to the largest window
// contained completely in that row when the grid geometry manager requests a
// size from the containing window. If only an option is specified, with no
// value, the current value of that option is returned. If only the container
// window and index is specified, all the current settings are returned in a
// list of “-option value” pairs.
//
// More information might be available at the [Tcl/Tk grid] page.
//
// [Tcl/Tk grid]: https://www.tcl.tk/man/tcl9.0/TkCmd/grid.html
func GridRowConfigure(w Widget, index any, options ...Opt) {
	switch x := index.(type) {
	case []int:
		s := fmt.Sprint(x)
		index = "{" + s[1:len(s)-1] + "}"
	default:
		index = tclSafeString(fmt.Sprint(index))
	}
	evalErr(fmt.Sprintf("grid rowconfigure %s %v %s", w, index, collect(options...)))
}

// Minsize option.
//
// Known uses:
//   - [GridColumnConfigure]
//   - [GridRowConfigure]
func Minsize(val ...any) Opt {
	return rawOption(fmt.Sprintf(`-minsize %s`, collectAny(val...)))
}

// Grid — Geometry manager that arranges widgets in a grid
//
// # Description
//
// Query or set the column properties of the index column of the geometry
// container, window. The valid options are -minsize, -weight, -uniform and
// -pad. If one or more options are provided, then index may be given as a list
// of column indices to which the configuration options will operate on.
// Indices may be integers, window names or the keyword all. For all the
// options apply to all columns currently occupied be content windows. For a
// window name, that window must be a content of this container and the options
// apply to all columns currently occupied be the content. The -minsize option
// sets the minimum size, in screen units, that will be permitted for this
// column. The -weight option (an integer value) sets the relative weight for
// apportioning any extra spaces among columns. A weight of zero (0) indicates
// the column will not deviate from its requested size. A column whose weight
// is two will grow at twice the rate as a column of weight one when extra
// space is allocated to the layout. The -uniform option, when a non-empty
// value is supplied, places the column in a uniform group with other columns
// that have the same value for -uniform. The space for columns belonging to a
// uniform group is allocated so that their sizes are always in strict
// proportion to their -weight values. See THE GRID ALGORITHM below for further
// details. The -pad option specifies the number of screen units that will be
// added to the largest window contained completely in that column when the
// grid geometry manager requests a size from the containing window. If only an
// option is specified, with no value, the current value of that option is
// returned. If only the container window and index is specified, all the
// current settings are returned in a list of “-option value” pairs.
//
// More information might be available at the [Tcl/Tk grid] page.
//
// [Tcl/Tk grid]: https://www.tcl.tk/man/tcl9.0/TkCmd/grid.html
func GridColumnConfigure(w Widget, index any, options ...Opt) {
	switch x := index.(type) {
	case []int:
		s := fmt.Sprint(x)
		index = "{" + s[1:len(s)-1] + "}"
	default:
		index = tclSafeString(fmt.Sprint(index))
	}
	evalErr(fmt.Sprintf("grid columnconfigure %s %v %s", w, index, collect(options...)))
}

// Pad option.
//
// The -pad option specifies the number of screen units that will be added to
// the largest window contained completely in that column/row when the grid
// geometry manager requests a size from the containing window.
//
// Known uses:
//   - [GridColumnConfigure] (command specific)
//   - [GridRowConfigure] (command specific)
func Pad(val any) Opt {
	return rawOption(fmt.Sprintf(`-pad %s`, tclSafeString(fmt.Sprint(val))))
}

// Configure alters the configuration of 'w' and returns 'w'.
func (w *Window) Configure(options ...Opt) *Window {
	options, tvs, vs := w.split(options)
	if len(options) != 0 {
		evalErr(fmt.Sprintf("%s configure %s", w, collect(options...)))
	}
	if len(tvs) != 0 {
		tvo := tvs[len(tvs)-1]
		tclVar := textVariables[w]
		if tclVar == "" {
			tclVar = fmt.Sprintf("textVar%d", id.Add(1))
			textVariables[w] = tclVar
			evalErr(fmt.Sprintf("%s configure -textvariable %s", w, tclVar))
		}
		evalErr(fmt.Sprintf("set %s %s", tclVar, tclSafeString(string(tvo))))
	}
	if len(vs) != 0 {
		vo := vs[len(vs)-1]
		variables[w] = vo
		if vo.tclName == "" {
			vo.tclName = fmt.Sprintf("goVar%d", id.Add(1))
		}
		evalErr(fmt.Sprintf("%s configure -variable %s", w, vo.tclName))
		evalErr(fmt.Sprintf("set %s %s", vo.tclName, tclSafeString(fmt.Sprint(vo.val))))
	}
	return w
}

// ttk::widget — Standard options and commands supported by Tk themed widgets
//
// # Description
//
// Modify or inquire widget state. If stateSpec is not "", sets the widget
// state: for each flag in stateSpec, sets the corresponding flag or clears it
// if prefixed by an exclamation point.  Returns a new state spec indicating
// which flags were changed.
//
// If stateSpec is "", returns a list of the currently-enabled state flags.
//
// # Widget States
//
// The widget state is a bitmap of independent state flags. Widget state flags include:
//
//   - active:
//     The mouse cursor is over the widget and pressing a mouse button will
//     cause some action to occur. (aka “prelight” (Gnome), “hot” (Windows),
//     “hover”).
//   - disabled:
//     Widget is disabled under program control (aka “unavailable”,
//     “inactive”).
//   - focus:
//     Widget has keyboard focus.
//   - pressed:
//     Widget is being pressed (aka “armed” in Motif).
//   - selected:
//     “On”, “true”, or “current” for things like checkbuttons and
//     radiobuttons.
//   - background:
//     Windows and the Mac have a notion of an “active” or foreground window.
//     The background state is set for widgets in a background window, and
//     cleared for those in the foreground window.
//   - readonly:
//     Widget should not allow user modification.
//   - alternate:
//     A widget-specific alternate display format. For example, used for
//     checkbuttons and radiobuttons in the “tristate” or “mixed” state, and
//     for buttons with -default active.
//   - invalid:
//     The widget's value is invalid. (Potential uses: scale widget value out
//     of bounds, entry widget value failed validation.)
//   - hover:
//     The mouse cursor is within the widget. This is similar to the active
//     state; it is used in some themes for widgets that provide distinct
//     visual feedback for the active widget in addition to the active element
//     within the widget.
//   - user1-user6
//     Freely usable for other purposes
//
// A state specification or stateSpec is a list of state names, optionally
// prefixed with an exclamation point (!) indicating that the bit is off.
func (w *Window) WidgetState(stateSpec string) (r string) {
	if stateSpec != "" {
		stateSpec = tclSafeString(stateSpec)
	}
	return evalErr(fmt.Sprintf("%s state %s", w, stateSpec))
}

// tk_messageBox — pops up a message window and waits for user response.
//
// # Description
//
// This procedure creates and displays a message window with an
// application-specified message, an icon and a set of buttons. Each of the
// buttons in the message window is identified by a unique symbolic name (see
// the -type options). After the message window is popped up, tk_messageBox
// waits for the user to select one of the buttons. Then it returns the
// symbolic name of the selected button.  The following optins are
// supported:
//
//   - [Command] handler
//
// Specifies the handler to invoke when the user closes the
// dialog. The actual command consists of string followed by a space and the
// name of the button clicked by the user to close the dialog. This is only
// available on Mac OS X.
//
//   - [Default] name
//
// Name gives the symbolic name of the default button for this message window (
// “ok”, “cancel”, and so on). See -type for a list of the symbolic names. If
// this option is not specified, the first button in the dialog will be made
// the default.
//
//   - [Detail] string
//
// Specifies an auxiliary message to the main message given by the -message
// option. The message detail will be presented beneath the main message and,
// where supported by the OS, in a less emphasized font than the main message.
//
//   - [Icon] iconImage
//
// Specifies an icon to display. IconImage must be one of the following: error,
// info, question or warning. If this option is not specified, then the info
// icon will be displayed.
//
//   - [Message] string
//
// Specifies the message to display in this message box. The default value is
// an empty string.
//
//   - [Parent] window
//
// Makes window the logical parent of the message box. The message box is
// displayed on top of its parent window.
//
//   - [Title] titleString
//
// Specifies a string to display as the title of the message box. The default
// value is an empty string.
//
//   - [Type] predefinedType
//
// Arranges for a predefined set of buttons to be displayed. The following
// values are possible for predefinedType:
//
//   - abortretryignore - Displays three buttons whose symbolic names are abort, retry and ignore.
//   - ok - Displays one button whose symbolic name is ok.
//   - okcancel - Displays two buttons whose symbolic names are ok and cancel.
//   - retrycancel - Displays two buttons whose symbolic names are retry and cancel.
//   - yesno - Displays two buttons whose symbolic names are yes and no.
//   - yesnocancel - Displays three buttons whose symbolic names are yes, no and cancel.
//
// More information might be available at the [Tcl/Tk messageBox] page.
//
// [Tcl/Tk messageBox]: https://www.tcl.tk/man/tcl9.0/TkCmd/messageBox.html
func MessageBox(options ...Opt) string {
	return evalErr(fmt.Sprintf("tk_messageBox %s", collect(options...)))
}

// Bell — Ring a display's bell
//
// # Description
//
// This command rings the bell on the display for window and returns an empty
// string. If the -displayof option is omitted, the display of the
// application's main window is used by default. The command uses the current
// bell-related settings for the display, which may be modified with programs
// such as xset.
//
// If -nice is not specified, this command also resets the screen saver for the
// screen. Some screen savers will ignore this, but others will reset so that
// the screen becomes visible again.
//
//   - [Displayof] window
//   - [Nice]
//
// More information might be available at the [Tcl/Tk bell] page.
//
// [Tcl/Tk bell]: https://www.tcl.tk/man/tcl9.0/TclCmd/bell.html
func Bell(options ...Opt) {
	evalErr(fmt.Sprintf("bell %s", collect(options...)))
}

// ChooseColor — pops up a dialog box for the user to select a color.
//
// # Description
//
// ChooseColor pops up a dialog box for the user to select a color. The
// following option-value pairs are possible as command line arguments:
//
//   - [Initialcolor] color
//
// Specifies the color to display in the color dialog when it pops up. color
// must be in a form acceptable to the Tk_GetColor function.
//
//   - [Parent] window
//
// Makes window the logical parent of the color dialog. The color dialog is
// displayed on top of its parent window.
//
//   - [Title] titleString
//
// Specifies a string to display as the title of the dialog box. If this option
// is not specified, then a default title will be displayed.
//
// If the user selects a color, ChooseColor will return the name of the
// color in a form acceptable to Tk_GetColor. If the user cancels the
// operation, both commands will return the empty string.
//
// More information might be available at the [Tcl/Tk choosecolor] page.
//
// [Tcl/Tk choosecolor]: https://www.tcl.tk/man/tcl9.0/TclCmd/chooseColor.html
func ChooseColor(options ...Opt) string {
	return evalErr(fmt.Sprintf("tk_chooseColor %s", collect(options...)))
}

// Busy — confine pointer events to a window sub-tree
//
// # Description
//
// The Busy command provides a simple means to block pointer events from Tk
// widgets, while overriding the widget's cursor with a configurable busy
// cursor. Note this command does not prevent keyboard events from being sent
// to the widgets made busy.
//
//   - [Cursor] cursorName
//
// Specifies the cursor to be displayed when the widget is made busy.
// CursorName can be in any form accepted by Tk_GetCursor. The default cursor
// is wait on Windows and watch on other platforms.
//
// More information might be available at the [Tcl/Tk busy] page.
//
// [Tcl/Tk update]: https://www.tcl.tk/man/tcl9.0/TclCmd/busy.html
func (w *Window) Busy(options ...Opt) {
	evalErr(fmt.Sprintf("tk busy %s %s", w, collect(options...)))
}

// BusyForget — undo Busy
//
// # Description
//
// Releases resources allocated by the [Window.Busy] command for window, including
// the transparent window. User events will again be received by window.
// Resources are also released when window is destroyed. Window must be the
// name of a widget specified in the hold operation, otherwise an error is
// reported.
//
// More information might be available at the [Tcl/Tk busy] page.
//
// [Tcl/Tk update]: https://www.tcl.tk/man/tcl9.0/TclCmd/busy.html
func (w *Window) BusyForget(options ...Opt) {
	evalErr(fmt.Sprintf("tk busy forget %s %s", w, collect(options...)))
}

// Update — Process pending events and idle callbacks
//
// More information might be available at the [Tcl/Tk update] page.
//
// [Tcl/Tk update]: https://www.tcl.tk/man/tcl9.0/TclCmd/update.html
func Update() {
	evalErr("update")
}

//TODO?
//  - [Command] string
//
// Specifies the prefix of a Tcl command to invoke when the user closes the
// dialog after having selected an item. This callback is not called if the
// user cancelled the dialog. The actual command consists of string followed by
// a space and the value selected by the user in the dialog. This is only
// available on Mac OS X.

// ChooseDirectory — pops up a dialog box for the user to select a directory.
//
// # Description
//
// The procedure tk_chooseDirectory pops up a dialog box for the user to select
// a directory. The following option-value pairs are possible as command line
// arguments:
//
//   - [Initialdir] dirname
//
// Specifies that the directories in directory should be displayed when the
// dialog pops up. If this parameter is not specified, the initial directory
// defaults to the current working directory on non-Windows systems and on
// Windows systems prior to Vista. On Vista and later systems, the initial
// directory defaults to the last user-selected directory for the application.
// If the parameter specifies a relative path, the return value will convert
// the relative path to an absolute path.
//
//   - [Message] string
//
// Specifies a message to include in the client area of the dialog. This is
// only available on Mac OS X.
//
//   - [Mustexist] boolean
//
// Specifies whether the user may specify non-existent directories. If this
// parameter is true, then the user may only select directories that already
// exist. The default value is false.
//
//   - [Parent] window
//
// Makes window the logical parent of the dialog. The dialog is displayed on
// top of its parent window. On Mac OS X, this turns the file dialog into a
// sheet attached to the parent window.
//
//   - [Title] titleString
//
// Specifies a string to display as the title of the dialog box. If this option
// is not specified, then a default title will be displayed.
//
// More information might be available at the [Tcl/Tk chooseDirectory] page.
//
// [Tcl/Tk chooseDirectory]: https://www.tcl.tk/man/tcl9.0/TkCmd/chooseDirectory.html
func ChooseDirectory(options ...Opt) string {
	return evalErr(fmt.Sprintf("tk_chooseDirectory %s", collect(options...)))
}

// ClipboardAppend — Manipulate Tk clipboard
//
// # Description
//
// This command provides a Tcl interface to the Tk clipboard, which stores data
// for later retrieval using the selection mechanism (via the -selection
// CLIPBOARD option). In order to copy data into the clipboard, clipboard clear
// must be called, followed by a sequence of one or more calls to clipboard
// append. To ensure that the clipboard is updated atomically, all appends
// should be completed before returning to the event loop.
//
// ClipboardAppend appends 'data' to the clipboard on window's display in the
// form given by type with the representation given by format and claims
// ownership of the clipboard on window's display.
//
//   - [Displayof] window
//
//   - [Format] format
//
// The format argument specifies the representation that should be used to
// transmit the selection to the requester (the second column of Table 2 of the
// ICCCM), and defaults to STRING. If format is STRING, the selection is
// transmitted as 8-bit ASCII characters. If format is ATOM, then the data is
// divided into fields separated by white space; each field is converted to its
// atom value, and the 32-bit atom value is transmitted instead of the atom
// name. For any other format, data is divided into fields separated by white
// space and each field is converted to a 32-bit integer; an array of integers
// is transmitted to the selection requester. Note that strings passed to
// clipboard append are concatenated before conversion, so the caller must take
// care to ensure appropriate spacing across string boundaries. All items
// appended to the clipboard with the same type must have the same format.
//
// The format argument is needed only for compatibility with clipboard
// requesters that do not use Tk. If the Tk toolkit is being used to retrieve
// the CLIPBOARD selection then the value is converted back to a string at the
// requesting end, so format is irrelevant.
//
//   - [Type] type
//
// Type specifies the form in which the selection is to be returned (the
// desired “target” for conversion, in ICCCM terminology), and should be an
// atom name such as STRING or FILE_NAME; see the Inter-Client Communication
// Conventions Manual for complete details. Type defaults to STRING.
//
// More information might be available at the [Tcl/Tk clipboard] page.
//
// [Tcl/Tk clipboard]: https://www.tcl.tk/man/tcl9.0/TkCmd/clipboard.html
func ClipboardAppend(data string, options ...Opt) {
	evalErr(fmt.Sprintf("clipboard append %s -- %s", collect(options...), tclSafeString(data)))
}

// ClipboardClear — Manipulate Tk clipboard
//
// # Description
//
// Claims ownership of the clipboard on window's display and removes any
// previous contents. Window defaults to App. Returns an empty string.
//
// More information might be available at the [Tcl/Tk clipboard] page.
//
// [Tcl/Tk clipboard]: https://www.tcl.tk/man/tcl9.0/TkCmd/clipboard.html
func ClipboardClear(options ...Opt) {
	evalErr(fmt.Sprintf("clipboard clear %s", collect(options...)))
}

// ClipboardGet — Manipulate Tk clipboard
//
// # Description
//
// Retrieve data from the clipboard on window's display. Window defaults to App.
//
//   - [Displayof] window
//
//   - [Type] type
//
// Type specifies the form in which the data is to be returned and should be an
// atom name such as STRING or FILE_NAME. Type defaults to STRING. This command
// is equivalent to [SelectionGet](Selection("CLIPBOARD").
//
// Note that on modern X11 systems, the most useful type to retrieve for
// transferred strings is not STRING, but rather UTF8_STRING.
//
// More information might be available at the [Tcl/Tk clipboard] page.
//
// [Tcl/Tk clipboard]: https://www.tcl.tk/man/tcl9.0/TkCmd/clipboard.html
func ClipboardGet(options ...Opt) string {
	return evalErr(fmt.Sprintf("clipboard get %s", collect(options...)))
}

var onceForceInit bool

func forceInit() {
	if onceForceInit {
		return
	}

	defer func() { onceForceInit = true }()

	evalErr("# forceInit() executed")
}

// ExitHandler returns a canned [Command] that destroys the [App].
func ExitHandler() Opt {
	forceInit()
	return exitHandler
}

// Exit provides a canned [Button] with default [Txt] "Exit", bound to the
// [ExitHandler].
//
// Use [Window.Exit] to create a Exit with a particular parent.
func Exit(options ...Opt) *ButtonWidget {
	return App.Exit(options...)
}

// Exit provides a canned [Button] with default [Txt] "Exit", bound to the
// [ExitHandler].
//
// The resulting [Window] is a child of 'w'
func (w *Window) Exit(options ...Opt) *ButtonWidget {
	switch {
	case isVNC:
		return Tooltip(w.Button(append([]Opt{Txt("Disconnect"), ExitHandler()}, options...)...), disconnectButtonTooltip).(*ButtonWidget)
	default:
		return Tooltip(w.Button(append([]Opt{Txt("Exit"), ExitHandler()}, options...)...), exitButtonTooltip).(*ButtonWidget)
	}
}

// TExit provides a canned [TButton] with default [Txt] "Exit", bound to the
// [ExitHandler].
//
// Use [Window.TExit] to create a TExit with a particular parent.
func TExit(options ...Opt) *TButtonWidget {
	return App.TExit(options...)
}

// TExit provides a canned [TButton] with default [Txt] "Exit", bound to the
// [ExitHandler].
//
// The resulting [Window] is a child of 'w'
func (w *Window) TExit(options ...Opt) *TButtonWidget {
	switch {
	case isVNC:
		return Tooltip(w.TButton(append([]Opt{Txt("Disconnect"), ExitHandler()}, options...)...), disconnectButtonTooltip).(*TButtonWidget)
	default:
		return Tooltip(w.TButton(append([]Opt{Txt("Exit"), ExitHandler()}, options...)...), exitButtonTooltip).(*TButtonWidget)
	}
}

var _ Opt = (*VariableOpt)(nil)

// VariableOpt is an Opt linking Go and Tcl variables.
type VariableOpt struct {
	tclName string
	val     any
}

// Set sets the value of the linked Tcl variable.
func (v *VariableOpt) Set(val any) {
	if v == nil || v.tclName == "" {
		fail(fmt.Errorf("%T not linked", v))
		return
	}

	evalErr(fmt.Sprintf("set %s %s", v.tclName, tclSafeString(fmt.Sprint(val))))
}

// Get return the value of the linked Tcl variable.
func (v *VariableOpt) Get() (r string) {
	if v == nil || v.tclName == "" {
		fail(fmt.Errorf("%T not linked", v))
		return
	}

	return evalErr(fmt.Sprintf("set %s", v.tclName))
}

func (*VariableOpt) optionString(*Window) string {
	panic("internal error") // Not supposed to be invoked.
}

// Variable option.
//
// Known uses:
//   - [Checkbutton] (widget specific)
//   - [MenuWidget.AddCascade] (command specific)
//   - [MenuWidget.AddCommand] (command specific)
//   - [MenuWidget.AddSeparator] (command specific)
//   - [Radiobutton] (widget specific)
//   - [Scale] (widget specific)
//   - [TCheckbutton] (widget specific)
//   - [TProgressbar] (widget specific)
//   - [TRadiobutton] (widget specific)
//   - [TScale] (widget specific)
func Variable(val any) *VariableOpt {
	return &VariableOpt{val: val}
}

// Variable — Get the configured option value.
//
// Known uses:
//   - [Checkbutton] (widget specific)
//   - [Radiobutton] (widget specific)
//   - [Scale] (widget specific)
//   - [TCheckbutton] (widget specific)
//   - [TProgressbar] (widget specific)
//   - [TRadiobutton] (widget specific)
//   - [TScale] (widget specific)
func (w *Window) Variable() string {
	if tclVar := variables[w]; tclVar != nil {
		return evalErr(fmt.Sprintf("set %s", tclVar.tclName))
	}

	return ""
}

type textVarOpt string

func (textVarOpt) optionString(*Window) string {
	panic("internal error") // Not supposed to be invoked.
}

// Textvariable option.
//
// Specifies the value to be displayed inside the widget.
// The way in which the string is displayed in the widget depends on the
// particular widget and may be determined by other options, such as
// -anchor or -justify.
//
// Known uses:
//   - [Button]
//   - [Checkbutton]
//   - [Entry]
//   - [Label]
//   - [Menubutton]
//   - [Message]
//   - [Radiobutton]
//   - [Spinbox]
//   - [TButton]
//   - [TCheckbutton]
//   - [TCombobox] (widget specific)
//   - [TEntry] (widget specific)
//   - [TLabel]
//   - [TMenubutton]
//   - [TRadiobutton]
func Textvariable(s string) Opt {
	return textVarOpt(s)
}

// Textvariable — Get the configured option value.
//
// Known uses:
//   - [Button]
//   - [Checkbutton]
//   - [Entry]
//   - [Label]
//   - [Menubutton]
//   - [Message]
//   - [Radiobutton]
//   - [Spinbox]
//   - [TButton]
//   - [TCheckbutton]
//   - [TCombobox] (widget specific)
//   - [TEntry] (widget specific)
//   - [TLabel]
//   - [TMenubutton]
//   - [TRadiobutton]
func (w *Window) Textvariable() (r string) {
	if tclVar := textVariables[w]; tclVar != "" {
		return evalErr(fmt.Sprintf("set %s", tclVar))
	}

	return ""
}

// Focus — Manage the input focus
//
// # Description
//
// The focus command is used to manage the Tk input focus. At any given time,
// one window on each display is designated as the focus window; any key press
// or key release events for the display are sent to that window. It is
// normally up to the window manager to redirect the focus among the top-level
// windows of a display. For example, some window managers automatically set
// the input focus to a top-level window whenever the mouse enters it; others
// redirect the input focus only when the user clicks on a window. Usually the
// window manager will set the focus only to top-level windows, leaving it up
// to the application to redirect the focus among the children of the
// top-level.
//
// Tk remembers one focus window for each top-level (the most recent descendant
// of that top-level to receive the focus); when the window manager gives the
// focus to a top-level, Tk automatically redirects it to the remembered
// window. Within a top-level Tk uses an explicit focus model by default.
// Moving the mouse within a top-level does not normally change the focus; the
// focus changes only when a widget decides explicitly to claim the focus
// (e.g., because of a button click), or when the user types a key such as Tab
// that moves the focus.
//
// The Tcl procedure tk_focusFollowsMouse may be invoked to create an implicit
// focus model: it reconfigures Tk so that the focus is set to a window
// whenever the mouse enters it. The Tcl procedures tk_focusNext and
// tk_focusPrev implement a focus order among the windows of a top-level; they
// are used in the default bindings for Tab and Shift-Tab, among other things.
//
// The focus command can take any of the following forms:
//
//	Focus()
//
// Returns the path name of the focus window on the display containing the
// application's main window, or an empty string if no window in this
// application has the focus on that display. Note: it is better to specify the
// display explicitly using -displayof (see below) so that the code will work
// in applications using multiple displays.
//
//	Focus(window)
//
// If the application currently has the input focus on window's display, this
// command resets the input focus for window's display to window and returns an
// empty string. If the application does not currently have the input focus on
// window's display, window will be remembered as the focus for its top-level;
// the next time the focus arrives at the top-level, Tk will redirect it to
// window. If window is an empty string then the command does nothing.
//
//	Focus(Displayof(window))
//
// Returns the name of the focus window on the display containing window. If
// the focus window for window's display is not in this application, the return
// value is an empty string.
//
//	Focus(Force(window))
//
// Sets the focus of window's display to window, even if the application does
// not currently have the input focus for the display. This command should be
// used sparingly, if at all. In normal usage, an application should not claim
// the focus for itself; instead, it should wait for the window manager to give
// it the focus. If window is an empty string then the command does nothing.
//
//	Focus(Lastfor(window))
//
// Returns the name of the most recent window to have the input focus among all
// the windows in the same top-level as window. If no window in that top-level
// has ever had the input focus, or if the most recent focus window has been
// deleted, then the name of the top-level is returned. The return value is the
// window that will receive the input focus the next time the window manager
// gives the focus to the top-level.
//
// # Quirks
//
// When an internal window receives the input focus, Tk does not actually set
// the X focus to that window; as far as X is concerned, the focus will stay on
// the top-level window containing the window with the focus. However, Tk
// generates FocusIn and FocusOut events just as if the X focus were on the
// internal window. This approach gets around a number of problems that would
// occur if the X focus were actually moved; the fact that the X focus is on
// the top-level is invisible unless you use C code to query the X server
// directly.
//
// More information might be available at the [Tcl/Tk focus] page.
//
// [Tcl/Tk focus]: https://www.tcl.tk/man/tcl9.0/TkCmd/focus.html
func Focus(options ...Opt) string {
	return evalErr(fmt.Sprintf("focus %s", collect(options...)))
}

// FontFace represents a Tk font.
type FontFace struct {
	name string
}

func (f *FontFace) optionString(_ *Window) (r string) {
	if f != nil {
		return f.name
	}

	return "font0" // does not exist
}

// String implements fmt.Stringer.
func (f *FontFace) String() string {
	return f.optionString(nil)
}

// NewFont — Create and inspect fonts.
//
// # Description
//
// Returns the amount in pixels that the tallest letter sticks up above the baseline
// of the font, plus any extra blank space added by the designer of the font.
//
// Additional information might be available at the [Tcl/Tk font] page.
//
// [Tcl/Tk font]: https://www.tcl.tk/man/tcl9.0/TkCmd/font.html
func (f *FontFace) MetricsAscent(window *Window) int {
	metric := evalErr(fmt.Sprintf("font metrics %s -displayof %s -ascent", f.name, window))
	return atoi(metric)
}

// NewFont — Create and inspect fonts.
//
// # Description
//
// Returns the largest amount in pixels that any letter sticks down below the baseline
// of the font, plus any extra blank space added by the designer of the font.
//
// Additional information might be available at the [Tcl/Tk font] page.
//
// [Tcl/Tk font]: https://www.tcl.tk/man/tcl9.0/TkCmd/font.html
func (f *FontFace) MetricsDescent(window *Window) int {
	metric := evalErr(fmt.Sprintf("font metrics %s -displayof %s -descent", f.name, window))
	return atoi(metric)
}

// NewFont — Create and inspect fonts.
//
// # Description
//
// Returns how far apart vertically in pixels two lines of text using the same font should
// be placed so that none of the characters in one line overlap any of the characters
// in the other line. This is generally the sum of the ascent above the baseline line plus
// the descent below the baseline.
//
// Additional information might be available at the [Tcl/Tk font] page.
//
// [Tcl/Tk font]: https://www.tcl.tk/man/tcl9.0/TkCmd/font.html
func (f *FontFace) MetricsLinespace(window *Window) int {
	metric := evalErr(fmt.Sprintf("font metrics %s -displayof %s -linespace", f.name, window))
	return atoi(metric)
}

// NewFont — Create and inspect fonts.
//
// # Description
//
// Returns a boolean flag that is true if this is a fixed-width font, where each normal character
// is the same width as all the other characters, or is false if this is a proportionally-spaced font,
// where individual characters have different widths. The widths of control characters, tab characters,
// and other non-printing characters are not included when calculating this value.
//
// Additional information might be available at the [Tcl/Tk font] page.
//
// [Tcl/Tk font]: https://www.tcl.tk/man/tcl9.0/TkCmd/font.html
func (f *FontFace) MetricsFixed(window *Window) bool {
	metric := evalErr(fmt.Sprintf("font metrics %s -displayof %s -fixed", f.name, window))
	return tclBool(metric)
}

// NewFont — Create and inspect fonts.
//
// # Description
//
// Creates a new font.
//
// The following options are supported on all platforms, and are used when
// creating/specifying a font:
//
//   - [Family] name
//
// The case-insensitive font family name. Tk guarantees to support the font
// families named Courier (a monospaced “typewriter” font), Times (a serifed
// “newspaper” font), and Helvetica (a sans-serif “European” font). The most
// closely matching native font family will automatically be substituted when
// one of the above font families is used. The name may also be the name of a
// native, platform-specific font family; in that case it will work as desired
// on one platform but may not display correctly on other platforms. If the
// family is unspecified or unrecognized, a platform-specific default font will
// be chosen.
//
//   - [Size] size
//
// The desired size of the font. If the size argument is a positive number, it
// is interpreted as a size in points. If size is a negative number, its
// absolute value is interpreted as a size in pixels. If a font cannot be
// displayed at the specified size, a nearby size will be chosen. If size is
// unspecified or zero, a platform-dependent default size will be chosen.
//
// Sizes should normally be specified in points so the application will remain
// the same ruler size on the screen, even when changing screen resolutions or
// moving scripts across platforms. However, specifying pixels is useful in
// certain circumstances such as when a piece of text must line up with respect
// to a fixed-size bitmap. The mapping between points and pixels is set when
// the application starts, based on properties of the installed monitor, but it
// can be overridden by calling the tk scaling command.
//
//   - [Weight] weight
//
// The nominal thickness of the characters in the font. The value normal
// specifies a normal weight font, while bold specifies a bold font. The
// closest available weight to the one specified will be chosen. The default
// weight is normal.
//
//   - [Slant] slant
//
// The amount the characters in the font are slanted away from the vertical.
// Valid values for slant are roman and italic. A roman font is the normal,
// upright appearance of a font, while an italic font is one that is tilted
// some number of degrees from upright. The closest available slant to the one
// specified will be chosen. The default slant is roman.
//
//   - [Underline] boolean
//
// The value is a boolean flag that specifies whether characters in this font
// should be underlined. The default value for underline is false.
//
//   - [Overstrike] boolean
//
// The value is a boolean flag that specifies whether a horizontal line should
// be drawn through the middle of characters in this font. The default value
// for overstrike is false.
//
// Additional information might be available at the [Tcl/Tk font] page.
//
// [Tcl/Tk font]: https://www.tcl.tk/man/tcl9.0/TkCmd/font.html
func NewFont(options ...Opt) *FontFace {
	nm := fmt.Sprintf("font%v", id.Add(1))
	code := ""
	configure := false
	if name := collectOne("-family", options...); name != "" &&
		strings.HasPrefix(name, "Tk") {
		code = "font actual " + name
		nm = name
		configure = true
	} else {
		code = fmt.Sprintf("font create %s %s", nm, collect(options...))
	}
	r, err := eval(code)
	if err == nil && configure {
		code := fmt.Sprintf("font configure %s %s", nm, collect(options...))
		r, err = eval(code)
	}
	if err != nil {
		fail(fmt.Errorf("code=%s -> r=%s err=%v", code, r, err))
		return nil
	}
	return &FontFace{name: nm}
}

// FontFamilies — Create and inspect fonts.
//
// # Description
//
// The return value is a list of the case-insensitive names of all font
// families that exist on window's display. If the Displayof argument is
// omitted, it defaults to the main window.
//
//   - [Displayof] window
//
// Additional information might be available at the [Tcl/Tk font] page.
//
// [Tcl/Tk font]: https://www.tcl.tk/man/tcl9.0/TkCmd/font.html
func FontFamilies(options ...Opt) []string {
	return parseList(evalErr(fmt.Sprintf("font families %s", collect(options...))))
}

func parseList(list string) (r []string) {
	if len(list) == 0 {
		return nil
	}

	cList, err := cString(list)
	if err != nil {
		fail(fmt.Errorf("failed to allocate C string for list : %w", err))
		return
	}
	defer free(cList)

	var _uintptr uintptr
	pointers, err := malloc(int(2 * unsafe.Sizeof(_uintptr)))
	if err != nil {
		fail(fmt.Errorf("failed to allocate memory for argc & argv pointers : %w", err))
		return
	}
	defer free(pointers)
	argcPtr := pointers
	argvPtr := pointers + unsafe.Sizeof(_uintptr)

	callSplitList(cList, argcPtr, argvPtr)
	argc := *((*int)(unsafe.Pointer(argcPtr)))
	argv := unsafe.Slice((*(**uintptr)(unsafe.Pointer(argvPtr))), argc)

	items := make([]string, argc)
	for i, arg := range argv {
		items[i] = goString(arg)
	}

	return items
}

// Delete — Manipulate fonts.
//
// # Description
//
// Delete the font. If there are widgets using the named font, the named font
// will not actually be deleted until all the instances are released. Those
// widgets will continue to display using the last known values for the named
// font. If a deleted named font is subsequently recreated with another call to
// font create, the widgets will use the new named font and redisplay
// themselves using the new attributes of that font.
//
// Additional information might be available at the [Tcl/Tk font] page.
//
// [Tcl/Tk font]: https://www.tcl.tk/man/tcl9.0/TkCmd/font.html
func (f *FontFace) Delete() {
	evalErr(fmt.Sprintf("font delete %s", f))
}

// Text — Create and manipulate 'text' hypertext editing widgets
//
// # Description
//
// Inserts all of the chars arguments just before the character at index. If
// index refers to the end of the text (the character after the last newline)
// then the new text is inserted just before the last newline instead. If there
// is a single chars argument and no tagList, then the new text will receive
// any tags that are present on both the character before and the character
// after the insertion point; if a tag is present on only one of these
// characters then it will not be applied to the new text. If tagList is
// specified then it consists of a list of tag names; the new characters will
// receive all of the tags in this list and no others, regardless of the tags
// present around the insertion point. If multiple chars-tagList argument pairs
// are present, they produce the same effect as if a separate pathName insert
// widget command had been issued for each pair, in order. The last tagList
// argument may be omitted.
//
// The value that is passed to Tcl/Tk for the 'index' argument is obtained by
// fmt.Sprint(index), enabling any custom index encoding via implementing
// fmt.Stringer.
//
// # Indices
//
// Many of the widget commands for texts take one or more indices as arguments.
// An index is a string used to indicate a particular place within a text, such
// as a place to insert characters or one endpoint of a range of characters to
// delete. Indices have the syntax
//
//	base modifier modifier modifier ...
//
// Where base gives a starting point and the modifiers adjust the index from
// the starting point (e.g. move forward or backward one character). Every
// index must contain a base, but the modifiers are optional. Most modifiers
// (as documented below) allow an optional submodifier. Valid submodifiers are
// any and display. If the submodifier is abbreviated, then it must be followed
// by whitespace, but otherwise there need be no space between the submodifier
// and the following modifier. Typically the display submodifier adjusts the
// meaning of the following modifier to make it refer to visual or non-elided
// units rather than logical units, but this is explained for each relevant
// case below. Lastly, where count is used as part of a modifier, it can be
// positive or negative, so “base - -3 lines” is perfectly valid (and
// equivalent to “base +3lines”).
//
// The base for an index must have one of the following forms:
//
//	"line.char"
//
// Indicates char'th character on line line. Lines are numbered from 1 for
// consistency with other UNIX programs that use this numbering scheme. Within
// a line, characters are numbered from 0. If char is end then it refers to the
// newline character that ends the line.
//
// This form of index can be passed as [LC]{line, char}.
//
//	"@x,y"
//
// Indicates the character that covers the pixel whose x and y coordinates
// within the text's window are x and y.
//
//	"end"
//
// Indicates the end of the text (the character just after the last newline).
//
//	"mark"
//
// Indicates the character just after the mark whose name is mark (see MARKS
// for details).
//
//	"tag.first"
//
// Indicates the first character in the text that has been tagged with tag.
// This form generates an error if no characters are currently tagged with tag.
//
//	"tag.last"
//
// Indicates the character just after the last one in the text that has been
// tagged with tag. This form generates an error if no characters are currently
// tagged with tag.
//
//	"pathName"
//
// Indicates the position of the embedded window whose name is pathName. This
// form generates an error if there is no embedded window by the given name.
//
//	"imageName"
//
// Indicates the position of the embedded image whose name is imageName. This
// form generates an error if there is no embedded image by the given name.
//
// If the base could match more than one of the above forms, such as a mark and
// imageName both having the same value, then the form earlier in the above
// list takes precedence. If modifiers follow the base index, each one of them
// must have one of the forms listed below. Keywords such as chars and wordend
// may be abbreviated as long as the abbreviation is unambiguous.
//
//	"+ count ?submodifier? chars"
//
// Adjust the index forward by count characters, moving to later lines in the
// text if necessary. If there are fewer than count characters in the text
// after the current index, then set the index to the last index in the text.
// Spaces on either side of count are optional. If the display submodifier is
// given, elided characters are skipped over without being counted. If any is
// given, then all characters are counted. For historical reasons, if neither
// modifier is given then the count actually takes place in units of index
// positions (see INDICES for details). This behaviour may be changed in a
// future major release, so if you need an index count, you are encouraged to
// use indices instead wherever possible.
//
//	"- count ?submodifier? chars"
//
// Adjust the index backward by count characters, moving to earlier lines in
// the text if necessary. If there are fewer than count characters in the text
// before the current index, then set the index to the first index in the text
// (1.0). Spaces on either side of count are optional. If the display
// submodifier is given, elided characters are skipped over without being
// counted. If any is given, then all characters are counted. For historical
// reasons, if neither modifier is given then the count actually takes place in
// units of index positions (see INDICES for details). This behavior may be
// changed in a future major release, so if you need an index count, you are
// encouraged to use indices instead wherever possible.
//
//	"+ count ?submodifier? indices"
//
// Adjust the index forward by count index positions, moving to later lines in
// the text if necessary. If there are fewer than count index positions in the
// text after the current index, then set the index to the last index position
// in the text. Spaces on either side of count are optional. Note that an index
// position is either a single character or a single embedded image or embedded
// window. If the display submodifier is given, elided indices are skipped over
// without being counted. If any is given, then all indices are counted; this
// is also the default behaviour if no modifier is given.
//
//	"- count ?submodifier? indices"
//
// Adjust the index backward by count index positions, moving to earlier lines
// in the text if necessary. If there are fewer than count index positions in
// the text before the current index, then set the index to the first index
// position (1.0) in the text. Spaces on either side of count are optional. If
// the display submodifier is given, elided indices are skipped over without
// being counted. If any is given, then all indices are counted; this is also
// the default behaviour if no modifier is given.
//
//	"+ count ?submodifier? lines"
//
// Adjust the index forward by count lines, retaining the same character
// position within the line. If there are fewer than count lines after the line
// containing the current index, then set the index to refer to the same
// character position on the last line of the text. Then, if the line is not
// long enough to contain a character at the indicated character position,
// adjust the character position to refer to the last character of the line
// (the newline). Spaces on either side of count are optional. If the display
// submodifier is given, then each visual display line is counted separately.
// Otherwise, if any (or no modifier) is given, then each logical line (no
// matter how many times it is visually wrapped) counts just once. If the
// relevant lines are not wrapped, then these two methods of counting are
// equivalent.
//
//	"- count ?submodifier? lines"
//
// Adjust the index backward by count logical lines, retaining the same
// character position within the line. If there are fewer than count lines
// before the line containing the current index, then set the index to refer to
// the same character position on the first line of the text. Then, if the line
// is not long enough to contain a character at the indicated character
// position, adjust the character position to refer to the last character of
// the line (the newline). Spaces on either side of count are optional. If the
// display submodifier is given, then each visual display line is counted
// separately. Otherwise, if any (or no modifier) is given, then each logical
// line (no matter how many times it is visually wrapped) counts just once. If
// the relevant lines are not wrapped, then these two methods of counting are
// equivalent.
//
//	"?submodifier? linestart"
//
// Adjust the index to refer to the first index on the line. If the display
// submodifier is given, this is the first index on the display line, otherwise
// on the logical line.
//
//	"?submodifier? lineend"
//
// Adjust the index to refer to the last index on the line (the newline). If
// the display submodifier is given, this is the last index on the display
// line, otherwise on the logical line.
//
//	"?submodifier? wordstart"
//
// Adjust the index to refer to the first character of the word containing the
// current index. A word consists of any number of adjacent characters that are
// letters, digits, or underscores, or a single character that is not one of
// these. If the display submodifier is given, this only examines non-elided
// characters, otherwise all characters (elided or not) are examined.
//
//	"?submodifier? wordend"
//
// Adjust the index to refer to the character just after the last one of the
// word containing the current index. If the current index refers to the last
// character of the text then it is not modified. If the display submodifier is
// given, this only examines non-elided characters, otherwise all characters
// (elided or not) are examined.
//
// If more than one modifier is present then they are applied in left-to-right
// order. For example, the index “end - 1 chars” refers to the next-to-last
// character in the text and “insert wordstart - 1 c” refers to the character
// just before the first one in the word containing the insertion cursor.
// Modifiers are applied one by one in this left to right order, and after each
// step the resulting index is constrained to be a valid index in the text
// widget. So, for example, the index “1.0 -1c +1c” refers to the index “2.0”.
//
// Where modifiers result in index changes by display lines, display chars or
// display indices, and the base refers to an index inside an elided tag, that
// base index is considered to be equivalent to the first following non-elided
// index.
//
// Insert returns its index argument.
//
// Additional information might be available at the [Tcl/Tk text] page.
//
// [Tcl/Tk text]: https://www.tcl.tk/man/tcl9.0/TkCmd/text.html
func (w *TextWidget) Insert(index any, chars string, options ...string) any {
	idx := fmt.Sprint(index)
	evalErr(fmt.Sprintf("%s insert %s %s %s", w, tclSafeString(idx), tclSafeString(chars), tclSafeStrings(options...)))
	return index
}

// Text — Create and manipulate 'text' hypertext editing widgets
//
// # Description
//
// This command associates a script with the tag given by tagName. Whenever the
// event sequence given by sequence occurs for a character that has been tagged
// with tagName, the script will be invoked. This widget command is similar to
// the bind command except that it operates on characters in a text rather than
// entire widgets. See the bind manual entry for complete details on the syntax
// of sequence and the substitutions performed on script before invoking it.
// A new binding is created, replacing any existing binding for the same sequence and tagName.
// The only events for which
// bindings may be specified are those related to the mouse and keyboard (such
// as Enter, Leave, Button, Motion, and Key) or virtual events. Mouse and
// keyboard event bindings for a text widget respectively use the current and
// insert marks described under MARKS above. An Enter event triggers for a tag
// when the tag first becomes present on the current character, and a Leave
// event triggers for a tag when it ceases to be present on the current
// character. Enter and Leave events can happen either because the current mark
// moved or because the character at that position changed. Note that these
// events are different than Enter and Leave events for windows. Mouse events
// are directed to the current character, while keyboard events are directed to
// the insert character. If a virtual event is used in a binding, that binding
// can trigger only if the virtual event is defined by an underlying
// mouse-related or keyboard-related event.
//
// It is possible for the current character to have multiple tags, and for each
// of them to have a binding for a particular event sequence. When this occurs,
// one binding is invoked for each tag, in order from lowest-priority to
// highest priority. If there are multiple matching bindings for a single tag,
// then the most specific binding is chosen (see the manual entry for the bind
// command for details). continue and break commands within binding scripts are
// processed in the same way as for bindings created with the bind command.
//
// If bindings are created for the widget as a whole using the bind command,
// then those bindings will supplement the tag bindings. The tag bindings will
// be invoked first, followed by bindings for the window as a whole.
//
// Additional information might be available at the [Tcl/Tk text] page.
//
// [Tcl/Tk text]: https://www.tcl.tk/man/tcl9.0/TkCmd/text.html
func (w *TextWidget) TagBind(tag, sequence string, handler any) string {
	return evalErr(fmt.Sprintf("%s tag bind %s %s %s", w, tclSafeString(tag), tclSafeString(sequence), newEventHandler("", handler).optionString(w.Window)))
}

// Text — Create and manipulate 'text' hypertext editing widgets
//
// # Description
//
// Counts the number of relevant things between the two indices. If index1 is
// after index2, the result will be a negative number (and this holds for each
// of the possible options). The actual items which are counted depend on the
// options given. The result is a list of integers, one for the result of each
// counting option given. Valid counting options are -chars, -displaychars,
// -displayindices, -displaylines, -indices, -lines, -xpixels and -ypixels. The
// default value, if no option is specified, is -indices. There is an
// additional possible option -update which is a modifier. If given (and if the
// text widget is managed by a geometry manager), then all subsequent options
// ensure that any possible out of date information is recalculated. This
// currently only has any effect for the -ypixels count (which, if -update is
// not given, will use the text widget's current cached value for each line).
// This -update option is obsoleted by pathName sync, pathName pendingsync and
// <<WidgetViewSync>>. The count options are interpreted as follows:
//
//   - [Chars] Count all characters, whether elided or not. Do not count
//     embedded windows or images.
//   - [Displaychars] Count all non-elided characters.
//   - [Displayindices] Count all non-elided characters, windows and images.
//   - [Displaylines] Count all display lines (i.e. counting one for each time
//     a line wraps) from the line of the first index up to, but not including
//     the display line of the second index. Therefore if they are both on the
//     same display line, zero will be returned. By definition displaylines are
//     visible and therefore this only counts portions of actual visible lines.
//   - [Indices] Count all characters and embedded windows or images (i.e.
//     everything which counts in text-widget index space), whether they are
//     elided or not.
//   - [Lines] Count all logical lines (irrespective of wrapping) from the line
//     of the first index up to, but not including the line of the second index.
//     Therefore if they are both on the same line, zero will be returned.
//     Logical lines are counted whether they are currently visible (non-elided)
//     or not.
//   - [Xpixels] Count the number of horizontal pixels from the first pixel of
//     the first index to (but not including) the first pixel of the second
//     index. To count the total desired width of the text widget (assuming
//     wrapping is not enabled), first find the longest line and then use “.text
//     count -xpixels "${line}.0" "${line}.0 lineend"”.
//   - [Ypixels] Count the number of vertical pixels from the first pixel of
//     the first index to (but not including) the first pixel of the second
//     index. If both indices are on the same display line, zero will be
//     returned. To count the total number of vertical pixels in the text widget,
//     use “.text count -ypixels 1.0 end”, and to ensure this is up to date, use
//     “.text count -update -ypixels 1.0 end”.
//
// The command returns a positive or negative integer corresponding to the
// number of items counted between the two indices. One such integer is
// returned for each counting option given, so a list is returned if more than
// one option was supplied. For example “.text count -xpixels -ypixels 1.3 4.5”
// is perfectly valid and will return a list of two elements.
//
// Additional information might be available at the [Tcl/Tk text] page.
//
// [Tcl/Tk text]: https://www.tcl.tk/man/tcl9.0/TkCmd/text.html
func (w *TextWidget) Count(options ...any) []string {
	return parseList(evalErr(fmt.Sprintf("%s count %s", w, collectAny(options...))))
}

// All option.
//
// Known uses:
//   - [TextWidget] (command specific)
func All() Opt {
	return rawOption("-all")
}

// Displayindices option.
//
// Known uses:
//   - [TextWidget] (command specific)
func Displayindices() Opt {
	return rawOption("-displayindices")
}

// Displaylines option.
//
// Known uses:
//   - [TextWidget] (command specific)
func Displaylines() Opt {
	return rawOption("-displaylines")
}

// Indices option.
//
// Known uses:
//   - [TextWidget] (command specific)
func Indices() Opt {
	return rawOption("-indices")
}

// Lines option.
//
// Known uses:
//   - [TextWidget] (command specific)
func Lines() Opt {
	return rawOption("-lines")
}

// Mark option.
//
// Known uses:
//   - [TextWidget] (command specific)
func Mark() Opt {
	return rawOption("-mark")
}

// Xpixels option.
//
// Known uses:
//   - [TextWidget] (command specific)
func Xpixels() Opt {
	return rawOption("-xpixels")
}

// Ypixels option.
//
// Known uses:
//   - [TextWidget] (command specific)
func Ypixels() Opt {
	return rawOption("-ypixels")
}

// Chars option.
//
// Known uses:
//   - [TextWidget] (command specific)
func Chars() Opt {
	return rawOption("-chars")
}

// Displaychars option.
//
// Known uses:
//   - [TextWidget] (command specific)
func Displaychars() Opt {
	return rawOption("-displaychars")
}

// Text — Create and manipulate 'text' hypertext editing widgets
//
// # Description
//
// Clears the undo and redo stacks.
//
// Additional information might be available at the [Tcl/Tk text] page.
//
// [Tcl/Tk text]: https://www.tcl.tk/man/tcl9.0/TkCmd/text.html
func (w *TextWidget) EditReset() {
	evalErr(fmt.Sprintf("%s edit reset", w))
}

// Text — Create and manipulate 'text' hypertext editing widgets
//
// # Description
//
// Sets the modified flag of the widget to 'v'.
//
// Additional information might be available at the [Tcl/Tk text] page.
//
// [Tcl/Tk text]: https://www.tcl.tk/man/tcl9.0/TkCmd/text.html
func (w *TextWidget) SetModified(v bool) {
	evalErr(fmt.Sprintf("%s edit modified %v", w, v))
}

// Text — Create and manipulate 'text' hypertext editing widgets
//
// # Description
//
// Returns the modified flag of the widget. The insert, delete, edit undo and
// edit redo commands or the user can set or clear the modified flag.
//
// Additional information might be available at the [Tcl/Tk text] page.
//
// [Tcl/Tk text]: https://www.tcl.tk/man/tcl9.0/TkCmd/text.html
func (w *TextWidget) Modified() bool {
	return tclBool(evalErr(fmt.Sprintf("%s edit modified", w)))
}

func tclBool(s string) bool {
	switch s {
	case "1", "true", "yes":
		return true
	case "0", "false", "no":
		return false
	default:
		fail(fmt.Errorf("unexpected Tcl bool: %q", s))
		return false
	}
}

// Text — Create and manipulate 'text' hypertext editing widgets
//
// # Description
//
// Returns a list whose elements are the names of all the tags that are active
// at the character position given by index. If index is omitted, then the
// return value will describe all of the tags that exist for the text (this
// includes all tags that have been named in a “pathName tag” widget command
// but have not been deleted by a “pathName tag delete” widget command, even if
// no characters are currently marked with the tag). The list will be sorted in
// order from lowest priority to highest priority.
//
// Additional information might be available at the [Tcl/Tk text] page.
//
// [Tcl/Tk text]: https://www.tcl.tk/man/tcl9.0/TkCmd/text.html
func (w *TextWidget) TagNames(index string) []string {
	if index != "" {
		index = tclSafeString(index)
	}
	return parseList(evalErr(fmt.Sprintf("%s tag names %s", w, index)))
}

// Text — Create and manipulate 'text' hypertext editing widgets
//
// # Description
//
// If direction is not specified, returns left or right to indicate which of
// its adjacent characters markName is attached to. If direction is specified,
// it must be left or right; the gravity of markName is set to the given value.
//
// Additional information might be available at the [Tcl/Tk text] page.
//
// [Tcl/Tk text]: https://www.tcl.tk/man/tcl9.0/TkCmd/text.html
func (w *TextWidget) MarkGravity(markName, direction string) {
	evalErr(fmt.Sprintf("%s mark gravity %s %s", w, tclSafeString(markName), tclSafeString(direction)))
}

// Text — Create and manipulate 'text' hypertext editing widgets
//
// # Description
//
// Sets the mark named markName to a position just before the character at
// index. If markName already exists, it is moved from its old position; if it
// does not exist, a new mark is created. This command returns an empty string.
//
// Additional information might be available at the [Tcl/Tk text] page.
//
// [Tcl/Tk text]: https://www.tcl.tk/man/tcl9.0/TkCmd/text.html
func (w *TextWidget) MarkSet(markName string, index any) {
	evalErr(fmt.Sprintf("%s mark set %s %s", w, tclSafeString(markName), tclSafeString(fmt.Sprint(index))))
}

// Text — Create and manipulate 'text' hypertext editing widgets
//
// # Description
//
// Remove the mark corresponding to each of the markName arguments. The removed
// marks will not be usable in indices and will not be returned by future calls
// to “pathName mark names”. This command returns an empty string.
//
// Additional information might be available at the [Tcl/Tk text] page.
//
// [Tcl/Tk text]: https://www.tcl.tk/man/tcl9.0/TkCmd/text.html
func (w *TextWidget) MarkUnset(markName ...string) {
	evalErr(fmt.Sprintf("%s mark unset %s", w, tclSafeStrings(markName...)))
}

// Text — Create and manipulate 'text' hypertext editing widgets
//
// # Description
//
// Returns a list whose elements are the names of all the marks that are
// currently set.
//
// Additional information might be available at the [Tcl/Tk text] page.
//
// [Tcl/Tk text]: https://www.tcl.tk/man/tcl9.0/TkCmd/text.html
func (w *TextWidget) MarkNames() []string {
	return parseList(evalErr(fmt.Sprintf("%s mark names", w)))
}

// Text — Create and manipulate 'text' hypertext editing widgets
//
// # Description
//
// Searches the text in pathName starting at index for a range of characters
// that matches pattern. If a match is found, the index of the first character
// in the match is returned as result; otherwise an empty string is returned.
// One or more of the following switches (or abbreviations thereof) may be
// specified to control the search:
//
// Additional information might be available at the [Tcl/Tk text] page.
//
// [Tcl/Tk text]: https://www.tcl.tk/man/tcl9.0/TkCmd/text.html
func (w *TextWidget) Search(options ...any) string {
	var switches Opts
	var args []any
	for _, v := range options {
		switch x := v.(type) {
		case Opt:
			switches = append(switches, x)
		default:
			args = append(args, v)
		}
	}
	return evalErr(fmt.Sprintf("%s search %s %s", w, collect(switches...), tclSafeList(args...)))
}

// Forward option.
//
// Known uses:
//   - [Text] (widget specific, applies to Search)
func Forward() Opt {
	return rawOption(fmt.Sprintf(`-forward`))
}

// Backward option.
//
// Known uses:
//   - [Text] (widget specific, applies to Search)
func Backward() Opt {
	return rawOption(fmt.Sprintf(`-backward`))
}

// Regexp option.
//
// Known uses:
//   - [Text] (widget specific, applies to Search)
func Regexp() Opt {
	return rawOption(fmt.Sprintf(`-regexp`))
}

// Nocase option.
//
// Known uses:
//   - [Text] (widget specific, applies to Search)
func Nocase() Opt {
	return rawOption(fmt.Sprintf(`-nocase`))
}

// Text — Create and manipulate 'text' hypertext editing widgets
//
// # Description
//
// Deletes all tag information for each of the tagName arguments. The command
// removes the tags from all characters in the file and also deletes any other
// information associated with the tags, such as bindings and display
// information. The command returns an empty string.
//
// Additional information might be available at the [Tcl/Tk text] page.
//
// [Tcl/Tk text]: https://www.tcl.tk/man/tcl9.0/TkCmd/text.html
func (w *TextWidget) TagDelete(tags ...string) {
	evalErr(fmt.Sprintf("%s tag delete %s", w, tclSafeStrings(tags...)))
}

// Text — Create and manipulate 'text' hypertext editing widgets
//
// # Description
//
// Clear makes 'w' empty.
func (w *TextWidget) Clear() {
	w.Delete("0.0", "end")
	w.TagDelete(w.TagNames("")...)
	w.MarkUnset(w.MarkNames()...)
}

// Text — Create and manipulate 'text' hypertext editing widgets
//
// # Description
//
// Copies the selection in the widget to the clipboard, if there is a selection.
//
// Additional information might be available at the [Tcl/Tk text] page.
//
// [Tcl/Tk text]: https://www.tcl.tk/man/tcl9.0/TkCmd/text.html
func (w *TextWidget) Copy() {
	evalErr(fmt.Sprintf("tk_textCopy %s", w))
}

// Text — Create and manipulate 'text' hypertext editing widgets
//
// # Description
//
// Copies the selection in the widget to the clipboard and deletes the
// selection. If there is no selection in the widget then these keys have no
// effect.
//
// Additional information might be available at the [Tcl/Tk text] page.
//
// [Tcl/Tk text]: https://www.tcl.tk/man/tcl9.0/TkCmd/text.html
func (w *TextWidget) Cut() {
	evalErr(fmt.Sprintf("tk_textCut %s", w))
}

// Text — Create and manipulate 'text' hypertext editing widgets
//
// # Description
//
// Inserts the contents of the clipboard at the position of the insertion cursor.
//
// Additional information might be available at the [Tcl/Tk text] page.
//
// [Tcl/Tk text]: https://www.tcl.tk/man/tcl9.0/TkCmd/text.html
func (w *TextWidget) Paste() {
	evalErr(fmt.Sprintf("tk_textPaste %s", w))
}

// Text — Create and manipulate 'text' hypertext editing widgets
//
// # Description
//
// Undoes the last edit action when the -undo option is true, and returns a
// list of indices indicating what ranges were changed by the undo operation.
// An edit action is defined as all the insert and delete commands that are
// recorded on the undo stack in between two separators. Generates an error
// when the undo stack is empty. Does nothing when the -undo option is false.
//
// Additional information might be available at the [Tcl/Tk text] page.
//
// [Tcl/Tk text]: https://www.tcl.tk/man/tcl9.0/TkCmd/text.html
func (w *TextWidget) Undo() {
	evalErr(fmt.Sprintf("%s edit undo", w))
}

// Text — Create and manipulate 'text' hypertext editing widgets
//
// # Description
//
// When the -undo option is true, reapplies the last undone edits provided no
// other edits were done since then, and returns a list of indices indicating
// what ranges were changed by the redo operation. Generates an error when the
// redo stack is empty. Does nothing when the -undo option is false.
//
// Additional information might be available at the [Tcl/Tk text] page.
//
// [Tcl/Tk text]: https://www.tcl.tk/man/tcl9.0/TkCmd/text.html
func (w *TextWidget) Redo() {
	evalErr(fmt.Sprintf("%s edit redo", w))
}

// Text — Create and manipulate 'text' hypertext editing widgets
//
// # Description
//
// Return a range of characters from the text. The return value will be all the
// characters in the text starting with the one whose index is index1 and
// ending just before the one whose index is index2 (the character at index2
// will not be returned). If index2 is omitted then the single character at
// index1 is returned. If there are no characters in the specified range (e.g.
// index1 is past the end of the file or index2 is less than or equal to
// index1) then an empty string is returned. If the specified range contains
// embedded windows, no information about them is included in the returned
// string. If multiple index pairs are given, multiple ranges of text will be
// returned in a list. Invalid ranges will not be represented with empty
// strings in the list. The ranges are returned in the order passed to pathName
// get. If the -displaychars option is given, then, within each range, only
// those characters which are not elided will be returned. This may have the
// effect that some of the returned ranges are empty strings.
//
// BUG(jnml) [TextWidget.Get] currently supports only one range.
//
// Additional information might be available at the [Tcl/Tk text] page.
//
// [Tcl/Tk text]: https://www.tcl.tk/man/tcl9.0/TkCmd/text.html
func (w *TextWidget) Get(options ...any) (r []string) {
	return []string{evalErr(fmt.Sprintf("%s get %s", w, collectAny(options...)))}
}

// Text is a shortcut for w.Get("1.0", "end-1c")[0].
func (w *TextWidget) Text() string {
	return w.Get("1.0", "end-1c")[0]
}

// Text — Create and manipulate 'text' hypertext editing widgets
//
// # Description
//
// Replaces the range of characters between index1 and index2 with the given
// characters and tags. See the section on pathName insert for an explanation
// of the handling of the tagList... arguments, and the section on pathName delete
// for an explanation of the handling of the indices. If index2 corresponds to
// an index earlier in the text than index1, an error will be generated.
//
// The deletion and insertion are arranged so that no unnecessary scrolling of
// the window or movement of insertion cursor occurs. In addition the undo/redo
// stack are correctly modified, if undo operations are active in the text widget.
// The command returns an empty string.
//
// Additional information might be available at the [Tcl/Tk text] page.
//
// [Tcl/Tk text]: https://www.tcl.tk/man/tcl9.0/TkCmd/text.html
func (w *TextWidget) Replace(index1 any, index2 any, chars string, options ...any) string {
	return evalErr(fmt.Sprintf("%s replace %s %s %s %s", w,
		tclSafeString(fmt.Sprint(index1)),
		tclSafeString(fmt.Sprint(index2)),
		tclSafeString(chars),
		collectAny(options...),
	))
}

// Text — Create and manipulate 'text' hypertext editing widgets
//
// # Description
//
// Return the contents of the text widget from index1 up to, but not including
// index2, including the text and information about marks, tags, and embedded
// windows. If index2 is not specified, then it defaults to one character past
// index1. The information is returned in the following format:
//
// key1 value1 index1 key2 value2 index2 ...
//
// The possible key values are text, mark, tagon, tagoff, image, and window.
// The corresponding value is the text, mark name, tag name, image name, or
// window name. The index information is the index of the start of the text,
// mark, tag transition, image or window. One or more of the following switches
// (or abbreviations thereof) may be specified to control the dump:
//
//   - [All] Return information about all elements: text, marks, tags, images and windows.
//     This is the default.
//   - [Command command] Instead of returning the information as the result of the
//     dump operation, invoke the command on each element of the text widget within the range. The command has three arguments appended to it before it is evaluated: the key, value, and index.
//   - [Image] Include information about images in the dump results.
//   - [Mark] Include information about marks in the dump results.
//   - [Tag] Include information about tag transitions in the dump results. Tag
//     information is returned as tagon and tagoff elements that indicate the begin
//     and end of each range of each tag, respectively.
//   - [Text] Include information about text in the dump results. The value is the
//     text up to the next element or the end of range indicated by index2. A text
//     element does not span newlines. A multi-line block of text that contains no
//     marks or tag transitions will still be dumped as a set of text segments that
//     each end with a newline. The newline is part of the value.
//   - [Window] Include information about embedded windows in the dump results. The
//     value of a window is its Tk pathname, unless the window has not been created yet.
//     (It must have a create script.) In this case an empty string is returned, and you
//     must query the window by its index position to get more information.
//
// Additional information might be available at the [Tcl/Tk text] page.
//
// [Tcl/Tk text]: https://www.tcl.tk/man/tcl9.0/TkCmd/text.html
func (w *TextWidget) Dump(options ...any) []string {
	return parseList(evalErr(fmt.Sprintf("%s dump %s", w, collectAny(options...))))
}

// LC encodes a text index consisting of a line and char number.
type LC struct {
	Line int // 1-based line number within the text content.
	Char int // 0-based char number within the line.
}

// String implements fmt.Stringer.
func (lc LC) String() string {
	return fmt.Sprintf("%d.%d", lc.Line, lc.Char)
}

// Text — Create and manipulate 'text' hypertext editing widgets
//
// # Description
//
// This command creates a new image annotation, which will appear in the text
// at the position given by index. Any number of option-value pairs may be
// specified to configure the annotation. Returns a unique identifier that may
// be used as an index to refer to this image. See EMBEDDED IMAGES for
// information on the options that are supported, and a description of the
// identifier returned.
//
// Additional information might be available at the [Tcl/Tk text] page.
//
// [Tcl/Tk text]: https://www.tcl.tk/man/tcl9.0/TkCmd/text.html
func (w *TextWidget) ImageCreate(index any, options ...Opt) {
	idx := fmt.Sprint(index)
	evalErr(fmt.Sprintf("%s image create %s %s", w, tclSafeString(idx), collect(options...)))
}

// Text — Create and manipulate 'text' hypertext editing widgets
//
// # Description
//
// This command creates a new window annotation, which will appear in the text
// at the position given by index. Any number of option-value pairs may be
// specified to configure the annotation. See EMBEDDED WINDOWS for information
// on the options that are supported. Returns an empty string.
//
// Additional information might be available at the [Tcl/Tk text] page.
//
// [Tcl/Tk text]: https://www.tcl.tk/man/tcl9.0/TkCmd/text.html
func (w *TextWidget) WindowCreate(index any, options ...Opt) {
	idx := fmt.Sprint(index)
	evalErr(fmt.Sprintf("%s window create %s %s", w, tclSafeString(idx), collect(options...)))
}

// Win option.
//
// Known uses:
//   - [Text] (widget specific, applies to embedded windows)
func Win(val any) Opt {
	return rawOption(fmt.Sprintf(`-window %s`, optionString(val)))
}

// Text — Create and manipulate 'text' hypertext editing widgets
//
// # Description
//
// Returns a list containing two elements, both of which are real fractions
// between 0 and 1. The first element gives the position of the first visible
// pixel of the first character (or image, etc) in the top line in the window,
// relative to the text as a whole (0.5 means it is halfway through the text,
// for example). The second element gives the position of the first pixel just
// after the last visible one in the bottom line of the window, relative to the
// text as a whole. These are the same values passed to scrollbars via the
// -yscrollcommand option.
//
// Additional information might be available at the [Tcl/Tk text] page.
//
// [Tcl/Tk text]: https://www.tcl.tk/man/tcl9.0/TkCmd/text.html
func (w *TextWidget) Yview() string {
	return evalErr(fmt.Sprintf("%s yview", w))
}

// Text — Create and manipulate 'text' hypertext editing widgets
//
// # Description
//
// Adjusts the view in the window so that the pixel given by fraction appears
// at the top of the top line of the window. Fraction is a fraction between 0
// and 1; 0 indicates the first pixel of the first character in the text, 0.33
// indicates the pixel that is one-third the way through the text; and so on.
// Values close to 1 will indicate values close to the last pixel in the text
// (1 actually refers to one pixel beyond the last pixel), but in such cases
// the widget will never scroll beyond the last pixel, and so a value of 1 will
// effectively be rounded back to whatever fraction ensures the last pixel is
// at the bottom of the window, and some other pixel is at the top.
//
// Additional information might be available at the [Tcl/Tk text] page.
//
// [Tcl/Tk text]: https://www.tcl.tk/man/tcl9.0/TkCmd/text.html
func (w *TextWidget) Yviewmoveto(fraction any) string {
	return evalErr(fmt.Sprintf("%s yview moveto %s", w, tclSafeString(fmt.Sprint(fraction))))
}

// Text — Create and manipulate 'text' hypertext editing widgets
//
// # Description
//
// Adjusts the view in the window so that the character given by index is
// completely visible. If index is already visible then the command does
// nothing. If index is a short distance out of view, the command adjusts the
// view just enough to make index visible at the edge of the window. If index
// is far out of view, then the command centers index in the window.
//
// Additional information might be available at the [Tcl/Tk text] page.
//
// [Tcl/Tk text]: https://www.tcl.tk/man/tcl9.0/TkCmd/text.html
func (w *TextWidget) See(index any) {
	evalErr(fmt.Sprintf("%s see %s", w, tclSafeString(fmt.Sprint(index))))
}

// Text — Create and manipulate 'text' hypertext editing widgets
//
// # Description
//
// Returns a list describing all of the ranges of text that have been tagged
// with tagName. The first two elements of the list describe the first tagged
// range in the text, the next two elements describe the second range, and so
// on. The first element of each pair contains the index of the first character
// of the range, and the second element of the pair contains the index of the
// character just after the last one in the range. If there are no characters
// tagged with tag then an empty string is returned.
//
// Additional information might be available at the [Tcl/Tk text] page.
//
// [Tcl/Tk text]: https://www.tcl.tk/man/tcl9.0/TkCmd/text.html
func (w *TextWidget) TagRanges(tagName string) (r []string) {
	return parseList(evalErr(fmt.Sprintf("%s tag ranges %s", w, tclSafeString(tagName))))
}

// Text — Create and manipulate 'text' hypertext editing widgets
//
// # Description
//
// Returns the position corresponding to index in the form line.char where line
// is the line number and char is the character number. Index may have any of
// the forms described under INDICES above.
//
// Additional information might be available at the [Tcl/Tk text] page.
//
// [Tcl/Tk text]: https://www.tcl.tk/man/tcl9.0/TkCmd/text.html
func (w *TextWidget) Index(index any) (r string) {
	return evalErr(fmt.Sprintf("%s index %s", w, tclSafeString(fmt.Sprint(index))))
}

// Text — Create and manipulate 'text' hypertext editing widgets
//
// # Description
//
// Returns a list containing two elements, both of which are real fractions
// between 0 and 1. The first element gives the position of the first visible
// pixel of the first character (or image, etc) in the top line in the window,
// relative to the text as a whole (0.5 means it is halfway through the text,
// for example). The second element gives the position of the first pixel just
// after the last visible one in the bottom line of the window, relative to the
// text as a whole. These are the same values passed to scrollbars via the
// -xscrollcommand option.
//
// Additional information might be available at the [Tcl/Tk text] page.
//
// [Tcl/Tk text]: https://www.tcl.tk/man/tcl9.0/TkCmd/text.html
func (w *TextWidget) Xview() string {
	return evalErr(fmt.Sprintf("%s xview", w))
}

// Text — Create and manipulate 'text' hypertext editing widgets
//
// # Description
//
// This command is similar to the pathName configure widget command except that
// it modifies options associated with the tag given by tagName instead of
// modifying options for the overall text widget. If no option is specified,
// the command returns a list describing all of the available options for
// tagName (see Tk_ConfigureInfo for information on the format of this list).
// If option is specified with no value, then the command returns a list
// describing the one named option (this list will be identical to the
// corresponding sublist of the value returned if no option is specified). If
// one or more option-value pairs are specified, then the command modifies the
// given option(s) to have the given value(s) in tagName; in this case the
// command returns an empty string. See TAGS above for details on the options
// available for tags.
//
// Additional information might be available at the [Tcl/Tk text] page.
//
// [Tcl/Tk text]: https://www.tcl.tk/man/tcl9.0/TkCmd/text.html
func (w *TextWidget) TagConfigure(tagName string, options ...Opt) {
	evalErr(fmt.Sprintf("%s tag configure %s %s", w, tclSafeString(tagName), collect(options...)))
}

// Text — Create and manipulate 'text' hypertext editing widgets
//
// # Description
//
// Select all text in 'w'.
func (w *TextWidget) SelectAll() {
	evalErr(fmt.Sprintf("%s tag add sel 1.0 end", w))
}

// Text — Create and manipulate 'text' hypertext editing widgets
//
// # Description
//
// Associate the tag name with all of the characters starting with index1
// and ending just before index2 (the character at index2 is not tagged). A
// single command may contain any number of index1-index2 pairs. If the last
// index2 is omitted then the single character at index1 is tagged. If there
// are no characters in the specified range (e.g. index1 is past the end of the
// file or index2 is less than or equal to index1) then the command has no
// effect.
//
// Additional information might be available at the [Tcl/Tk text] page.
//
// [Tcl/Tk text]: https://www.tcl.tk/man/tcl9.0/TkCmd/text.html
func (w *TextWidget) TagAdd(tagName string, indexes ...any) string {
	evalErr(fmt.Sprintf("%s tag add %s %s", w, tclSafeString(tagName), collectAny(indexes...)))
	return tagName
}

// Text — Create and manipulate 'text' hypertext editing widgets
//
// # Description
//
// Remove the tag tagName from all of the characters starting at index1 and
// ending just before index2 (the character at index2 is not affected). A single
// command may contain any number of index1-index2 pairs. If the last index2 is
// omitted then the tag is removed from the single character at index1. If there
// are no characters in the specified range (e.g. index1 is past the end of the
// file or index2 is less than or equal to index1) then the command has no effect.
// This command returns an empty string.
//
// Additional information might be available at the [Tcl/Tk text] page.
//
// [Tcl/Tk text]: https://www.tcl.tk/man/tcl9.0/TkCmd/text.html
func (w *TextWidget) TagRemove(tagName string, indexes ...any) string {
	return evalErr(fmt.Sprintf("%s tag remove %s %s", w, tclSafeString(tagName), collectAny(indexes...)))
}

// Text — Create and manipulate 'text' hypertext editing widgets
//
// # Description
//
// Delete a range of characters from the text. If both index1 and index2 are
// specified, then delete all the characters starting with the one given by
// index1 and stopping just before index2 (i.e. the character at index2 is not
// deleted). If index2 does not specify a position later in the text than
// index1 then no characters are deleted. If index2 is not specified then the
// single character at index1 is deleted. Attempts to delete characters in a
// way that would leave the text without a newline as the last character will
// be tweaked by the text widget to avoid this. In particular, deletion of
// complete lines of text up to the end of the text will also delete the
// newline character just before the deleted block so that it is replaced by
// the new final newline of the text widget. The command returns an empty
// string. If more indices are given, multiple ranges of text will be deleted.
// All indices are first checked for validity before any deletions are made.
// They are sorted and the text is removed from the last range to the first
// range so deleted text does not cause an undesired index shifting
// side-effects. If multiple ranges with the same start index are given, then
// the longest range is used. If overlapping ranges are given, then they will
// be merged into spans that do not cause deletion of text outside the given
// ranges due to text shifted during deletion.
//
// Additional information might be available at the [Tcl/Tk text] page.
//
// [Tcl/Tk text]: https://www.tcl.tk/man/tcl9.0/TkCmd/text.html
func (w *TextWidget) Delete(options ...any) {
	evalErr(fmt.Sprintf("%s delete %s", w, collectAny(options...)))
}

// Text — Create and manipulate 'text' hypertext editing widgets
//
// # Description
//
// InsertML inserts 'ml' at the end of 'w', interpreting it as a HTML-like
// markup language.
//
// It recognizes and treats accordingly the <br> tag.
//
// The <img> tag is reserved for embedded images. You can inline the image
// directly:
//
//	InsertML("Hello", NewPhoto(...), Align("top"), "world!")
//
// The <embed> tag is reserved for embedded widgets. You can inline a widget
// directly:
//
//	InsertML("Hello", Button(Txt("Foo")), Align("center"), "world!")
//
// The <pre> tag works similarly to HTML, ie. white space and line breaks are
// kept. To make the content of a <pre> rendered in monospace, configure the
// tag, for example:
//
//	t.TagConfigure("pre", Font(CourierFont(), 10)
//
// Other ML-tags are used as names of configured 'w' tags, if configured,
// ignored otherwise.
//
// Example usage in _examples/embed.go.
func (w *TextWidget) InsertML(list ...any) {
	var ml bytes.Buffer
	for i := 0; i < len(list); i++ {
		switch x := list[i].(type) {
		case string:
			ml.WriteString(x)
		case *Img:
			var opts Opts
			for j := i + 1; j < len(list); j++ {
				x, ok := list[j].(Opt)
				if !ok {
					break
				}

				opts = append(opts, x)
			}
			fmt.Fprintf(&ml, "<img src=%q", x)
			for _, v := range opts {
				fmt.Fprintf(&ml, " opt=%q", v.optionString(w.Window))
				i++
			}
			ml.WriteString(">")
		case Widget:
			var opts Opts
			for j := i + 1; j < len(list); j++ {
				x, ok := list[j].(Opt)
				if !ok {
					break
				}

				opts = append(opts, x)
			}
			fmt.Fprintf(&ml, "<embed src=%q", x)
			for _, v := range opts {
				fmt.Fprintf(&ml, " opt=%q", v.optionString(w.Window))
				i++
			}
			ml.WriteString(">")
		}
	}
	doc, err := html.Parse(&ml)
	if err != nil {
		fail(err)
		return
	}

	var tags []string
	var body int
	k := TkScaling() * 72 / 600
	walk(0, doc, func(lvl int, n *html.Node) bool {
		switch n.Type {
		case html.TextNode:
			if lvl < len(tags) {
				tags = tags[:lvl]
			}
			tags := tags[body+1:]
			for _, v := range tags {
				if v == "pre" {
					evalErr(fmt.Sprintf("%s insert end %s %s", w, tclSafeString(unescapeML(n.Data)), tclSafeStrings(tags...)))
					return true
				}
			}

			ids, toks := tokenize(n.Data)
			for i, id := range ids {
				switch id {
				default:
					if s := toks[i]; strings.HasPrefix(s, "$$") && strings.HasSuffix(s, "$$") || strings.HasPrefix(s, "$") && strings.HasSuffix(s, "$") {
						img := NewPhoto(Data(TeX(s, k)))
						evalErr(fmt.Sprintf("%v image create end -image %s -align top", w, img))
						break
					}

					fallthrough
				case 0:
					evalErr(fmt.Sprintf("%s insert end %s {%s}", w, tclFromElementNode(unescapeML(toks[i])), tclSafeStrings(tags...)))
				}
			}
		case html.ElementNode:
			switch n.Data {
			case "br":
				evalErr(fmt.Sprintf("%s insert end \\n {%s}", w, tclSafeStrings(tags...)))
			case "body":
				tags = append(tags, n.Data)
				body = lvl
			case "img":
				var src string
				var opts []string
				for _, v := range n.Attr {
					switch v.Key {
					case "src":
						src = v.Val
					case "opt":
						opts = append(opts, v.Val)
					}
				}
				evalErr(fmt.Sprintf("%v image create end -image %s %s", w, src, strings.Join(opts, " ")))
			case "embed":
				var src string
				var opts []string
				for _, v := range n.Attr {
					switch v.Key {
					case "src":
						src = v.Val
					case "opt":
						opts = append(opts, v.Val)
					}
				}
				evalErr(fmt.Sprintf("%v window create end -window %s %s", w, src, strings.Join(opts, " ")))
			default:
				tags = append(tags, n.Data)
			}
		}
		return true
	})
}

func unescapeML(s string) string {
	s = strings.ReplaceAll(s, "\\$", "$")
	s = strings.ReplaceAll(s, "&lt;", "<")
	s = strings.ReplaceAll(s, "&gt;", ">")
	return s
}

func tokenize(s string) (ids []int, toks []string) {
	for {
		id, len := mlToken(s)
		if len == 0 {
			return ids, toks
		}

		ids = append(ids, id)
		toks = append(toks, s[:len])
		s = s[len:]
	}
}

func tclFromElementNode(s string) string {
	a := strings.Fields(s)
	var prefix, suffix string
	if s != "" {
		switch s[0] {
		case '\n', ' ', '\t':
			prefix = " "
		}
		switch s[len(s)-1] {
		case '\n', ' ', '\t':
			suffix = " "
		}
	}
	for i, v := range a {
		a[i] = tclSafeInBraces(v)
	}
	r := fmt.Sprintf("{%s%s%s}", prefix, strings.Join(a, " "), suffix)
	return r
}

var badMLChars = [...]bool{
	'}':  true,
	'\\': true,
	'\n': true,
	'\r': true,
	'\t': true,
}

// Align option.
//
// Known uses:
//   - [Text] (widget specific, applies to embedded images)
func Align(val any) Opt {
	return rawOption(fmt.Sprintf(`-align %s`, optionString(val)))
}

func walk(lvl int, n *html.Node, visitor func(lvl int, n *html.Node) (dive bool)) {
	for ; n != nil; n = n.NextSibling {
		if visitor(lvl, n) {
			walk(lvl+1, n.FirstChild, visitor)
		}
	}
}

// Fontchooser — control font selection dialog
//
// # Description
//
// The tk fontchooser command controls the Tk font selection dialog. It uses
// the native platform font selection dialog where available, or a dialog
// implemented in Tcl otherwise.
//
// Unlike most of the other Tk dialog commands, tk fontchooser does not return
// an immediate result, as on some platforms (Mac OS X) the standard font
// dialog is modeless while on others (Windows) it is modal. To get the
// user-selected font use FontchooserFont() from a handler assigned via
// [Command].
//
// Set one or more of the configurations options below (analogous to Tk widget configuration).
//
//   - [Parent]
//
// Specifies/returns the logical parent window of the font selection dialog
// (similar to the -parent option to other dialogs). The font selection dialog
// is hidden if it is visible when the parent window is destroyed.
//
//   - [Title]
//
// Specifies/returns the title of the dialog. Has no effect on platforms where
// the font selection dialog does not support titles.
//
//   - [FontFace]
//
// Specifies/returns the font that is currently selected in the dialog if it is
// visible, or that will be initially selected when the dialog is shown (if
// supported by the platform). Can be set to the empty string to indicate that
// no font should be selected. Fonts can be specified in any form given by the
// "FONT DESCRIPTION" section in the font manual page.
//
//   - [Command]
//
// Specifies the command called when a font selection has been made by the
// user. To obtain the font description, call [FontchooserFont] from the
// handler.
//
// Additional information might be available at the [Tcl/Tk fontchooser] page.
//
// [Tcl/Tk fontchooser]: https://www.tcl.tk/man/tcl9.0/TkCmd/fontchooser.html
func Fontchooser(options ...Opt) {
	evalErr(fmt.Sprintf("tk fontchooser configure %s", collect(options...)))
}

// FontchooserFont — control font selection dialog
//
// # Description
//
// Returns the selected font description in the form
//
//	family size style...
//
// Additional information might be available at the [Tcl/Tk fontchooser] page.
//
// [Tcl/Tk fontchooser]: https://www.tcl.tk/man/tcl9.0/TkCmd/fontchooser.html
func FontchooserFont() []string {
	return parseList(evalErr("tk fontchooser config -font"))
}

// FontchooserShow — control font selection dialog
//
// # Description
//
// Show the font selection dialog. Depending on the platform, may return
// immediately or only once the dialog has been withdrawn.
//
// Additional information might be available at the [Tcl/Tk fontchooser] page.
//
// [Tcl/Tk fontchooser]: https://www.tcl.tk/man/tcl9.0/TkCmd/fontchooser.html
func FontchooserShow() {
	evalErr("tk  fontchooser show")
}

// FontchooserHide — control font selection dialog
//
// # Description
//
// Hide the font selection dialog if it is visible and cause any pending tk
// fontchooser show command to return.
//
// Additional information might be available at the [Tcl/Tk fontchooser] page.
//
// [Tcl/Tk fontchooser]: https://www.tcl.tk/man/tcl9.0/TkCmd/fontchooser.html
func FontchooserHide() {
	evalErr("tk fontchooser hide")
}

// GetOpenFile — pop up a dialog box for the user to select a file to open.
//
// # Description
//
// GetOpenFile pops up a dialog box for the user to select a file to open. The
// function is usually associated with the Open command in the File menu. Its
// purpose is for the user to select an existing file only. If the user enters
// a non-existent file, the dialog box gives the user an error prompt and
// requires the user to give an alternative selection. If an application allows
// the user to create new files, it should do so by providing a separate New
// menu command.
//
//   - [Defaultextension] extension
//
// Specifies a string that will be appended to the filename if the user enters
// a filename without an extension. The default value is the empty string,
// which means no extension will be appended to the filename in any case. This
// option is ignored on Mac OS X, which does not require extensions to
// filenames, and the UNIX implementation guesses reasonable values for this
// from the -filetypes option when this is not supplied.
//
//   - [Filetypes] filePatternList ([][FileType])
//
// If a File types listbox exists in the file dialog on the particular
// platform, this option gives the filetypes in this listbox. When the user
// choose a filetype in the listbox, only the files of that type are listed. If
// this option is unspecified, or if it is set to the empty list, or if the
// File types listbox is not supported by the particular platform then all
// files are listed regardless of their types. See the section SPECIFYING FILE
// PATTERNS below for a discussion on the contents of filePatternList.
//
//   - [Initialdir] directory
//
// Specifies that the files in directory should be displayed when the dialog
// pops up. If this parameter is not specified, the initial directory defaults
// to the current working directory on non-Windows systems and on Windows
// systems prior to Vista. On Vista and later systems, the initial directory
// defaults to the last user-selected directory for the application. If the
// parameter specifies a relative path, the return value will convert the
// relative path to an absolute path.
//
//   - [Initialfile] filename
//
// Specifies a filename to be displayed in the dialog when it pops up.
//
//   - [Multiple] boolean
//
// Allows the user to choose multiple files from the Open dialog.
//
//   - [Parent] window
//
// Makes window the logical parent of the file dialog. The file dialog is
// displayed on top of its parent window. On Mac OS X, this turns the file
// dialog into a sheet attached to the parent window.
//
//   - [Title] titleString
//
// Specifies a string to display as the title of the dialog box. If this option
// is not specified, then a default title is displayed.
//
// Additional information might be available at the [Tcl/Tk getopenfile] page.
//
// [Tcl/Tk getopenfile]: https://www.tcl.tk/man/tcl9.0/TkCmd/getOpenFile.html
func GetOpenFile(options ...Opt) (r []string) {
	return parseList(evalErr(fmt.Sprintf("tk_getOpenFile %s", collect(options...))))
}

// FileType specifies a single file type for the [Filetypes] option.
type FileType struct {
	TypeName   string   // Eg. "Go files"
	Extensions []string // Eg. []string{".go"}
	MacType    string   // Eg. "TEXT"
}

// GetSaveFile — pop up a dialog box for the user to select a file to save.
//
// # Description
//
// GetSaveFile pops up a dialog box for the user to select a file to save.
//
// The functio is usually associated with the Save as command in the File menu.
// If the user enters a file that already exists, the dialog box prompts the
// user for confirmation whether the existing file should be overwritten or
// not.
//
//   - [Confirmoverwrite] boolean
//
// Configures how the Save dialog reacts when the selected file already exists,
// and saving would overwrite it. A true value requests a confirmation dialog
// be presented to the user. A false value requests that the overwrite take
// place without confirmation. Default value is true.
//
//   - [Defaultextension] extension
//
// Specifies a string that will be appended to the filename if the user enters
// a filename without an extension. The default value is the empty string,
// which means no extension will be appended to the filename in any case. This
// option is ignored on Mac OS X, which does not require extensions to
// filenames, and the UNIX implementation guesses reasonable values for this
// from the -filetypes option when this is not supplied.
//
//   - [Filetypes] filePatternList ([][FileType])
//
// If a File types listbox exists in the file dialog on the particular
// platform, this option gives the filetypes in this listbox. When the user
// choose a filetype in the listbox, only the files of that type are listed. If
// this option is unspecified, or if it is set to the empty list, or if the
// File types listbox is not supported by the particular platform then all
// files are listed regardless of their types. See the section SPECIFYING FILE
// PATTERNS below for a discussion on the contents of filePatternList.
//
//   - [Initialdir] directory
//
// Specifies that the files in directory should be displayed when the dialog
// pops up. If this parameter is not specified, the initial directory defaults
// to the current working directory on non-Windows systems and on Windows
// systems prior to Vista. On Vista and later systems, the initial directory
// defaults to the last user-selected directory for the application. If the
// parameter specifies a relative path, the return value will convert the
// relative path to an absolute path.
//
//   - [Initialfile] filename
//
// Specifies a filename to be displayed in the dialog when it pops up.
//
//   - [Parent] window
//
// Makes window the logical parent of the file dialog. The file dialog is
// displayed on top of its parent window. On Mac OS X, this turns the file
// dialog into a sheet attached to the parent window.
//
//   - [Title] titleString
//
// Specifies a string to display as the title of the dialog box. If this option
// is not specified, then a default title is displayed.
//
// Additional information might be available at the [Tcl/Tk getopenfile] page.
//
// [Tcl/Tk getopenfile]: https://www.tcl.tk/man/tcl9.0/TkCmd/getOpenFile.html
func GetSaveFile(options ...Opt) string {
	return evalErr(fmt.Sprintf("tk_getSaveFile %s", collect(options...)))
}

// Place — Geometry manager for fixed or rubber-sheet placement
//
// # Description
//
// The placer is a geometry manager for Tk. It provides simple fixed placement
// of windows, where you specify the exact size and location of one window,
// called the content, within another window, called the container. The placer
// also provides rubber-sheet placement, where you specify the size and
// location of the content in terms of the dimensions of the container, so that
// the content changes size and location in response to changes in the size of
// the container. Lastly, the placer allows you to mix these styles of
// placement so that, for example, the content has a fixed width and height but
// is centered inside the container.
//
// The first argument must be a *Window.
//
// The following options are supported:
//
//   - [Anchor] where
//
// Where specifies which point of window is to be positioned at the (x,y)
// location selected by the -x, -y, -relx, and -rely options. The anchor point
// is in terms of the outer area of window including its border, if any. Thus
// if where is se then the lower-right corner of window's border will appear at
// the given (x,y) location in the container. The anchor position defaults to
// nw.
//
//   - [Bordermode] mode
//
// Mode determines the degree to which borders within the container are used in
// determining the placement of the content. The default and most common value
// is inside. In this case the placer considers the area of the container to be
// the innermost area of the container, inside any border: an option of -x 0
// corresponds to an x-coordinate just inside the border and an option of
// -relwidth 1.0 means window will fill the area inside the container's border.
//
// If mode is outside then the placer considers the area of the container to
// include its border; this mode is typically used when placing window outside
// its container, as with the options -x 0 -y 0 -anchor ne. Lastly, mode may be
// specified as ignore, in which case borders are ignored: the area of the
// container is considered to be its official X area, which includes any
// internal border but no external border. A bordermode of ignore is probably
// not very useful.
//
//   - [Height] size
//
// Size specifies the height for window in screen units (i.e. any of the forms
// accepted by Tk_GetPixels). The height will be the outer dimension of window
// including its border, if any. If size is an empty string, or if no -height
// or -relheight option is specified, then the height requested internally by
// the window will be used.
//
//   - [In] container
//
// Container specifies the path name of the window relative to which window is
// to be placed. Container must either be window's parent or a descendant of
// window's parent. In addition, container and window must both be descendants
// of the same top-level window. These restrictions are necessary to guarantee
// that window is visible whenever container is visible. If this option is not
// specified then the other window defaults to window's parent.
//
//   - [Relheight] size
//
// Size specifies the height for window. In this case the height is specified
// as a floating-point number relative to the height of the container: 0.5
// means window will be half as high as the container, 1.0 means window will
// have the same height as the container, and so on. If both -height and
// -relheight are specified for a content, their values are summed. For
// example, -relheight 1.0 -height -2 makes the content 2 pixels shorter than
// the container.
//
//   - [Relwidth] size
//
// Size specifies the width for window. In this case the width is specified as
// a floating-point number relative to the width of the container: 0.5 means
// window will be half as wide as the container, 1.0 means window will have the
// same width as the container, and so on. If both -width and -relwidth are
// specified for a content, their values are summed. For example, -relwidth 1.0
// -width 5 makes the content 5 pixels wider than the container.
//
//   - [Relx] location
//
// Location specifies the x-coordinate within the container window of the
// anchor point for window. In this case the location is specified in a
// relative fashion as a floating-point number: 0.0 corresponds to the left
// edge of the container and 1.0 corresponds to the right edge of the
// container. Location need not be in the range 0.0-1.0. If both -x and -relx
// are specified for a content then their values are summed. For example, -relx
// 0.5 -x -2 positions the left edge of the content 2 pixels to the left of the
// center of its container.
//
//   - [Rely] location
//
// Location specifies the y-coordinate within the container window of the
// anchor point for window. In this case the value is specified in a relative
// fashion as a floating-point number: 0.0 corresponds to the top edge of the
// container and 1.0 corresponds to the bottom edge of the container. Location
// need not be in the range 0.0-1.0. If both -y and -rely are specified for a
// content then their values are summed. For example, -rely 0.5 -x 3 positions
// the top edge of the content 3 pixels below the center of its container.
//
//   - [Width] size
//
// Size specifies the width for window in screen units (i.e. any of the forms
// accepted by Tk_GetPixels). The width will be the outer width of window
// including its border, if any. If size is an empty string, or if no -width or
// -relwidth option is specified, then the width requested internally by the
// window will be used.
//
//   - [X] location
//
// Location specifies the x-coordinate within the container window of the
// anchor point for window. The location is specified in screen units (i.e. any
// of the forms accepted by Tk_GetPixels) and need not lie within the bounds of
// the container window.
//
//   - [Y] location
//
// Location specifies the y-coordinate within the container window of the
// anchor point for window. The location is specified in screen units (i.e. any
// of the forms accepted by Tk_GetPixels) and need not lie within the bounds of
// the container window.
//
// If the same value is specified separately with two different options, such
// as -x and -relx, then the most recent option is used and the older one is
// ignored.
//
// Additional information might be available at the [Tcl/Tk place] page.
//
// [Tcl/Tk place]: https://www.tcl.tk/man/tcl9.0/TkCmd/place.htm
func Place(options ...Opt) {
	autocenterDisabled = true
	evalErr(fmt.Sprintf("place %s", collect(options...)))
}

// Lower — Change a window's position in the stacking order
//
// # Description
//
// If the belowThis argument is nil then the command lowers window so that
// it is below all of its siblings in the stacking order (it will be obscured
// by any siblings that overlap it and will not obscure any siblings). If
// belowThis is specified then it must be the path name of a window that is
// either a sibling of window or the descendant of a sibling of window. In this
// case the lower command will insert window into the stacking order just below
// belowThis (or the ancestor of belowThis that is a sibling of window); this
// could end up either raising or lowering window.
//
// All toplevel windows may be restacked with respect to each other, whatever
// their relative path names, but the window manager is not obligated to
// strictly honor requests to restack.
//
// Additional information might be available at the [Tcl/Tk lower] page.
//
// [Tcl/Tk lower]: https://www.tcl.tk/man/tcl9.0/TkCmd/lower.html
func (w *Window) Lower(belowThis Widget) {
	b := ""
	if belowThis != nil {
		b = belowThis.optionString(nil)
	}
	evalErr(fmt.Sprintf("lower %s %s", w, b))
}

// Raise — Change a window's position in the stacking order
//
// # Description
//
// If the aboveThis argument is nil then the command raises window so that it
// is above all of its siblings in the stacking order (it will not be obscured
// by any siblings and will obscure any siblings that overlap it). If aboveThis
// is specified then it must be the path name of a window that is either a
// sibling of window or the descendant of a sibling of window. In this case the
// raise command will insert window into the stacking order just above
// aboveThis (or the ancestor of aboveThis that is a sibling of window); this
// could end up either raising or lowering window.
//
// All toplevel windows may be restacked with respect to each other, whatever
// their relative path names, but the window manager is not obligated to
// strictly honor requests to restack.
//
// On macOS raising an iconified toplevel window causes it to be deiconified.
//
// Additional information might be available at the [Tcl/Tk raise] page.
//
// [Tcl/Tk raise]: https://www.tcl.tk/man/tcl9.0/TkCmd/raise.html
func (w *Window) Raise(aboveThis Widget) {
	b := ""
	if aboveThis != nil {
		b = aboveThis.optionString(nil)
	}
	evalErr(fmt.Sprintf("raise %s %s", w, b))
}

// MenuItem represents an entry on a menu.
type MenuItem struct {
	id string
}

// String implements fmt.Stringer.
func (m *MenuItem) String() string {
	return m.optionString(nil)
}

func (m *MenuItem) optionString(_ *Window) string {
	if m != nil {
		return m.id
	}

	return "mnu_non_existing"
}

// tk_popup — Post a popup menu
//
// # Description
//
// This procedure posts a menu at a given position on the screen and configures
// Tk so that the menu and its cascaded children can be traversed with the
// mouse or the keyboard. Menu is the name of a menu widget and x and y are the
// root coordinates at which to display the menu. If entry is omitted or an
// empty string, the menu's upper left corner is positioned at the given point.
// Otherwise entry gives the index of an entry in menu and the menu will be
// positioned so that the entry is positioned over the given point.
//
// Additional information might be available at the [Tcl/Tk popup] page.
//
// [Tcl/Tk popup]: https://www.tcl.tk/man/tk9.0/TkCmd/popup.htm
func Popup(menu *Window, x, y int, entry any) {
	var s string
	if entry != nil && entry != "" {
		s = collectAny(entry)
	}
	evalErr(fmt.Sprintf("tk_popup %s %v %v %s", menu, x, y, s))
}

// Menu — Create and manipulate 'menu' widgets and menubars
//
// # Description
//
// Add a new radiobutton entry to the bottom of the menu.
//
// Additional information might be available at the [Tcl/Tk menu] page.
//
// [Tcl/Tk menu]: https://www.tcl.tk/man/tk9.0/TkCmd/menu.htm
func (w *MenuWidget) AddRadiobutton(options ...Opt) *MenuItem {
	return &MenuItem{id: evalErr(fmt.Sprintf("%s add radiobutton %s", w, winCollect(w.Window, options...)))}
}

// Menu — Create and manipulate 'menu' widgets and menubars
//
// # Description
//
// Add a new checkbutton entry to the bottom of the menu.
//
// Additional information might be available at the [Tcl/Tk menu] page.
//
// [Tcl/Tk menu]: https://www.tcl.tk/man/tk9.0/TkCmd/menu.htm
func (w *MenuWidget) AddCheckbutton(options ...Opt) *MenuItem {
	return &MenuItem{id: evalErr(fmt.Sprintf("%s add checkbutton %s", w, winCollect(w.Window, options...)))}
}

// Menu — Create and manipulate 'menu' widgets and menubars
//
// # Description
//
// Add a new command entry to the bottom of the menu.
//
// Additional information might be available at the [Tcl/Tk menu] page.
//
// [Tcl/Tk menu]: https://www.tcl.tk/man/tk9.0/TkCmd/menu.htm
func (w *MenuWidget) AddCommand(options ...Opt) *MenuItem {
	return &MenuItem{id: evalErr(fmt.Sprintf("%s add command %s", w, winCollect(w.Window, options...)))}
}

// Menu — Create and manipulate 'menu' widgets and menubars
//
// # Description
//
// Add a new cascade entry to the end of the menu.
//
// Additional information might be available at the [Tcl/Tk menu] page.
//
// [Tcl/Tk menu]: https://www.tcl.tk/man/tk9.0/TkCmd/menu.htm
func (w *MenuWidget) AddCascade(options ...Opt) *MenuItem {
	return &MenuItem{id: evalErr(fmt.Sprintf("%s add cascade %s", w, winCollect(w.Window, options...)))}
}

// Menu — Create and manipulate 'menu' widgets and menubars
//
// # Description
//
// Add a new separator entry to the bottom of the menu.
//
// Additional information might be available at the [Tcl/Tk menu] page.
//
// [Tcl/Tk menu]: https://www.tcl.tk/man/tk9.0/TkCmd/menu.htm
func (w *MenuWidget) AddSeparator(options ...Opt) *MenuItem {
	return &MenuItem{id: evalErr(fmt.Sprintf("%s add separator %s", w, winCollect(w.Window, options...)))}
}

// Menu — Create and manipulate 'menu' widgets and menubars
//
// # Description
//
// Invoke the action of the menu entry. See the sections on the individual
// entries above for details on what happens. If the menu entry is disabled
// then nothing happens. If the entry has a command associated with it then the
// result of that command is returned as the result of the invoke widget
// command. Otherwise the result is an empty string. Note: invoking a menu
// entry does not automatically unpost the menu; the default bindings normally
// take care of this before invoking the invoke widget command.
//
// Additional information might be available at the [Tcl/Tk menu] page.
//
// [Tcl/Tk menu]: https://www.tcl.tk/man/tk9.0/TkCmd/menu.htm
func (w *MenuWidget) Invoke(index uint) {
	evalErr(fmt.Sprintf("%s invoke %d", w, index))
}

// Menu — Create and manipulate 'menu' widgets and menubars
//
// # Description
//
// This command is similar to the configure command, except that it applies to
// the options for an individual entry, whereas configure applies to the
// options for the menu as a whole. Options may have any of the values
// described in the MENU ENTRY OPTIONS section below. If options are specified,
// options are modified as indicated in the command and the command returns an
// empty string. If no options are specified, returns a list describing the
// current options for entry index (see Tk_ConfigureInfo for information on the
// format of this list).
//
// Additional information might be available at the [Tcl/Tk menu] page.
//
// [Tcl/Tk menu]: https://www.tcl.tk/man/tk9.0/TkCmd/menu.htm
func (w *MenuWidget) EntryConfigure(index uint, options ...Opt) {
	evalErr(fmt.Sprintf("%s entryconfigure %d %s", w, index, winCollect(w.Window, options...)))
}

// Menu — Create and manipulate 'menu' widgets and menubars
//
// # Description
//
// Returns the id of the menu entry given by index. This is the identifier that was
// assigned to the entry when it was created using the add or insert widget command.
// Returns an empty string for the tear-off entry, or if index is equivalent to {}.
//
// Note that while the raw tk command returns an id, that alone is of no practical use
// to the caller. So a MenuItem representing the id is returned.
//
// Additional information might be available at the [Tcl/Tk menu] page.
//
// [Tcl/Tk menu]: https://www.tcl.tk/man/tk9.0/TkCmd/menu.htm
func (m *MenuWidget) Id(index int) *MenuItem {
	id := evalErr(fmt.Sprintf("%s id %d", m, index))
	return &MenuItem{id: id}
}

// Menu — Create and manipulate 'menu' widgets and menubars
//
// # Description
//
// Returns the numerical index corresponding to a MenuItem.
//
// Additional information might be available at the [Tcl/Tk menu] page.
//
// [Tcl/Tk menu]: https://www.tcl.tk/man/tk9.0/TkCmd/menu.htm
func (m *MenuWidget) Index(item *MenuItem) int {
	index := evalErr(fmt.Sprintf("%s index %s", m, item.id))
	return atoi(index)
}

// TScrollbar — Control the viewport of a scrollable widget
//
// # Description
//
// This command is normally invoked by the scrollbar's associated widget from
// an -xscrollcommand or -yscrollcommand callback. Specifies the visible range
// to be displayed. first and last are real fractions between 0 and 1.
//
// More information might be available at the [Tcl/Tk ttk_scrollbar] page.
//
// [Tcl/Tk ttk_scrollbar]: https://www.tcl.tk/man/tk9.0/TkCmd/ttk_scrollbar.htm
func (w *TScrollbarWidget) Set(firstLast string) {
	evalErr(fmt.Sprintf("%s set %s", w, firstLast))
}

// tk — Manipulate Tk internal state
//
// # Description
//
// Sets and queries the current scaling factor used by Tk to convert between
// physical units (for example, points, inches, or millimeters) and pixels. The
// number argument is a floating point number that specifies the number of
// pixels per point on window's display. If the window argument is omitted, it
// defaults to the main window. If the number argument is omitted, the current
// value of the scaling factor is returned.
//
// A “point” is a unit of measurement equal to 1/72 inch. A scaling factor of
// 1.0 corresponds to 1 pixel per point, which is equivalent to a standard 72
// dpi monitor. A scaling factor of 1.25 would mean 1.25 pixels per point,
// which is the setting for a 90 dpi monitor; setting the scaling factor to
// 1.25 on a 72 dpi monitor would cause everything in the application to be
// displayed 1.25 times as large as normal. The initial value for the scaling
// factor is set when the application starts, based on properties of the
// installed monitor, but it can be changed at any time. Measurements made
// after the scaling factor is changed will use the new scaling factor, but it
// is undefined whether existing widgets will resize themselves dynamically to
// accommodate the new scaling factor.
//
// - [Displayof] window
//
// - Number
//
// Additional information might be available at the [Tcl/Tk tk] page.
//
// [Tcl/Tk tk]: https://www.tcl.tk/man/tcl9.0/TkCmd/tk.html
func TkScaling(options ...any) float64 {
	var a []Opt
	for _, v := range options {
		switch x := v.(type) {
		case Opts:
			a = append(a, x)
		case Opt:
			a = append(a, x)
		default:
			a = append(a, stringOption(fmt.Sprint(x)))
		}
	}
	if s := evalErr(fmt.Sprintf("tk scaling %s", collect(a...))); s != "" {
		n, err := strconv.ParseFloat(s, 64)
		if err == nil {
			return n
		}

		fail(err)
	}
	return 1
}

// Font option.
//
// Specifies the font to use when drawing text inside the widget.
// The value may have any of the forms described in the font manual
// page under FONT DESCRIPTION.
//
// Known uses:
//   - [Button]
//   - [Checkbutton]
//   - [Entry]
//   - [Fontchooser] (command specific)
//   - [Label]
//   - [Labelframe]
//   - [Listbox]
//   - [MenuWidget.AddCascade] (command specific)
//   - [MenuWidget.AddCommand] (command specific)
//   - [MenuWidget.AddSeparator] (command specific)
//   - [Menu]
//   - [Menubutton]
//   - [Message]
//   - [Radiobutton]
//   - [Scale]
//   - [Spinbox]
//   - [TEntry]
//   - [TLabel]
//   - [TProgressbar]
//   - [TextWidget.TagConfigure] (command specific)
//   - [Text]
func Font(list ...any) Opt {
	return rawOption(fmt.Sprintf(`-font {%s}`, tclSafeList(list...)))
}

// Font — Get the configured option value.
//
// Known uses:
//   - [Button]
//   - [Checkbutton]
//   - [Entry]
//   - [Label]
//   - [Labelframe]
//   - [Listbox]
//   - [Menu]
//   - [Menubutton]
//   - [Message]
//   - [Radiobutton]
//   - [Scale]
//   - [Spinbox]
//   - [TEntry]
//   - [TLabel]
//   - [TProgressbar]
//   - [Text]
func (w *Window) Font() string {
	return evalErr(fmt.Sprintf(`%s cget -font`, w))
}

// + ttk::style configure style ?-option ?value option value...? ?
// + ttk::style element args
// 	+ ttk::style element create elementName type ?args...?
// 	+ ttk::style element names
// 	+ ttk::style element options element
// + ttk::style layout style ?layoutSpec?
// + ttk::style lookup style -option ?state ?default??
// + ttk::style map style ?-option { statespec value... }?
// - ttk::style theme args
// 	- ttk::style theme create themeName ?-parent basedon? ?-settings script... ?
// 	+ ttk::style theme names
// 	- ttk::style theme settings themeName script
// 	+ ttk::style theme styles ?themeName?
// 	+ ttk::style theme use ?themeName?

// ttk::style — Manipulate style database
//
// # Description
//
// Sets the default value of the specified option(s) in style. If style does
// not exist, it is created. Example:
//
//	StyleConfigure(".", Font("times"), Background(LightBlue))
//
// If only style and -option are specified, get the default value for option
// -option of style style. Example:
//
//	StyleConfigure(".", Font)
//
// If only style is specified, get the default value for all options of style
// style. Example:
//
//	StyleConfigure(".")
//
// Additional information might be available at the [Tcl/Tk style] page.
// There's also a [Styles and Themes] tutorial at tkdoc.com.
//
// [Tcl/Tk style]: https://www.tcl.tk/man/tcl9.0/TkCmd/ttk_style.html
// [Styles and Themes]: https://tkdocs.com/tutorial/styles.html
func StyleConfigure(style string, options ...any) []string {
	if len(options) == 0 {
		return parseList(evalErr(fmt.Sprintf("ttk::style configure %s", tclSafeString(style))))
	}

	if len(options) == 1 {
		o := options[0]
		if x, ok := o.(Opt); ok {
			return []string{evalErr(fmt.Sprintf("ttk::style configure %s %s", tclSafeString(style), x.optionString(nil)))}
		}

		if s := funcToTclOption(o); s != "" {
			return []string{evalErr(fmt.Sprintf("ttk::style configure %s %s", tclSafeString(style), s))}
		}

		return nil
	}

	var a []string
	for _, v := range options {
		switch x := v.(type) {
		case Opt:
			a = append(a, x.optionString(nil))
		default:
			fail(fmt.Errorf("expected Opt: %T", x))
			return nil
		}
	}

	return []string{evalErr(fmt.Sprintf("ttk::style configure %s %s", tclSafeString(style), strings.Join(a, " ")))}
}

func funcToTclOption(fn any) string {
	t := reflect.TypeOf(fn)
	if t.Kind() != reflect.Func {
		return ""
	}

	s := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
	if s == "" {
		fail(fmt.Errorf("failed to determine function name"))
		return ""
	}

	a := strings.Split(s, ".")
	if len(a) == 0 {
		fail(fmt.Errorf("failed to determine function name: %q", s))
		return ""
	}

	s = strings.ToLower(a[len(a)-1])
	if r, ok := replaceOpt[s]; ok {
		return "-" + r
	}

	return "-" + s
}

var replaceOpt = map[string]string{
	"btn": "button",
	"lbl": "label",
	"mnu": "menu",
	"msg": "message",
	"txt": "text",
}

// ttk::style — Manipulate style database
//
// # Description
//
// Creates a new element in the current theme of type type. The only
// cross-platform built-in element type is image (see ttk_image(n)) but themes
// may define other element types (see Ttk_RegisterElementFactory). On suitable
// versions of Windows an element factory is registered to create Windows theme
// elements (see ttk_vsapi(n)). Examples:
//
//	StyleElementCreate("TSpinbox.uparrow", "from", "default") // Inherit the existing element from theme 'default'.
//
//	StyleElementCreate("Red.Corner.TButton.indicator", "image", NewPhoto(File("red_corner.png")), Width(10))
//
//	imageN := NewPhoto(...)
//	StyleElementCreate("TCheckbutton.indicator", "image", image5, "disabled selected", image6, "disabled alternate",
//		image8, "disabled", image9, "alternate", image7, "!selected", image4, Width(20), Border(4), Sticky("w"))
//
// After the type "image" comes a list of one or more images. Every image is
// optionally followed by a space separated list of states the image applies
// to. An exclamation mark before the state is a negation.
//
// Additional information might be available at the [Tcl/Tk modifying a button]
// and [Tcl/Tk style] pages.
// There's also a [Styles and Themes] tutorial at tkdoc.com.
//
// [Tcl/Tk style]: https://www.tcl.tk/man/tcl9.0/TkCmd/ttk_style.html
// [Tcl/Tk modifying a button]: https://wiki.tcl-lang.org/page/Tutorial%3A+Modifying+a+ttk+button%27s+style
// [Styles and Themes]: https://tkdocs.com/tutorial/styles.html
func StyleElementCreate(elementName, typ string, options ...any) string {
	// ttk::style element create TRadiobutton.indicator image {pyimage11 {disabled selected} pyimage12 disabled pyimage13 !selected pyimage10} -width 20 -border 4 -sticky w
	var a, b []string
	for _, v := range options {
		switch x := v.(type) {
		case *Img:
			a = append(a, x.optionString(nil))
		case Opt:
			b = append(b, x.optionString(nil))
		default:
			switch s := strings.Fields(fmt.Sprint(v)); {
			case len(s) == 1:
				a = append(a, tclSafeInBraces(s[0]))
			default:
				for i, v := range s {
					s[i] = tclSafeInBraces(v)
				}
				a = append(a, fmt.Sprintf("{%s}", strings.Join(s, " ")))
			}
		}
	}
	switch {
	case len(a) == 2 && a[0] == "from":
		return evalErr(fmt.Sprintf("ttk::style element create from %s", a[1]))
	default:
		return evalErr(fmt.Sprintf("ttk::style element create {%s} image {%s} %s", tclSafeInBraces(elementName), strings.Join(a, " "), strings.Join(b, " ")))
	}
}

// ttk::style — Manipulate style database
//
// # Description
//
// Returns the list of elements defined in the current theme.
//
// Additional information might be available at the [Tcl/Tk style] page.
// There's also a [Styles and Themes] tutorial at tkdoc.com.
//
// [Tcl/Tk style]: https://www.tcl.tk/man/tcl9.0/TkCmd/ttk_style.html
// [Styles and Themes]: https://tkdocs.com/tutorial/styles.html
func StyleElementNames() []string {
	return parseList(evalErr("ttk::style element names"))
}

// ttk::style — Manipulate style database
//
// # Description
//
// Returns the list of element's options.
//
// Additional information might be available at the [Tcl/Tk style] page.
// There's also a [Styles and Themes] tutorial at tkdoc.com.
//
// [Tcl/Tk style]: https://www.tcl.tk/man/tcl9.0/TkCmd/ttk_style.html
// [Styles and Themes]: https://tkdocs.com/tutorial/styles.html
func StyleElementOptions(element string) []string {
	return parseList(evalErr(fmt.Sprintf("ttk::style element options %s", tclSafeString(element))))
}

// ttk::style — Manipulate style database
//
// # Description
//
// Define the widget layout for style style. See LAYOUTS below for the format
// of layoutSpec. If layoutSpec is omitted, return the layout specification for
// style style. Example:
//
//	StyleLayout("Red.Corner.Button",
//		"Button.border", Sticky("nswe"), Border(1), Children(
//			"Button.focus", Sticky("nswe"), Children(
//				"Button.padding", Sticky("nswe"), Children(
//					"Button.label", Sticky("nswe"),
//					"Red.Corner.TButton.indicator", Side("right"), Sticky("ne")))))
//
// Additional information might be available at the [Tcl/Tk modifying a button]
// and [Tcl/Tk style] pages.
// There's also a [Styles and Themes] tutorial at tkdoc.com.
//
// [Tcl/Tk style]: https://www.tcl.tk/man/tcl9.0/TkCmd/ttk_style.html
// [Tcl/Tk modifying a button]: https://wiki.tcl-lang.org/page/Tutorial%3A+Modifying+a+ttk+button%27s+style
// [Styles and Themes]: https://tkdocs.com/tutorial/styles.html
func StyleLayout(style string, options ...any) string {
	if len(options) == 0 {
		return evalErr(fmt.Sprintf("ttk::style layout %s", tclSafeString(style)))
	}

	evalErr(fmt.Sprintf("ttk::style layout %s %s", tclSafeString(style), children("", options...)))
	return ""
}

// ttk::style — Manipulate style database
//
// # Description
//
// Returns the value specified for -option in style style in state state, using
// the standard lookup rules for element options. state is a list of state
// names; if omitted, it defaults to all bits off (the “normal” state). If the
// default argument is present, it is used as a fallback value in case no
// specification for -option is found. If style does not exist, it is created.
// For example,
//
//	StyleLookup("TButton", Font)
//
// may return "TkDefaultFont", depending on the operating system, theme in use
// and the configured style options.
//
// Additional information might be available at the [Tcl/Tk style] page.
// There's also a [Styles and Themes] tutorial at tkdoc.com.
//
// [Tcl/Tk style]: https://www.tcl.tk/man/tcl9.0/TkCmd/ttk_style.html
// [Styles and Themes]: https://tkdocs.com/tutorial/styles.html
func StyleLookup(style string, options ...any) string {
	for i, v := range options {
		if s := funcToTclOption(v); s != "" {
			options[i] = rawOption(s)
		}
	}

	return evalErr(fmt.Sprintf("ttk::style lookup %s %s", tclSafeString(style), collectAny(options...)))
}

// ttk::style — Manipulate style database
//
// # Description
//
// Sets dynamic (state dependent) values of the specified option(s) in style.
// Each statespec / value pair is examined in order; the value corresponding to
// the first matching statespec is used. If style does not exist, it is
// created. If only style and -option are specified, get the dynamic values for
// option -option of style style. If only style is specified, get the dynamic
// values for all options of style 'style'.
//
// With no options the function returns the currently configured style map for
// 'style'.  For example,
//
//	StyleMap("TButton")
//
// may return "-relief {{!disabled pressed} sunken}", depending on the
// operating system, theme in use and the configured style options.
//
// Setting a style map is done by providing a list of options, each option is
// followed by a list of states and a value. For example:
//
//	StyleMap("TButton",
//		Background, "disabled", "#112233", "active", "#445566",
//		Foreground, "disabled", "#778899",
//	        Relief, "pressed", "!disabled", "sunken")
//
// # Widget states
//
// The widget state is a bitmap of independent state flags.
//
//   - active - The mouse cursor is over the widget and pressing a mouse button
//     will cause some action to occur
//   - alternate - A widget-specific alternate display format
//   - background - Windows and Mac have a notion of an “active” or foreground
//     window. The background state is set for widgets in a background window,
//     and cleared for those in the foreground window
//   - disabled - Widget is disabled under program control
//   - focus - Widget has keyboard focus
//   - invalid - The widget’s value is invalid
//   - pressed - Widget is being pressed
//   - readonly - Widget should not allow user modification
//   - selected - “On”, “true”, or “current” for things like Checkbuttons and
//     radiobuttons
//   - hover - The mouse cursor is within the widget. This is similar to the
//     active state; it is used in some themes for widgets that provide distinct
//     visual feedback for the active widget in addition to the active element
//     within the widget.
//   - user1-user6 - Freely usable for other purposes
//
// A state specification is a sequence of state names, optionally prefixed with
// an exclamation point indicating that the bit is off.
//
// Additional information might be available at the [Tcl/Tk style] page.
// There's also a [Styles and Themes] tutorial at tkdoc.com.
//
// [Tcl/Tk style]: https://www.tcl.tk/man/tcl9.0/TkCmd/ttk_style.html
// [Styles and Themes]: https://tkdocs.com/tutorial/styles.html
func StyleMap(style string, options ...any) string {
	if len(options) == 0 {
		return evalErr(fmt.Sprintf("ttk::style map %s", tclSafeString(style)))
	}

	a, err := parseStyleMapOpts(options...)
	if err != nil {
		fail(fmt.Errorf("parsing StyleMap options: %v", err))
		return ""
	}

	evalErr(fmt.Sprintf("ttk::style map %s %s", tclSafeString(style), strings.Join(a, " ")))
	return ""
}

func parseStyleMapOpts(in ...any) (r []string, err error) {
	for len(in) != 0 {
		opt := funcToTclOption(in[0])
		if opt == "" {
			return nil, fmt.Errorf("expected option, eg. 'Relief' (the function, not Relief(argument)), got %T", in[0])
		}

		in = in[1:]
		r = append(r, opt)
		var list []string
		for {
			var states []string
			for len(in) != 0 {
				s, ok := in[0].(string)
				if !ok || !isState(s) {
					break
				}

				states = append(states, s)
				in = in[1:]
			}
			if len(in) == 0 {
				return nil, fmt.Errorf("missing option value")
			}

			val := tclSafeString(fmt.Sprint(in[0]))
			in = in[1:]

			var s string
			switch len(states) {
			case 0:
				// nop
			case 1:
				s = states[0]
			default:
				s = fmt.Sprintf("{%s}", strings.Join(states, " "))
			}
			list = append(list, fmt.Sprintf("%s %s", s, val))
			if len(in) == 0 || funcToTclOption(in[0]) != "" {
				break
			}
		}
		r = append(r, fmt.Sprintf("{%s}", strings.Join(list, " ")))
	}
	return r, nil
}

// https://www.tcl-lang.org/man/tcl9.0/TkCmd/ttk_widget.html#M34
func isState(s string) bool {
	if len(s) == 0 {
		return false
	}

	if s[0] == '!' {
		s = s[1:]
	}

	switch s {
	case
		"active",
		"disabled",
		"focus",
		"pressed",
		"selected",
		"background",
		"readonly",
		"alternate",
		"invalid",
		"hover",
		"user1",
		"user2",
		"user3",
		"user4",
		"user5",
		"user6":

		return true
	default:
		return false
	}
}

// Border option.
//
// Known uses:
//   - [StyleLayout]
func Border(val any) Opt {
	return rawOption(fmt.Sprintf(`-border %s`, optionString(val)))
}

// Focuscolor option.
//
// Known uses:
//   - [StyleConfigure]
func Focuscolor(val any) Opt {
	return rawOption(fmt.Sprintf(`-focuscolor %s`, optionString(val)))
}

// Focusthickness option.
//
// Known uses:
//   - [StyleConfigure]
func Focusthickness(val any) Opt {
	return rawOption(fmt.Sprintf(`-focusthickness %s`, optionString(val)))
}

// Focussolid option.
//
// Known uses:
//   - [StyleConfigure]
func Focussolid(val any) Opt {
	return rawOption(fmt.Sprintf(`-focussolid %s`, optionString(val)))
}

// Children option.
//
// Known uses:
//   - [StyleLayout]
//
// Children describes children of a style layout.
func Children(list ...any) Opt {
	return children("-children", list...)
}

func children(prefixed string, list ...any) Opt {
	var a []string
	for _, v := range list {
		a = append(a, fmt.Sprint(v))
	}
	return rawOption(fmt.Sprintf(" %s {%s}", prefixed, strings.Join(a, " ")))
}

// ttk::style — Manipulate style database
//
// # Description
//
// Returns a list of all known themes.
//
// Additional information might be available at the [Tcl/Tk style] page.
//
// [Tcl/Tk style]: https://www.tcl.tk/man/tcl9.0/TkCmd/ttk_style.html
func StyleThemeNames() []string {
	return parseList(evalErr("ttk::style theme names"))
}

// ttk::style — Manipulate style database
//
// # Description
//
// Returns a list of all styles in themeName. If themeName is omitted, the
// current theme is used.
//
// Additional information might be available at the [Tcl/Tk style] page.
//
// [Tcl/Tk style]: https://www.tcl.tk/man/tcl9.0/TkCmd/ttk_style.html
func StyleThemeStyles(themeName ...string) []string {
	var s string
	if len(themeName) != 0 {
		s = tclSafeString(themeName[0])
	}
	return parseList(evalErr(fmt.Sprintf("ttk::style theme styles %s", s)))
}

// ttk::style — Manipulate style database
//
// # Description
//
// Without a argument the result is the name of the current theme. Otherwise
// this command sets the current theme to themeName, and refreshes all widgets.
//
// Additional information might be available at the [Tcl/Tk style] page.
// There's also a [Styles and Themes] tutorial at tkdoc.com.
//
// This mechanism is separate from the RegisterTheme/ActivateTheme one.  Use of
// this function is best suited for the built in themes like "clam" etc.
//
// [Tcl/Tk style]: https://www.tcl.tk/man/tcl9.0/TkCmd/ttk_style.html
// [Styles and Themes]: https://tkdocs.com/tutorial/styles.html
func StyleThemeUse(themeName ...string) string {
	var s string
	if len(themeName) != 0 {
		s = tclSafeString(themeName[0])
	}
	return evalErr(fmt.Sprintf("ttk::style theme use %s", s))
}

// CourierFont returns "{courier new}" on Windows and "courier" elsewhere.
func CourierFont() string {
	if goos == "windows" {
		return "courier new"
	}

	return "courier"
}

// button — Create and manipulate 'button' action widgets
//
// # Description
//
// Invoke the Tcl command associated with the button, if there is one. The
// return value is the return value from the Tcl command, or an empty string if
// there is no command associated with the button. This command is ignored if
// the button's state is disabled.
//
// Additional information might be available at the [Tcl/Tk button] page.
//
// [Tcl/Tk button]: https://www.tcl.tk/man/tcl9.0/TkCmd/button.html
func (w *ButtonWidget) Invoke() string {
	return evalErr(fmt.Sprintf("%s invoke", w))
}

// TButton — Widget that issues a command when pressed
//
// # Description
//
// Invokes the command associated with the button.
//
// Additional information might be available at the [Tcl/Tk ttk::button] page.
//
// [Tcl/Tk ttk::button]: https://www.tcl.tk/man/tcl9.0/TkCmd/ttk_button.html
func (w *TButtonWidget) Invoke() string {
	return evalErr(fmt.Sprintf("%s invoke", w))
}

// TButton — Widget that issues a command when pressed
//
// # Description
//
// Shiftrelief specifies how far the button contents are shifted down and right
// in the pressed state. This action provides additional skeuomorphic feedback.
func Shiftrelief(val any) Opt {
	return rawOption(fmt.Sprintf(`-shiftrelief %s`, optionString(val)))
}

// Bordercolor — Styling widgets
//
// Bordercolor is a styling option of one or more widgets. Please see [Changing
// Widget Colors] for details.
//
// [Changing Widget Colors]: https://wiki.tcl-lang.org/page/Changing+Widget+Colors
func Bordercolor(val any) Opt {
	return rawOption(fmt.Sprintf(`-bordercolor %s`, optionString(val)))
}

// Darkcolor — Styling widgets
//
// Darkcolor is a styling option of one or more widgets. Please see [Changing
// Widget Colors] for details.
//
// [Changing Widget Colors]: https://wiki.tcl-lang.org/page/Changing+Widget+Colors
func Darkcolor(val any) Opt {
	return rawOption(fmt.Sprintf(`-darkcolor %s`, optionString(val)))
}

// Lightcolor — Styling widgets
//
// Lightcolor is a styling option of one or more widgets. Please see [Changing
// Widget Colors] for details.
//
// [Changing Widget Colors]: https://wiki.tcl-lang.org/page/Changing+Widget+Colors
func Lightcolor(val any) Opt {
	return rawOption(fmt.Sprintf(`-lightcolor %s`, optionString(val)))
}

// Indicatorbackground — Styling widgets
//
// Indicatorbackground is a styling option of one or more widgets. Please see
// [Changing Widget Colors] for details.
//
// [Changing Widget Colors]: https://wiki.tcl-lang.org/page/Changing+Widget+Colors
func Indicatorbackground(val any) Opt {
	return rawOption(fmt.Sprintf(`-indicatorbackground %s`, optionString(val)))
}

// Indicatorcolor — Styling widgets
//
// Indicatorcolor is a styling option of one or more widgets. Please see
// [Changing Widget Colors] for details.
//
// [Changing Widget Colors]: https://wiki.tcl-lang.org/page/Changing+Widget+Colors
func Indicatorcolor(val any) Opt {
	return rawOption(fmt.Sprintf(`-indicatorcolor %s`, optionString(val)))
}

// Indicatormargin — Styling widgets
//
// Indicatormargin is a styling option of one or more widgets. Please see
// [Changing Widget Colors] for details.
//
// [Changing Widget Colors]: https://wiki.tcl-lang.org/page/Changing+Widget+Colors
func Indicatormargin(val any) Opt {
	return rawOption(fmt.Sprintf(`-indicatormargin %s`, optionString(val)))
}

// Indicatorrelief — Styling widgets
//
// Indicatorrelief is a styling option of one or more widgets. Please see
// [Changing Widget Colors] for details.
//
// [Changing Widget Colors]: https://wiki.tcl-lang.org/page/Changing+Widget+Colors
func Indicatorrelief(val any) Opt {
	return rawOption(fmt.Sprintf(`-indicatorrelief %s`, optionString(val)))
}

// Arrowcolor — Styling widgets
//
// Arrowcolor is a styling option of one or more widgets. Please see
// [Changing Widget Colors] for details.
//
// [Changing Widget Colors]: https://wiki.tcl-lang.org/page/Changing+Widget+Colors
func Arrowcolor(val any) Opt {
	return rawOption(fmt.Sprintf(`-arrowcolor %s`, optionString(val)))
}

// Arrowsize — Styling widgets
//
// Arrowsize is a styling option of one or more widgets. Please see
// [Changing Widget Colors] for details.
//
// [Changing Widget Colors]: https://wiki.tcl-lang.org/page/Changing+Widget+Colors
func Arrowsize(val any) Opt {
	return rawOption(fmt.Sprintf(`-arrosize %s`, optionString(val)))
}

// Focusfill — Styling widgets
//
// Focusfill is a styling option of a ttk::combobox.
// More information might be available at the [Tcl/Tk ttk_combobox] page.
//
// [Tcl/Tk ttk_combobox]: https://www.tcl.tk/man/tcl9.0/TkCmd/ttk_combobox.html
func Focusfill(val any) Opt {
	return rawOption(fmt.Sprintf(`-focusfill %s`, optionString(val)))
}

// Fieldbackground — Styling widgets
//
// Fieldbackground is a styling option of one or more widgets. Please see
// [Changing Widget Colors] for details.
//
// [Changing Widget Colors]: https://wiki.tcl-lang.org/page/Changing+Widget+Colors
func Fieldbackground(val any) Opt {
	return rawOption(fmt.Sprintf(`-fieldbackground %s`, optionString(val)))
}

// Insertcolor — Styling widgets
//
// Insertcolor is a styling option of one or more widgets. Please see
// [Changing Widget Colors] for details.
//
// [Changing Widget Colors]: https://wiki.tcl-lang.org/page/Changing+Widget+Colors
func Insertcolor(val any) Opt {
	return rawOption(fmt.Sprintf(`-insertcolor %s`, optionString(val)))
}

// Postoffset — Styling widgets
//
// Postoffset is a styling option of a ttk::combobox.
// More information might be available at the [Tcl/Tk ttk_combobox] page.
//
// [Tcl/Tk ttk_combobox]: https://www.tcl.tk/man/tcl9.0/TkCmd/ttk_combobox.html
func Postoffset(val any) Opt {
	return rawOption(fmt.Sprintf(`-postoffset %s`, optionString(val)))
}

// Labelmargins — Styling widgets
//
// Labelmargins is a styling option of a ttk::labelframe.
// More information might be available at the [Tcl/Tk ttk_labelframe] page.
//
// [Tcl/Tk ttk_labelframe]: https://www.tcl.tk/man/tcl9.0/TkCmd/ttk_labelframe.html
func Labelmargins(val any) Opt {
	return rawOption(fmt.Sprintf(`-labelmargins %s`, optionString(val)))
}

// Labeloutside — Styling widgets
//
// Labeloutside is a styling option of a ttk::labelframe.
// More information might be available at the [Tcl/Tk ttk_labelframe] page.
//
// [Tcl/Tk ttk_labelframe]: https://www.tcl.tk/man/tcl9.0/TkCmd/ttk_labelframe.html
func Labeloutside(val any) Opt {
	return rawOption(fmt.Sprintf(`-labeloutside %s`, optionString(val)))
}

// Tabmargins — Styling widgets
//
// Tabmargins is a styling option of a ttk::notebook.
// More information might be available at the [Tcl/Tk ttk_notebook] page.
//
// [Tcl/Tk ttk_notebook]: https://www.tcl.tk/man/tcl9.0/TkCmd/ttk_notebook.html
func Tabmargins(val any) Opt {
	return rawOption(fmt.Sprintf(`-tabmargins %s`, optionString(val)))
}

// Tabposition — Styling widgets
//
// Tabposition is a styling option of a ttk::notebook.
//
// More information might be available at the [Tcl/Tk ttk_notebook] page.
//
// [Tcl/Tk ttk_notebook]: https://www.tcl.tk/man/tcl9.0/TkCmd/ttk_notebook.html
func Tabposition(val any) Opt {
	return rawOption(fmt.Sprintf(`-tabposition %s`, optionString(val)))
}

// Sashthickness — Styling widgets
//
// Sashthickness is a styling option of one or more widgets. Please see
// [Changing Widget Colors] for details.
//
// [Changing Widget Colors]: https://wiki.tcl-lang.org/page/Changing+Widget+Colors
func Sashthickness(val any) Opt {
	return rawOption(fmt.Sprintf(`-sashthickness %s`, optionString(val)))
}

// panedwindow — Create and manipulate 'panedwindow' split container widgets
//
// # Description
//
// Add one or more windows to the panedwindow, each in a separate pane. The
// arguments consist of the names of one or more windows followed by pairs of
// arguments that specify how to manage the windows. Option may have any of the
// values accepted by the configure subcommand.
//
// More information might be available at the [Tcl/Tk panedwindow] page.
//
// [Tcl/Tk panedwindow]: https://www.tcl.tk/man/tcl9.0/TkCmd/panedwindow.html
func (w *PanedwindowWidget) Add(subwindow *Window, options ...Opt) {
	evalErr(fmt.Sprintf("%s add %s %s", w, subwindow, collect(options...)))
}

// panedwindow — Create and manipulate 'panedwindow' split container widgets
//
// # Description
//
// Query or modify the management options for window. If no option is
// specified, returns a list describing all of the available options for
// pathName (see Tk_ConfigureInfo for information on the format of this list).
// If option is specified with no value, then the command returns a list
// describing the one named option (this list will be identical to the
// corresponding sublist of the value returned if no option is specified). If
// one or more option-value pairs are specified, then the command modifies the
// given widget option(s) to have the given value(s); in this case the command
// returns an empty string.
//
// More information might be available at the [Tcl/Tk panedwindow] page.
//
// [Tcl/Tk panedwindow]: https://www.tcl.tk/man/tcl9.0/TkCmd/panedwindow.html
func (w *PanedwindowWidget) Paneconfigure(subwindow *Window, options ...Opt) string {
	return evalErr(fmt.Sprintf("%s paneconfigure %s %s", w, subwindow, collect(options...)))
}

// Stretch option.
//
// # Description
//
// Controls how extra space is allocated to each of the panes. When is one of
// always, first, last, middle, and never. The panedwindow will calculate the
// required size of all its panes. Any remaining (or deficit) space will be
// distributed to those panes marked for stretching. The space will be
// distributed based on each panes current ratio of the whole. The when values
// have the following definition:
//
//   - always: This pane will always stretch.
//   - first: Only if this pane is the first pane (left-most or top-most) will it stretch.
//   - last: Only if this pane is the last pane (right-most or bottom-most) will it stretch. This is the default value.
//   - middle: Only if this pane is not the first or last pane will it stretch.
//   - never: This pane will never stretch.
//
// Known uses:
//   - [PanedwindowWidget] (command specific)
func Stretch(when string) Opt {
	return rawOption(fmt.Sprintf(`-stretch %s`, tclSafeString(when)))
}

// panedwindow — Create and manipulate 'panedwindow' split container widgets
//
// # Description
//
// Query a management option for window. Option may be any value allowed by the
// paneconfigure subcommand.
//
// More information might be available at the [Tcl/Tk panedwindow] page.
//
// [Tcl/Tk panedwindow]: https://www.tcl.tk/man/tcl9.0/TkCmd/panedwindow.html
func (w *PanedwindowWidget) Panecget(subwindow *Window, opt any) string {
	return evalErr(fmt.Sprintf("%s panecget %s %s", w, subwindow, funcToTclOption(opt)))
}

// ttk::panedwindow — Multi-pane container window
//
// # Description
//
// Adds a new pane to the window. See PANE OPTIONS for the list of available options.
//
// More information might be available at the [Tcl/Tk ttk_panedwindow] page.
//
// [Tcl/Tk ttk_panedwindow]: https://www.tcl.tk/man/tcl9.0/TkCmd/ttk_panedwindow.html
func (w *TPanedwindowWidget) Add(subwindow *Window, options ...Opt) {
	evalErr(fmt.Sprintf("%s add %s %s", w, subwindow, collect(options...)))
}

// Gripsize — Styling widgets
//
// Gripsize is a styling option of a ttk::panedwindow and ttk::scrollbar.
// More information might be available at the [Tcl/Tk ttk_panedwindow] or
// [Tcl/Tk ttk_scrollbar] page.
//
// [Tcl/Tk ttk_panedwindow]: https://www.tcl.tk/man/tcl9.0/TkCmd/ttk_panedwindow.html
// [Tcl/Tk ttk_scrollbar]: https://www.tcl.tk/man/tcl9.0/TkCmd/ttk_scrollbar.html
func Gripsize(val any) Opt {
	return rawOption(fmt.Sprintf(`-tabposition %s`, optionString(val)))
}

// Maxphase — Styling widgets
//
// Maxphase is a styling option of a ttk::progressbar.
// More information might be available at the [Tcl/Tk ttk_progressbar] page.
//
// [Tcl/Tk ttk_progressbar]: https://www.tcl.tk/man/tcl9.0/TkCmd/ttk_progressbar.html
func Maxphase(val any) Opt {
	return rawOption(fmt.Sprintf(`-maxphase %s`, optionString(val)))
}

// Period — Styling widgets
//
// Period is a styling option of a ttk::progressbar.
// More information might be available at the [Tcl/Tk ttk_progressbar] page.
//
// [Tcl/Tk ttk_progressbar]: https://www.tcl.tk/man/tcl9.0/TkCmd/ttk_progressbar.html
func Period(val any) Opt {
	return rawOption(fmt.Sprintf(`-period %s`, optionString(val)))
}

// Groovewidth — Styling widgets
//
// Groovewidth is a styling option of a ttk::scale.
// More information might be available at the [Tcl/Tk ttk_scale] page.
//
// [Tcl/Tk ttk_scale]: https://www.tcl.tk/man/tcl9.0/TkCmd/ttk_scale.html
func Groovewidth(val any) Opt {
	return rawOption(fmt.Sprintf(`-groovewidth %s`, optionString(val)))
}

// Sliderwidth — Styling widgets
//
// Sliderwidth is a styling option of a ttk::scale.
// More information might be available at the [Tcl/Tk ttk_scale] page.
//
// [Tcl/Tk ttk_scale]: https://www.tcl.tk/man/tcl9.0/TkCmd/ttk_scale.html
func Sliderwidth(val any) Opt {
	return rawOption(fmt.Sprintf(`-sliderwidth %s`, optionString(val)))
}

// Troughrelief — Styling widgets
//
// Troughrelief is a styling option of one or more widgets. Please see
// [Changing Widget Colors] for details.
//
// [Changing Widget Colors]: https://wiki.tcl-lang.org/page/Changing+Widget+Colors
func Troughrelief(val any) Opt {
	return rawOption(fmt.Sprintf(`-troughrelief %s`, optionString(val)))
}

// Indent — Styling widgets
//
// Indent is a styling option of a ttk::treeview.
// More information might be available at the [Tcl/Tk ttk_treeview] page.
//
// [Tcl/Tk ttk_treeview]: https://www.tcl.tk/man/tcl9.0/TkCmd/ttk_treeview.html
func Indent(val any) Opt {
	return rawOption(fmt.Sprintf(`-indent %s`, optionString(val)))
}

// Columnseparatorwidth — Styling widgets
//
// Columnseparatorwidth is a styling option of a ttk::treeview.
// More information might be available at the [Tcl/Tk ttk_treeview] page.
//
// [Tcl/Tk ttk_treeview]: https://www.tcl.tk/man/tcl9.0/TkCmd/ttk_treeview.html
func Columnseparatorwidth(val any) Opt {
	return rawOption(fmt.Sprintf(`-columnseparatorwidth %s`, optionString(val)))
}

// Rowheight — Styling widgets
//
// Rowheight is a styling option of one or more widgets. Please see
// [Changing Widget Colors] for details.
//
// [Changing Widget Colors]: https://wiki.tcl-lang.org/page/Changing+Widget+Colors
func Rowheight(val any) Opt {
	return rawOption(fmt.Sprintf(`-rowheight %s`, optionString(val)))
}

// Stripedbackground — Styling widgets
//
// Stripedbackground is a styling option of a ttk::treeview.
// More information might be available at the [Tcl/Tk ttk_treeview] page.
//
// [Tcl/Tk ttk_treeview]: https://www.tcl.tk/man/tcl9.0/TkCmd/ttk_treeview.html
func Stripedbackground(val any) Opt {
	return rawOption(fmt.Sprintf(`-stripedbackground %s`, optionString(val)))
}

// Indicatormargins — Styling widgets
//
// Indicatormargins is a styling option of a ttk::treeview.
// More information might be available at the [Tcl/Tk ttk_treeview] page.
//
// [Tcl/Tk ttk_treeview]: https://www.tcl.tk/man/tcl9.0/TkCmd/ttk_treeview.html
func Indicatormargins(val any) Opt {
	return rawOption(fmt.Sprintf(`-indicatormargins %s`, optionString(val)))
}

// Indicatorsize — Styling widgets
//
// Indicatorsize is a styling option of a ttk::treeview.
// More information might be available at the [Tcl/Tk ttk_treeview] page.
//
// [Tcl/Tk ttk_treeview]: https://www.tcl.tk/man/tcl9.0/TkCmd/ttk_treeview.html
func Indicatorsize(val any) Opt {
	return rawOption(fmt.Sprintf(`-indicatorsize %s`, optionString(val)))
}

// Topmost option.
//
// Known uses:
//   - [WmAttributes] (command specific)
func Topmost(v bool) Opt {
	return rawOption(fmt.Sprintf(`-topmost %s`, optionString(v)))
}

// Type option.
//
// Known uses:
//   - [ClipboardAppend] (command specific)
//   - [ClipboardGet] (command specific)
//   - [Menu] (widget specific)
//   - [MessageBox] (command specific)
//   - [WmAttributes] (command specific)
func Type(val any) Opt {
	return rawOption(fmt.Sprintf(`-type %s`, optionString(val)))
}

// Type — Get the configured option value.
//
// Known uses:
//   - [Menu] (widget specific)
//   - [WmAttributes] (command specific)
func (w *Window) Type() string {
	return evalErr(fmt.Sprintf(`%s cget -type`, w))
}

// wm — Communicate with window manager
//
// # Description
//
//   - wm attributes window
//   - wm attributes window ?option?
//   - wm attributes window ?option value option value...?
//
// This subcommand returns or sets platform specific attributes associated with
// a window. The first form returns a list of the platform specific flags and
// their values. The second form returns the value for the specific option. The
// third form sets one or more of the values. The values are as follows:
//
// [Topmost]: Specifies whether this is a topmost window (displays above all other windows).
//
// [Type]: Requests that the window should be interpreted by the window manager
// as being of the specified type(s). This may cause the window to be decorated
// in a different way or otherwise managed differently, though exactly what
// happens is entirely up to the window manager. A list of types may be used,
// in order of preference. The following values are mapped to constants defined
// in the EWMH specification (using others is possible, but not advised):
//
//   - "desktop"
//     Indicates a desktop feature.
//   - "dock"
//     Indicates a dock/panel feature.
//   - "toolbar"
//     Indicates a toolbar window that should be acting on behalf of another
//     window, as indicated with wm transient.
//   - "menu"
//     Indicates a torn-off menu that should be acting on behalf of another
//     window, as indicated with wm transient.
//   - "utility"
//     Indicates a utility window (e.g., palette or toolbox) that should be acting
//     on behalf of another window, as indicated with wm transient.
//   - "splash"
//     Indicates a splash screen, displayed during application start up.
//   - "dialog"
//     Indicates a general dialog window, that should be acting on behalf of
//     another window, as indicated with wm transient.
//   - "dropdownMenu"
//     Indicates a menu summoned from a menu bar, which should usually also be set
//     to be override-redirected (with wm overrideredirect).
//   - "popupMenu"
//     Indicates a popup menu, which should usually also be set to be
//     override-redirected (with wm overrideredirect).
//   - "tooltip"
//     Indicates a tooltip window, which should usually also be set to be
//     override-redirected (with wm overrideredirect).
//   - "notification"
//     Indicates a window that provides a background notification of some event,
//     which should usually also be set to be override-redirected (with wm
//     overrideredirect).
//   - "combo"
//     Indicates the drop-down list of a combobox widget, which should usually
//     also be set to be override-redirected (with wm overrideredirect).
//   - "dnd"
//     Indicates a window that represents something being dragged, which should
//     usually also be set to be override-redirected (with wm overrideredirect).
//   - "normal"
//     Indicates a window that has no special interpretation.
//
// More information might be available at the [Tcl/Tk wm] page.
//
// [Tcl/Tk wm]: https://www.tcl.tk/man/tcl9.0/TkCmd/wm.html
func WmAttributes(w *Window, options ...any) string {
	switch len(options) {
	case 0:
		return evalErr(fmt.Sprintf("wm attributes %s", w))
	case 1:
		if s := funcToTclOption(options[0]); s != "" {
			return evalErr(fmt.Sprintf("wm attributes %s %s", w, s))
		}

		fallthrough
	default:
		return evalErr(fmt.Sprintf("wm attributes %s %s", w, collectAny(options...)))
	}
}

// wm — Communicate with window manager
//
// # Description
//
// This command is used to manage window manager protocols. The name argument
// in the wm protocol command is the name of an atom corresponding to a window
// manager protocol. Examples include WM_DELETE_WINDOW or WM_SAVE_YOURSELF or
// WM_TAKE_FOCUS.  A window manager protocol is a class of messages sent from a
// window manager to a Tk application outside of the normal event processing
// system. The main example is the WM_DELETE_WINDOW protocol; these messages
// are sent when the user clicks the close widget in the title bar of a window.
// Handlers for window manager protocols are installed with the wm protocol
// command. As a rule, if no handler has been installed for a protocol by the
// wm protocol command then all messages of that protocol are ignored. The
// WM_DELETE_WINDOW protocol is an exception to this rule. At start-up Tk
// installs a handler for this protocol, which responds by destroying the
// window. The wm protocol command can be used to replace this default handler
// by one which responds differently.
//
// A common requirement is to make sure that even if the user clicks the
// application's close button (X) a close handling function will be
// called (e.g., to prompt to save unsaved changes). For example, to
// ensure a custom func onQuit() function is called use:
//
//	WmProtocol(App, "WM_DELETE_WINDOW", onQuit)
//
// The list of available window manager protocols depends on the window
// manager, but all window managers supported by Tk provide WM_DELETE_WINDOW.
// On the Windows platform, a WM_SAVE_YOURSELF message is sent on user logout
// or system restart.
//
// If both name and command are specified, then command becomes the handler for
// the protocol specified by name. The atom for name will be added to window's
// WM_PROTOCOLS property to tell the window manager that the application has a
// handler for the protocol specified by name, and command will be invoked in
// the future whenever the window manager sends a message of that protocol to
// the Tk application. In this case the wm protocol command returns an empty
// string. If name is specified but command is not (is nil), then the current
// handler for name is returned, or an empty string if there is no handler
// defined for name (as a special case, the default handler for
// WM_DELETE_WINDOW is not returned). If command is specified as an empty
// string then the atom for name is removed from the WM_PROTOCOLS property of
// window and the handler is destroyed; an empty string is returned. Lastly, if
// neither name nor command is specified, the wm protocol command returns a
// list of all of the protocols for which handlers are currently defined for
// window.
//
// More information might be available at the [Tcl/Tk wm] page.
//
// [Tcl/Tk wm]: https://www.tcl.tk/man/tcl9.0/TkCmd/wm.html
func WmProtocol(w *Window, name string, command any) string {
	switch {
	case command == nil:
		return evalErr(fmt.Sprintf("wm protocol %s %s", w, tclSafeString(name)))
	case command == "":
		return evalErr(fmt.Sprintf("wm protocol %s %s {}", w, tclSafeString(name)))
	default:
		switch x := command.(type) {
		case *eventHandler:
			x.tcl = ""
			return evalErr(fmt.Sprintf("wm protocol %s %s %s", w, tclSafeString(name), collect(x)))
		default:
			return evalErr(fmt.Sprintf("wm protocol %s %s %s", w, tclSafeString(name), newEventHandler("", command).optionString(w)))
		}
	}
}

// wm — Communicate with window manager
//
// # Description
//
// If container is specified, then the window manager is informed that window
// is a transient window (e.g. pull-down menu) working on behalf of container
// (where container is the path name for a top-level window). If container is
// specified as an empty string then window is marked as not being a transient
// window any more. Otherwise the command returns the path name of window's
// current container, or an empty string if window is not currently a transient
// window. A transient window will mirror state changes in the container and
// inherit the state of the container when initially mapped. The directed graph
// with an edge from each transient to its container must be acyclic. In
// particular, it is an error to attempt to make a window a transient of
// itself. The window manager may also decorate a transient window differently,
// removing some features normally present (e.g., minimize and maximize
// buttons) though this is entirely at the discretion of the window manager.
//
// More information might be available at the [Tcl/Tk wm] page.
//
// [Tcl/Tk wm]: https://www.tcl.tk/man/tcl9.0/TkCmd/wm.html
func WmTransient(options ...Opt) string {
	return evalErr(fmt.Sprintf("wm transient %s", collect(options...)))
}

// wm — Communicate with window manager
//
// # Description
//
// Arrange for window to be iconified. It window has not yet been mapped for
// the first time, this command will arrange for it to appear in the iconified
// state when it is eventually mapped.
//
// More information might be available at the [Tcl/Tk wm] page.
//
// [Tcl/Tk wm]: https://www.tcl.tk/man/tcl9.0/TkCmd/wm.html
func WmIconify(w *Window) {
	if w == App {
		appIconified = true
		appDeiconified = false
	}
	wmIconify(w)
}

func wmIconify(w *Window) {
	evalErr(fmt.Sprintf("wm iconify %s", w))
}

// wm — Communicate with window manager
//
// # Description
//
// Arrange for window to be displayed in normal (non-iconified) form. This is
// done by mapping the window. If the window has never been mapped then this
// command will not map the window, but it will ensure that when the window is
// first mapped it will be displayed in de-iconified form. On Windows, a
// deiconified window will also be raised and be given the focus (made the
// active window). Returns an empty string.
//
// More information might be available at the [Tcl/Tk wm] page.
//
// [Tcl/Tk wm]: https://www.tcl.tk/man/tcl9.0/TkCmd/wm.html
func WmDeiconify(w *Window) {
	if w == App {
		appDeiconified = true
		appIconified = false
	}
	wmDeiconify(w)
}

func wmDeiconify(w *Window) {
	evalErr(fmt.Sprintf("wm deiconify %s", w))
}

// wm — Communicate with window manager
//
// # Description
//
// Arranges for window to be withdrawn from the screen. This causes the window
// to be unmapped and forgotten about by the window manager. If the window has
// never been mapped, then this command causes the window to be mapped in the
// withdrawn state. Not all window managers appear to know how to handle
// windows that are mapped in the withdrawn state. Note that it sometimes seems
// to be necessary to withdraw a window and then re-map it (e.g. with wm
// deiconify) to get some window managers to pay attention to changes in window
// attributes such as group.
//
// More information might be available at the [Tcl/Tk wm] page.
//
// [Tcl/Tk wm]: https://www.tcl.tk/man/tcl9.0/TkCmd/wm.html
func WmWithdraw(w *Window) {
	if w == App {
		appWithdrawn = true
	}
	wmWithdraw(w)
}

func wmWithdraw(w *Window) {
	evalErr(fmt.Sprintf("wm withdraw %s", w))
}

// wm — Communicate with window manager
//
// # Description
//
// If newGeometry is specified, then the geometry of window is changed and an
// empty string is returned. Otherwise the current geometry for window is
// returned (this is the most recent geometry specified either by manual
// resizing or in a wm geometry command). NewGeometry has the form
// =widthxheight±x±y, where any of =, widthxheight, or ±x±y may be omitted.
// Width and height are positive integers specifying the desired dimensions of
// window. If window is gridded (see GRIDDED GEOMETRY MANAGEMENT below) then
// the dimensions are specified in grid units; otherwise they are specified in
// pixel units.
//
// X and y specify the desired location of window on the screen, in pixels. If
// x is preceded by +, it specifies the number of pixels between the left edge
// of the screen and the left edge of window's border; if preceded by - then x
// specifies the number of pixels between the right edge of the screen and the
// right edge of window's border. If y is preceded by + then it specifies the
// number of pixels between the top of the screen and the top of window's
// border; if y is preceded by - then it specifies the number of pixels between
// the bottom of window's border and the bottom of the screen.
//
// If newGeometry is specified as an empty string then any existing
// user-specified geometry for window is cancelled, and the window will revert
// to the size requested internally by its widgets.
//
// Note that this is related to winfo geometry, but not the same. That can only
// query the geometry, and always reflects Tk's current understanding of the
// actual size and location of window, whereas wm geometry allows both setting
// and querying of the window manager's understanding of the size and location
// of the window. This can vary significantly, for example to reflect the
// addition of decorative elements to window such as title bars, and window
// managers are not required to precisely follow the requests made through this
// command.
//
// More information might be available at the [Tcl/Tk wm] page.
//
// [Tcl/Tk wm]: https://www.tcl.tk/man/tcl9.0/TkCmd/wm.html
func WmGeometry(w *Window, geometry ...string) string {
	// https://gitlab.com/cznic/tk9.0/-/commit/94774b2ee20e8d417f8eb12fa838eb9dcf58e8c6#note_2314759986
	switch {
	case len(geometry) == 0:
		return evalErr(fmt.Sprintf("wm geometry %s", w))
	default:
		if w == App {
			autocenterDisabled = true
		}
		return evalErr(fmt.Sprintf("wm geometry %s %s", w, tclSafeString(geometry[0])))
	}
}

// wm — Communicate with window manager
//
// # Description
//
// Width and height give the maximum permissible
// dimensions for window. For gridded windows the dimensions are specified in
// grid units; otherwise they are specified in pixel units. The window manager
// will restrict the window's dimensions to be less than or equal to width and
// height. If width and height are specified, then the command returns an empty
// string. Otherwise it returns a Tcl list with two elements, which are the
// maximum width and height currently in effect. The maximum size defaults to
// the size of the screen. See the sections on geometry management below for
// more information.
//
// More information might be available at the [Tcl/Tk wm] page.
//
// [Tcl/Tk wm]: https://www.tcl.tk/man/tcl9.0/TkCmd/wm.html
func WmSetMaxSize(w *Window, width, height int) {
	evalErr(fmt.Sprintf("wm maxsize %s %v %v", w, width, height))
}

// wm — Communicate with window manager
//
// # Description
//
// Returns the  maximum width and height currently in effect. The maximum size defaults to
// the size of the screen. See the sections on geometry management below for
// more information.
//
// More information might be available at the [Tcl/Tk wm] page.
//
// [Tcl/Tk wm]: https://www.tcl.tk/man/tcl9.0/TkCmd/wm.html
func WmMaxSize(w *Window) (width, height int) {
	a := strings.Fields(evalErr(fmt.Sprintf("wm maxsize %s", w)))
	if len(a) != 2 {
		return -1, -1
	}

	var err error
	if width, err = strconv.Atoi(a[0]); err != nil {
		return -1, -1
	}

	if height, err = strconv.Atoi(a[1]); err != nil {
		return -1, -1
	}

	return width, height
}

// wm — Communicate with window manager
//
// # Description
//
// If width and height are specified, they give the minimum permissible
// dimensions for window. For gridded windows the dimensions are specified in
// grid units; otherwise they are specified in pixel units. The window manager
// will restrict the window's dimensions to be greater than or equal to width
// and height. If width and height are specified, then the command returns an
// empty string. Otherwise it returns a Tcl list with two elements, which are
// the minimum width and height currently in effect. The minimum size defaults
// to one pixel in each dimension. See the sections on geometry management
// below for more information.
//
// More information might be available at the [Tcl/Tk wm] page.
//
// [Tcl/Tk wm]: https://www.tcl.tk/man/tcl9.0/TkCmd/wm.html
func WmMinSize(w *Window, widthHeight ...any) (r string) {
	arg := ""
	if len(widthHeight) >= 2 {
		arg = fmt.Sprintf("%s %s", tclSafeString(fmt.Sprint(widthHeight[0])), tclSafeString(fmt.Sprint(widthHeight[1])))
	}
	return evalErr(fmt.Sprintf("wm minsize %s %s", w, arg))
}

// Initalize enforces the parts of package initialization that are otherwise
// done lazily. The function may panic if ErrorMode is PanicOnError.
func Initialize() {
	lazyInit()
}

// ttk::notebook — Multi-paned container widget
//
// # Description
//
// Adds a new tab to the notebook. See TAB OPTIONS for the list of available
// options. If window is currently managed by the notebook but hidden, it is
// restored to its previous position.
//
// More information might be available at the [Tcl/Tk TNotebook] page.
//
// [Tcl/Tk TNotebook]: https://www.tcl.tk/man/tcl9.0/TkCmd/ttk_notebook.html
func (w *TNotebookWidget) Add(options ...Opt) {
	evalErr(fmt.Sprintf("%s add %v", w, winCollect(w.Window, options...)))
}

// TNotebook — Multi-paned container widget
//
// # Description
//
// Selects the specified tab. The associated content window will be displayed,
// and the previously-selected window (if different) is unmapped. If tabid is
// omitted, returns the widget name of the currently selected pane.
//
// More information might be available at the [Tcl/Tk TNotebook] page.
//
// [Tcl/Tk TNotebook]: https://www.tcl.tk/man/tcl9.0/TkCmd/ttk_notebook.html
func (w *TNotebookWidget) Select(tabid any) string {
	var arg string
	if tabid != nil && tabid != "" {
		arg = tclSafeString(fmt.Sprint(tabid))
	}
	return evalErr(fmt.Sprintf("%s select %s", w, arg))
}

// TNotebook — Multi-paned container widget
//
// # Description
//
// Returns the list of windows managed by the notebook, in the index order of
// their associated tabs.
//
// More information might be available at the [Tcl/Tk TNotebook] page.
//
// [Tcl/Tk TNotebook]: https://www.tcl.tk/man/tcl9.0/TkCmd/ttk_notebook.html
func (w *TNotebookWidget) Tabs() (r []*Window) {
	a := parseList(evalErr(fmt.Sprintf("%s tabs", w)))
	for _, v := range a {
		r = append(r, windowIndex[v])
	}
	return r
}

// tk_dialog — Create modal dialog and wait for response
//
// # Description
//
// This procedure is part of the Tk script library. It is largely deprecated by
// the tk_messageBox. Its arguments describe a dialog box:
//
//   - window - Name of top-level window to use for dialog. Any existing window
//     by this name is destroyed.
//
//   - title - Text to appear in the window manager's title bar for the dialog.
//
//   - text - Message to appear in the top portion of the dialog box.
//
//   - bitmap - If non-empty, specifies a bitmap (in a form suitable for
//     Tk_GetBitmap) to display in the top portion of the dialog, to the left of
//     the text. If this is an empty string then no bitmap is displayed in the
//     dialog.
//
//   - defaultButton - If this is an integer greater than or equal to zero, then it
//     gives the index of the button that is to be the default button for the
//     dialog (0 for the leftmost button, and so on). If negative or an empty
//     string then there will not be any default button.
//
//   - buttons - There will be one button for each of these arguments. Each string
//     specifies text to display in a button, in order from left to right.
//
// After creating a dialog box, tk_dialog waits for the user to select one of
// the buttons either by clicking on the button with the mouse or by typing
// return to invoke the default button (if any). Then it returns the index of
// the selected button: 0 for the leftmost button, 1 for the button next to it,
// and so on. If the dialog's window is destroyed before the user selects one
// of the buttons, then -1 is returned.
//
// While waiting for the user to respond, tk_dialog sets a local grab. This
// prevents the user from interacting with the application in any way except to
// invoke the dialog box.
// func Dialog(window *Window, title, text, bitmap string, defaultButton int, buttons ...string) (r int) {
// 	s := evalErr(fmt.Sprintf("tk_dialog %s %s %s", window, tclSafeStrings(title, text, bitmap), tclSafeStrings(buttons...)))
// 	if s == "" {
// 		return -1
// 	}
//
// 	var err error
// 	if r, err = strconv.Atoi(s); err != nil {
// 		fail(err)
// 		return -1
// 	}
//
// 	return r
// }

type Ticker struct {
	eh *eventHandler
}

func NewTicker(d time.Duration, handler func()) (r *Ticker, err error) {
	eh := newEventHandler("", handler)
	nm := fmt.Sprintf("ticker%v", id.Add(1))
	if _, err = eval(fmt.Sprintf(`proc %s {} {
	after %v {
		eventDispatcher %v
		%[1]s
	}
}
%[1]s
`, nm, d.Milliseconds(), eh.id)); err != nil {
		return nil, err
	}

	return &Ticker{eh: eh}, nil
}

// ttk::checkbutton — On/off widget
//
// # Description
//
// Toggles between the selected and deselected states and evaluates the
// associated -command. If the widget is currently selected, sets the -variable
// to the -offvalue and deselects the widget; otherwise, sets the -variable to
// the -onvalue. Returns the result of the -command.
//
// More information might be available at the [Tcl/Tk TCheckbutton] page.
//
// [Tcl/Tk TCheckbutton]: https://www.tcl.tk/man/tcl9.0/TkCmd/ttk_checkbutton.html
func (w *TCheckbuttonWidget) Invoke() string {
	return evalErr(fmt.Sprintf("%s invoke", w))
}

// // ttk::checkbutton — On/off widget
// //
// // # Description
// //
// // Set the on/off state of 'w'.
// func (w *TCheckbuttonWidget) Set(on bool) string {
// 	return evalErr(fmt.Sprintf("set %s %v", w.variable(cfgVar), on))
// }
//
// func cfgVar(w *Window, nm string) {
// 	w.Configure(Variable(nm))
// }
//
// func (w *Window) variable(reg func(w *Window, nm string)) (r string) {
// 	if r, ok := variables[w]; ok {
// 		return r
// 	}
//
// 	r = fmt.Sprintf("::tk9var%d", id.Add(1))
// 	variables[w] = r
// 	if reg != nil {
// 		reg(w, r)
// 	}
// 	return r
// }
//
// // ttk::checkbutton — On/off widget
// //
// // # Description
// //
// // Get the on/off state of 'w'.
// func (w *TCheckbuttonWidget) Get() bool {
// 	return toBool(evalErr(fmt.Sprintf("return $%s", w.variable(cfgVar))))
// }
//
// func toBool(s string) (r bool) {
// 	switch s {
// 	case "false", "0":
// 		return false
// 	case "true", "1":
// 		return true
// 	}
//
// 	fail(fmt.Errorf("invalid boolean: %q", s))
// 	return false
// }

// checkbutton — Create and manipulate 'checkbutton' boolean selection widgets
//
// # Description
//
// Selects the checkbutton and sets the associated variable to its “on” value.
//
// More information might be available at the [Tcl/Tk checkbutton] page.
//
// [Tcl/Tk checkbutton]: https://www.tcl.tk/man/tcl9.0/TkCmd/checkbutton.html
func (w *CheckbuttonWidget) Select() {
	evalErr(fmt.Sprintf("%s select", w))
}

// checkbutton — Create and manipulate 'checkbutton' boolean selection widgets
//
// # Description
//
// Does just what would have happened if the user invoked the checkbutton with
// the mouse: toggle the selection state of the button and invoke the Tcl
// command associated with the checkbutton, if there is one. The return value
// is the return value from the Tcl command, or an empty string if there is no
// command associated with the checkbutton. This command is ignored if the
// checkbutton's state is disabled.
//
// More information might be available at the [Tcl/Tk checkbutton] page.
//
// [Tcl/Tk checkbutton]: https://www.tcl.tk/man/tcl9.0/TkCmd/checkbutton.html
func (w *CheckbuttonWidget) Invoke() string {
	return evalErr(fmt.Sprintf("%s invoke", w))
}

// checkbutton — Create and manipulate 'checkbutton' boolean selection widgets
//
// # Description
//
// Deselects the checkbutton and sets the associated variable to its “off” value.
//
// More information might be available at the [Tcl/Tk checkbutton] page.
//
// [Tcl/Tk checkbutton]: https://www.tcl.tk/man/tcl9.0/TkCmd/checkbutton.html
func (w *CheckbuttonWidget) Deselect() {
	evalErr(fmt.Sprintf("%s deselect", w))
}

// font — Create and inspect fonts.
//
// # Description
//
// Query or modify the desired attributes for the named font called fontname.
// If no option is specified, returns a list describing all the options and
// their values for fontname. If a single option is specified with no value,
// then returns the current value of that attribute. If one or more
// option-value pairs are specified, then the command modifies the given named
// font to have the given values; in this case, all widgets using that font
// will redisplay themselves using the new attributes for the font. See FONT
// OPTIONS below for a list of the possible attributes.
//
// Note that on Aqua/macOS, the system fonts (see PLATFORM SPECIFIC FONTS
// below) may not be actually altered because they are implemented by the
// system theme. To achieve the effect of modification, use font actual to get
// their configuration and font create to synthesize a copy of the font which
// can be modified.
//
// More information might be available at the [Tcl/Tk font] page.
//
// [Tcl/Tk font]: https://www.tcl.tk/man/tcl9.0/TkCmd/font.html
func FontConfigure(name string, options ...any) []string {
	for i, v := range options {
		if s := funcToTclOption(v); s != "" {
			options[i] = rawOption(s)
		}
	}
	return parseList(evalErr(fmt.Sprintf("font configure %s %s", tclSafeString(name), collectAny(options...))))
}

// Data option.
//
// Known uses:
//   - [NewBitmap] (command specific)
//   - [NewPhoto] (command specific)
func Data(val any) Opt {
	return rawOption(fmt.Sprintf(`-data %s`, optionString(val)))
}

// winfo — Return window-related information
//
// # Description
//
// Returns a slice of all the children of window. Top-level windows are
// returned as children of their logical parents. The list is in stacking
// order, with the lowest window first, except for Top-level windows which are
// not returned in stacking order. Use the wm stackorder command to query the
// stacking order of Top-level windows.
//
// More information might be available at the [Tcl/Tk winfo] page.
//
// [Tcl/Tk winfo]: https://www.tcl.tk/man/tcl9.0/TkCmd/winfo.html
func WinfoChildren(w *Window) (r []*Window) {
	for _, v := range parseList(evalErr(fmt.Sprintf("winfo children %s", w))) {
		if w := windowIndex[v]; w != nil {
			r = append(r, w)
		}
	}
	return r
}

// winfo — Return window-related information
//
// # Description
//
// Returns a decimal string giving the height of window's screen, in pixels.
//
// More information might be available at the [Tcl/Tk winfo] page.
//
// [Tcl/Tk winfo]: https://www.tcl.tk/man/tcl9.0/TkCmd/winfo.html
func WinfoScreenHeight(w *Window) string {
	return evalErr(fmt.Sprintf("winfo screenheight %s", w))
}

// winfo — Return window-related information
//
// # Description
//
// Returns a decimal string giving the width of window's screen, in pixels.
//
// More information might be available at the [Tcl/Tk winfo] page.
//
// [Tcl/Tk winfo]: https://www.tcl.tk/man/tcl9.0/TkCmd/winfo.html
func WinfoScreenWidth(w *Window) string {
	return evalErr(fmt.Sprintf("winfo screenwidth %s", w))
}

// winfo — Return window-related information
//
// # Description
//
// Returns a decimal string giving window's height in pixels. When a window is
// first created its height will be 1 pixel; the height will eventually be
// changed by a geometry manager to fulfil the window's needs. If you need the
// true height immediately after creating a widget, invoke update to force the
// geometry manager to arrange it, or use winfo reqheight to get the window's
// requested height instead of its actual height.
//
// More information might be available at the [Tcl/Tk winfo] page.
//
// [Tcl/Tk winfo]: https://www.tcl.tk/man/tcl9.0/TkCmd/winfo.html
func WinfoHeight(w *Window) string {
	return evalErr(fmt.Sprintf("winfo height %s", w))
}

// winfo — Return window-related information
//
// # Description
//
// Returns a decimal string giving window's width in pixels. When a window is
// first created its width will be 1 pixel; the width will eventually be
// changed by a geometry manager to fulfil the window's needs. If you need the
// true width immediately after creating a widget, invoke update to force the
// geometry manager to arrange it, or use winfo reqwidth to get the window's
// requested width instead of its actual width.
//
// More information might be available at the [Tcl/Tk winfo] page.
//
// [Tcl/Tk winfo]: https://www.tcl.tk/man/tcl9.0/TkCmd/winfo.html
func WinfoWidth(w *Window) string {
	return evalErr(fmt.Sprintf("winfo width %s", w))
}

// tooltip — Tooltip management
//
// # Description
//
// Prevents the specified widgets from showing tooltips. pattern is a glob
// pattern and defaults to matching all widgets.
//
// More information might be available at the [Tklib tooltip] page.
//
// [Tklib tooltip]: https://core.tcl-lang.org/tklib/doc/trunk/embedded/md/tklib/files/modules/tooltip/tooltip.md
func TooltipClear(pattern string) {
	s := ""
	if pattern != "" {
		s = tclSafeString(pattern)
	}
	evalErr(fmt.Sprintf("tooltip::tooltip clear %s", s))
}

// tooltip — Tooltip management
//
// # Description
//
// Queries or modifies the configuration options of the tooltip. The supported
// options are -backgroud, -foreground and -font. If one option is specified with
// no value, returns the value of that option. Otherwise, sets the given
// options to the corresponding values.
//
// More information might be available at the [Tklib tooltip] page.
//
// [Tklib tooltip]: https://core.tcl-lang.org/tklib/doc/trunk/embedded/md/tklib/files/modules/tooltip/tooltip.md
func TooltipConfigure(options ...any) string {
	return evalErr(fmt.Sprintf("tooltip::tooltip configure %s", collectAny(options)))
}

// tooltip — Tooltip management
//
// # Description
//
// Query or set the hover delay. This is the interval that the pointer must
// remain over the widget before the tooltip is displayed. The delay is
// specified in milliseconds and must be greater than or equal to 50 ms. With
// a negative argument the current delay is returned.
//
// More information might be available at the [Tklib tooltip] page.
//
// [Tklib tooltip]: https://core.tcl-lang.org/tklib/doc/trunk/embedded/md/tklib/files/modules/tooltip/tooltip.md
func TooltipDelay(delay time.Duration) string {
	s := ""
	if delay >= 0 {
		s = optionString(delay)
	}
	return evalErr(fmt.Sprintf("tooltip::tooltip delay %s", s))
}

// tooltip — Tooltip management
//
// # Description
//
// Enable or disable fading of the tooltip. The fading is enabled by default on
// Win32 and Aqua. The tooltip will fade away on Leave events instead
// disappearing.
//
// More information might be available at the [Tklib tooltip] page.
//
// [Tklib tooltip]: https://core.tcl-lang.org/tklib/doc/trunk/embedded/md/tklib/files/modules/tooltip/tooltip.md
func TooltipFade(v bool) string {
	return evalErr(fmt.Sprintf("tooltip::tooltip fade %v", v))
}

// tooltip — Tooltip management
//
// # Description
//
// # Disable all tooltips
//
// More information might be available at the [Tklib tooltip] page.
//
// [Tklib tooltip]: https://core.tcl-lang.org/tklib/doc/trunk/embedded/md/tklib/files/modules/tooltip/tooltip.md
func TooltipOff(v bool) string {
	return evalErr("tooltip::tooltip off")
}

// tooltip — Tooltip management
//
// # Description
//
// Enables tooltips for defined widgets.
//
// More information might be available at the [Tklib tooltip] page.
//
// [Tklib tooltip]: https://core.tcl-lang.org/tklib/doc/trunk/embedded/md/tklib/files/modules/tooltip/tooltip.md
func TooltipOn(v bool) string {
	return evalErr("tooltip::tooltip on")
}

// tooltip — Tooltip management
//
// # Description
//
// This command arranges for widget 'w' to display a tooltip with a
// message.
//
// If the specified widget is a menu, canvas, listbox, ttk::treeview,
// ttk::notebook or text widget then additional options are used to tie the
// tooltip to specific menu, canvas or listbox items, ttk::treeview items or
// column headings, ttk::notebook tabs, or text widget tags.
//
//   - [Heading] columnId: This option is used to set a tooltip for a
//     ttk::treeview column heading. The column does not need to already exist.
//     You should not use the same identifiers for columns and items in a widget
//     for which you are using tooltips as their tooltips will be mixed. The
//     widget must be a ttk::treeview widget.
//
//   - [Image] image: The specified (photo) image will be displayed to the left
//     of the primary tooltip message.
//
//   - [Index] index: This option is used to set a tooltip on a menu item. The
//     index may be either the entry index or the entry label. The widget must be
//     a menu widget but the entries do not have to exist when the tooltip is
//     set.
//
//   - [Info] info: The specified info text will be displayed as additional
//     information below the primary tooltip message.
//
//   - [Items] items: This option is used to set a tooltip for canvas, listbox
//     or ttk::treview items. For the canvas widget, the item must already be
//     present in the canvas and will be found with a find withtag lookup. For
//     listbox and ttk::treview widgets the item(s) may be created later but the
//     programmer is responsible for managing the link between the listbox or
//     ttk::treview item index and the corresponding tooltip. If the listbox or
//     ttk::treview items are re-ordered, the tooltips will need amending.
//
//     If the widget is not a canvas, listbox or ttk::treview then an error is
//     raised.
//
//   - [Tab] tabId: The -tab option can be used to set a tooltip for a
//     ttk::notebook tab. The tab should already be present when this command is
//     called, or an error will be returned. The widget must be a ttk::notebook
//     widget.
//
//   - [Tag] name: The -tag option can be used to set a tooltip for a text
//     widget tag. The tag should already be present when this command is called,
//     or an error will be returned. The widget must be a text widget.
//
//   - "--": The -- option marks the end of options. The argument following
//     this one will be treated as message even if it starts with a -.
//
// Tooltip returns 'w'.
//
// More information might be available at the [Tklib tooltip] page.
//
// [Tklib tooltip]: https://core.tcl-lang.org/tklib/doc/trunk/embedded/md/tklib/files/modules/tooltip/tooltip.md
func Tooltip(w Widget, options ...any) (r Widget) {
	evalErr(fmt.Sprintf("tooltip::tooltip %s %s", w, collectAny(options...)))
	return w
}

// Heading option.
//
// Known uses:
//   - [Tooltip] (command specific)
func Heading(columnId string) Opt {
	return rawOption(fmt.Sprintf(`-heading %s`, optionString(columnId)))
}

// Index option.
//
// Known uses:
//   - [Tooltip] (command specific)
func Index(index any) Opt {
	return rawOption(fmt.Sprintf(`-index %s`, optionString(index)))
}

// Info option.
//
// Known uses:
//   - [Tooltip] (command specific)
func Info(info string) Opt {
	return rawOption(fmt.Sprintf(`-info %s`, optionString(info)))
}

// Items option.
//
// Known uses:
//   - [Tooltip] (command specific)
func Items(items ...any) Opt {
	return rawOption(fmt.Sprintf(`-items %s`, collectAny(items...)))
}

// Tab option.
//
// Known uses:
//   - [Tooltip] (command specific)
func Tab(tabId any) Opt {
	return rawOption(fmt.Sprintf(`-tab %s`, collectAny(tabId)))
}

// Tag option.
//
// Known uses:
//   - [Tooltip] (command specific)
//   - [TextWidget] (command specific)
func Tag(name string) Opt {
	return rawOption(fmt.Sprintf(`-tag %s`, optionString(name)))
}

// ttk::combobox — text field with popdown selection list
//
// # Description
//
// If newIndex is supplied, sets the combobox value to the element at position
// newIndex in the list of -values (in addition to integers, the end index is
// supported and indicates the last element of the list, moreover the same
// simple interpretation as for the command string index is supported, with
// simple integer index arithmetic and indexing relative to end). Otherwise,
// returns the index of the current value in the list of -values or {} if the
// current value does not appear in the list.
//
// More information might be available at the [Tcl/Tk combobox] page.
//
// [Tcl/Tk combobox]: https://tcl.tk/man/tcl9.0/TkCmd/ttk_combobox.html
func (w *TComboboxWidget) Current(newIndex any) (r string) {
	var arg string
	if newIndex != nil {
		arg = tclSafeString(fmt.Sprint(newIndex))
	}
	return evalErr(fmt.Sprintf("%s current %s", w, arg))
}

// ttk::treeview — hierarchical multicolumn data display widget
//
// # Description
//
// Returns the integer index of item within its parent's list of children or -1
// otherwise.
//
// More information might be available at the [Tcl/Tk treeview] page.
//
// [Tcl/Tk treeview]: https://tcl.tk/man/tcl9.0/TkCmd/ttk_treeview.html
func (w *TTreeviewWidget) Index(item any) (r int) {
	s := evalErr(fmt.Sprintf("%s index %s", w, tclSafeString(fmt.Sprint(item))))
	if n, err := strconv.ParseInt(s, 10, 32); err == nil {
		return int(n)
	}

	return -1
}

// ttk::treeview — hierarchical multicolumn data display widget
//
// # Description
//
// If newchildren is not specified, returns the list of children belonging to
// item.
//
// If newchildren is specified, replaces item's child list with newchildren.
// Items in the old child list not present in the new child list are detached
// from the tree. None of the items in newchildren may be an ancestor of item.
//
// More information might be available at the [Tcl/Tk treeview] page.
//
// [Tcl/Tk treeview]: https://tcl.tk/man/tcl9.0/TkCmd/ttk_treeview.html
func (w *TTreeviewWidget) Children(item any, newChildren ...any) (r []string) {
	switch {
	case len(newChildren) == 0:
		return parseList(evalErr(fmt.Sprintf("%s children %s", w, tclSafeString(fmt.Sprint(item)))))
	default:
		return parseList(evalErr(fmt.Sprintf("%s children %s {%s}", w, tclSafeString(fmt.Sprint(item)), tclSafeList(flat(newChildren...)))))
	}
}

// ttk::treeview — hierarchical multicolumn data display widget
//
// # Description
//
// Returns one of:
//
//   - heading
//
//     Tree heading area; use [pathname identify column x y] to determine the
//     heading number.
//
//   - separator
//
//     Space between two column headings; [pathname identify column x y] will
//     return the display column identifier of the heading to left of the
//     separator.
//
//   - tree
//
//     The tree area.
//
//   - cell
//
//     A data cell.
//
// More information might be available at the [Tcl/Tk treeview] page.
//
// [Tcl/Tk treeview]: https://tcl.tk/man/tcl9.0/TkCmd/ttk_treeview.html
func (w *TTreeviewWidget) IdentifyRegion(x, y int) (r string) {
	return evalErr(fmt.Sprintf("%s identify region %v %v", w, x, y))
}

// ttk::treeview — hierarchical multicolumn data display widget
//
// # Description
//
// Returns the item ID of the item at position x y.
//
// More information might be available at the [Tcl/Tk treeview] page.
//
// [Tcl/Tk treeview]: https://tcl.tk/man/tcl9.0/TkCmd/ttk_treeview.html
func (w *TTreeviewWidget) IdentifyItem(x, y int) (r string) {
	return evalErr(fmt.Sprintf("%s identify item %v %v", w, x, y))
}

// ttk::treeview — hierarchical multicolumn data display widget
//
// # Description
//
// Returns the display column identifier of the cell at position x. The tree
// column has ID #0.
//
// More information might be available at the [Tcl/Tk treeview] page.
//
// [Tcl/Tk treeview]: https://tcl.tk/man/tcl9.0/TkCmd/ttk_treeview.html
func (w *TTreeviewWidget) IdentifyColumn(x, y int) (r string) {
	return evalErr(fmt.Sprintf("%s identify column %v %v", w, x, y))
}

// ttk::treeview — hierarchical multicolumn data display widget
//
// # Description
//
// Returns the cell identifier of the cell at position x, y. A cell identifier
// is a list of item ID and column ID.
//
// More information might be available at the [Tcl/Tk treeview] page.
//
// [Tcl/Tk treeview]: https://tcl.tk/man/tcl9.0/TkCmd/ttk_treeview.html
func (w *TTreeviewWidget) IdentifyCell(x, y int) (r []string) {
	return parseList(evalErr(fmt.Sprintf("%s identify cell %v %v", w, x, y)))
}

// ttk::treeview — hierarchical multicolumn data display widget
//
// # Description
//
// Returns the element at position x, y.
//
// More information might be available at the [Tcl/Tk treeview] page.
//
// [Tcl/Tk treeview]: https://tcl.tk/man/tcl9.0/TkCmd/ttk_treeview.html
func (w *TTreeviewWidget) IdentifyElement(x, y int) (r string) {
	return evalErr(fmt.Sprintf("%s identify element %v %v", w, x, y))
}

// ttk::treeview — hierarchical multicolumn data display widget
//
// # Description
//
// Deletes each of the items in itemList and all of their descendants. The root
// item may not be deleted. See also: detach.
//
// More information might be available at the [Tcl/Tk treeview] page.
//
// [Tcl/Tk treeview]: https://tcl.tk/man/tcl9.0/TkCmd/ttk_treeview.html
func (w *TTreeviewWidget) Delete(itemList ...any) {
	itemList = flat(itemList...)
	if len(itemList) == 0 {
		return
	}

	evalErr(fmt.Sprintf("%s delete {%v}", w, tclSafeList(itemList...)))
}

func flat(list ...any) (r []any) {
	for _, v := range list {
		switch x := v.(type) {
		case []string:
			for _, v := range x {
				r = append(r, v)
			}
		case []any:
			for _, v := range x {
				r = append(r, flat(v))
			}
		default:
			r = append(r, v)
		}
	}
	return r
}

// ttk::treeview — hierarchical multicolumn data display widget
//
// # Description
//
// Returns the ID of the parent of item, or "" if item is at the top level of
// the hierarchy.
//
// More information might be available at the [Tcl/Tk treeview] page.
//
// [Tcl/Tk treeview]: https://tcl.tk/man/tcl9.0/TkCmd/ttk_treeview.html
func (w *TTreeviewWidget) Parent(item any) (r string) {
	if r = evalErr(fmt.Sprintf("%s parent %s", w, tclSafeString(fmt.Sprint(item)))); r == "{}" {
		r = ""
	}
	return r
}

// ttk::treeview — hierarchical multicolumn data display widget
//
// # Description
//
// Query or modify the options for the specified column. If no -option is
// specified, returns a dictionary of option/value pairs. If a single -option
// is specified, returns the value of that option. Otherwise, the options are
// updated with the specified values. The following options may be set on each
// column:
//
//   - id name:
//     The column name. This is a read-only option. For example, [$pathname
//     column #n -id] returns the data column associated with display column n.
//     The tree column has -id #0.
//   - anchor anchor:
//     Specifies how the text in this column should be aligned with respect to
//     the cell. Anchor is one of n, ne, e, se, s, sw, w, nw, or center.
//   - minwidth minwidth:
//     The minimum width of the column in pixels. The treeview widget will not
//     make the column any smaller than -minwidth when the widget is resized or
//     the user drags a heading column separator. Default is 20 pixels.
//   - separator boolean:
//     Specifies whether or not a column separator should be drawn to the right
//     of the column. Default is false.
//   - stretch boolean:
//     Specifies whether or not the column width should be adjusted when the
//     widget is resized or the user drags a heading column separator. Boolean
//     may have any of the forms accepted by Tcl_GetBoolean. By default columns
//     are stretchable.
//     -width width:
//     The width of the column in pixels. Default is 200 pixels. The specified
//     column width may be changed by Tk in order to honor -stretch and/or
//     -minwidth, or when the widget is resized or the user drags a heading
//     column separator.
//
// Use pathname "#0" to configure the tree column.
//
// More information might be available at the [Tcl/Tk treeview] page.
//
// [Tcl/Tk treeview]: https://tcl.tk/man/tcl9.0/TkCmd/ttk_treeview.html
func (w *TTreeviewWidget) Column(column any, options ...Opt) (r string) {
	return evalErr(fmt.Sprintf("%s column %s %s", w, tclSafeString(fmt.Sprint(column)), collect(options...)))
}

// ttk::treeview — hierarchical multicolumn data display widget
//
// # Description
//
// Query or modify the heading options for the specified column. Valid options
// are:
//
//   - text text:
//     The text to display in the column heading.
//   - image imageName:
//     Specifies an image to display to the right of the column heading.
//   - anchor anchor:
//     Specifies how the heading text should be aligned. One of the standard Tk anchor values.
//   - command script:
//     A script to evaluate when the heading label is pressed.
//
// Use pathname heading "#0" to configure the tree column heading.
//
// More information might be available at the [Tcl/Tk treeview] page.
//
// [Tcl/Tk treeview]: https://tcl.tk/man/tcl9.0/TkCmd/ttk_treeview.html
func (w *TTreeviewWidget) Heading(column any, options ...Opt) (r string) {
	return evalErr(fmt.Sprintf("%s heading %s %s", w, tclSafeString(fmt.Sprint(column)), collect(options...)))
}

// ttk::treeview — hierarchical multicolumn data display widget
//
// # Description
//
// Creates a new item. parent is the item ID of the parent item, or the empty
// string {} to create a new top-level item. index is an integer, or the value
// end, specifying where in the list of parent's children to insert the new
// item. If index is negative or zero, the new node is inserted at the
// beginning; if index is greater than or equal to the current number of
// children, it is inserted at the end. If -id is specified, it is used as the
// item identifier; id must not already exist in the tree. Otherwise, a new
// unique identifier is generated.
//
// Insert returns the item identifier of the newly created item. See
// ITEM OPTIONS for the list of available options.
//
// More information might be available at the [Tcl/Tk treeview] page.
//
// [Tcl/Tk treeview]: https://tcl.tk/man/tcl9.0/TkCmd/ttk_treeview.html
func (w *TTreeviewWidget) Insert(parent, index any, options ...Opt) (r string) {
	return evalErr(fmt.Sprintf("%s insert %s %s %s", w, tclSafeString(fmt.Sprint(parent)), tclSafeString(fmt.Sprint(index)), collect(options...)))
}

// Id option.
//
// Known uses:
//   - [TTreeviewWidget] (command specific)
//
// More information might be available at the [Tcl/Tk treeview] page.
//
// [Tcl/Tk treeview]: https://tcl.tk/man/tcl9.0/TkCmd/ttk_treeview.html
func Id(val any) Opt {
	return rawOption(fmt.Sprintf(`-id %s`, optionString(val)))
}

// Open option.
//
// Known uses:
//   - [TTreeviewWidget] (command specific)
//
// More information might be available at the [Tcl/Tk treeview] page.
//
// [Tcl/Tk treeview]: https://tcl.tk/man/tcl9.0/TkCmd/ttk_treeview.html
func Open(val bool) Opt {
	return rawOption(fmt.Sprintf(`-open %v`, val))
}

// ttk::treeview — hierarchical multicolumn data display widget
//
// # Description
//
// Query or modify the options for the specified item. If no -option is
// specified, returns a dictionary of option/value pairs. If a single -option
// is specified, returns the value of that option. Otherwise, the item's
// options are updated with the specified values. See ITEM OPTIONS for the list
// of available options.  pathname move item parent index
//
// More information might be available at the [Tcl/Tk treeview] page.
//
// [Tcl/Tk treeview]: https://tcl.tk/man/tcl9.0/TkCmd/ttk_treeview.html
func (w *TTreeviewWidget) Item(item any, options ...any) (r string) {
	if len(options) == 1 {
		if s := funcToTclOption(options[0]); s != "" {
			return evalErr(fmt.Sprintf("%s item %s %s", w, tclSafeString(fmt.Sprint(item)), s))
		}
	}

	return evalErr(fmt.Sprintf("%s item %s %s", w, tclSafeString(fmt.Sprint(item)), collectAny(options...)))
}

// ttk::treeview — hierarchical multicolumn data display widget
//
// # Description
//
// Adds the specified tag to each of the listed items. If tag is already
// present for a particular item, then the -tags for that item are unchanged.
//
// More information might be available at the [Tcl/Tk treeview] page.
//
// [Tcl/Tk treeview]: https://tcl.tk/man/tcl9.0/TkCmd/ttk_treeview.html
func (w *TTreeviewWidget) TagAdd(tag string, items ...any) {
	evalErr(fmt.Sprintf("%s tag add %s %s", w, tclSafeString(tag), collectAny(items...)))
}

// ttk::treeview — hierarchical multicolumn data display widget
//
// # Description
//
// If item is specified, sets the focus item to item. Otherwise, returns the
// current focus item, or {} if there is none.
//
// More information might be available at the [Tcl/Tk treeview] page.
//
// [Tcl/Tk treeview]: https://tcl.tk/man/tcl9.0/TkCmd/ttk_treeview.html
func (w *TTreeviewWidget) Focus(item ...any) (r string) {
	switch len(item) {
	case 0:
		return evalErr(fmt.Sprintf("%s focus", w))
	default:
		return evalErr(fmt.Sprintf("%s focus %s", w, tclSafeString(fmt.Sprint(item[0]))))
	}
}

// ttk::treeview — hierarchical multicolumn data display widget
//
// # Description
//
// Query or modify the options for the specified tagName. If one or more
// option/value pairs are specified, sets the value of those options for the
// specified tag. If a single option is specified, returns the value of that
// option (or the empty string if the option has not been specified for
// tagName). With no additional arguments, returns a dictionary of the option
// settings for tagName. See TAG OPTIONS for the list of available options.
//
// More information might be available at the [Tcl/Tk treeview] page.
//
// [Tcl/Tk treeview]: https://tcl.tk/man/tcl9.0/TkCmd/ttk_treeview.html
func (w *TTreeviewWidget) TagConfigure(tag string, options ...any) (r string) {
	if len(options) == 1 {
		if s := funcToTclOption(options[0]); s != "" {
			return evalErr(fmt.Sprintf("%s tag configure %s %s", w, tclSafeString(fmt.Sprint(tag)), s))
		}
	}

	return evalErr(fmt.Sprintf("%s tag configure %s %s", w, tclSafeString(tag), collectAny(options...)))
}

// ttk::treeview — hierarchical multicolumn data display widget
//
// # Description
//
// Manages item selection. Item selection is independent from cell selection
// handled by the cellselection command. If selop is not specified, returns the
// list of selected items. Otherwise, selop is one of the following:
//
//   - set itemList:
//     itemList becomes the new selection.
//   - add itemList:
//     Add itemList to the selection.
//   - remove itemList:
//     Remove itemList from the selection.
//   - toggle itemList:
//     Toggle the selection state of each item in itemList.
//
// More information might be available at the [Tcl/Tk treeview] page.
//
// [Tcl/Tk treeview]: https://tcl.tk/man/tcl9.0/TkCmd/ttk_treeview.html
func (w *TTreeviewWidget) Selection(selop string, itemList ...any) (r []string) {
	var s string
	if len(itemList) != 0 {
		s = fmt.Sprintf("{%s}", tclSafeList(flat(itemList...)...))
	}
	if selop == "" {
		return parseList(evalErr(fmt.Sprintf("%s selection %s", w, s)))
	}

	evalErr(fmt.Sprintf("%s selection %s %s", w, tclSafeString(fmt.Sprint(selop)), s))
	return nil
}

// ttk::treeview — hierarchical multicolumn data display widget
//
// # Description
//
// Ensure that item is visible: sets all of item's ancestors to -open true, and
// scrolls the widget if necessary so that item is within the visible portion
// of the tree.
//
// More information might be available at the [Tcl/Tk treeview] page.
//
// [Tcl/Tk treeview]: https://tcl.tk/man/tcl9.0/TkCmd/ttk_treeview.html
func (w *TTreeviewWidget) See(item any) {
	evalErr(fmt.Sprintf("%s see %s", w, tclSafeString(fmt.Sprint(item))))
}

// ttk::scale — Create and manipulate a scale widget
//
// # Description
//
// Get the current value of the -value option, or the value corresponding to
// the coordinates x,y if they are specified. X and y are pixel coordinates
// relative to the scale widget origin.
//
// More information might be available at the [Tcl/Tk scale] page.
//
// [Tcl/Tk scale]: https://tcl.tk/man/tcl9.0/TkCmd/ttk_scale.html
func (w *TScaleWidget) Get(xy ...any) (r string) {
	arg := ""
	if len(xy) >= 2 {
		arg = fmt.Sprintf("%s %s", tclSafeString(fmt.Sprint(xy[0])), tclSafeString(fmt.Sprint(xy[1])))
	}
	return evalErr(fmt.Sprintf("%s get %s", w, arg))
}

// tk_optionMenu — Create an option menubutton and its menu
//
// # Description
//
// This procedure creates an option menubutton whose name is pathName, plus an
// associated menu. Together they allow the user to select one of the values
// given by the value arguments. The current value will be stored in the global
// variable whose name is given by varName and it will also be displayed as the
// label in the option menubutton. The user can click on the menubutton to
// display a menu containing all of the values and thereby select a new value.
// Once a new value is selected, it will be stored in the variable and appear
// in the option menubutton. The current value can also be changed by setting
// the variable.
//
// The return value from tk_optionMenu is the name of the menu associated with
// pathName, so that the caller can change its configuration options or
// manipulate it in other ways.
//
// More information might be available at the [Tcl/Tk option menu] page.
//
// [Tcl/Tk option menu]: https://tcl.tk/man/tcl9.0/TkCmd/optionMenu.html
func OptionMenu(varName *VariableOpt, options ...any) (r *OptionMenuWidget) {
	return App.OptionMenu(varName, options...)
}

// tk_optionMenu — Create an option menubutton and its menu
//
// The resulting [Widget] is a child of 'w'
//
// For details please see [Button]
func (w *Window) OptionMenu(varName *VariableOpt, options ...any) (r *OptionMenuWidget) {
	r = &OptionMenuWidget{Window: w.newChild0("optionmenu")}
	for i, v := range options {
		if s := fmt.Sprint(v); s == "" {
			options[i] = "{}"
		}
	}
	r.name = evalErr(fmt.Sprintf("tk_optionMenu %s %s %s", r, varName.tclName, tclSafeList(options...)))
	return r
}

// OptionMenuWidget represents the Tcl/Tk option menu.
//
// More information might be available at the [Tcl/Tk option menu] page.
//
// [Tcl/Tk option menu]: https://tcl.tk/man/tcl9.0/TkCmd/optionMenu.html
type OptionMenuWidget struct {
	*Window
	name string
}

// Name returns the menu name of 'w'.
func (w *OptionMenuWidget) Name() string {
	return w.name
}

// grab — Confine pointer and keyboard events to a window sub-tree
//
// # Description
//
// Same as GrabSet().
//
// More information might be available at the [Tcl/Tk grab] page.
//
// [Tcl/Tk grab]: https://tcl.tk/man/tcl9.0/TkCmd/grab.html
func Grab(options ...Opt) {
	GrabSet(options...)
}

// grab — Confine pointer and keyboard events to a window sub-tree
//
// # Description
//
// If window is specified, returns the name of the current grab window in this
// application for window's display, or an empty string if there is no such
// window. If window is omitted, the command returns a list whose elements are
// all of the windows grabbed by this application for all displays, or an empty
// string if the application has no grabs.
//
// More information might be available at the [Tcl/Tk grab] page.
//
// [Tcl/Tk grab]: https://tcl.tk/man/tcl9.0/TkCmd/grab.html
func GrabCurrent(options ...Opt) []string {
	return parseList(evalErr(fmt.Sprintf("grab current %s", collect(options...))))
}

// grab — Confine pointer and keyboard events to a window sub-tree
//
// # Description
//
// Releases the grab on window if there is one, otherwise does nothing.
//
// More information might be available at the [Tcl/Tk grab] page.
//
// [Tcl/Tk grab]: https://tcl.tk/man/tcl9.0/TkCmd/grab.html
func GrabRelease(w Opt) {
	evalErr(fmt.Sprintf("grab release %s", w))
}

// grab — Confine pointer and keyboard events to a window sub-tree
//
// # Description
//
// Sets a grab on window. If -global is specified then the grab is global,
// otherwise it is local. If a grab was already in effect for this application
// on window's display then it is automatically released. If there is already a
// grab on window and it has the same global/local form as the requested grab,
// then the command does nothing.
//
// More information might be available at the [Tcl/Tk grab] page.
//
// [Tcl/Tk grab]: https://tcl.tk/man/tcl9.0/TkCmd/grab.html
func GrabSet(options ...Opt) {
	evalErr(fmt.Sprintf("grab set %s", collect(options...)))
}

// grab — Confine pointer and keyboard events to a window sub-tree
//
// # Description
//
// Returns none if no grab is currently set on window, local if a local grab is
// set on window, and global if a global grab is set.
//
// More information might be available at the [Tcl/Tk grab] page.
//
// [Tcl/Tk grab]: https://tcl.tk/man/tcl9.0/TkCmd/grab.html
func GrabStatus(w Opt) string {
	return evalErr(fmt.Sprintf("grab status %s", w))
}

// Global option.
//
// Known uses:
//   - [Grab] (command specific)
//   - [GrabSet] (command specific)
func Global() Opt {
	return rawOption("-global")
}

// after — Execute a command after a time delay
//
// # Description
//
// If script is not specified, the command sleeps for duration 'ms' and then
// returns. Negative duration 'ms' is treated as zero. While the command is
// sleeping the application does not respond to events.
//
// If script is specified, the command returns immediately, but it arranges for
// a Tcl command to be executed ms milliseconds later as an event handler. The
// command will be executed exactly once, at the given time. The command will
// be executed at global level (outside the context of any Tcl procedure). If
// an error occurs while executing the delayed command then the background
// error will be reported by the command registered with interp bgerror. The
// after command returns an identifier that can be used to cancel the delayed
// command using after cancel. A ms value of 0 (or negative) queues the event
// immediately with priority over other event types (if not installed withn an
// event proc, which will wait for next round of events).
//
// More information might be available at the [Tcl/Tk after] page.
//
// [Tcl/Tk grab]: https://tcl.tk/man/tcl9.0/TkCmd/after.html
func TclAfter(ms time.Duration, script ...any) string {
	switch {
	case len(script) == 0:
		return evalErr(fmt.Sprintf("after %v", optionString(ms)))
	default:
		return evalErr(fmt.Sprintf("after %v %s", optionString(ms), newEventHandler("", script[0]).optionString(nil)))
	}
}

// after — Execute a command after a time delay
//
// # Description
//
// Cancels the execution of a delayed command that was previously scheduled. Id
// indicates which command should be canceled; it must have been the return
// value from a previous after command. If the command given by id has already
// been executed then the after cancel command has no effect.
//
// More information might be available at the [Tcl/Tk after] page.
//
// [Tcl/Tk grab]: https://tcl.tk/man/tcl9.0/TkCmd/after.html
func TclAfterCancel(id string) {
	evalErr(fmt.Sprintf("after cancel %s", tclSafeString(id)))
}

// after — Execute a command after a time delay
//
// # Description
//
// Concatenates the script arguments together with space separators (just as in
// the concat command), and arranges for the resulting script to be evaluated
// later as an idle callback. The script will be run exactly once, the next
// time the event loop is entered and there are no events to process. The
// command returns an identifier that can be used to cancel the delayed command
// using after cancel. If an error occurs while executing the script then the
// background error will be reported by the command registered with interp
// bgerror.
//
// More information might be available at the [Tcl/Tk after] page.
//
// [Tcl/Tk grab]: https://tcl.tk/man/tcl9.0/TkCmd/after.html
func TclAfterIdle(script any) string {
	switch x := script.(type) {
	case nil:
		return ""
	case *eventHandler:
		x.tcl = ""
		return evalErr(fmt.Sprintf("after idle %s", collect(x)))
	default:
		return evalErr(fmt.Sprintf("after idle %s", newEventHandler("", script).optionString(nil)))
	}
}

// A helper for ActivateTheme/ActivateExtension.
func isCalledFromMain() bool {
	pcs := make([]uintptr, 20)
	n := runtime.Callers(3, pcs)
	frames := runtime.CallersFrames(pcs[:n])
	for {
		frame, more := frames.Next()
		if strings.HasPrefix(frame.Function, "main.") {
			return true
		}

		if !more {
			return false
		}
	}
}

// Entry — Create and manipulate 'entry' one-line text entry widgets
//
// # Description
//
// Arrange for the insertion cursor to be displayed just before the character
// given by index. Returns an empty string.
//
// More information might be available at the [Tcl/Tk entry] page.
//
// [Tcl/Tk entry]: https://www.tcl.tk/man/tcl9.0/TkCmd/entry.html
func (w *EntryWidget) Icursor(index any) (r string) {
	return evalErr(fmt.Sprintf("%s icursor %s", w, tclSafeString(fmt.Sprint(index))))
}

// TEntry — Editable text field widget
//
// # Description
//
// Arrange for the insertion cursor to be displayed just before the character
// given by index. Returns an empty string.
//
// More information might be available at the [Tcl/Tk ttk_entry] page.
//
// [Tcl/Tk ttk_entry]: https://www.tcl.tk/man/tcl9.0/TkCmd/ttk_entry.html
func (w *TEntryWidget) Icursor(index any) (r string) {
	return evalErr(fmt.Sprintf("%s icursor %s", w, tclSafeString(fmt.Sprint(index))))
}

// listbox — Create and manipulate 'listbox' item list widgets
//
// # Description
//
// If last is omitted, returns the contents of the listbox element indicated by
// first, or an empty string if first refers to a non-existent element. If last
// is specified, the command returns a list whose elements are all of the
// listbox elements between first and last, inclusive. Both first and last may
// have any of the standard forms for indices.
//
// More information might be available at the [Tcl/Tk listbox] page.
//
// [Tcl/Tk listbox]: https://www.tcl.tk/man/tcl9.0/TkCmd/listbox.html
func (w *ListboxWidget) Get(first any, last ...any) (r []string) {
	switch len(last) {
	case 0:
		return []string{evalErr(fmt.Sprintf("%s get %s", w, tclSafeString(fmt.Sprint(first))))}
	default:
		return parseList(evalErr(fmt.Sprintf("%s get %s %s", w, tclSafeString(fmt.Sprint(first)), tclSafeString(fmt.Sprint(last[0])))))
	}
}

// listbox — Create and manipulate 'listbox' item list widgets
//
// # Description
//
// Inserts zero or more new elements in the list just before the element given
// by index. If index is specified as end then the new elements are added to
// the end of the list.
//
// More information might be available at the [Tcl/Tk listbox] page.
//
// [Tcl/Tk listbox]: https://www.tcl.tk/man/tcl9.0/TkCmd/listbox.html
func (w *ListboxWidget) Insert(index any, elements ...any) {
	evalErr(fmt.Sprintf("%s insert %s %s", w, tclSafeString(fmt.Sprint(index)), tclSafeList(elements...)))
}

// listbox — Create and manipulate 'listbox' item list widgets
//
// # Description
//
// Returns a list containing the numerical indices of all of the elements in the
// listbox that are currently selected. If there are no elements selected in the
// listbox then a zero length slice is returned.
//
// More information might be available at the [Tcl/Tk listbox] page.
//
// [Tcl/Tk listbox]: https://www.tcl.tk/man/tcl9.0/TkCmd/listbox.html
func (w *ListboxWidget) Curselection() (r []int) {
	indicesStrings := parseList(evalErr(fmt.Sprintf("%s curselection", w)))
	indices := make([]int, len(indicesStrings))
	for i, index := range indicesStrings {
		indices[i] = atoi(index)
	}
	return indices
}

// ttk::progressbar — Provide progress feedback
//
// # Description
//
// Begin autoincrement mode: schedules a recurring timer event that calls step
// every interval milliseconds. If omitted, interval defaults to 50
// milliseconds (20 steps/second).
//
// More information might be available at the [Tcl/Tk progressbar] page.
//
// [Tcl/Tk progressbar]: https://www.tcl.tk/man/tcl9.0/TkCmd/ttk_progressbar.html
func (w *TProgressbarWidget) Start(interval ...time.Duration) {
	s := ""
	if len(interval) != 0 {
		s = fmt.Sprint(int(interval[0] / time.Millisecond))
	}
	evalErr(fmt.Sprintf("%s start %s", w, s))
}

// ttk::progressbar — Provide progress feedback
//
// # Description
//
// Increments the -value by amount. amount defaults to 1.0 if omitted.
//
// More information might be available at the [Tcl/Tk progressbar] page.
//
// [Tcl/Tk progressbar]: https://www.tcl.tk/man/tcl9.0/TkCmd/ttk_progressbar.html
func (w *TProgressbarWidget) Step(amount ...float64) {
	s := ""
	if len(amount) != 0 {
		s = fmt.Sprint(amount[0])
	}
	evalErr(fmt.Sprintf("%s step %s", w, s))
}

// ttk::progressbar — Provide progress feedback
//
// # Description
//
// Stop autoincrement mode: cancels any recurring timer event initiated by
// pathName start.
//
// More information might be available at the [Tcl/Tk progressbar] page.
//
// [Tcl/Tk progressbar]: https://www.tcl.tk/man/tcl9.0/TkCmd/ttk_progressbar.html
func (w *TProgressbarWidget) Stop() {
	evalErr(fmt.Sprintf("%s stop", w))
}

func goString(p uintptr) string { // Result can be retained.
	if p == 0 {
		return ""
	}

	p0 := p
	var n int
	for ; *(*byte)(unsafe.Pointer(p)) != 0; n++ {
		p++
	}
	if n != 0 {
		return string(unsafe.Slice((*byte)(unsafe.Pointer(p0)), n))
	}

	return ""
}

// Values option.
//
// Known uses:
//   - [Spinbox] (widget specific)
//   - [TCombobox] (widget specific)
//   - [TSpinbox] (widget specific)
func Values(val any) Opt {
	switch x := val.(type) {
	case []string:
		return rawOption(fmt.Sprintf(`-values {%s}`, tclSafeStrings(x...)))
	case []any:
		return rawOption(fmt.Sprintf(`-values {%s}`, tclSafeList(x...)))
	default:
		return rawOption(fmt.Sprintf(`-values %s`, optionString(val)))
	}
}

// Uniform option.
//
// Known uses:
//   - [GridColumnConfigure] (command specific)
//   - [GridRowConfigure] (command specific)
//
// More information might be available at the [Tcl/Tk grid] page.
//
// [Tcl/Tk grid]: https://www.tcl-lang.org/man/tcl9.0/TkCmd/grid.html#M8
func Uniform(val any) Opt {
	// Credit: Complex_Signal2842, https://www.reddit.com/r/golang/comments/1jcpp3d/comment/mja2l9s/?context=3
	return rawOption(fmt.Sprintf(`-uniform %s`, tclSafeString(fmt.Sprint(val))))
}
