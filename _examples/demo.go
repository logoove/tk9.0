package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"os/exec"
	"time"

	. "modernc.org/tk9.0"
)

var (
	demos = []struct {
		name string
		x, y int
	}{
		{"photo.go", 1200, 50},
		{"tex.go", 1500, 0},
		{"splot.go", 470, 90},
		{"tori.go", 1050, 400},
		{"svg.go", 0, 0},
		{"font.go", 0, 500},
		{"text.go", 60, 400},
		{"embed.go", 460, 700},
		{"calc.go", 540, 0},
		{"b5.go", 840, 0},
	}
	sleep = flag.Duration("t", time.Second, "")
)

func main() {
	os.Setenv("TK9_DEMO", "1")
	w, h := WmMaxSize(App)
	kx := float64(w) / 1920
	ky := float64(h) / 1080
	var cmds []*exec.Cmd
	for _, v := range demos {
		args := []string{"run", v.name, fmt.Sprintf("+%v+%v", round(kx*float64(v.x)), round(ky*float64(v.y)))}
		fmt.Println(args)
		cmd := exec.Command("go", args...)
		cmds = append(cmds, cmd)
		fmt.Println("starting", v.name)
		cmd.Start()
		time.Sleep(*sleep)
	}
}

func round(n float64) int {
	return int(math.Round(n))
}
