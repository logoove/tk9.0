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

//go:embed embed/darwin/amd64/lib.zip
var libZip []byte

// Keep last for internal/shasig.go to update.
var shasig = map[string]string{
	"libtcl9.0.dylib":               "e8b73144f68acb79f7d57e78e86612229cdd87a4622f2bc66cb90004e05ea259",
	"libtcl9jpegtcl9.6.0.dylib":     "5abc4a63ad1513963f51597cc4e8d6fcba942576183b895bbb53ee7fc0a8c6ac",
	"libtcl9pngtcl1.6.44.dylib":     "72e4c3d063e1a1b9594ef80a627035a974a3f8df010216540f6d7df51e8c8403",
	"libtcl9tifftcl4.7.0.dylib":     "7cf58c15c4acd58cc694ca6d739ead08da664d5c5d810af497e86eae43cc109d",
	"libtcl9tk9.0.dylib":            "a3bf659cf22bed7d015747e343033b4061fabb9f338d568600f89c649ebba1e9",
	"libtcl9tkimg2.0.1.dylib":       "fddd4d4ddc23a9a877258aa7f9ca7179585343b3beb48a9ef7ac2a5275341e4d",
	"libtcl9tkimgbmp2.0.1.dylib":    "e6fafd6ac968d38aa63f3a4bfbbb1f9f9a0bca04c67df8511e7f12548338163d",
	"libtcl9tkimgdted2.0.1.dylib":   "479b1d115b0624b3544da39af500464e4f6bc131a52c69fb7b48c81f3f28ec0d",
	"libtcl9tkimgflir2.0.1.dylib":   "a585bf1ca181d1b41e5753a976cb55134414d06b1897dea1e809ae890d51d4d5",
	"libtcl9tkimggif2.0.1.dylib":    "ffeac35cbe146942a28b6a19ee11b472c4a43b572bc084e7ec0789e814e02a95",
	"libtcl9tkimgico2.0.1.dylib":    "a007a7206060d3927bfe0d12569f878e02793b48223e3e168b62a5b774ac15e3",
	"libtcl9tkimgjpeg2.0.1.dylib":   "2cc5793ca16512301cff58df3908cb7337fcfbbc3cee832b10e81bf2432ebca4",
	"libtcl9tkimgpcx2.0.1.dylib":    "a5579dfc99cb46c0fbbe6b1b63e4767bf0f60286a0e8a939ce7d3e20a6eaa0d0",
	"libtcl9tkimgpixmap2.0.1.dylib": "a3d8d3dde4799729b4c1b77a6f730139b7c20c8cd7e250038244d73a9591379d",
	"libtcl9tkimgpng2.0.1.dylib":    "9e757d3b5e7bb8703a3bb7540b2ddb18e6a6381a1fe118131a2150db1d6554e5",
	"libtcl9tkimgppm2.0.1.dylib":    "b8eacb609006afae5973b989ed8fe94928ca8d4dbdc5ab4353179a669f2c863b",
	"libtcl9tkimgps2.0.1.dylib":     "0ea91f531e53d5771b9a93a89011a85d7c72164dfcddefe98c941980c8c8af95",
	"libtcl9tkimgraw2.0.1.dylib":    "0836fe10f5e5d96fd2ac546c29fb14cecd1c214c27670b6c2d79a349874d55d8",
	"libtcl9tkimgsgi2.0.1.dylib":    "aceb59be66711bae1489f5a09990b134c0537b1cc629b31d60285a4b74505044",
	"libtcl9tkimgsun2.0.1.dylib":    "5775c46164f1359a0f8ba61d9bbdfc3bf16f041307e1f6641b75ab8789cfe913",
	"libtcl9tkimgtga2.0.1.dylib":    "a4142b529d43a5e2e7061253c97059f4f2e929ddea1f8c50cdd10a96c4c71e75",
	"libtcl9tkimgtiff2.0.1.dylib":   "7f59cbb29d7fe53a4d80c5e44f48547ce815d3a6cd1bb509ec9d7c37623a148f",
	"libtcl9tkimgwindow2.0.1.dylib": "2be07494e1c797cb37f3776b1ddc08db7793b93bd608bfcd1e1aca6351efa913",
	"libtcl9tkimgxbm2.0.1.dylib":    "ac73c5a6552dcac8deedefc9cb07ed89ed86b3d7a51487d3506c1531e49f38bf",
	"libtcl9tkimgxpm2.0.1.dylib":    "7ba32f728870352954273a01a98ae51817839b1c853ad8b41444ab962263f570",
	"libtcl9zlibtcl1.3.1.dylib":     "3e0eb48ae321198a9b0fb7c7324a874d92ca25f55c91f0a28a88300293b92593",
	"libtk9.0.1.zip":                "3c68df95197c9835d77fb0362a4661fcb62b1b05f0769ee7191f6c5ede94c40d",
}
