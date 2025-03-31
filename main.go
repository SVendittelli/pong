package main

import (
	"image/color"
	"log"
	"math/rand"

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

	villainDir      int
	villainCoolDown int
}

func (g *Game) Init() {
	g.speed = 2
	g.playerY = (screenHeight - paddleHeight) / 2
	g.villainY = (screenHeight - paddleHeight) / 2

	g.villainDir = 1
	g.villainCoolDown = 0
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

	if g.villainCoolDown > 0 {
		g.villainY += g.speed * g.villainDir
		g.villainCoolDown--
	} else {
		g.villainDir = (1 - 2*rand.Intn(2))
		g.villainCoolDown = rand.Intn(30) + 10 // Random cooldown between 10 and 40 frames
	}

	// Clamp the vertical location of the paddle within the bounds of the screen
	g.playerY = Clamp(g.playerY, offsetVertical, screenHeight-paddleHeight-offsetVertical)
	g.villainY = Clamp(g.villainY, offsetVertical, screenHeight-paddleHeight-offsetVertical)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)

	ebitenutil.DebugPrint(screen, "PONG")

	i := ebiten.NewImage(paddleWidth, paddleHeight)
	i.Fill(color.White)
	op := &ebiten.DrawImageOptions{}

	// Draw the player's paddle
	op.GeoM.Reset()
	op.GeoM.Translate(float64(offsetHorizonal), float64(g.playerY))
	screen.DrawImage(i, op)

	// Draw the villain's paddle
	op.GeoM.Reset()
	op.GeoM.Translate(float64(offsetHorizonal+(screenWidth-((2*offsetHorizonal)+paddleWidth))), float64(g.villainY))
	screen.DrawImage(i, op)
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
