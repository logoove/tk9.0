// Copyright 2025 The tk9.0-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tk9_0 // import "modernc.org/tk9.0"

import (
	"errors"
	"fmt"
)

// OperationCallback is the signature for a WidgetProxy operation callback.
type OperationCallback func(args []string)

// widgetProxy wraps a widget window, providing the means to intercept its
// internal operations, and modify their behaviour.
// For example, it can provide access to a TextWidget's internal 'insert' and
// 'delete' operations.
type widgetProxy struct {
	window       *Window
	originalPath string
	commandName  uintptr
	operations   map[string]OperationCallback
}

// newWidgetProxy creates a proxy that facilitates hooking in to
// a widget window's internal operations.
func newWidgetProxy(window *Window) widgetProxy {
	if !extensionInitialized("eval") {
		fail(errors.New("use of WidgetProxy requires the 'eval' extension to be enabled, but it is not"))
	}

	proxy := widgetProxy{
		window:       window,
		originalPath: window.String() + "_original",
		operations:   make(map[string]OperationCallback),
	}

	// Rename the Tcl command within Tcl:
	evalErr(fmt.Sprintf("rename %s %s", window, proxy.originalPath))
	// Register an event dispatcher for the window.
	proxy.registerEventDispatcher()

	return proxy
}

// Close undoes the wrapping of its Window.
// All registered operations are unregistered.
func (proxy *widgetProxy) Close() {
	// Unregister all registered operations.
	for operation := range proxy.operations {
		proxy.Unregister(operation)
	}

	// Restore the original widget Tcl command.
	proxy.unregisterEventDispatcher()
	// Restore the Tcl command within Tcl:
	evalErr(fmt.Sprintf("rename %s %s", proxy.originalPath, proxy.window))
}

// Register registers a callback for an operation supported by the wrapped Window.
//
// The operation name is widget-specific. It is not possible to validate whether
// the operation name is support by the wrapped widget. If an unsupported
// operation name is used it will be silently ignored.
//
// The operation's arguments are passed to the callback. The callback may perform
// the operation on the wrapped widget by calling the EvalWrapped method with the
// arguments. The arguments may be modified if desired.
func (proxy *widgetProxy) Register(operation string, callback OperationCallback) {
	proxy.operations[operation] = callback
}

// Unregister unregisters an operation callback.
// If there is no registered callback for the operation, it is silently ignored.
func (proxy *widgetProxy) Unregister(operation string) {
	delete(proxy.operations, operation)
}

// EvalWrapped evaluates the arguments as a raw Tcl string against the wrapped Window.
func (proxy *widgetProxy) EvalWrapped(args []string) {
	evalErr(fmt.Sprintf("%s %s", proxy.originalPath, tclSafeStrings(args...)))
}

// TextWidgetProxy wraps a TextWidget.
// It provides the ability to intercept the widget's internal operations (such as
// 'insert' and 'delete'), and modify their behaviour.
type TextWidgetProxy struct {
	*TextWidget
	widgetProxy
}

// NewTextWidgetProxy creates a TextWidgetProxy, wrapping the
// provided TextWidget.
//
// Note that because the EvalWrapped method evaluates raw Tcl, the "eval" extension
// must be enabled in the main package before calling this function.
func NewTextWidgetProxy(widget *TextWidget) TextWidgetProxy {
	return TextWidgetProxy{
		TextWidget:  widget,
		widgetProxy: newWidgetProxy(widget.Window),
	}
}
