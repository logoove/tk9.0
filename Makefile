# Copyright 2024 The tk9.0-go Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

.PHONY:	all clean edit editor test work w65 lib_win lib_linux_ccgo lib_linux_purego \
	lib_darwin lib_freebsd build_all_targets demo examples xvfb

TCL_TAR = tcl-core9.0.1-src.tar.gz
TCL_TAR_URL = http://prdownloads.sourceforge.net/tcl/$(TCL_TAR)
TK_TAR = tk9.0.1-src.tar.gz
TK_TAR_URL = http://prdownloads.sourceforge.net/tcl/$(TK_TAR)
TK_IMG_TAR = Img-2.0.1.tar.gz
TK_IMG_URL = http://prdownloads.sourceforge.net/tkimg/$(TK_IMG_TAR)
GOOS = $(shell go env GOOS)
GOARCH = $(shell go env GOARCH)
WIN32 = embed/windows/386
WIN64 = embed/windows/amd64
WINARM64 = embed/windows/arm64
GOMAXPROCS = $(shell go run internal/cpus.go 2>&1)
PWD = $(shell pwd)

all:
	golint 2>&1
	staticcheck 2>&1
	$(shell for f in _examples/*.go ; do go build -o /dev/null $$f ; done)

build_all_targets:
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build ./...
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go test -o /dev/null -c
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build ./...
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go test -o /dev/null -c
	GOOS=freebsd GOARCH=amd64 CGO_ENABLED=0 go build -gcflags="github.com/ebitengine/purego/internal/fakecgo=-std" ./...
	GOOS=freebsd GOARCH=amd64 CGO_ENABLED=0 go test -o /dev/null -c -gcflags="github.com/ebitengine/purego/internal/fakecgo=-std"
	GOOS=freebsd GOARCH=arm64 CGO_ENABLED=0 go build -gcflags="github.com/ebitengine/purego/internal/fakecgo=-std" ./...
	GOOS=freebsd GOARCH=arm64 CGO_ENABLED=0 go test -o /dev/null -c -gcflags="github.com/ebitengine/purego/internal/fakecgo=-std"
	GOOS=linux GOARCH=386 CGO_ENABLED=0 go build ./...
	GOOS=linux GOARCH=386 CGO_ENABLED=0 go test -o /dev/null -c
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build ./...
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go test -o /dev/null -c
	GOOS=linux GOARCH=arm CGO_ENABLED=0 go build ./...
	GOOS=linux GOARCH=arm CGO_ENABLED=0 go test -o /dev/null -c
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build ./...
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go test -o /dev/null -c
	GOOS=linux GOARCH=loong64 CGO_ENABLED=0 go build ./...
	GOOS=linux GOARCH=loong64 CGO_ENABLED=0 go test -o /dev/null -c
	GOOS=linux GOARCH=ppc64le CGO_ENABLED=0 go build ./...
	GOOS=linux GOARCH=ppc64le CGO_ENABLED=0 go test -o /dev/null -c
	GOOS=linux GOARCH=riscv64 CGO_ENABLED=0 go build ./...
	GOOS=linux GOARCH=riscv64 CGO_ENABLED=0 go test -o /dev/null -c
	GOOS=linux GOARCH=s390x CGO_ENABLED=0 go build ./...
	GOOS=linux GOARCH=s390x CGO_ENABLED=0 go test -o /dev/null -c
	GOOS=windows GOARCH=386 CGO_ENABLED=0 go build ./...
	GOOS=windows GOARCH=386 CGO_ENABLED=0 go test -o /dev/null -c
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build ./...
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go test -o /dev/null -c
	GOOS=windows GOARCH=arm64 CGO_ENABLED=0 go build ./...
	GOOS=windows GOARCH=arm64 CGO_ENABLED=0 go test -o /dev/null -c

clean:
	rm -f log-* cpu.test mem.test *.out go.work*
	go clean

download:
	@if [ ! -f $(TCL_TAR) ]; then wget $(TCL_TAR_URL) ; fi
	@if [ ! -f $(TK_TAR) ]; then wget $(TK_TAR_URL) ; fi
	@if [ ! -f $(TK_IMG_TAR) ]; then wget $(TK_IMG_URL) ; fi

edit:
	@if [ -f "Session.vim" ]; then gvim -S & else gvim -p Makefile go.mod builder.json *.go & fi

editor:
	go test -vet=off -c -o /dev/null
	go build -v  -o /dev/null generator.go
	@go run generator.go > /dev/null
	gofmt -l -s -w .
	go build -v  -o /dev/null
	go build -v  -o /dev/null ./vnc
	go build -v  -o /dev/null ./themes/azure

test:
	go test -vet=off -v -timeout 24h -count=1

work:
	rm -f go.work*
	go work init
	go work use .

lib_win: lib_win64 lib_win32 lib_winarm64
	git status

lib_win64: download
	if [ "$(GOOS)" != "linux" ]; then exit 1 ; fi
	if [ "$(GOARCH)" != "amd64" ]; then exit 1 ; fi
	rm -rf tcl9.0.1/ tk9.0.1/ Img-2.0.1/ $(WIN64)
	mkdir -p $(WIN64)
	tar xf $(TCL_TAR)
	tar xf $(TK_TAR)
	tar xf $(TK_IMG_TAR)
	patch Img-2.0.1/tiff/tiff.c internal/tiff.c.patch
	sh -c "cd tcl9.0.1/win ; ./configure --build=x86_64-linux-gnu --host=x86_64-w64-mingw32"
	make -C tcl9.0.1/win -j$(GOMAXPROCS)
	cp -v tcl9.0.1/win/*.dll $(WIN64)
	sh -c "cd tk9.0.1/win ; ./configure --build=x86_64-linux-gnu --host=x86_64-w64-mingw32 --with-tcl=$(PWD)/tcl9.0.1/win"
	make -C tk9.0.1/win -j$(GOMAXPROCS)
	cp -v tk9.0.1/win/tcl9tk90.dll tk9.0.1/win/libtk9.0.1.zip $(WIN64)
	sh -c "cd Img-2.0.1 ; ./configure --build=x86_64-linux-gnu --host=x86_64-w64-mingw32 --with-tcl=$(PWD)/tcl9.0.1/win  --with-tk=$(PWD)/tk9.0.1/win"
	make -C Img-2.0.1 -j$(GOMAXPROCS)
	find Img-2.0.1 -name \*.dll -exec cp {} $(WIN64) \;
	go run internal/shasig.go -goos=windows -goarch=amd64 tk_windows_amd64.go
	gofmt -l -s -w tk_windows_amd64.go
	zip -j $(WIN64)/lib.zip.tmp $(WIN64)/*.dll $(WIN64)/*.zip
	rm -f $(WIN64)/*.dll $(WIN64)/*.zip
	mv $(WIN64)/lib.zip.tmp $(WIN64)/lib.zip
	rm -rf tcl9.0.1/ tk9.0.1/ Img-2.0.1/
	GOOS=windows GOARCH=amd64 go build -v
	GOOS=windows GOARCH=amd64 ./unconvert.sh
	git status

lib_win32: download
	if [ "$(GOOS)" != "linux" ]; then exit 1 ; fi
	if [ "$(GOARCH)" != "amd64" ]; then exit 1 ; fi
	rm -rf tcl9.0.1/ tk9.0.1/ Img-2.0.1/ $(WIN32)
	mkdir -p $(WIN32)
	tar xf $(TCL_TAR)
	tar xf $(TK_TAR)
	tar xf $(TK_IMG_TAR)
	patch Img-2.0.1/tiff/tiff.c internal/tiff.c.patch
	sh -c "cd tcl9.0.1/win ; ./configure --build=x86_64-linux-gnu --host=i686-w64-mingw32"
	make -C tcl9.0.1/win -j$(GOMAXPROCS)
	cp -v tcl9.0.1/win/*.dll $(WIN32)
	sh -c "cd tk9.0.1/win ; ./configure --build=x86_64-linux-gnu --host=i686-w64-mingw32 --with-tcl=$(PWD)/tcl9.0.1/win"
	make -C tk9.0.1/win -j$(GOMAXPROCS)
	cp -v tk9.0.1/win/tcl9tk90.dll tk9.0.1/win/libtk9.0.1.zip $(WIN32)
	sh -c "cd Img-2.0.1 ; ./configure --build=x86_64-linux-gnu --host=i686-w64-mingw32 --with-tcl=$(PWD)/tcl9.0.1/win  --with-tk=$(PWD)/tk9.0.1/win"
	make -C Img-2.0.1 -j$(GOMAXPROCS)
	find Img-2.0.1 -name \*.dll -exec cp {} $(WIN32) \;
	go run internal/shasig.go -goos=windows -goarch=386 tk_windows_386.go
	gofmt -l -s -w tk_windows_386.go
	zip -j $(WIN32)/lib.zip.tmp $(WIN32)/*.dll $(WIN32)/*.zip
	rm -f $(WIN32)/*.dll $(WIN32)/*.zip
	mv $(WIN32)/lib.zip.tmp $(WIN32)/lib.zip
	rm -rf tcl9.0.1/ tk9.0.1/ Img-2.0.1/
	GOOS=windows GOARCH=386 go build -v
	GOOS=windows GOARCH=386 ./unconvert.sh
	git status

lib_winarm64: download
	if [ "$(GOOS)" != "linux" ]; then exit 1 ; fi
	if [ "$(GOARCH)" != "amd64" ]; then exit 1 ; fi
	rm -rf tcl9.0.1/ tk9.0.1/ Img-2.0.1/ $(WINARM64)
	mkdir -p $(WINARM64)
	tar xf $(TCL_TAR)
	tar xf $(TK_TAR)
	tar xf $(TK_IMG_TAR)
	patch Img-2.0.1/tiff/tiff.c internal/tiff.c.patch
	patch Img-2.0.1/base/tkimg.h internal/tkimg.h.patch
	sh -c "cd tcl9.0.1/win ; ./configure --build=x86_64-linux-gnu --host=aarch64-w64-mingw32 --enable-64bit=arm64"
	make -C tcl9.0.1/win -j$(GOMAXPROCS)
	cp -v tcl9.0.1/win/*.dll $(WINARM64)
	sh -c "cd tk9.0.1/win ; ./configure --build=x86_64-linux-gnu --host=aarch64-w64-mingw32 --enable-64bit=arm64 --with-tcl=$(PWD)/tcl9.0.1/win"
	make -C tk9.0.1/win -j$(GOMAXPROCS)
	cp -v tk9.0.1/win/tcl9tk90.dll tk9.0.1/win/libtk9.0.1.zip $(WINARM64)
	sh -c "cd Img-2.0.1 ; ./configure --build=x86_64-linux-gnu --host=aarch64-w64-mingw32 --enable-64bit=arm64 --with-tcl=$(PWD)/tcl9.0.1/win  --with-tk=$(PWD)/tk9.0.1/win"
	make -C Img-2.0.1 -j$(GOMAXPROCS)
	find Img-2.0.1 -name \*.dll -exec cp {} $(WINARM64) \;
	go run internal/shasig.go -goos=windows -goarch=arm64 tk_windows_arm64.go
	gofmt -l -s -w tk_windows_arm64.go
	zip -j $(WINARM64)/lib.zip.tmp $(WINARM64)/*.dll $(WINARM64)/*.zip
	rm -f $(WINARM64)/*.dll $(WINARM64)/*.zip
	mv $(WINARM64)/lib.zip.tmp $(WINARM64)/lib.zip
	rm -rf tcl9.0.1/ tk9.0.1/ Img-2.0.1/
	GOOS=windows GOARCH=arm64 go build -v
	GOOS=windows GOARCH=arm64 ./unconvert.sh
	git status

lib_linux_ccgo: download
	if [ "$(GOOS)" != "linux" ]; then exit 1 ; fi
	rm -rf Img-2.0.1/ internal/img/img_$(GOOS)_$(GOARCH).go
	mkdir -p embed/$(GOOS)/$(GOARCH)
	tar xf $(TK_IMG_TAR)
	patch Img-2.0.1/compat/libpng/pngpriv.h internal/pngpriv.h.patch
	patch Img-2.0.1/compat/libtiff/libtiff/tif_dirinfo.c internal/tif_dirinfo.c.patch
	sh -c "cd Img-2.0.1 ; ./configure \
		--with-tcl=$(HOME)/.config/ccgo/v4/libtcl9.0/linux/$(GOARCH)/tcl9.0.1/unix/ \
		--with-tk=$(HOME)/.config/ccgo/v4/libtk9.0/linux/$(GOARCH)/tk9.0.1/unix/"
	ccgo \
		-ignore-unsupported-alignment \
		-exec make -C Img-2.0.1 -j$(GOMAXPROCS)
	mkdir -p internal/img
	./img_ccgo.sh
	rm -rf Img-2.0.1/
	go build -v
	./unconvert.sh
	git status

lib_linux_purego: download
	if [ "$(GOOS)" != "linux" ]; then exit 1 ; fi
	rm -rf tcl9.0.1/ tk9.0.1/ Img-2.0.1/ embed/$(GOOS)/$(GOARCH)
	mkdir -p embed/$(GOOS)/$(GOARCH)
	tar xf $(TCL_TAR)
	tar xf $(TK_TAR)
	tar xf $(TK_IMG_TAR)
	sh -c "cd tcl9.0.1/unix ; ./configure --disable-dll-unloading"
	make -C tcl9.0.1/unix -j$(GOMAXPROCS)
	cp -v tcl9.0.1/unix/*.so embed/$(GOOS)/$(GOARCH)
	sh -c "cd tk9.0.1/unix ; ./configure --with-tcl=$(PWD)/tcl9.0.1/unix"
	make -C tk9.0.1/unix -j$(GOMAXPROCS)
	cp -v tk9.0.1/unix/*.so tk9.0.1/unix/libtk9.0.1.zip embed/$(GOOS)/$(GOARCH)
	sh -c "cd Img-2.0.1 ; ./configure --with-tcl=$(PWD)/tcl9.0.1/unix  --with-tk=$(PWD)/tk9.0.1/unix"
	make -C Img-2.0.1 -j$(GOMAXPROCS)
	find Img-2.0.1 -name \*.so -exec cp {} embed/$(GOOS)/$(GOARCH) \;
	go run internal/shasig.go - tk_$(GOOS)_$(GOARCH).go
	gofmt -l -s -w tk_$(GOOS)_$(GOARCH).go
	zip -j embed/$(GOOS)/$(GOARCH)/lib.zip.tmp embed/$(GOOS)/$(GOARCH)/*
	rm -f embed/$(GOOS)/$(GOARCH)/*.so embed/$(GOOS)/$(GOARCH)/*.zip
	mv embed/$(GOOS)/$(GOARCH)/lib.zip.tmp embed/$(GOOS)/$(GOARCH)/lib.zip
	rm -rf tcl9.0.1/ tk9.0.1/ Img-2.0.1/
	go build -v
	./unconvert.sh
	git status

lib_darwin: download
	if [ "$(GOOS)" != "darwin" ]; then exit 1 ; fi
	rm -rf tcl9.0.1/ tk9.0.1/ Img-2.0.1/ embed/$(GOOS)/$(GOARCH)
	mkdir -p embed/$(GOOS)/$(GOARCH)
	tar xf $(TCL_TAR)
	tar xf $(TK_TAR)
	tar xf $(TK_IMG_TAR)
	sh -c "cd tcl9.0.1/unix ; ./configure"
	make -C tcl9.0.1/unix -j$(GOMAXPROCS)
	cp -v tcl9.0.1/unix/*.dylib embed/$(GOOS)/$(GOARCH)
	sh -c "cd tk9.0.1/unix ; ./configure --with-tcl=$(PWD)/tcl9.0.1/unix --enable-aqua"
	make -C tk9.0.1/unix -j$(GOMAXPROCS)
	cp -v tk9.0.1/unix/*.dylib tk9.0.1/unix/libtk9.0.1.zip embed/$(GOOS)/$(GOARCH)
	sh -c "cd Img-2.0.1 ; ./configure --with-tcl=$(PWD)/tcl9.0.1/unix  --with-tk=$(PWD)/tk9.0.1/unix"
	make -C Img-2.0.1 -j$(GOMAXPROCS)
	find Img-2.0.1 -name \*.dylib -exec cp {} embed/$(GOOS)/$(GOARCH) \;
	go run internal/shasig.go - tk_$(GOOS)_$(GOARCH).go
	gofmt -l -s -w tk_$(GOOS)_$(GOARCH).go
	zip -j embed/$(GOOS)/$(GOARCH)/lib.zip.tmp embed/$(GOOS)/$(GOARCH)/*
	rm -f embed/$(GOOS)/$(GOARCH)/*.dylib embed/$(GOOS)/$(GOARCH)/*.zip
	mv embed/$(GOOS)/$(GOARCH)/lib.zip.tmp embed/$(GOOS)/$(GOARCH)/lib.zip
	rm -rf tcl9.0.1/ tk9.0.1/ Img-2.0.1/
	go build -v
	./unconvert.sh
	git status

# use gmake
lib_freebsd: download
	if [ "$(GOOS)" != "freebsd" ]; then exit 1 ; fi
	rm -rf tcl9.0.1/ tk9.0.1/ Img-2.0.1/ embed/$(GOOS)/$(GOARCH)
	mkdir -p embed/$(GOOS)/$(GOARCH)
	tar xf $(TCL_TAR)
	tar xf $(TK_TAR)
	tar xf $(TK_IMG_TAR)
	sh -c "cd tcl9.0.1/unix ; ./configure --disable-dll-unloading"
	make -C tcl9.0.1/unix -j$(GOMAXPROCS)
	cp -v tcl9.0.1/unix/*.so embed/$(GOOS)/$(GOARCH)
	sh -c "cd tk9.0.1/unix ; ./configure --with-tcl=$(PWD)/tcl9.0.1/unix"
	make -C tk9.0.1/unix -j$(GOMAXPROCS)
	cp -v tk9.0.1/unix/*.so tk9.0.1/unix/libtk9.0.1.zip embed/$(GOOS)/$(GOARCH)
	sh -c "cd Img-2.0.1 ; ./configure --with-tcl=$(PWD)/tcl9.0.1/unix  --with-tk=$(PWD)/tk9.0.1/unix"
	make -C Img-2.0.1 -j$(GOMAXPROCS)
	find Img-2.0.1 -name \*.so -exec cp {} embed/$(GOOS)/$(GOARCH) \;
	go run internal/shasig.go - tk_$(GOOS)_$(GOARCH).go
	gofmt -l -s -w tk_$(GOOS)_$(GOARCH).go
	zip -j embed/$(GOOS)/$(GOARCH)/lib.zip.tmp embed/$(GOOS)/$(GOARCH)/*
	rm -f embed/$(GOOS)/$(GOARCH)/*.so embed/$(GOOS)/$(GOARCH)/*.zip
	mv embed/$(GOOS)/$(GOARCH)/lib.zip.tmp embed/$(GOOS)/$(GOARCH)/lib.zip
	rm -rf tcl9.0.1/ tk9.0.1/ Img-2.0.1/
	go build -v
	./unconvert.sh
	git status

demo:
	cd _examples && go run demo.go

examples:
	./examples.sh

xvfb:
	MODERNC_BUILDER=1 DISPLAY= XVFB_DISPLAY=:6060 go test -v -timeout 24h
