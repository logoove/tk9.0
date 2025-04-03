// Copyright 2024 The tk9.0-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tk9_0 // import "modernc.org/tk9.0"

import _ "embed"

const (
	tclBin = "libtcl9.0.so"
	tkBin  = "libtcl9tk9.0.so"
)

var moreDLLs = []*dllInfo{
	{"libtcl9tkimg2.0.1.so", "Tkimg_Init"},
	{"libtcl9jpegtcl9.6.0.so", "Jpegtcl_Init"},
	{"libtcl9tkimgjpeg2.0.1.so", "Tkimgjpeg_Init"},
	{"libtcl9tkimgbmp2.0.1.so", "Tkimgbmp_Init"},
	{"libtcl9tkimgico2.0.1.so", "Tkimgico_Init"},
	{"libtcl9tkimgpcx2.0.1.so", "Tkimgpcx_Init"},
	{"libtcl9tkimgxpm2.0.1.so", "Tkimgxpm_Init"},
	{"libtcl9zlibtcl1.3.1.so", "Zlibtcl_Init"},
	{"libtcl9pngtcl1.6.44.so", "Pngtcl_Init"},
	{"libtcl9tkimgpng2.0.1.so", "Tkimgpng_Init"},
	{"libtcl9tkimgppm2.0.1.so", "Tkimgppm_Init"},
	{"libtcl9tkimgtga2.0.1.so", "Tkimgtga_Init"},
	{"libtcl9tifftcl4.7.0.so", "Tifftcl_Init"},
	{"libtcl9tkimgtiff2.0.1.so", "Tkimgtiff_Init"},
	{"libtcl9tkimgxbm2.0.1.so", "Tkimgxbm_Init"},
}

//go:embed embed/freebsd/arm64/lib.zip
var libZip []byte

// Keep last for internal/shasig.go to update.
var shasig = map[string]string{
	"libtcl9.0.so":               "39e6cd7d3eaa959ba0e1b712419d809df9f6088b7b0d6dae258c923d417ed30f",
	"libtcl9jpegtcl9.6.0.so":     "f3cbdfaa0051d40ee4635e4e7a8cd0d1a68e2fb8b63150cf009ef64b4e764c12",
	"libtcl9pngtcl1.6.44.so":     "a39081a5ecdf8ec00f2ed08b40ab8c54265cad41b3a6d9a76a0d12cae4efa409",
	"libtcl9tifftcl4.7.0.so":     "8a7b79b7fb60462b3fceeae408d1ef4093f91ba623e17b4bbacaf46d41eefe2b",
	"libtcl9tk9.0.so":            "d77d80549d567915333312946b5cacf5582614b7b78aa2b8cc78bdd75cbc86dc",
	"libtcl9tkimg2.0.1.so":       "506f5454924b82f9b44077d579f0ebbe48425686dbd8c1847f48dfd8ad260ef8",
	"libtcl9tkimgbmp2.0.1.so":    "2c1912d6eb7391cfb58ac6b8787b9331b367a9a5a7b32764005d3a3d0aaaa810",
	"libtcl9tkimgdted2.0.1.so":   "3c414db98fe83afd269267de56f3dbc48119b3aa1570bbcc32a441b07afd7cfc",
	"libtcl9tkimgflir2.0.1.so":   "4902385ffb0e8cf71f50694a646e97c59fde9ffff41d5a466f405170e5c12151",
	"libtcl9tkimggif2.0.1.so":    "7501129c9407d897cbe686b4d4696fd2fcd4484abc6b12fc8d8113b871c366d2",
	"libtcl9tkimgico2.0.1.so":    "39ef080558937c42c6154a36a8a9f60599a2ad405943afa9397849008245425a",
	"libtcl9tkimgjpeg2.0.1.so":   "e5d70999ee86dfe52d4569626b3e7f476d9febcbc617550c554fb76c2faebb59",
	"libtcl9tkimgpcx2.0.1.so":    "a91ba2fb430f55c2ca841568f64549568a316439d993118b0342eb1efad8b471",
	"libtcl9tkimgpixmap2.0.1.so": "213d00c838190914220622525aec23101d5e0ce9acb88a48728a12784bf398a5",
	"libtcl9tkimgpng2.0.1.so":    "9553fb0f0468844737dad10685f864a95d1d379c82f7154779d96404ce3a7088",
	"libtcl9tkimgppm2.0.1.so":    "d2b006b7d1b2ca2add8287b05843f78a53f2c3579f46b86093d970410f725f2c",
	"libtcl9tkimgps2.0.1.so":     "d37d5612589ad69df6ee477d0470a854731731ef7a722958aa1897a43efd2197",
	"libtcl9tkimgraw2.0.1.so":    "3c80162b6fed698419e3424a4781b4b13610f1bf1d39bbe8f11811c592d99613",
	"libtcl9tkimgsgi2.0.1.so":    "1a6fa89e30ec2aa6052d4e4aba24116bc5c79163a52c8ae0c6568cb47040f0ce",
	"libtcl9tkimgsun2.0.1.so":    "b91f2770d4f5b8ad5a8b189c38c22962527fe0c6107534f202907db584bc7925",
	"libtcl9tkimgtga2.0.1.so":    "83c17d90bc62fb006108c6d0e8acef91a30e36633329c10e29f8a16319cb5643",
	"libtcl9tkimgtiff2.0.1.so":   "075681e206044897f838830b2f15ef765105727bf9e5c417836f6dd8a4ba23d1",
	"libtcl9tkimgwindow2.0.1.so": "1a3ff82f1b654ab9ad1fcf20b9997e216de3d5c4410c8dfc444fba5dd2d40ab6",
	"libtcl9tkimgxbm2.0.1.so":    "9b2ab09addc990dd3b285a0a41bf3dba5050bc7973ddf1afef055a0e3e26e757",
	"libtcl9tkimgxpm2.0.1.so":    "eaba13cd9ef1b94e3e916836bf9021d10eabeced867198977bbe12adbd54fe58",
	"libtcl9zlibtcl1.3.1.so":     "d1d23f583618092b9ae3a2d57176c6fa628ce25f1f6a3aae1fc49f56cafdfc51",
	"libtk9.0.1.zip":             "3559e57962125b8d1b2ec1c72259eb12e59b0b9ef70168c12046b00a52d4f026",
}
