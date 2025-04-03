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

//go:embed embed/freebsd/amd64/lib.zip
var libZip []byte

// Keep last for internal/shasig.go to update.
var shasig = map[string]string{
	"libtcl9.0.so":               "19cb1117f5cc920c0971ce46137a9221836d65f152103dad179e2f83cda52039",
	"libtcl9jpegtcl9.6.0.so":     "b3fda78806e463704cd60ae4f3ec5e4e29ce3fe43c97b9847bec16afc9b70f9f",
	"libtcl9pngtcl1.6.44.so":     "cc4e66ac1a4d649455d218e545346011fed66f23f565ea08713ff35081d184ab",
	"libtcl9tifftcl4.7.0.so":     "da044ffd8998cdc7a9fc3c9b8a7662159e9c1412585f22e5e3f7b3bbd03ff5c3",
	"libtcl9tk9.0.so":            "d212b0f1bbf732d3e725cf8af59882988a4aadd2f26ef2e94240dccf7e1e1f4e",
	"libtcl9tkimg2.0.1.so":       "b8165577d31e8ed02dcb5187bca4604d1ad5dcc55144d27939dfa293e3d4aadd",
	"libtcl9tkimgbmp2.0.1.so":    "852fbb1d6b88fed45b83d262588cb7d7c24937cfa661a8c382a2b4bad564b024",
	"libtcl9tkimgdted2.0.1.so":   "91ac47f1c4f042a06928804a4e78ca19881fa7e8737be624f3a6057506b3b010",
	"libtcl9tkimgflir2.0.1.so":   "0ebc7062433945e1f33726ff05d07bfd8049f8b8b4f8e76461aeee298dbe31bc",
	"libtcl9tkimggif2.0.1.so":    "d16dd0b7e6f0207cd12b27cdb63d9851e39b3c97e0be9ae339d014faeb28f423",
	"libtcl9tkimgico2.0.1.so":    "5d3b45c427e0ba43939effc1533b8f995776b776de0d93346f39523b029ce0e7",
	"libtcl9tkimgjpeg2.0.1.so":   "4d3e21793b53852b3daef89ce0e0968dd2c1a00f1d255dd6ba63fd98dc2fbb0f",
	"libtcl9tkimgpcx2.0.1.so":    "1dd19709cf044fba2d5fd85a089c3095fa0998fce2f00f12338e6433aa488473",
	"libtcl9tkimgpixmap2.0.1.so": "0fbfd8c1799c95a86423d27857187a99813c82049ac97a025215b5b8b7b9c2b4",
	"libtcl9tkimgpng2.0.1.so":    "dabceeea6e70d9854f6e4f98d24339008495d75d94f257ebe33aec2e2fddbf75",
	"libtcl9tkimgppm2.0.1.so":    "39858f9201dd08de1935701875d8f526605d93add6df1dabdc563ecaebb4ba8e",
	"libtcl9tkimgps2.0.1.so":     "9a63964ed2fa4aa385ed1462a64fcd83721c9b94be9b2ef05f6798d665719e74",
	"libtcl9tkimgraw2.0.1.so":    "c050e178390105ddc90a81b1f1ed44daec6b6a12e6bede710a6fa1edac322c39",
	"libtcl9tkimgsgi2.0.1.so":    "110ed2e94f7a765bd02b3077246ce29530830dba3e733570150c380679148594",
	"libtcl9tkimgsun2.0.1.so":    "47a0a63fd5e9355445c8eac746a559ca6cd209c9cb13cb1c60fa07fbeb2077b7",
	"libtcl9tkimgtga2.0.1.so":    "7c0e85c5dae29e9ca50b72c5791de0048965eb5f5412a823175c75c96a585b64",
	"libtcl9tkimgtiff2.0.1.so":   "d38e894aadfc76aa5b6a83309da1632fee92b5c2a1da97ff69c7edef20551d9e",
	"libtcl9tkimgwindow2.0.1.so": "60725aa1c5b42b6fe95328689c810bb6d915a00fe046363e944533817671f4f6",
	"libtcl9tkimgxbm2.0.1.so":    "ba0e16acb51f23edfc0c5a2a1b1aa07fe5236770b21c8280d3a31b347fc87142",
	"libtcl9tkimgxpm2.0.1.so":    "7eee21cbe868874a94bb724408ee80d28478b6f6b0b89ee2e08e6d64f15d8c37",
	"libtcl9zlibtcl1.3.1.so":     "f76dd378ecc042cc9b39dce48eb824e0dbb4d425f5e5c4a168fefe68b2a30fef",
	"libtk9.0.1.zip":             "15fcdbd88286ac5c1f32e92f795111305975ff5582f64d5257c66ab8f6ba71b0",
}
