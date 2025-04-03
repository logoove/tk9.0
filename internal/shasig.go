// Copyright 2025 The tk9.0-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

////go:build ignore

package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
)

var (
	goos   = runtime.GOOS
	goarch = runtime.GOARCH
)

func main() {
	flag.StringVar(&goos, "goos", runtime.GOOS, "")
	flag.StringVar(&goarch, "goarch", runtime.GOARCH, "")
	flag.Parse()
	m, err := filepath.Glob(filepath.Join("embed", goos, goarch, "*"))
	if err != nil {
		panic(err.Error())
	}

	fn := fmt.Sprintf("tk_%s_%s.go", goos, goarch)
	src, err := os.ReadFile(fn)
	if err != nil {
		panic(err.Error())
	}

	src = src[:bytes.Index(src, []byte("var shasig"))]
	b := bytes.NewBuffer(src)
	b.WriteString("var shasig = map[string]string{\n")
	sort.Strings(m)
	for _, v := range m {
		out, err := exec.Command("sha256sum", v).CombinedOutput()
		if err != nil {
			panic(err.Error())
		}

		sig := strings.Fields(string(out))[0]
		fmt.Fprintf(b, "\t%q: %q,\n", filepath.Base(v), sig)
	}
	b.WriteString("}\n")
	if err := os.WriteFile(fn, b.Bytes(), 0660); err != nil {
		panic(err.Error())
	}
}
