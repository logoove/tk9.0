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

//go:embed embed/windows/amd64/lib.zip
var libZip []byte

// Keep last for internal/shasig.go to update.
var shasig = map[string]string{
	"libtk9.0.1.zip":         "acfeb6e324ccdc0cd1d25d45441ca204ab601f212ee55a7add9c558f9c2ef716",
	"libtommath.dll":         "2d760fefb452665b6af8c8d9d29f3a8378f10fc0847cdd9938ea0cb5edf1d573",
	"tcl90.dll":              "8bd6cea53df4cf7982cf4b600b9304aa2f4a7865e7708d9927f6fd27777e0af9",
	"tcl9dde14.dll":          "6d77cc449c711e0f7247f4f6685943d10b3bee63e6c48bb67df4bda744edfa6f",
	"tcl9jpegtcl960.dll":     "7a76e25128c3f58d90288c03a9487a87480e7311c03bf9d86925d3b686c2a339",
	"tcl9pngtcl1644.dll":     "6d5680f577139e703d788eb4e4cae20d8e3328c70e14a5afb08617f882c6d1af",
	"tcl9registry13.dll":     "cea564632500e1d9ea059f0e2888f81db607833445b8f2c3298ab343c8323950",
	"tcl9tifftcl470.dll":     "d38b3b24927461996dfffdc370c16bf7b378c440a686d3978c74fb8930252810",
	"tcl9tk90.dll":           "811a5a9d58ffbfb9758075416e4723b05b238df2347de96776cd74da45844ac9",
	"tcl9tkimg201.dll":       "4d7f3e3a8c14d30e1580958706977119c73c93637cf326b57172d93ee7a61968",
	"tcl9tkimgbmp201.dll":    "471ca49fa65c37e9ad45cf220d37e34e15c960074dbb97e08a01eb09ec0d8d54",
	"tcl9tkimgdted201.dll":   "9d965fd12a6ca857672587fb063e42aa3811d555e10fd20cbb386f397c61b3ec",
	"tcl9tkimgflir201.dll":   "933a0bd03acb60793ecbd678f3142ce7ecf75de3bf811a6eafa635f9ffa6e276",
	"tcl9tkimggif201.dll":    "b5210227cbd3f6c5517713083f882a4c7f76ea7f492b8414c7d6fdb6bbbe4393",
	"tcl9tkimgico201.dll":    "8a16f355cdadc5e9278035308a0074b45fee9f761a3b5a741347525a5129e0df",
	"tcl9tkimgjpeg201.dll":   "282ed8c7e9768c9419e9ee21e5617644393a3ab89a1bcb2c1dab7b0b4eb34772",
	"tcl9tkimgpcx201.dll":    "20e2d759ac3e21744ef5cf3720bacb170b4e0ef30d50120e9d0f48aaf36ee12c",
	"tcl9tkimgpixmap201.dll": "7fd81261deead568411cf5f04ed130590d5d03d1c6d264a39a130c67edb237fa",
	"tcl9tkimgpng201.dll":    "ebfea4fede8c9329f4ed8eff60a572b23c782be9885c8d906c25ad0b04773640",
	"tcl9tkimgppm201.dll":    "564600528e65bd891f30a441cb7de5124e961b6748abaf3149456583bdf5636d",
	"tcl9tkimgps201.dll":     "d6cb58429e000269d678926c57542f90804d04270ad1593cb87dda62d98ff678",
	"tcl9tkimgraw201.dll":    "f27444a15f9aa9e43de856fcd97d545c607a8f8bd6d9f48053ccf94198e3267e",
	"tcl9tkimgsgi201.dll":    "9cafb67123a4b0d831c38d8f21360a7b5dfd52f739302689df5036b6ae8a1fa1",
	"tcl9tkimgsun201.dll":    "a5d7e0f0a5af313d4d0e3b56ed4909f2ab6725ed012e9eb7f3e101660511a031",
	"tcl9tkimgtga201.dll":    "f319a233d27b34268e89151eae2258e95b992983b2dffe33a42c22d9c3cbdbff",
	"tcl9tkimgtiff201.dll":   "b285caac087ea68aa6fb7f3cc16953a6ca281d2ee50d100e1b986bd545a3a93c",
	"tcl9tkimgwindow201.dll": "37d723aea212c831e33665c59f991a9c50cf2dfb50fddd6f715a70e6ef20eb9f",
	"tcl9tkimgxbm201.dll":    "087e0d363e81f9e1f4abc1b68f90977499fc2461f4f8daafdc83cbbfaaf0e388",
	"tcl9tkimgxpm201.dll":    "8e1c889f9a66603f73bd309597814a61ac2448e9e9097ce6ae72df32ae8175cf",
	"tcl9zlibtcl131.dll":     "7affcaf3f8ec4cf67b048e85b97a9110711b584c0f63875a7abdb5db0f3bafa1",
	"tcldde14.dll":           "99d1b099e2d941d84fbccda8809f9666529fd469b4ad7d8091c6911c0ca1280b",
	"tclregistry13.dll":      "088ab64e79ed1f6aaebe64f066c5f1bed72ee31916799efb75ed29b1f55e3128",
	"zlib1.dll":              "04117778e255ed158cf6a4a1e51aa40f49124d9035208218fbfebbe565cf254d",
}
