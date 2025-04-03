// Copyright 2024 The tk9.0-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build linux || freebsd

package vnc // import "modernc.org/tk9.0/vnc"

import (
	"bufio"
	"context"
	"embed"
	"errors"
	"fmt"
	"html/template"
	"io"
	"math"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/mileusna/useragent"
	"modernc.org/mathutil"
	"modernc.org/opt"
	"modernc.org/tk9.0"
)

func init() {
	checkServe()
}

const (
	// ScaleEnvVarMobile, if a valid (floating point) number, overrides the
	// TK9_SCALE value for mobile clients.
	ScaleEnvVarMobile = "TK9_VNC_SCALE_MOBILE"

	// EnvVarVNC is non-blank when app runs via VNC.
	EnvVarVNC = "TK9_VNC"

	// EnvVarMobile is non-blank when a mobile client connects.
	EnvVarMobile = "TK9_VNC_MOBILE"

	// EnvVarInstanceStart is set to the unix milliseconds when the client instance
	// started.
	EnvVarInstanceStart = "TK9_VNC_INSTANCE_START"

	// EnvVarInstanceStat is set to the unix milliseconds of the client instance
	// stat mtime.
	EnvVarInstanceStat = "TK9_VNC_INSTANCE_STAT"

	// defaultlMaxXServerNumber is the default value of MaxXServerNumber.
	defaultlMaxXServerNumber = 75
	// defaultPort is the default value of the -vnc.port CLI flag
	defaultPort = 1221

	depth = 16 // Xvfb screen depth

	scaleEnvVar = tk9_0.ScaleEnvVar
)

// maxXServerNumber limits the number X servers in [1..maxXServerNumber]
var maxXServerNumber = defaultlMaxXServerNumber

var (
	appBin        string
	appTitle      string
	clients       = newClientRegister()
	dbg           = os.Getenv("TK9_VNC_DEBUG") != ""
	mu            sync.Mutex
	prng          *prng32
	vncArgs       []string
	vncHTML       *template.Template
	websockifyBin = "websockify"
	x11vncBin     = "x11vnc"
	xvfbBin       = "Xvfb"

	//go:embed embed
	assets embed.FS
)

func log(s string, args ...any) {
	s = strings.TrimSpace(s) + "\n"
	fmt.Fprintf(os.Stderr, s, args...)
}

func existingXServers() (r map[int]struct{}, err error) {
	m, err := filepath.Glob(filepath.Join(os.TempDir(), "\\.X*"))
	if len(m) == 0 || err != nil {
		log("%v", err)
		return nil, err
	}

	r = map[int]struct{}{0: {}}
	for _, v := range m {
		b := filepath.Base(v)
		b = strings.TrimLeft(b, ".X")
		b = strings.TrimRight(b, "-lock")
		if b != "" {
			if n, err := strconv.ParseInt(b, 10, 32); err == nil {
				r[int(n)] = struct{}{}
			}

		}
	}
	return r, nil
}

func allocDisplay() (r int, err error) {
	mu.Lock()

	defer mu.Unlock()

	ex, err := existingXServers()
	if err != nil {
		log("%v", err)
		return 0, err
	}

	for i := 1; i < maxXServerNumber; i++ {
		if _, ok := ex[i]; !ok {
			return i, nil
		}
	}

	return 0, fmt.Errorf("cannot find free X server number")
}

func (f *flags) start(bin string, args []string, pipe bool, env map[string]string) (cmd *exec.Cmd, cancel context.CancelFunc, stdout io.ReadCloser, err error) {
	return f.start0(bin, args, pipe, false, env)
}

func (f *flags) start0(bin string, args []string, pipe, silent bool, env map[string]string) (cmd *exec.Cmd, cancel context.CancelFunc, stdout io.ReadCloser, err error) {
	if !silent && f.verbose {
		defer func() {
			pid := -1
			if cmd != nil {
				pid = cmd.Process.Pid
			}
			fmt.Fprintf(os.Stderr, "bin=%q args=%v -> PID=%v err=%v\n", bin, args, pid, err)
		}()
	}
	ctx, cancel := context.WithCancel(context.Background())
	doCancel := true

	defer func() {
		if doCancel {
			cancel()
		}
	}()

	cmd = exec.CommandContext(ctx, bin, args...)
	if env != nil {
		cmd.Env = os.Environ()
		for k, v := range env {
			cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
		}
	}
	cmd.WaitDelay = time.Second
	if pipe {
		if stdout, err = cmd.StdoutPipe(); err != nil {
			log("%v", err)
			return nil, nil, nil, err
		}
	}
	if !silent && f.verbose {
		cmd.Stderr = os.Stderr
	}
	if err := cmd.Start(); err != nil {
		log("%v", err)
		return nil, nil, nil, err
	}

	go func() {
		<-ctx.Done()
	}()

	doCancel = false
	return cmd, cancel, stdout, nil
}

func (f *flags) startX11vnc(args []string) (cmd *exec.Cmd, cancel context.CancelFunc, port int, err error) {
	cmd, cancel, stdout, err := f.start(x11vncBin, args, true, nil)
	if err != nil {
		log("%v", err)
		return nil, nil, 0, err
	}

	sc := bufio.NewReader(stdout)
	s, err := sc.ReadString('\n')
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "PORT=")
	n, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		log("%v", err)
		return nil, nil, port, err
	}

	port = int(n)
	return cmd, cancel, port, nil
}

type flags struct {
	noVNCQuality int           // 0-9
	pollInterval time.Duration // Check client connected at pollInterval+rnd(pollVariance)
	pollVariance time.Duration

	port uint16

	serve          bool
	verbose        bool
	x11vncPassword bool
}

func newFlags() *flags {
	return &flags{
		noVNCQuality: -1,
		port:         defaultPort,
		pollInterval: 30 * time.Second,
		pollVariance: time.Minute,
	}
}

func parseFlags(in []string) (r *flags, out []string, err error) {
	set := opt.NewSet()
	r = newFlags()
	set.Arg("vnc.port", false, func(opt, arg string) error {
		n, err := strconv.ParseUint(arg, 10, 16)
		if err == nil {
			r.port = uint16(n)
		}
		return err
	})
	set.Arg("vnc.poll.interval", false, func(opt, arg string) error {
		if n, err := time.ParseDuration(arg); err == nil {
			r.pollInterval = n
		}
		return nil
	})
	set.Arg("vnc.poll.variance", false, func(opt, arg string) error {
		if n, err := time.ParseDuration(arg); err == nil {
			r.pollVariance = n
		}
		return nil
	})
	set.Arg("vnc.quality", false, func(opt, arg string) error {
		n, err := strconv.ParseUint(arg, 10, 16)
		if err == nil && n <= 9 {
			r.noVNCQuality = int(n)
		}
		return nil
	})
	set.Opt("vnc.nopw", func(opt string) error { r.x11vncPassword = false; return nil })
	set.Opt("vnc.serve", func(opt string) error { r.serve = true; return nil })
	set.Opt("vnc.usepw", func(opt string) error { r.x11vncPassword = true; return nil })
	set.Opt("vnc.verbose", func(opt string) error { r.verbose = true; return nil })
	err = errors.Join(err, set.Parse(in, func(opt string) error {
		switch {
		case strings.HasPrefix(opt, "-vnc."):
			return fmt.Errorf("unknown vnc option: %s", opt)
		default:
			out = append(out, opt)
			return nil
		}
	}))
	return r, out, err
}

type prng32 struct {
	sync.Mutex
	prng *mathutil.FC32
}

func newPrng32() (r *prng32, err error) {
	prng, err := mathutil.NewFC32(math.MinInt32, math.MaxInt32, true)
	if err != nil {
		log("%v", err)
		return nil, err
	}

	return &prng32{prng: prng}, nil
}

func (p *prng32) id() uint {
	p.Lock()

	defer p.Unlock()

	return uint(uint32(p.prng.Next()))
}

func lookPath(s string) (r string) {
	r = s
	if s, err := exec.LookPath(s); err == nil {
		r = s
	}
	return r
}

func checkServe() {
	flags, out, err := parseFlags(os.Args)
	os.Args = out
	if err != nil {
		log("%v", err)
		return
	}

	if !flags.serve {
		return
	}

	appTitle = filepath.Base(os.Args[0])
	appTitle = strings.TrimSuffix(appTitle, ".exe")
	xvfbBin = lookPath(xvfbBin)
	websockifyBin = lookPath(websockifyBin)
	x11vncBin = lookPath(x11vncBin)
	b, err := assets.ReadFile("embed/vnc.html")
	if err != nil {
		log("%v", err)
		return
	}

	if vncHTML, err = template.New("vnc").Parse(string(b)); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	if prng, err = newPrng32(); err != nil {
		log("%v", err)
		return
	}

	prng.prng.Seed(time.Now().UnixNano())
	if appBin, err = os.Executable(); err != nil {
		log("%v", err)
		return
	}

	http.Handle("/", flags)
	if flags.verbose {
		fmt.Fprintf(os.Stderr, "HTTP server listening at :%d\n", flags.port)
	}
	if err = http.ListenAndServe(fmt.Sprintf(":%d", flags.port), nil); err != nil {
		log("%v", err)
		os.Exit(1)
	}
}

func (f *flags) err(w http.ResponseWriter, code int) {
	http.Error(w, http.StatusText(code), code)
}

func (f *flags) ServeHTTP(w http.ResponseWriter, rq *http.Request) {
	if rq.Method != "GET" {
		f.err(w, http.StatusMethodNotAllowed)
		return
	}

	if dbg {
		trc("GET %q %q %q", rq.URL.Path, rq.Host, rq.Header)
	}
	switch {
	case rq.URL.Path == "/":
		f.connect(w)
	case
		strings.HasPrefix(rq.URL.Path, "/core/"),
		strings.HasPrefix(rq.URL.Path, "/favicon"),
		strings.HasPrefix(rq.URL.Path, "/vendor/"):

		p := path.Join("embed", strings.ReplaceAll(rq.URL.Path, "/vendor/", "/vendor_/"))
		if dbg {
			f, err := assets.Open(p)
			trc("%q -> (%p, %v)", p, f, err)
			if err == nil {
				f.Close()
			}
		}
		http.ServeFileFS(w, rq, assets, p)
	default:
		a := strings.Split(rq.URL.Path, "_")
		if len(a) != 3 {
			f.err(w, http.StatusBadRequest)
			return
		}

		clientID := a[0]
		width := a[1]
		height := a[2]
		c := clients.get(clientID)

		defer c.Unlock()

		if c.isConnected || c.disconnected {
			f.connect(w)
			return
		}

		host := rq.Host
		if x := strings.IndexByte(host, ':'); x != 0 {
			host = host[:x]
		}
		c.flags = f
		c.connect(w, host, clientID, width, height, rq)
	}
}

func (f *flags) connect(w http.ResponseWriter) {
	fmt.Fprintf(w, `<!DOCTYPE html>
<html lang="en">
<head>
    <script>
	    function bodyOnload() {
		    window.location.assign(%s/%d_${window.innerWidth}_${window.innerHeight}%[1]s)
	    }
    </script>
</head>
<body style="background-color:#eee;margin:0;min-width:100vw;min-height:100vh" onload="bodyOnload();">
</body>
</html>`, "`", prng.id())
}

type client struct {
	sync.Mutex

	appCancel        context.CancelFunc
	appCmd           *exec.Cmd
	display          string // :1, :2, ...
	flags            *flags
	id               string
	port             int
	websockifyCancel context.CancelFunc
	websockifyCmd    *exec.Cmd
	x11vncCancel     context.CancelFunc
	x11vncCmd        *exec.Cmd
	xvfbCancel       context.CancelFunc
	xvfbCmd          *exec.Cmd

	disconnected bool
	isConnected  bool
}

func newClient() *client {
	return &client{}
}

func cancel(f ...func()) {
	for _, g := range f {
		if g != nil {
			g()
		}
	}
}

func (c *client) after() <-chan time.Time {
	return time.After(c.flags.pollInterval + time.Duration(rand.Int63n(int64(c.flags.pollVariance))))
}

func (c *client) poll() {
	for ch := c.after(); ; ch = c.after() {
		select {
		case <-ch:
			c.Lock()

			if c.disconnected {
				clients.delete(c.id)
				c.Unlock()
				return
			}

			cmd, cancel, stdout, err := c.flags.start0(x11vncBin, []string{"-query", "client_count", "-display", c.display}, true, true, nil)
			if err != nil {
				log("%v", err)
				c.Unlock()
				continue
			}

			sc := bufio.NewReader(stdout)
			s, err := sc.ReadString('\n')
			cancel()
			cmd.Wait()
			if err != nil {
				log("%v", err)
				c.Unlock()
				continue
			}

			if s = strings.TrimSpace(s); strings.HasSuffix(s, ":0") {
				c.Unlock()
				c.disconnect()
				return
			}

			c.Unlock()
		}
	}
}

func (c *client) disconnect1(cf *context.CancelFunc, pcmd **exec.Cmd) {
	defer func() {
		*cf = nil
		*pcmd = nil
	}()

	if cancel := *cf; cancel != nil {
		cancel()
	}
	if cmd := *pcmd; pcmd != nil {
		go cmd.Wait()
	}
}

func (c *client) disconnect() {
	c.Lock()

	defer func() {
		c.Unlock()
	}()

	if !c.isConnected || c.disconnected {
		return
	}

	defer func() {
		c.disconnected = true
		c.isConnected = false
	}()

	if c.flags.verbose {
		fmt.Fprintf(os.Stderr, "disconnecting DISPLAY=%v\n", c.display)
	}
	c.disconnect1(&c.appCancel, &c.appCmd)
	c.disconnect1(&c.websockifyCancel, &c.websockifyCmd)
	c.disconnect1(&c.x11vncCancel, &c.x11vncCmd)
	c.disconnect1(&c.xvfbCancel, &c.xvfbCmd)
	tmp := os.TempDir()
	if !strings.HasPrefix(tmp, "/") || tmp == "/" {
		panic(todo("internal error"))
	}

	display := c.display[1:]
	arg := fmt.Sprintf("rm -rf '%s'", filepath.Join(tmp, fmt.Sprintf(".X%s-lock", display)))
	if c.flags.verbose {
		fmt.Fprintf(os.Stderr, "exec `sh -c %s`\n", arg)
	}
	exec.Command("sh", "-c", arg).Run()
}

func (c *client) connect(w http.ResponseWriter, host, id, width, height string, rq *http.Request) {
	defer func() {
		if c.isConnected {
			go c.poll()
			return
		}

		cancel(c.appCancel, c.websockifyCancel, c.x11vncCancel, c.xvfbCancel)
		c.appCancel = nil
		c.appCmd = nil
		c.websockifyCancel = nil
		c.websockifyCmd = nil
		c.x11vncCancel = nil
		c.x11vncCmd = nil
		c.xvfbCancel = nil
		c.xvfbCmd = nil
	}()

	c.id = id
	displayNum, err := allocDisplay()
	if err != nil {
		log("%v", err)
		http.Error(w, "cannot allocate new X server", http.StatusTooManyRequests)
		return
	}

	display := fmt.Sprintf(":%d", displayNum)
	args := []string{display, "-screen", "0", fmt.Sprintf("%sx%sx%d", width, height, depth)}
	if c.xvfbCmd, c.xvfbCancel, _, err = c.flags.start(xvfbBin, args, false, nil); err != nil {
		log("%v", err)
		http.Error(w, "cannot create new X server", http.StatusFailedDependency)
		return
	}

	args = []string{"-display", display, "-forever", "-autoport", "5900", "-noshm"}
	switch {
	case c.flags.x11vncPassword:
		args = append(args, "-usepw")
	default:
		args = append(args, "-nopw")
	}
	if c.x11vncCmd, c.x11vncCancel, c.port, err = c.flags.startX11vnc(args); err != nil {
		log("%v", err)
		http.Error(w, "cannot create new VNC server", http.StatusFailedDependency)
		return
	}

	args = []string{fmt.Sprint(c.port), fmt.Sprintf("localhost:%d", c.port)}
	if c.websockifyCmd, c.websockifyCancel, _, err = c.flags.start(websockifyBin, args, false, nil); err != nil {
		log("%v", err)
		http.Error(w, "cannot start websockify", http.StatusFailedDependency)
		return
	}

	tArgs := struct {
		Port    int
		Quality int
		Title   string
	}{
		Port:    c.port,
		Quality: c.flags.noVNCQuality,
		Title:   appTitle,
	}
	if err = vncHTML.Execute(w, tArgs); err != nil {
		log("%v", err)
		http.Error(w, "cannot execute html template", http.StatusInternalServerError)
		return
	}

	mu.Lock()

	defer mu.Unlock()

	m := map[string]string{EnvVarMobile: ""}
	if isMobileClient(rq) {
		m[EnvVarMobile] = "1"
		if s := os.Getenv(ScaleEnvVarMobile); s != "" {
			m[scaleEnvVar] = s
		}
	}
	m["DISPLAY"] = display
	m[EnvVarVNC] = "1"
	m["TK9_VNC_WIDTH"] = fmt.Sprint(width)
	m["TK9_VNC_HEIGHT"] = fmt.Sprint(height)
	m["TK9_VNC_DEPTH"] = fmt.Sprint(depth)
	m[EnvVarInstanceStart] = fmt.Sprint(time.Now().UTC().UnixMilli())
	if fi, err := os.Stat(appBin); err == nil {
		m[EnvVarInstanceStat] = fmt.Sprint(fi.ModTime().UTC().UnixMilli())
	}
	if c.appCmd, c.appCancel, _, err = c.flags.start(appBin, os.Args[1:], false, m); err != nil {
		log("%v", err)
		http.Error(w, "cannot start new application instance", http.StatusFailedDependency)
		return
	}

	c.display = display

	go func(c *client) {
		pid := c.appCmd.Process.Pid
		c.appCmd.Wait()
		if c.flags.verbose {
			fmt.Fprintf(os.Stderr, "client exited: DISPLAY=%v PID=%v\n", c.display, pid)
		}
		c.disconnect()
	}(c)

	if c.flags.verbose {
		fmt.Fprintf(os.Stderr, "VNC server for DISPLAY=%s listening on :%d\n", display, c.port)
	}
	c.isConnected = true
}

func isMobileClient(rq *http.Request) (r bool) {
	if dbg {
		defer func() { trc("", r) }()
	}

	for _, v := range rq.Header["User-Agent"] {
		if useragent.Parse(v).Mobile {
			return true
		}
	}
	return false
}

type clientRegister struct {
	sync.Mutex
	m map[string]*client
}

func newClientRegister() *clientRegister {
	return &clientRegister{m: map[string]*client{}}
}

func (c *clientRegister) get(id string) (r *client) {
	c.Lock()

	defer c.Unlock()

	if r = c.m[id]; r == nil {
		r = newClient()
		c.m[id] = r
	}
	r.Lock()
	return r
}

func (c *clientRegister) delete(id string) {
	c.Lock()

	defer c.Unlock()

	delete(c.m, id)
}
