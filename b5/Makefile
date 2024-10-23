# Copyright 2024 The tk9.0-go Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

.PHONY:	all clean edit editor work

all: editor

clean:
	go clean
	rm -f *~ *.test *.out

edit:
	gvim -p Makefile *.go &

editor:
	gofmt -l -s -w *.go
	go test -c -o /dev/null

work:
	rm -f go.work*
	go work init
	go work use . ~/src/modernc.org/tk9.0
