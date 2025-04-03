// Copyright 2025 The tk9.0-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build linux && (386 || arm || loong64 || ppc64le || riscv64 || s390x)

package img

import "unsafe"

// Workaround ccgo inability to handle cyclic initializers.
func init() {
	_tiffFields[48].Ffield_subfields = uintptr(unsafe.Pointer(&_tiffFieldArray))
}
