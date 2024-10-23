// Copyright 2024 The tk9.0-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tk9_0 // import "modernc.org/tk9.0"

import _ "embed"

const (
	tclBin = "libtcl9.0.so"
	tkBin  = "libtcl9tk9.0.so"
)

//go:embed embed/freebsd/amd64/lib.zip
var libZip []byte

var shasig = map[string]string{
	// embed/freebsd/amd64/lib.zip
	"libtcl9.0.so":    "c63e002ae89e18019a7abddfa3b694248f244381f01e0785ca640a2a107c95be",
	"libtcl9tk9.0.so": "a9b8188197b997d3f459295decacbe58f283783a019b2c85a585df2761c517f4",
	"libtk9.0.0.zip":  "d9c5ea5d8f80d353dc2bdf2b0231c8ac33cfd136df0a11f5eb4886ff9488d8fc",
}
