// Copyright 2024 The tk9.0-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tk9_0 // import "modernc.org/tk9.0"

import _ "embed"

const (
	tclBin = "libtcl9.0.so"
	tkBin  = "libtcl9tk9.0.so"
)

//go:embed embed/linux/amd64/lib.zip
var libZip []byte

var shasig = map[string]string{
	// embed/linux/amd64/lib.zip
	"libtcl9.0.so":    "e693dc74ea2039474ce42dc53eb2774e883b0903e7dbf3dad3d9040c665dcd1a",
	"libtcl9tk9.0.so": "e55f7910c1477dfc259bc7cace99a8184fd3c091e9b22377b02da11ca9f896e4",
	"libtk9.0.0.zip":  "c666cf661852f896d83669a34fe5eb12620d0aae1517841082eb38093563b73c",
}
