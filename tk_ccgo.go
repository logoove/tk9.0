// Copyright 2024 The tk9.0-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build linux && (386 || arm || loong64 || ppc64le || riscv64 || s390x)

package tk9_0 // import "modernc.org/tk9.0"

import (
	_ "embed"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"modernc.org/libc"
	libtcl "modernc.org/libtcl9.0"
	tcllib "modernc.org/libtcl9.0/library"
	libtk "modernc.org/libtk9.0"
	tklib "modernc.org/libtk9.0/library"
	tcl "modernc.org/tcl9.0"
)

const (
	tclLibZip        = "tcl_library.zip"
	tclLibMountPoint = "/lib/tcl"
	tkLibZip         = "tk_library.zip"
	tkLibMountPoint  = "/lib/tk"
)

var (
	interp *tcl.Interp

	shasig = map[string]string{
		// other
		"tcl_library.zip": "ef851d549039c822cd06af0c657c8173006eae90f997bdae11c60c0bdc5a0c1c",
		"tk_library.zip":  "2afaf3ccb4521fe44d330b4da077d7d433d377f9ffc56f5ce8decd1689e00352",
	}
)

func init() {
	if isBuilder {
		return
	}

	runtime.LockOSThread()
}

func lazyInit() {
	if initialized {
		return
	}

	initialized = true

	defer commonLazyInit()

	var cacheDir string
	if cacheDir, Error = getCacheDir(); Error != nil {
		return
	}

	tls := libc.NewTLS()
	zf1 := filepath.Join(cacheDir, tclLibZip)
	zf2 := filepath.Join(cacheDir, tkLibZip)
	var cs uintptr
	if cs, Error = libc.CString(fmt.Sprintf(`
zipfs mount %s %s
zipfs mount %s %s
`, zf1, tclLibMountPoint, zf2, tkLibMountPoint)); Error != nil {
		return
	}

	p := libtcl.XTcl_SetPreInitScript(tls, cs)
	if p != 0 {
		panic(todo("Tcl_SetPreInitScript internal error: %s", libc.GoString(p)))
	}

	if interp, Error = tcl.NewInterp(map[string]string{
		"tcl_library": fmt.Sprintf("//zipfs:%s/library", tclLibMountPoint),
		"tk_library":  fmt.Sprintf("//zipfs:%s/library", tkLibMountPoint),
	}); Error != nil {
		return
	}

	if rc := libtk.XTk_Init(interp.TLS(), interp.Handle()); rc != libtk.TCL_OK {
		interp.Close()
		Error = fmt.Errorf("failed to initialize the Tk subsystem")
		return
	}

	if Error = interp.RegisterCommand("eventDispatcher", eventDispatcher, nil, nil); Error == nil {
		setDefaults()
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

	zf := filepath.Join(tmp, tclLibZip)
	if err = os.WriteFile(zf, []byte(tcllib.Zip), 0660); err != nil {
		return "", err
	}

	zf = filepath.Join(tmp, tkLibZip)
	if err = os.WriteFile(zf, []byte(tklib.Zip), 0660); err != nil {
		return "", err
	}

	if err = os.Rename(tmp, r); err == nil {
		return r, nil
	}

	cleanupDirs = append(cleanupDirs, tmp)
	return tmp, nil
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

	return interp.Eval(code, tcl.EvalDirect)
}

func eventDispatcher(data any, interp *tcl.Interp, argv []string) int {
	id, e, err := newEvent(argv[1])
	if err != nil {
		interp.SetResult(fmt.Sprintf("eventDispatcher internal error: argv1=`%s`", argv[1]))
		return tcl_error
	}

	h := handlers[int32(id)]
	e.W = h.w
	if len(argv) > 2 { // eg.: ["eventDispatcher", "42", "0.1", "0.9"]
		e.args = argv[2:]
	}
	switch h.callback(e); {
	case e.Err != nil:
		interp.SetResult(tclSafeString(e.Err.Error()))
		return libtcl.TCL_ERROR
	default:
		interp.SetResult(e.Result)
		return e.returnCode
	}
}

// Finalize releases all resources held, if any. This may include temporary
// files. Finalize is intended to be called on process shutdown only.
func Finalize() (err error) {
	if finished.Swap(1) != 0 {
		return
	}

	defer runtime.UnlockOSThread()

	if interp != nil {
		err = interp.Close()
		interp = nil
	}
	for _, v := range cleanupDirs {
		err = errors.Join(err, os.RemoveAll(v))
	}
	return err
}
