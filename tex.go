// The original code in this file comes from
//
//	https://git.sr.ht/~sbinet/star-tex/tree/main/item/cmd/dvi-cnv
//
// and is
//
// Copyright ©2021 The star-tex Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE-STAR-TEX file.
//
// Modifications are
//
// Copyright 2024 The tk9.0-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tk9_0 // import "modernc.org/tk9.0"

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"strings"

	"github.com/disintegration/imaging"
	"golang.org/x/image/font"
	fixedmath "golang.org/x/image/math/fixed"
	"modernc.org/knuth/dvi"
	"modernc.org/knuth/font/fixed"
	"modernc.org/knuth/font/pkf"
	"modernc.org/knuth/kpath"
	"modernc.org/knuth/tex"
)

const (
	shrink = 1
)

var (
	_ dvi.Renderer = (*renderer)(nil)
)

type fntkey struct {
	name string
	size fixed.Int12_20
}

type renderer struct {
	bkg   color.Color
	bound image.Rectangle
	ctx   kpath.Context
	err   error
	faces map[fntkey]font.Face
	final image.Image
	img   *image.RGBA
	page  int
	post  dvi.CmdPost
	pre   dvi.CmdPre
	scale float64

	conv  float32 // converts DVI units to pixels
	dpi   float32
	tconv float32 // converts unmagnified DVI units to pixels

	bounded bool
}

func newRenderer(ctx kpath.Context, scale float64) *renderer {
	return &renderer{ctx: ctx, faces: make(map[fntkey]font.Face), scale: scale}
}

func (pr *renderer) Init(pre *dvi.CmdPre, post *dvi.CmdPost) {
	pr.pre = *pre
	pr.post = *post
	if pr.dpi == 0 {
		pr.dpi = 600
	}
	res := pr.dpi
	conv := float32(pr.pre.Num) / 254000.0 * (res / float32(pr.pre.Den))
	pr.tconv = conv
	pr.conv = conv * float32(pr.pre.Mag) / 1000.0
	conv = 1/(float32(pre.Num)/float32(pre.Den)*(float32(pre.Mag)/1000.0)*(pr.dpi*shrink/254000.0)) + 0.5
	if pr.bkg == nil {
		pr.bkg = color.Transparent
	}
}

func (pr *renderer) BOP(bop *dvi.CmdBOP) {
	if pr.err != nil {
		return
	}

	pr.page = int(bop.C0)
	bnd := image.Rect(0, 0, int(pr.pixels(int32(pr.post.Width))), int(pr.pixels(int32(pr.post.Height))))
	pr.img = image.NewRGBA(bnd)
	draw.Draw(pr.img, bnd, image.NewUniform(pr.bkg), image.Point{}, draw.Over)
}

func (pr *renderer) DrawGlyph(x, y int32, font dvi.Font, glyph rune, c color.Color) {
	if pr.err != nil {
		return
	}

	dot := fixedmath.Point26_6{X: fixedmath.I(int(pr.pixels(x))), Y: fixedmath.I(int(pr.pixels(y)))}
	face, ok := pr.face(font)
	if !ok {
		return
	}

	dr, mask, maskp, _, ok := face.Glyph(dot, glyph)
	if !ok {
		pr.setErr(fmt.Errorf("could not find glyph 0x%02x", glyph))
		return
	}

	draw.DrawMask(pr.img, dr, image.NewUniform(c), image.Point{}, mask, maskp, draw.Over)
	pr.union(dr)
}

func (pr *renderer) union(r image.Rectangle) {
	switch pr.bounded {
	case true:
		pr.bound = pr.bound.Union(r)
	default:
		pr.bound = r
		pr.bounded = true
	}
}

func (pr *renderer) DrawRule(x, y, w, h int32, c color.Color) {
	if pr.err != nil {
		return
	}

	r := image.Rect(int(pr.pixels(x)), int(pr.pixels(y)), int(pr.pixels(x+w)), int(pr.pixels(y-h)))
	draw.Draw(pr.img, r, image.NewUniform(c), image.Point{}, draw.Over)
	pr.union(r)
}

func (pr *renderer) EOP() {
	if pr.err != nil {
		return
	}

	pr.final = pr.img.SubImage(pr.bound)
	if pr.scale != 1.0 {
		pr.final = imaging.Resize(pr.final, int(float64(pr.bound.Max.X-pr.bound.Min.X)*pr.scale+0.5), 0, imaging.Lanczos)
	}
}

func (pr *renderer) setErr(err error) {
	if pr.err == nil {
		pr.err = err
	}
}

func (pr *renderer) face(fnt dvi.Font) (font.Face, bool) {
	key := fntkey{
		name: fnt.Name(),
		size: fnt.Size(),
	}
	if f, ok := pr.faces[key]; ok {
		return f, ok
	}

	fname, err := pr.ctx.Find(key.name + ".pk")
	if err != nil {
		pr.setErr(fmt.Errorf("could not find font face %q: %+v", key.name, err))
		return nil, false
	}

	f, err := pr.ctx.Open(fname)
	if err != nil {
		pr.setErr(fmt.Errorf("could not open font face %q: %+v", key.name, err))
		return nil, false
	}

	defer f.Close()

	pk, err := pkf.Parse(f)
	if err != nil {
		pr.setErr(fmt.Errorf("could not parse font face %q: %+v", key.name, err))
		return nil, false
	}

	tfm := fnt.Metrics()
	if tfm.Checksum() != pk.Checksum() {
		pr.setErr(fmt.Errorf(
			"TFM and PK checksum do not match for %q: tfm=0x%x, pk=0x%x",
			key.name,
			tfm.Checksum(),
			pk.Checksum(),
		))
		return nil, false
	}

	face := pkf.NewFace(pk, tfm, &pkf.FaceOptions{
		Size: tfm.DesignSize().Float64(),
		DPI:  float64(pr.dpi),
	})
	pr.faces[key] = face
	return face, true
}

func (pr *renderer) pixels(v int32) int32 {
	x := pr.conv * float32(v)
	return roundF32(x / shrink)
}

func roundF32(v float32) int32 {
	if v > 0 {
		return int32(v + 0.5)
	}
	return int32(v - 0.5)
}

func dvi2png(r io.Reader, scale float64) (_ []byte, err error) {
	img, err := dvi2img(r, scale)
	if err != nil {
		return nil, err
	}

	var out bytes.Buffer
	if err = png.Encode(&out, img); err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}

func dvi2img(r io.Reader, scale float64) (img image.Image, err error) {
	ctx := kpath.New()
	renderer := newRenderer(ctx, scale)
	vm := dvi.NewMachine(
		dvi.WithContext(ctx),
		dvi.WithRenderer(renderer),
		dvi.WithHandlers(dvi.NewColorHandler(ctx)),
		dvi.WithOffsetX(0),
		dvi.WithOffsetY(0),
	)

	raw, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("could not read DVI program file: %w", err)
	}

	prog, err := dvi.Compile(raw)
	if err != nil {
		return nil, fmt.Errorf("could not compile DVI program: %w", err)
	}

	err = vm.Run(prog)
	if err != nil {
		return nil, fmt.Errorf("could not interpret DVI program: %w", err)
	}

	if renderer.err != nil {
		return nil, fmt.Errorf("could not render DVI program: %w", renderer.err)
	}

	return renderer.final, nil
}

func tex2dvi(src string) (dvi *bytes.Buffer, err error) {
	// To get rid of the page number rendered by default, the function prepends
	// "\footline={}\n" to src. Also, "\n\bye\n" is appended to 'src' to make it a
	// complete TeX document.
	var stdout, stderr, b bytes.Buffer
	nm := wmTitle
	switch {
	case nm == "plain":
		nm += "_"
	case nm == "":
		nm = "x"
	}
	if err = tex.Main(
		strings.NewReader(fmt.Sprintf("\\input plain \\input %s", nm)),
		&stdout,
		&stderr,
		tex.WithInputFile(nm+".tex", strings.NewReader(fmt.Sprintf("\\footline={}\n%s\n\\bye\n", src))),
		tex.WithDVIFile(&b),
		tex.WithLogFile(io.Discard),
	); err != nil {
		a := []string{err.Error()}
		if b := stdout.Bytes(); len(b) != 0 {
			a = append(a, string(b))
		}
		if b := stderr.Bytes(); len(b) != 0 {
			a = append(a, string(b))
		}
		return nil, fmt.Errorf("%s", strings.Join(a, "\n"))
	}

	return &b, nil
}

func sanitizeTeX(s string) (r string) {
	s = strings.TrimSpace(s)
	a := strings.Fields(s)
	return strings.Join(a, " ")
}

// TeX is like Tex2 but report errors using [ErrorMode] and [Error].
func TeX(src string, scale float64) (png []byte) {
	b, err := TeX2(src, scale)
	if err != nil {
		fail(err)
		return nil
	}

	return b
}

// TeX2 renders TeX 'src' as a png file that shows the TeX "snippet" in a fixed
// 600 dpi resolution. The result is afterwards resized using the 'scale'
// factor. Scale factor 1.0 means no resize.
//
// Only plain Tex and a subset of some of the default Computer Modern fonts are
// supported. Many small fonts are not available.
func TeX2(src string, scale float64) (png []byte, err error) {
	if src = sanitizeTeX(src); src == "" {
		return nil, fmt.Errorf("empty TeX code")
	}

	dvi, err := tex2dvi(src)
	if err != nil {
		return nil, err
	}

	return dvi2png(dvi, scale)
}

// TeXImg renders is line TeX but returns an [image.Image].
func TeXImg(src string, scale float64) (img image.Image) {
	img, err := TeXImg2(src, scale)
	if err != nil {
		fail(err)
		return nil
	}

	return img
}

// TeXImg2 renders is line TeX2 but returns an [image.Image].
func TeXImg2(src string, scale float64) (img image.Image, err error) {
	if src = sanitizeTeX(src); src == "" {
		return nil, fmt.Errorf("empty TeX code")
	}

	dvi, err := tex2dvi(src)
	if err != nil {
		return nil, err
	}

	return dvi2img(dvi, scale)
}
