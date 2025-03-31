package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  int = 640
	screenHeight int = 480

	paddleWidth  int = 10
	paddleHeight int = 40

	offsetHorizonal int = 20
	offsetVertical  int = 20
)

type Game struct {
	counter int
	diff    int
}

func (g *Game) Init() {
	g.counter = 20
	g.diff = 2
}

func NewGame() ebiten.Game {
	g := &Game{}
	g.Init()
	return g
}

func (g *Game) Update() error {
	g.counter += g.diff
	if g.counter < offsetVertical {
		g.counter = offsetVertical
		g.diff *= -1
	} else if g.counter > screenHeight-paddleHeight-offsetVertical {
		g.counter = screenHeight - paddleHeight - offsetVertical
		g.diff *= -1
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)

	ebitenutil.DebugPrint(screen, "PONG")

	i := ebiten.NewImage(paddleWidth, paddleHeight)
	i.Fill(color.White)
	op := &ebiten.DrawImageOptions{}

	for idx := range []string{"player", "villain"} {
		op.GeoM.Reset()
		op.GeoM.Translate(float64(offsetHorizonal+(screenWidth-((2*offsetHorizonal)+paddleWidth))*idx), float64(g.counter))
		screen.DrawImage(i, op)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Pong")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
