package ui

import (
	"fmt"
	"image/color"

	"GoGame/internal/game"
	"GoGame/internal/player"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var endTurnButton *widget.Button
var newGameButton *widget.Button

func SetupUI(g *game.Game) {
	g.ScoreLabel = widget.NewLabel(fmt.Sprintf("Score - %s: %d, %s: %d", g.Player1.Name, g.Player1.Score, g.Player2.Name, g.Player2.Score))
	phaseLabel := widget.NewLabel("Current Phase: Draw")

	player1Field := createPlayerField(g, &g.Player1, true)
	player2Field := createPlayerField(g, &g.Player2, false)

	endTurnButton = widget.NewButton("End Turn", func() {
		g.EndTurn <- true
	})
	endTurnButton.Disable()

	newGameButton = widget.NewButton("New Game", func() {
		startNewGame(g)
	})

	gameBoard := container.NewVBox(
		player2Field,
		widget.NewSeparator(),
		player1Field,
		container.NewHBox(endTurnButton, newGameButton),
	)

	content := container.NewBorder(container.NewVBox(g.ScoreLabel, phaseLabel), nil, nil, nil, gameBoard)

	g.GetWindow().SetContent(content)
	g.GetWindow().Resize(fyne.NewSize(1920, 1080))

	// Set up the UIUpdate function
	g.UIUpdate = func() {
		g.GetWindow().Canvas().Refresh(content)
		updatePhaseLabel(g, phaseLabel)
		updateEndTurnButton(g)
		updatePlayerFields(g, player1Field, player2Field)
	}
}

func startNewGame(g *game.Game) {
	g.Reset()
	g.UIUpdate()
	go g.GameLoop()
}

func updatePhaseLabel(g *game.Game, label *widget.Label) {
	var phase string
	switch g.CurrentPhase {
	case game.DrawPhase:
		phase = "Draw"
	case game.PlayPhase:
		phase = "Play"
	case game.EndPhase:
		phase = "End"
	}
	label.SetText(fmt.Sprintf("Current Phase: %s", phase))
}

func updateEndTurnButton(g *game.Game) {
	if g.CurrentPlayer == &g.Player1 && g.CurrentPhase == game.PlayPhase {
		endTurnButton.Enable()
	} else {
		endTurnButton.Disable()
	}
}

func updatePlayerFields(g *game.Game, player1Field, player2Field *fyne.Container) {
	updatePlayerField(g, &g.Player1, player1Field)
	updatePlayerField(g, &g.Player2, player2Field)
}

func updatePlayerField(g *game.Game, player *player.Player, field *fyne.Container) {
	// Update player card
	playerCard := field.Objects[2].(*fyne.Container).Objects[1].(*fyne.Container)
	updatePlayerCard(player, playerCard)

	// Update hand
	handCards := field.Objects[1].(*fyne.Container)
	updateHandCards(g, player, handCards)
}

func createPlayerField(g *game.Game, player *player.Player, isBottom bool) *fyne.Container {
	playerCard := createPlayerCard(player)
	cardSpaces := createCardSpaces()
	
	var handCards *fyne.Container
	var deck, discardPile *widget.Button
	
	if isBottom {
		handCards = createHandCards(g, player)
		deck = widget.NewButton(fmt.Sprintf("Deck (%d)", len(g.Deck)), func() {
			showDeckInfo(g)
		})
		discardPile = widget.NewButton(fmt.Sprintf("Discard (%d)", len(g.DiscardPile)), func() {
			showDiscardPileInfo(g)
		})
	} else {
		handCards = container.NewHBox() // Empty container for opponent's hand
		deck = widget.NewButton(fmt.Sprintf("Deck (%d)", len(g.Deck)), nil) // Placeholder for opponent's deck
		discardPile = widget.NewButton(fmt.Sprintf("Discard (%d)", len(g.DiscardPile)), nil) // Placeholder for opponent's discard pile
	}

	field := container.New(layout.NewBorderLayout(nil, handCards, deck, discardPile),
		deck, discardPile, handCards,
		container.NewHBox(cardSpaces[0], playerCard, cardSpaces[1]))

	return field
}

func createPlayerCard(player *player.Player) *fyne.Container {
	nameLabel := widget.NewLabel(player.Name)
	healthLabel := widget.NewLabel(fmt.Sprintf("Health: %d/%d", player.Health, player.MaxHealth))
	manaLabel := widget.NewLabel(fmt.Sprintf("Mana: %d/%d", player.Mana, player.MaxMana))
	armorLabel := widget.NewLabel(fmt.Sprintf("Armor: %d", player.Armor))
	ringLabel := widget.NewLabel(fmt.Sprintf("Ring: %s", getItemName(player.Ring)))
	necklaceLabel := widget.NewLabel(fmt.Sprintf("Necklace: %s", getItemName(player.Necklace)))
	weaponLabel := widget.NewLabel(fmt.Sprintf("Weapon: %s", getItemName(player.Weapon)))
	
	statsButton := widget.NewButton("View Stats", func() {
		showPlayerStats(player)
	})

	return container.NewVBox(
		nameLabel,
		healthLabel,
		manaLabel,
		armorLabel,
		ringLabel,
		necklaceLabel,
		weaponLabel,
		statsButton,
	)
}

func updatePlayerCard(player *player.Player, card *fyne.Container) {
	card.Objects[1].(*widget.Label).SetText(fmt.Sprintf("Health: %d/%d", player.Health, player.MaxHealth))
	card.Objects[2].(*widget.Label).SetText(fmt.Sprintf("Mana: %d/%d", player.Mana, player.MaxMana))
	card.Objects[3].(*widget.Label).SetText(fmt.Sprintf("Armor: %d", player.Armor))
	card.Objects[4].(*widget.Label).SetText(fmt.Sprintf("Ring: %s", getItemName(player.Ring)))
	card.Objects[5].(*widget.Label).SetText(fmt.Sprintf("Necklace: %s", getItemName(player.Necklace)))
	card.Objects[6].(*widget.Label).SetText(fmt.Sprintf("Weapon: %s", getItemName(player.Weapon)))
}

func getItemName(item *player.Item) string {
	if item == nil {
		return "None"
	}
	return item.Name
}

func createCardSpaces() [2]*fyne.Container {
	leftSpace := container.NewVBox()
	rightSpace := container.NewVBox()

	for i := 0; i < 3; i++ {
		leftSpace.Add(createCardSlot())
		rightSpace.Add(createCardSlot())
	}

	return [2]*fyne.Container{leftSpace, rightSpace}
}

func createCardSlot() *fyne.Container {
	slot := canvas.NewRectangle(color.NRGBA{R: 204, G: 204, B: 204, A: 76})
	slot.SetMinSize(fyne.NewSize(100, 150))

	return container.NewMax(slot, widget.NewLabel(""))
}

func createHandCards(g *game.Game, player *player.Player) *fyne.Container {
	handCards := container.NewHBox()

	for i, card := range player.Hand {
		cardButton := widget.NewButton(card.GetInfo(), func(i int) func() {
			return func() {
				playCard(g, player, i)
			}
		}(i))
		handCards.Add(cardButton)
	}

	return handCards
}

func updateHandCards(g *game.Game, player *player.Player, handCards *fyne.Container) {
	handCards.RemoveAll()
	for i, card := range player.Hand {
		cardButton := widget.NewButton(card.GetInfo(), func(i int) func() {
			return func() {
				playCard(g, player, i)
			}
		}(i))
		handCards.Add(cardButton)
	}
}

func playCard(g *game.Game, player *player.Player, cardIndex int) {
	if g.CurrentPlayer == player && g.CurrentPhase == game.PlayPhase {
		g.PlayCard(player, cardIndex)
		showRoundResult(g)

		// Check if the game is over
		if g.CheckGameOver() {
			showGameResult(g)
		}

		// Update UI
		g.UIUpdate()
	}
}

func showRoundResult(g *game.Game) {
	lastPlay := g.LastPlay
	message := fmt.Sprintf(
		"%s played %s\n%s",
		g.CurrentPlayer.Name, lastPlay.PlayerCard.GetInfo(),
		lastPlay.Message,
	)

	dialog := widget.NewLabel(message)
	popUp := widget.NewPopUp(dialog, g.GetWindow().Canvas())
	popUp.Show()
}

func showGameResult(g *game.Game) {
	var winner string
	if g.Player1.Health <= 0 || len(g.Player1.Hand) == 0 {
		winner = g.Player2.Name
	} else if g.Player2.Health <= 0 || len(g.Player2.Hand) == 0 {
		winner = g.Player1.Name
	} else if g.Player1.Score > g.Player2.Score {
		winner = g.Player1.Name
	} else if g.Player2.Score > g.Player1.Score {
		winner = g.Player2.Name
	} else {
		winner = "It's a tie!"
	}

	message := fmt.Sprintf("Game Over!\n%s wins!\n\nFinal Score:\n%s: %d (Health: %d)\n%s: %d (Health: %d)", 
		winner, g.Player1.Name, g.Player1.Score, g.Player1.Health, g.Player2.Name, g.Player2.Score, g.Player2.Health)

	dialog := widget.NewLabel(message)
	popUp := widget.NewPopUp(dialog, g.GetWindow().Canvas())
	popUp.Show()
}

func showDeckInfo(g *game.Game) {
	message := fmt.Sprintf("Remaining cards in deck: %d", len(g.Deck))
	dialog := widget.NewLabel(message)
	popUp := widget.NewPopUp(dialog, g.GetWindow().Canvas())
	popUp.Show()
}

func showDiscardPileInfo(g *game.Game) {
	message := fmt.Sprintf("Cards in discard pile: %d\n\n", len(g.DiscardPile))
	for _, card := range g.DiscardPile {
		message += card.GetInfo() + "\n\n"
	}
	dialog := widget.NewLabel(message)
	popUp := widget.NewPopUp(dialog, g.GetWindow().Canvas())
	popUp.Show()
}

func showPlayerStats(player *player.Player) {
	message := fmt.Sprintf("Player: %s\n\n", player.Name)
	message += fmt.Sprintf("Health: %d/%d\n", player.Health, player.MaxHealth)
	message += fmt.Sprintf("Mana: %d/%d\n", player.Mana, player.MaxMana)
	message += fmt.Sprintf("Armor: %d\n", player.Armor)
	message += fmt.Sprintf("Ring: %s\n", getItemName(player.Ring))
	message += fmt.Sprintf("Necklace: %s\n", getItemName(player.Necklace))
	message += fmt.Sprintf("Weapon: %s\n", getItemName(player.Weapon))
	message += fmt.Sprintf("Total Bonus: %d", player.GetTotalBonus())

	dialog := widget.NewLabel(message)
	popUp := widget.NewPopUp(dialog, fyne.CurrentApp().Driver().AllWindows()[0].Canvas())
	popUp.Show()
}