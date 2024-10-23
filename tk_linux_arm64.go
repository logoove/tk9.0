// Copyright 2024 The tk9.0-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tk9_0 // import "modernc.org/tk9.0"

import _ "embed"

const (
	tclBin = "libtcl9.0.so"
	tkBin  = "libtcl9tk9.0.so"
)

//go:embed embed/linux/arm64/lib.zip
var libZip []byte

var shasig = map[string]string{
	// embed/linux/arm64/lib.zip
	"libtcl9.0.so":    "02a04c7126deff142efe0176990909c538dc750ab6c64849e5215f264afcb558",
	"libtcl9tk9.0.so": "3d48a698cfd884064a01ff8bf606610407280b8692d87d4c50ee4731979e6746",
	"libtk9.0.0.zip":  "327e0672fe208d18b5ca09d5008ef1181db977175dda8b1e5084a32990c70452",
}
