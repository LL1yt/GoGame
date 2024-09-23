package player

import (
	"GoGame/internal/card"
)

// Item представляет предмет экипировки
type Item struct {
	Name        string
	Description string
	Bonus       int
}

// Player представляет игрока
type Player struct {
	Name      string
	Hand      []card.Card
	Score     int
	Health    int
	MaxHealth int
	Mana      int
	MaxMana   int
	Armor     int
	Ring      *Item
	Necklace  *Item
	Weapon    *Item
}

// NewPlayer создает нового игрока
func NewPlayer(name string) *Player {
	return &Player{
		Name:      name,
		Health:    100,
		MaxHealth: 100,
		Mana:      10,
		MaxMana:   10,
	}
}

// EquipItem экипирует предмет в соответствующий слот
func (p *Player) EquipItem(item *Item, slot string) bool {
	switch slot {
	case "ring":
		p.Ring = item
	case "necklace":
		p.Necklace = item
	case "weapon":
		p.Weapon = item
	default:
		return false
	}
	return true
}

// UnequipItem снимает предмет из указанного слота
func (p *Player) UnequipItem(slot string) *Item {
	var item *Item
	switch slot {
	case "ring":
		item = p.Ring
		p.Ring = nil
	case "necklace":
		item = p.Necklace
		p.Necklace = nil
	case "weapon":
		item = p.Weapon
		p.Weapon = nil
	}
	return item
}

// AddArmor добавляет броню игроку
func (p *Player) AddArmor(amount int) {
	p.Armor += amount
	if p.Armor < 0 {
		p.Armor = 0
	}
}

// TakeDamage наносит урон игроку
func (p *Player) TakeDamage(amount int) {
	if amount <= p.Armor {
		p.Armor -= amount
	} else {
		remainingDamage := amount - p.Armor
		p.Armor = 0
		p.Health -= remainingDamage
		if p.Health < 0 {
			p.Health = 0
		}
	}
}

// Heal восстанавливает здоровье игрока
func (p *Player) Heal(amount int) {
	p.Health += amount
	if p.Health > p.MaxHealth {
		p.Health = p.MaxHealth
	}
}

// UseMana использует ману игрока
func (p *Player) UseMana(amount int) bool {
	if p.Mana >= amount {
		p.Mana -= amount
		return true
	}
	return false
}

// RestoreMana восстанавливает ману игрока
func (p *Player) RestoreMana(amount int) {
	p.Mana += amount
	if p.Mana > p.MaxMana {
		p.Mana = p.MaxMana
	}
}

// GetTotalBonus возвращает сумму бонусов от всех экипированных предметов
func (p *Player) GetTotalBonus() int {
	bonus := 0
	if p.Ring != nil {
		bonus += p.Ring.Bonus
	}
	if p.Necklace != nil {
		bonus += p.Necklace.Bonus
	}
	if p.Weapon != nil {
		bonus += p.Weapon.Bonus
	}
	return bonus
}