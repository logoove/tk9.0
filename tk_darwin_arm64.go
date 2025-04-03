// Copyright 2024 The tk9.0-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tk9_0 // import "modernc.org/tk9.0"

import _ "embed"

const (
	tclBin = "libtcl9.0.dylib"
	tkBin  = "libtcl9tk9.0.dylib"
)

var moreDLLs = []*dllInfo{
	{"libtcl9tkimg2.0.1.dylib", "Tkimg_Init"},
	{"libtcl9jpegtcl9.6.0.dylib", "Jpegtcl_Init"},
	{"libtcl9tkimgjpeg2.0.1.dylib", "Tkimgjpeg_Init"},
	{"libtcl9tkimgbmp2.0.1.dylib", "Tkimgbmp_Init"},
	{"libtcl9tkimgico2.0.1.dylib", "Tkimgico_Init"},
	{"libtcl9tkimgpcx2.0.1.dylib", "Tkimgpcx_Init"},
	{"libtcl9tkimgxpm2.0.1.dylib", "Tkimgxpm_Init"},
	{"libtcl9zlibtcl1.3.1.dylib", "Zlibtcl_Init"},
	{"libtcl9pngtcl1.6.44.dylib", "Pngtcl_Init"},
	{"libtcl9tkimgpng2.0.1.dylib", "Tkimgpng_Init"},
	{"libtcl9tkimgppm2.0.1.dylib", "Tkimgppm_Init"},
	{"libtcl9tkimgtga2.0.1.dylib", "Tkimgtga_Init"},
	{"libtcl9tifftcl4.7.0.dylib", "Tifftcl_Init"},
	{"libtcl9tkimgtiff2.0.1.dylib", "Tkimgtiff_Init"},
	{"libtcl9tkimgxbm2.0.1.dylib", "Tkimgxbm_Init"},
}

//go:embed embed/darwin/arm64/lib.zip
var libZip []byte

// Keep last for internal/shasig.go to update.
var shasig = map[string]string{
	"libtcl9.0.dylib":               "3ff7244441b91539315f7e25a4eb177eefe5d68277df4c1ad9fa946c873564c2",
	"libtcl9jpegtcl9.6.0.dylib":     "66148da05d3bc9360df9f5813f9d5f2f6283f7f6f0d1a31222f752965616bc9c",
	"libtcl9pngtcl1.6.44.dylib":     "3d5019074d6ca62f6d14b96d436415cbce7cb43dadc1b2f35e2f5382d634d573",
	"libtcl9tifftcl4.7.0.dylib":     "f47f223f2c28f318a4d1752db1fd0848c76fbf9b89d210c26834f35c77dca833",
	"libtcl9tk9.0.dylib":            "6acc300a4beb8a83ff80cd790e3f81b2c4c2976b0c15dd6fa16354cee2df7c9c",
	"libtcl9tkimg2.0.1.dylib":       "8899025f1e2ae2c9ce536d3ee6be8d31fed921000c01a62f669b4ab6b12e521d",
	"libtcl9tkimgbmp2.0.1.dylib":    "eb7c982836501e8be77c849edebb3df625d67ff8841c7e7358a7b6cdcb756a78",
	"libtcl9tkimgdted2.0.1.dylib":   "96ad7e52117fce20fc2bbf3abb5e68a57c5652f79375aa28d4a78880fbf1f650",
	"libtcl9tkimgflir2.0.1.dylib":   "711df5120cc0a853ef38d7130a53dcf91265b847317a392c9c4438592a07299a",
	"libtcl9tkimggif2.0.1.dylib":    "6a88c7844b99d1b649ded0f8e6294f816dd3ac113497b4237a42d8ae846fb421",
	"libtcl9tkimgico2.0.1.dylib":    "daf75e7a8d6ce601515ec72aca4e02492ddde4a5c23b55843b9ab0c90c8e96bf",
	"libtcl9tkimgjpeg2.0.1.dylib":   "05616a9999a561dc722dd0ccac81b4738d13052c4a7571749347d72cafe66f2f",
	"libtcl9tkimgpcx2.0.1.dylib":    "f0ba6390a598557c853260f57a572ec109d47ae0a4fc0d9d5978d72ae486945e",
	"libtcl9tkimgpixmap2.0.1.dylib": "a0006cc3cf26d677a9a7b92cb6d6fd3aeba11fd3d0c72ee8e3daf92142e7fe6f",
	"libtcl9tkimgpng2.0.1.dylib":    "62e59b7e1d95cf3e30555b6e99923202f4c8dc8c5ff05bf0ef998edc52f2846f",
	"libtcl9tkimgppm2.0.1.dylib":    "7f3090d82df279b3789f5c85b5a38a87cce8b21783a98b659d6889aa68af39ea",
	"libtcl9tkimgps2.0.1.dylib":     "5f72dcf0b73f4428418ca49549a314305fec4229dd726e34e8a6a0f1a5eda456",
	"libtcl9tkimgraw2.0.1.dylib":    "52fa5796b872382c35689b9bacc6f38410637b79cec9af9ad89882fcdd0294df",
	"libtcl9tkimgsgi2.0.1.dylib":    "fe696f578e06a4dd1f188ff0151e78d799b321a0f76c6e90a8ad214c3db1a0f5",
	"libtcl9tkimgsun2.0.1.dylib":    "f0ed33fa149198ed7408e93c3a575bebc3be4bf495113d3368a6c2bec663526f",
	"libtcl9tkimgtga2.0.1.dylib":    "b30f2212fd23d658f8062e5470335b88f962835ccfaca4b1c8b1ab13975c85ae",
	"libtcl9tkimgtiff2.0.1.dylib":   "5443091049c5f8a7085da4e2c0ec94e1fb8815851db2a303021e69805fa4c62f",
	"libtcl9tkimgwindow2.0.1.dylib": "ead56047915f98c29ff6faa05d2ba7de2c92249c64218682f40a5d5b14f7f820",
	"libtcl9tkimgxbm2.0.1.dylib":    "8a3d5d36f2dcf4c167a8475c69ef53e4cda876fa928f81c5eca91f91cd4524fe",
	"libtcl9tkimgxpm2.0.1.dylib":    "e108833cd6963d0089d570d2a8c743cdbd315303c704155a3c3d968ccd2e10b5",
	"libtcl9zlibtcl1.3.1.dylib":     "6b8632785e5e8df0a768e61aa59ffdb6edbb858120b5b2d37e7c5e0756550095",
	"libtk9.0.1.zip":                "6e8e50cab4500600a2cae65938238806e10bb7a123104d80a127d115b247fae9",
}
