# Copyright 2024 The tk9.0-go Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

.PHONY:	all clean edit editor test watch work build_all_targets

all:
	golint 2>&1
	staticcheck 2>&1

build_all_targets:
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go test -o /dev/null -c
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go test -o /dev/null -c
	GOOS=freebsd GOARCH=amd64 CGO_ENABLED=0 go build -gcflags="github.com/ebitengine/purego/internal/fakecgo=-std"
	GOOS=freebsd GOARCH=amd64 CGO_ENABLED=0 go test -o /dev/null -c -gcflags="github.com/ebitengine/purego/internal/fakecgo=-std"
	GOOS=freebsd GOARCH=arm64 CGO_ENABLED=0 go build -gcflags="github.com/ebitengine/purego/internal/fakecgo=-std"
	GOOS=freebsd GOARCH=arm64 CGO_ENABLED=0 go test -o /dev/null -c -gcflags="github.com/ebitengine/purego/internal/fakecgo=-std"
	GOOS=linux GOARCH=386 CGO_ENABLED=0 go build
	GOOS=linux GOARCH=386 CGO_ENABLED=0 go test -o /dev/null -c
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go test -o /dev/null -c
	GOOS=linux GOARCH=arm CGO_ENABLED=0 go build
	GOOS=linux GOARCH=arm CGO_ENABLED=0 go test -o /dev/null -c
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go test -o /dev/null -c
	GOOS=linux GOARCH=loong64 CGO_ENABLED=0 go build
	GOOS=linux GOARCH=loong64 CGO_ENABLED=0 go test -o /dev/null -c
	GOOS=linux GOARCH=ppc64le CGO_ENABLED=0 go build
	GOOS=linux GOARCH=ppc64le CGO_ENABLED=0 go test -o /dev/null -c
	GOOS=linux GOARCH=riscv64 CGO_ENABLED=0 go build
	GOOS=linux GOARCH=riscv64 CGO_ENABLED=0 go test -o /dev/null -c
	GOOS=linux GOARCH=s390x CGO_ENABLED=0 go build
	GOOS=linux GOARCH=s390x CGO_ENABLED=0 go test -o /dev/null -c
	GOOS=windows GOARCH=386 CGO_ENABLED=0 go build
	GOOS=windows GOARCH=386 CGO_ENABLED=0 go test -o /dev/null -c
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go test -o /dev/null -c
	GOOS=windows GOARCH=arm64 CGO_ENABLED=0 go build
	GOOS=windows GOARCH=arm64 CGO_ENABLED=0 go test -o /dev/null -c

clean:
	rm -f log-* cpu.test mem.test *.out go.work*
	go clean

edit:
	@if [ -f "Session.vim" ]; then gvim -S & else gvim -p Makefile *.go & fi

editor:
	go test -c -o /dev/null
	gofmt -l -s -w .
	go build -v  -o /dev/null

test:
	go test -vet=off -v -timeout 24h -count=1

watch:
	watch 'ls -la /tmp/.X*-lock ; pgrep -l vfb\|vnc\|websock'

work:
	rm -f go.work*
	go work init
	go work use .
