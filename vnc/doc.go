// Copyright 2024 The tk9.0-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package vnc provides tk9.0 applications with a built-in VNC over websockets
// server.
//
// # Supported targets
//
// This package works only on the [targets supported by tk9.0] having X11 as a
// backend.  Currently that means Linux and FreeBSD. This package still builds
// on other taregts but there it does nothing.
//
// # Browser requirements
//
// The client web page uses [noVNC]. Quoting from the respective [README.md]
// (retrieved 2024-11-27):
//
// noVNC uses many modern web technologies so a formal requirement list is not
// available. However these are the minimum versions we are currently aware of:
//
//   - Chrome 89, Firefox 89, Safari 15, Opera 75, Edge 89
//
// # Run time requirements
//
// This package needs to be able to execute multiple instances of [Xvfb],
// [x11vnc] and [websockify].
//
// # How to use it
//
// To add the VNC server to an existing tk9.0 application add
//
//	import _ "modernc.org/tk9.0/vnc"
//
// to the package import clause. The application continues to work as before but it will
// now check at initialization for -vnc.* CLI flags.
//
// # Configuration flags
//
//   - -vnc.nopw
//
//     Clients will be allowed to connect without password. This is the
//     default, but it is not a recommended flag for other use than development
//     and connecting from localhost. Override by -vnc.usepw.
//
//   - -vnc.poll.interval <time.Duration>
//
//     The VNC server will poll connected clients at this interval to discover
//     disconnected clients. Defaults to 30s.
//
//   - -vnc.poll.variance <time.Duration>
//
//     A random value in [0, variance] is added to the vnc.poll.interval value
//     on each poll cycle. Defaults to 1m.
//
//   - -vnc.port <port num>
//
//     Set the web server port to <port num>. Defaults to 1221. Web clients connect by
//     opening, for example http://localhost:1221.
//
//   - -vnc.quality <num>
//
//     If present, set VNC quality to <num>. Must be in [0, 9]. Default is
//     provided by [x11vnc] and it seems to be 6.
//
//   - -vnc.serve
//
//     This flags stops the application from starting. Instead it starts the
//     VNC server, listening for client connections. When a client connects,
//     the server will start a new app instance, pass all non -vnc.* flags to
//     it and this new instance is the one that shows in the web browser
//     window. The default is to start the app normally.
//
//   - -vnc.usepw
//
//     Require a password to connect. See man 1 [x11vnc] for the details about
//     setting a password on the server machine. Defaults to false.
//
//   - -vnc.verbose
//
//     Produce a lot of additional output on stderr. Can be useful when
//     debugging connection, permission or other issues. Defaults to false.
//
// # How it works
//
// This package is inspired by [Jeff Smith's] [CloudTk] but does not use any of its code.
//
// The VNC server starts a new [Xvfb], [x11vnc] and [websockify] instance per
// connecting web client. The web client is initially served a small web page
// that determines the dimensions of the browser window and then redirects the
// client to a [noVNC] page, connected to a new app instance running on
// the server using a properly sized X11 virtual frame buffer.
//
// The rest is network traffic.
//
// [CloudTk]: https://cloudtk.tcl-lang.org/
// [Jeff Smith's]: https://wiki.tcl-lang.org/page/Jeff+Smith
// [README.md]: https://github.com/novnc/noVNC/blob/master/README.md
// [Xvfb]: https://en.wikipedia.org/wiki/Xvfb
// [noVNC]: https://github.com/novnc/noVNC
// [targets supported by tk9.0]: https://pkg.go.dev/modernc.org/tk9.0#hdr-Supported_targets
// [websockify]: https://github.com/novnc/websockify
// [x11vnc]: https://en.wikipedia.org/wiki/X11vnc
package vnc // import "modernc.org/tk9.0/vnc"
