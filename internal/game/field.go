package game

import (
	"GoGame/internal/card"
	"GoGame/internal/player"
)

// CardSlot представляет собой место для одной карты на поле
type CardSlot struct {
	Card     *card.Card
	IsOccupied bool
}

// PlayerField представляет половину игрового поля, принадлежащую одному игроку
type PlayerField struct {
	Player        *player.Player
	LeftCards     [3]CardSlot  // Левая сторона от карточки игрока
	RightCards    [3]CardSlot  // Правая сторона от карточки игрока
	PlayerCard    CardSlot     // Карточка самого игрока в центре
}

// GameField представляет всё игровое поле
type GameField struct {
	PlayerField    PlayerField
	OpponentField  PlayerField
	PlayerDeck     []card.Card
	OpponentDeck   []card.Card
	PlayerDiscard  []card.Card
	OpponentDiscard []card.Card
}

// NewGameField создает новое игровое поле
func NewGameField(player, opponent *player.Player) *GameField {
	return &GameField{
		PlayerField: PlayerField{
			Player: player,
		},
		OpponentField: PlayerField{
			Player: opponent,
		},
	}
}

// PlaceCard размещает карту в указанном слоте
func (pf *PlayerField) PlaceCard(card *card.Card, position int, isLeft bool) bool {
	var targetSlots *[3]CardSlot
	if isLeft {
		targetSlots = &pf.LeftCards
	} else {
		targetSlots = &pf.RightCards
	}

	if position < 0 || position >= len(*targetSlots) {
		return false
	}

	if targetSlots[position].IsOccupied {
		return false
	}

	targetSlots[position].Card = card
	targetSlots[position].IsOccupied = true
	return true
}

// RemoveCard удаляет карту из указанного слота
func (pf *PlayerField) RemoveCard(position int, isLeft bool) *card.Card {
	var targetSlots *[3]CardSlot
	if isLeft {
		targetSlots = &pf.LeftCards
	} else {
		targetSlots = &pf.RightCards
	}

	if position < 0 || position >= len(*targetSlots) {
		return nil
	}

	if !targetSlots[position].IsOccupied {
		return nil
	}

	card := targetSlots[position].Card
	targetSlots[position].Card = nil
	targetSlots[position].IsOccupied = false
	return card
}