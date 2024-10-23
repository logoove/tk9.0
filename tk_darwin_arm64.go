// Copyright 2024 The tk9.0-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tk9_0 // import "modernc.org/tk9.0"

import _ "embed"

const (
	tclBin = "libtcl9.0.dylib"
	tkBin  = "libtcl9tk9.0.dylib"
)

//go:embed embed/darwin/arm64/lib.zip
var libZip []byte

var shasig = map[string]string{
	// embed/darwin/arm64/lib.zip
	"libtcl9.0.dylib":    "707bb2ec114901c5821ed7c6dfb4bfc9b351c5d45398a9704b624586231a250f",
	"libtcl9tk9.0.dylib": "face40c56012e61bfff95970de07ccda03a31b60e2ed4f371f93c19d4259b49b",
	"libtk9.0.0.zip":     "9df27e04a9ca76b080520de045200a7bb1e9a62f738828765da40784a4221252",
}
