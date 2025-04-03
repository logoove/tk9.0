// Copyright 2025 The tk9.0-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build linux && (386 || arm || loong64 || ppc64le || riscv64 || s390x)

package tk9_0 // import "modernc.org/tk9.0"

import (
	"fmt"

	tcl "modernc.org/tcl9.0"
)

// func (proxy *widgetProxy) eventDispatcher(clientData, in uintptr, argc int32, argv uintptr) uintptr {
func (proxy *widgetProxy) eventDispatcher(clientData any, in *tcl.Interp, args []string) int {
	// Expect at least arguments for the path and the operation.
	if len(args) < 2 {
		setResult(fmt.Sprintf("WidgetProxy eventDispatcher internal error: args=%q", args))
		return tcl_error
	}

	args = args[1:]
	operation := args[0]
	if callback, ok := proxy.operations[operation]; ok {
		// Dispatch to the operation's registered callback.
		callback(args)
	} else {
		// Process as normal.
		proxy.EvalWrapped(args)
	}

	return tcl_ok
}

// Create a new Tcl command whose name is the widget's pathname, and
// whose action is to dispatch on the operation passed to the widget:
func (proxy *widgetProxy) registerEventDispatcher() {
	if err := interp.RegisterCommand(proxy.window.String(), proxy.eventDispatcher, nil, nil); err != nil {
		fail(fmt.Errorf("registering widget proxy event dispatcher proxy failed: %v", err))
	}
}

func (proxy *widgetProxy) unregisterEventDispatcher() {
	evalErr(fmt.Sprintf("rename %s {}", proxy.window))
}
