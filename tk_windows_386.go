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

//go:embed embed/windows/386/lib.zip
var libZip []byte

// Keep last for internal/shasig.go to update.
var shasig = map[string]string{
	"libtk9.0.1.zip":         "3b7dfe192a84461f83684dd758946082a43ed915b557c0d00b4d76ec29af6d64",
	"libtommath.dll":         "7ff97843cde97215fcf4f087d61044cda01286630b486398117967e577e039e3",
	"tcl90.dll":              "74f7b7d745aedca8465cb7eca6a5b1ca725755bc744cb088fdbcbc8357074fcc",
	"tcl9dde14.dll":          "6cb4970ed71ab865b19953e098c89be3ee74253f25ede18a95151d8af90b3f71",
	"tcl9jpegtcl960.dll":     "54daa13553ed1356ee8a33405209c8b7f17c71ed924dbaac2b0df8a354336269",
	"tcl9pngtcl1644.dll":     "3822977fd3d318c790ef5bd586f24ead8d025b0c01e9cd03c2dd692cd9876d5a",
	"tcl9registry13.dll":     "3311be4270ecab7b8fb969289d0f27c79d7cfd4c29cf4cbc05aeebf5de398687",
	"tcl9tifftcl470.dll":     "5ae45cfe013709ad4856fcee4f0143c0aecfcaee1708c8e3b1d409e57c37fd71",
	"tcl9tk90.dll":           "0bd5ca4c43355e48c23aa68dab9471bc6af717b98fb4586135d2e077a31976f7",
	"tcl9tkimg201.dll":       "cb67d432508f11dbe810f27bf2175838d302d6b198020e7ed1a89d5e9a88af86",
	"tcl9tkimgbmp201.dll":    "20e5e58bcee6c60d9d5388c1542afc400faa2dcbfdc2ec7fe3cd5028a0e57205",
	"tcl9tkimgdted201.dll":   "9a0e1641e8b7381ce7e0564701d151089f04307344dab5fd23177bf95e101783",
	"tcl9tkimgflir201.dll":   "ea9c41b9bb57527907ef328d60e1bebd1d7d39dbe3e607fe83b9783e4fa631d9",
	"tcl9tkimggif201.dll":    "7b5f6381d65c8ff30d244d8f1657c94457437cb701f5bfd1d055efe91f11a320",
	"tcl9tkimgico201.dll":    "8f48c20ac7e247e816374e0d4c7daccd4461f76419e670079daf6c6796715540",
	"tcl9tkimgjpeg201.dll":   "4c8daece894a1e3bed2067f9762537ef376293a0b7e13165f2b8a83b99bd0c80",
	"tcl9tkimgpcx201.dll":    "3588f36eda5ea07fb9228e03c5490622f8620e99d53249cc26948517d3cef115",
	"tcl9tkimgpixmap201.dll": "81640481d77ff1eb75e70fe63e9843920fcda54977622b808f8f9a8fc5285c00",
	"tcl9tkimgpng201.dll":    "0e232bea1c4afb8b07e899c99f638d6e33a270914c37586553738338ae97c3e2",
	"tcl9tkimgppm201.dll":    "6b54f027f84e8de0a0ec60411ccf920e47b2783459201b8f5e322a0742769ad4",
	"tcl9tkimgps201.dll":     "3817052af0b42ee727f72393cadc875f29f7d7112fcd3d475278e3a64d71e035",
	"tcl9tkimgraw201.dll":    "e55aef30e1e95b5c0c7ea9421fe379885830000744a36c51f29d14f66545a374",
	"tcl9tkimgsgi201.dll":    "20d3e3b78f3dd601ae3be4434c34579cc73e73ec75e218e47288c57adf0a6210",
	"tcl9tkimgsun201.dll":    "972677f91776c00847e02395a42cc348ece9894285bd0f968906bba7b6901049",
	"tcl9tkimgtga201.dll":    "e87b4fedcbcda6d9bb2d033991f5c9fd7964585266f2eee3f76c9923d8614a45",
	"tcl9tkimgtiff201.dll":   "d4a60e982ab45438f2f89d2ac6b7765aa49713f46e0b3db2cc7788713c40de4a",
	"tcl9tkimgwindow201.dll": "1dbc2bdadddee11453d72efa9400752f9aa05cfa8f4e6cfe5f4f60fde01b50f7",
	"tcl9tkimgxbm201.dll":    "ebc86d74e026931a01ac68ac853963a0b8b359733c2f59cfaa3c2a3b46f32292",
	"tcl9tkimgxpm201.dll":    "25cacba9940acc23a274f7d7af2450d582a6fd57e11b1f5be9b4fd5734eb16aa",
	"tcl9zlibtcl131.dll":     "d7631eb5a6f962fa5b05764246dfedc601715965ad20a3bf45394f3499290e65",
	"tcldde14.dll":           "c7d0d266b3c94c0d57904713fb38b6172991457092ea100bb34327f0889c1d9b",
	"tclregistry13.dll":      "4fe68da94eb9961a3d632b67b1d2dda7caaf933f683eef2523c58ef70c1baf2c",
	"zlib1.dll":              "60f637680d84a0717cbee4cbf219b6215ca1f21fe0b32c8de2819c328c72ef15",
}
