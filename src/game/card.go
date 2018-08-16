
package game

import (
	"strings"
	"fmt"
)

const (
	CARD_TYPE_CREATURE	        = 0
	CARD_TYPE_ITEM_GREEN	    = 1
	CARD_TYPE_ITEM_RED	        = 2
	CARD_TYPE_ITEM_BLUE	        = 3

	CARD_ABILITY_BREAKTHROUGH   = 0x100000
	CARD_ABILITY_CHARGE         = 0x010000
	CARD_ABILITY_DRAIN       	= 0x001000
	CARD_ABILITY_GUARD          = 0x000100
	CARD_ABILITY_LETHAL         = 0x000010
	CARD_ABILITY_WARD           = 0x000001

	MAX_ABILITIES				= 6
)


type Card struct {
	CardNumber 				int
	Id    					int
	Location				int
	Type  					int
	Cost  					int
	Attack					int
	Defense					int
	Abilities				int
	HealthChange			int
	OpponentHealthChange 	int
	CardDraw				int
	Charge					int
	Attacked				bool
}

func NewCard(cardNumber int,
			id int,
			cost int,
			attack int,
			defense int,
			abilities string,
			heroHealthChange int,
			opponentHealthChange int,
			cardDraw int,
			type_ int, 
) *Card {
	new_card := &Card{
		CardNumber: 			cardNumber,
		Id:             		id,
		Location:				0,
		Type:           		type_,
		Cost:           		cost,
		Attack:         		attack,
		Defense:        		defense,
		Abilities:      		0,
		HealthChange: 			heroHealthChange,
		OpponentHealthChange: 	opponentHealthChange,
		CardDraw:				cardDraw,
		Charge: 				0,
		Attacked: 				false,
	}
	for _, c := range(strings.Split(abilities, "")) {
		switch c  {
		case "W":
			new_card.EnableAbility(CARD_ABILITY_WARD)
		case "D":
			new_card.EnableAbility(CARD_ABILITY_DRAIN)
		case "L":
			new_card.EnableAbility(CARD_ABILITY_LETHAL)
		case "B":
			new_card.EnableAbility(CARD_ABILITY_BREAKTHROUGH)
		case "C":
			new_card.EnableAbility(CARD_ABILITY_CHARGE)
		case "G":
			new_card.EnableAbility(CARD_ABILITY_GUARD)
		}
	}
	return new_card
}

func (c *Card) Raw() []interface{} {
	return []interface{} {
		c.CardNumber,
		c.Id,
		c.Location,
		c.Type,
		c.Cost,
		c.Attack,
		c.Defense,
		c.Abilities,
		c.HealthChange,
		c.OpponentHealthChange,
		c.CardDraw,
	}
}
func (c *Card) Copy() *Card {
	new_c := &Card{
		CardNumber: 	c.CardNumber,
		Id:             c.Id,
		Location:		c.Location,
		Type:           c.Type,
		Cost:           c.Cost,
		Attack:         c.Attack,
		Defense:        c.Defense,
		HealthChange: 	c.HealthChange,
		Abilities:		c.Abilities,
		OpponentHealthChange: c.OpponentHealthChange,
		CardDraw:		c.CardDraw,
		Charge:			c.Charge,
		Attacked:		c.Attacked,
	}

	return new_c
}

func (c *Card) IsAbleTo(ability int) bool {
	return c.Abilities & ability > 0
}
func (c *Card) DisableAbility(ability int) {
	c.Abilities ^= ability
}
func (c *Card) EnableAbility(ability int) {
	c.Abilities |= ability
}

func typeToString(ct int) string {
	switch ct {
	case CARD_TYPE_CREATURE:
		return "C"
	case CARD_TYPE_ITEM_BLUE:
		return "IB"
	case CARD_TYPE_ITEM_GREEN:
		return "IG"
	case CARD_TYPE_ITEM_RED:
		return "IR"
	}
	return ""
}
func abilitiesToString(a int) string {
	str := make([]string, MAX_ABILITIES)
	if a & CARD_ABILITY_BREAKTHROUGH > 0 { str[0] = "B" } else { str[0] = "-" }
	if a & CARD_ABILITY_CHARGE > 0 { str[1] = "C" } else { str[1] = "-" }
	if a & CARD_ABILITY_DRAIN > 0 { str[2] = "D" } else { str[2] = "-" }
	if a & CARD_ABILITY_GUARD > 0 { str[3] = "G" } else { str[3] = "-" }
	if a & CARD_ABILITY_LETHAL > 0 { str[4] = "L" } else { str[4] = "-" }
	if a & CARD_ABILITY_WARD > 0 { str[5] = "W" } else { str[5] = "-" }

	return strings.Join(str, "")
}

func (c *Card) toString() string {
	str := fmt.Sprintf("%d %d %s %d %d %d %s %d %d %d %d %t",
						c.CardNumber,
						c.Id,
						typeToString(c.Type),
						c.Cost,
						c.Attack,
						c.Defense,
						abilitiesToString(c.Abilities),
						c.HealthChange,
						c.OpponentHealthChange,
						c.CardDraw,
						c.Charge,
						c.Attacked)
	return str
}
