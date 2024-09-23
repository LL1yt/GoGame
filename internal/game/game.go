package game

import (
	"fmt"
	"math/rand"
	"time"

	"GoGame/internal/card"
	"GoGame/internal/player"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type GamePhase int

const (
	DrawPhase GamePhase = iota
	PlayPhase
	EndPhase
)

type Game struct {
	Player1       player.Player
	Player2       player.Player
	CurrentPlayer *player.Player
	ScoreLabel    *widget.Label
	window        fyne.Window
	LastPlay      PlayResult
	Deck          []card.Card
	DiscardPile   []card.Card
	TurnCount     int
	GameOver      bool
	UIUpdate      func()
	CurrentPhase  GamePhase
	EndTurn       chan bool
}

type PlayResult struct {
	PlayerCard   card.Card
	OpponentCard card.Card
	Message      string
}

func NewGame() *Game {
	player1, player2 := initializePlayers()
	game := &Game{
		Player1:      player1,
		Player2:      player2,
		Deck:         InitializeDeck(),
		CurrentPhase: DrawPhase,
		EndTurn:      make(chan bool),
	}
	game.CurrentPlayer = &game.Player1
	return game
}

func (g *Game) Reset() {
	player1, player2 := initializePlayers()
	g.Player1 = player1
	g.Player2 = player2
	g.CurrentPlayer = &g.Player1
	g.Deck = InitializeDeck()
	g.DiscardPile = []card.Card{}
	g.TurnCount = 0
	g.GameOver = false
	g.CurrentPhase = DrawPhase
	g.EndTurn = make(chan bool)
	g.DealInitialHands()
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

func (g *Game) PlayCard(player *player.Player, cardIndex int) {
	if cardIndex >= len(player.Hand) {
		return
	}

	playerCard := player.Hand[cardIndex]
	player.Hand = append(player.Hand[:cardIndex], player.Hand[cardIndex+1:]...)

	// Play the card
	playerCard.Play(player)

	message := fmt.Sprintf("%s played %s", player.Name, playerCard.GetInfo())
	g.UpdateScore()

	g.LastPlay = PlayResult{
		PlayerCard: playerCard,
		Message:    message,
	}

	// Add played card to discard pile
	g.DiscardPile = append(g.DiscardPile, playerCard)
}

func (g *Game) PlayRandomCard(player *player.Player) {
	if len(player.Hand) == 0 {
		return
	}

	cardIndex := rand.Intn(len(player.Hand))
	g.PlayCard(player, cardIndex)
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

func (g *Game) SwitchTurn() {
	if g.CurrentPlayer == &g.Player1 {
		g.CurrentPlayer = &g.Player2
	} else {
		g.CurrentPlayer = &g.Player1
	}
	g.TurnCount++
	g.CurrentPhase = DrawPhase
}

func (g *Game) GameLoop() {
	for !g.GameOver {
		g.CurrentPhase = DrawPhase
		g.DrawCard(g.CurrentPlayer)
		g.CurrentPhase = PlayPhase

		if g.CurrentPlayer == &g.Player1 {
			// Player 1's turn (human player)
			fmt.Printf("%s's turn\n", g.CurrentPlayer.Name)
			// Wait for the player to end their turn
			<-g.EndTurn
		} else {
			// Player 2's turn (opponent)
			g.PlayRandomCard(g.CurrentPlayer)
			time.Sleep(1 * time.Second) // Simulate opponent thinking
		}

		g.CurrentPhase = EndPhase
		// Update UI
		if g.UIUpdate != nil {
			g.UIUpdate()
		}

		// Check for game over conditions
		if g.CheckGameOver() {
			g.GameOver = true
		} else {
			g.SwitchTurn()
		}
	}

	// Game over
	g.DetermineWinner()
}

func (g *Game) CheckGameOver() bool {
	return len(g.Player1.Hand) == 0 || len(g.Player2.Hand) == 0 || g.Player1.Health <= 0 || g.Player2.Health <= 0 || g.TurnCount >= 20
}

func (g *Game) DetermineWinner() {
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

	fmt.Printf("Game Over! %s wins!\n", winner)
	fmt.Printf("Final Score: %s: %d, %s: %d\n", g.Player1.Name, g.Player1.Score, g.Player2.Name, g.Player2.Score)

	// Update UI with game over state
	if g.UIUpdate != nil {
		g.UIUpdate()
	}
}