package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/diegoxter/blockgame/game"
	"github.com/hajimehoshi/ebiten/v2"
)


const (
	tileSize     = 16
	screenWidth  = 320
	screenHeight = 240
)


func main() {
	g, err := game.NewGame()
	if err != nil {
		log.Fatal(err)
	}

	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("veGame!")

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGINT,
		syscall.SIGTERM)
	go func() {
		<-sigc

		g.Exit()
	}()
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
