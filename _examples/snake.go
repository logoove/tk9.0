// Based on https://g.co/gemini/share/993bd674bf40

package main

import (
	"fmt"
	"image"
	"math/rand"
	"time"

	. "modernc.org/tk9.0"
	_ "modernc.org/tk9.0/themes/azure"
	_ "modernc.org/tk9.0/vnc"
)

const (
	width    = 20 * gridSize
	height   = 15 * gridSize
	gridSize = 20
)

var game *Game

type Game struct {
	bodyColor Opt
	canvas    *CanvasWidget
	drawCycle int
	dx, dy    int
	food      image.Point
	headColor Opt
	points    int
	score     *Window
	speed     *Window
	gear      int
	snake     []image.Point
	t         *Ticker
	tick      int
	isRunning bool
}

func round(n, to int) int {
	return n - n%to
}

func NewGame() *Game {
	App.WmTitle("Snake Game")
	canvas := Canvas(Background(Linen), Width(width), Height(height))
	g := &Game{
		bodyColor: Fill(Green),
		canvas:    canvas,
		headColor: Fill(Yellow),
		score:     Label(Justify("center")).Window,
		speed:     TLabel(Justify("center"), Txt("Speed 1 (F1-F4 to change)")).Window,
	}
	Pack(g.speed, canvas, g.score, Padx("1m"), Pady("2m"), Ipadx("1m"), Ipady("1m"))
	Bind(App, "<Key>", Command(g.handleKeyPress))
	Bind(g.score, "<Button-1>", Command(func() {
		if !g.isRunning {
			g.init()
		}

	}))
	g.init()
	g.t, _ = NewTicker(100*time.Millisecond, func() {
		g.tick++
		if g.gear < 4 && g.tick&3 > g.gear-1 {
			return
		}
		if g.isRunning {
			g.moveSnake()
			g.draw()
		}
	})
	return g
}

func (g *Game) init() {
	g.dx = gridSize
	g.dy = 0
	g.setSpeed(1)
	g.setPoints(0)
	g.snake = []image.Point{{round(width/2, gridSize), round(height/2, gridSize)}}
	g.isRunning = true
	g.spawnFood()
	g.draw()
}

func (g *Game) setPoints(n int) {
	g.points = n
	g.score.Configure(Txt(fmt.Sprintf("Score: %v", g.points)))
}

func (g *Game) moveSnake() {
	head := g.snake[0]
	newHead := image.Point{head.X + g.dx, head.Y + g.dy}
	if newHead.X < 0 || newHead.X >= width || newHead.Y < 0 || newHead.Y >= height || g.collidesWithSelf(newHead) {
		g.isRunning = false
		g.score.Configure(Txt(fmt.Sprintf("Score: %v\nGame Over!\nClick to restart", g.points)))
		return
	}

	g.snake = append([]image.Point{newHead}, g.snake[:len(g.snake)-1]...)
	if newHead == g.food {
		g.snake = append([]image.Point{newHead}, g.snake...)
		g.spawnFood()
		g.setPoints(g.points + g.gear)
	}
}

func (g *Game) draw() {
	g.canvas.Delete("all")
	g.drawCycle++
	for i, p := range g.snake {
		switch {
		case i == 0:
			g.canvas.CreateRectangle(p.X, p.Y, p.X+gridSize, p.Y+gridSize, g.headColor)
		default:
			switch {
			case (i+g.drawCycle)&1 == 0:
				g.canvas.CreatePolygon(p.X, p.Y, p.X+gridSize/2, p.Y, p.X+gridSize, p.Y+gridSize/2, p.X+gridSize, p.Y+gridSize, p.X+gridSize/2, p.Y+gridSize, p.X, p.Y+gridSize/2, g.bodyColor)
			default:
				g.canvas.CreatePolygon(p.X+gridSize, p.Y, p.X+gridSize/2, p.Y, p.X, p.Y+gridSize/2, p.X, p.Y+gridSize, p.X+gridSize/2, p.Y+gridSize, p.X+gridSize, p.Y+gridSize/2, g.bodyColor)
			}
		}
	}
	g.canvas.CreateRectangle(g.food.X, g.food.Y, g.food.X+gridSize, g.food.Y+gridSize, Fill(Red))

}

func (g *Game) spawnFood() {
	g.food = image.Point{
		rand.Intn(width/gridSize) * gridSize,
		rand.Intn(height/gridSize) * gridSize,
	}
}

func (g *Game) handleKeyPress(e *Event) {
	switch e.Keysym {
	case "Left":
		g.dx, g.dy = -gridSize, 0
	case "Right":
		g.dx, g.dy = gridSize, 0
	case "Up":
		g.dx, g.dy = 0, -gridSize
	case "Down":
		g.dx, g.dy = 0, gridSize
	case "F1":
		g.setSpeed(1)
	case "F2":
		g.setSpeed(2)
	case "F3":
		g.setSpeed(3)
	case "F4":
		g.setSpeed(4)
	}
}

func (g *Game) setSpeed(n int) {
	g.gear = n
	g.speed.Configure(Txt(fmt.Sprintf("Speed: %v", n)))
}

func (g *Game) collidesWithSelf(head image.Point) bool {
	for _, p := range g.snake {
		if p == head {
			return true
		}
	}
	return false
}

func main() {
	ActivateTheme("azure light")
	rand.Seed(time.Now().UnixNano())
	game = NewGame()
	App.SetResizable(false, false)
	App.Wait()
}
