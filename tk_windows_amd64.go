// Copyright 2024 The tk9.0-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tk9_0 // import "modernc.org/tk9.0"

import _ "embed"

const (
	tclBin = "tcl90.dll"
	tkBin  = "tcl9tk90.dll"
)

//go:embed embed/windows/amd64/lib.zip
var libZip []byte

var shasig = map[string]string{
	// embed/windows/amd64/lib.zip
	"libtommath.dll":     "2d760fefb452665b6af8c8d9d29f3a8378f10fc0847cdd9938ea0cb5edf1d573",
	"tcl90.dll":          "ffe73bcaf947e361561d71a8d93e525c637ca4e2d0b4b4c14c5e0df9756e92b0",
	"tcl9dde14.dll":      "752748e6975bc56cb941e29c291c18db67e33216b026d49962bf62042584b50e",
	"tcl9registry13.dll": "0943f57b7bf4a5433a660f5ae6f252fa2f67603af01e916725c48f4fdd6ed658",
	"tcl9tk90.dll":       "38efd85ac62473ea9615d272b1b923ee3fe8fe9e35d6fa9ebab3f7a0d4aaf795",
	"tcldde14.dll":       "71ee58d7ddd6213fc2a3d05f089296453c4c8b3de9b7268c2256bbbbab386aab",
	"tclregistry13.dll":  "e0bfebf0377f3b1671f7d501bb14683cadf382d8181936448d98864eb137acb2",
	"zlib1.dll":          "04117778e255ed158cf6a4a1e51aa40f49124d9035208218fbfebbe565cf254d",
	"libtk9.0.0.zip":     "b8873abe39f903e8c0e6e0669083421e020646af9946e16baf30897a64dd5460",
}
