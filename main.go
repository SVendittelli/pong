package main

import (
	"bytes"
	"image/color"
	"log"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	raudio "github.com/SVendittelli/pong/resources/audio"
)

const (
	screenWidth  int = 640
	screenHeight int = 480

	paddleWidth  int = 10
	paddleHeight int = 40

	offsetHorizonal int = 20
	offsetVertical  int = 0

	// Ball size
	ballWidth  int = 10
	ballHeight int = 10
)

type Mode int

const (
	ModeTitle Mode = iota
	ModeGame
	ModeGameOver
)

type Game struct {
	mode Mode

	paddleSpeed int
	playerY     int
	villainY    int

	villainDir      int
	villainCoolDown int

	ballX        int
	ballY        int
	ballVelX     float64
	ballVelY     float64
	ballMaxSpeed float64

	audioContext *audio.Context
	bouncePlayer *audio.Player
}

func (g *Game) Init() {
	g.mode = ModeTitle

	g.paddleSpeed = 3
	g.playerY = (screenHeight - paddleHeight) / 2
	g.villainY = (screenHeight - paddleHeight) / 2

	g.villainDir = 0
	g.villainCoolDown = 0

	g.ballMaxSpeed = 4
	g.ballX = screenWidth / 2
	g.ballY = screenHeight / 2
	g.ballVelX = float64(rand.Intn(2)*2-1) * ((g.ballMaxSpeed-1)*rand.Float64() + 1)
	g.ballVelY = g.ballMaxSpeed * (rand.Float64()*2 - 1)

	if g.audioContext == nil {
		g.audioContext = audio.NewContext(48000)
	}

	bounceD, err := wav.DecodeF32(bytes.NewReader(raudio.Bounce_wav))
	if err != nil {
		log.Fatal(err)
	}
	g.bouncePlayer, err = g.audioContext.NewPlayerF32(bounceD)
	if err != nil {
		log.Fatal(err)
	}
}

func NewGame() ebiten.Game {
	g := &Game{}
	g.Init()
	return g
}

func (g *Game) IsUpPressed() bool {
	return ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyK)
}

func (g *Game) IsDownPressed() bool {
	return ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyJ)
}

func (g *Game) PlayBounce() error {
	if err := g.bouncePlayer.Rewind(); err != nil {
		return err
	}
	g.bouncePlayer.Play()
	return nil
}

func (g *Game) Update() error {
	if g.mode == ModeTitle {
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.mode = ModeGame
		}
		return nil
	}

	if g.mode == ModeGameOver {
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.Init()
		}
		return nil
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.Init()
		return nil
	}

	if g.IsUpPressed() {
		g.playerY -= g.paddleSpeed
	} else if g.IsDownPressed() {
		g.playerY += g.paddleSpeed
	}

	if g.villainCoolDown > 0 {
		g.villainY += g.paddleSpeed * g.villainDir
		g.villainCoolDown--
	} else {
		g.villainDir = (1 - 2*rand.Intn(2))
		g.villainCoolDown = rand.Intn(30) + 10 // Random cooldown between 10 and 40 frames
	}

	// Clamp the vertical location of the paddle within the bounds of the screen
	g.playerY = Clamp(g.playerY, offsetVertical, screenHeight-paddleHeight-offsetVertical)
	g.villainY = Clamp(g.villainY, offsetVertical, screenHeight-paddleHeight-offsetVertical)

	// Update the ball position
	g.ballX += int(g.ballVelX)
	g.ballY += int(g.ballVelY)

	g.ballX = Clamp(g.ballX, ballWidth/2, screenWidth-ballWidth/2)
	g.ballY = Clamp(g.ballY, ballWidth/2, screenHeight-ballWidth/2)

	// Check for vertical collision with the walls
	if g.ballY == ballHeight/2 || g.ballY == screenHeight-ballHeight/2 {
		err := g.PlayBounce()
		if err != nil {
			return err
		}
		g.ballVelY *= -1
	}

	// Check for horizontal collision with the paddles
	if math.Abs(float64(g.ballX-ballWidth/2-(offsetHorizonal+paddleWidth))) <= 1 && g.ballY+ballHeight/2 >= g.playerY && g.ballY-ballHeight/2 <= g.playerY+paddleHeight {
		err := g.PlayBounce()
		if err != nil {
			return err
		}
		g.ballVelX *= -1
	} else if math.Abs(float64(g.ballX+ballWidth/2-(screenWidth-offsetHorizonal-paddleWidth))) <= 1 && g.ballY+ballHeight/2 >= g.villainY && g.ballY-ballHeight/2 <= g.villainY+paddleHeight {
		err := g.PlayBounce()
		if err != nil {
			return err
		}
		g.ballVelX *= -1
	}

	// Check for game over condition
	if g.ballX <= ballWidth/2 || g.ballX >= screenWidth-ballWidth/2 {
		g.mode = ModeGameOver
	}

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

	// Draw the ball
	i = ebiten.NewImage(ballWidth, ballHeight)
	i.Fill(color.White)
	op = &ebiten.DrawImageOptions{}

	op.GeoM.Translate(float64(g.ballX-ballWidth/2), float64(g.ballY-ballHeight/2))
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
		panic(err)
	}
}
