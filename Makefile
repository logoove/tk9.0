# Copyright 2024 The tk9.0-go Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

.PHONY:	all clean edit editor test work w65 lib_win lib_linux lib_darwin lib_freebsd \
	build_all_targets demo

TAR = tcl-core9.0.0-src.tar.gz
URL = http://prdownloads.sourceforge.net/tcl/$(TAR)
TAR2 = tk9.0.0-src.tar.gz
URL2 = http://prdownloads.sourceforge.net/tcl/$(TAR2)
GOOS = $(shell go env GOOS)
GOARCH = $(shell go env GOARCH)
WIN32 = embed/windows/386
WIN64 = embed/windows/amd64
WINARM64 = embed/windows/arm64

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

download:
	@if [ ! -f $(TAR) ]; then wget $(URL) ; fi
	@if [ ! -f $(TAR2) ]; then wget $(URL2) ; fi

edit:
	@if [ -f "Session.vim" ]; then gvim -S & else gvim -p Makefile go.mod builder.json *.go & fi

editor:
	go test -c -o /dev/null
	go build -v  -o /dev/null generator.go
	go run generator.go
	gofmt -l -s -w .
	go build -v  -o /dev/null
	$(shell for f in _examples/*.go ; do go build -o /dev/null $$f ; done)

test:
	go test -vet=off -v -timeout 24h -count=1

work:
	rm -f go.work*
	go work init
	go work use .
	go work use ../cc/v4
	go work use ../ccgo/v3
	go work use ../ccgo/v4
	go work use ../libc
	go work use ../libtcl9.0
	go work use ../libtk9.0
	go work use ../libz
	go work use ../tcl9.0
	go work use ../ngrab

win65:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go test -o /dev/null -c
	rsync \
		-avP \
		--rsync-path='wsl rsync' \
		--exclude .git/ \
		--exclude \*.gz \
		--exclude html/ \
		.  \
		win65:src/modernc.org/tk9.0

lib_win: lib_win64 lib_win32 lib_winarm64

lib_win64: download
	if [ "$(GOOS)" != "linux" ]; then exit 1 ; fi
	if [ "$(GOARCH)" != "amd64" ]; then exit 1 ; fi
	rm -rf ~/tmp/tcl9* ~/tmp/tk9* $(WIN64)
	mkdir -p $(WIN64)
	tar xf $(TAR) -C ~/tmp
	tar xf $(TAR2) -C ~/tmp
	sh -c "cd ~/tmp/tcl9.0.0/win ; ./configure --build=x86_64-linux-gnu --host=x86_64-w64-mingw32"
	make -C ~/tmp/tcl9.0.0/win -j12
	cp -v ~/tmp/tcl9.0.0/win/*.dll $(WIN64)
	sh -c "cd ~/tmp/tk9.0.0/win ; ./configure  --build=x86_64-linux-gnu --host=x86_64-w64-mingw32 --with-tcl=$$HOME/tmp/tcl9.0.0/win"
	make -C ~/tmp/tk9.0.0/win -j12
	cp -v ~/tmp/tk9.0.0/win/tcl9tk90.dll ~/tmp/tk9.0.0/win/libtk9.0.0.zip $(WIN64)
	zip -j $(WIN64)/lib.zip.tmp $(WIN64)/*.dll $(WIN64)/*.zip
	rm -f $(WIN64)/*.dll $(WIN64)/*.zip
	mv $(WIN64)/lib.zip.tmp $(WIN64)/lib.zip

lib_win32: download
	if [ "$(GOOS)" != "linux" ]; then exit 1 ; fi
	if [ "$(GOARCH)" != "amd64" ]; then exit 1 ; fi
	rm -rf ~/tmp/tcl9* ~/tmp/tk9* $(WIN32)
	mkdir -p $(WIN32)
	tar xf $(TAR) -C ~/tmp
	tar xf $(TAR2) -C ~/tmp
	sh -c "cd ~/tmp/tcl9.0.0/win ; ./configure --build=x86_64-linux-gnu --host=i686-w64-mingw32"
	make -C ~/tmp/tcl9.0.0/win -j12
	cp -v ~/tmp/tcl9.0.0/win/*.dll ~/tmp/tcl9.0.0/win/tcl90.dll $(WIN32)
	sh -c "cd ~/tmp/tk9.0.0/win ; ./configure  --build=x86_64-linux-gnu --host=i686-w64-mingw32 --with-tcl=$$HOME/tmp/tcl9.0.0/win"
	make -C ~/tmp/tk9.0.0/win -j12
	cp -v ~/tmp/tk9.0.0/win/tcl9tk90.dll ~/tmp/tk9.0.0/win/libtk9.0.0.zip $(WIN32)
	zip -j $(WIN32)/lib.zip.tmp $(WIN32)/*.dll $(WIN32)/*.zip
	rm -f $(WIN32)/*.dll $(WIN32)/*.zip
	mv $(WIN32)/lib.zip.tmp $(WIN32)/lib.zip

lib_winarm64: download
	if [ "$(GOOS)" != "linux" ]; then exit 1 ; fi
	if [ "$(GOARCH)" != "amd64" ]; then exit 1 ; fi
	rm -rf ~/tmp/tcl9* ~/tmp/tk9* $(WINARM64)
	mkdir -p $(WINARM64)
	tar xf $(TAR) -C ~/tmp
	tar xf $(TAR2) -C ~/tmp
	sh -c "cd ~/tmp/tcl9.0.0/win ; ./configure --build=x86_64-linux-gnu --host=aarch64-w64-mingw32 --enable-64bit=arm64"
	sh -c "cd ~/tmp/tcl9.0.0/win ; sed -i 's/-DHAVE_CPUID=1/-UHAVE_CPUID/g' *"
	make -C ~/tmp/tcl9.0.0/win -j12
	cp -v ~/tmp/tcl9.0.0/win/*.dll $(WINARM64)
	sh -c "cd ~/tmp/tk9.0.0/win ; ./configure --build=x86_64-linux-gnu --host=aarch64-w64-mingw32 --with-tcl=$$HOME/tmp/tcl9.0.0/win  --enable-64bit=arm64"
	make -C ~/tmp/tk9.0.0/win -j12
	cp -v ~/tmp/tk9.0.0/win/tcl9tk90.dll ~/tmp/tk9.0.0/win/libtk9.0.0.zip $(WINARM64)
	zip -j $(WINARM64)/lib.zip.tmp $(WINARM64)/*.dll $(WINARM64)/*.zip
	rm -f $(WINARM64)/*.dll $(WINARM64)/*.zip
	mv $(WINARM64)/lib.zip.tmp $(WINARM64)/lib.zip

lib_linux: download
	if [ "$(GOOS)" != "linux" ]; then exit 1 ; fi
	rm -rf ~/tmp/tcl9* ~/tmp/tk9* embed/$(GOOS)/$(GOARCH)
	mkdir -p embed/$(GOOS)/$(GOARCH)
	tar xf $(TAR) -C ~/tmp
	tar xf $(TAR2) -C ~/tmp
	sh -c "cd ~/tmp/tcl9.0.0/unix ; ./configure --disable-dll-unloading"
	make -C ~/tmp/tcl9.0.0/unix -j2
	cp -v ~/tmp/tcl9.0.0/unix/libtcl9.0.so embed/$(GOOS)/$(GOARCH)
	sh -c "cd ~/tmp/tk9.0.0/unix ; ./configure --with-tcl=$$HOME/tmp/tcl9.0.0/unix"
	make -C ~/tmp/tk9.0.0/unix -j2
	cp -v ~/tmp/tk9.0.0/unix/libtcl9tk9.0.so ~/tmp/tk9.0.0/unix/libtk9.0.0.zip embed/$(GOOS)/$(GOARCH)
	zip -j embed/$(GOOS)/$(GOARCH)/lib.zip.tmp embed/$(GOOS)/$(GOARCH)/*.so embed/$(GOOS)/$(GOARCH)/*.zip
	rm -f embed/$(GOOS)/$(GOARCH)/*.so embed/$(GOOS)/$(GOARCH)/*.zip
	mv embed/$(GOOS)/$(GOARCH)/lib.zip.tmp embed/$(GOOS)/$(GOARCH)/lib.zip

lib_darwin: download
	if [ "$(GOOS)" != "darwin" ]; then exit 1 ; fi
	rm -rf ~/tmp/tcl9* ~/tmp/tk9* embed/$(GOOS)/$(GOARCH)
	mkdir -p embed/$(GOOS)/$(GOARCH)
	tar xf $(TAR) -C ~/tmp
	tar xf $(TAR2) -C ~/tmp
	sh -c "cd ~/tmp/tcl9.0.0/unix ; ./configure"
	make -C ~/tmp/tcl9.0.0/unix -j2
	cp -v ~/tmp/tcl9.0.0/unix/libtcl9.0.dylib embed/$(GOOS)/$(GOARCH)
	sh -c "cd ~/tmp/tk9.0.0/unix ; ./configure --with-tcl=$$HOME/tmp/tcl9.0.0/unix --enable-aqua"
	make -C ~/tmp/tk9.0.0/unix -j2
	cp -v ~/tmp/tk9.0.0/unix/libtcl9tk9.0.dylib ~/tmp/tk9.0.0/unix/libtk9.0.0.zip embed/$(GOOS)/$(GOARCH)
	zip -j embed/$(GOOS)/$(GOARCH)/lib.zip.tmp embed/$(GOOS)/$(GOARCH)/*.dylib embed/$(GOOS)/$(GOARCH)/*.zip
	rm -f embed/$(GOOS)/$(GOARCH)/*.dylib embed/$(GOOS)/$(GOARCH)/*.zip
	mv embed/$(GOOS)/$(GOARCH)/lib.zip.tmp embed/$(GOOS)/$(GOARCH)/lib.zip

# use gmake
lib_freebsd: download
	if [ "$(GOOS)" != "freebsd" ]; then exit 1 ; fi
	rm -rf ~/tmp/tcl9* ~/tmp/tk9* embed/$(GOOS)/$(GOARCH)
	mkdir -p embed/$(GOOS)/$(GOARCH)
	tar xf $(TAR) -C ~/tmp
	tar xf $(TAR2) -C ~/tmp
	sh -c "cd ~/tmp/tcl9.0.0/unix ; ./configure --disable-dll-unloading"
	gmake -C ~/tmp/tcl9.0.0/unix -j2
	cp -v ~/tmp/tcl9.0.0/unix/libtcl9.0.so embed/$(GOOS)/$(GOARCH)
	sh -c "cd ~/tmp/tk9.0.0/unix ; ./configure --with-tcl=$$HOME/tmp/tcl9.0.0/unix"
	gmake -C ~/tmp/tk9.0.0/unix -j2
	cp -v ~/tmp/tk9.0.0/unix/libtcl9tk9.0.so ~/tmp/tk9.0.0/unix/libtk9.0.0.zip embed/$(GOOS)/$(GOARCH)
	zip -j embed/$(GOOS)/$(GOARCH)/lib.zip.tmp embed/$(GOOS)/$(GOARCH)/*.so embed/$(GOOS)/$(GOARCH)/*.zip
	rm -f embed/$(GOOS)/$(GOARCH)/*.so embed/$(GOOS)/$(GOARCH)/*.zip
	mv embed/$(GOOS)/$(GOARCH)/lib.zip.tmp embed/$(GOOS)/$(GOARCH)/lib.zip

demo:
	cd _examples && go run demo.go
