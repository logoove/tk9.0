// Copyright 2024 The tk9.0-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !tk.dmesg
// +build !tk.dmesg

package tk9_0 // import "modernc.org/tk9.0"

const dmesgs = false

func dmesg(s string, args ...interface{}) {}
