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

//go:embed embed/linux/amd64/lib.zip
var libZip []byte

// Keep last for internal/shasig.go to update.
var shasig = map[string]string{
	"libtcl9.0.so":               "199c015b3cc49dd91a91440d3507831d71444eddf1c183b2ed8df9f48d7bb626",
	"libtcl9jpegtcl9.6.0.so":     "c3651f8adfe85434bfc0ca45bfce639c2c980d69ce48d0d3b1f49c60625b885f",
	"libtcl9pngtcl1.6.44.so":     "f9bf0d551a3ffc79d8545068e9064d0f31f0dffef908769eea80dcc283a65695",
	"libtcl9tifftcl4.7.0.so":     "e07d6917a5d3ba638242688b449655cf165ac73c5f771fd40e63719781b33c95",
	"libtcl9tk9.0.so":            "0f650d9191fe5ee9815a392c4a0fe2c643a8d481328bb94fcb4f5346abb1257a",
	"libtcl9tkimg2.0.1.so":       "50b8eec87b7b461933bb04b4cdab42f82ada1e009defe9ecab6cde9d483ae535",
	"libtcl9tkimgbmp2.0.1.so":    "6976463614feb84f886e43a821d922d7aede71413ed3b7bdc7620f5e14a3de86",
	"libtcl9tkimgdted2.0.1.so":   "7b6d65f576bd3fed71d5a043d6f3ec9ce97602d49697433c6b4d04d1b0528a37",
	"libtcl9tkimgflir2.0.1.so":   "1d53f97cf67ad6f1634b17f910a191eb871a9762dc8ee50ecf74b4bd37699aaa",
	"libtcl9tkimggif2.0.1.so":    "2e7f8bf3a24b8010f0f9bcc26728fc6e5e8df0069261077a74e58ebfec869844",
	"libtcl9tkimgico2.0.1.so":    "90032113838ad09b6b987c4e573fbcb3d8ed993465dc28c220d004b83eeb9983",
	"libtcl9tkimgjpeg2.0.1.so":   "b8326d56d26160e51bc31cc29c060f70d2a7f6fe105554dcd47f0049e44e560a",
	"libtcl9tkimgpcx2.0.1.so":    "f7c231d6a67ab2608dd77c200ce9f7437d2eec6a762ea1b1bb32a1232ce77296",
	"libtcl9tkimgpixmap2.0.1.so": "9b361125a8cf82aa274ae15ac1245df3234a22348f0df166229f96be92f13ea6",
	"libtcl9tkimgpng2.0.1.so":    "e91f9f63f726baeba81d60d119ea74b3d3be60b1861b2a383bafdb68e2fd40c9",
	"libtcl9tkimgppm2.0.1.so":    "d7605e5fedb3e73cf044b6fec60c0ac9824bb1abccb4d0970dc08c0f93e587fe",
	"libtcl9tkimgps2.0.1.so":     "fa1b0fdc69fa318ef8514b91b7ef7642652fbc285024e21d362c50c52f25302c",
	"libtcl9tkimgraw2.0.1.so":    "2ac4aac32b8cd799664ea8c26774c5c35f293615f5906a959dd649278e0bbbf8",
	"libtcl9tkimgsgi2.0.1.so":    "3750775950213a46c12efbea5089eb602e28dae8e8f85696f0d749d2bb8bf18f",
	"libtcl9tkimgsun2.0.1.so":    "f39d7ce28fba6716d68f9d166214e084782ff2d5315685820eac0b16d8aab857",
	"libtcl9tkimgtga2.0.1.so":    "913180257b165fd376bf8a7d000e7fa2fefbf4058c535535e529f80371d8f7b1",
	"libtcl9tkimgtiff2.0.1.so":   "ca6abbb8853b0de36e18231c50d67fd392c14aa39857dc4385443258c6190b3f",
	"libtcl9tkimgwindow2.0.1.so": "0b472db46efabea2b15b8b844268a0b01a84d2bfe5180c9745d0f036d49cfc81",
	"libtcl9tkimgxbm2.0.1.so":    "65a16d18db9cc1f5e4cbb524a68c6c093d96c9ad271d18799a74e5d31f83916d",
	"libtcl9tkimgxpm2.0.1.so":    "0a11e3c5005b6afc1aea4d21af326c0dd074718b5e7dba86446a193257dcd8d7",
	"libtcl9zlibtcl1.3.1.so":     "52e18151eeee25bccda53704e21355e67e2577c1c5d60f7fd70b4bcda7575948",
	"libtk9.0.1.zip":             "9d9dd98fa801a71b5d2894e33111a3857cf10213d563f5114461091f07c3df46",
}
