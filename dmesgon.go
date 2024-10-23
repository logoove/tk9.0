// Copyright 2024 The tk9.0-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build tk.dmesg
// +build tk.dmesg

package tk9_0 // import "modernc.org/tk9.0"

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const dmesgs = true

var (
	pid  = fmt.Sprintf("[%v %v] ", os.Getpid(), filepath.Base(os.Args[0]))
	logf *os.File
)

func init() {
	fn := filepath.Join(os.TempDir(), fmt.Sprintf("%s-%v-%s", filepath.Base(os.Args[0]), os.Getpid(), time.Now().Format("20060102-150405")))
	var err error
	if logf, err = os.OpenFile(fn, os.O_APPEND|os.O_CREATE|os.O_WRONLY|os.O_SYNC, 0660); err != nil {
		panic(err.Error())
	}

	dmesg("enter")
	fmt.Println(fn)
}

func dmesg(s string, args ...interface{}) {
	if s == "" {
		s = strings.Repeat("%v ", len(args))
	}
	s = fmt.Sprintf(pid+s, args...)
	s = fmt.Sprintf("%s %v", s, []string{origin(2), origin(3), origin(4)})
	switch {
	case len(s) != 0 && s[len(s)-1] == '\n':
		fmt.Fprint(logf, s)
	default:
		fmt.Fprintln(logf, s)
	}
}
