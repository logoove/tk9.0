// Copyright 2024 The tk9.0-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tk9_0 // import "modernc.org/tk9.0"

import _ "embed"

const (
	tclBin = "tcl90.dll"
	tkBin  = "tcl9tk90.dll"
)

var moreDLLs = []*dllInfo{
	{"tcl9tkimg201.dll", "Tkimg_Init"},
	{"tcl9jpegtcl960.dll", "Jpegtcl_Init"},
	{"tcl9tkimgjpeg201.dll", "Tkimgjpeg_Init"},
	{"tcl9tkimgbmp201.dll", "Tkimgbmp_Init"},
	{"tcl9tkimgico201.dll", "Tkimgico_Init"},
	{"tcl9tkimgpcx201.dll", "Tkimgpcx_Init"},
	{"tcl9tkimgxpm201.dll", "Tkimgxpm_Init"},
	{"tcl9zlibtcl131.dll", "Zlibtcl_Init"},
	{"tcl9pngtcl1644.dll", "Pngtcl_Init"},
	{"tcl9tkimgpng201.dll", "Tkimgpng_Init"},
	{"tcl9tkimgppm201.dll", "Tkimgppm_Init"},
	{"tcl9tkimgtga201.dll", "Tkimgtga_Init"},
	{"tcl9tifftcl470.dll", "Tifftcl_Init"},
	{"tcl9tkimgtiff201.dll", "Tkimgtiff_Init"},
	{"tcl9tkimgxbm201.dll", "Tkimgxbm_Init"},
}

//go:embed embed/windows/arm64/lib.zip
var libZip []byte

// Keep last for internal/shasig.go to update.
var shasig = map[string]string{
	"libtk9.0.1.zip":         "623a92281779f6e5d5ba3048829b731c31e1a03c1a5aa385462a11e589f053d2",
	"libtommath.dll":         "907b9c7860fc07231f1e238551715e5d813283807f52dc383dae0cb47a879d29",
	"tcl90.dll":              "946550f143038d61ea20a64087c1f2b6a5c533bb365468fae08d9092c71960e6",
	"tcl9dde14.dll":          "7bbd008b09e0b1e74e5df69152686edf0b79893301986a8a794a41d1fb862504",
	"tcl9jpegtcl960.dll":     "340199f135062d5238395129acd19d8a2dcbc9a6b1e0a411b3fb187ffb8df696",
	"tcl9pngtcl1644.dll":     "d9a49f2839587b9256d1432d5eaf92ea51fc6f14e9ee6e84df7a786363f4f48e",
	"tcl9registry13.dll":     "63910b8dd3bb088e50d447d55e8bc046ca1f5743423e61e86cbf80bbdc68329a",
	"tcl9tifftcl470.dll":     "7ecd680308305bbf464e6dc98ca401f8cc237aaeee6e3d4b86b242dba4bf4ded",
	"tcl9tk90.dll":           "1471e1cf40833072b6d337fb6d1c3ecb70c14b92d75b2290c1b5fed5985e68bf",
	"tcl9tkimg201.dll":       "e12ee2f9b59dbcf61277519542b6be089bc59f2022be2f11233ac9a69d2648cb",
	"tcl9tkimgbmp201.dll":    "173ba3302df5897e704086273512800eafc69527309c0d3b613fe714b96060fe",
	"tcl9tkimgdted201.dll":   "9d2d9b894369e71cf70318eca059ba9196125904e68b649017b53fcc941e1feb",
	"tcl9tkimgflir201.dll":   "cf077f9f23ad0974eb85bdf465b70f972b7188c59cc157b5ffb4e88816eddde8",
	"tcl9tkimggif201.dll":    "0daf21816a66679db35fcd9a3227944a8e903fecab1366933c33df610c5e85d7",
	"tcl9tkimgico201.dll":    "a2d1bbc413f5187631d43e5f162a515d654539f623be93623080b2df8d635ea4",
	"tcl9tkimgjpeg201.dll":   "43a28dd8e07001dd52acc62db13d8f822ab7711a0d86a1b10f3cb6f371d64fd4",
	"tcl9tkimgpcx201.dll":    "e58eaaeb88ed22894ea885fc0028f6407cdc404745400715c9e3aa86a4acdce6",
	"tcl9tkimgpixmap201.dll": "d49867ce7d2eb9ab1dcaa35e6822e378718bb9b80e789ba2727306f9f3e53db9",
	"tcl9tkimgpng201.dll":    "06b9936ec0bb5c2597811259b86cdcdd46cf17bf863531c843e2942f3418f1ea",
	"tcl9tkimgppm201.dll":    "9c593eba5b880197a8290c095d1427da477bf8f1d6e61811773d468e810d012b",
	"tcl9tkimgps201.dll":     "efd74739195ec4718d5d168613c6c8eca652f4923fdf91f03b84e926ecb92919",
	"tcl9tkimgraw201.dll":    "7bd610472b99681a27bd2c0ed734fbd64409b2999814368d56ef8a20838ece76",
	"tcl9tkimgsgi201.dll":    "2c9bc0554d9966deaa63c00373a8aaafe2bd3d25ec9701b3654efa3bcddec8c6",
	"tcl9tkimgsun201.dll":    "bde8a35106891c403baf7866d4142b57be47bcf8cac715c0994d3ebe44887ec2",
	"tcl9tkimgtga201.dll":    "945ce75c8fa7be85a4e62ec7cbe9ead83569b1ed9109d7b957ea8560f49e5848",
	"tcl9tkimgtiff201.dll":   "278b5cfe96633075c21035fcec639e7c95431cb962d072755bd9fb70fd33bf95",
	"tcl9tkimgwindow201.dll": "a7ead48068338f8a00dbaf4fd62d7ac4561b2acbbb753290e592ed4aa3582cc5",
	"tcl9tkimgxbm201.dll":    "f1f29fae705575d49401963789f487a6e9f0419ec0c50a5be8d01d184241b390",
	"tcl9tkimgxpm201.dll":    "26fcc4b1e24ee3aaefb1195975708b5e95a10c6146fcdebf5ce184ece26c7a21",
	"tcl9zlibtcl131.dll":     "415d7e671d45af39cd7162e819a5ddc8996540853083cdd04bf622e2e8814505",
	"tcldde14.dll":           "9c204cf4bb1d845bbf1f238bfaa1ce8128b2dc5f3ecd8774f6edecb2ea16cb99",
	"tclregistry13.dll":      "107d21531a752f7f4504dab3ddbc6b5e4d1980701a71d2975ed722312adfd2c1",
	"zlib1.dll":              "6f10a76dcc2c831d1f08d98c0b345afa0911bec0238fcba357b612ccc6ab5d81",
}
