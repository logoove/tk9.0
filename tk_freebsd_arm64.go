// Copyright 2024 The tk9.0-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tk9_0 // import "modernc.org/tk9.0"

import _ "embed"

const (
	tclBin = "libtcl9.0.so"
	tkBin  = "libtcl9tk9.0.so"
)

//go:embed embed/freebsd/arm64/lib.zip
var libZip []byte

var shasig = map[string]string{
	// embed/freebsd/arm64/lib.zip
	"libtcl9.0.so":    "53641e10ee9405c4951edb90fb6d1efa4f9f98eb1d6da711597e5d6c4a8e1d4a",
	"libtcl9tk9.0.so": "1db9adbd1463368f6992f86fd27d94ad88013ed5e4035ce30dd651bed1b2f6b8",
	"libtk9.0.0.zip":  "7ea3ca76b18051b63aeba1b0ddf1bec037891ebb5be100215384742d3af1641b",
}
