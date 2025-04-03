GOOS=$(go env GOOS)
GOARCH=$(go env GOARCH)
ccgo $(find Img-2.0.1/ -name \*.o.go | sort) \
	$HOME/.config/ccgo/v4/libtcl9.0/linux/"$GOARCH"/tcl9.0.1/unix/tclStubLib.o \
	$HOME/.config/ccgo/v4/libtk9.0/linux/"$GOARCH"/tk9.0.1/unix/tkStubLib.o \
	--package-name img \
	--prefix-external=X \
	--prefix-field=F \
	--prefix-static-internal=_ \
	--prefix-static-none=_ \
	--prefix-tagged-struct=T \
	--prefix-tagged-union=T \
	--prefix-typename=T \
	--prefix-undefined=_ \
	-ignore-link-errors \
	-lX11 \
	-o internal/img/img_"$GOOS"_"$GOARCH".go \
