// Copyright 2024 The tk9.0-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tk9_0 // import "modernc.org/tk9.0"

import _ "embed"

const (
	tclBin = "tcl90.dll"
	tkBin  = "tcl9tk90.dll"
)

//go:embed embed/windows/arm64/lib.zip
var libZip []byte

var shasig = map[string]string{
	// embed/windows/arm64/lib.zip
	"libtommath.dll":     "907b9c7860fc07231f1e238551715e5d813283807f52dc383dae0cb47a879d29",
	"tcl90.dll":          "688e841263aab47bb1acc823da57b644cba64dabf47fb306258f85c86ca71742",
	"tcl9dde14.dll":      "1976dcce97481fecf6e4f66639ef109598d6d23b5c9168bd76d4536777ed38b4",
	"tcl9registry13.dll": "1738fd8e4a73524d915070413282c2bc33bf6637418fab87ad03a3665b7b58b5",
	"tcl9tk90.dll":       "ea2afba46ee9952cdd537824a3b72162d06fe3fe015343df43cc87e0f75f1cf8",
	"tcldde14.dll":       "8a47b46c88ed0f6a5eb3c9d1109ab60d2596570cccc8f4b2283215bb836b7d0f",
	"tclregistry13.dll":  "c3d3e8ae5b19380c1d38896543a428199415cd2f11b866ad5f9db62487d73259",
	"zlib1.dll":          "6f10a76dcc2c831d1f08d98c0b345afa0911bec0238fcba357b612ccc6ab5d81",
	"libtk9.0.0.zip":     "2ed66d57f34e4e0be9d6526e7a0428b6ccb2644477de4a62c85639af2e7143aa",
}
