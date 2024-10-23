// Copyright 2024 The tk9.0-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tk9_0 // import "modernc.org/tk9.0"

import _ "embed"

const (
	tclBin = "libtcl9.0.dylib"
	tkBin  = "libtcl9tk9.0.dylib"
)

//go:embed embed/darwin/amd64/lib.zip
var libZip []byte

var shasig = map[string]string{
	// embed/darwin/amd64/lib.zip
	"libtcl9.0.dylib":    "cff8126844e7628e0141a914ba900c7be001712c7ce4e0ed409bf8d1db041c48",
	"libtcl9tk9.0.dylib": "7f75e08ec9b6226900d226859ae8a9068718a29cf85601d58f33af6389eadd26",
	"libtk9.0.0.zip":     "83cb3515f07b5adc1aa8876960c92ba25bf1ceb76e02e26b6792c7a3ecf4a5ce",
}
