// Copyright 2024 The tk9.0-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tk9_0 // import "modernc.org/tk9.0"

import (
	_ "embed"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"unsafe"

	"github.com/evilsocket/islazy/zip"
	"golang.org/x/sys/windows"
	"modernc.org/memory"
)

// trcw prints and return caller's position and an optional message tagged with TRC. Output goes to stderr.
//
//lint:ignore U1000 debug helper
func trcw(s string, args ...interface{}) string {
	switch {
	case s == "":
		s = fmt.Sprintf(strings.Repeat("%v ", len(args)), args...)
	default:
		s = fmt.Sprintf(s, args...)
	}
	r := fmt.Sprintf("%s: TRC(id=%v tid=%v) %s", origin(2), goroutineID(), windows.GetCurrentThreadId(), s)
	fmt.Fprintf(os.Stderr, "%s\n", r)
	os.Stderr.Sync()
	return r
}

var (
	// No mutex, the package must be used by a single goroutine only.
	allocator memory.Allocator

	createCommandProc *windows.Proc
	deleteCommandProc *windows.Proc
	evalExProc        *windows.Proc
	getObjResultProc  *windows.Proc
	getStringProc     *windows.Proc
	interp            uintptr
	newStringObjProc  *windows.Proc
	runCmdProxy       = windows.NewCallbackCDecl(eventDispatcher)
	setObjResultProc  *windows.Proc
	splitListProc     *windows.Proc
	tclDll            *windows.DLL
	tkDll             *windows.DLL
)

func lazyInit() {
	if initialized {
		return
	}

	runtime.LockOSThread()
	// trcw("LockOSThread")
	initialized = true

	defer commonLazyInit()

	var cacheDir string
	if cacheDir, Error = getCacheDir(); Error != nil {
		return
	}

	if bindLibs(cacheDir); Error != nil {
		return
	}

	var nm uintptr
	if nm, Error = cString("eventDispatcher"); Error != nil {
		return
	}

	cmd, _, _ := createCommandProc.Call(interp, nm, runCmdProxy, 0, 0)
	if cmd == 0 {
		Error = fmt.Errorf("registering event dispatcher proxy failed: %v", getObjResultProc)
		return
	}

	setDefaults()
}

func bindLibs(cacheDir string) {
	var wd string
	if wd, Error = os.Getwd(); Error != nil {
		return
	}

	defer func() {
		Error = errors.Join(Error, os.Chdir(wd))
	}()

	if Error = os.Chdir(cacheDir); Error != nil {
		return
	}

	if tclDll, Error = windows.LoadDLL(tclBin); Error != nil {
		return
	}

	if tkDll, Error = windows.LoadDLL(tkBin); Error != nil {
		return
	}

	var tclCreateInterp, tclInit, tkInit *windows.Proc
	if tclCreateInterp, Error = tclDll.FindProc("Tcl_CreateInterp"); Error != nil {
		return
	}

	if tclInit, Error = tclDll.FindProc("Tcl_Init"); Error != nil {
		return
	}

	if createCommandProc, Error = tclDll.FindProc("Tcl_CreateCommand"); Error != nil {
		return
	}

	if deleteCommandProc, Error = tclDll.FindProc("Tcl_DeleteCommand"); Error != nil {
		return
	}

	if evalExProc, Error = tclDll.FindProc("Tcl_EvalEx"); Error != nil {
		return
	}

	if setObjResultProc, Error = tclDll.FindProc("Tcl_SetObjResult"); Error != nil {
		return
	}

	if splitListProc, Error = tclDll.FindProc("Tcl_SplitList"); Error != nil {
		return
	}

	if getObjResultProc, Error = tclDll.FindProc("Tcl_GetObjResult"); Error != nil {
		return
	}

	if getStringProc, Error = tclDll.FindProc("Tcl_GetString"); Error != nil {
		return
	}

	if newStringObjProc, Error = tclDll.FindProc("Tcl_NewStringObj"); Error != nil {
		return
	}

	if tkInit, Error = tkDll.FindProc("Tk_Init"); Error != nil {
		return
	}

	if interp, _, _ = tclCreateInterp.Call(); interp == 0 {
		Error = fmt.Errorf("failed to create a Tcl interpreter")
		return
	}

	if r, _, _ := tclInit.Call(interp); r != tcl_ok {
		Error = fmt.Errorf("failed to initialize the Tcl interpreter")
		return
	}

	if _, Error := eval(fmt.Sprintf("zipfs mount lib%s.zip /app", libVersion)); Error != nil {
		return
	}

	if r, _, _ := tkInit.Call(interp); r != tcl_ok {
		Error = fmt.Errorf("failed to initialize Tk")
		return
	}

	for _, dll := range moreDLLs {
		var handle *windows.DLL
		if handle, Error = windows.LoadDLL(dll.dll); Error != nil {
			return
		}

		var initProc *windows.Proc
		if initProc, Error = handle.FindProc(dll.initProc); Error != nil {
			return
		}

		r, _, _ := initProc.Call(interp)
		if r != tcl_ok {
			Error = fmt.Errorf("failed to initialize %s: %s", dll.dll, tclResult())
			return
		}
	}
}

func getCacheDir() (r string, err error) {
	if r, err = os.UserCacheDir(); err != nil {
		return "", err
	}

	r0 := filepath.Join(r, "modernc.org", libVersion, goos)
	r = filepath.Join(r0, goarch)
	fi, err := os.Stat(r)
	if err == nil && fi.IsDir() {
		if checkSig(r, shasig) {
			return r, nil
		}

		os.RemoveAll(r) // Tampered or corrupted.
	}

	os.MkdirAll(r0, 0700)
	tmp, err := os.MkdirTemp(r0, "")
	if err != nil {
		return "", err
	}

	zf := filepath.Join(tmp, "lib.zip")
	if err = os.WriteFile(zf, libZip, 0660); err != nil {
		return "", err
	}

	if _, err = zip.Unzip(zf, tmp); err != nil {
		os.Remove(zf)
		return "", err
	}

	os.Remove(zf)
	if err = os.Rename(tmp, r); err == nil {
		return r, nil
	}

	cleanupDirs = append(cleanupDirs, tmp)
	return tmp, nil
}

func tclResult() string {
	r, _, _ := getObjResultProc.Call(interp)
	if r == 0 {
		return ""
	}

	if r, _, _ = getStringProc.Call(r); r != 0 {
		return goString(r)
	}

	return ""
}

func cString(s string) (r uintptr, err error) {
	if s == "" {
		return 0, nil
	}

	if r, err = allocator.UintptrMalloc(len(s) + 1); err != nil {
		return 0, err
	}

	copy(unsafe.Slice((*byte)(unsafe.Pointer(r)), len(s)), s)
	*(*byte)(unsafe.Add(unsafe.Pointer(r), len(s))) = 0
	return r, nil
}

func setResult(s string) (err error) {
	cs, err := cString(s)
	if err != nil {
		return err
	}

	defer allocator.UintptrFree(cs)

	obj, _, _ := newStringObjProc.Call(cs, uintptr(len(s)))
	if obj == 0 {
		return fmt.Errorf("OOM")
	}

	setObjResultProc.Call(interp, obj)
	return nil
}

func goTransientString(p uintptr) (r string) { // Result cannot be retained.
	if p == 0 {
		return ""
	}

	var n uintptr
	for p := p; *(*byte)(unsafe.Pointer(p + n)) != 0; n++ {
	}
	return string(unsafe.Slice((*byte)(unsafe.Pointer(p)), n))
}

func eval(code string) (r string, err error) {
	if dmesgs {
		defer func() {
			dmesg("code=%s -> r=%v err=%v", code, r, err)
		}()
	}

	if !initialized {
		lazyInit()
		if Error != nil {
			return "", Error
		}
	}

	cs, err := cString(code)
	if err != nil {
		return "", err
	}

	defer allocator.UintptrFree(cs)

	// trcw("code=%s", code)
	switch r0, _, _ := evalExProc.Call(interp, cs, uintptr(len(code)), tcl_eval_direct); r0 {
	case tcl_ok, tcl_return:
		// trcw("%s->%s, nil", code, tclResult())
		return tclResult(), nil
	default:
		// trcw("%s->{}, %s", code, tclResult())
		return "", fmt.Errorf("%s", tclResult())
	}
}

func eventDispatcher(clientData, in uintptr, argc int32, argv uintptr) uintptr {
	if argc < 2 {
		setResult(fmt.Sprintf("eventDispatcher internal error: argc=%v", argc))
		return tcl_error
	}

	arg1 := goTransientString(*(*uintptr)(unsafe.Pointer(argv + unsafe.Sizeof(uintptr(0)))))
	id, e, err := newEvent(arg1)
	if err != nil {
		setResult(fmt.Sprintf("eventDispatcher internal error: argv[1]=%q, err=%v", arg1, err))
		return tcl_error
	}

	h := handlers[int32(id)]
	e.W = h.w
	for i := int32(2); i < argc; i++ {
		e.args = append(e.args, goString(*(*uintptr)(unsafe.Pointer(argv + uintptr(i)*unsafe.Sizeof(uintptr(0))))))
	}
	switch h.callback(e); {
	case e.Err != nil:
		setResult(tclSafeString(e.Err.Error()))
		return tcl_error
	default:
		if setResult(e.Result) != nil {
			return tcl_error
		}

		return uintptr(e.returnCode)
	}
}

// Finalize releases all resources held, if any. This may include temporary
// files. Finalize is intended to be called on process shutdown only.
func Finalize() (err error) {
	if finished.Swap(1) != 0 {
		return
	}

	defer runtime.UnlockOSThread()

	for _, v := range cleanupDirs {
		err = errors.Join(err, os.RemoveAll(v))
	}
	return err
}

func callSplitList(cList uintptr, argcPtr uintptr, argvPtr uintptr) (r1 uintptr, r2 uintptr, err error) {
	return splitListProc.Call(interp, cList, argcPtr, argvPtr)
}

// Internal malloc enabling parseList() in tk.go to not care about the target
// specific implemetations.
func malloc(sz int) (r uintptr, err error) {
	return allocator.UintptrMalloc(sz)
}

// Internal free enabling parseList() in tk.go to not care about the target
// specific implemetations.
func free(p uintptr) (err error) {
	return allocator.UintptrFree(p)
}
