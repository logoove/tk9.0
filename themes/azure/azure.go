// Copyright 2024 The tk9.0-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package azure provides the [Azure Tk theme].
//
// To make the theme available in an application:
//
//	import _ "modernc.org/tk9.0/theme/azure"
//
// See the [modernc.org/tk9.0.Theme] documentation for information about
// activating themes at run time.
//
// [Azure Tk theme]: https://github.com/rdbende/Azure-ttk-theme
package azure // import "modernc.org/tk9.0/theme/azure"

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	. "modernc.org/tk9.0"
)

const (
	// Update when azure.zip file changed.
	version = "v0.1.0"
)

var (
	_ Theme = (*theme)(nil)

	// The azure.zip file is downloaded from
	//
	//	https://github.com/rdbende/Azure-ttk-theme/archive/refs/heads/main.zip.
	//
	// The theme/{light,dark}.tcl files were patched like this:
	//
	//	s/package require Tk 8.6/package require Tk 9.0/
	//
	// No other changes were necessary in this case.
	//
	//go:embed embed/azure.zip
	zip []byte

	initialized bool
)

func init() {
	RegisterTheme("Azure light", newTheme("set_theme light"))
	RegisterTheme("Azure dark", newTheme("set_theme dark"))
}

func setup(context ThemeContext) (err error) {
	defer func() {
		initialized = true
	}()

	if initialized {
		return nil
	}

	root, err := os.UserCacheDir()
	if err != nil {
		return err
	}

	const azure = "azure.zip"
	dir, err := mkzip(filepath.Join(root, "modernc.org", "tk9.0.1", "themes", "azure", version), azure, zip)
	if err != nil {
		return err
	}

	var wd string
	if wd, err = os.Getwd(); err != nil {
		return err
	}

	defer func() {
		err = errors.Join(err, os.Chdir(wd))
	}()

	if err = os.Chdir(dir); err != nil {
		return
	}

	mount := "/themes/azure"
	if _, err = context.Eval(fmt.Sprintf("zipfs mount %s %s", azure, mount)); err != nil {
		return err
	}

	if _, err = context.Eval(fmt.Sprintf("source //zipfs:%s/Azure-ttk-theme-main/azure.tcl", mount)); err != nil {
		return err
	}

	return nil
}

func mkzip(dir, base string, zip []byte) (r string, err error) {
	if _, err = os.Stat(dir); err != nil {
		if !os.IsNotExist(err) {
			return "", err
		}

		if err = os.MkdirAll(dir, 0770); err != nil {
			return "", err
		}
	}

	path := filepath.Join(dir, base)
	if _, err = os.Stat(path); err != nil {
		if !os.IsNotExist(err) {
			return "", err
		}

		return dir, os.WriteFile(path, zip, 0660)
	}

	b, err := os.ReadFile(path)
	if err == nil {
		if bytes.Equal(b, zip) {
			return dir, nil
		}
	}

	os.Remove(path)
	if err := os.WriteFile(path, zip, 0660); err == nil {
		return dir, nil
	}

	dir, err = os.MkdirTemp("", "azure-theme")
	if err != nil {
		return "", err
	}

	path = filepath.Join(dir, base)
	return dir, os.WriteFile(path, zip, 0660)
}

type theme struct {
	activate string
}

func newTheme(activate string) *theme {
	return &theme{
		activate: activate,
	}
}

func (t *theme) Activate(context ThemeContext) (err error) {
	_, err = context.Eval(t.activate)
	return err
}

func (t *theme) Deactivate(context ThemeContext) error {
	// nop
	return nil
}

func (t *theme) Finalize(context ThemeContext) error {
	// nop
	return nil
}

func (t *theme) Initialize(context ThemeContext) error {
	// nop
	return setup(context)
}
