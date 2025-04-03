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
	"modernc.org/tk9.0/internal/img"
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
		"tcl_library.zip": "1849c8e8df2e23cdaf904bd04f1316be29473612a215e84c9f9f8ba144d16b2f",
		"tk_library.zip":  "ea619ae0c921446db3659cbfc4efa2c700f2531c9a20ce9029b603b629c29711",
	}
)

func lazyInit() {
	if initialized {
		return
	}

	runtime.LockOSThread()
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

	tls = interp.TLS()
	h := interp.Handle()
	if rc := libtk.XTk_Init(tls, h); rc != libtk.TCL_OK {
		interp.Close()
		Error = fmt.Errorf("failed to initialize the Tk subsystem")
		return
	}

	if Error = interp.RegisterCommand("eventDispatcher", eventDispatcher, nil, nil); Error == nil {
		setDefaults()
	}

	if rc := img.XTkimg_Init(tls, h); rc != 0 {
		Error = fmt.Errorf("failed to initialize the img subsystem: Tkimg_Init")
		return
	}

	if rc := img.XJpegtcl_Init(tls, h); rc != 0 {
		Error = fmt.Errorf("failed to initialize the img subsystem: Jpegtcl_Init")
		return
	}

	if rc := img.XTkimgjpeg_Init(tls, h); rc != 0 {
		Error = fmt.Errorf("failed to initialize the img subsystem: Tkimgjpeg_Init")
		return
	}

	if rc := img.XTkimgbmp_Init(tls, h); rc != 0 {
		Error = fmt.Errorf("failed to initialize the img subsystem: Tkimgbmp_Init")
		return
	}

	if rc := img.XTkimgico_Init(tls, h); rc != 0 {
		Error = fmt.Errorf("failed to initialize the img subsystem: Tkimgico_Init")
		return
	}

	if rc := img.XTkimgpcx_Init(tls, h); rc != 0 {
		Error = fmt.Errorf("failed to initialize the img subsystem: Tkimgpcx_Init")
		return
	}

	if rc := img.XTkimgxpm_Init(tls, h); rc != 0 {
		Error = fmt.Errorf("failed to initialize the img subsystem: Tkimgxpm_Init")
		return
	}

	if rc := img.XZlibtcl_Init(tls, h); rc != 0 {
		Error = fmt.Errorf("failed to initialize the img subsystem: Zlibtcl_Init")
		return
	}

	if rc := img.XPngtcl_Init(tls, h); rc != 0 {
		Error = fmt.Errorf("failed to initialize the img subsystem: Pngtcl_Init")
		return
	}

	if rc := img.XTkimgpng_Init(tls, h); rc != 0 {
		Error = fmt.Errorf("failed to initialize the img subsystem: Tkimgpng_Init")
		return
	}

	if rc := img.XTkimgppm_Init(tls, h); rc != 0 {
		Error = fmt.Errorf("failed to initialize the img subsystem: Tkimgppm_Init")
		return
	}

	if rc := img.XTkimgtga_Init(tls, h); rc != 0 {
		Error = fmt.Errorf("failed to initialize the img subsystem: Tkimgtga_Init")
		return
	}

	if rc := img.XTifftcl_Init(tls, h); rc != 0 {
		Error = fmt.Errorf("failed to initialize the img subsystem: Tifftcl_Init")
		return
	}

	if rc := img.XTkimgtiff_Init(tls, h); rc != 0 {
		Error = fmt.Errorf("failed to initialize the img subsystem: Tkimgtiff_Init")
		return
	}

	if rc := img.XTkimgxbm_Init(tls, h); rc != 0 {
		Error = fmt.Errorf("failed to initialize the img subsystem: Tkimgxbm_Init")
		return
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

func setResult(s string) (err error) {
	return interp.SetResult(s)
}

func cString(s string) (r uintptr, err error) {
	return libc.CString(s)
}

func callSplitList(cList uintptr, argcPtr uintptr, argvPtr uintptr) (r1 uintptr, r2 uintptr, err uintptr) {
	rc := libtcl.XTclSplitList(interp.TLS(), interp.Handle(), cList, argcPtr, argvPtr) // .SyscallN(splitListProc, interp, cList, argcPtr, argvPtr)
	if rc == tcl_error {
		err = libtcl.TCL_ERROR
	}
	return uintptr(rc), 0, err
}

var oom = errors.New("OOM")

// Internal malloc enabling parseList() in tk.go to not care about the target
// specific implemetations.
func malloc(sz int) (r uintptr, err error) {
	if r = libc.Xmalloc(interp.TLS(), libc.Tsize_t(sz)); r == 0 {
		err = oom
	}
	return r, err
}

// Internal free enabling parseList() in tk.go to not care about the target
// specific implemetations.
func free(p uintptr) (err error) {
	libc.Xfree(interp.TLS(), p)
	return nil
}
