// Copyright 2025 The tk9.0-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tk9_0 // import "modernc.org/tk9.0"

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

func (proxy *widgetProxy) eventDispatcher(clientData, in uintptr, argc int32, argv uintptr) uintptr {
	// Expect at least arguments for the path and the operation.
	if argc < 2 {
		setResult(fmt.Sprintf("WidgetProxy eventDispatcher internal error: argc=%v", argc))
		return tcl_error
	}

	argsPointers := unsafe.Slice((*uintptr)(unsafe.Pointer(argv)), argc)
	argsPointers = argsPointers[1:] // skip path
	args := make([]string, len(argsPointers))
	for i := 0; i < len(argsPointers); i++ {
		args[i] = goString(argsPointers[i])
	}

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
	runCmdProxy := windows.NewCallback(proxy.eventDispatcher)
	if proxy.commandName, Error = cString(proxy.window.String()); Error != nil {
		return
	}
	cmd, _, _ := createCommandProc.Call(interp, proxy.commandName, runCmdProxy, 0, 0)
	if cmd == 0 {
		fail(fmt.Errorf("registering widget proxy event dispatcher proxy failed: %v", getObjResultProc))
		return
	}
}

func (proxy *widgetProxy) unregisterEventDispatcher() {
	defer func() {
		allocator.UintptrFree(proxy.commandName)
		proxy.commandName = 0 // nil
	}()
	_, _, _ = deleteCommandProc.Call(interp, proxy.commandName)
}
