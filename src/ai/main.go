package main

import (
	"fmt"
	"math/rand"
	"time"
	"os"
	"strings"

	ag "agents"
	g  "game"
)

var Hero *g.Player	= g.NewPlayer(1, g.STARTING_LIFE, g.STARTING_MANA, g.STARTING_RUNES)
var Vilain *g.Player	= g.NewPlayer(2, g.STARTING_LIFE, g.STARTING_MANA, g.STARTING_RUNES)

func pickRandomCard() *g.Card {
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)
	n := random.Intn(len(g.CARDS))
	c := g.CARDS[n]
	return c
}
func emptyDeck(p1 *g.Player) {
	p1.Deck = make([]*g.Card, 0)
}
func emptyRunes(p1 *g.Player) {
	p1.Runes = 0
}
func loadDeck(p1 *g.Player, start_id int) {
	if p1 != nil {
		for i := 0 ; i < g.STARTING_CARD ; i++ {
			c := pickRandomCard()
			c.Id = start_id + i
			p1.DeckAddCard(c)
		}
	}
}
func loadHand(p1 *g.Player, n int) {
	if p1 != nil { p1.DrawN(n) }
}
func initState() *g.State {
	p1 := Hero.Copy()
	p2 := Vilain.Copy()

	s  := g.NewState(p2, p1)

	loadDeck(p1, 1)
	loadDeck(p2, g.STARTING_CARD + 1)

	loadHand(p1, 4)
	loadHand(p2, 5)

	s.NextTurn()
	return s
}


func main() {

	ai := ag.NewAI()
	ai.LoadAgentMCTS()
	ai.LoadAgentRandom()

	s := initState()

	for i := 0 ; s.GameOver() == nil ; i++ {
		var moves []*g.Move;

		fmt.Println("============== BEGIN (",i,") ===============")
		s.Print()
		if s.Hero().Id == 1 {
			moves = ai.Think("MCTS", s)
		} else {
			moves = ai.Think("RANDOM", s)
		}
		fmt.Println("====================================")
		fmt.Println("[MOVES]", len(moves))
		for _, m := range(moves) {
			s.Move(m)
		}
		s.NextTurn()
	}
	winner := s.GameOver()
	if winner != nil {
		fmt.Println("The winner is player", winner.Id)
	}

}

/* STARTING */

func mainCG() {

	var hero, vilain *g.Player
	var phero, pvilain *g.Player
	var step_draft bool
	var previous_board []int = make([]int, 0)
	var round int = 0

	ai := ag.NewAI()
	ai.LoadAgentMCTS()
	ai.LoadAgentRandom()
	ai.LoadAgentDraft()

	for {

		var players []*g.Player

		players = make([]*g.Player, 2)
		step_draft = false

		for i := 0; i < 2; i++ {
			var playerHealth, playerMana, playerDeck, playerRune int
			fmt.Scan(&playerHealth, &playerMana, &playerDeck, &playerRune)

			players[i] = g.NewPlayer(i, playerHealth, playerMana, playerRune)
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

		var draft []*g.Card = make([]*g.Card, 0)

		for i := 0; i < cardCount; i++ {
			var cardNumber, instanceId, location, cardType, cost, attack, defense int
			var abilities string
			var myHealthChange, opponentHealthChange, cardDraw int

			var card *g.Card

			fmt.Scan(&cardNumber, &instanceId, &location, &cardType, &cost, &attack, &defense, &abilities, &myHealthChange, &opponentHealthChange, &cardDraw)

			card = g.NewCard(cardNumber, instanceId, cost, attack, defense, abilities, myHealthChange, opponentHealthChange, cardDraw, cardType)
			switch cardType {
			case g.CARD_TYPE_CREATURE:
				if card.IsAbleTo(g.CARD_ABILITY_CHARGE) {
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


		init_state := g.NewState(hero, vilain)

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

			var moves []*g.Move;
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


