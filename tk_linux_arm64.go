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

//go:embed embed/linux/arm64/lib.zip
var libZip []byte

// Keep last for internal/shasig.go to update.
var shasig = map[string]string{
	"libtcl9.0.so":               "5e3bc6ef81a9d50bda00f4963ba2dcbd32931d3906fb2e34ff1f759bb23dddca",
	"libtcl9jpegtcl9.6.0.so":     "db007e358a449e14db1d07ee6c63ae8db005af96ee35660ef5e64ad55808770d",
	"libtcl9pngtcl1.6.44.so":     "7a76119cd3ce5068dbe56320fcb24eddc1808f0f3b4c35be4b2ff78a1faf8e2b",
	"libtcl9tifftcl4.7.0.so":     "538e34f7a747fb97439d977bfe06ab77176b1f0374a85d9b764b768dcb83787f",
	"libtcl9tk9.0.so":            "9c6d686ddff15af8cc9ec86508f07262d9ecc208c64919353f8c017a88958313",
	"libtcl9tkimg2.0.1.so":       "4091326f80467d7fade8b9194b40dce6d6b85c37be719a63beaa8ddcc1344d33",
	"libtcl9tkimgbmp2.0.1.so":    "5c879d1ae14b6e26051f49f1f412a79e25da0110bafe63cdd0c2f13cf764d6f1",
	"libtcl9tkimgdted2.0.1.so":   "66951bbee812d68d67281a2a98f03d00fb33b169516a8c8d1ccc4523161cb86e",
	"libtcl9tkimgflir2.0.1.so":   "244b640ed2ae26d0b95c53a6ce3b6e67a404aee318449c9eb8529f6dc509a775",
	"libtcl9tkimggif2.0.1.so":    "3fcf81d45371b8fdaf488aff67c47efaf8e80da259c881088195330a5d942869",
	"libtcl9tkimgico2.0.1.so":    "73f0512c06ff02b9893f18fdf3607d6e80a154dd79c76f6a095df13112a57000",
	"libtcl9tkimgjpeg2.0.1.so":   "153f383f34ce61bd20ad509b844bbd4c7a008dc6ef3ada60acc4082d3e5ec35f",
	"libtcl9tkimgpcx2.0.1.so":    "9a07985f97bc6f46bf03e7eef6fbab8663aac84f21198e35ac327cc7947b4e28",
	"libtcl9tkimgpixmap2.0.1.so": "acb7dc6f3dbcba863271a117a72e428f41312868f7726d8fc39d1976819e3af7",
	"libtcl9tkimgpng2.0.1.so":    "de81a3a583758081c3081e973123a73410ad28100a66edbc9bf17d77038f5350",
	"libtcl9tkimgppm2.0.1.so":    "c53df11e4185ba03c45c0ee01b8e0406b034301e1c756fdd285edee32493c32e",
	"libtcl9tkimgps2.0.1.so":     "5ad4d4ca56763b79093e344056e56b1e0e969f77f4f6ab78ebd6b33fb7ffd32d",
	"libtcl9tkimgraw2.0.1.so":    "deb56e2bf69b3bd7f9a9dc4e8ab60f0bb1d860df0eeec09c61ebea9f710d6038",
	"libtcl9tkimgsgi2.0.1.so":    "01b944bac0f1a237358392fd1b133c334a031bb1892092a698f50c4a6a0a2680",
	"libtcl9tkimgsun2.0.1.so":    "660d1b0eb45a3f30ecd8a1f38329bc5240cf6f7d34cdf5cd155351277f26f37c",
	"libtcl9tkimgtga2.0.1.so":    "0e0b1338f9ba8208afe516e8445d4d3d1d1a98d6e882e096559380769cc9abef",
	"libtcl9tkimgtiff2.0.1.so":   "35b17f74ce097779204b67a70f32b7cc61cecbbfd098e0bc54d25314ff7a752f",
	"libtcl9tkimgwindow2.0.1.so": "a7bb644972a7c608e258af2b18451db3cbb2f9088ccebec153fd3beaaef50a6e",
	"libtcl9tkimgxbm2.0.1.so":    "820e609a143b3ef106fd6e07fbdbf7ef4918a3b5db20f40d1885ac2156fa45d7",
	"libtcl9tkimgxpm2.0.1.so":    "fc0637d04313316a1a44d186dc23f84812344c3aef339fbe6562605026d722dc",
	"libtcl9zlibtcl1.3.1.so":     "8024fd432408a31204cd3788e9768de2ef1e27aa800f3c46b0b9af9960a76cb6",
	"libtk9.0.1.zip":             "b173dcbdc2529e302159e557b44b85fd71e01c66409b4db6a04b507a1ba0682b",
}
