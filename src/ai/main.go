package main

import (
	"fmt"
	"math/rand"
	"time"

	ag "agents"
	g  "game"
)

var Hero *g.Player 		= g.NewPlayer(1, g.STARTING_LIFE, g.STARTING_MANA, g.STARTING_RUNES)
var Vilain *g.Player 	= g.NewPlayer(2, g.STARTING_LIFE, g.STARTING_MANA, g.STARTING_RUNES)

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

	s := initState()

	for i := 0 ; i < 10 ; i++ {
		fmt.Println("============== BEGIN (",i,") ===============")
		s.Print()
		moves := ai.Think("MCTS", s)	
		fmt.Println("====================================")
		fmt.Println("[MCTS] Moves:", len(moves))
		for _, m := range(moves) {
			s.Move(m)
		}
		s.NextTurn()
	}

}
