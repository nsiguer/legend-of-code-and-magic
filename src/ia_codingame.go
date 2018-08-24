package main

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"os"
	"reflect"
	"strings"
	"time"
)
const (
	MAX_MANA                  = 12
	MAX_PLAYERS               = 2
	MIN_PLAYERS               = 2
	MAX_HAND_CARD             = 8
	MAX_BOARD_CARD            = 6
	STARTING_MANA             = 0
	STARTING_LIFE             = 30
	STARTING_CARD             = 30
	STARTING_RUNES            = 25
	STARTING_CARDS            = 30
	STEP_RUNE                 = 5
	DRAFT_PICK                = 3
	WEIGHT_BREAKTHROUGH       = 1
	WEIGHT_CHARGE             = 0.5
	WEIGHT_DRAIN              = 1.5
	WEIGHT_GUARD              = 2
	WEIGHT_LETHAL             = 2
	WEIGHT_WARD               = 1
	MOVE_PASS                 = 0
	MOVE_PICK                 = 1
	MOVE_SUMMON               = 2
	MOVE_ATTACK               = 3
	MOVE_USE                  = 4
	OUTCOME_WIN               = 100
	OUTCOME_LOSE              = -100
	CARD_TYPE_CREATURE        = 0
	CARD_TYPE_ITEM_GREEN      = 1
	CARD_TYPE_ITEM_RED        = 2
	CARD_TYPE_ITEM_BLUE       = 3
	CARD_ABILITY_BREAKTHROUGH = 0x100000
	CARD_ABILITY_CHARGE       = 0x010000
	CARD_ABILITY_DRAIN        = 0x001000
	CARD_ABILITY_GUARD        = 0x000100
	CARD_ABILITY_LETHAL       = 0x000010
	CARD_ABILITY_WARD         = 0x000001
	MAX_ABILITIES             = 6
	BIAS_PARAMETER            = 0.7
	MCTS_ITERATION            = 150
	MCTS_SIMULATION           = 100
	MCTS_TIMEOUT              = 65
)
/* STARTING */

func mainCG() {

	var hero, vilain *Player
	var phero, pvilain *Player
	var step_draft bool
	var previous_board []int = make([]int, 0)
	var round int = 0

	ai := NewAI()
	ai.LoadAgentMCTS()
	ai.LoadAgentRandom()
	ai.LoadAgentDraft()

	for {

		var players []*Player

		players = make([]*Player, 2)
		step_draft = false

		for i := 0; i < 2; i++ {
			var playerHealth, playerMana, playerDeck, playerRune int
			fmt.Scan(&playerHealth, &playerMana, &playerDeck, &playerRune)

			players[i] = NewPlayer(i, playerHealth, playerMana, playerRune)
		}

		hero	= players[0]
		vilain	= players[1]

		if phero != nil {
			hero.Deck = phero.Deck
			vilain.Deck = pvilain.Deck
		}

		var opponentHand int
		fmt.Scan(&opponentHand)

		var cardCount int
		fmt.Scan(&cardCount)

		var draft []*Card = make([]*Card, 0)

		for i := 0; i < cardCount; i++ {
			var cardNumber, instanceId, location, cardType, cost, attack, defense int
			var abilities string
			var myHealthChange, opponentHealthChange, cardDraw int

			var card *Card

			fmt.Scan(&cardNumber, &instanceId, &location, &cardType, &cost, &attack, &defense, &abilities, &myHealthChange, &opponentHealthChange, &cardDraw)

			card = NewCard(cardNumber, instanceId, cost, attack, defense, abilities, myHealthChange, opponentHealthChange, cardDraw, cardType)
			switch cardType {
			case CARD_TYPE_CREATURE:
				if card.IsAbleTo(CARD_ABILITY_CHARGE) {
					card.Charge = 1
				}
			}
			//fmt.Fprintln(os.Stderr, card)

			// Draft step
			if instanceId == -1 {
				step_draft = true
				draft = append(draft, card)

			} else {
				step_draft = false
				switch location {
				case -1:
					vilain.Board = append(vilain.Board, card)
				case 0:
					hero.Pick(card)
				case 1:
					if len(previous_board) > 0 {
						exist, _ := in_array(card.Id, previous_board)
						if exist { card.Charge = 1 }
					}
					hero.Board = append(hero.Board, card)
				}
			}
		}

		// Remove new card draw from the deck
		for _, c1 := range(hero.Hand) {
			if phero != nil {
				found := false
				for _, c2 := range(phero.Hand) {
					if c1.Id == c2.Id {
						found = true
						break
					}
				}
				if ! found {
					hero.DeckRemoveCard(c1)
				}
			} else {
				hero.DeckRemoveCard(c1)
			}
		}

		for _, c1 := range(vilain.Board) {
			if pvilain != nil {
				found := false
				for _, c2 := range(pvilain.Board) {
					if c1.Id == c2.Id {
						found = true
						break
					}
				}
				if ! found {
					vilain.DeckRemoveCard(c1)
				}
			} else {
				vilain.DeckRemoveCard(c1)
			}
		}


		init_state := NewState(hero, vilain)

		if step_draft {
			init_state.Draft = draft
			moves := ai.Think("DRAFT", init_state)
			str_moves := make([]string, 0)
			for _, m := range(moves) {
				fmt.Fprintln(os.Stderr, "Move:", m.ToString())
				str_moves = append(str_moves, m.ToString())
				card_pick := draft[m.Params[0]].Copy()
				hero.DeckAddCard(card_pick)
				vilain.DeckAddCard(card_pick)
			}
			fmt.Println(strings.Join(str_moves, ";"))
			//fmt.Println("PASS")
			//hero.Pick(draft)
		} else {

			var moves []*Move;
			init_state.Print()
			hero_state := init_state.Hero()
			moves = ai.Think("MCTS", init_state)
			str_moves := make([]string, 0)
			for _, m := range(moves) {
				init_state.Move(m)
				fmt.Fprintln(os.Stderr, "Move:", m.ToString())
				str_moves = append(str_moves, m.ToString())
			}
			fmt.Println(strings.Join(str_moves, ";"))

			previous_board = make([]int, len(hero_state.Board))
			for i, c := range(hero_state.Board) {
				previous_board[i] = c.Id
			}
		}

		phero	= hero
		pvilain = vilain
		round++
	}
}


/* STARTING */

func in_array(v interface{}, in interface{}) (ok bool, i int) {
	val := reflect.Indirect(reflect.ValueOf(in))
	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		for ; i < val.Len(); i++ {
			if ok = v == val.Index(i).Interface(); ok {
				return
			}
		}
	}
	return
}


func iterate_combinations(elems []*Move) [][]*Move{
	moves := make([][]*Move, 0)
	n := len(elems)
	for num:=0;num < (1 << uint(n));num++ {
		combination := []*Move{}
		for ndx:=0;ndx<n;ndx++ {
			// (is the bit "on" in this number?)
			if num & (1 << uint(ndx)) != 0 {
				// (then add it to the combination)
				combination = append(combination, elems[ndx])
			}
		}
		if len(combination) > 0 {
			moves = append(moves, combination)
		
		}
	}
	return moves
}

func permutations(arr []*Move) ([][]*Move){
    var helper func([]*Move, int)
    res := make([][]*Move, 0)

    helper = func(arr []*Move, n int){
        if n == 1{
            tmp := make([]*Move, len(arr))
            copy(tmp, arr)
            res = append(res, tmp)
        } else {
			
            for i := 0; i < n; i++{
                helper(arr, n - 1)
                if n % 2 == 1{
                    tmp := arr[i]
                    arr[i] = arr[n - 1]
                    arr[n - 1] = tmp
                } else {
                    tmp := arr[0]
                    arr[0] = arr[n - 1]
                    arr[n - 1] = tmp
                }
            }
        }
    }
    helper(arr, len(arr))
    return res
}

func del_in_int_array(v int, a []int) (ok bool, r []int) {
    idx := -1
    for i, c := range(a) { if v == c { idx = i ; break } }
    if idx != -1 {
        r = append(a[:idx], a[idx+1:]...)
    }
    return
}
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
/* STARTING */

var CARDS = []*Card{
	NewCard(1, -1, 1, 2, 1, "------", 1, 0, 0, CARD_TYPE_CREATURE),
	NewCard(2, -1, 1, 1, 2, "------", 0, -1, 0, CARD_TYPE_CREATURE),
	NewCard(3, -1, 1, 2, 2, "------", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(4, -1, 2, 1, 5, "------", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(5, -1, 2, 4, 1, "------", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(6, -1, 2, 3, 2, "------", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(7, -1, 2, 2, 2, "-----W", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(8, -1, 2, 2, 3, "------", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(9, -1, 3, 3, 4, "------", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(10, -1, 3, 3, 1, "--D---", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(11, -1, 3, 5, 2, "------", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(12, -1, 3, 2, 5, "------", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(13, -1, 4, 5, 3, "------", 1, -1, 0, CARD_TYPE_CREATURE),
	NewCard(14, -1, 4, 9, 1, "------", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(15, -1, 4, 4, 5, "------", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(16, -1, 4, 6, 2, "------", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(17, -1, 4, 4, 5, "------", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(18, -1, 4, 7, 4, "------", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(19, -1, 5, 5, 6, "------", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(20, -1, 5, 8, 2, "------", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(21, -1, 5, 6, 5, "------", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(22, -1, 6, 7, 5, "------", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(23, -1, 7, 8, 8, "------", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(24, -1, 1, 1, 1, "------", 0, -1, 0, CARD_TYPE_CREATURE),
	NewCard(25, -1, 2, 3, 1, "------", -2, -2, 0, CARD_TYPE_CREATURE),
	NewCard(26, -1, 2, 3, 2, "------", 0, -1, 0, CARD_TYPE_CREATURE),
	NewCard(27, -1, 2, 2, 2, "------", 2, 0, 0, CARD_TYPE_CREATURE),
	NewCard(28, -1, 2, 1, 2, "------", 0, 0, 1, CARD_TYPE_CREATURE),
	NewCard(29, -1, 2, 2, 1, "------", 0, 0, 1, CARD_TYPE_CREATURE),
	NewCard(30, -1, 3, 4, 2, "------", 0, -2, 0, CARD_TYPE_CREATURE),
	NewCard(31, -1, 3, 3, 1, "------", 0, -1, 0, CARD_TYPE_CREATURE),
	NewCard(32, -1, 3, 3, 2, "------", 0, 0, 1, CARD_TYPE_CREATURE),
	NewCard(33, -1, 4, 4, 3, "------", 0, 0, 1, CARD_TYPE_CREATURE),
	NewCard(34, -1, 5, 3, 5, "------", 0, 0, 1, CARD_TYPE_CREATURE),
	NewCard(35, -1, 6, 5, 2, "B-----", 0, 0, 1, CARD_TYPE_CREATURE),
	NewCard(36, -1, 6, 4, 4, "------", 0, 0, 2, CARD_TYPE_CREATURE),
	NewCard(37, -1, 6, 5, 7, "------", 0, 0, 1, CARD_TYPE_CREATURE),
	NewCard(38, -1, 1, 1, 3, "--D---", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(39, -1, 1, 2, 1, "--D---", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(40, -1, 3, 2, 3, "--DG--", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(41, -1, 3, 2, 2, "-CD---", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(42, -1, 4, 4, 2, "--D---", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(43, -1, 6, 5, 5, "--D---", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(44, -1, 6, 3, 7, "--D-L-", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(45, -1, 6, 6, 5, "B-D---", -3, 0, 0, CARD_TYPE_CREATURE),
	NewCard(46, -1, 9, 7, 7, "--D---", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(47, -1, 2, 1, 5, "--D---", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(48, -1, 1, 1, 1, "----L-", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(49, -1, 2, 1, 2, "---GL-", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(50, -1, 3, 3, 2, "----L-", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(51, -1, 4, 3, 5, "----L-", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(52, -1, 4, 2, 4, "----L-", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(53, -1, 4, 1, 1, "-C--L-", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(54, -1, 3, 2, 2, "----L-", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(55, -1, 2, 0, 5, "---G--", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(56, -1, 4, 2, 7, "------", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(57, -1, 4, 1, 8, "------", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(58, -1, 6, 5, 6, "B-----", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(59, -1, 7, 7, 7, "------", 1, -1, 0, CARD_TYPE_CREATURE),
	NewCard(60, -1, 7, 4, 8, "------", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(61, -1, 9, 10, 10, "------", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(62, -1, 12, 12, 12, "B--G--", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(63, -1, 2, 0, 4, "---G-W", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(64, -1, 2, 1, 1, "---G-W", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(65, -1, 2, 2, 2, "-----W", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(66, -1, 5, 5, 1, "-----W", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(67, -1, 6, 5, 5, "-----W", 0, -2, 0, CARD_TYPE_CREATURE),
	NewCard(68, -1, 6, 7, 5, "-----W", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(69, -1, 3, 4, 4, "B-----", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(70, -1, 4, 6, 3, "B-----", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(71, -1, 4, 3, 2, "BC----", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(72, -1, 4, 5, 3, "B-----", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(73, -1, 4, 4, 4, "B-----", 4, 0, 0, CARD_TYPE_CREATURE),
	NewCard(74, -1, 5, 5, 4, "B--G--", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(75, -1, 5, 6, 5, "B-----", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(76, -1, 6, 5, 5, "B-D---", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(77, -1, 7, 7, 7, "B-----", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(78, -1, 8, 5, 5, "B-----", 0, -5, 0, CARD_TYPE_CREATURE),
	NewCard(79, -1, 8, 8, 8, "B-----", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(80, -1, 8, 8, 8, "B--G--", 0, 0, 1, CARD_TYPE_CREATURE),
	NewCard(81, -1, 9, 6, 6, "BC----", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(82, -1, 7, 5, 5, "B-D--W", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(83, -1, 0, 1, 1, "-C----", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(84, -1, 2, 1, 1, "-CD--W", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(85, -1, 3, 2, 3, "-C----", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(86, -1, 3, 1, 5, "-C----", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(87, -1, 4, 2, 5, "-C-G--", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(88, -1, 5, 4, 4, "-C----", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(89, -1, 5, 4, 1, "-C----", 2, 0, 0, CARD_TYPE_CREATURE),
	NewCard(90, -1, 8, 5, 5, "-C----", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(91, -1, 0, 1, 2, "---G--", 0, 1, 0, CARD_TYPE_CREATURE),
	NewCard(92, -1, 1, 0, 1, "---G--", 2, 0, 0, CARD_TYPE_CREATURE),
	NewCard(93, -1, 1, 2, 1, "---G--", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(94, -1, 2, 1, 4, "---G--", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(95, -1, 2, 2, 3, "---G--", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(96, -1, 2, 3, 2, "---G--", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(97, -1, 3, 3, 3, "---G--", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(98, -1, 3, 2, 4, "---G--", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(99, -1, 3, 2, 5, "---G--", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(100, -1, 3, 1, 6, "---G--", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(101, -1, 4, 3, 4, "---G--", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(102, -1, 4, 3, 3, "---G--", 0, -1, 0, CARD_TYPE_CREATURE),
	NewCard(103, -1, 4, 3, 6, "---G--", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(104, -1, 4, 4, 4, "---G--", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(105, -1, 5, 4, 6, "---G--", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(106, -1, 5, 5, 5, "---G--", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(107, -1, 5, 3, 3, "---G--", 3, 0, 0, CARD_TYPE_CREATURE),
	NewCard(108, -1, 5, 2, 6, "---G--", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(109, -1, 5, 5, 6, "------", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(110, -1, 5, 0, 9, "---G--", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(111, -1, 6, 6, 6, "---G--", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(112, -1, 6, 4, 7, "---G--", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(113, -1, 6, 2, 4, "---G--", 4, 0, 0, CARD_TYPE_CREATURE),
	NewCard(114, -1, 7, 7, 7, "---G--", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(115, -1, 8, 5, 5, "---G-W", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(116, -1, 12, 8, 8, "BCDGLW", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(117, -1, 1, 1, 1, "B-----", 0, 0, 0, CARD_TYPE_ITEM_GREEN),
	NewCard(118, -1, 0, 0, 3, "------", 0, 0, 0, CARD_TYPE_ITEM_GREEN),
	NewCard(119, -1, 1, 1, 2, "------", 0, 0, 0, CARD_TYPE_ITEM_GREEN),
	NewCard(120, -1, 2, 1, 0, "----L-", 0, 0, 0, CARD_TYPE_ITEM_GREEN),
	NewCard(121, -1, 2, 0, 3, "------", 0, 0, 1, CARD_TYPE_ITEM_GREEN),
	NewCard(122, -1, 2, 1, 3, "---G--", 0, 0, 0, CARD_TYPE_ITEM_GREEN),
	NewCard(123, -1, 2, 4, 0, "------", 0, 0, 0, CARD_TYPE_ITEM_GREEN),
	NewCard(124, -1, 3, 2, 1, "--D---", 0, 0, 0, CARD_TYPE_ITEM_GREEN),
	NewCard(125, -1, 3, 1, 4, "------", 0, 0, 0, CARD_TYPE_ITEM_GREEN),
	NewCard(126, -1, 3, 2, 3, "------", 0, 0, 0, CARD_TYPE_ITEM_GREEN),
	NewCard(127, -1, 3, 0, 6, "------", 0, 0, 0, CARD_TYPE_ITEM_GREEN),
	NewCard(128, -1, 4, 4, 3, "------", 0, 0, 0, CARD_TYPE_ITEM_GREEN),
	NewCard(129, -1, 4, 2, 5, "------", 0, 0, 0, CARD_TYPE_ITEM_GREEN),
	NewCard(130, -1, 4, 0, 6, "------", 4, 0, 0, CARD_TYPE_ITEM_GREEN),
	NewCard(131, -1, 4, 4, 1, "------", 0, 0, 0, CARD_TYPE_ITEM_GREEN),
	NewCard(132, -1, 5, 3, 3, "B-----", 0, 0, 0, CARD_TYPE_ITEM_GREEN),
	NewCard(133, -1, 5, 4, 0, "-----W", 0, 0, 0, CARD_TYPE_ITEM_GREEN),
	NewCard(134, -1, 4, 2, 2, "------", 0, 0, 1, CARD_TYPE_ITEM_GREEN),
	NewCard(135, -1, 6, 5, 5, "------", 0, 0, 0, CARD_TYPE_ITEM_GREEN),
	NewCard(136, -1, 0, 1, 1, "------", 0, 0, 0, CARD_TYPE_ITEM_GREEN),
	NewCard(137, -1, 2, 0, 0, "-----W", 0, 0, 0, CARD_TYPE_ITEM_GREEN),
	NewCard(138, -1, 2, 0, 0, "---G--", 0, 0, 1, CARD_TYPE_ITEM_GREEN),
	NewCard(139, -1, 4, 0, 0, "----LW", 0, 0, 0, CARD_TYPE_ITEM_GREEN),
	NewCard(140, -1, 2, 0, 0, "-C----", 0, 0, 0, CARD_TYPE_ITEM_GREEN),
	NewCard(141, -1, 0, -1, -1, "------", 0, 0, 0, CARD_TYPE_ITEM_RED),
	NewCard(142, -1, 0, 0, 0, "BCDGLW", 0, 0, 0, CARD_TYPE_ITEM_RED),
	NewCard(143, -1, 0, 0, 0, "---G--", 0, 0, 0, CARD_TYPE_ITEM_RED),
	NewCard(144, -1, 1, 0, -2, "------", 0, 0, 0, CARD_TYPE_ITEM_RED),
	NewCard(145, -1, 3, -2, -2, "------", 0, 0, 0, CARD_TYPE_ITEM_RED),
	NewCard(146, -1, 4, -2, -2, "------", 0, -2, 0, CARD_TYPE_ITEM_RED),
	NewCard(147, -1, 2, 0, -1, "------", 0, 0, 1, CARD_TYPE_ITEM_RED),
	NewCard(148, -1, 2, 0, -2, "BCDGLW", 0, 0, 0, CARD_TYPE_ITEM_RED),
	NewCard(149, -1, 3, 0, 0, "BCDGLW", 0, 0, 1, CARD_TYPE_ITEM_RED),
	NewCard(150, -1, 2, 0, -3, "------", 0, 0, 0, CARD_TYPE_ITEM_RED),
	NewCard(151, -1, 5, 0, -99, "BCDGLW", 0, 0, 0, CARD_TYPE_ITEM_RED),
	NewCard(152, -1, 7, 0, -7, "------", 0, 0, 1, CARD_TYPE_ITEM_RED),
	NewCard(153, -1, 2, 0, 0, "------", 5, 0, 0, CARD_TYPE_ITEM_BLUE),
	NewCard(154, -1, 2, 0, 0, "------", 0, -2, 1, CARD_TYPE_ITEM_BLUE),
	NewCard(155, -1, 3, 0, -3, "------", 0, -1, 0, CARD_TYPE_ITEM_BLUE),
	NewCard(156, -1, 3, 0, 0, "------", 3, -3, 0, CARD_TYPE_ITEM_BLUE),
	NewCard(157, -1, 3, 0, -1, "------", 1, 0, 1, CARD_TYPE_ITEM_BLUE),
	NewCard(158, -1, 3, 0, -4, "------", 0, 0, 0, CARD_TYPE_ITEM_BLUE),
	NewCard(159, -1, 4, 0, -3, "------", 3, 0, 0, CARD_TYPE_ITEM_BLUE),
	NewCard(160, -1, 2, 0, 0, "------", 2, -2, 0, CARD_TYPE_ITEM_BLUE),
}
 /* STARTING */

/************************************/
/*           MOVE STRUCT            */
/************************************/
type Move struct {
	Cost			int
	Type			int
	Params			[]int
	probability		float64
}
func NewMove(type_, cost int, probability float64, params []int) *Move {
	return &Move {
		Cost: cost,
		Type: type_,
		probability: probability,
		Params: params,
	}
}
func (m *Move) Copy() *Move {
	move := &Move {
		Cost: m.Cost,
		Type: m.Type,
		probability: m.probability,
		Params: make([]int, len(m.Params)),
	}
	copy(move.Params, m.Params)
	return move
}
func (m *Move) Probability() float64 {
	return m.probability
}
func (m *Move) ToString() string {
	if m == nil {
		return ""
	}
	var str []string = make([]string, 0)
	switch m.Type {
	case MOVE_PICK:
		str = append(str, "PICK")
	case MOVE_PASS:
		str = append(str, "PASS")
	case MOVE_ATTACK:
		str = append(str, "ATTACK")
	case MOVE_USE:
		str = append(str, "USE")
	case MOVE_SUMMON:
		str = append(str, "SUMMON")
	}
	for _, p := range(m.Params) {
		str = append(str, fmt.Sprintf("%d", p))
	}
	return strings.Join(str, " ")
}
/************************************/
/*          STATE STRUCT            */
/************************************/
type State struct {
	Players		[]*Player
	AMoves		[]*Move
	Turn		int

	Draft		[]*Card
}
func NewState(hero, vilain *Player) *State {
	new_state := &State{
		Players: make([]*Player, 2),
		AMoves: nil,
		Turn: 0,
		Draft: nil,
	}
	new_state.Players[0] = hero
	new_state.Players[1] = vilain

	return new_state
}

// TODO
func LoadState() {}

func (s *State) Raw() []interface{} {
	r := make([]interface{}, 0)
	no_card := make([]interface{}, len(CARDS[0].Raw()))
	for i := 0 ; i < len(CARDS[0].Raw()) ; i++ {
		no_card[i] = -1
	}
	for p := 0 ; p < MAX_PLAYERS ; p++ {
		player := s.Players[p]
		r = append(r, player.Raw()...)
	}
	return r
}
func (s *State) Copy() (state *State) {
	state = NewState(nil, nil)
	for i, _ := range(s.Players) {
		state.Players[i] = s.Players[i].Copy()
	}
	return state
}

func (s *State) PrintHero() { s.PrintPlayer(s.Hero()) }
func (s *State) PrintVilain() { s.PrintPlayer(s.Vilain()) }
func (s *State) PrintHeroHand() { s.PrintHand(s.Hero()) }
func (s *State) PrintHeroBoard() { s.PrintBoard(s.Hero()) }
func (s *State) PrintVilainHand() { s.PrintHand(s.Vilain()) }
func (s *State) PrintVilainBoard() { s.PrintBoard(s.Vilain()) }
func (s *State) Print() {
	s.PrintHero()
	s.PrintHeroHand()
	s.PrintHeroBoard()
	s.PrintVilain()
	s.PrintVilainHand()
	s.PrintVilainBoard()
}
func (s *State) Hero() (*Player) {
	return s.Players[0]
}
func (s *State) Vilain() (*Player) {
	return s.Players[1]
}
func (s *State) SwapPlayers() {
	s.Players[0], s.Players[1] = s.Players[1], s.Players[0]
}

func (s *State) NextTurn() (err error) {
	s.SwapPlayers()
	s.Hero().IncreaseMana()
	s.Hero().ReloadMana()
	s.Hero().UpdateBoard()

	err = s.Hero().DrawStackCards()
	err = s.Hero().Draw()

	if err != nil {
		return err
	}
	s.AMoves = nil
	
	return nil
}

func (s *State) PrintHand(p *Player) {
	if p != nil {
		for _, c := range(p.Hand) {
			fmt.Fprintln(os.Stderr, "[GAME] Hand", c.Cost, c.ToString())
		}
	}
}
func (s *State) PrintBoard(p *Player) {
	if p != nil {
		for _, c := range(p.Board) {
			fmt.Fprintln(os.Stderr, "[GAME] Board", c.ToString())
		}
	}
}
func (s *State) PrintPlayer(p *Player) {
	if p != nil {
		fmt.Fprintln(os.Stderr, "[GAME] Player", p.Id, "L:", p.Life, "M:", p.Mana, "/", p.MaxMana, "D:", p.DeckCount(), "R:", p.Runes)
	}
}

func (s *State) DealDamageMonster(m1 *Card, dmg int) (int, error) {
	if m1 == nil || m1.Type != CARD_TYPE_CREATURE {
		return 0, fmt.Errorf("Card", m1, "is not a creature")
	}

	if dmg < 0 {
		return 0, fmt.Errorf("Cannot deal negative damage")
	}
	if dmg > 0 && m1.IsAbleTo(CARD_ABILITY_WARD) {
		m1.DisableAbility(CARD_ABILITY_WARD)
		return 0, nil
	}
	m1.Defense -= dmg

	return dmg, nil
}

func (s *State) CreatureFight(c1, c2 *Card) (error) {
	var dead1, dead2 bool

	dead1, dead2 = false, false

	dmg1, err1 := s.DealDamageMonster(c2, c1.Attack)

	if err1 != nil { return err1 }
	if dmg1 > 0 && c1.IsAbleTo(CARD_ABILITY_DRAIN) { s.Hero().GainLife(dmg1) }
	if dmg1 > 0 && c1.IsAbleTo(CARD_ABILITY_BREAKTHROUGH) { s.Vilain().ReceiveDamage(dmg1) }
	if dmg1 > 0 && c1.IsAbleTo(CARD_ABILITY_LETHAL) { dead1 = true }
	
	dead1 = dead1 || (c2.Defense <= 0)

	dmg2, err2 := s.DealDamageMonster(c1, c2.Attack)

	if err2 != nil { return err2 }
	if dmg2 > 0 && c2.IsAbleTo(CARD_ABILITY_LETHAL) { dead2 = true }
	dead2 = dead2 || c1.Defense <= 0

	if dead1 { s.Vilain().BoardRemoveCard(c2.Id) }
	if dead2 { s.Hero().BoardRemoveCard(c1.Id) }

	return nil
}
func (s *State) CreatureBoostAttack(c *Card, bonus int) (error) {
	if c == nil || c.Type != CARD_TYPE_CREATURE {
		return fmt.Errorf("Can't boost attack on non creature card")
	}
	if bonus == 0 {
		return nil
	}
	c.Attack += bonus

	if c.Attack < 0 {
		c.Attack = 0
	}
	return nil
}
func (s *State) CreatureBoostDefense(c *Card, bonus int) (error) {
	if c == nil || c.Type != CARD_TYPE_CREATURE {
		return fmt.Errorf("Can't boost defense on non creature card")
	}
	if bonus == 0 {
		return nil
	}
	c.Defense += bonus
	return nil
}
func (s *State) CreatureBoostAbilities(c *Card, abilities int) (error) {
	if c == nil || c.Type != CARD_TYPE_CREATURE {
		return fmt.Errorf("Can't boost abilities on non creature card")
	}
	c.EnableAbility(abilities)
	return nil
}

func (s *State) MoveUse(id1, id2 int) (error) {
	c1 := s.Hero().HandGetCard(id1)
	if c1 == nil {
		return fmt.Errorf("Use %d. Card doesn't exist in Hand", id1)
	}

	err := s.Hero().HandPlayCard(id1)
	if err != nil {
		return err
	}

	s.Hero().GainLife(c1.HealthChange)
	s.Vilain().ReceiveDamage(-c1.OpponentHealthChange)
	
	switch c1.Type {
	case CARD_TYPE_ITEM_BLUE:
		if id2 != -1 {
			c2 := s.Vilain().BoardGetCard(id2)
			if c2 == nil {
				return fmt.Errorf("Use %d. Card doesn't exist in Board on Player %d", id2, s.Hero().Id)	
			}
			_, err2 := s.DealDamageMonster(c2, -c1.Defense)
			if err2 != nil { return err2 }
			dead2 := c2.Defense <= 0
			if dead2 { s.Vilain().BoardRemoveCard(id2) }
		}

	case CARD_TYPE_ITEM_GREEN:
		c2 := s.Hero().BoardGetCard(id2)
		if c2 == nil {
			return fmt.Errorf("Use %d. Card doesn't exist in Board on Player %d", id2, s.Hero().Id)	
		}
		s.CreatureBoostAttack(c2, c1.Attack)
		s.CreatureBoostDefense(c2, c1.Defense)
		c2.EnableAbility(c1.Abilities)
		
	case CARD_TYPE_ITEM_RED:
		c2 := s.Vilain().BoardGetCard(id2)
		if c2 == nil {
			return fmt.Errorf("Use %d. Card doesn't exist in Board on Player %d", id2, s.Vilain().Id)	
		}
		s.CreatureBoostAttack(c2, c1.Attack)
		s.CreatureBoostDefense(c2, c1.Defense)
		c2.DisableAbility(c1.Abilities)

		if c2.Defense <= 0 { s.Vilain().BoardRemoveCard(id2) }
	}

	return nil
}
func (s *State) MoveSummon(id1 int) error {
	if len(s.Hero().Board) >= MAX_BOARD_CARD {
		return fmt.Errorf("Max board card reach %d", MAX_BOARD_CARD)
	}

	c1  := s.Hero().HandGetCard(id1)
	err := s.Hero().HandPlayCard(id1)

	if err != nil {
		return err
	}
	if c1.Type != CARD_TYPE_CREATURE {
		return fmt.Errorf("Can't summon card type %s", typeToString(c1.Type))
	}

	s.Hero().GainLife(c1.HealthChange)
	s.Vilain().ReceiveDamage(-c1.OpponentHealthChange)

	if c1.IsAbleTo(CARD_ABILITY_CHARGE) { c1.Charge = 1 }

	return nil
}

func (s *State) MoveAttackPolicy(id1, id2 int) (error) {
	c1 := s.Hero().BoardGetCard(id1)
	if c1 == nil {
	    return fmt.Errorf("[GAME][ATTACK] Unknow card %d", id1)
	}
	if c1.Attacked || c1.Charge == 0 {
	    return fmt.Errorf("[GAME][ATTACK] Move ATTACK %d %d not permitted", id1, id2)
	}
	guards := s.Vilain().BoardGetGuardsId()

	c2 := s.Vilain().BoardGetCard(id2)
	if len(guards) > 0 {
		if ! (c2 != nil && c2.IsAbleTo(CARD_ABILITY_GUARD)) { 
	   		return fmt.Errorf("[GAME][ATTACK] Move ATTACK %d %d not permitted", id1, id2)
		}
	}
	return nil
}
func (s *State) MoveAttack(id1, id2 int) (err error) {
	var err_str string
	c1 := s.Hero().BoardGetCard(id1)
	if c1 == nil {
		err_str = fmt.Sprintf("[GAME][ATTACK] Current player %d don't have card %d", s.Hero().Id, id1)
		return errors.New(err_str)
	}

	err = s.MoveAttackPolicy(id1, id2)
	if err != nil {
		return err
	}

	c1a := c1.Attack

	if id2 == -1 {
		s.Vilain().ReceiveDamage(c1a)
	} else {
		c2 := s.Vilain().BoardGetCard(id2)
		if c2 == nil {
			err_str = fmt.Sprintf("MoveAttack: Current oppoent %d don't have card %d", s.Vilain().Id, id2)
			return errors.New(err_str)
		}
		err = s.CreatureFight(c1, c2)
	}

	c1.Attacked = true
	return err
}
func (s *State) MovePick(id int) (err error) {
	if len(s.Draft) == 0 {
		return fmt.Errorf("[GAME][DRAFT] wrong action pick")
	}
	if id < len(s.Draft) {
		card := s.Draft[id].Copy()
		card.Id = s.Hero().DeckCount() + s.Vilain().DeckCount() + 1
		s.Hero().Pick(s.Draft[id])
	}
	return nil
}
func (s *State) UpdateAvailableMoves(optimized bool) {
	
	s.AMoves = make([]*Move, 0)

	mu := s.AvailableMovesUse(optimized)
	s.AMoves = append(s.AMoves, mu...)

	ms := s.AvailableMovesSummon(optimized)
	s.AMoves = append(s.AMoves, ms...)
	
	mb := s.AvailableMovesBoard(optimized)	
	s.AMoves = append(s.AMoves, mb...)
	    
	/*
    ms := s.AvailableMovesSummon(optimized)
    s.AMoves = append(s.AMoves, ms...)
	if len(s.AMoves) == 0 {
	    mb := s.AvailableMovesBoard(optimized)	
	    s.AMoves = append(s.AMoves, mb...)
	    if len(s.AMoves) == 0 {

	    }
	}
	*/
}
func (s *State) CopyAvailableMoves() []*Move {
	var moves []*Move = nil

	if s.AvailableMoves == nil {
		return nil
	}

	moves = make([]*Move, len(s.AMoves))
	for i, m := range(s.AMoves) {
		moves[i] = m.Copy()
	}

	return moves
}
func (s *State) AvailableMovesBoardOne(c *Card) []*Move {
	var move *Move
	moves := make([]*Move, 0)

	if c == nil {
		return moves
	}
	if c.Attack <= 0 || c.Attacked || c.Charge <= 0 {
		return moves
	}

	guards := s.Vilain().BoardGetGuardsId()

	if len(guards) > 0 {
		for _, v := range(guards) {
			move = NewMove(MOVE_ATTACK, 0, 1, []int{c.Id, v.Id})
			moves = append(moves, move)
		}
	} else {
		move = NewMove(MOVE_ATTACK, 0, 1, []int{c.Id, -1})
		moves = append(moves, move)

		for _, v := range(s.Vilain().Board) {
			move = NewMove(MOVE_ATTACK, 0, 1, []int{c.Id, v.Id})
			moves = append(moves, move)
		} 
	}
	return moves
}
func (s *State) AvailableMovesBoard(optimized bool) []*Move {
	moves := make([]*Move, 0)

	for _, c := range(s.Hero().Board) {
		moves = append(moves, s.AvailableMovesBoardOne(c)...)
	}
	if len(moves) == 0 {
		return moves
	}
	//
	// Optimized Action. Only return the move
	// of the best board evaluation.
	//
	if optimized {
		combinations := iterate_combinations(moves)
		best_move_id := -1
		best_move_score := -1.0
		for i, ms := range(combinations) {
			tmp_state := s.Copy()
			for _, m := range(ms) {	tmp_state.Move(m) }
			score := tmp_state.Evaluate()
			if score > best_move_score {
				best_move_score = score
				best_move_id = i
			}
		}
		if best_move_id != 0 {
			return combinations[best_move_id]
		}
	}
	
	return moves
}
func (s *State) AvailableMovesUse(optimized bool) []*Move {
	var move *Move

	moves := make([]*Move, 0)
	
	for _, h := range(s.Hero().Hand) {
		if h.Cost > s.Hero().Mana { continue }
		switch h.Type {
		case CARD_TYPE_ITEM_BLUE:
			move = NewMove(MOVE_USE, h.Cost, 1, []int{h.Id, -1})
			moves = append(moves, move)
			if h.Defense < 0 {
				for _, c := range(s.Vilain().Board) {
					move = NewMove(MOVE_USE, h.Cost, 1, []int{h.Id, c.Id})
					moves = append(moves, move)
				}	
			}
		case CARD_TYPE_ITEM_GREEN:
			for _, c := range(s.Hero().Board) {
				move = NewMove(MOVE_USE, h.Cost, 1, []int{h.Id, c.Id})
				moves = append(moves, move)
			}
		case CARD_TYPE_ITEM_RED:
			for _, c := range(s.Vilain().Board) {
				move = NewMove(MOVE_USE, h.Cost, 1, []int{h.Id, c.Id})
				moves = append(moves, move)
			}
		}
	}

	if len(moves) == 0 {
		return moves
	}
	// Precalculate State
	if optimized {
		filter_moves := make([][]*Move, 0)
		combinations := iterate_combinations(moves)
		for _, moves_cb := range(combinations) {
			var cost int = 0
			for _, move_cb := range(moves_cb) {
				cost += move_cb.Cost
			}
			if cost <= s.Hero().Mana {
				filter_moves = append(filter_moves, moves_cb)
			}
		}

		best_move_id := -1
		best_move_score := -1.0
		for i, ms := range(filter_moves) {
			tmp_state := s.Copy()
			for _, m := range(ms) {	tmp_state.Move(m) }
			score := tmp_state.Evaluate()
			if score > best_move_score {
				best_move_score = score
				best_move_id = i
			}
		}


		if best_move_id != -1 {
			return filter_moves[best_move_id]
		}
	}
	
	return moves
}
func (s *State) AvailableMovesSummon(optimized bool) []*Move {
	var move *Move

	moves := make([]*Move, 0)

	for _, h := range(s.Hero().Hand) {
		if h.Cost > s.Hero().Mana { continue }
		switch h.Type {
		case CARD_TYPE_CREATURE:
			move = NewMove(MOVE_SUMMON, h.Cost, 1, []int{h.Id})
			moves = append(moves, move)
		}
	}
	if len(moves) == 0 {
		return moves
	}
	if optimized {
		filter_moves := make([][]*Move, 0)
		combinations := iterate_combinations(moves)

		// Keep only combinations that < Player current Mana
		for _, moves_cb := range(combinations) {
			var cost int = 0
			for _, move_cb := range(moves_cb) {
				cost += move_cb.Cost
			}
			if cost <= s.Hero().Mana {
				filter_moves = append(filter_moves, moves_cb)
			}
		}

		best_move_id := -1
		best_move_score := -1.0
		for i, tmp_moves := range(filter_moves) {
			total_power := 0.0
			for _, m := range(tmp_moves) { 
				c := s.Hero().HandGetCard(m.Params[0])
				total_power += float64(c.Attack) + float64(c.Defense)
				if c.IsAbleTo(CARD_ABILITY_BREAKTHROUGH) { total_power += float64(WEIGHT_BREAKTHROUGH) }
				if c.IsAbleTo(CARD_ABILITY_GUARD) { total_power += float64(WEIGHT_GUARD) } 
				if c.IsAbleTo(CARD_ABILITY_LETHAL) { total_power += float64(WEIGHT_LETHAL) } 
				if c.IsAbleTo(CARD_ABILITY_DRAIN) { total_power += float64(WEIGHT_DRAIN) } 
				if c.IsAbleTo(CARD_ABILITY_WARD) { total_power += float64(WEIGHT_WARD) } 
				if c.IsAbleTo(CARD_ABILITY_CHARGE) { total_power += float64(WEIGHT_CHARGE) } 
			
			}
			if best_move_score < total_power {
				best_move_score = total_power
				best_move_id = i
			}
		}
	
		if best_move_id != - 1{
			return filter_moves[best_move_id]
		}
	}
	return moves
}

func (s *State) AvailableMoves() []*Move {
	if s.AMoves != nil {
		return s.AMoves
	}

	s.UpdateAvailableMoves(false)

	if len(s.AMoves) == 0 {
		s.AMoves = append(s.AMoves, NewMove(MOVE_PASS, 0, 0, nil))
	}

	return s.AMoves
}
func (s *State) Evaluate() float64 {

	
	var score float64 = 0
	var hmana, htma, hpow, hdef, vpow, vdef float64
	var vmana, vtma, hl, vl, hm, vm float64

	//hc = float64(len(s.Hero().Hand))
	hl = float64(s.Hero().Life)
	hm = float64(len(s.Hero().Board))
	for _, c := range(s.Hero().Board) {
		hpow += float64(c.Attack)
		if c.Defense>=4 {
			htma += float64(c.Defense)
		}
		hmana += float64(c.Cost)
		hdef += float64(c.Defense)
	}

	//vc = float64(len(s.Vilain().Hand))
	vl = float64(s.Vilain().Life)
	vm = float64(len(s.Vilain().Board))
	for _, c := range(s.Vilain().Board) {
		vpow += float64(c.Attack)
		if c.Defense>=4{
			vtma += float64(c.Attack)
		}
		vmana += float64(c.Cost)
		vdef += float64(c.Defense)
	}

	minion_advantage := hm - vm
    tough_minion_advantage := htma-vtma

	if vl <= 0 {
		return OUTCOME_WIN
	} else if hl <= 0 {
		return OUTCOME_LOSE
	}

	score += minion_advantage * minion_advantage
	score += (hpow + hdef  - (vpow + vdef))
	score += tough_minion_advantage
	score += hl - vl

	return score 
}
func (s *State) IsEndTurn() bool {
	s.UpdateAvailableMoves(false)
	for _, m := range(s.AvailableMoves()) {
		if m.Cost <= s.Hero().Mana && m.Type != MOVE_PASS {
			return false
		}
	}
	return true
}
func (s *State) RandomMove() (*Move, error) {
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)

	l := len(s.AvailableMoves())
	if l > 0 {
		n := random.Intn(l)
		move := s.AvailableMoves()[n]
		//fmt.Println("[MCTS] Random move", move.ToString(), "in", s.AvailableMoves())
		e := s.Move(move)
		return move, e
	}
	return nil, nil
}
func (s *State) Move(m *Move) (error) {
	var err error

	if m == nil {
		return fmt.Errorf("Move is nil")
	}

	switch m.Type {
	case MOVE_PASS:
		if len(s.Draft) > 0 {
			s.MovePick(0)
		}
	case MOVE_PICK:
		err = s.MovePick(m.Params[0])
	case MOVE_SUMMON:
		err = s.MoveSummon(m.Params[0])
	case MOVE_ATTACK:
		err = s.MoveAttack(m.Params[0], m.Params[1])
	case MOVE_USE:
		err = s.MoveUse(m.Params[0], m.Params[1])
	}

	return err 
}
func (s *State) GameOver() *Player {
	if s.Hero().Life <= 0 {
		return s.Vilain()
	} else if s.Vilain().Life <= 0 {
		return s.Hero()
	}
	return nil
}

/* STARTING */

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
	c.Abilities &= ^ability
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

func (c *Card) ToString() string {
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
/* STARTING */

type AgentRandom struct {

}

func NewAgentRandom() *AgentRandom {
	return &AgentRandom{}
}

func (a *AgentRandom) Name() string { return "RANDOM" }
func (a *AgentRandom) Think(s *State) []*Move {

	moves := s.AvailableMoves()

	if len(moves) == 0 {
		moves = append(moves, NewMove(MOVE_PASS, 0, 0, nil))
	}
	return moves
}

/* STARTING */
type AgentDraft struct {

}

func NewAgentDraft() *AgentDraft {
	return &AgentDraft{}
}

func Gauss(x, m, v float64) float64 {
	exp := -math.Pow((x - m), 2) / (2 * v * v)
	exp = math.Exp(exp)
	return exp / (v * math.Sqrt(2 * math.Pi))
}
func (a *AgentDraft) Name() string { return "DRAFT" }
func (a *AgentDraft) Think(s *State) []*Move {
	var maps_cost map[int]int = make(map[int]int)
	var maps_type map[int]int = make(map[int]int)

	ref_cost := map[int]int {
		0: 0,
		1: 0,
		2: 10,
		3: 10,
		4: 0,
		5: 1,
		6: 8,
		7: 1,
	}

	ref_type := map[int]int {
		CARD_TYPE_CREATURE: 23,
		CARD_TYPE_ITEM_BLUE: 0,
		CARD_TYPE_ITEM_GREEN: 7,
		CARD_TYPE_ITEM_RED: 0,
	}

	for _, c := range(s.Hero().Deck) {
		maps_cost[c.Cost]++
		maps_type[c.Type]++
	}

	id := -1
	best_score := -1.0
	score := make([]float64, len(s.Draft))
	for i, c := range(s.Draft) {
		sum_abilities := 0.0
		/*
		for _, a := range(c.Abilities) { if a != "-" { sum_abilities += 2 }}
		sum_abilities += float64(c.CardDraw)
		sum_abilities += float64(c.HealthChange) / 5.0
		sum_abilities -= float64(c.OpponentHealthChange) / 5.0600
		*/

		cost := 7
		if c.Cost < 7  {
			cost = c.Cost
		}

		switch c.Type {
		case CARD_TYPE_CREATURE:

			if c.IsAbleTo(CARD_ABILITY_GUARD) { sum_abilities += 4 }
			if c.IsAbleTo(CARD_ABILITY_WARD) { sum_abilities += 3 }
			if c.IsAbleTo(CARD_ABILITY_LETHAL) { sum_abilities += 3 }
			if c.IsAbleTo(CARD_ABILITY_CHARGE) { sum_abilities += 0 }
			if c.IsAbleTo(CARD_ABILITY_DRAIN) { sum_abilities += 2 }
			if c.IsAbleTo(CARD_ABILITY_BREAKTHROUGH) { sum_abilities += 1 }

			score[i] = (float64(c.Attack) * 0.9 + float64(c.Defense) * 1.1) / (float64(c.Cost) + 1)
			score[i] *= (1 + sum_abilities / 9 )
			score[i] += (float64(ref_cost[cost]) - float64(maps_cost[cost])) / 1.5
			score[i] += (float64(ref_type[c.Type]) - float64(maps_type[c.Type])) / 4
			score[i] += Gauss(float64(c.Attack - c.Defense), 0.0, 2.2) * 20
			if c.Cost >= 10 || c.Cost <= 3 {
				score[i] /= (1 + float64(c.Cost)/100.0)
			}
		case CARD_TYPE_ITEM_BLUE:
			score[i] = 0
			// score[i] *= CARDS_VALUE[c.CardNumber - 1]
		case CARD_TYPE_ITEM_GREEN:
			score[i] = (float64(c.Attack)  + float64(c.Defense)) / (float64(c.Cost) + 1)
			if c.IsAbleTo(CARD_ABILITY_GUARD) { sum_abilities += 2 }
			if c.IsAbleTo(CARD_ABILITY_WARD) { sum_abilities += 2 }
			if c.IsAbleTo(CARD_ABILITY_LETHAL) { sum_abilities += 1 }
			if c.IsAbleTo(CARD_ABILITY_CHARGE) { sum_abilities += 0 }
			if c.IsAbleTo(CARD_ABILITY_DRAIN) { sum_abilities += 0.5 }
			if c.IsAbleTo(CARD_ABILITY_BREAKTHROUGH) { sum_abilities += 0.3 }

			score[i] += sum_abilities / 10
			score[i] += (float64(ref_cost[cost]) - float64(maps_cost[cost])) / 4
			score[i] += (float64(ref_type[c.Type]) - float64(maps_type[c.Type]))
			//score[i] += Gauss(float64(c.Attack - c.Defense), 0.0, 1.8) * 20

		default:
			score[i] += (float64(ref_cost[cost]) - float64(maps_cost[cost]))
			score[i] += (float64(ref_type[c.Type]) - float64(maps_type[c.Type])) * 2

		}
		fmt.Fprintln(os.Stderr, "Score", i, "->", score[i])
		if score[i] > best_score {
			best_score = score[i]
			id = i
		}
	}
	if id == -1 {
		return []*Move{NewMove(MOVE_PICK, 0, 1, []int{0})}
	}
	return []*Move{NewMove(MOVE_PICK, 0, 1, []int{id})}
}
/* STARTING */
type Agent interface {
	Name()			string
	Think(s *State) []*Move
}

type AI struct {
	Agents	[]Agent
}

func NewAI() *AI {
	return &AI{
		Agents: make([]Agent, 0),
	}
}
func (ai *AI) LoadAgentRandom() {
	ai.Agents = append(ai.Agents, NewAgentRandom())
}
func (ai *AI) LoadAgentMCTS() {
	ai.Agents = append(ai.Agents, NewAgentMCTS())
}

func (ai *AI) LoadAgentDraft() {
	ai.Agents = append(ai.Agents, NewAgentDraft())
}

func (ai *AI) GetAgent(name_agent string) Agent {
	for _, a := range(ai.Agents) {
		if a.Name() == name_agent { return a }
	}
	return nil
}
func (ai *AI) Think(name_agent string, s *State) []*Move {
	a := ai.GetAgent(name_agent)
	if a != nil {
		return a.Think(s)
	}
	return []*Move{}
}
/* STARTING */
type AgentMCTS struct {

}

func NewAgentMCTS() *AgentMCTS {
	return &AgentMCTS{}
}

func (a *AgentMCTS) Name() string { return "MCTS" }
func (a *AgentMCTS) Think(s *State) []*Move {
	node, _ := MonteCarloTimeout(s, BIAS_PARAMETER, MCTS_ITERATION, MCTS_SIMULATION, MCTS_TIMEOUT)
	moves := make([]*Move, 0)



	for ; node != nil && len(node.Children) > 0 ; {
		var score  float64 = -200
		if node.EndTurn { break }

		for _, child := range node.Children {
			child_score := MCCalculateScore(child, BIAS_PARAMETER)
			if child_score > score {
				score = child_score
				node = child
			}
		}
		moves = append(moves, node.ByMove)
	}

	if len(moves) == 0 {
		moves = append(moves, NewMove(MOVE_PASS, 0, 0, nil))
	}
	return moves
}



type Node struct {
	Id		 	int
	Parent   		*Node
	Children 		[]*Node
	State    		*State
	Outcome     	float64
	Visits   		int
	ByMove			*Move
	EndTurn  		bool
	UnexploreMoves	[]*Move
}

func NewNode(parent *Node, state *State, move *Move) *Node {
	return &Node{
		Id:		  1,
		Parent:   parent,
		Children: nil,
		State:    state,
		Outcome:     0,
		Visits:   0,
		ByMove: move,
		EndTurn: false,
		UnexploreMoves: nil,
	}
}
func (n *Node) Count() int {
	count := 1 
	for _, c := range(n.Children) {
		count += c.Count()
	}
	return count
}
func (n *Node) DotPrintNode(id int, f *os.File) (error) {
	var col_name []string = []string{
		"PLAYER",
		"HEALTH",
		"MANA",
		"CARDS",
		"RUNES",
	}


	graph := fmt.Sprintf("%d [shape=none, label=<", id)
	f.WriteString(graph)
	f.WriteString("<TABLE BORDER=\"0\" CELLBORDER=\"1\" CELLSPACING=\"0\" CELLPADDING=\"4\">")
	f.WriteString("<TR>")
	f.WriteString("<TD>SCORE</TD>")
	f.WriteString("<TD>OUTCOME</TD>")
	f.WriteString("<TD>VISITS</TD>")
	f.WriteString("<TD>ENDTURN</TD>")
	f.WriteString("<TD>UM</TD>")
	f.WriteString("</TR>")
	f.WriteString("<TR>")
	f.WriteString(fmt.Sprintf("<TD>%f</TD>", MCCalculateScore(n, BIAS_PARAMETER)))
	f.WriteString(fmt.Sprintf("<TD>%f</TD>", n.Outcome))
	f.WriteString(fmt.Sprintf("<TD>%d</TD>", n.Visits))
	f.WriteString(fmt.Sprintf("<TD>%t</TD>", n.EndTurn))
	f.WriteString(fmt.Sprintf("<TD>%d</TD>", len(n.UnexploreMoves)))
	f.WriteString("</TR>")
	f.WriteString("<TR>")

	for _, t := range(col_name) {
		f.WriteString("<TD>")
		f.WriteString(t)
		f.WriteString("</TD>")
	}
	f.WriteString("</TR>")

	for _, p := range(n.State.Players) {
		f.WriteString("<TR>")
		f.WriteString("<TD>")

		f.WriteString(fmt.Sprintf("%d", p.Id))
		f.WriteString("</TD>")
		for i, v := range(p.Raw()) {
			f.WriteString("<TD>")
			if i == 1 {
				f.WriteString(fmt.Sprintf("%d/%d", p.Mana, p.MaxMana))
			} else {
				f.WriteString(fmt.Sprintf("%d", v))
			}
			f.WriteString("</TD>")
		}
		f.WriteString("</TR>")
		for _, c := range(p.Hand) {
			f.WriteString("<TR>")
			f.WriteString(fmt.Sprintf("<TD COLSPAN=\"%d\">", len(col_name)))
			f.WriteString(fmt.Sprintf("Hand %s", c.ToString()))
			f.WriteString("</TD>")
			f.WriteString("</TR>")
		}
		for _, c := range(p.Board) {
			f.WriteString("<TR>")
			f.WriteString(fmt.Sprintf("<TD COLSPAN=\"%d\">", len(col_name)))
			f.WriteString(fmt.Sprintf("Board %s", c.ToString()))
			f.WriteString("</TD>")
			f.WriteString("</TR>")
		}
	}
	f.WriteString("</TABLE>>]\n")


	for i, c := range(n.Children) {
		str := fmt.Sprintf("%d -- %d  [label=\"%s\"]", id, id * 100 + 1 + i, c.ByMove.ToString())
		c.DotPrintNode(id * 100 + 1 + i, f)
		_, _ = f.WriteString(str)
	}
	return nil
}

func (n *Node) ExportGraph(filename string) (error) {
	var err error


	f, err1 := os.Create(fmt.Sprintf("./%s.dot", filename))
	if err1 != nil {
		fmt.Println(err1)
		return err1
	}

	_, err = f.WriteString("graph generate {\n")
	n.DotPrintNode(1, f)
	_, err = f.WriteString("}\n")

	defer f.Close()
	return err
}
func (n *Node) DeleteUnexploreMoves(m *Move) (error) {
	idx := -1
	if n.UnexploreMoves == nil || len(n.UnexploreMoves) == 0 {
		return fmt.Errorf("Can't remove Move in empty list")
	}
	for i, tmp_m := range(n.UnexploreMoves) {
		if  tmp_m == m {
			idx = i
			break
		}
	}
	len_moves := len(n.UnexploreMoves)
	if idx != -1 && len_moves > 1 {
		n.UnexploreMoves = append(n.UnexploreMoves[:idx], n.UnexploreMoves[idx+1:]...)
	} else if idx != -1 {
		n.UnexploreMoves = make([]*Move, 0)
	}
	return nil
}

func (n *Node) UpdateScore(score float64) {
	n.Visits++
	n.Outcome += score
}

func MonteCarloTimeout(state *State, bias float64, iteration, simulation, timeout int) (*Node, error) {
	var node *Node
	var done chan bool = make(chan bool)

	now := time.Now()
	to	:= time.Nanosecond * time.Duration(timeout) * (1000000)
	time.AfterFunc(to, func () {
		done <- true
	})

	var running = true
	var root_node *Node

	root_node = NewNode(nil, state, nil)

	for ; running ; {

		node = MCSelection(root_node, bias)
		node = MCExpansion(node)
		score := MCSimulation(node.State, simulation)
		MCBackPropagation(node, score)

		select {
		case <-done:
			fmt.Fprintln(os.Stderr, "Timeout", time.Since(now))
			running = false
			break
		default:
		}
	}
	fmt.Fprintln(os.Stderr, "Node count:", root_node.Count())

	return root_node, nil
}
func MonteCarlo(state *State, bias float64, iteration, simulation int) (*Node, error) {
	var node *Node

	if iteration <= 0 { return nil, fmt.Errorf("Iteration should be > 0") }
	if simulation < -1 { return nil, fmt.Errorf("Simulation should be > 0") }

	var root_node *Node
	root_node = NewNode(nil, state, nil)

	for i := 0; i < iteration; i++ {

		node = MCSelection(root_node, bias)
		node = MCExpansion(node)
		/*
		if len(root_node.UnexploreMoves) == 0 && len(root_node.Children) == 1 {


			return node, nil
		}
		*/
		score := MCSimulation(node.State, simulation)
		MCBackPropagation(node, score)
	}

	//root_node.ExportGraph("./games/graph-test")
	fmt.Fprintln(os.Stderr, "Node count:", root_node.Count())
	return root_node, nil
}

func MCSelection(node *Node, bias float64) *Node {
	var candidate_node *Node

	if node.UnexploreMoves == nil {
		node.State.AvailableMoves()
		node.UnexploreMoves = node.State.CopyAvailableMoves()
	}


	if len(node.UnexploreMoves) == 0 && node.Children != nil && len(node.Children) > 0 {
		candidate_node = nil
		score := -1.0
		for _, n := range node.Children {
			child_score := MCCalculateScore(n, bias)
			if child_score > score || candidate_node == nil {
				score = child_score
				candidate_node = n
			}
		}
		if candidate_node == nil {
			return node
		}
		return MCSelection(candidate_node, bias)
	}
	return node
}
func MCCalculateScore(node *Node, bias float64) float64 {
	if node.Parent == nil {
		return 0
	}
	exploitScore := float64(node.Outcome) / float64(node.Visits)
	exploreScore := math.Sqrt(2 * math.Log(float64(node.Parent.Visits)) / float64(node.Visits))
	exploreScore = bias * exploreScore

	return exploitScore + exploreScore
}
func MCExpansion(node *Node) *Node {

	if len(node.UnexploreMoves) == 0 {
		return node
	}


	new_state := node.State.Copy()

	// Pick random move
	source 	:= rand.NewSource(time.Now().UnixNano())
	random 	:= rand.New(source)
	rmove 	:= node.UnexploreMoves[random.Intn(len(node.UnexploreMoves))]

	new_state.Move(rmove)
	new_node := NewNode(node, new_state, rmove)
	node.DeleteUnexploreMoves(rmove)

	new_node.Parent = node
	node.Children = append(node.Children, new_node)
	if new_state.IsEndTurn() || rmove.Type == MOVE_PASS {
		new_node.EndTurn = true
		new_state.NextTurn()
	}

	return new_node
}
func MCSimulation(state *State, simulation int) float64 {

	var moves []*Move = nil


	source 	:= rand.NewSource(time.Now().UnixNano())
	random 	:= rand.New(source)

	simulate_state := state.Copy()

	for i := 0 ; simulate_state.GameOver() != nil && (simulation == -1 || i < simulation) ; i++  {
		moves = simulate_state.AvailableMoves()
		if moves != nil || len(moves) == 0 {
			break
		}
		move := moves[random.Intn(len(moves))]
		simulate_state.Move(move)
		if simulate_state.IsEndTurn() {
			simulate_state.NextTurn()
		}
	}
	return simulate_state.Evaluate()
}

func MCBackPropagation(node *Node, score float64) *Node {
	id_hero := node.State.Hero().Id
	for node.Parent != nil {
		if node.State.Hero().Id == id_hero {
			node.UpdateScore(score)
		} else {
			node.UpdateScore(-score)
		}
		node = node.Parent
	}
	node.Visits++
	return node
}
func main() {
	mainCG()
}
