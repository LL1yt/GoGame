package main

import (
	"math/rand"
	"time"

	"GoGame/internal/game"
	"GoGame/internal/ui"

	"fyne.io/fyne/v2/app"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	a := app.New()
	w := a.NewWindow("Go Card Game")

	g := game.NewGame()
	g.SetWindow(w)

	// Initialize the deck and deal initial hands
	g.DealInitialHands()

	ui.SetupUI(g)

	// Start the game loop in a separate goroutine
	go g.GameLoop()

	w.ShowAndRun()
}