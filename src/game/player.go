package game

import (
	"fmt"
)
const (
	MAX_MANA			= 12
	MAX_PLAYERS			= 2
	MIN_PLAYERS			= 2
	MAX_HAND_CARD		= 8
	MAX_BOARD_CARD		= 6


	STARTING_MANA		= 0
	STARTING_LIFE		= 30
	STARTING_CARD		= 30
	STARTING_RUNES		= 25
	STARTING_CARDS		= 30

	STEP_RUNE			= 5

	DRAFT_PICK			= 3
)

/* STARTING */

type Player struct {
	Id    			int
	Deck  			[]*Card
	Life  			int
	Mana  			int
	MaxMana			int
	Board 			[]*Card
	Hand  			[]*Card
	Runes			int
	StackCard 		int
}


func NewPlayer(id, life, max_mana, runes int) *Player {
	return &Player{
		Id:     	id,
		Mana:   	max_mana,
		MaxMana: 	max_mana,
		Life:   	life,
		Deck: 		make([]*Card, 0),
		Board:  	make([]*Card, 0),
		Hand:   	make([]*Card, 0),
		Runes: 		runes,
		StackCard:	0,
	}

}
func (p *Player) Raw() []interface{} {
	return []interface{}{
		p.Life,
		p.Mana,
		len(p.Deck),
		p.Runes,
		p.StackCard,
	}
}
func (p *Player) Copy() *Player {
	new_player := &Player{
		Id:     	p.Id,
		Mana:   	p.Mana,
		MaxMana: 	p.MaxMana,
		Life:   	p.Life,
		Deck: 		make([]*Card, len(p.Deck)),
		Board:  	make([]*Card, len(p.Board)),
		Hand:   	make([]*Card, len(p.Hand)),
		Runes: 		p.Runes,
		StackCard:	p.StackCard,
	}
	for i, _ := range(p.Board) {
		new_player.Board[i] = p.Board[i].Copy()
	}
	for i, _ := range(p.Hand) {
		new_player.Hand[i] = p.Hand[i].Copy()
	}
	for i, _ := range(p.Deck) {
		new_player.Deck[i] = p.Deck[i].Copy()
	}
	return new_player
}
func (p *Player) DeckCount() int {
	return len(p.Deck)
}
func (p *Player) DeckAddCard(c *Card) {
	if p.Deck != nil {
		p.Deck = append(p.Deck, c)
	}
}

func (p *Player) DeckRemoveCard(c *Card) {
	if p.Deck != nil {
		idx := -1
		for i, c1 := range(p.Deck) {
			if c1.CardNumber == c.CardNumber {
				idx = i
				break
			}
		}
		if idx != -1 {
			p.Deck = append(p.Deck[:idx], p.Deck[idx+1:]...)
		}
	}
}
func (p *Player) SetMaxMana(mana int) {
	if mana <= MAX_MANA && mana >= 0 {
		p.MaxMana = mana
	}
}
func (p *Player) SetLife(life int) {
	p.Life = life
}
func (p *Player) DrawN(n int) (err error) {
	for i := 0 ; i < n ; i++ {
		err = p.Draw()
		if err != nil {
			return err
		}
	}
	return nil
}
func (p *Player) Draw() (error) {

	if len(p.Hand) >= MAX_HAND_CARD {
		//fmt.Println("[GAME][DECK] Maximum card hand reach", MAX_HAND_CARD)
		return fmt.Errorf("[GAME][DECK] Maximum card hand reach %d", MAX_HAND_CARD)
	}
	if len(p.Deck) == 0 {
		p.LifeToNextRune()
		return fmt.Errorf("[GAME][DECK] No more card in the deck")
	}
	card := p.Deck[0]
	p.Deck = p.Deck[1:]
	p.Hand = append(p.Hand, card)

	return nil
}
func (p *Player) StackDraw() {
	p.StackCard++
}
func (p *Player) StackDrawN(n int) {
	for i := 0 ; i < n ; i++ {
		p.StackDraw()
	}
}
func (p *Player) DrawStackCards() (err error) {
	max := p.StackCard
	for i := 0 ; i < max ; i++ {
		err = p.Draw()
		if err != nil {
			return err
		}
		p.StackCard--
	}
	return nil

}
func (p *Player) LifeToNextRune() {
	if p.Runes > 0 {
		if p.Life > p.Runes {
			p.ReceiveDamage(p.Life - p.Runes)
		} else {
			for ; p.Life < p.Runes && p.Runes >= 0 ; p.Runes -= STEP_RUNE {
				damage := p.Runes - p.Life
				if damage > STEP_RUNE {
					damage = STEP_RUNE
				}
				p.ReceiveDamage(damage)
			}
		}
	} else {
		p.ReceiveDamage(p.Life)
	}

}
func (p *Player) Pick(c *Card) {
	if p.HandGetCard(c.Id) == nil {
		if len(p.Hand) < MAX_HAND_CARD {
			p.Hand = append(p.Hand, c)
		} else {
			//fmt.Println("GAME: Max hand card reach", MAX_HAND_CARD)
		}

	} else {
		//fmt.Println("GAME: Card", c, "already exist in", p.Hand)
	}
}
func (p *Player) HandGetCard(id int) *Card {
	for _, c := range(p.Hand) {
		if c.Id == id {
			return c
		}
	}
	return nil
}
func (p *Player) HandRemoveCard(id int) (error) {
	idx := -1
	for i, c := range(p.Hand) {
		if c.Id == id {
			idx = i
			break
		}
	}
	if idx == -1 {
		return fmt.Errorf("RemoveCard: Hand with id %i not in hand", idx)
	}

	l := len(p.Hand)
	p.Hand[idx] = p.Hand[l - 1]
	if l >= 2 {
		p.Hand = p.Hand[:l - 1]
	} else if l == 1 {
		p.Hand = make([]*Card, 0)
	}
	return nil
}
func (p *Player) HandPlayCard(id int) (error) {
	c := p.HandGetCard(id)
	if c == nil {
		return fmt.Errorf("[PLAYER] Card with id %d is not present in hand of player %d", id, p.Id)
	}

	if c.Cost > p.Mana {
		return fmt.Errorf("[PLAYER] No more mana (%d) for playing card %d with cost %d", p.Mana, id, c.Cost)
	}

	p.Mana -= c.Cost

	switch c.Type {
	case CARD_TYPE_CREATURE:
		if len(p.Board) < MAX_BOARD_CARD {
			if c.IsAbleTo(CARD_ABILITY_CHARGE) {
				c.Charge = 1
			}
			p.Board = append(p.Board, c)
		} else {
			return fmt.Errorf("Max card in board reach")
		}
	case CARD_TYPE_ITEM_BLUE:
	case CARD_TYPE_ITEM_GREEN:
	case CARD_TYPE_ITEM_RED:

	default:
		return fmt.Errorf("Unkow card type %d for player %d", c.Type, p.Id)
	}

	////////fmt.Fprintln(os.Stderr, "[GAME][PLAYER] Player", p.Id, "play card", c)
	p.HandRemoveCard(id)
	if c.CardDraw > 0 {
		p.StackCard += c.CardDraw
	}
	return nil
}
func (p *Player) BoardGetCard(id int) *Card {
	for _, c := range(p.Board) {
		if c.Id == id {
			return c
		}
	}
	return nil
}
type FuncCardFilter func (c *Card) bool

func (p *Player) BoardGet(f FuncCardFilter) []*Card {
	ids := make([]*Card, 0)
	for _, c := range(p.Board) {
		if f(c) {
			ids = append(ids, c)
		}
	}
	return ids
}
func (p *Player) BoardGetGuardsId() []*Card {
	return p.BoardGet(func (c *Card) bool {
		return c.IsAbleTo(CARD_ABILITY_GUARD)
	})
}
func (p *Player) BoardRemoveCard(id int) error {
	idx := -1
	for i, b := range p.Board {
		if b.Id == id {
			idx = i
			break
		}
	}

	if idx == -1 {
		return fmt.Errorf("RemoveCard: Hand with id %i not in hand", idx)
	}

	l := len(p.Board)
	if l >= 2 {
		p.Board[idx] = p.Board[l - 1]
		p.Board = p.Board[:l - 1]
	} else if l == 1 {
		p.Board = make([]*Card, 0)
	}

	//fmt.Println("[GAME][DAMAGE] Monster", id, "has been killed")
	return nil

}
func (p *Player) ReloadMana() {
	p.Mana = p.MaxMana
}
func (p *Player) UpdateBoard() {
	for _, c := range(p.Board) {
		c.Attacked 	= false
		c.Charge 	= 1
	}
}
func (p *Player) GainLife(life int) {
	if life > 0 {
		////////fmt.Fprintln(os.Stderr, "[GAME][HEALTH] Player", p.Id, "Gain", life, "life")
		p.SetLife(p.Life + life)
	} else if life < 0 {
		p.ReceiveDamage(-life)
	}
}
func (p *Player) ReceiveDamage(damage int) {
	if damage > 0 {
		////////fmt.Fprintln(os.Stderr, "[GAME][HEALTH] Player", p.Id, "Receive", damage, "damage")
		p.SetLife(p.Life - damage)
	}
}
func (p *Player) IncreaseMana() {
	if p.MaxMana < MAX_MANA {
		p.MaxMana++
	}
}
