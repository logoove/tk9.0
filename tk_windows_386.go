// Copyright 2024 The tk9.0-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tk9_0 // import "modernc.org/tk9.0"

import _ "embed"

const (
	tclBin = "tcl90.dll"
	tkBin  = "tcl9tk90.dll"
)

//go:embed embed/windows/386/lib.zip
var libZip []byte

var shasig = map[string]string{
	// embed/windows/386/lib.zip
	"libtommath.dll":     "7ff97843cde97215fcf4f087d61044cda01286630b486398117967e577e039e3",
	"tcl90.dll":          "9f2493060f307a27e117998160b19ad1b6ddae999ce71eea431f8299f223277d",
	"tcl9dde14.dll":      "6c218c19c18ff297185386d7e361ea6ce4d56ce22eb5d2551ea767b61e5c8fe1",
	"tcl9registry13.dll": "ac04ac882f9dba038974ff1ebb8b26c5384c8e8b84d62caf7a8872dd5a25e535",
	"tcl9tk90.dll":       "00b1ce0f72a297d28fdd6bb24507937c14105304eb36ea85a739b52c0cbcc7f6",
	"tcldde14.dll":       "0d31eb9d3c21c36514c7ff5e53ab05d3851f2bd6a6feddf9b0390c76452cb2c1",
	"tclregistry13.dll":  "4441756a0805ec11db9f334a574fe6fb4593a63d204a632e0c4af0591264abfa",
	"zlib1.dll":          "60f637680d84a0717cbee4cbf219b6215ca1f21fe0b32c8de2819c328c72ef15",
	"libtk9.0.0.zip":     "515d0f0a22c1faa7e59389601d903571b5437cd98bdf6192d794a4d5dd2acf8e",
}
