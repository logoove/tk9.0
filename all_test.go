// Copyright 2024 The tk9.0-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tk9_0 // import "modernc.org/tk9.0"

import (
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"syscall"
	"testing"
	"time"

	_ "github.com/adrg/xdg"       // generator.go
	_ "github.com/expr-lang/expr" // examples
	_ "golang.org/x/net/html"     // generator.go
	_ "modernc.org/ngrab/lib"     // generator.go
	_ "modernc.org/rec/lib"       // generator.go
)

const (
	xvfbDisplayVar = "XVFB_DISPLAY"
)

var (
	display = os.Getenv("DISPLAY")
	re      *regexp.Regexp
)

func TestMain(m *testing.M) {
	if Error != nil {
		fmt.Fprintln(os.Stderr, Error)
		os.Exit(1)
	}

	if display == "" && (goos == "linux" || goos == "freebsd") {
		if s := os.Getenv(xvfbDisplayVar); s != "" {
			display = s
			os.Setenv("DISPLAY", s)
		}
	}
	oRe := flag.String("re", "", "")
	flag.Parse()
	if *oRe != "" {
		re = regexp.MustCompile(*oRe)
	}
	rc := m.Run()
	Finalize()
	os.Exit(rc)
}

func sys(arg0 string, args ...string) (r []byte, err error) {
	return exec.Command(arg0, args...).CombinedOutput()
}

// Commit 18c4e94e171d4 diff
//
// -               `([^$]|\\\$)*`,           // Not TeX, incl. "\$"
// -               `\$([^$]|\\\$)*\$`,       // $TeX$ or $$TeX$$, incl. $Te\$X$
// -               `\$\$?([^$]|\\\$)*\$\$?`, // $TeX$ or $$TeX$$, incl. $Te\$X$
// +               `([^$]|\\\$)*`,         // Not TeX, incl. "\$"
// +               `\$([^$]|\\\$)*\$`,     // $TeX$, incl. $Te\$X$
// +               `\$\$([^$]|\\\$)*\$\$`, // $$TeX$$, incl. $Te\$X$
func TestTokenizer(t *testing.T) {
	for i, test := range []struct {
		s    string
		ids  []int
		toks []string
	}{
		{},
		{"a", []int{0}, []string{"a"}},
		{"\\$", []int{0}, []string{"\\$"}},
		{"\\$\\$", []int{0}, []string{"\\$\\$"}},
		{"\\$\\$\\$", []int{0}, []string{"\\$\\$\\$"}},

		{"\\$\\$\\$\\$", []int{0}, []string{"\\$\\$\\$\\$"}},
		{"a\\$", []int{0}, []string{"a\\$"}},
		{"a\\$\\$", []int{0}, []string{"a\\$\\$"}},
		{"a\\$\\$\\$", []int{0}, []string{"a\\$\\$\\$"}},
		{"a\\$\\$\\$\\$", []int{0}, []string{"a\\$\\$\\$\\$"}},

		{"$a$", []int{1}, []string{"$a$"}},
		// Not valid since 18c4e94e171d4 {"$$a$", []int{2}, []string{"$$a$"}},
		{"$$a$$", []int{2}, []string{"$$a$$"}},
		// Not valid since 18c4e94e171d4 {"$a$$", []int{2}, []string{"$a$$"}},
		{"x$a$", []int{0, 1}, []string{"x", "$a$"}},

		// Not valid since 18c4e94e171d4 {"x$$a$", []int{0, 2}, []string{"x", "$$a$"}},
		{"x$$a$$", []int{0, 2}, []string{"x", "$$a$$"}},
		// Not valid since 18c4e94e171d4 {"x$a$$", []int{0, 2}, []string{"x", "$a$$"}},
		{"x$a$y", []int{0, 1, 0}, []string{"x", "$a$", "y"}},
		// Not valid since 18c4e94e171d4 {"x$$a$y", []int{0, 2, 0}, []string{"x", "$$a$", "y"}},

		{"x$$a$$y", []int{0, 2, 0}, []string{"x", "$$a$$", "y"}},
		// Not valid since 18c4e94e171d4 {"x$a$$y", []int{0, 2, 0}, []string{"x", "$a$$", "y"}},
		// Not valid since 18c4e94e171d4 {"x\\$0$a\\$1b$$\\$y", []int{0, 2, 0}, []string{"x\\$0", "$a\\$1b$$", "\\$y"}},
	} {
		ids, toks := tokenize(test.s)
		if g, e := fmt.Sprintf("%v %q", ids, toks), fmt.Sprintf("%v %q", test.ids, test.toks); g != e {
			t.Errorf("#%3v: `%s`\ngot %s\nexp %s", i, test.s, g, e)
		}
	}
}

// Credits: https://gitlab.com/cznic/tk9.0/-/issues/51#note_2374472931
func TestParseList(t *testing.T) {
	switch target {
	case "linux/s390x":
		t.Skipf("this test is known to not work on %s VM", target)
	}

	Initialize()
	tests := []struct {
		name     string
		in       string
		expected []string
	}{
		{"empty", "", []string{}},
		{"one item", "abc", []string{"abc"}},
		{"multiple items", "abc def ghi", []string{"abc", "def", "ghi"}},
		{"multiple inter item spaces", "abc   def", []string{"abc", "def"}},
		{"leading and trailing spaces", "  abc def  ", []string{"abc", "def"}},
		{"delimited item at start", "{ab c} def ghi", []string{"ab c", "def", "ghi"}},
		{"delimited item in middle", "abc {de f} ghi", []string{"abc", "de f", "ghi"}},
		{"delimited item at end", "abc def {gh i}", []string{"abc", "def", "gh i"}},
		{"all items delimited", "{abc} {def} {ghi}", []string{"abc", "def", "ghi"}},
		{"delimited with leading and trailing space", " {ab c} def {gh i}  ", []string{"ab c", "def", "gh i"}},
		{"whitespace in items", "{ab c} {de\tf} {gh\ni}", []string{"ab c", "de\tf", "gh\ni"}},
		{"braces in items", `ab\{c de\}f`, []string{`ab{c`, "de}f"}},
		{"backslash not escaping a brace", `{ab\c}`, []string{"ab\\c"}},
		{"whitespace in items", "{ab c} {de\tf} {gh\ni}", []string{"ab c", "de\tf", "gh\ni"}},
		{"braces in elements", "a{b c}d e{f} g{{h i{}j k}{l }m}}", []string{"a{b", "c}d", "e{f}", "g{{h", "i{}j", "k}{l", "}m}}"}},
		{"nested list", "{abc {def ghi}} jkl", []string{"abc {def ghi}", "jkl"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			list := parseList(test.in)
			if slices.Compare(list, test.expected) != 0 {
				t.Errorf("got %#v, expected %#v", list, test.expected)
			}
		})
	}
}

func TestExamples(t *testing.T) {
	if !isBuilder {
		t.Skip("not a builder")
	}

	blacklist := map[string]struct{}{
		"demo.go":        {}, // executes multiple other examples
		"fontmetrics.go": {}, // non GUI example
		"ring.go":        {}, // expects arguments
	}

	t.Logf("DISPLAY=%s XVFB_DISPLAY=%s display=%s", os.Getenv("DISPLAY"), os.Getenv(xvfbDisplayVar), display)
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	wd, err = filepath.Abs(wd)
	if err != nil {
		t.Fatal(err)
	}

	wd = filepath.Join(wd, "_examples")
	const retries = 1
	switch goos {
	case "linux", "freebsd":
		if display == "" {
			t.Fatal("DISPLAY=")
		}
	case "windows":
		blacklist["dialog.go"] = struct{}{} // uses X11 specific stuff
	}
	switch target {
	case "linux/s390x":
		blacklist["font.go"] = struct{}{}          // Looks like a qemu issue.
		blacklist["winfoChildren.go"] = struct{}{} // Looks like a qemu issue.
	case "windows/amd64":
		blacklist["photo_gif.go"] = struct{}{}  // See #66
		blacklist["photo_gif2.go"] = struct{}{} // See #66
		blacklist["tablelist.go"] = struct{}{}  // See #66
	case "windows/arm64":
		blacklist["splot.go"] = struct{}{}       // No gnuplot on this builder.
		blacklist["tori.go"] = struct{}{}        // No gnuplot on this builder.
		blacklist["tori_canvas.go"] = struct{}{} // No gnuplot on this builder.
	case "windows/386":
		blacklist["widgetproxy.go"] = struct{}{} // See #54
		blacklist["splot.go"] = struct{}{}       // No gnuplot on this builder.
		blacklist["tori.go"] = struct{}{}        // No gnuplot on this builder.
		blacklist["tori_canvas.go"] = struct{}{} // No gnuplot on this builder.
	}

	tmpDir := t.TempDir()
	m, err := filepath.Glob(filepath.Join("_examples", "*.go"))
	if err != nil {
		t.Fatal(err)
	}

next:
	for i, v := range m {
		if re != nil && !re.MatchString(v) {
			continue
		}

		base := filepath.Base(v)
		if _, ok := blacklist[base]; ok {
			t.Logf("SKIP %v (%v/%v)", v, i+1, len(m))
			continue
		}

		bin := filepath.Join(tmpDir, fmt.Sprintf("prog%v", i))
		if goos == "windows" {
			bin += ".exe"
		}

		if _, err := sys("go", "build", "-o", bin, v); err != nil {
			t.Error(err)
			continue
		}

		var j int
		for j = 0; j < retries; j++ {
			if err = testExample(t, wd, tmpDir, bin); err == nil {
				if testing.Verbose() {
					t.Logf("PASS %v (%v/%v)", v, i+1, len(m))
				}
				continue next
			}
		}

		t.Errorf("%v: FAIL %v (tries=%v)", v, err, j)
	}
}

func testExample(t *testing.T, wd, tmpDir, bin string) (err error) {
	token := fmt.Sprint(time.Now().UnixNano())
	os.Setenv(testHookWaitVar, token)

	var stdout, stderr bytes.Buffer
	cmd := exec.Command(bin)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Dir = wd
	t0 := time.Now()
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	crashCheckDuration := 30 * time.Second
	crashCheckTimer := time.NewTimer(crashCheckDuration)

	waitChan := make(chan error, 1) // Buffered channel to prevent goroutine leak

	go func() {
		waitChan <- cmd.Wait()
	}()

	select {
	case err := <-waitChan:
		crashCheckTimer.Stop() // Stop timer if process exits before timeout.
		if err != nil {
			if _, ok := err.(*exec.ExitError); ok {
				return fmt.Errorf("process crashed within the check duration (after %s).\nstdout=%s\nstderr=%s", time.Since(t0), stdout.Bytes(), stderr.Bytes())
			} else {
				return fmt.Errorf("process exited with error: %v", err)
			}
		} else {
			return fmt.Errorf("process unexpectedly exited normally")
		}
	case <-crashCheckTimer.C:
		if goos == "windows" {
			if err := cmd.Process.Kill(); err != nil {
				return fmt.Errorf("error killing process: %v", err)
			}
		} else {
			if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
				if err := cmd.Process.Kill(); err != nil {
					return fmt.Errorf("error killing process: %v", err)
				}
			}
		}
		err := <-waitChan // Wait for the process to fully exit after the kill signal.
		if err != nil {
			if exitError, ok := err.(*exec.ExitError); ok {
				status := exitError.Sys().(syscall.WaitStatus)
				if !status.Signaled() && goos == "linux" {
					return fmt.Errorf("process exited with error: %v", err)
				}
			} else {
				return fmt.Errorf("process exited with error: %v", err)
			}
		}
	}

	if g, ge, e := strings.TrimSpace(string(stdout.Bytes())), strings.TrimSpace(string(stderr.Bytes())), token; g != e {
		return fmt.Errorf("expected process stdout=%s, got stdout=%s, stderr=%s", e, g, ge)
	}

	return nil
}
