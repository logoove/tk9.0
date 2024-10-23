// Copyright 2024 The tk9.0-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package tk9.0 is a CGo-free, cross platform GUI toolkit for Go. It is
// similar to [Tkinter] for Python.
//
// # Hello world
//
// Also available in _examples/hello.go
//
//	package main
//
//	import . "modernc.org/tk9.0"
//
//	func main() {
//		Pack(Button(Txt("Hello"), Command(func() { Destroy(App) })))
//		App.Wait()
//	}
//
// To execute the above program on any supported target issue something like
//
//	$ CGO_ENABLED=0 go run hello.go
//
// The CGO_ENABLED=0 is optional and here it only demonstrates the program can
// be built without CGo.
//
// # In action
//
//   - [equ] A Plain TeX math editor.
//
// # Debugging
//
// Consider this program in _examples/debugging.go:
//
//	// Build this program using -tags=tk.dmesg
//	package main
//
//	import . "modernc.org/tk9.0"
//
//	func main() {
//		Pack(
//			TButton(Txt("Hello"), Command(func() { Destroy(App) })),
//			Ipadx(10), Ipady(5), Padx(20), Pady(10),
//		)
//		App.Wait()
//	}
//
// Execute the program using the tags as indicated, then close the window or
// click the Hello button. With the tk.dmesg tag the package initialization
// prints the debug messages path. So we can view it, for example, like this:
//
//	$ go run -tags=tk.dmesg _examples/debugging.go | tee log
//	...
//	/tmp/debugging-18876-20240928-163046
//	$ cat /tmp/debugging-18876-20240928-163046
//	[18876 debugging] enter [dmesgon.go:32:0 proc.go:7278:doInit1 proc.go:7245:doInit]
//	...
//	[18876 debugging] code=wm iconphoto . img2 -> r= err=<nil> [tk_purego.go:225:eval tk_purego.go:225:eval tk.go:354:evalErr]
//	[18876 debugging] code=wm title . debugging -> r= err=<nil> [tk_purego.go:225:eval tk_purego.go:225:eval tk.go:354:evalErr]
//	[18876 debugging] code=. configure -padx 4m -pady 3m -> r= err=<nil> [tk_purego.go:225:eval tk_purego.go:225:eval tk.go:354:evalErr]
//	[18876 debugging] code=tk::PlaceWindow . center -> r= err=<nil> [tk_purego.go:225:eval tk_purego.go:225:eval tk.go:354:evalErr]
//	[18876 debugging] code=ttk::button ..tbutton4 -text Hello -command {eventDispatcher 3} -> r=.tbutton4 err=<nil> [tk_purego.go:225:eval tk_purego.go:225:eval tk.go:342:newChild]
//	[18876 debugging] code=pack .tbutton4 -ipadx 10 -ipady 5 -padx 20 -pady 10 -> r= err=<nil> [tk_purego.go:225:eval tk_purego.go:225:eval tk.go:354:evalErr]
//	[18876 debugging] code=destroy . -> r= err=<nil> [tk_purego.go:225:eval tk_purego.go:225:eval tk.go:354:evalErr]
//	[18876 debugging] code=tkwait window . -> r= err=<nil> [tk_purego.go:225:eval tk_purego.go:225:eval tk.go:354:evalErr]
//	$
//
// 18876 was the process PID in this particular run. Using the tags allows to
// inspect the Tcl/Tk code executed during the lifetime of the process.
//
// # Supported targets
//
// These combinations of GOOS and GOARCH are currently supported
//
//	OS      Arch
//	-------------
//	darwin  amd64
//	darwin  arm64
//	freebsd amd64
//	freebsd arm64
//	linux   386
//	linux   amd64
//	linux   arm
//	linux   arm64
//	linux   loong64
//	linux   ppc64le
//	linux   riscv64
//	linux   s390x
//	windows 386
//	windows amd64
//	windows arm64
//
// Specific to [FreeBSD]:
//
// When building with cross-compiling or CGO_ENABLED=0, add the following
// argument to `go` so that these symbols are defined by making fakecgo the
// Cgo.
//
//	-gcflags="github.com/ebitengine/purego/internal/fakecgo=-std"
//
// # Builders
//
// Builder results available at [modern-c.appspot.com].
//
// # Runtime dependencies
//
//   - [Img.Graph] and [CanvasWidget.Graph] require the gnuplot executable
//     available in $PATH.
//
// # Completeness
//
// At the moment the package is a [MVP] allowing to build at least some simple,
// yet useful programs. The full Tk API is not yet usable. Please report
// needed, but non-exposed Tk features at the [issue tracker], thanks.
//
// Providing feedback about the missing building blocks, bugs and your user
// experience is invaluable in helping this package to eventually reach version
// 1. See also [RERO].
//
// # Error handling
//
// The [ErrorMode] variable selects the behaviour on errors for certain
// functions that do not return error.
//
// When ErrorMode is PanicOnError, the default, errors will panic, providing a stack trace.
//
// When ErrorMode is CollectErrors, errors will be recorded using [errors.Join] in
// the [Error] variable.  Even if a function does not return error, it is still
// possible to handle errors in the usual way when needed, except that Error is
// now a static variable. That's a problem in the general case, but less so in
// this package that must be used from a single goroutine only, as
// documented elsewhere.
//
//	// Explicit error handling.
//	ErrorMode = CollectErrors
//	if SomeFunnction(someArgument); Error != nil {
//		... error handling goes here
//	}
//
// This is obviously a compromise enabling to have a way to check for errors
// and, at the same time, the ability to write concise code like:
//
//	// Deferred error handling.
//	if b = Button(Txt("Foo"), Padx(5), Pady(10)); Error != nil {
//		...
//	}
//
// There are altogether four different places where the call to the Button
// function can produce errors as additionally to the call itself, every of its
// three arguments can independently fail as well.  Checking each and one of
// them separately is not always necessary in GUI code. But the explicit option
// in the first example is still available when needed.
//
// # Package initialization
//
// Package initialization is done lazily. This saves noticeable additional
// startup time and avoids screen flicker in hybrid programs that use the GUI
// only on demand.
//
// Early package initialization can be enforced by [Initialize].
//
// Initialization will fail if a Unix process starts on a machine with no
// X server or the process is started in a way that it has no access to the X
// server. On the other hand, this package may work on Unix machines with no X
// server if the process is started remotely using '$ ssh -X foo@bar' and the X
// forwarding is enabled/supported.
//
// Darwin port uses the macOS GUI API and does not use X11.
//
// # The options pattern
//
// Zero or more options can be specified when creating a widget. For example
//
//	b := Button(Txt("Hello"), OverRelief("flat"))
//
// or
//
//	lbl := myFrame.Label(State("disabled"), Width(200))
//
// # Widget path names, image and font names
//
// Tcl/Tk uses widget pathnames, image and font names explicitly set by user
// code.  This package generates those names automatically and they are not
// directly needed in code that uses this package.
//
// # Renamed options
//
// There is, for a example, a Tcl/tk 'text' widget and a '-text' option. This
// package exports the widget as type 'TextWidget', its constructor as function
// 'Text' and the option as function 'Txt'. The complete list is:
//
//   - [Button]  option is renamed to [Btn]
//   - [Label]   option is renamed to [Lbl]
//   - [Menu]    option is renamed to [Mnu]
//   - [Message] option is renamed to [Msg]
//   - [Text]    option is renamed to [Txt]
//
// # OS thread
//
// This package should be used from the same goroutine that initialized the
// package. Package initialization performs a runtime.LockOSThread, meaning
// func main() will start execuing locked on the same OS thread.
//
// # Event handlers
//
// The Command() and similar options expect an argument that must be one of:
//
// - An EventHandler or a function literal of the same signature.
//
// - A func(). This can be used when the handler does not need the associated
// Event instance.
//
// # Specially handled types
//
//   - [time.Duration]
//
// When passing an argument of type [time.Durarion] to a function accepting
// 'any', the duration is converted to an integer number of milliseconds.
//
//   - []byte
//
// When passing an argument of type []byte to a function accepting 'any', the
// byte slice is converted to a [encoding/base64] encoded string.
//
//   - []FileType
//
// When passing an argument of type []FileType to a function accepting 'any',
// the slice is converted to the representation the Tcl/Tk -filetypes option
// expects.
//
// # Tcl/Tk code
//
// At least some minimal knowledge of reading Tcl/Tk code is probably required
// for using this package and/or using the related documentation. However you
// will not need to write any Tcl code and you do not need to care about the
// grammar of Tcl words/string literals and how it differs from Go.
//
// There are several Tcl/Tk tutorials available, for example at
// [tutorialspoint].
//
// # Hacking
//
// Merge requests for known issues are always welcome.
//
// Please send merge requests for new features/APIs after filling and
// discussing the additions/changes at the [issue tracker] first.
//
// # Notes
//
// Most of the documentation is generated directly from the Tcl/Tk
// documentation and may not be entirely correct for the Go package. Those
// parts hopefully still serve as a quick/offline Tcl/Tk reference.
//
// # Additional copyrights
//
// Parts of the documentation are copied and/or modified from the [tcl.tk
// site], see the LICENSE-TCLTK file for details.
//
// Parts of the documentation are copied and/or modified from the [tkinter.ttk site]
// which is
//
//	----------------------------------------------------------------------------
//	© Copyright 2001-2024, Python Software Foundation, licensed under the Python
//	Software Foundation License Version 2.
//	----------------------------------------------------------------------------
//
// # Sponsorship
//
// You can support the maintenance and further development of this package at
// [jnml's LiberaPay] (using PayPal).
//
// # "alt" theme style guide
//
// "Checkbutton.indicator" style element options:
//
//   - [Background]
//   - [Bordercolor]
//   - [Foreground]
//   - [Indicatorcolor]
//   - [Indicatormargin]
//   - [Lightcolor]
//   - [Shadecolor]
//
// "Combobox.downarrow" style element options:
//
//   - [Arrowcolor]
//   - [Arrowsize]
//   - [Background]
//   - [Bordercolor]
//   - [Relief]
//
// "Menubutton.indicator" style element options:
//
//   - [Arrowcolor]
//   - [Arrowsize]
//   - [Direction]
//
// "Radiobutton.indicator" style element options:
//
//   - [Background]
//   - [Bordercolor]
//   - [Foreground]
//   - [Indicatorcolor]
//   - [Indicatormargin]
//   - [Lightcolor]
//   - [Shadecolor]
//
// "Spinbox.downarrow" style element options:
//
//   - [Arrowcolor]
//   - [Arrowsize]
//   - [Background]
//   - [Bordercolor]
//   - [Relief]
//
// "Spinbox.uparrow" style element options:
//
//   - [Arrowcolor]
//   - [Arrowsize]
//   - [Background]
//   - [Bordercolor]
//   - [Relief]
//
// "Treeitem.indicator" style element options:
//
//   - [Foreground]
//   - [Indicatormargins]
//   - [Size]
//
// "arrow" style element options:
//
//   - [Arrowcolor]
//   - [Arrowsize]
//   - [Background]
//   - [Bordercolor]
//   - [Relief]
//
// "border" style element options:
//
//   - [Background]
//   - [Bordercolor]
//   - [Borderwidth]
//   - [Default]
//   - [Relief]
//
// "downarrow" style element options:
//
//   - [Arrowcolor]
//   - [Arrowsize]
//   - [Background]
//   - [Bordercolor]
//   - [Relief]
//
// "field" style element options:
//
//   - [Bordercolor]
//   - [Fieldbackground]
//   - [Focuscolor]
//   - [Focuswidth]
//
// "leftarrow" style element options:
//
//   - [Arrowcolor]
//   - [Arrowsize]
//   - [Background]
//   - [Bordercolor]
//   - [Relief]
//
// "rightarrow" style element options:
//
//   - [Arrowcolor]
//   - [Arrowsize]
//   - [Background]
//   - [Bordercolor]
//   - [Relief]
//
// "slider" style element options:
//
//   - [Background]
//   - [Bordercolor]
//   - [Borderwidth]
//   - [Orient]
//   - [Sliderrelief]
//   - [Sliderthickness]
//
// "thumb" style element options:
//
//   - [Background]
//   - [Bordercolor]
//   - [Orient]
//   - [Relief]
//   - [Width]
//
// "uparrow" style element options:
//
//   - [Arrowcolor]
//   - [Arrowsize]
//   - [Background]
//   - [Bordercolor]
//   - [Relief]
//
// "alt" theme style list
//
//	.
//
// Style map: -foreground {disabled #a3a3a3} -background {disabled #d9d9d9 active #ececec} -embossed {disabled 1}
//
//	ComboboxPopdownFrame
//
// Layout: ComboboxPopdownFrame.border -sticky nswe
//
//	Heading
//
// Layout: Treeheading.cell -sticky nswe Treeheading.border -sticky nswe -children {Treeheading.padding -sticky nswe -children {Treeheading.image -side right -sticky {} Treeheading.text -sticky we}}
//
//	Item
//
// Layout: Treeitem.padding -sticky nswe -children {Treeitem.indicator -side left -sticky {} Treeitem.image -side left -sticky {} Treeitem.text -sticky nswe}
//
//	Separator
//
// Layout: Treeitem.separator -sticky nswe
//
//	TButton
//
// Layout: Button.border -sticky nswe -border 1 -children {Button.focus -sticky nswe -children {Button.padding -sticky nswe -children {Button.label -sticky nswe}}}
//
// Style map: -highlightcolor {alternate black} -relief { {pressed !disabled} sunken {active !disabled} raised }
//
//	TCheckbutton
//
// Layout: Checkbutton.padding -sticky nswe -children {Checkbutton.indicator -side left -sticky {} Checkbutton.focus -side left -sticky w -children {Checkbutton.label -sticky nswe}}
//
// Style map: -indicatorcolor {pressed #d9d9d9 alternate #aaaaaa disabled #d9d9d9}
//
//	TCombobox
//
// Layout: Combobox.field -sticky nswe -children {Combobox.downarrow -side right -sticky ns Combobox.padding -sticky nswe -children {Combobox.textarea -sticky nswe}}
//
// Style map: -fieldbackground {readonly #d9d9d9 disabled #d9d9d9} -arrowcolor {disabled #a3a3a3}
//
//	TEntry
//
// Layout: Entry.field -sticky nswe -border 1 -children {Entry.padding -sticky nswe -children {Entry.textarea -sticky nswe}}
//
// Style map: -fieldbackground {readonly #d9d9d9 disabled #d9d9d9}
//
//	TLabelframe
//
// Layout: Labelframe.border -sticky nswe
//
//	TMenubutton
//
// Layout: Menubutton.border -sticky nswe -children {Menubutton.focus -sticky nswe -children {Menubutton.indicator -side right -sticky {} Menubutton.padding -sticky we -children {Menubutton.label -side left -sticky {}}}}
//
//	TNotebook
//
// Layout: Notebook.client -sticky nswe
//
//	TNotebook.Tab
//
// Layout: Notebook.tab -sticky nswe -children {Notebook.padding -side top -sticky nswe -children {Notebook.focus -side top -sticky nswe -children {Notebook.label -side top -sticky {}}}}
//
// Style map: -expand {selected {1.5p 1.5p 0.75p 0}} -background {selected #d9d9d9}
//
//	TProgressbar
//
// -
//
//	TRadiobutton
//
// Layout: Radiobutton.padding -sticky nswe -children {Radiobutton.indicator -side left -sticky {} Radiobutton.focus -side left -sticky {} -children {Radiobutton.label -sticky nswe}}
//
// Style map: -indicatorcolor {pressed #d9d9d9 alternate #aaaaaa disabled #d9d9d9}
//
//	TScale
//
// -
//
//	TScrollbar
//
// -
//
//	TSpinbox
//
// Layout: Spinbox.field -side top -sticky we -children {null -side right -sticky {} -children {Spinbox.uparrow -side top -sticky e Spinbox.downarrow -side bottom -sticky e} Spinbox.padding -sticky nswe -children {Spinbox.textarea -sticky nswe}}
//
// Style map: -fieldbackground {readonly #d9d9d9 disabled #d9d9d9} -arrowcolor {disabled #a3a3a3}
//
//	Tab
//
// Layout: Notebook.tab -sticky nswe -children {Notebook.padding -side top -sticky nswe -children {Notebook.focus -side top -sticky nswe -children {Notebook.label -side top -sticky {}}}}
//
//	Toolbutton
//
// Layout: Toolbutton.border -sticky nswe -children {Toolbutton.focus -sticky nswe -children {Toolbutton.padding -sticky nswe -children {Toolbutton.label -sticky nswe}}}
//
// Style map: -relief {disabled flat selected sunken pressed sunken active raised} -background {pressed #c3c3c3 active #ececec}
//
//	Treeview
//
// Layout: Treeview.field -sticky nswe -border 1 -children {Treeview.padding -sticky nswe -children {Treeview.treearea -sticky nswe}}
//
// Style map: -foreground {disabled #a3a3a3 selected #ffffff} -background {disabled #d9d9d9 selected #4a6984}
//
//	Treeview.Separator
//
// Layout: Treeitem.separator -sticky nswe
//
// # "aqua" theme style guide
//
// "Button.button" style element options:
//
// "Checkbutton.button" style element options:
//
// "Combobox.button" style element options:
//
// "DisclosureButton.button" style element options:
//
// "Entry.field" style element options:
//
//   - [Background]
//   - [Fieldbackground]
//
// "GradientButton.button" style element options:
//
// "HelpButton.button" style element options:
//
// "Horizontal.Scrollbar.leftarrow" style element options:
//
//   - [Orient]
//
// "Horizontal.Scrollbar.rightarrow" style element options:
//
//   - [Orient]
//
// "Horizontal.Scrollbar.thumb" style element options:
//
//   - [Orient]
//
// "Horizontal.Scrollbar.trough" style element options:
//
//   - [Orient]
//
// "InlineButton.button" style element options:
//
// "Labelframe.border" style element options:
//
// "Menubutton.button" style element options:
//
// "Notebook.client" style element options:
//
// "Notebook.tab" style element options:
//
// "Progressbar.track" style element options:
//
//   - [Maximum]
//   - [Mode]
//   - [Orient]
//   - [Phase]
//   - [Value]
//
// "Radiobutton.button" style element options:
//
// "RecessedButton.button" style element options:
//
// "RoundedRectButton.button" style element options:
//
// "Scale.slider" style element options:
//
// "Scale.trough" style element options:
//
//   - [From]
//   - [Orient]
//   - [To]
//   - [Value]
//
// "Searchbox.field" style element options:
//
//   - [Background]
//   - [Fieldbackground]
//
// "SidebarButton.button" style element options:
//
// "Spinbox.downarrow" style element options:
//
// "Spinbox.field" style element options:
//
//   - [Background]
//   - [Fieldbackground]
//
// "Spinbox.uparrow" style element options:
//
// "Toolbar.background" style element options:
//
// "Toolbutton.border" style element options:
//
// "Treeheading.cell" style element options:
//
// "Treeitem.indicator" style element options:
//
// "Treeview.treearea" style element options:
//
// "Vertical.Scrollbar.downarrow" style element options:
//
//   - [Orient]
//
// "Vertical.Scrollbar.thumb" style element options:
//
//   - [Orient]
//
// "Vertical.Scrollbar.trough" style element options:
//
//   - [Orient]
//
// "Vertical.Scrollbar.uparrow" style element options:
//
//   - [Orient]
//
// "background" style element options:
//
// "field" style element options:
//
//   - [Fieldbackground]
//
// "fill" style element options:
//
// "hseparator" style element options:
//
// "separator" style element options:
//
// "sizegrip" style element options:
//
// "vseparator" style element options:
//
// "aqua" theme style list
//
//	.
//
// Style map: -selectforeground { background systemSelectedTextColor !focus systemSelectedTextColor} -foreground { disabled systemDisabledControlTextColor background systemLabelColor} -selectbackground { background systemSelectedTextBackgroundColor !focus systemSelectedTextBackgroundColor}
//
//	DisclosureButton
//
// Layout: DisclosureButton.button -sticky nswe
//
//	GradientButton
//
// Layout: GradientButton.button -sticky nswe -children {Button.padding -sticky nswe -children {Button.label -sticky nswe}}
//
//	Heading
//
// Layout: Treeheading.cell -sticky nswe Treeheading.image -side right -sticky {} Treeheading.text -side top -sticky {}
//
//	HelpButton
//
// Layout: HelpButton.button -sticky nswe
//
//	Horizontal.TScrollbar
//
// Layout: Horizontal.Scrollbar.trough -sticky we -children {Horizontal.Scrollbar.thumb -sticky nswe Horizontal.Scrollbar.rightarrow -side right -sticky {} Horizontal.Scrollbar.leftarrow -side right -sticky {}}
//
//	ImageButton
//
// Layout: Button.padding -sticky nswe -children {Button.label -sticky nswe}
//
// Style map: -foreground { pressed systemLabelColor !pressed systemSecondaryLabelColor }
//
//	InlineButton
//
// Layout: InlineButton.button -sticky nswe -children {Button.padding -sticky nswe -children {Button.label -sticky nswe}}
//
// Style map: -foreground { disabled systemWindowBackgroundColor }
//
//	Item
//
// Layout: Treeitem.padding -sticky nswe -children {Treeitem.indicator -side left -sticky {} Treeitem.image -side left -sticky {} Treeitem.text -side left -sticky {}}
//
//	Label
//
// Layout: Label.fill -sticky nswe -children {Label.text -sticky nswe}
//
//	RecessedButton
//
// Layout: RecessedButton.button -sticky nswe -children {Button.padding -sticky nswe -children {Button.label -sticky nswe}}
//
// Style map: -font { selected RecessedFont active RecessedFont pressed RecessedFont } -foreground { {disabled selected} systemWindowBackgroundColor3 {disabled !selected} systemDisabledControlTextColor selected systemTextBackgroundColor active white pressed white }
//
//	RoundedRectButton
//
// Layout: RoundedRectButton.button -sticky nswe -children {Button.padding -sticky nswe -children {Button.label -sticky nswe}}
//
//	Searchbox
//
// Layout: Searchbox.field -sticky nswe -border 1 -children {Entry.padding -sticky nswe -children {Entry.textarea -sticky nswe}}
//
//	SidebarButton
//
// Layout: SidebarButton.button -sticky nswe -children {Button.padding -sticky nswe -children {Button.label -sticky nswe}}
//
// Style map: -foreground { {disabled selected} systemWindowBackgroundColor3 {disabled !selected} systemDisabledControlTextColor selected systemTextColor active systemTextColor pressed systemTextColor }
//
//	TButton
//
// Layout: Button.button -sticky nswe -children {Button.padding -sticky nswe -children {Button.label -sticky nswe}}
//
// Style map: -foreground { pressed white {alternate !pressed !background} white disabled systemDisabledControlTextColor}
//
//	TCheckbutton
//
// Layout: Checkbutton.button -sticky nswe -children {Checkbutton.padding -sticky nswe -children {Checkbutton.label -side left -sticky {}}}
//
//	TCombobox
//
// Layout: Combobox.button -sticky nswe -children {Combobox.padding -sticky nswe -children {Combobox.textarea -sticky nswe}}
//
// Style map: -foreground { disabled systemDisabledControlTextColor } -selectbackground { !focus systemUnemphasizedSelectedTextBackgroundColor }
//
//	TEntry
//
// Layout: Entry.field -sticky nswe -border 1 -children {Entry.padding -sticky nswe -children {Entry.textarea -sticky nswe}}
//
// Style map: -foreground { disabled systemDisabledControlTextColor } -selectbackground { !focus systemUnemphasizedSelectedTextBackgroundColor }
//
//	TLabelframe
//
// Layout: Labelframe.border -sticky nswe
//
//	TLabelframe.Label
//
// Layout: Label.fill -sticky nswe -children {Label.text -sticky nswe}
//
//	TMenubutton
//
// Layout: Menubutton.button -sticky nswe -children {Menubutton.padding -sticky nswe -children {Menubutton.label -side left -sticky {}}}
//
//	TNotebook
//
// Layout: Notebook.client -sticky nswe
//
//	TNotebook.Tab
//
// Layout: Notebook.tab -sticky nswe -children {Notebook.padding -sticky nswe -children {Notebook.label -sticky nswe}}
//
// Style map: -foreground { {background !selected} systemControlTextColor {background selected} black {!background selected} systemSelectedTabTextColor disabled systemDisabledControlTextColor}
//
//	TProgressbar
//
// Layout: Progressbar.track -sticky nswe
//
//	TRadiobutton
//
// Layout: Radiobutton.button -sticky nswe -children {Radiobutton.padding -sticky nswe -children {Radiobutton.label -side left -sticky {}}}
//
//	TScrollbar
//
// -
//
//	TSpinbox
//
// Layout: Spinbox.buttons -side right -sticky {} -children {Spinbox.uparrow -side top -sticky e Spinbox.downarrow -side bottom -sticky e} Spinbox.field -sticky we -children {Spinbox.textarea -sticky we}
//
// Style map: -foreground { disabled systemDisabledControlTextColor } -selectbackground { !focus systemUnemphasizedSelectedTextBackgroundColor }
//
//	Tab
//
// Layout: Notebook.tab -sticky nswe -children {Notebook.padding -sticky nswe -children {Notebook.label -sticky nswe}}
//
//	Toolbar
//
// Layout: Toolbar.background -sticky nswe
//
//	Toolbutton
//
// Layout: Toolbutton.border -sticky nswe -children {Toolbutton.focus -sticky nswe -children {Toolbutton.padding -sticky nswe -children {Toolbutton.label -sticky nswe}}}
//
//	Treeview
//
// Layout: Treeview.field -sticky nswe -children {Treeview.padding -sticky nswe -children {Treeview.treearea -sticky nswe}}
//
// Style map: -background { selected systemSelectedTextBackgroundColor }
//
//	Vertical.TScrollbar
//
// Layout: Vertical.Scrollbar.trough -sticky ns -children {Vertical.Scrollbar.thumb -sticky nswe Vertical.Scrollbar.downarrow -side bottom -sticky {} Vertical.Scrollbar.uparrow -side bottom -sticky {}}
//
// # "clam" theme style guide
//
// "Checkbutton.indicator" style element options:
//
//   - [Indicatorbackground]
//   - [Indicatorforeground]
//   - [Indicatormargin]
//   - [Lowerbordercolor]
//   - [Upperbordercolor]
//
// "Combobox.field" style element options:
//
//   - [Bordercolor]
//   - [Fieldbackground]
//   - [Lightcolor]
//
// "Radiobutton.indicator" style element options:
//
//   - [Indicatorbackground]
//   - [Indicatorforeground]
//   - [Indicatormargin]
//   - [Lowerbordercolor]
//   - [Upperbordercolor]
//
// "Spinbox.downarrow" style element options:
//
//   - [Arrowcolor]
//   - [Arrowsize]
//   - [Background]
//   - [Bordercolor]
//   - [Darkcolor]
//   - [Gripsize]
//   - [Lightcolor]
//   - [Orient]
//   - [Sliderlength]
//   - [Troughcolor]
//
// "Spinbox.uparrow" style element options:
//
//   - [Arrowcolor]
//   - [Arrowsize]
//   - [Background]
//   - [Bordercolor]
//   - [Darkcolor]
//   - [Gripsize]
//   - [Lightcolor]
//   - [Orient]
//   - [Sliderlength]
//   - [Troughcolor]
//
// "arrow" style element options:
//
//   - [Arrowcolor]
//   - [Arrowsize]
//   - [Background]
//   - [Bordercolor]
//   - [Darkcolor]
//   - [Gripsize]
//   - [Lightcolor]
//   - [Orient]
//   - [Sliderlength]
//   - [Troughcolor]
//
// "bar" style element options:
//
//   - [Arrowcolor]
//   - [Arrowsize]
//   - [Background]
//   - [Bordercolor]
//   - [Darkcolor]
//   - [Gripsize]
//   - [Lightcolor]
//   - [Orient]
//   - [Sliderlength]
//   - [Troughcolor]
//
// "border" style element options:
//
//   - [Bordercolor]
//   - [Borderwidth]
//   - [Darkcolor]
//   - [Lightcolor]
//   - [Relief]
//
// "client" style element options:
//
//   - [Background]
//   - [Bordercolor]
//   - [Darkcolor]
//   - [Lightcolor]
//
// "downarrow" style element options:
//
//   - [Arrowcolor]
//   - [Arrowsize]
//   - [Background]
//   - [Bordercolor]
//   - [Darkcolor]
//   - [Gripsize]
//   - [Lightcolor]
//   - [Orient]
//   - [Sliderlength]
//   - [Troughcolor]
//
// "field" style element options:
//
//   - [Bordercolor]
//   - [Fieldbackground]
//   - [Lightcolor]
//
// "hgrip" style element options:
//
//   - [Bordercolor]
//   - [Gripsize]
//   - [Lightcolor]
//
// "leftarrow" style element options:
//
//   - [Arrowcolor]
//   - [Arrowsize]
//   - [Background]
//   - [Bordercolor]
//   - [Darkcolor]
//   - [Gripsize]
//   - [Lightcolor]
//   - [Orient]
//   - [Sliderlength]
//   - [Troughcolor]
//
// "pbar" style element options:
//
//   - [Arrowcolor]
//   - [Arrowsize]
//   - [Background]
//   - [Bordercolor]
//   - [Darkcolor]
//   - [Gripsize]
//   - [Lightcolor]
//   - [Orient]
//   - [Sliderlength]
//   - [Troughcolor]
//
// "rightarrow" style element options:
//
//   - [Arrowcolor]
//   - [Arrowsize]
//   - [Background]
//   - [Bordercolor]
//   - [Darkcolor]
//   - [Gripsize]
//   - [Lightcolor]
//   - [Orient]
//   - [Sliderlength]
//   - [Troughcolor]
//
// "slider" style element options:
//
//   - [Arrowcolor]
//   - [Arrowsize]
//   - [Background]
//   - [Bordercolor]
//   - [Darkcolor]
//   - [Gripsize]
//   - [Lightcolor]
//   - [Orient]
//   - [Sliderlength]
//   - [Troughcolor]
//
// "tab" style element options:
//
//   - [Background]
//   - [Bordercolor]
//   - [Darkcolor]
//   - [Lightcolor]
//
// "thumb" style element options:
//
//   - [Arrowcolor]
//   - [Arrowsize]
//   - [Background]
//   - [Bordercolor]
//   - [Darkcolor]
//   - [Gripsize]
//   - [Lightcolor]
//   - [Orient]
//   - [Sliderlength]
//   - [Troughcolor]
//
// "trough" style element options:
//
//   - [Arrowcolor]
//   - [Arrowsize]
//   - [Background]
//   - [Bordercolor]
//   - [Darkcolor]
//   - [Gripsize]
//   - [Lightcolor]
//   - [Orient]
//   - [Sliderlength]
//   - [Troughcolor]
//
// "uparrow" style element options:
//
//   - [Arrowcolor]
//   - [Arrowsize]
//   - [Background]
//   - [Bordercolor]
//   - [Darkcolor]
//   - [Gripsize]
//   - [Lightcolor]
//   - [Orient]
//   - [Sliderlength]
//   - [Troughcolor]
//
// "vgrip" style element options:
//
//   - [Bordercolor]
//   - [Gripsize]
//   - [Lightcolor]
//
// "clam" theme style list
//
//	.
//
// Style map: -selectforeground {!focus white} -foreground {disabled #999999} -selectbackground {!focus #9e9a91} -background {disabled #dcdad5 active #eeebe7}
//
//	ComboboxPopdownFrame
//
// Layout: ComboboxPopdownFrame.border -sticky nswe
//
//	Heading
//
// Layout: Treeheading.cell -sticky nswe Treeheading.border -sticky nswe -children {Treeheading.padding -sticky nswe -children {Treeheading.image -side right -sticky {} Treeheading.text -sticky we}}
//
//	Horizontal.Sash
//
// Layout: Sash.hsash -sticky nswe -children {Sash.hgrip -sticky nswe}
//
//	Item
//
// Layout: Treeitem.padding -sticky nswe -children {Treeitem.indicator -side left -sticky {} Treeitem.image -side left -sticky {} Treeitem.text -sticky nswe}
//
//	Sash
//
// -
//
//	Separator
//
// Layout: Treeitem.separator -sticky nswe
//
//	TButton
//
// Layout: Button.border -sticky nswe -border 1 -children {Button.focus -sticky nswe -children {Button.padding -sticky nswe -children {Button.label -sticky nswe}}}
//
// Style map: -lightcolor {pressed #bab5ab} -background {disabled #dcdad5 pressed #bab5ab active #eeebe7} -bordercolor {alternate #000000} -darkcolor {pressed #bab5ab}
//
//	TCheckbutton
//
// Layout: Checkbutton.padding -sticky nswe -children {Checkbutton.indicator -side left -sticky {} Checkbutton.focus -side left -sticky w -children {Checkbutton.label -sticky nswe}}
//
// Style map: -indicatorbackground {pressed #dcdad5 {!disabled alternate} #5895bc {disabled alternate} #a0a0a0 disabled #dcdad5}
//
//	TCombobox
//
// Layout: Combobox.downarrow -side right -sticky ns Combobox.field -sticky nswe -children {Combobox.padding -sticky nswe -children {Combobox.textarea -sticky nswe}}
//
// Style map: -foreground {{readonly focus} #ffffff} -fieldbackground {{readonly focus} #4a6984 readonly #dcdad5} -background {active #eeebe7 pressed #eeebe7} -bordercolor {focus #4a6984} -arrowcolor {disabled #999999}
//
//	TEntry
//
// Layout: Entry.field -sticky nswe -border 1 -children {Entry.padding -sticky nswe -children {Entry.textarea -sticky nswe}}
//
// Style map: -lightcolor {focus #6f9dc6} -background {readonly #dcdad5} -bordercolor {focus #4a6984}
//
//	TLabelframe
//
// Layout: Labelframe.border -sticky nswe
//
//	TMenubutton
//
// Layout: Menubutton.border -sticky nswe -children {Menubutton.focus -sticky nswe -children {Menubutton.indicator -side right -sticky {} Menubutton.padding -sticky we -children {Menubutton.label -side left -sticky {}}}}
//
//	TNotebook.Tab
//
// Layout: Notebook.tab -sticky nswe -children {Notebook.padding -side top -sticky nswe -children {Notebook.focus -side top -sticky nswe -children {Notebook.label -side top -sticky {}}}}
//
// Style map: -lightcolor {selected #eeebe7 {} #cfcdc8} -padding {selected {4.5p 3p 4.5p 1.5p}} -background {selected #dcdad5 {} #bab5ab}
//
//	TProgressbar
//
// -
//
//	TRadiobutton
//
// Layout: Radiobutton.padding -sticky nswe -children {Radiobutton.indicator -side left -sticky {} Radiobutton.focus -side left -sticky {} -children {Radiobutton.label -sticky nswe}}
//
// Style map: -indicatorbackground {pressed #dcdad5 {!disabled alternate} #5895bc {disabled alternate} #a0a0a0 disabled #dcdad5}
//
//	TScale
//
// -
//
//	TScrollbar
//
// -
//
//	TSpinbox
//
// Layout: Spinbox.field -side top -sticky we -children {null -side right -sticky {} -children {Spinbox.uparrow -side top -sticky e Spinbox.downarrow -side bottom -sticky e} Spinbox.padding -sticky nswe -children {Spinbox.textarea -sticky nswe}}
//
// Style map: -background {readonly #dcdad5} -bordercolor {focus #4a6984} -arrowcolor {disabled #999999}
//
//	Tab
//
// Layout: Notebook.tab -sticky nswe -children {Notebook.padding -side top -sticky nswe -children {Notebook.focus -side top -sticky nswe -children {Notebook.label -side top -sticky {}}}}
//
//	Toolbutton
//
// Layout: Toolbutton.border -sticky nswe -children {Toolbutton.focus -sticky nswe -children {Toolbutton.padding -sticky nswe -children {Toolbutton.label -sticky nswe}}}
//
// Style map: -lightcolor {pressed #bab5ab} -relief {disabled flat selected sunken pressed sunken active raised} -background {disabled #dcdad5 pressed #bab5ab active #eeebe7} -darkcolor {pressed #bab5ab}
//
//	Treeview
//
// Layout: Treeview.field -sticky nswe -border 1 -children {Treeview.padding -sticky nswe -children {Treeview.treearea -sticky nswe}}
//
// Style map: -foreground {disabled #999999 selected #ffffff} -background {disabled #dcdad5 selected #4a6984} -bordercolor {focus #4a6984}
//
//	Treeview.Separator
//
// Layout: Treeitem.separator -sticky nswe
//
//	Vertical.Sash
//
// Layout: Sash.vsash -sticky nswe -children {Sash.vgrip -sticky nswe}
//
// # "classic" theme style guide
//
// "Button.border" style element options:
//
//   - [Background]
//   - [Borderwidth]
//   - [Default]
//   - [Relief]
//
// "Checkbutton.indicator" style element options:
//
//   - [Background]
//   - [Borderwidth]
//   - [Indicatorcolor]
//   - [Indicatormargin]
//   - [Indicatorrelief]
//   - [Indicatorsize]
//
// "Combobox.downarrow" style element options:
//
//   - [Arrowcolor]
//   - [Arrowsize]
//   - [Background]
//   - [Borderwidth]
//   - [Relief]
//
// "Menubutton.indicator" style element options:
//
//   - [Background]
//   - [Indicatorborderwidth]
//   - [Indicatorheight]
//   - [Indicatormargin]
//   - [Indicatorrelief]
//   - [Indicatorwidth]
//
// "Radiobutton.indicator" style element options:
//
//   - [Background]
//   - [Borderwidth]
//   - [Indicatorcolor]
//   - [Indicatormargin]
//   - [Indicatorrelief]
//   - [Indicatorsize]
//
// "Spinbox.downarrow" style element options:
//
//   - [Arrowcolor]
//   - [Arrowsize]
//   - [Background]
//   - [Borderwidth]
//   - [Relief]
//
// "Spinbox.uparrow" style element options:
//
//   - [Arrowcolor]
//   - [Arrowsize]
//   - [Background]
//   - [Borderwidth]
//   - [Relief]
//
// "arrow" style element options:
//
//   - [Arrowsize]
//   - [Background]
//   - [Borderwidth]
//   - [Relief]
//
// "downarrow" style element options:
//
//   - [Arrowsize]
//   - [Background]
//   - [Borderwidth]
//   - [Relief]
//
// "highlight" style element options:
//
//   - [Default]
//   - [Highlightcolor]
//   - [Highlightthickness]
//
// "hsash" style element options:
//
//   - [Background]
//   - [Handlepad]
//   - [Handlesize]
//   - [Sashpad]
//   - [Sashrelief]
//   - [Sashthickness]
//
// "leftarrow" style element options:
//
//   - [Arrowsize]
//   - [Background]
//   - [Borderwidth]
//   - [Relief]
//
// "rightarrow" style element options:
//
//   - [Arrowsize]
//   - [Background]
//   - [Borderwidth]
//   - [Relief]
//
// "slider" style element options:
//
//   - [Background]
//   - [Orient]
//   - [Sliderborderwidth]
//   - [Sliderlength]
//   - [Sliderrelief]
//   - [Sliderthickness]
//
// "uparrow" style element options:
//
//   - [Arrowsize]
//   - [Background]
//   - [Borderwidth]
//   - [Relief]
//
// "vsash" style element options:
//
//   - [Background]
//   - [Handlepad]
//   - [Handlesize]
//   - [Sashpad]
//   - [Sashrelief]
//   - [Sashthickness]
//
// "classic" theme style list
//
//	.
//
// Style map: -highlightcolor {focus black} -foreground {disabled #a3a3a3} -background {disabled #d9d9d9 active #ececec}
//
//	ComboboxPopdownFrame
//
// Layout: ComboboxPopdownFrame.border -sticky nswe
//
//	Heading
//
// Layout: Treeheading.cell -sticky nswe Treeheading.border -sticky nswe -children {Treeheading.padding -sticky nswe -children {Treeheading.image -side right -sticky {} Treeheading.text -sticky we}}
//
//	Horizontal.TScale
//
// Layout: Horizontal.Scale.highlight -sticky nswe -children {Horizontal.Scale.trough -sticky nswe -children {Horizontal.Scale.slider -side left -sticky {}}}
//
//	Item
//
// Layout: Treeitem.padding -sticky nswe -children {Treeitem.indicator -side left -sticky {} Treeitem.image -side left -sticky {} Treeitem.text -sticky nswe}
//
//	Sash
//
// -
//
//	Separator
//
// Layout: Treeitem.separator -sticky nswe
//
//	TButton
//
// Layout: Button.highlight -sticky nswe -children {Button.border -sticky nswe -border 1 -children {Button.padding -sticky nswe -children {Button.label -sticky nswe}}}
//
// Style map: -relief {{!disabled pressed} sunken}
//
//	TCheckbutton
//
// Layout: Checkbutton.highlight -sticky nswe -children {Checkbutton.border -sticky nswe -children {Checkbutton.padding -sticky nswe -children {Checkbutton.indicator -side left -sticky {} Checkbutton.label -side left -sticky nswe}}}
//
// Style map: -indicatorrelief {alternate raised selected sunken pressed sunken} -indicatorcolor {pressed #d9d9d9 alternate #b05e5e selected #b03060}
//
//	TCombobox
//
// Layout: Combobox.highlight -sticky nswe -children {Combobox.field -sticky nswe -children {Combobox.downarrow -side right -sticky ns Combobox.padding -sticky nswe -children {Combobox.textarea -sticky nswe}}}
//
// Style map: -fieldbackground {readonly #d9d9d9 disabled #d9d9d9}
//
//	TEntry
//
// Layout: Entry.highlight -sticky nswe -children {Entry.field -sticky nswe -border 1 -children {Entry.padding -sticky nswe -children {Entry.textarea -sticky nswe}}}
//
// Style map: -fieldbackground {readonly #d9d9d9 disabled #d9d9d9}
//
//	TLabelframe
//
// Layout: Labelframe.border -sticky nswe
//
//	TMenubutton
//
// Layout: Menubutton.highlight -sticky nswe -children {Menubutton.border -sticky nswe -children {Menubutton.indicator -side right -sticky {} Menubutton.padding -sticky we -children {Menubutton.label -sticky {}}}}
//
//	TNotebook.Tab
//
// Layout: Notebook.tab -sticky nswe -children {Notebook.padding -side top -sticky nswe -children {Notebook.focus -side top -sticky nswe -children {Notebook.label -side top -sticky {}}}}
//
// Style map: -background {selected #d9d9d9}
//
//	TProgressbar
//
// -
//
//	TRadiobutton
//
// Layout: Radiobutton.highlight -sticky nswe -children {Radiobutton.border -sticky nswe -children {Radiobutton.padding -sticky nswe -children {Radiobutton.indicator -side left -sticky {} Radiobutton.label -side left -sticky nswe}}}
//
// Style map: -indicatorrelief {alternate raised selected sunken pressed sunken} -indicatorcolor {pressed #d9d9d9 alternate #b05e5e selected #b03060}
//
//	TScale
//
// Style map: -sliderrelief {{pressed !disabled} sunken}
//
//	TScrollbar
//
// Style map: -relief {{pressed !disabled} sunken}
//
//	TSpinbox
//
// Layout: Spinbox.highlight -sticky nswe -children {Spinbox.field -sticky nswe -children {null -side right -sticky {} -children {Spinbox.uparrow -side top -sticky e Spinbox.downarrow -side bottom -sticky e} Spinbox.padding -sticky nswe -children {Spinbox.textarea -sticky nswe}}}
//
// Style map: -fieldbackground {readonly #d9d9d9 disabled #d9d9d9}
//
//	Tab
//
// Layout: Notebook.tab -sticky nswe -children {Notebook.padding -side top -sticky nswe -children {Notebook.focus -side top -sticky nswe -children {Notebook.label -side top -sticky {}}}}
//
//	Toolbutton
//
// Layout: Toolbutton.focus -sticky nswe -children {Toolbutton.border -sticky nswe -children {Toolbutton.padding -sticky nswe -children {Toolbutton.label -sticky nswe}}}
//
// Style map: -relief {disabled flat selected sunken pressed sunken active raised} -background {pressed #b3b3b3 active #ececec}
//
//	Treeview
//
// Layout: Treeview.highlight -sticky nswe -children {Treeview.field -sticky nswe -border 1 -children {Treeview.padding -sticky nswe -children {Treeview.treearea -sticky nswe}}}
//
// Style map: -foreground {disabled #a3a3a3 selected #000000} -background {disabled #d9d9d9 selected #c3c3c3}
//
//	Treeview.Separator
//
// Layout: Treeitem.separator -sticky nswe
//
//	Vertical.TScale
//
// Layout: Vertical.Scale.highlight -sticky nswe -children {Vertical.Scale.trough -sticky nswe -children {Vertical.Scale.slider -side top -sticky {}}}
//
// # "default" theme style guide
//
// "" style element options:
//
// "Checkbutton.indicator" style element options:
//
//   - [Bordercolor]
//   - [Indicatorbackground]
//   - [Indicatorforeground]
//   - [Indicatormargin]
//
// "Combobox.downarrow" style element options:
//
//   - [Arrowcolor]
//   - [Arrowsize]
//   - [Background]
//   - [Borderwidth]
//   - [Relief]
//
// "Menubutton.indicator" style element options:
//
//   - [Arrowcolor]
//   - [Arrowpadding]
//   - [Arrowsize]
//
// "Radiobutton.indicator" style element options:
//
//   - [Bordercolor]
//   - [Indicatorbackground]
//   - [Indicatorforeground]
//   - [Indicatormargin]
//
// "Spinbox.downarrow" style element options:
//
//   - [Arrowcolor]
//   - [Arrowsize]
//   - [Background]
//   - [Borderwidth]
//   - [Relief]
//
// "Spinbox.uparrow" style element options:
//
//   - [Arrowcolor]
//   - [Arrowsize]
//   - [Background]
//   - [Borderwidth]
//   - [Relief]
//
// "Treeheading.cell" style element options:
//
//   - [Background]
//   - [Rownumber]
//
// "Treeitem.indicator" style element options:
//
//   - [Foreground]
//   - [Indicatormargins]
//   - [Indicatorsize]
//
// "Treeitem.row" style element options:
//
//   - [Background]
//   - [Rownumber]
//
// "Treeitem.separator" style element options:
//
//   - [Background]
//   - [Rownumber]
//
// "arrow" style element options:
//
//   - [Arrowcolor]
//   - [Arrowsize]
//   - [Background]
//   - [Borderwidth]
//   - [Relief]
//
// "background" style element options:
//
//   - [Background]
//
// "border" style element options:
//
//   - [Background]
//   - [Borderwidth]
//   - [Relief]
//
// "client" style element options:
//
//   - [Background]
//   - [Borderwidth]
//   - [Highlight]
//   - [Highlightcolor]
//
// "ctext" style element options:
//
//   - [Anchor]
//   - [Embossed]
//   - [Font]
//   - [Foreground]
//   - [Justify]
//   - [Text]
//   - [Underline]
//   - [Width]
//   - [Wraplength]
//
// "downarrow" style element options:
//
//   - [Arrowcolor]
//   - [Arrowsize]
//   - [Background]
//   - [Borderwidth]
//   - [Relief]
//
// "field" style element options:
//
//   - [Borderwidth]
//   - [Fieldbackground]
//   - [Focuscolor]
//   - [Focuswidth]
//
// "fill" style element options:
//
//   - [Background]
//
// "focus" style element options:
//
//   - [Focuscolor]
//   - [Focussolid]
//   - [Focusthickness]
//
// "hsash" style element options:
//
//   - [Sashthickness]
//
// "hseparator" style element options:
//
//   - [Background]
//   - [Orient]
//
// "image" style element options:
//
//   - [Background]
//   - [Image]
//   - [Stipple]
//
// "indicator" style element options:
//
// "label" style element options:
//
//   - [Anchor]
//   - [Background]
//   - [Compound]
//   - [Embossed]
//   - [Font]
//   - [Foreground]
//   - [Image]
//   - [Justify]
//   - [Space]
//   - [Stipple]
//   - [Text]
//   - [Underline]
//   - [Width]
//   - [Wraplength]
//
// "leftarrow" style element options:
//
//   - [Arrowcolor]
//   - [Arrowsize]
//   - [Background]
//   - [Borderwidth]
//   - [Relief]
//
// "padding" style element options:
//
//   - [Padding]
//   - [Relief]
//   - [Shiftrelief]
//
// "pbar" style element options:
//
//   - [Background]
//   - [Barsize]
//   - [Borderwidth]
//   - [Orient]
//   - [Pbarrelief]
//   - [Thickness]
//
// "rightarrow" style element options:
//
//   - [Arrowcolor]
//   - [Arrowsize]
//   - [Background]
//   - [Borderwidth]
//   - [Relief]
//
// "separator" style element options:
//
//   - [Background]
//   - [Orient]
//
// "sizegrip" style element options:
//
//   - [Background]
//   - [Gripsize]
//
// "slider" style element options:
//
//   - [Bordercolor]
//   - [Innercolor]
//   - [Orient]
//   - [Outercolor]
//
// "tab" style element options:
//
//   - [Background]
//   - [Borderwidth]
//   - [Highlight]
//   - [Highlightcolor]
//
// "text" style element options:
//
//   - [Anchor]
//   - [Embossed]
//   - [Font]
//   - [Foreground]
//   - [Justify]
//   - [Text]
//   - [Underline]
//   - [Width]
//   - [Wraplength]
//
// "textarea" style element options:
//
//   - [Font]
//   - [Width]
//
// "thumb" style element options:
//
//   - [Background]
//   - [Borderwidth]
//   - [Orient]
//   - [Relief]
//   - [Width]
//
// "treearea" style element options:
//
// "trough" style element options:
//
//   - [Groovewidth]
//   - [Orient]
//   - [Troughborderwidth]
//   - [Troughcolor]
//   - [Troughrelief]
//
// "uparrow" style element options:
//
//   - [Arrowcolor]
//   - [Arrowsize]
//   - [Background]
//   - [Borderwidth]
//   - [Relief]
//
// "vsash" style element options:
//
//   - [Sashthickness]
//
// "vseparator" style element options:
//
//   - [Background]
//   - [Orient]
//
// "default" theme style list
//
//	.
//
// Style map: -foreground {disabled #a3a3a3} -background {disabled #edeceb active #ececec}
//
//	Cell
//
// Layout: Treedata.padding -sticky nswe -children {Treeitem.image -side left -sticky {} Treeitem.text -sticky nswe}
//
//	ComboboxPopdownFrame
//
// Layout: ComboboxPopdownFrame.border -sticky nswe
//
//	Heading
//
// Layout: Treeheading.cell -sticky nswe Treeheading.border -sticky nswe -children {Treeheading.padding -sticky nswe -children {Treeheading.image -side right -sticky {} Treeheading.text -sticky we}}
//
//	Horizontal.Sash
//
// Layout: Sash.hsash -sticky we
//
//	Horizontal.TProgressbar
//
// Layout: Horizontal.Progressbar.trough -sticky nswe -children {Horizontal.Progressbar.pbar -side left -sticky ns Horizontal.Progressbar.ctext -side left -sticky {}}
//
//	Horizontal.TScale
//
// Layout: Horizontal.Scale.focus -sticky nswe -children {Horizontal.Scale.padding -sticky nswe -children {Horizontal.Scale.trough -sticky nswe -children {Horizontal.Scale.slider -side left -sticky {}}}}
//
//	Horizontal.TScrollbar
//
// Layout: Horizontal.Scrollbar.trough -sticky we -children {Horizontal.Scrollbar.leftarrow -side left -sticky {} Horizontal.Scrollbar.rightarrow -side right -sticky {} Horizontal.Scrollbar.thumb -sticky nswe}
//
//	Item
//
// Layout: Treeitem.padding -sticky nswe -children {Treeitem.indicator -side left -sticky {} Treeitem.image -side left -sticky {} Treeitem.text -sticky nswe}
//
//	Label
//
// Layout: Label.fill -sticky nswe -children {Label.text -sticky nswe}
//
//	Row
//
// Layout: Treeitem.row -sticky nswe
//
//	Sash
//
// -
//
//	Separator
//
// Layout: Treeitem.separator -sticky nswe
//
//	TButton
//
// Layout: Button.border -sticky nswe -border 1 -children {Button.focus -sticky nswe -children {Button.padding -sticky nswe -children {Button.label -sticky nswe}}}
//
// Style map: -relief {{!disabled pressed} sunken}
//
//	TCheckbutton
//
// Layout: Checkbutton.padding -sticky nswe -children {Checkbutton.indicator -side left -sticky {} Checkbutton.focus -side left -sticky w -children {Checkbutton.label -sticky nswe}}
//
// Style map: -indicatorbackground {{alternate disabled} #a3a3a3 {alternate pressed} #5895bc alternate #4a6984 {selected disabled} #a3a3a3 {selected pressed} #5895bc selected #4a6984 disabled #edeceb pressed #c3c3c3}
//
//	TCombobox
//
// Layout: Combobox.field -sticky nswe -children {Combobox.downarrow -side right -sticky ns Combobox.padding -sticky nswe -children {Combobox.textarea -sticky nswe}}
//
// Style map: -fieldbackground {readonly #edeceb disabled #edeceb} -arrowcolor {disabled #a3a3a3}
//
//	TEntry
//
// Layout: Entry.field -sticky nswe -border 1 -children {Entry.padding -sticky nswe -children {Entry.textarea -sticky nswe}}
//
// Style map: -fieldbackground {readonly #edeceb disabled #edeceb}
//
//	TFrame
//
// Layout: Frame.border -sticky nswe
//
//	TLabel
//
// Layout: Label.border -sticky nswe -border 1 -children {Label.padding -sticky nswe -border 1 -children {Label.label -sticky nswe}}
//
//	TLabelframe
//
// Layout: Labelframe.border -sticky nswe
//
//	TMenubutton
//
// Layout: Menubutton.border -sticky nswe -children {Menubutton.focus -sticky nswe -children {Menubutton.indicator -side right -sticky {} Menubutton.padding -sticky we -children {Menubutton.label -side left -sticky {}}}}
//
// Style map: -arrowcolor {disabled #a3a3a3}
//
//	TNotebook
//
// Layout: Notebook.client -sticky nswe
//
//	TNotebook.Tab
//
// Layout: Notebook.tab -sticky nswe -children {Notebook.padding -side top -sticky nswe -children {Notebook.focus -side top -sticky nswe -children {Notebook.label -side top -sticky {}}}}
//
// Style map: -highlightcolor {selected #4a6984} -highlight {selected 1} -background {selected #edeceb}
//
//	TPanedwindow
//
// Layout: Panedwindow.background -sticky {}
//
//	TProgressbar
//
// -
//
//	TRadiobutton
//
// Layout: Radiobutton.padding -sticky nswe -children {Radiobutton.indicator -side left -sticky {} Radiobutton.focus -side left -sticky {} -children {Radiobutton.label -sticky nswe}}
//
// Style map: -indicatorbackground {{alternate disabled} #a3a3a3 {alternate pressed} #5895bc alternate #4a6984 {selected disabled} #a3a3a3 {selected pressed} #5895bc selected #4a6984 disabled #edeceb pressed #c3c3c3}
//
//	TScale
//
// Style map: -outercolor {active #ececec}
//
//	TScrollbar
//
// Style map: -arrowcolor {disabled #a3a3a3}
//
//	TSeparator
//
// Layout: Separator.separator -sticky nswe
//
//	TSizegrip
//
// Layout: Sizegrip.sizegrip -side bottom -sticky se
//
//	TSpinbox
//
// Layout: Spinbox.field -side top -sticky we -children {null -side right -sticky {} -children {Spinbox.uparrow -side top -sticky e Spinbox.downarrow -side bottom -sticky e} Spinbox.padding -sticky nswe -children {Spinbox.textarea -sticky nswe}}
//
// Style map: -fieldbackground {readonly #edeceb disabled #edeceb} -arrowcolor {disabled #a3a3a3}
//
//	Tab
//
// Layout: Notebook.tab -sticky nswe -children {Notebook.padding -side top -sticky nswe -children {Notebook.focus -side top -sticky nswe -children {Notebook.label -side top -sticky {}}}}
//
//	Toolbutton
//
// Layout: Toolbutton.border -sticky nswe -children {Toolbutton.focus -sticky nswe -children {Toolbutton.padding -sticky nswe -children {Toolbutton.label -sticky nswe}}}
//
// Style map: -relief {disabled flat selected sunken pressed sunken active raised} -background {pressed #c3c3c3 active #ececec}
//
//	Treeview
//
// Layout: Treeview.field -sticky nswe -border 1 -children {Treeview.padding -sticky nswe -children {Treeview.treearea -sticky nswe}}
//
// Style map: -foreground {disabled #a3a3a3 selected #ffffff} -background {disabled #edeceb selected #4a6984}
//
//	Treeview.Separator
//
// Layout: Treeitem.separator -sticky nswe
//
//	Vertical.Sash
//
// Layout: Sash.vsash -sticky ns
//
//	Vertical.TProgressbar
//
// Layout: Vertical.Progressbar.trough -sticky nswe -children {Vertical.Progressbar.pbar -side bottom -sticky we}
//
//	Vertical.TScale
//
// Layout: Vertical.Scale.focus -sticky nswe -children {Vertical.Scale.padding -sticky nswe -children {Vertical.Scale.trough -sticky nswe -children {Vertical.Scale.slider -side top -sticky {}}}}
//
//	Vertical.TScrollbar
//
// Layout: Vertical.Scrollbar.trough -sticky ns -children {Vertical.Scrollbar.uparrow -side top -sticky {} Vertical.Scrollbar.downarrow -side bottom -sticky {} Vertical.Scrollbar.thumb -sticky nswe}PASS
//
// # "vista" theme style guide
//
// "Combobox.background" style element options:
//
// "Combobox.border" style element options:
//
// "Combobox.rightdownarrow" style element options:
//
// "ComboboxPopdownFrame.background" style element options:
//
// "Entry.background" style element options:
//
// "Entry.field" style element options:
//
// "Horizontal.Progressbar.pbar" style element options:
//
// "Horizontal.Scale.slider" style element options:
//
// "Horizontal.Scrollbar.grip" style element options:
//
// "Horizontal.Scrollbar.leftarrow" style element options:
//
// "Horizontal.Scrollbar.rightarrow" style element options:
//
// "Horizontal.Scrollbar.thumb" style element options:
//
// "Horizontal.Scrollbar.trough" style element options:
//
// "Menubutton.dropdown" style element options:
//
// "Spinbox.background" style element options:
//
// "Spinbox.downarrow" style element options:
//
// "Spinbox.field" style element options:
//
// "Spinbox.innerbg" style element options:
//
// "Spinbox.uparrow" style element options:
//
// "Vertical.Progressbar.pbar" style element options:
//
// "Vertical.Scale.slider" style element options:
//
// "Vertical.Scrollbar.downarrow" style element options:
//
// "Vertical.Scrollbar.grip" style element options:
//
// "Vertical.Scrollbar.thumb" style element options:
//
// "Vertical.Scrollbar.trough" style element options:
//
// "Vertical.Scrollbar.uparrow" style element options:
//
// "vista" theme style list
//
//	.
//
// Style map: -foreground {disabled SystemGrayText}
//
//	ComboboxPopdownFrame
//
// Layout: ComboboxPopdownFrame.background -sticky nswe -border 1 -children {ComboboxPopdownFrame.padding -sticky nswe}
//
//	Heading
//
// Layout: Treeheading.cell -sticky nswe Treeheading.border -sticky nswe -children {Treeheading.padding -sticky nswe -children {Treeheading.image -side right -sticky {} Treeheading.text -sticky we}}
//
//	Horizontal.TProgressbar
//
// Layout: Horizontal.Progressbar.trough -sticky nswe -children {Horizontal.Progressbar.pbar -side left -sticky ns Horizontal.Progressbar.ctext -sticky nswe}
//
//	Horizontal.TScale
//
// Layout: Scale.focus -sticky nswe -children {Horizontal.Scale.trough -sticky nswe -children {Horizontal.Scale.track -sticky we Horizontal.Scale.slider -side left -sticky {}}}
//
//	Item
//
// Layout: Treeitem.padding -sticky nswe -children {Treeitem.indicator -side left -sticky {} Treeitem.image -side left -sticky {} Treeitem.text -sticky nswe}
//
//	Label
//
// Layout: Label.fill -sticky nswe -children {Label.text -sticky nswe}
//
//	Separator
//
// Layout: Treeitem.separator -sticky nswe
//
//	TButton
//
// Layout: Button.button -sticky nswe -children {Button.focus -sticky nswe -children {Button.padding -sticky nswe -children {Button.label -sticky nswe}}}
//
//	TCheckbutton
//
// Layout: Checkbutton.padding -sticky nswe -children {Checkbutton.indicator -side left -sticky {} Checkbutton.focus -side left -sticky w -children {Checkbutton.label -sticky nswe}}
//
//	TCombobox
//
// Layout: Combobox.border -sticky nswe -children {Combobox.rightdownarrow -side right -sticky ns Combobox.padding -sticky nswe -children {Combobox.background -sticky nswe -children {Combobox.focus -sticky nswe -children {Combobox.textarea -sticky nswe}}}}
//
// Style map: -focusfill {{readonly focus} SystemHighlight} -foreground {disabled SystemGrayText {readonly focus} SystemHighlightText} -selectforeground {!focus SystemWindowText} -selectbackground {!focus SystemWindow}
//
//	TEntry
//
// Layout: Entry.field -sticky nswe -children {Entry.background -sticky nswe -children {Entry.padding -sticky nswe -children {Entry.textarea -sticky nswe}}}
//
// Style map: -selectforeground {!focus SystemWindowText} -selectbackground {!focus SystemWindow}
//
//	TLabelframe.Label
//
// Layout: Label.fill -sticky nswe -children {Label.text -sticky nswe}
//
//	TMenubutton
//
// Layout: Menubutton.dropdown -side right -sticky ns Menubutton.button -sticky nswe -children {Menubutton.padding -sticky we -children {Menubutton.label -sticky {}}}
//
//	TNotebook
//
// Layout: Notebook.client -sticky nswe
//
//	TNotebook.Tab
//
// Layout: Notebook.tab -sticky nswe -children {Notebook.padding -side top -sticky nswe -children {Notebook.focus -side top -sticky nswe -children {Notebook.label -side top -sticky {}}}}
//
// Style map: -expand {selected {2 2 2 2}}
//
//	TProgressbar
//
// -
//
//	TRadiobutton
//
// Layout: Radiobutton.padding -sticky nswe -children {Radiobutton.indicator -side left -sticky {} Radiobutton.focus -side left -sticky {} -children {Radiobutton.label -sticky nswe}}
//
//	TScale
//
// -
//
//	TSpinbox
//
// Layout: Spinbox.field -sticky nswe -children {Spinbox.background -sticky nswe -children {Spinbox.padding -sticky nswe -children {Spinbox.innerbg -sticky nswe -children {Spinbox.textarea -sticky nswe}} Spinbox.uparrow -side top -sticky nse Spinbox.downarrow -side bottom -sticky nse}}
//
// Style map: -selectforeground {!focus SystemWindowText} -selectbackground {!focus SystemWindow}
//
//	Tab
//
// Layout: Notebook.tab -sticky nswe -children {Notebook.padding -side top -sticky nswe -children {Notebook.focus -side top -sticky nswe -children {Notebook.label -side top -sticky {}}}}
//
//	Toolbutton
//
// Layout: Toolbutton.border -sticky nswe -children {Toolbutton.focus -sticky nswe -children {Toolbutton.padding -sticky nswe -children {Toolbutton.label -sticky nswe}}}
//
//	Treeview
//
// Layout: Treeview.field -sticky nswe -border 1 -children {Treeview.padding -sticky nswe -children {Treeview.treearea -sticky nswe}}
//
// Style map: -foreground {disabled SystemGrayText selected SystemHighlightText} -background {disabled SystemButtonFace selected SystemHighlight}
//
//	Treeview.Separator
//
// Layout: Treeitem.separator -sticky nswe
//
//	Vertical.TProgressbar
//
// Layout: Vertical.Progressbar.trough -sticky nswe -children {Vertical.Progressbar.pbar -side bottom -sticky we}
//
//	Vertical.TScale
//
// Layout: Scale.focus -sticky nswe -children {Vertical.Scale.trough -sticky nswe -children {Vertical.Scale.track -sticky ns Vertical.Scale.slider -side top -sticky {}}}
//
// # "winnative" theme style guide
//
// "Button.border" style element options:
//
//   - [Default]
//   - [Highlightcolor]
//   - [Relief]
//
// "Checkbutton.indicator" style element options:
//
// "Combobox.focus" style element options:
//
//   - [Focusfill]
//
// "ComboboxPopdownFrame.border" style element options:
//
//   - [Background]
//   - [Borderwidth]
//   - [Relief]
//
// "Radiobutton.indicator" style element options:
//
// "Scrollbar.trough" style element options:
//
// "Spinbox.downarrow" style element options:
//
// "Spinbox.uparrow" style element options:
//
// "border" style element options:
//
//   - [Relief]
//
// "client" style element options:
//
// "downarrow" style element options:
//
// "field" style element options:
//
//   - [Fieldbackground]
//
// "focus" style element options:
//
// "leftarrow" style element options:
//
// "rightarrow" style element options:
//
// "sizegrip" style element options:
//
// "slider" style element options:
//
//   - [Orient]
//
// "tab" style element options:
//
//   - [Background]
//   - [Borderwidth]
//
// "thumb" style element options:
//
//   - [Orient]
//
// "uparrow" style element options:
//
// "winnative" theme style list
//
//	.
//
// Style map: -foreground {disabled SystemGrayText} -embossed {disabled 1}
//
//	ComboboxPopdownFrame
//
// Layout: ComboboxPopdownFrame.border -sticky nswe
//
//	Heading
//
// Layout: Treeheading.cell -sticky nswe Treeheading.border -sticky nswe -children {Treeheading.padding -sticky nswe -children {Treeheading.image -side right -sticky {} Treeheading.text -sticky we}}
//
//	Item
//
// Layout: Treeitem.padding -sticky nswe -children {Treeitem.indicator -side left -sticky {} Treeitem.image -side left -sticky {} Treeitem.text -sticky nswe}
//
//	Label
//
// Layout: Label.fill -sticky nswe -children {Label.text -sticky nswe}
//
//	Separator
//
// Layout: Treeitem.separator -sticky nswe
//
//	TButton
//
// Layout: Button.border -sticky nswe -children {Button.padding -sticky nswe -children {Button.label -sticky nswe}}
//
// Style map: -relief {{!disabled pressed} sunken}
//
//	TCheckbutton
//
// Layout: Checkbutton.padding -sticky nswe -children {Checkbutton.indicator -side left -sticky {} Checkbutton.focus -side left -sticky w -children {Checkbutton.label -sticky nswe}}
//
//	TCombobox
//
// Layout: Combobox.field -sticky nswe -children {Combobox.downarrow -side right -sticky ns Combobox.padding -sticky nswe -children {Combobox.focus -sticky nswe -children {Combobox.textarea -sticky nswe}}}
//
// Style map: -focusfill {{readonly focus} SystemHighlight} -foreground {disabled SystemGrayText {readonly focus} SystemHighlightText} -selectforeground {!focus SystemWindowText} -fieldbackground {readonly SystemButtonFace disabled SystemButtonFace} -selectbackground {!focus SystemWindow}
//
//	TEntry
//
// Layout: Entry.field -sticky nswe -border 1 -children {Entry.padding -sticky nswe -children {Entry.textarea -sticky nswe}}
//
// Style map: -selectforeground {!focus SystemWindowText} -selectbackground {!focus SystemWindow} -fieldbackground {readonly SystemButtonFace disabled SystemButtonFace}
//
//	TLabelframe
//
// Layout: Labelframe.border -sticky nswe
//
//	TLabelframe.Label
//
// Layout: Label.fill -sticky nswe -children {Label.text -sticky nswe}
//
//	TMenubutton
//
// Layout: Menubutton.border -sticky nswe -children {Menubutton.focus -sticky nswe -children {Menubutton.indicator -side right -sticky {} Menubutton.padding -sticky we -children {Menubutton.label -side left -sticky {}}}}
//
//	TNotebook
//
// Layout: Notebook.client -sticky nswe
//
//	TNotebook.Tab
//
// Layout: Notebook.tab -sticky nswe -children {Notebook.padding -side top -sticky nswe -children {Notebook.focus -side top -sticky nswe -children {Notebook.label -side top -sticky {}}}}
//
// Style map: -expand {selected {2 2 2 0}}
//
//	TProgressbar
//
// -
//
//	TRadiobutton
//
// Layout: Radiobutton.padding -sticky nswe -children {Radiobutton.indicator -side left -sticky {} Radiobutton.focus -side left -sticky {} -children {Radiobutton.label -sticky nswe}}
//
//	TScale
//
// -
//
//	TSpinbox
//
// Layout: Spinbox.field -side top -sticky we -children {null -side right -sticky {} -children {Spinbox.uparrow -side top -sticky e Spinbox.downarrow -side bottom -sticky e} Spinbox.padding -sticky nswe -children {Spinbox.textarea -sticky nswe}}
//
//	Tab
//
// Layout: Notebook.tab -sticky nswe -children {Notebook.padding -side top -sticky nswe -children {Notebook.focus -side top -sticky nswe -children {Notebook.label -side top -sticky {}}}}
//
//	Toolbutton
//
// Layout: Toolbutton.border -sticky nswe -children {Toolbutton.focus -sticky nswe -children {Toolbutton.padding -sticky nswe -children {Toolbutton.label -sticky nswe}}}
//
// Style map: -relief {disabled flat selected sunken pressed sunken active raised}
//
//	Treeview
//
// Layout: Treeview.field -sticky nswe -border 1 -children {Treeview.padding -sticky nswe -children {Treeview.treearea -sticky nswe}}
//
// Style map: -foreground {disabled SystemGrayText selected SystemHighlightText} -background {disabled SystemButtonFace selected SystemHighlight}
//
//	Treeview.Separator
//
// Layout: Treeitem.separator -sticky nswe
//
// # "xpnative" theme style guide
//
// "Button.button" style element options:
//
// "Checkbutton.indicator" style element options:
//
// "Combobox.downarrow" style element options:
//
// "Combobox.field" style element options:
//
// "Entry.field" style element options:
//
// "Horizontal.Progressbar.pbar" style element options:
//
// "Horizontal.Progressbar.trough" style element options:
//
// "Horizontal.Scale.slider" style element options:
//
// "Horizontal.Scale.track" style element options:
//
// "Horizontal.Scrollbar.grip" style element options:
//
// "Horizontal.Scrollbar.thumb" style element options:
//
// "Horizontal.Scrollbar.trough" style element options:
//
// "Labelframe.border" style element options:
//
// "Menubutton.button" style element options:
//
// "Menubutton.dropdown" style element options:
//
// "NotebookPane.background" style element options:
//
// "Radiobutton.indicator" style element options:
//
// "Scale.trough" style element options:
//
// "Scrollbar.downarrow" style element options:
//
// "Scrollbar.leftarrow" style element options:
//
// "Scrollbar.rightarrow" style element options:
//
// "Scrollbar.uparrow" style element options:
//
// "Spinbox.downarrow" style element options:
//
// "Spinbox.field" style element options:
//
// "Spinbox.uparrow" style element options:
//
// "Toolbutton.border" style element options:
//
// "Treeheading.border" style element options:
//
// "Treeitem.indicator" style element options:
//
// "Treeview.field" style element options:
//
// "Vertical.Progressbar.pbar" style element options:
//
// "Vertical.Progressbar.trough" style element options:
//
// "Vertical.Scale.slider" style element options:
//
// "Vertical.Scale.track" style element options:
//
// "Vertical.Scrollbar.grip" style element options:
//
// "Vertical.Scrollbar.thumb" style element options:
//
// "Vertical.Scrollbar.trough" style element options:
//
// "client" style element options:
//
// "sizegrip" style element options:
//
// "tab" style element options:
//
// "xpnative" theme style list
//
//	.
//
// Style map: -foreground {disabled SystemGrayText}
//
//	Heading
//
// Layout: Treeheading.cell -sticky nswe Treeheading.border -sticky nswe -children {Treeheading.padding -sticky nswe -children {Treeheading.image -side right -sticky {} Treeheading.text -sticky we}}
//
//	Horizontal.TScale
//
// Layout: Scale.focus -sticky nswe -children {Horizontal.Scale.trough -sticky nswe -children {Horizontal.Scale.track -sticky we Horizontal.Scale.slider -side left -sticky {}}}
//
//	Horizontal.TScrollbar
//
// Layout: Horizontal.Scrollbar.trough -sticky we -children {Horizontal.Scrollbar.leftarrow -side left -sticky {} Horizontal.Scrollbar.rightarrow -side right -sticky {} Horizontal.Scrollbar.thumb -sticky nswe -unit 1 -children {Horizontal.Scrollbar.grip -sticky {}}}
//
//	Item
//
// Layout: Treeitem.padding -sticky nswe -children {Treeitem.indicator -side left -sticky {} Treeitem.image -side left -sticky {} Treeitem.text -sticky nswe}
//
//	Label
//
// Layout: Label.fill -sticky nswe -children {Label.text -sticky nswe}
//
//	Separator
//
// Layout: Treeitem.separator -sticky nswe
//
//	TButton
//
// Layout: Button.button -sticky nswe -children {Button.focus -sticky nswe -children {Button.padding -sticky nswe -children {Button.label -sticky nswe}}}
//
//	TCheckbutton
//
// Layout: Checkbutton.padding -sticky nswe -children {Checkbutton.indicator -side left -sticky {} Checkbutton.focus -side left -sticky w -children {Checkbutton.label -sticky nswe}}
//
//	TCombobox
//
// Layout: Combobox.field -sticky nswe -children {Combobox.downarrow -side right -sticky ns Combobox.padding -sticky nswe -children {Combobox.focus -sticky nswe -children {Combobox.textarea -sticky nswe}}}
//
// Style map: -focusfill {{readonly focus} SystemHighlight} -foreground {disabled SystemGrayText {readonly focus} SystemHighlightText} -selectforeground {!focus SystemWindowText} -selectbackground {!focus SystemWindow}
//
//	TEntry
//
// Layout: Entry.field -sticky nswe -border 1 -children {Entry.padding -sticky nswe -children {Entry.textarea -sticky nswe}}
//
// Style map: -selectforeground {!focus SystemWindowText} -selectbackground {!focus SystemWindow}
//
//	TLabelframe.Label
//
// Layout: Label.fill -sticky nswe -children {Label.text -sticky nswe}
//
//	TMenubutton
//
// Layout: Menubutton.dropdown -side right -sticky ns Menubutton.button -sticky nswe -children {Menubutton.padding -sticky we -children {Menubutton.label -sticky {}}}
//
//	TNotebook
//
// Layout: Notebook.client -sticky nswe
//
//	TNotebook.Tab
//
// Layout: Notebook.tab -sticky nswe -children {Notebook.padding -side top -sticky nswe -children {Notebook.focus -side top -sticky nswe -children {Notebook.label -side top -sticky {}}}}
//
// Style map: -expand {selected {2 2 2 2}}
//
//	TProgressbar
//
// -
//
//	TRadiobutton
//
// Layout: Radiobutton.padding -sticky nswe -children {Radiobutton.indicator -side left -sticky {} Radiobutton.focus -side left -sticky {} -children {Radiobutton.label -sticky nswe}}
//
//	TScale
//
// -
//
//	TScrollbar
//
// -
//
//	TSpinbox
//
// Layout: Spinbox.field -side top -sticky we -children {null -side right -sticky {} -children {Spinbox.uparrow -side top -sticky e Spinbox.downarrow -side bottom -sticky e} Spinbox.padding -sticky nswe -children {Spinbox.textarea -sticky nswe}}
//
// Style map: -selectforeground {!focus SystemWindowText} -selectbackground {!focus SystemWindow}
//
//	Tab
//
// Layout: Notebook.tab -sticky nswe -children {Notebook.padding -side top -sticky nswe -children {Notebook.focus -side top -sticky nswe -children {Notebook.label -side top -sticky {}}}}
//
//	Toolbutton
//
// Layout: Toolbutton.border -sticky nswe -children {Toolbutton.focus -sticky nswe -children {Toolbutton.padding -sticky nswe -children {Toolbutton.label -sticky nswe}}}
//
//	Treeview
//
// Layout: Treeview.field -sticky nswe -border 1 -children {Treeview.padding -sticky nswe -children {Treeview.treearea -sticky nswe}}
//
// Style map: -foreground {disabled SystemGrayText selected SystemHighlightText} -background {disabled SystemButtonFace selected SystemHighlight}
//
//	Treeview.Separator
//
// Layout: Treeitem.separator -sticky nswe
//
//	Vertical.TScale
//
// Layout: Scale.focus -sticky nswe -children {Vertical.Scale.trough -sticky nswe -children {Vertical.Scale.track -sticky ns Vertical.Scale.slider -side top -sticky {}}}
//
//	Vertical.TScrollbar
//
// Layout: Vertical.Scrollbar.trough -sticky ns -children {Vertical.Scrollbar.uparrow -side top -sticky {} Vertical.Scrollbar.downarrow -side bottom -sticky {} Vertical.Scrollbar.thumb -sticky nswe -unit 1 -children {Vertical.Scrollbar.grip -sticky {}}}PASS
//
// [FreeBSD]: https://github.com/ebitengine/purego/blob/7402fed73989eaf478f4f7905862d0f04537ac8c/internal/fakecgo/freebsd.go#L15
// [MVP]: https://en.wikipedia.org/wiki/Minimum_viable_product
// [RERO]: https://en.wikipedia.org/wiki/Release_early,_release_often
// [Tkinter]: https://en.wikipedia.org/wiki/Tkinter
// [equ]: https://pkg.go.dev/modernc.org/equ
// [issue tracker]: https://gitlab.com/cznic/tk9.0/-/issues
// [jnml's LiberaPay]: https://liberapay.com/jnml/donate
// [modern-c.appspot.com]: https://modern-c.appspot.com/-/builder/?importpath=modernc.org%2ftk9.0
// [tcl.tk site]: https://www.tcl.tk/man/tcl9.0/TkCmd/index.html
// [tkinter.ttk site]: https://docs.python.org/3/library/tkinter.ttk.html
// [tutorialspoint]: https://www.tutorialspoint.com/tcl-tk/tk_overview.htm
package tk9_0 // import "modernc.org/tk9.0"
