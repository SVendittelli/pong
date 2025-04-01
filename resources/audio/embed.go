package audio

import (
	_ "embed"
)

var (
	//go:embed bounce.wav
	Bounce_wav []byte

	//go:embed game_over.wav
	GameOver_wav []byte

	//go:embed title.wav
	Title_wav []byte
)
