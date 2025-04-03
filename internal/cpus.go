// Copyright 2025 The tk9.0-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore

package main

import "runtime"

func main() {
	println(runtime.GOMAXPROCS(-1))
}
