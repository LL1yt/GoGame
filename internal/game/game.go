package game

import (
	"fmt"
	"math/rand"

	"GoGame/internal/card"
	"GoGame/internal/player"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type Game struct {
	Player1     player.Player
	Player2     player.Player
	ScoreLabel  *widget.Label
	window      fyne.Window
	LastPlay    PlayResult
	Deck        []card.Card
	DiscardPile []card.Card
}

type PlayResult struct {
	PlayerCard   card.Card
	OpponentCard card.Card
	Message      string
}

func NewGame() *Game {
	player1, player2 := initializePlayers()
	
	return &Game{
		Player1: player1,
		Player2: player2,
		Deck:    InitializeDeck(),
	}
}

func (g *Game) SetWindow(w fyne.Window) {
	g.window = w
}

func (g *Game) GetWindow() fyne.Window {
	return g.window
}

func (g *Game) UpdateScore() {
	g.ScoreLabel.SetText(fmt.Sprintf("Score - %s: %d, %s: %d", g.Player1.Name, g.Player1.Score, g.Player2.Name, g.Player2.Score))
}

func (g *Game) PlayCard(player, opponent *player.Player, cardIndex int) {
	if cardIndex >= len(player.Hand) {
		return
	}

	playerCard := player.Hand[cardIndex]
	player.Hand = append(player.Hand[:cardIndex], player.Hand[cardIndex+1:]...)

	opponentCardIndex := rand.Intn(len(opponent.Hand))
	opponentCard := opponent.Hand[opponentCardIndex]
	opponent.Hand = append(opponent.Hand[:opponentCardIndex], opponent.Hand[opponentCardIndex+1:]...)

	// Play the cards
	playerCard.Play(player)
	opponentCard.Play(opponent)

	message := g.DetermineRoundWinner(player, opponent, playerCard, opponentCard)
	g.UpdateScore()

	g.LastPlay = PlayResult{
		PlayerCard:   playerCard,
		OpponentCard: opponentCard,
		Message:      message,
	}

	// Add played cards to discard pile
	g.DiscardPile = append(g.DiscardPile, playerCard, opponentCard)

	// Draw new cards
	g.DrawCard(player)
	g.DrawCard(opponent)
}

func (g *Game) DetermineRoundWinner(player, opponent *player.Player, playerCard, opponentCard card.Card) string {
	playerPower := playerCard.Power + player.GetTotalBonus()
	opponentPower := opponentCard.Power + opponent.GetTotalBonus()

	if playerPower > opponentPower {
		player.Score++
		return fmt.Sprintf("%s won the round!", player.Name)
	} else if opponentPower > playerPower {
		opponent.Score++
		return fmt.Sprintf("%s won the round!", opponent.Name)
	}
	return "The round ended in a tie!"
}

func InitializeDeck() []card.Card {
	deck := []card.Card{
		card.CreateBasicUnitCard("Soldier", 1),
		card.CreateBasicUnitCard("Archer", 2),
		card.CreateBasicUnitCard("Knight", 3),
		card.CreateBasicUnitCard("Mage", 4),
		card.CreateBasicUnitCard("Dragon", 5),
		card.CreateBasicUnitCard("Hero", 6),
		card.CreateBasicUnitCard("Commander", 7),
		card.CreateBasicUnitCard("Wizard", 8),
		card.CreateBasicUnitCard("Titan", 9),
		card.CreateBasicUnitCard("Legend", 10),
	}

	// Add some spell cards
	deck = append(deck, card.CreateSpellCard("Fireball", "Deal 3 damage to the opponent", func(target interface{}) {
		if player, ok := target.(*player.Player); ok {
			player.TakeDamage(3)
		}
	}))

	deck = append(deck, card.CreateSpellCard("Heal", "Restore 3 health", func(target interface{}) {
		if player, ok := target.(*player.Player); ok {
			player.Heal(3)
		}
	}))

	// Add some item cards
	deck = append(deck, card.CreateItemCard("Shield", 1, "Increase armor by 2", func(target interface{}) {
		if player, ok := target.(*player.Player); ok {
			player.AddArmor(2)
		}
	}))

	return deck
}

func initializePlayers() (player.Player, player.Player) {
	player1 := player.NewPlayer("Player 1")
	player2 := player.NewPlayer("Player 2")

	return *player1, *player2
}

func (g *Game) DrawCard(player *player.Player) {
	if len(g.Deck) == 0 {
		g.ShuffleDiscardPileToDeck()
	}
	if len(g.Deck) > 0 {
		card := g.Deck[0]
		g.Deck = g.Deck[1:]
		player.Hand = append(player.Hand, card)
	}
}

func (g *Game) ShuffleDiscardPileToDeck() {
	g.Deck = append(g.Deck, g.DiscardPile...)
	g.DiscardPile = []card.Card{}
	rand.Shuffle(len(g.Deck), func(i, j int) {
		g.Deck[i], g.Deck[j] = g.Deck[j], g.Deck[i]
	})
}

func (g *Game) DealInitialHands() {
	for i := 0; i < 5; i++ {
		g.DrawCard(&g.Player1)
		g.DrawCard(&g.Player2)
	}
}