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

func SetupUI(g *game.Game) {
	g.ScoreLabel = widget.NewLabel(fmt.Sprintf("Score - %s: %d, %s: %d", g.Player1.Name, g.Player1.Score, g.Player2.Name, g.Player2.Score))

	player1Field := createPlayerField(g, &g.Player1, true)
	player2Field := createPlayerField(g, &g.Player2, false)

	gameBoard := container.NewVBox(
		player2Field,
		widget.NewSeparator(),
		player1Field,
	)

	content := container.NewBorder(g.ScoreLabel, nil, nil, nil, gameBoard)

	g.GetWindow().SetContent(content)
	g.GetWindow().Resize(fyne.NewSize(1920, 1080))
}

func createPlayerField(g *game.Game, player *player.Player, isBottom bool) *fyne.Container {
	playerCard := createPlayerCard(player)
	cardSpaces := createCardSpaces(g, player)
	
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

func getItemName(item *player.Item) string {
	if item == nil {
		return "None"
	}
	return item.Name
}

func createCardSpaces(g *game.Game, player *player.Player) [2]*fyne.Container {
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

func playCard(g *game.Game, player *player.Player, cardIndex int) {
	opponent := &g.Player2
	if player == &g.Player2 {
		opponent = &g.Player1
	}

	g.PlayCard(player, opponent, cardIndex)
	showRoundResult(g, player, opponent)

	// Check if the game is over
	if len(player.Hand) == 0 || len(opponent.Hand) == 0 {
		showGameResult(g)
	}

	// Update UI
	g.GetWindow().Content().Refresh()
}

func showRoundResult(g *game.Game, player, opponent *player.Player) {
	lastPlay := g.LastPlay
	message := fmt.Sprintf(
		"%s played %s\n%s played %s\n%s",
		player.Name, lastPlay.PlayerCard.GetInfo(),
		opponent.Name, lastPlay.OpponentCard.GetInfo(),
		lastPlay.Message,
	)

	dialog := widget.NewLabel(message)
	popUp := widget.NewPopUp(dialog, g.GetWindow().Canvas())
	popUp.Show()
}

func showGameResult(g *game.Game) {
	var winner string
	if g.Player1.Score > g.Player2.Score {
		winner = g.Player1.Name
	} else if g.Player2.Score > g.Player1.Score {
		winner = g.Player2.Name
	} else {
		winner = "It's a tie!"
	}

	message := fmt.Sprintf("Game Over!\n%s wins!\n\nFinal Score:\n%s: %d\n%s: %d", 
		winner, g.Player1.Name, g.Player1.Score, g.Player2.Name, g.Player2.Score)

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