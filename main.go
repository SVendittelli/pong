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
	speed    int
	playerY  int
	villainY int
}

func (g *Game) Init() {
	g.speed = 2
	g.playerY = (screenHeight - paddleHeight) / 2
	g.villainY = 20
}

func NewGame() ebiten.Game {
	g := &Game{}
	g.Init()
	return g
}

func (g *Game) isUpPressed() bool {
	return ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyK)
}

func (g *Game) isDownPressed() bool {
	return ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyJ)
}

func (g *Game) Update() error {
	if g.isUpPressed() {
		g.playerY -= g.speed
	} else if g.isDownPressed() {
		g.playerY += g.speed
	}
	// Clamp the vertical location of the paddle within the bounds of the screen
	if g.playerY < offsetVertical {
		g.playerY = offsetVertical
	} else if g.playerY > screenHeight-paddleHeight-offsetVertical {
		g.playerY = screenHeight - paddleHeight - offsetVertical
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
		op.GeoM.Translate(float64(offsetHorizonal+(screenWidth-((2*offsetHorizonal)+paddleWidth))*idx), float64(g.playerY))
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
