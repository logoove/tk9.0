// Copyright 2024 The tk9.0-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tk9_0 // import "modernc.org/tk9.0"

import (
	"fmt"
	"strings"
)

type cursor string

func (c cursor) optionString(*Window) string {
	return fmt.Sprintf(`-cursor %s`, string(c))
}

// Named cursor constants recognized on all platforms.
//
// Writing, for example, 'Arrow' as an option is the same as writing
// 'Cursor("arrow")', but the compiler will catch any typos in the cursor name.
const (
	Arrow             = cursor("arrow")
	BasedArrowDown    = cursor("based_arrow_down")
	BasedArrowUp      = cursor("based_arrow_up")
	Boat              = cursor("boat")
	Bogosity          = cursor("bogosity")
	BottomLeftCorner  = cursor("bottom_left_corner")
	BottomRightCorner = cursor("bottom_right_corner")
	BottomSide        = cursor("bottom_side")
	BottomTee         = cursor("bottom_tee")
	BoxSpiral         = cursor("box_spiral")
	CenterPtr         = cursor("center_ptr")
	Circle            = cursor("circle")
	Clock             = cursor("clock")
	CoffeeMug         = cursor("coffee_mug")
	Cross             = cursor("cross")
	CrossReverse      = cursor("cross_reverse")
	Crosshair         = cursor("crosshair")
	CursorIcon        = cursor("icon")
	DiamondCross      = cursor("diamond_cross")
	Dot               = cursor("dot")
	Dotbox            = cursor("dotbox")
	DoubleArrow       = cursor("double_arrow")
	DraftLarge        = cursor("draft_large")
	DraftSmall        = cursor("draft_small")
	DrapedBox         = cursor("draped_box")
	Exchange          = cursor("exchange")
	Fleur             = cursor("fleur")
	Gobbler           = cursor("gobbler")
	Gumby             = cursor("gumby")
	Hand1             = cursor("hand1")
	Hand2             = cursor("hand2")
	Heart             = cursor("heart")
	IronCross         = cursor("iron_cross")
	LeftPtr           = cursor("left_ptr")
	LeftSide          = cursor("left_side")
	LeftTee           = cursor("left_tee")
	Leftbutton        = cursor("leftbutton")
	LlAngle           = cursor("ll_angle")
	LrAngle           = cursor("lr_angle")
	Man               = cursor("man")
	Middlebutton      = cursor("middlebutton")
	Mouse             = cursor("mouse")
	None              = cursor("none")
	Pencil            = cursor("pencil")
	Pirate            = cursor("pirate")
	Plus              = cursor("plus")
	QuestionArrow     = cursor("question_arrow")
	RightPtr          = cursor("right_ptr")
	RightSide         = cursor("right_side")
	RightTee          = cursor("right_tee")
	Rightbutton       = cursor("rightbutton")
	RtlLogo           = cursor("rtl_logo")
	Sailboat          = cursor("sailboat")
	SbDownArrow       = cursor("sb_down_arrow")
	SbHDoubleArrow    = cursor("sb_h_double_arrow")
	SbLeftArrow       = cursor("sb_left_arrow")
	SbRightArrow      = cursor("sb_right_arrow")
	SbUpArrow         = cursor("sb_up_arrow")
	SbVDoubleArrow    = cursor("sb_v_double_arrow")
	Shuttle           = cursor("shuttle")
	Sizing            = cursor("sizing")
	Spider            = cursor("spider")
	Spraycan          = cursor("spraycan")
	Star              = cursor("star")
	Target            = cursor("target")
	Tcross            = cursor("tcross")
	TopLeftArrow      = cursor("top_left_arrow")
	TopLeftCorner     = cursor("top_left_corner")
	TopRightCorner    = cursor("top_right_corner")
	TopSide           = cursor("top_side")
	TopTee            = cursor("top_tee")
	Trek              = cursor("trek")
	UlAngle           = cursor("ul_angle")
	Umbrella          = cursor("umbrella")
	UrAngle           = cursor("ur_angle")
	Watch             = cursor("watch")
	XCursor           = cursor("X_cursor")
	Xterm             = cursor("xterm")
)

// Additional Windows cursors.
const (
	CursorSize = cursor("size")
	No         = cursor("no")
	SizeNeSw   = cursor("size_ne_sw")
	SizeNs     = cursor("size_ns")
	SizeNwSe   = cursor("size_nw_se")
	SizeWe     = cursor("size_we")
	Starting   = cursor("starting")
	Uparrow    = cursor("uparrow")
)

// Additional macOS cursors.
const (
	Aliasarrow            = cursor("aliasarrow")
	Bucket                = cursor("bucket")
	Cancel                = cursor("cancel")
	Closedhand            = cursor("closedhand")
	Contextualmenuarrow   = cursor("contextualmenuarrow")
	Copyarrow             = cursor("copyarrow")
	Countingdownhand      = cursor("countingdownhand")
	Countingupanddownhand = cursor("countingupanddownhand")
	Countinguphand        = cursor("countinguphand")
	CrossHair             = cursor("cross-hair")
	CursorText            = cursor("text")
	Eyedrop               = cursor("eyedrop")
	EyedropFull           = cursor("eyedrop-full")
	Fist                  = cursor("fist")
	Hand                  = cursor("hand")
	Help                  = cursor("help")
	Movearrow             = cursor("movearrow")
	Notallowed            = cursor("notallowed")
	Openhand              = cursor("openhand")
	Pointinghand          = cursor("pointinghand")
	Poof                  = cursor("poof")
	Resize                = cursor("resize")
	Resizebottomleft      = cursor("resizebottomleft")
	Resizebottomright     = cursor("resizebottomright")
	Resizedown            = cursor("resizedown")
	Resizeleft            = cursor("resizeleft")
	Resizeleftright       = cursor("resizeleftright")
	Resizeright           = cursor("resizeright")
	Resizetopleft         = cursor("resizetopleft")
	Resizetopright        = cursor("resizetopright")
	Resizeup              = cursor("resizeup")
	Resizeupdown          = cursor("resizeupdown")
	Spinning              = cursor("spinning")
	ZoomIn                = cursor("zoom-in")
	ZoomOut               = cursor("zoom-out")
)

// Additional macOS and Windows cursors.
const (
	Wait = cursor("wait")
)

// Named colors. Writing, for example, 'Background(Agua)' is the same as
// writing 'Background("agua")' but the compiler can catch any typos in the
// color name.
const (
	Agua                 = "agua"                 // R:0 G:255 B:255
	AliceBlue            = "AliceBlue"            // R:240 G:248 B:255
	AntiqueWhite         = "AntiqueWhite"         // R:250 G:235 B:215
	AntiqueWhite1        = "AntiqueWhite1"        // R:255 G:239 B:219
	AntiqueWhite2        = "AntiqueWhite2"        // R:238 G:223 B:204
	AntiqueWhite3        = "AntiqueWhite3"        // R:205 G:192 B:176
	AntiqueWhite4        = "AntiqueWhite4"        // R:139 G:131 B:120
	Aquamarine           = "aquamarine"           // R:127 G:255 B:212
	Aquamarine1          = "aquamarine1"          // R:127 G:255 B:212
	Aquamarine2          = "aquamarine2"          // R:118 G:238 B:198
	Aquamarine3          = "aquamarine3"          // R:102 G:205 B:170
	Aquamarine4          = "aquamarine4"          // R:69 G:139 B:116
	Azure                = "azure"                // R:240 G:255 B:255
	Azure1               = "azure1"               // R:240 G:255 B:255
	Azure2               = "azure2"               // R:224 G:238 B:238
	Azure3               = "azure3"               // R:193 G:205 B:205
	Azure4               = "azure4"               // R:131 G:139 B:139
	Beige                = "beige"                // R:245 G:245 B:220
	Bisque               = "bisque"               // R:255 G:228 B:196
	Bisque1              = "bisque1"              // R:255 G:228 B:196
	Bisque2              = "bisque2"              // R:238 G:213 B:183
	Bisque3              = "bisque3"              // R:205 G:183 B:158
	Bisque4              = "bisque4"              // R:139 G:125 B:107
	Black                = "black"                // R:0 G:0 B:0
	BlanchedAlmond       = "BlanchedAlmond"       // R:255 G:235 B:205
	Blue                 = "blue"                 // R:0 G:0 B:255
	Blue1                = "blue1"                // R:0 G:0 B:255
	Blue2                = "blue2"                // R:0 G:0 B:238
	Blue3                = "blue3"                // R:0 G:0 B:205
	Blue4                = "blue4"                // R:0 G:0 B:139
	BlueViolet           = "BlueViolet"           // R:138 G:43 B:226
	Brown                = "brown"                // R:165 G:42 B:42
	Brown1               = "brown1"               // R:255 G:64 B:64
	Brown2               = "brown2"               // R:238 G:59 B:59
	Brown3               = "brown3"               // R:205 G:51 B:51
	Brown4               = "brown4"               // R:139 G:35 B:35
	Burlywood            = "burlywood"            // R:222 G:184 B:135
	Burlywood1           = "burlywood1"           // R:255 G:211 B:155
	Burlywood2           = "burlywood2"           // R:238 G:197 B:145
	Burlywood3           = "burlywood3"           // R:205 G:170 B:125
	Burlywood4           = "burlywood4"           // R:139 G:115 B:85
	CadetBlue            = "CadetBlue"            // R:95 G:158 B:160
	CadetBlue1           = "CadetBlue1"           // R:152 G:245 B:255
	CadetBlue2           = "CadetBlue2"           // R:142 G:229 B:238
	CadetBlue3           = "CadetBlue3"           // R:122 G:197 B:205
	CadetBlue4           = "CadetBlue4"           // R:83 G:134 B:139
	Chartreuse           = "chartreuse"           // R:127 G:255 B:0
	Chartreuse1          = "chartreuse1"          // R:127 G:255 B:0
	Chartreuse2          = "chartreuse2"          // R:118 G:238 B:0
	Chartreuse3          = "chartreuse3"          // R:102 G:205 B:0
	Chartreuse4          = "chartreuse4"          // R:69 G:139 B:0
	Chocolate            = "chocolate"            // R:210 G:105 B:30
	Chocolate1           = "chocolate1"           // R:255 G:127 B:36
	Chocolate2           = "chocolate2"           // R:238 G:118 B:33
	Chocolate3           = "chocolate3"           // R:205 G:102 B:29
	Chocolate4           = "chocolate4"           // R:139 G:69 B:19
	Coral                = "coral"                // R:255 G:127 B:80
	Coral1               = "coral1"               // R:255 G:114 B:86
	Coral2               = "coral2"               // R:238 G:106 B:80
	Coral3               = "coral3"               // R:205 G:91 B:69
	Coral4               = "coral4"               // R:139 G:62 B:47
	CornflowerBlue       = "CornflowerBlue"       // R:100 G:149 B:237
	Cornsilk             = "cornsilk"             // R:255 G:248 B:220
	Cornsilk1            = "cornsilk1"            // R:255 G:248 B:220
	Cornsilk2            = "cornsilk2"            // R:238 G:232 B:205
	Cornsilk3            = "cornsilk3"            // R:205 G:200 B:177
	Cornsilk4            = "cornsilk4"            // R:139 G:136 B:120
	Crymson              = "crymson"              // R:220 G:20 B:60
	Cyan                 = "cyan"                 // R:0 G:255 B:255
	Cyan1                = "cyan1"                // R:0 G:255 B:255
	Cyan2                = "cyan2"                // R:0 G:238 B:238
	Cyan3                = "cyan3"                // R:0 G:205 B:205
	Cyan4                = "cyan4"                // R:0 G:139 B:139
	DarkBlue             = "DarkBlue"             // R:0 G:0 B:139
	DarkCyan             = "DarkCyan"             // R:0 G:139 B:139
	DarkGoldenrod        = "DarkGoldenrod"        // R:184 G:134 B:11
	DarkGoldenrod1       = "DarkGoldenrod1"       // R:255 G:185 B:15
	DarkGoldenrod2       = "DarkGoldenrod2"       // R:238 G:173 B:14
	DarkGoldenrod3       = "DarkGoldenrod3"       // R:205 G:149 B:12
	DarkGoldenrod4       = "DarkGoldenrod4"       // R:139 G:101 B:8
	DarkGray             = "DarkGray"             // R:169 G:169 B:169
	DarkGreen            = "DarkGreen"            // R:0 G:100 B:0
	DarkGrey             = "DarkGrey"             // R:169 G:169 B:169
	DarkKhaki            = "DarkKhaki"            // R:189 G:183 B:107
	DarkMagenta          = "DarkMagenta"          // R:139 G:0 B:139
	DarkOliveGreen       = "DarkOliveGreen"       // R:85 G:107 B:47
	DarkOliveGreen1      = "DarkOliveGreen1"      // R:202 G:255 B:112
	DarkOliveGreen2      = "DarkOliveGreen2"      // R:188 G:238 B:104
	DarkOliveGreen3      = "DarkOliveGreen3"      // R:162 G:205 B:90
	DarkOliveGreen4      = "DarkOliveGreen4"      // R:110 G:139 B:61
	DarkOrange           = "DarkOrange"           // R:255 G:140 B:0
	DarkOrange1          = "DarkOrange1"          // R:255 G:127 B:0
	DarkOrange2          = "DarkOrange2"          // R:238 G:118 B:0
	DarkOrange3          = "DarkOrange3"          // R:205 G:102 B:0
	DarkOrange4          = "DarkOrange4"          // R:139 G:69 B:0
	DarkOrchid           = "DarkOrchid"           // R:153 G:50 B:204
	DarkOrchid1          = "DarkOrchid1"          // R:191 G:62 B:255
	DarkOrchid2          = "DarkOrchid2"          // R:178 G:58 B:238
	DarkOrchid3          = "DarkOrchid3"          // R:154 G:50 B:205
	DarkOrchid4          = "DarkOrchid4"          // R:104 G:34 B:139
	DarkRed              = "DarkRed"              // R:139 G:0 B:0
	DarkSalmon           = "DarkSalmon"           // R:233 G:150 B:122
	DarkSeaGreen         = "DarkSeaGreen"         // R:143 G:188 B:143
	DarkSeaGreen1        = "DarkSeaGreen1"        // R:193 G:255 B:193
	DarkSeaGreen2        = "DarkSeaGreen2"        // R:180 G:238 B:180
	DarkSeaGreen3        = "DarkSeaGreen3"        // R:155 G:205 B:155
	DarkSeaGreen4        = "DarkSeaGreen4"        // R:105 G:139 B:105
	DarkSlateBlue        = "DarkSlateBlue"        // R:72 G:61 B:139
	DarkSlateGray        = "DarkSlateGray"        // R:47 G:79 B:79
	DarkSlateGray1       = "DarkSlateGray1"       // R:151 G:255 B:255
	DarkSlateGray2       = "DarkSlateGray2"       // R:141 G:238 B:238
	DarkSlateGray3       = "DarkSlateGray3"       // R:121 G:205 B:205
	DarkSlateGray4       = "DarkSlateGray4"       // R:82 G:139 B:139
	DarkSlateGrey        = "DarkSlateGrey"        // R:47 G:79 B:79
	DarkTurquoise        = "DarkTurquoise"        // R:0 G:206 B:209
	DarkViolet           = "DarkViolet"           // R:148 G:0 B:211
	DeepPink             = "DeepPink"             // R:255 G:20 B:147
	DeepPink1            = "DeepPink1"            // R:255 G:20 B:147
	DeepPink2            = "DeepPink2"            // R:238 G:18 B:137
	DeepPink3            = "DeepPink3"            // R:205 G:16 B:118
	DeepPink4            = "DeepPink4"            // R:139 G:10 B:80
	DeepSkyBlue          = "DeepSkyBlue"          // R:0 G:191 B:255
	DeepSkyBlue1         = "DeepSkyBlue1"         // R:0 G:191 B:255
	DeepSkyBlue2         = "DeepSkyBlue2"         // R:0 G:178 B:238
	DeepSkyBlue3         = "DeepSkyBlue3"         // R:0 G:154 B:205
	DeepSkyBlue4         = "DeepSkyBlue4"         // R:0 G:104 B:139
	DimGray              = "DimGray"              // R:105 G:105 B:105
	DimGrey              = "DimGrey"              // R:105 G:105 B:105
	DodgerBlue           = "DodgerBlue"           // R:30 G:144 B:255
	DodgerBlue1          = "DodgerBlue1"          // R:30 G:144 B:255
	DodgerBlue2          = "DodgerBlue2"          // R:28 G:134 B:238
	DodgerBlue3          = "DodgerBlue3"          // R:24 G:116 B:205
	DodgerBlue4          = "DodgerBlue4"          // R:16 G:78 B:139
	Firebrick            = "firebrick"            // R:178 G:34 B:34
	Firebrick1           = "firebrick1"           // R:255 G:48 B:48
	Firebrick2           = "firebrick2"           // R:238 G:44 B:44
	Firebrick3           = "firebrick3"           // R:205 G:38 B:38
	Firebrick4           = "firebrick4"           // R:139 G:26 B:26
	FloralWhite          = "FloralWhite"          // R:255 G:250 B:240
	ForestGreen          = "ForestGreen"          // R:34 G:139 B:34
	Fuchsia              = "fuchsia"              // R:255 G:0 B:255
	Gainsboro            = "gainsboro"            // R:220 G:220 B:220
	GhostWhite           = "GhostWhite"           // R:248 G:248 B:255
	Gold                 = "gold"                 // R:255 G:215 B:0
	Gold1                = "gold1"                // R:255 G:215 B:0
	Gold2                = "gold2"                // R:238 G:201 B:0
	Gold3                = "gold3"                // R:205 G:173 B:0
	Gold4                = "gold4"                // R:139 G:117 B:0
	Goldenrod            = "goldenrod"            // R:218 G:165 B:32
	Goldenrod1           = "goldenrod1"           // R:255 G:193 B:37
	Goldenrod2           = "goldenrod2"           // R:238 G:180 B:34
	Goldenrod3           = "goldenrod3"           // R:205 G:155 B:29
	Goldenrod4           = "goldenrod4"           // R:139 G:105 B:20
	Gray                 = "gray"                 // R:128 G:128 B:128
	Gray0                = "gray0"                // R:0 G:0 B:0
	Gray1                = "gray1"                // R:3 G:3 B:3
	Gray10               = "gray10"               // R:26 G:26 B:26
	Gray100              = "gray100"              // R:255 G:255 B:255
	Gray11               = "gray11"               // R:28 G:28 B:28
	Gray12               = "gray12"               // R:31 G:31 B:31
	Gray13               = "gray13"               // R:33 G:33 B:33
	Gray14               = "gray14"               // R:36 G:36 B:36
	Gray15               = "gray15"               // R:38 G:38 B:38
	Gray16               = "gray16"               // R:41 G:41 B:41
	Gray17               = "gray17"               // R:43 G:43 B:43
	Gray18               = "gray18"               // R:46 G:46 B:46
	Gray19               = "gray19"               // R:48 G:48 B:48
	Gray2                = "gray2"                // R:5 G:5 B:5
	Gray20               = "gray20"               // R:51 G:51 B:51
	Gray21               = "gray21"               // R:54 G:54 B:54
	Gray22               = "gray22"               // R:56 G:56 B:56
	Gray23               = "gray23"               // R:59 G:59 B:59
	Gray24               = "gray24"               // R:61 G:61 B:61
	Gray25               = "gray25"               // R:64 G:64 B:64
	Gray26               = "gray26"               // R:66 G:66 B:66
	Gray27               = "gray27"               // R:69 G:69 B:69
	Gray28               = "gray28"               // R:71 G:71 B:71
	Gray29               = "gray29"               // R:74 G:74 B:74
	Gray3                = "gray3"                // R:8 G:8 B:8
	Gray30               = "gray30"               // R:77 G:77 B:77
	Gray31               = "gray31"               // R:79 G:79 B:79
	Gray32               = "gray32"               // R:82 G:82 B:82
	Gray33               = "gray33"               // R:84 G:84 B:84
	Gray34               = "gray34"               // R:87 G:87 B:87
	Gray35               = "gray35"               // R:89 G:89 B:89
	Gray36               = "gray36"               // R:92 G:92 B:92
	Gray37               = "gray37"               // R:94 G:94 B:94
	Gray38               = "gray38"               // R:97 G:97 B:97
	Gray39               = "gray39"               // R:99 G:99 B:99
	Gray4                = "gray4"                // R:10 G:10 B:10
	Gray40               = "gray40"               // R:102 G:102 B:102
	Gray41               = "gray41"               // R:105 G:105 B:105
	Gray42               = "gray42"               // R:107 G:107 B:107
	Gray43               = "gray43"               // R:110 G:110 B:110
	Gray44               = "gray44"               // R:112 G:112 B:112
	Gray45               = "gray45"               // R:115 G:115 B:115
	Gray46               = "gray46"               // R:117 G:117 B:117
	Gray47               = "gray47"               // R:120 G:120 B:120
	Gray48               = "gray48"               // R:122 G:122 B:122
	Gray49               = "gray49"               // R:125 G:125 B:125
	Gray5                = "gray5"                // R:13 G:13 B:13
	Gray50               = "gray50"               // R:127 G:127 B:127
	Gray51               = "gray51"               // R:130 G:130 B:130
	Gray52               = "gray52"               // R:133 G:133 B:133
	Gray53               = "gray53"               // R:135 G:135 B:135
	Gray54               = "gray54"               // R:138 G:138 B:138
	Gray55               = "gray55"               // R:140 G:140 B:140
	Gray56               = "gray56"               // R:143 G:143 B:143
	Gray57               = "gray57"               // R:145 G:145 B:145
	Gray58               = "gray58"               // R:148 G:148 B:148
	Gray59               = "gray59"               // R:150 G:150 B:150
	Gray6                = "gray6"                // R:15 G:15 B:15
	Gray60               = "gray60"               // R:153 G:153 B:153
	Gray61               = "gray61"               // R:156 G:156 B:156
	Gray62               = "gray62"               // R:158 G:158 B:158
	Gray63               = "gray63"               // R:161 G:161 B:161
	Gray64               = "gray64"               // R:163 G:163 B:163
	Gray65               = "gray65"               // R:166 G:166 B:166
	Gray66               = "gray66"               // R:168 G:168 B:168
	Gray67               = "gray67"               // R:171 G:171 B:171
	Gray68               = "gray68"               // R:173 G:173 B:173
	Gray69               = "gray69"               // R:176 G:176 B:176
	Gray7                = "gray7"                // R:18 G:18 B:18
	Gray70               = "gray70"               // R:179 G:179 B:179
	Gray71               = "gray71"               // R:181 G:181 B:181
	Gray72               = "gray72"               // R:184 G:184 B:184
	Gray73               = "gray73"               // R:186 G:186 B:186
	Gray74               = "gray74"               // R:189 G:189 B:189
	Gray75               = "gray75"               // R:191 G:191 B:191
	Gray76               = "gray76"               // R:194 G:194 B:194
	Gray77               = "gray77"               // R:196 G:196 B:196
	Gray78               = "gray78"               // R:199 G:199 B:199
	Gray79               = "gray79"               // R:201 G:201 B:201
	Gray8                = "gray8"                // R:20 G:20 B:20
	Gray80               = "gray80"               // R:204 G:204 B:204
	Gray81               = "gray81"               // R:207 G:207 B:207
	Gray82               = "gray82"               // R:209 G:209 B:209
	Gray83               = "gray83"               // R:212 G:212 B:212
	Gray84               = "gray84"               // R:214 G:214 B:214
	Gray85               = "gray85"               // R:217 G:217 B:217
	Gray86               = "gray86"               // R:219 G:219 B:219
	Gray87               = "gray87"               // R:222 G:222 B:222
	Gray88               = "gray88"               // R:224 G:224 B:224
	Gray89               = "gray89"               // R:227 G:227 B:227
	Gray9                = "gray9"                // R:23 G:23 B:23
	Gray90               = "gray90"               // R:229 G:229 B:229
	Gray91               = "gray91"               // R:232 G:232 B:232
	Gray92               = "gray92"               // R:235 G:235 B:235
	Gray93               = "gray93"               // R:237 G:237 B:237
	Gray94               = "gray94"               // R:240 G:240 B:240
	Gray95               = "gray95"               // R:242 G:242 B:242
	Gray96               = "gray96"               // R:245 G:245 B:245
	Gray97               = "gray97"               // R:247 G:247 B:247
	Gray98               = "gray98"               // R:250 G:250 B:250
	Gray99               = "gray99"               // R:252 G:252 B:252
	Green                = "green"                // R:0 G:128 B:0
	Green1               = "green1"               // R:0 G:255 B:0
	Green2               = "green2"               // R:0 G:238 B:0
	Green3               = "green3"               // R:0 G:205 B:0
	Green4               = "green4"               // R:0 G:139 B:0
	GreenYellow          = "GreenYellow"          // R:173 G:255 B:47
	Grey                 = "grey"                 // R:128 G:128 B:128
	Grey0                = "grey0"                // R:0 G:0 B:0
	Grey1                = "grey1"                // R:3 G:3 B:3
	Grey10               = "grey10"               // R:26 G:26 B:26
	Grey100              = "grey100"              // R:255 G:255 B:255
	Grey11               = "grey11"               // R:28 G:28 B:28
	Grey12               = "grey12"               // R:31 G:31 B:31
	Grey13               = "grey13"               // R:33 G:33 B:33
	Grey14               = "grey14"               // R:36 G:36 B:36
	Grey15               = "grey15"               // R:38 G:38 B:38
	Grey16               = "grey16"               // R:41 G:41 B:41
	Grey17               = "grey17"               // R:43 G:43 B:43
	Grey18               = "grey18"               // R:46 G:46 B:46
	Grey19               = "grey19"               // R:48 G:48 B:48
	Grey2                = "grey2"                // R:5 G:5 B:5
	Grey20               = "grey20"               // R:51 G:51 B:51
	Grey21               = "grey21"               // R:54 G:54 B:54
	Grey22               = "grey22"               // R:56 G:56 B:56
	Grey23               = "grey23"               // R:59 G:59 B:59
	Grey24               = "grey24"               // R:61 G:61 B:61
	Grey25               = "grey25"               // R:64 G:64 B:64
	Grey26               = "grey26"               // R:66 G:66 B:66
	Grey27               = "grey27"               // R:69 G:69 B:69
	Grey28               = "grey28"               // R:71 G:71 B:71
	Grey29               = "grey29"               // R:74 G:74 B:74
	Grey3                = "grey3"                // R:8 G:8 B:8
	Grey30               = "grey30"               // R:77 G:77 B:77
	Grey31               = "grey31"               // R:79 G:79 B:79
	Grey32               = "grey32"               // R:82 G:82 B:82
	Grey33               = "grey33"               // R:84 G:84 B:84
	Grey34               = "grey34"               // R:87 G:87 B:87
	Grey35               = "grey35"               // R:89 G:89 B:89
	Grey36               = "grey36"               // R:92 G:92 B:92
	Grey37               = "grey37"               // R:94 G:94 B:94
	Grey38               = "grey38"               // R:97 G:97 B:97
	Grey39               = "grey39"               // R:99 G:99 B:99
	Grey4                = "grey4"                // R:10 G:10 B:10
	Grey40               = "grey40"               // R:102 G:102 B:102
	Grey41               = "grey41"               // R:105 G:105 B:105
	Grey42               = "grey42"               // R:107 G:107 B:107
	Grey43               = "grey43"               // R:110 G:110 B:110
	Grey44               = "grey44"               // R:112 G:112 B:112
	Grey45               = "grey45"               // R:115 G:115 B:115
	Grey46               = "grey46"               // R:117 G:117 B:117
	Grey47               = "grey47"               // R:120 G:120 B:120
	Grey48               = "grey48"               // R:122 G:122 B:122
	Grey49               = "grey49"               // R:125 G:125 B:125
	Grey5                = "grey5"                // R:13 G:13 B:13
	Grey50               = "grey50"               // R:127 G:127 B:127
	Grey51               = "grey51"               // R:130 G:130 B:130
	Grey52               = "grey52"               // R:133 G:133 B:133
	Grey53               = "grey53"               // R:135 G:135 B:135
	Grey54               = "grey54"               // R:138 G:138 B:138
	Grey55               = "grey55"               // R:140 G:140 B:140
	Grey56               = "grey56"               // R:143 G:143 B:143
	Grey57               = "grey57"               // R:145 G:145 B:145
	Grey58               = "grey58"               // R:148 G:148 B:148
	Grey59               = "grey59"               // R:150 G:150 B:150
	Grey6                = "grey6"                // R:15 G:15 B:15
	Grey60               = "grey60"               // R:153 G:153 B:153
	Grey61               = "grey61"               // R:156 G:156 B:156
	Grey62               = "grey62"               // R:158 G:158 B:158
	Grey63               = "grey63"               // R:161 G:161 B:161
	Grey64               = "grey64"               // R:163 G:163 B:163
	Grey65               = "grey65"               // R:166 G:166 B:166
	Grey66               = "grey66"               // R:168 G:168 B:168
	Grey67               = "grey67"               // R:171 G:171 B:171
	Grey68               = "grey68"               // R:173 G:173 B:173
	Grey69               = "grey69"               // R:176 G:176 B:176
	Grey7                = "grey7"                // R:18 G:18 B:18
	Grey70               = "grey70"               // R:179 G:179 B:179
	Grey71               = "grey71"               // R:181 G:181 B:181
	Grey72               = "grey72"               // R:184 G:184 B:184
	Grey73               = "grey73"               // R:186 G:186 B:186
	Grey74               = "grey74"               // R:189 G:189 B:189
	Grey75               = "grey75"               // R:191 G:191 B:191
	Grey76               = "grey76"               // R:194 G:194 B:194
	Grey77               = "grey77"               // R:196 G:196 B:196
	Grey78               = "grey78"               // R:199 G:199 B:199
	Grey79               = "grey79"               // R:201 G:201 B:201
	Grey8                = "grey8"                // R:20 G:20 B:20
	Grey80               = "grey80"               // R:204 G:204 B:204
	Grey81               = "grey81"               // R:207 G:207 B:207
	Grey82               = "grey82"               // R:209 G:209 B:209
	Grey83               = "grey83"               // R:212 G:212 B:212
	Grey84               = "grey84"               // R:214 G:214 B:214
	Grey85               = "grey85"               // R:217 G:217 B:217
	Grey86               = "grey86"               // R:219 G:219 B:219
	Grey87               = "grey87"               // R:222 G:222 B:222
	Grey88               = "grey88"               // R:224 G:224 B:224
	Grey89               = "grey89"               // R:227 G:227 B:227
	Grey9                = "grey9"                // R:23 G:23 B:23
	Grey90               = "grey90"               // R:229 G:229 B:229
	Grey91               = "grey91"               // R:232 G:232 B:232
	Grey92               = "grey92"               // R:235 G:235 B:235
	Grey93               = "grey93"               // R:237 G:237 B:237
	Grey94               = "grey94"               // R:240 G:240 B:240
	Grey95               = "grey95"               // R:242 G:242 B:242
	Grey96               = "grey96"               // R:245 G:245 B:245
	Grey97               = "grey97"               // R:247 G:247 B:247
	Grey98               = "grey98"               // R:250 G:250 B:250
	Grey99               = "grey99"               // R:252 G:252 B:252
	Honeydew             = "honeydew"             // R:240 G:255 B:240
	Honeydew1            = "honeydew1"            // R:240 G:255 B:240
	Honeydew2            = "honeydew2"            // R:224 G:238 B:224
	Honeydew3            = "honeydew3"            // R:193 G:205 B:193
	Honeydew4            = "honeydew4"            // R:131 G:139 B:131
	HotPink              = "HotPink"              // R:255 G:105 B:180
	HotPink1             = "HotPink1"             // R:255 G:110 B:180
	HotPink2             = "HotPink2"             // R:238 G:106 B:167
	HotPink3             = "HotPink3"             // R:205 G:96 B:144
	HotPink4             = "HotPink4"             // R:139 G:58 B:98
	IndianRed            = "IndianRed"            // R:205 G:92 B:92
	IndianRed1           = "IndianRed1"           // R:255 G:106 B:106
	IndianRed2           = "IndianRed2"           // R:238 G:99 B:99
	IndianRed3           = "IndianRed3"           // R:205 G:85 B:85
	IndianRed4           = "IndianRed4"           // R:139 G:58 B:58
	Indigo               = "indigo"               // R:75 G:0 B:130
	Ivory                = "ivory"                // R:255 G:255 B:240
	Ivory1               = "ivory1"               // R:255 G:255 B:240
	Ivory2               = "ivory2"               // R:238 G:238 B:224
	Ivory3               = "ivory3"               // R:205 G:205 B:193
	Ivory4               = "ivory4"               // R:139 G:139 B:131
	Khaki                = "khaki"                // R:240 G:230 B:140
	Khaki1               = "khaki1"               // R:255 G:246 B:143
	Khaki2               = "khaki2"               // R:238 G:230 B:133
	Khaki3               = "khaki3"               // R:205 G:198 B:115
	Khaki4               = "khaki4"               // R:139 G:134 B:78
	Lavender             = "lavender"             // R:230 G:230 B:250
	LavenderBlush        = "LavenderBlush"        // R:255 G:240 B:245
	LavenderBlush1       = "LavenderBlush1"       // R:255 G:240 B:245
	LavenderBlush2       = "LavenderBlush2"       // R:238 G:224 B:229
	LavenderBlush3       = "LavenderBlush3"       // R:205 G:193 B:197
	LavenderBlush4       = "LavenderBlush4"       // R:139 G:131 B:134
	LawnGreen            = "LawnGreen"            // R:124 G:252 B:0
	LemonChiffon         = "LemonChiffon"         // R:255 G:250 B:205
	LemonChiffon1        = "LemonChiffon1"        // R:255 G:250 B:205
	LemonChiffon2        = "LemonChiffon2"        // R:238 G:233 B:191
	LemonChiffon3        = "LemonChiffon3"        // R:205 G:201 B:165
	LemonChiffon4        = "LemonChiffon4"        // R:139 G:137 B:112
	LightBlue            = "LightBlue"            // R:173 G:216 B:230
	LightBlue1           = "LightBlue1"           // R:191 G:239 B:255
	LightBlue2           = "LightBlue2"           // R:178 G:223 B:238
	LightBlue3           = "LightBlue3"           // R:154 G:192 B:205
	LightBlue4           = "LightBlue4"           // R:104 G:131 B:139
	LightCoral           = "LightCoral"           // R:240 G:128 B:128
	LightCyan            = "LightCyan"            // R:224 G:255 B:255
	LightCyan1           = "LightCyan1"           // R:224 G:255 B:255
	LightCyan2           = "LightCyan2"           // R:209 G:238 B:238
	LightCyan3           = "LightCyan3"           // R:180 G:205 B:205
	LightCyan4           = "LightCyan4"           // R:122 G:139 B:139
	LightGoldenrod       = "LightGoldenrod"       // R:238 G:221 B:130
	LightGoldenrod1      = "LightGoldenrod1"      // R:255 G:236 B:139
	LightGoldenrod2      = "LightGoldenrod2"      // R:238 G:220 B:130
	LightGoldenrod3      = "LightGoldenrod3"      // R:205 G:190 B:112
	LightGoldenrod4      = "LightGoldenrod4"      // R:139 G:129 B:76
	LightGoldenrodYellow = "LightGoldenrodYellow" // R:250 G:250 B:210
	LightGray            = "LightGray"            // R:211 G:211 B:211
	LightGreen           = "LightGreen"           // R:144 G:238 B:144
	LightGrey            = "LightGrey"            // R:211 G:211 B:211
	LightPink            = "LightPink"            // R:255 G:182 B:193
	LightPink1           = "LightPink1"           // R:255 G:174 B:185
	LightPink2           = "LightPink2"           // R:238 G:162 B:173
	LightPink3           = "LightPink3"           // R:205 G:140 B:149
	LightPink4           = "LightPink4"           // R:139 G:95 B:101
	LightSalmon          = "LightSalmon"          // R:255 G:160 B:122
	LightSalmon1         = "LightSalmon1"         // R:255 G:160 B:122
	LightSalmon2         = "LightSalmon2"         // R:238 G:149 B:114
	LightSalmon3         = "LightSalmon3"         // R:205 G:129 B:98
	LightSalmon4         = "LightSalmon4"         // R:139 G:87 B:66
	LightSeaGreen        = "LightSeaGreen"        // R:32 G:178 B:170
	LightSkyBlue         = "LightSkyBlue"         // R:135 G:206 B:250
	LightSkyBlue1        = "LightSkyBlue1"        // R:176 G:226 B:255
	LightSkyBlue2        = "LightSkyBlue2"        // R:164 G:211 B:238
	LightSkyBlue3        = "LightSkyBlue3"        // R:141 G:182 B:205
	LightSkyBlue4        = "LightSkyBlue4"        // R:96 G:123 B:139
	LightSlateBlue       = "LightSlateBlue"       // R:132 G:112 B:255
	LightSlateGray       = "LightSlateGray"       // R:119 G:136 B:153
	LightSlateGrey       = "LightSlateGrey"       // R:119 G:136 B:153
	LightSteelBlue       = "LightSteelBlue"       // R:176 G:196 B:222
	LightSteelBlue1      = "LightSteelBlue1"      // R:202 G:225 B:255
	LightSteelBlue2      = "LightSteelBlue2"      // R:188 G:210 B:238
	LightSteelBlue3      = "LightSteelBlue3"      // R:162 G:181 B:205
	LightSteelBlue4      = "LightSteelBlue4"      // R:110 G:123 B:139
	LightYellow          = "LightYellow"          // R:255 G:255 B:224
	LightYellow1         = "LightYellow1"         // R:255 G:255 B:224
	LightYellow2         = "LightYellow2"         // R:238 G:238 B:209
	LightYellow3         = "LightYellow3"         // R:205 G:205 B:180
	LightYellow4         = "LightYellow4"         // R:139 G:139 B:122
	Lime                 = "lime"                 // R:0 G:255 B:0
	LimeGreen            = "LimeGreen"            // R:50 G:205 B:50
	Linen                = "linen"                // R:250 G:240 B:230
	Magenta              = "magenta"              // R:255 G:0 B:255
	Magenta1             = "magenta1"             // R:255 G:0 B:255
	Magenta2             = "magenta2"             // R:238 G:0 B:238
	Magenta3             = "magenta3"             // R:205 G:0 B:205
	Magenta4             = "magenta4"             // R:139 G:0 B:139
	Maroon               = "maroon"               // R:128 G:0 B:0
	Maroon1              = "maroon1"              // R:255 G:52 B:179
	Maroon2              = "maroon2"              // R:238 G:48 B:167
	Maroon3              = "maroon3"              // R:205 G:41 B:144
	Maroon4              = "maroon4"              // R:139 G:28 B:98
	MediumAquamarine     = "MediumAquamarine"     // R:102 G:205 B:170
	MediumBlue           = "MediumBlue"           // R:0 G:0 B:205
	MediumOrchid         = "MediumOrchid"         // R:186 G:85 B:211
	MediumOrchid1        = "MediumOrchid1"        // R:224 G:102 B:255
	MediumOrchid2        = "MediumOrchid2"        // R:209 G:95 B:238
	MediumOrchid3        = "MediumOrchid3"        // R:180 G:82 B:205
	MediumOrchid4        = "MediumOrchid4"        // R:122 G:55 B:139
	MediumPurple         = "MediumPurple"         // R:147 G:112 B:219
	MediumPurple1        = "MediumPurple1"        // R:171 G:130 B:255
	MediumPurple2        = "MediumPurple2"        // R:159 G:121 B:238
	MediumPurple3        = "MediumPurple3"        // R:137 G:104 B:205
	MediumPurple4        = "MediumPurple4"        // R:93 G:71 B:139
	MediumSeaGreen       = "MediumSeaGreen"       // R:60 G:179 B:113
	MediumSlateBlue      = "MediumSlateBlue"      // R:123 G:104 B:238
	MediumSpringGreen    = "MediumSpringGreen"    // R:0 G:250 B:154
	MediumTurquoise      = "MediumTurquoise"      // R:72 G:209 B:204
	MediumVioletRed      = "MediumVioletRed"      // R:199 G:21 B:133
	MidnightBlue         = "MidnightBlue"         // R:25 G:25 B:112
	MintCream            = "MintCream"            // R:245 G:255 B:250
	MistyRose            = "MistyRose"            // R:255 G:228 B:225
	MistyRose1           = "MistyRose1"           // R:255 G:228 B:225
	MistyRose2           = "MistyRose2"           // R:238 G:213 B:210
	MistyRose3           = "MistyRose3"           // R:205 G:183 B:181
	MistyRose4           = "MistyRose4"           // R:139 G:125 B:123
	Moccasin             = "moccasin"             // R:255 G:228 B:181
	NavajoWhite          = "NavajoWhite"          // R:255 G:222 B:173
	NavajoWhite1         = "NavajoWhite1"         // R:255 G:222 B:173
	NavajoWhite2         = "NavajoWhite2"         // R:238 G:207 B:161
	NavajoWhite3         = "NavajoWhite3"         // R:205 G:179 B:139
	NavajoWhite4         = "NavajoWhite4"         // R:139 G:121 B:94
	Navy                 = "navy"                 // R:0 G:0 B:128
	NavyBlue             = "NavyBlue"             // R:0 G:0 B:128
	OldLace              = "OldLace"              // R:253 G:245 B:230
	Olive                = "olive"                // R:128 G:128 B:0
	OliveDrab            = "OliveDrab"            // R:107 G:142 B:35
	OliveDrab1           = "OliveDrab1"           // R:192 G:255 B:62
	OliveDrab2           = "OliveDrab2"           // R:179 G:238 B:58
	OliveDrab3           = "OliveDrab3"           // R:154 G:205 B:50
	OliveDrab4           = "OliveDrab4"           // R:105 G:139 B:34
	Orange               = "orange"               // R:255 G:165 B:0
	Orange1              = "orange1"              // R:255 G:165 B:0
	Orange2              = "orange2"              // R:238 G:154 B:0
	Orange3              = "orange3"              // R:205 G:133 B:0
	Orange4              = "orange4"              // R:139 G:90 B:0
	OrangeRed            = "OrangeRed"            // R:255 G:69 B:0
	OrangeRed1           = "OrangeRed1"           // R:255 G:69 B:0
	OrangeRed2           = "OrangeRed2"           // R:238 G:64 B:0
	OrangeRed3           = "OrangeRed3"           // R:205 G:55 B:0
	OrangeRed4           = "OrangeRed4"           // R:139 G:37 B:0
	Orchid               = "orchid"               // R:218 G:112 B:214
	Orchid1              = "orchid1"              // R:255 G:131 B:250
	Orchid2              = "orchid2"              // R:238 G:122 B:233
	Orchid3              = "orchid3"              // R:205 G:105 B:201
	Orchid4              = "orchid4"              // R:139 G:71 B:137
	PaleGoldenrod        = "PaleGoldenrod"        // R:238 G:232 B:170
	PaleGreen            = "PaleGreen"            // R:152 G:251 B:152
	PaleGreen1           = "PaleGreen1"           // R:154 G:255 B:154
	PaleGreen2           = "PaleGreen2"           // R:144 G:238 B:144
	PaleGreen3           = "PaleGreen3"           // R:124 G:205 B:124
	PaleGreen4           = "PaleGreen4"           // R:84 G:139 B:84
	PaleTurquoise        = "PaleTurquoise"        // R:175 G:238 B:238
	PaleTurquoise1       = "PaleTurquoise1"       // R:187 G:255 B:255
	PaleTurquoise2       = "PaleTurquoise2"       // R:174 G:238 B:238
	PaleTurquoise3       = "PaleTurquoise3"       // R:150 G:205 B:205
	PaleTurquoise4       = "PaleTurquoise4"       // R:102 G:139 B:139
	PaleVioletRed        = "PaleVioletRed"        // R:219 G:112 B:147
	PaleVioletRed1       = "PaleVioletRed1"       // R:255 G:130 B:171
	PaleVioletRed2       = "PaleVioletRed2"       // R:238 G:121 B:159
	PaleVioletRed3       = "PaleVioletRed3"       // R:205 G:104 B:127
	PaleVioletRed4       = "PaleVioletRed4"       // R:139 G:71 B:93
	PapayaWhip           = "PapayaWhip"           // R:255 G:239 B:213
	PeachPuff            = "PeachPuff"            // R:255 G:218 B:185
	PeachPuff1           = "PeachPuff1"           // R:255 G:218 B:185
	PeachPuff2           = "PeachPuff2"           // R:238 G:203 B:173
	PeachPuff3           = "PeachPuff3"           // R:205 G:175 B:149
	PeachPuff4           = "PeachPuff4"           // R:139 G:119 B:101
	Peru                 = "peru"                 // R:205 G:133 B:63
	Pink                 = "pink"                 // R:255 G:192 B:203
	Pink1                = "pink1"                // R:255 G:181 B:197
	Pink2                = "pink2"                // R:238 G:169 B:184
	Pink3                = "pink3"                // R:205 G:145 B:158
	Pink4                = "pink4"                // R:139 G:99 B:108
	Plum                 = "plum"                 // R:221 G:160 B:221
	Plum1                = "plum1"                // R:255 G:187 B:255
	Plum2                = "plum2"                // R:238 G:174 B:238
	Plum3                = "plum3"                // R:205 G:150 B:205
	Plum4                = "plum4"                // R:139 G:102 B:139
	PowderBlue           = "PowderBlue"           // R:176 G:224 B:230
	Purple               = "purple"               // R:128 G:0 B:128
	Purple1              = "purple1"              // R:155 G:48 B:255
	Purple2              = "purple2"              // R:145 G:44 B:238
	Purple3              = "purple3"              // R:125 G:38 B:205
	Purple4              = "purple4"              // R:85 G:26 B:139
	Red                  = "red"                  // R:255 G:0 B:0
	Red1                 = "red1"                 // R:255 G:0 B:0
	Red2                 = "red2"                 // R:238 G:0 B:0
	Red3                 = "red3"                 // R:205 G:0 B:0
	Red4                 = "red4"                 // R:139 G:0 B:0
	RosyBrown            = "RosyBrown"            // R:188 G:143 B:143
	RosyBrown1           = "RosyBrown1"           // R:255 G:193 B:193
	RosyBrown2           = "RosyBrown2"           // R:238 G:180 B:180
	RosyBrown3           = "RosyBrown3"           // R:205 G:155 B:155
	RosyBrown4           = "RosyBrown4"           // R:139 G:105 B:105
	RoyalBlue            = "RoyalBlue"            // R:65 G:105 B:225
	RoyalBlue1           = "RoyalBlue1"           // R:72 G:118 B:255
	RoyalBlue2           = "RoyalBlue2"           // R:67 G:110 B:238
	RoyalBlue3           = "RoyalBlue3"           // R:58 G:95 B:205
	RoyalBlue4           = "RoyalBlue4"           // R:39 G:64 B:139
	SaddleBrown          = "SaddleBrown"          // R:139 G:69 B:19
	Salmon               = "salmon"               // R:250 G:128 B:114
	Salmon1              = "salmon1"              // R:255 G:140 B:105
	Salmon2              = "salmon2"              // R:238 G:130 B:98
	Salmon3              = "salmon3"              // R:205 G:112 B:84
	Salmon4              = "salmon4"              // R:139 G:76 B:57
	SandyBrown           = "SandyBrown"           // R:244 G:164 B:96
	SeaGreen             = "SeaGreen"             // R:46 G:139 B:87
	SeaGreen1            = "SeaGreen1"            // R:84 G:255 B:159
	SeaGreen2            = "SeaGreen2"            // R:78 G:238 B:148
	SeaGreen3            = "SeaGreen3"            // R:67 G:205 B:128
	SeaGreen4            = "SeaGreen4"            // R:46 G:139 B:87
	Seashell             = "seashell"             // R:255 G:245 B:238
	Seashell1            = "seashell1"            // R:255 G:245 B:238
	Seashell2            = "seashell2"            // R:238 G:229 B:222
	Seashell3            = "seashell3"            // R:205 G:197 B:191
	Seashell4            = "seashell4"            // R:139 G:134 B:130
	Sienna               = "sienna"               // R:160 G:82 B:45
	Sienna1              = "sienna1"              // R:255 G:130 B:71
	Sienna2              = "sienna2"              // R:238 G:121 B:66
	Sienna3              = "sienna3"              // R:205 G:104 B:57
	Sienna4              = "sienna4"              // R:139 G:71 B:38
	Silver               = "silver"               // R:192 G:192 B:192
	SkyBlue              = "SkyBlue"              // R:135 G:206 B:235
	SkyBlue1             = "SkyBlue1"             // R:135 G:206 B:255
	SkyBlue2             = "SkyBlue2"             // R:126 G:192 B:238
	SkyBlue3             = "SkyBlue3"             // R:108 G:166 B:205
	SkyBlue4             = "SkyBlue4"             // R:74 G:112 B:139
	SlateBlue            = "SlateBlue"            // R:106 G:90 B:205
	SlateBlue1           = "SlateBlue1"           // R:131 G:111 B:255
	SlateBlue2           = "SlateBlue2"           // R:122 G:103 B:238
	SlateBlue3           = "SlateBlue3"           // R:105 G:89 B:205
	SlateBlue4           = "SlateBlue4"           // R:71 G:60 B:139
	SlateGray            = "SlateGray"            // R:112 G:128 B:144
	SlateGray1           = "SlateGray1"           // R:198 G:226 B:255
	SlateGray2           = "SlateGray2"           // R:185 G:211 B:238
	SlateGray3           = "SlateGray3"           // R:159 G:182 B:205
	SlateGray4           = "SlateGray4"           // R:108 G:123 B:139
	SlateGrey            = "SlateGrey"            // R:112 G:128 B:144
	Snow                 = "snow"                 // R:255 G:250 B:250
	Snow1                = "snow1"                // R:255 G:250 B:250
	Snow2                = "snow2"                // R:238 G:233 B:233
	Snow3                = "snow3"                // R:205 G:201 B:201
	Snow4                = "snow4"                // R:139 G:137 B:137
	SpringGreen          = "SpringGreen"          // R:0 G:255 B:127
	SpringGreen1         = "SpringGreen1"         // R:0 G:255 B:127
	SpringGreen2         = "SpringGreen2"         // R:0 G:238 B:118
	SpringGreen3         = "SpringGreen3"         // R:0 G:205 B:102
	SpringGreen4         = "SpringGreen4"         // R:0 G:139 B:69
	SteelBlue            = "SteelBlue"            // R:70 G:130 B:180
	SteelBlue1           = "SteelBlue1"           // R:99 G:184 B:255
	SteelBlue2           = "SteelBlue2"           // R:92 G:172 B:238
	SteelBlue3           = "SteelBlue3"           // R:79 G:148 B:205
	SteelBlue4           = "SteelBlue4"           // R:54 G:100 B:139
	Tan                  = "tan"                  // R:210 G:180 B:140
	Tan1                 = "tan1"                 // R:255 G:165 B:79
	Tan2                 = "tan2"                 // R:238 G:154 B:73
	Tan3                 = "tan3"                 // R:205 G:133 B:63
	Tan4                 = "tan4"                 // R:139 G:90 B:43
	Teal                 = "teal"                 // R:0 G:128 B:128
	Thistle              = "thistle"              // R:216 G:191 B:216
	Thistle1             = "thistle1"             // R:255 G:225 B:255
	Thistle2             = "thistle2"             // R:238 G:210 B:238
	Thistle3             = "thistle3"             // R:205 G:181 B:205
	Thistle4             = "thistle4"             // R:139 G:123 B:139
	Tomato               = "tomato"               // R:255 G:99 B:71
	Tomato1              = "tomato1"              // R:255 G:99 B:71
	Tomato2              = "tomato2"              // R:238 G:92 B:66
	Tomato3              = "tomato3"              // R:205 G:79 B:57
	Tomato4              = "tomato4"              // R:139 G:54 B:38
	Turquoise            = "turquoise"            // R:64 G:224 B:208
	Turquoise1           = "turquoise1"           // R:0 G:245 B:255
	Turquoise2           = "turquoise2"           // R:0 G:229 B:238
	Turquoise3           = "turquoise3"           // R:0 G:197 B:205
	Turquoise4           = "turquoise4"           // R:0 G:134 B:139
	Violet               = "violet"               // R:238 G:130 B:238
	VioletRed            = "VioletRed"            // R:208 G:32 B:144
	VioletRed1           = "VioletRed1"           // R:255 G:62 B:150
	VioletRed2           = "VioletRed2"           // R:238 G:58 B:140
	VioletRed3           = "VioletRed3"           // R:205 G:50 B:120
	VioletRed4           = "VioletRed4"           // R:139 G:34 B:82
	Wheat                = "wheat"                // R:245 G:222 B:179
	Wheat1               = "wheat1"               // R:255 G:231 B:186
	Wheat2               = "wheat2"               // R:238 G:216 B:174
	Wheat3               = "wheat3"               // R:205 G:186 B:150
	Wheat4               = "wheat4"               // R:139 G:126 B:102
	White                = "white"                // R:255 G:255 B:255
	WhiteSmoke           = "WhiteSmoke"           // R:245 G:245 B:245
	Yellow               = "yellow"               // R:255 G:255 B:0
	Yellow1              = "yellow1"              // R:255 G:255 B:0
	Yellow2              = "yellow2"              // R:238 G:238 B:0
	Yellow3              = "yellow3"              // R:205 G:205 B:0
	Yellow4              = "yellow4"              // R:139 G:139 B:0
	YellowGreen          = "YellowGreen"          // R:154 G:205 B:50
)

// Additional system colors available on macOS.
const (
	SystemActiveAreaFill                   = "systemActiveAreaFill"
	SystemAlertBackgroundActive            = "systemAlertBackgroundActive"
	SystemAlertBackgroundInactive          = "systemAlertBackgroundInactive"
	SystemAlternatePrimaryHighlightColor   = "systemAlternatePrimaryHighlightColor"
	SystemAppleGuideCoachmark              = "systemAppleGuideCoachmark"
	SystemBevelActiveDark                  = "systemBevelActiveDark"
	SystemBevelActiveLight                 = "systemBevelActiveLight"
	SystemBevelInactiveDark                = "systemBevelInactiveDark"
	SystemBevelInactiveLight               = "systemBevelInactiveLight"
	SystemBlack                            = "systemBlack"
	SystemButtonActiveDarkHighlight        = "systemButtonActiveDarkHighlight"
	SystemButtonActiveDarkShadow           = "systemButtonActiveDarkShadow"
	SystemButtonActiveLightHighlight       = "systemButtonActiveLightHighlight"
	SystemButtonActiveLightShadow          = "systemButtonActiveLightShadow"
	SystemButtonFaceActive                 = "systemButtonFaceActive"
	SystemButtonFaceInactive               = "systemButtonFaceInactive"
	SystemButtonFacePressed                = "systemButtonFacePressed"
	SystemButtonFrame                      = "systemButtonFrame"
	SystemButtonFrameActive                = "systemButtonFrameActive"
	SystemButtonFrameInactive              = "systemButtonFrameInactive"
	SystemButtonInactiveDarkHighlight      = "systemButtonInactiveDarkHighlight"
	SystemButtonInactiveDarkShadow         = "systemButtonInactiveDarkShadow"
	SystemButtonInactiveLightHighlight     = "systemButtonInactiveLightHighlight"
	SystemButtonInactiveLightShadow        = "systemButtonInactiveLightShadow"
	SystemButtonPressedDarkHighlight       = "systemButtonPressedDarkHighlight"
	SystemButtonPressedDarkShadow          = "systemButtonPressedDarkShadow"
	SystemButtonPressedLightHighlight      = "systemButtonPressedLightHighlight"
	SystemButtonPressedLightShadow         = "systemButtonPressedLightShadow"
	SystemChasingArrows                    = "systemChasingArrows"
	SystemControlAccentColor               = "systemControlAccentColor"
	SystemControlTextColor                 = "systemControlTextColor"
	SystemDialogBackgroundActive           = "systemDialogBackgroundActive"
	SystemDialogBackgroundInactive         = "systemDialogBackgroundInactive"
	SystemDisabledControlTextColor         = "systemDisabledControlTextColor"
	SystemDocumentWindowBackground         = "systemDocumentWindowBackground"
	SystemDragHilite                       = "systemDragHilite"
	SystemDrawerBackground                 = "systemDrawerBackground"
	SystemFinderWindowBackground           = "systemFinderWindowBackground"
	SystemFocusHighlight                   = "systemFocusHighlight"
	SystemHighlightAlternate               = "systemHighlightAlternate"
	SystemHighlightSecondary               = "systemHighlightSecondary"
	SystemIconLabelBackground              = "systemIconLabelBackground"
	SystemIconLabelBackgroundSelected      = "systemIconLabelBackgroundSelected"
	SystemLabelColor                       = "systemLabelColor"
	SystemLinkColor                        = "systemLinkColor"
	SystemListViewBackground               = "systemListViewBackground"
	SystemListViewColumnDivider            = "systemListViewColumnDivider"
	SystemListViewEvenRowBackground        = "systemListViewEvenRowBackground"
	SystemListViewOddRowBackground         = "systemListViewOddRowBackground"
	SystemListViewSeparator                = "systemListViewSeparator"
	SystemListViewSortColumnBackground     = "systemListViewSortColumnBackground"
	SystemMenuActive                       = "systemMenuActive"
	SystemMenuBackground                   = "systemMenuBackground"
	SystemMenuBackgroundSelected           = "systemMenuBackgroundSelected"
	SystemModelessDialogBackgroundActive   = "systemModelessDialogBackgroundActive"
	SystemModelessDialogBackgroundInactive = "systemModelessDialogBackgroundInactive"
	SystemMovableModalBackground           = "systemMovableModalBackground"
	SystemNotificationWindowBackground     = "systemNotificationWindowBackground"
	SystemPlaceholderTextColor             = "systemPlaceholderTextColor"
	SystemPopupArrowActive                 = "systemPopupArrowActive"
	SystemPopupArrowInactive               = "systemPopupArrowInactive"
	SystemPopupArrowPressed                = "systemPopupArrowPressed"
	SystemPrimaryHighlightColor            = "systemPrimaryHighlightColor"
	SystemScrollBarDelimiterActive         = "systemScrollBarDelimiterActive"
	SystemScrollBarDelimiterInactive       = "systemScrollBarDelimiterInactive"
	SystemSecondaryHighlightColor          = "systemSecondaryHighlightColor"
	SystemSelectedTabTextColor             = "systemSelectedTabTextColor"
	SystemSelectedTextBackgroundColor      = "systemSelectedTextBackgroundColor"
	SystemSelectedTextColor                = "systemSelectedTextColor"
	SystemSeparatorColor                   = "systemSeparatorColor"
	SystemSheetBackground                  = "systemSheetBackground"
	SystemSheetBackgroundOpaque            = "systemSheetBackgroundOpaque"
	SystemSheetBackgroundTransparent       = "systemSheetBackgroundTransparent"
	SystemStaticAreaFill                   = "systemStaticAreaFill"
	SystemTextBackgroundColor              = "systemTextBackgroundColor"
	SystemTextColor                        = "systemTextColor"
	SystemToolbarBackground                = "systemToolbarBackground"
	SystemTransparent                      = "systemTransparent"
	SystemUtilityWindowBackgroundActive    = "systemUtilityWindowBackgroundActive"
	SystemUtilityWindowBackgroundInactive  = "systemUtilityWindowBackgroundInactive"
	SystemWhite                            = "systemWhite"
	SystemWindowBackgroundColor            = "systemWindowBackgroundColor"
	SystemWindowBackgroundColor1           = "systemWindowBackgroundColor1"
	SystemWindowBackgroundColor2           = "systemWindowBackgroundColor2"
	SystemWindowBackgroundColor3           = "systemWindowBackgroundColor3"
	SystemWindowBackgroundColor4           = "systemWindowBackgroundColor4"
	SystemWindowBackgroundColor5           = "systemWindowBackgroundColor5"
	SystemWindowBackgroundColor6           = "systemWindowBackgroundColor6"
	SystemWindowBackgroundColor7           = "systemWindowBackgroundColor7"
	SystemWindowBody                       = "systemWindowBody"
)

// Additional system colors available on Windows.  Note that the actual color
// values depend on the currently active OS theme.
const (
	System3dDarkShadow        = "system3dDarkShadow"
	System3dLight             = "system3dLight"
	SystemActiveBorder        = "systemActiveBorder"
	SystemActiveCaption       = "systemActiveCaption"
	SystemAppWorkspace        = "systemAppWorkspace"
	SystemBackground          = "systemBackground"
	SystemButtonHighlight     = "systemButtonHighlight"
	SystemButtonShadow        = "systemButtonShadow"
	SystemButtonText          = "systemButtonText"
	SystemCaptionText         = "systemCaptionText"
	SystemDisabledText        = "systemDisabledText"
	SystemGrayText            = "systemGrayText"
	SystemHighlightText       = "systemHighlightText"
	SystemInactiveBorder      = "systemInactiveBorder"
	SystemInactiveCaption     = "systemInactiveCaption"
	SystemInactiveCaptionText = "systemInactiveCaptionText"
	SystemInfoBackground      = "systemInfoBackground"
	SystemInfoText            = "systemInfoText"
	SystemMenuText            = "systemMenuText"
	SystemScrollbar           = "systemScrollbar"
	SystemWindow              = "systemWindow"
	SystemWindowFrame         = "systemWindowFrame"
	SystemWindowText          = "systemWindowText"
)

// Additional system colors available both on macOS and Windows.
const (
	SystemButtonFace = "systemButtonFace"
	SystemHighlight  = "systemHighlight"
	SystemMenu       = "systemMenu"
)

// Default generic font names
const (
	DefaultFont      = "TkDefaultFont"      // Default for items not otherwise specified.
	TextFont         = "TkTextFont"         // Used for entry widgets, listboxes, etc.
	FixedFont        = "TkFixedFont"        // A standard fixed-width font.
	MenuFont         = "TkMenuFont"         // The font used for menu items.
	HeadingFont      = "TkHeadingFont"      // Font for column headings in lists and tables.
	CaptionFont      = "TkCaptionFont"      // A font for window and dialog caption bars.
	SmallCaptionFont = "TkSmallCaptionFont" // A smaller caption font for tool dialogs.
	IconFont         = "TkIconFont"         // A font for icon captions.
	TooltipFont      = "TkTooltipFont"      // A font for tooltips.
)

// Common Tk-specific words.
//
// Although Go's style guide recommends MixedCase, these are all in ALL_CAPS
// to avoid namespace collisions. For the pack fill options a prefix is used
// to avoid colliding with the X() and Y() option functions.
//
// See https://gitlab.com/cznic/tk9.0/-/issues/25
const (
	// Font attributes
	NORMAL     = "normal"
	BOLD       = "bold"
	ITALIC     = "italic"
	ROMAN      = "roman"
	UNDERLINE  = "underline"
	OVERSTRIKE = "overstrike"
	// Common font names (for Courier use the CourierFont() function)
	HELVETICA = "helvetica"
	TIMES     = "times"
	// Text justify attributes (also used in other contexts)
	CENTER = "center" // also used as an anchor
	LEFT   = "left"   // also used as a pack side option
	RIGHT  = "right"  // also used as a pack side option
	// Text wrapping
	NONE = "none" // Also used as a Treeview select mode
	CHAR = "char"
	WORD = "word"
	// Text end position
	END = "end"
	// MessageBox icon names
	INFO     = "info"
	QUESTION = "question"
	WARNING  = "warning"
	ERROR    = "error"
	// Anchor and sticky options (CENTER is also an anchor option)
	N    = "n"
	S    = "s"
	W    = "w"
	E    = "e"
	NEWS = "nswe"
	WE   = "we"
	NS   = "ns"
	// Pack fill options
	FILL_X    = "x"
	FILL_Y    = "y"
	FILL_BOTH = "both"
	// Pack side options (can also use LEFT and RIGHT)
	TOP    = "top"
	BOTTOM = "bottom"
	// Orientation (e.g., for TPanedWindow)
	VERTICAL   = "vertical"
	HORIZONTAL = "horizontal"
	// Select mode
	EXTENDED = "extended"
	BROWSE   = "browse"
	// Select type
	CELL = "cell"
	ITEM = "item"
	// Relief
	FLAT   = "flat"
	GROOVE = "groove"
	RAISED = "raised"
	RIDGE  = "ridge"
	SOLID  = "solid"
	SUNKEN = "sunken"
	// Window Manager protocols
	WM_TAKE_FOCUS    = "WM_TAKE_FOCUS"
	WM_DELETE_WINDOW = "WM_DELETE_WINDOW"
)

// Modifier is a bit field representing 0 or more modifiers.
type Modifier int

func (mods Modifier) String() string {
	var names []string
	for _, mod := range modifierNames {
		if mods&mod.modifier == mod.modifier {
			names = append(names, mod.name)
		}
	}
	return strings.Join(names, "+")
}

const (
	ModifierNone  Modifier = 0
	ModifierShift Modifier = 1 << (iota - 1)
	ModifierLock
	ModifierControl
	ModifierMod1
	ModifierMod2
	ModifierMod3
	ModifierMod4
	ModifierMod5
	ModifierButton1
	ModifierButton2
	ModifierButton3
	ModifierButton4
	ModifierButton5

	ModifierAlt     = ModifierMod1
	ModifierNumlock = ModifierMod2
	ModifierWindows = ModifierMod4
	ModifierSuper   = ModifierMod4
)

var modifierNames = []struct {
	modifier Modifier
	name     string
}{
	{ModifierShift, "Shift"},
	{ModifierLock, "Lock"},
	{ModifierControl, "Control"},
	{ModifierMod1, "Mod1"},
	{ModifierMod2, "Mod2"},
	{ModifierMod3, "Mod3"},
	{ModifierMod4, "Mod4"},
	{ModifierMod5, "Mod5"},
	{ModifierButton1, "Button1"},
	{ModifierButton2, "Button2"},
	{ModifierButton3, "Button3"},
	{ModifierButton4, "Button4"},
	{ModifierButton5, "Button5"},
}
