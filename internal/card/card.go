package card

import (
	"fmt"
)

type CardType int

const (
	UnitCard CardType = iota
	SpellCard
	ItemCard
)

type Card struct {
	Name        string
	Power       int
	Type        CardType
	Description string
	Effect      func(interface{}) // This will be used to implement special card effects
}

func (c *Card) Play(target interface{}) {
	fmt.Printf("Playing card: %s\n", c.Name)
	if c.Effect != nil {
		c.Effect(target)
	}
}

func (c *Card) GetInfo() string {
	return fmt.Sprintf("%s (Power: %d)\nType: %s\n%s", c.Name, c.Power, c.getTypeString(), c.Description)
}

func (c *Card) getTypeString() string {
	switch c.Type {
	case UnitCard:
		return "Unit"
	case SpellCard:
		return "Spell"
	case ItemCard:
		return "Item"
	default:
		return "Unknown"
	}
}

// CreateBasicUnitCard creates a basic unit card with no special effects
func CreateBasicUnitCard(name string, power int) Card {
	return Card{
		Name:        name,
		Power:       power,
		Type:        UnitCard,
		Description: fmt.Sprintf("A basic unit with %d power.", power),
	}
}

// CreateSpellCard creates a spell card with a custom effect
func CreateSpellCard(name string, description string, effect func(interface{})) Card {
	return Card{
		Name:        name,
		Power:       0,
		Type:        SpellCard,
		Description: description,
		Effect:      effect,
	}
}

// CreateItemCard creates an item card with a custom effect
func CreateItemCard(name string, power int, description string, effect func(interface{})) Card {
	return Card{
		Name:        name,
		Power:       power,
		Type:        ItemCard,
		Description: description,
		Effect:      effect,
	}
}