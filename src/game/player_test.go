package game

import (
	"testing"
	"reflect"
)

var Player1 *Player = NewPlayer(1, 30, 0, 25)

func TestPlayerCopy(t *testing.T) {
	p := Player1.Copy()
	if ! reflect.DeepEqual(p, Player1) {
		t.Errorf("Hero copy is not the same got: %v, want: %v.", p, Player1)
	}
	
}
func TestPlayerHP(t *testing.T) {
	if Player1.Life != 30 {
		t.Errorf("Hero life got: %d, want: %d.", Player1.Life, 30)
	}
}
func TestPlayerHPDecrease(t *testing.T) {
	p := Player1.Copy()
	p.ReceiveDamage(10)
	if p.Life != 20 {
		t.Errorf("Hero life got: %d, want: %d.", p.Life, 20)
	}
}
func TestPlayerHPNotDecrease(t *testing.T) {
	p := Player1.Copy()
	p.ReceiveDamage(-10)
	if p.Life != 30 {
		t.Errorf("Hero life got: %d, want: %d.", p.Life, 30)
	}
}
func TestPlayerMana(t *testing.T) {
	if Player1.Mana != 0 {
		t.Errorf("Hero Mana got: %d, want: %d.", Player1.Mana, 0)
	}
}
func TestPlayerManaIncrease(t *testing.T) {
	p := Player1.Copy()
	p.IncreaseMana()
	if p.Mana != 0 {
		t.Errorf("Hero Mana got: %d, want: %d.", p.Mana, 0)
	}
}

func TestPlayerManaIncreaseMax(t *testing.T) {
	p := Player1.Copy()
	for i := 0 ; i < MAX_MANA + 1 ; i++ {
		p.IncreaseMana()
	}
	p.ReloadMana()
	if p.Mana != MAX_MANA {
		t.Errorf("Hero Mana got: %d, want: %d.", p.Mana, MAX_MANA)
	}
}
func TestPlayerRunes(t *testing.T) {
	if Player1.Runes != 25 {
		t.Errorf("Hero Runes got: %d, want: %d.", Player1.Mana, 0)
	}
}
func TestPlayerCardAdd(t *testing.T) {
	p := Player1.Copy()
	p.DeckAddCard(CARDS[0].Copy())
	if p.DeckCount() != 1 {
		t.Errorf("Hero Deck got: %d, want: %d.", p.DeckCount(), 1)
	}	
}
func TestPlayerCardDraw(t *testing.T) {
	p := Player1.Copy()
	c := CARDS[0].Copy()
	c.Id = 1
	p.DeckAddCard(c)
	p.Draw()
	if p.DeckCount() != 0 {
		t.Errorf("Hero Deck got: %d, want: %d.", p.DeckCount(), 0)
	}
	if len(p.Hand) != 1 {
		t.Errorf("Hero Hand got: %d, want: %d.", len(p.Hand), 1)
	}
	if ! reflect.DeepEqual(c, p.Hand[0]) {
		t.Errorf("Hero Hand got wrong card: %v, want: %v.", c, p.Hand[0])
	}		
}
func TestPlayerCardNoDraw(t *testing.T) {
	p := Player1.Copy()
	err := p.Draw()
	if len(p.Hand) != 0 {
		t.Errorf("Hero Hand got: %d, want: %d.", len(p.Hand), 0)
	}
	if err == nil {
		t.Errorf("Hero Player.Draw() should return a error")
	}
}
func TestPlayerCardDrawMax(t *testing.T) {
	p := Player1.Copy()
	for i := 0 ; i < MAX_HAND_CARD + 1 ; i++ {
		p.DeckAddCard(CARDS[i].Copy())
	}
	if p.DeckCount() != MAX_HAND_CARD + 1{
		t.Errorf("Hero Deck got: %d, want: %d.", p.DeckCount(), MAX_HAND_CARD + 1)
	}
	p.DrawN(MAX_HAND_CARD + 1)
	if len(p.Hand) != MAX_HAND_CARD {
		t.Errorf("Hero Hand got: %d, want: %d.", len(p.Hand), MAX_HAND_CARD)
	}	
}

func TestPlayerCardPick(t *testing.T) {
	p := Player1.Copy()
	c := CARDS[0].Copy()
	c.Id = 1
	p.Pick(c)
	
	if len(p.Hand) != 1 {
		t.Errorf("Hero Hand got: %d, want: %d.", len(p.Hand), MAX_HAND_CARD)
	}

	if ! reflect.DeepEqual(c, p.Hand[0]) {
		t.Errorf("Hero Hand Card are different, got: %v, want: %v.", c, p.Hand[0])
	}
}

func TestPlayerCardHandGet(t *testing.T) {
	p := Player1.Copy()
	c := CARDS[0].Copy()
	c.Id = 1
	p.Pick(c)
	
	if c != p.HandGetCard(c.Id) {
		t.Errorf("Hero Hand Card are different, got: %v, want: %v.", c, p.HandGetCard(c.Id))
	}
}

func TestPlayerCardHandRemove(t *testing.T) {
	p := Player1.Copy()
	c := CARDS[0].Copy()
	c.Id = 1
	p.Pick(c)
	p.HandRemoveCard(c.Id)
	if len(p.Hand) != 0 {
		t.Errorf("Hero Hand got: %v, want: [].", p.Hand)
	}
}
func TestPlayerCardHandPlayCardCost(t *testing.T) {
	var err error
	p := Player1.Copy()
	c := CARDS[0].Copy()
	c.Id = 1
	p.SetMaxMana(c.Cost - 1)
	p.ReloadMana()
	p.Pick(c)
	err = p.HandPlayCard(c.Id)
	if len(p.Hand) != 1 {
		t.Errorf("Hero Hand got: %d, want: %d.", len(p.Hand), 1)
	}
	if err == nil {
		t.Errorf("Hero Hand not enought mana should return error")
	}
}
func TestPlayerCardHandPlayCreature(t *testing.T) {
	p := Player1.Copy()
	c := CARDS[0].Copy()
	c.Id = 1
	p.SetMaxMana(c.Cost)
	p.ReloadMana()
	if p.Mana != c.Cost {
		t.Errorf("Hero Mana didn't be use, got: %d, want: %d.", p.Mana, 0)
	}

	p.Pick(c)
	if len(p.Hand) != 1 {
		t.Errorf("Hero Hand got: %d, want: 1.", len(p.Hand), 1)
	}

	p.HandPlayCard(c.Id)
	if len(p.Board) != 1 {
		t.Errorf("Hero Board got: %d, want: %d.", len(p.Board), 1)
	}
	if len(p.Hand) != 0 {
		t.Errorf("Hero Hand got: %d, want: d.", len(p.Hand), 0)
	}
}
func TestPlayerCardHandPlayCreatureMax(t *testing.T) {
	var err error
	p := Player1.Copy()
	for i := 0 ; i < MAX_BOARD_CARD + 1; i++ {
		c := CARDS[0].Copy()
		c.Id = i + 1
		p.SetMaxMana(c.Cost)
		p.ReloadMana()
		p.Pick(c)
		err = p.HandPlayCard(c.Id)
	}

	if len(p.Board) != MAX_BOARD_CARD {
		t.Errorf("Hero Board got: %d, want: %d.", len(p.Board), MAX_BOARD_CARD)
	}
	if err == nil {
		t.Errorf("Hero Hand not enought mana should return error")
	}
}
func TestPlayerCardHandPlayItem(t *testing.T) {
	p := Player1.Copy()
	c := CARDS[len(CARDS) - 1].Copy()
	c.Id = 1
	p.Pick(c)
	if len(p.Hand) != 1 {
		t.Errorf("Hero Hand got: %d, want: %d.", len(p.Hand), 1)
	}
	p.SetMaxMana(c.Cost)
	p.ReloadMana()
	p.HandPlayCard(c.Id)
	if len(p.Hand) != 0 {
		t.Errorf("Hero Hand got: %d, want: %d.", len(p.Hand), 0)
	}
}
func TestPlayerCardBoardGetCreature(t *testing.T) {
	p := Player1.Copy()
	c := CARDS[0].Copy()
	c.Id = 1

	p.SetMaxMana(c.Cost)
	p.Pick(c)
	p.HandPlayCard(c.Id)
	if c != p.HandGetCard(c.Id) {
		t.Errorf("Hero Board got: %v, want: %v.", p.HandGetCard(c.Id), c)
	}
}
func TestPlayerCardBoardGetCreatureGuards(t *testing.T) {
	p := Player1.Copy()
	c1 := CARDS[61].Copy()
	c1.Id = 1

	c2 := CARDS[62].Copy()
	c2.Id = 2

	c3 := CARDS[60].Copy()
	c3.Id = 3

	p.Board = append(p.Board, c1)
	p.Board = append(p.Board, c2)
	p.Board = append(p.Board, c3)
	
	guards := p.BoardGetGuardsId()

	if len(guards) != 2 {
		t.Errorf("Hero Board got: %d, want: %d.", len(guards), 2)
	}
}
func TestPlayerCardBoardRemoveCreature(t *testing.T) {
	p := Player1.Copy()
	c := CARDS[0].Copy()
	c.Id = 1

	p.SetMaxMana(c.Cost)
	p.Pick(c)
	p.HandPlayCard(c.Id)
	p.BoardRemoveCard(c.Id)
	if len(p.Board) != 0 {
		t.Errorf("Hero Board got: %d, want: %d.", len(p.Board), 0)
	}
}

func TestPlayerCardHandStackDraw(t *testing.T) {
	var err error

	p := Player1.Copy()
	c := CARDS[79].Copy()
	c.Id = 1

	p.SetMaxMana(c.Cost)
	p.ReloadMana()
	p.Pick(c)
	err = p.HandPlayCard(c.Id)
	if p.StackCard != 1 {
		t.Errorf("Hero stack card got: %d, want: %d.", p.StackCard, 1)
	}
	if err != nil {
		t.Errorf("%s", err)
	}
}

func TestPlayerCardHandStackDrawN(t *testing.T) {
	var err error

	p := Player1.Copy()
	c := CARDS[35].Copy()
	c.Id = 1

	p.SetMaxMana(c.Cost)
	p.ReloadMana()
	p.Pick(c)
	err = p.HandPlayCard(c.Id)
	if p.StackCard != 2 {
		t.Errorf("Hero stack card got: %d, want: %d.", p.StackCard, 2)
	}
	if err != nil {
		t.Errorf("%s", err)
	}
}

