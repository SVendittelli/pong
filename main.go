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
)

type Game struct{}

func NewGame() ebiten.Game {
	g := &Game{}
	return g
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x82, 0xa3, 0xff, 0xff})

	ebitenutil.DebugPrint(screen, "PONG")

	i := ebiten.NewImage(10, 30)
	i.Fill(color.White)

	for idx := range []string{"player", "villain"} {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Reset()
		op.GeoM.Translate(float64(10+screenWidth*idx-30*idx), 20)
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
