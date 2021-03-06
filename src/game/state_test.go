package game

import (
	"testing"
	//"reflect"
	"math/rand"
	"time"
)

var Hero *Player 		= NewPlayer(1, STARTING_LIFE, STARTING_MANA, STARTING_RUNES)
var Vilain *Player 		= NewPlayer(2, STARTING_LIFE, STARTING_MANA, STARTING_RUNES)

func pickRandomCard() *Card {
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)
	n := random.Intn(len(CARDS))
	c := CARDS[n]
	return c
}
func emptyDeck(p1 *Player) {
	p1.Deck = make([]*Card, 0)
}
func emptyRunes(p1 *Player) {
	p1.Runes = 0
}
func loadDeck(p1 *Player, start_id int) {
	if p1 != nil {
		for i := 0 ; i < STARTING_CARD ; i++ {
			c := pickRandomCard()
			c.Id = start_id + i
			p1.DeckAddCard(c)
		}
	}
}
func initState() *State {
	p1 := Hero.Copy()
	p2 := Vilain.Copy()

	s  := NewState(p1, p2)

	loadDeck(p1, 1)
	loadDeck(p2, STARTING_CARD + 1)

	return s
}

func TestInitState(t *testing.T) {
	s := initState()

	if s.Hero().DeckCount() != STARTING_CARD {
		t.Errorf("Hero deck count, got: %d, want: %d.", s.Hero().DeckCount(), STARTING_CARD)
	}
	if s.Vilain().DeckCount() != STARTING_CARD {
		t.Errorf("Vilain deck count, got: %d, want: %d.", s.Vilain().DeckCount(), STARTING_CARD)
	}
}
func TestStateNewTurn(t *testing.T) {
	s := initState()

	c1 := CARDS[40].Copy()
	c1.Id = 61

	
	s.Vilain().SetMaxMana(c1.Cost)
	s.Vilain().ReloadMana()
	s.Vilain().Pick(c1)

	previous_vilain := s.Vilain().Copy()
	previous_hero 	:= s.Hero().Copy()

	s.NextTurn()

	if s.Hero().Id != previous_vilain.Id {
		t.Errorf("Hero player id, got: %d, want: %d.", s.Hero().Id, previous_vilain.Id)
	}
	if s.Vilain().Id != previous_hero.Id {
		t.Errorf("Vilain player id, got: %d, want: %d.", s.Vilain().Id, previous_hero.Id)
	}

	if s.Hero().MaxMana != (previous_vilain.MaxMana + 1) {
		t.Errorf("Hero mana, got: %d, want: %d.", s.Hero().MaxMana, previous_vilain.MaxMana + 1)
	}
	if s.Hero().MaxMana != s.Hero().Mana {
		t.Errorf("Hero mana, got: %d, want: %d.", s.Hero().MaxMana, s.Hero().Mana)
	}

	if len(s.Hero().Hand) != (len(previous_vilain.Hand) + 1) {
		t.Errorf("Hero hand, got: %d, want: %d.", len(s.Hero().Hand), (len(previous_vilain.Hand) + 1))
	}
	if s.Hero().DeckCount() != (previous_vilain.DeckCount() - 1) {
		t.Errorf("Hero deck, got: %d, want: %d.", s.Hero().DeckCount(), (previous_vilain.DeckCount() - 1))
	}

	for _, c := range(s.Hero().Board) {
		if c.Attacked {
			t.Errorf("Hero board. Monster should be able to attack, got: %t, want: %t.", c.Attacked, false)
		}
		if c.Charge > 0 {
			t.Errorf("Hero board. Monster should be able to attack, got: %d, want: %d.", c.Charge, 1)
		}
	}
}
func TestStateNewTurnStackDraw(t *testing.T) {
	var err error

	s := initState()

	previous_vilain := s.Vilain().Copy()
	s.Vilain().StackDraw()
	err = s.NextTurn()

	if len(s.Hero().Hand) != (len(previous_vilain.Hand) + 2) {
		t.Errorf("Hero hand, got: %d, want: %d.", len(s.Hero().Hand), (len(previous_vilain.Hand) + 2))
	}
	if s.Hero().DeckCount() != (previous_vilain.DeckCount() - 2) {
		t.Errorf("Hero deck, got: %d, want: %d.", s.Hero().DeckCount(), (previous_vilain.DeckCount() - 2))
	}
	if err != nil {
		t.Errorf("%s", err)
	}
}
func TestStateNewTurnStackDrawN(t *testing.T) {
	var err error

	s := initState()

	previous_vilain := s.Vilain().Copy()
	s.Vilain().StackDrawN(2)
	err = s.NextTurn()

	if len(s.Hero().Hand) != (len(previous_vilain.Hand) + 3) {
		t.Errorf("Hero hand, got: %d, want: %d.", len(s.Hero().Hand), (len(previous_vilain.Hand) + 3))
	}
	if s.Hero().DeckCount() != (previous_vilain.DeckCount() - 3) {
		t.Errorf("Hero deck, got: %d, want: %d.", s.Hero().DeckCount(), (previous_vilain.DeckCount() - 3))
	}
	if err != nil {
		t.Errorf("%s", err)
	}
}
func TestStateSwapPlayer(t *testing.T) {
	s := initState()
	s.SwapPlayers()

	if s.Hero().Id != 2 {
		t.Errorf("Hero player id, got: %d, want: %d.", s.Hero().Id, 2)
	}
	if s.Vilain().Id != 1 {
		t.Errorf("Vilain player id, got: %d, want: %d.", s.Vilain().Id, 1)
	}
}
func TestStateDrawNoCard(t *testing.T) {
	var err error

	s := initState()

	previous_vilain := s.Vilain().Copy()
	emptyDeck(s.Vilain())

	err = s.NextTurn()

	if s.Hero().Life != (previous_vilain.Life - STEP_RUNE) {
		t.Errorf("Hero life, got: %d, want: %d.", s.Hero().Life, previous_vilain.Life - STEP_RUNE)
	}

	if err == nil {
		t.Errorf("Hero NextTurn should return a error")
	}
}
func TestStateDrawNoCardNoRune(t *testing.T) {
	var err error

	s := initState()

	emptyDeck(s.Vilain())
	emptyRunes(s.Vilain())

	err = s.NextTurn()

	if s.Hero().Life != 0 {
		t.Errorf("Hero life, got: %d, want: %d.", s.Hero().Life, 0)
	}

	if err == nil {
		t.Errorf("Hero NextTurn should return a error")
	}
}
func TestStateMoveAttackNoCharge(t *testing.T) {
	var err error

	s := initState()
	c := CARDS[30].Copy()
	c.Id = 61

	previous_vilain := s.Vilain().Copy()
	s.Hero().SetMaxMana(c.Cost)
	s.Hero().ReloadMana()
	s.Hero().Pick(c)
	s.Hero().HandPlayCard(c.Id)
	
	move_attack := NewMove(MOVE_ATTACK, 0, 1, []int{c.Id, -1})
	err = s.Move(move_attack)

	if s.Vilain().Life != previous_vilain.Life {
		t.Errorf("Hero action attack vilain life, got: %d, want: %d.", s.Vilain().Life, previous_vilain.Life)
	}

	if err == nil {
		t.Errorf("Hero action attack not permitted. Monster doesn't have Charge")
	}
}
func TestStateMoveAttackCharge(t *testing.T) {
	var err error

	s := initState()
	c := CARDS[40].Copy()
	c.Id = 61

	previous_vilain := s.Vilain().Copy()
	s.Hero().SetMaxMana(c.Cost)
	s.Hero().ReloadMana()
	s.Hero().Pick(c)
	s.Hero().HandPlayCard(c.Id)
	
	move_attack := NewMove(MOVE_ATTACK, 0, 1, []int{c.Id, -1})
	err = s.Move(move_attack)

	if s.Vilain().Life != previous_vilain.Life - c.Attack {
		t.Errorf("Hero action attack vilain life, got: %d, want: %d.", s.Vilain().Life, previous_vilain.Life - c.Attack)
	}
	
	if err != nil {
		t.Errorf("Hero action attack permitted. Monster should be able to attack. %v. %s", err)
	}
}
func TestStateMoveAttackAlreadyAttack(t *testing.T) {
	var err error

	s := initState()
	c := CARDS[40].Copy()
	c.Id = 61
	c.Attacked = true

	previous_vilain := s.Vilain().Copy()
	s.Hero().SetMaxMana(c.Cost)
	s.Hero().ReloadMana()
	s.Hero().Pick(c)
	s.Hero().HandPlayCard(c.Id)
	
	move_attack := NewMove(MOVE_ATTACK, 0, 1, []int{c.Id, -1})
	err = s.Move(move_attack)

	if s.Vilain().Life != previous_vilain.Life {
		t.Errorf("Hero action attack vilain life, got: %d, want: %d.", s.Vilain().Life, previous_vilain.Life)
	}	
	
	if err == nil {
		t.Errorf("Hero action attack not permitted. Monster doesn't have Charge")
	}
}
func TestStateMoveAttackOpponentWithGuards(t *testing.T) {
	var err error

	s := initState()
	c1 := CARDS[40].Copy()
	c1.Id = 61

	c2 := CARDS[39].Copy()
	c2.Id = 62

	previous_vilain := s.Vilain().Copy()

	s.Hero().SetMaxMana(c1.Cost)
	s.Hero().ReloadMana()
	s.Hero().Pick(c1)
	s.Hero().HandPlayCard(c1.Id)
	
	s.Vilain().SetMaxMana(c2.Cost)
	s.Vilain().ReloadMana()
	s.Vilain().Pick(c2)
	s.Vilain().HandPlayCard(c2.Id)
	
	move_attack := NewMove(MOVE_ATTACK, 0, 1, []int{c1.Id, -1})
	err = s.Move(move_attack)

	if s.Vilain().Life != previous_vilain.Life {
		t.Errorf("Hero action attack vilain life, got: %d, want: %d.", s.Vilain().Life, previous_vilain.Life)
	}	
	
	if err == nil {
		t.Errorf("Hero action attack not permitted. Vilain have guards")
	}
}
func TestStateMoveAttackCreatureWithGuards(t *testing.T) {
	var err error

	s := initState()
	c1 := CARDS[40].Copy()
	c1.Id = 61

	c2 := CARDS[39].Copy()
	c2.Id = 62

	c3 := CARDS[38].Copy()
	c3.Id = 63

	s.Hero().SetMaxMana(c1.Cost)
	s.Hero().ReloadMana()
	s.Hero().Pick(c1)
	s.Hero().HandPlayCard(c1.Id)
	
	s.Vilain().SetMaxMana(c2.Cost)
	s.Vilain().ReloadMana()
	s.Vilain().Pick(c2)
	s.Vilain().HandPlayCard(c2.Id)

	s.Vilain().SetMaxMana(c3.Cost)
	s.Vilain().ReloadMana()
	s.Vilain().Pick(c3)
	s.Vilain().HandPlayCard(c3.Id)
	
	
	move_attack := NewMove(MOVE_ATTACK, 0, 1, []int{c1.Id, c3.Id})
	err = s.Move(move_attack)

	cv := s.Vilain().BoardGetCard(c3.Id)

	if cv.Defense < CARDS[38].Defense {
		t.Errorf("Hero action attack vilain life, got: %d, want: %d.", cv.Defense, CARDS[38].Defense)
	}	
	if err == nil {
		t.Errorf("Hero action attack not permitted. Vilain have guards")
	}
}
func TestStateMoveAttackCreatureWithWard(t *testing.T) {
	var err error

	s := initState()
	ref_drain := CARDS[83]
	c1 := ref_drain.Copy()
	c1.Id = 61

	ref_ward := CARDS[81]
	c2 := ref_ward.Copy()
	c2.Id = 62


	previous_hero := s.Hero().Copy()

	s.Hero().SetMaxMana(c1.Cost)
	s.Hero().ReloadMana()
	s.Hero().Pick(c1)
	s.Hero().HandPlayCard(c1.Id)
	
	s.Vilain().SetMaxMana(c2.Cost)
	s.Vilain().ReloadMana()
	s.Vilain().Pick(c2)
	s.Vilain().HandPlayCard(c2.Id)

	
	
	move_attack := NewMove(MOVE_ATTACK, 0, 1, []int{c1.Id, c2.Id})
	err = s.Move(move_attack)

	if c2.Defense != ref_ward.Defense {
		t.Errorf("Hero action attack vilain ward %d, got: %d, want: %d.", c2.CardNumber, c2.Defense, ref_ward.Defense)
	}
	if c2.IsAbleTo(CARD_ABILITY_WARD) {
		t.Errorf("Hero action attack should remoev ward, got: %s, want: %s.", abilitiesToString(c2.Abilities), "-W")
	}
	if s.Hero().Life != previous_hero.Life {
		t.Errorf("Hero action attack vilain with drain %d, got: %d, want: %d.", c2.CardNumber, s.Hero().Life, previous_hero.Life)
	}	
	if err != nil {
		t.Errorf("Hero action attack are permitted. %s", err)
	}
}
func TestStateMoveAttackGuardWithGuards(t *testing.T) {
	var err error

	s := initState()
	c1 := CARDS[40].Copy()
	c1.Id = 61

	c2 := CARDS[39].Copy()
	c2.Id = 62

	c3 := CARDS[38].Copy()
	c3.Id = 63

	s.Hero().SetMaxMana(c1.Cost)
	s.Hero().ReloadMana()
	s.Hero().Pick(c1)
	s.Hero().HandPlayCard(c1.Id)
	
	s.Vilain().SetMaxMana(c2.Cost)
	s.Vilain().ReloadMana()
	s.Vilain().Pick(c2)
	s.Vilain().HandPlayCard(c2.Id)

	s.Vilain().SetMaxMana(c3.Cost)
	s.Vilain().ReloadMana()
	s.Vilain().Pick(c3)
	s.Vilain().HandPlayCard(c3.Id)
	
	
	move_attack := NewMove(MOVE_ATTACK, 0, 1, []int{c1.Id, c2.Id})
	err = s.Move(move_attack)

	cv := s.Vilain().BoardGetCard(c2.Id)

	if cv.Defense != (CARDS[39].Defense - c1.Attack) {
		t.Errorf("Hero action attack monster %d, got: %d, want: %d.", cv.Id, cv.Defense, (CARDS[39].Defense - c1.Attack))
	}
	if err != nil {
		t.Errorf("Monster should be able to attack. got: %s", err)
	}

}
func TestStateMoveSummon(t *testing.T) {
	var err error

	s := initState()
	c1 := CARDS[40].Copy()
	c1.Id = 61

	
	s.Hero().SetMaxMana(c1.Cost)
	s.Hero().ReloadMana()
	s.Hero().Pick(c1)

	previous_hero := s.Hero().Copy()

	move_summon := NewMove(MOVE_SUMMON, c1.Cost, 1, []int{c1.Id})
	err = s.Move(move_summon)

	if len(s.Hero().Hand) != len(previous_hero.Hand) - 1 {
		t.Errorf("Hero action summon monster %d, got: %d, want: %d.", c1.Id, len(s.Hero().Hand), len(previous_hero.Hand) - 1)
	}
	if len(s.Hero().Board) != len(previous_hero.Board) + 1 {
		t.Errorf("Hero action summon monster %d, got: %d, want: %d.", c1.Id, len(s.Hero().Board), len(previous_hero.Board) + 1)
	}
	if err != nil {
		t.Errorf("Monster could be summon. got: %s", err)
	}
}
func TestStateMoveSummonNoMana(t *testing.T) {
	var err error

	s := initState()
	c1 := CARDS[40].Copy()
	c1.Id = 61

	s.Hero().SetMaxMana(0)
	s.Hero().ReloadMana()
	s.Hero().Pick(c1)
	previous_hero := s.Hero().Copy()

	move_summon := NewMove(MOVE_SUMMON, c1.Cost, 1, []int{c1.Id})
	err = s.Move(move_summon)

	if len(s.Hero().Hand) != len(previous_hero.Hand) {
		t.Errorf("Hero action summon monster %d Hand, got: %d, want: %d.", c1.Id, len(s.Hero().Hand), len(previous_hero.Hand))
	}
	if len(s.Hero().Board) != len(previous_hero.Board) {
		t.Errorf("Hero action summon monster %d Board, got: %d, want: %d.", c1.Id, len(s.Hero().Board), len(previous_hero.Board))
	}
	if err == nil {
		t.Errorf("Monster shouldn't be summon")
	}

}
func TestStateMoveSummonMaxBoard(t *testing.T) {
	var err error

	s := initState()

	for i := 0 ; i < MAX_BOARD_CARD; i++ {
		c1 := CARDS[40].Copy()
		c1.Id = 61 + i

		s.Hero().SetMaxMana(c1.Cost)
		s.Hero().ReloadMana()
		s.Hero().Pick(c1)

		move_summon := NewMove(MOVE_SUMMON, c1.Cost, 1, []int{c1.Id})
		err = s.Move(move_summon)
	
		if len(s.Hero().Hand) != 0 {
			t.Errorf("Hero action summon monster %d - iteration: %d, got: %d, want: %d. %d %s", c1.Id, i + 1, len(s.Hero().Hand), 0, err)
		}
	}
	
	c1 := CARDS[40].Copy()
	c1.Id = 61 + MAX_BOARD_CARD
	s.Hero().SetMaxMana(c1.Cost)
	s.Hero().ReloadMana()
	s.Hero().Pick(c1)
	move_summon := NewMove(MOVE_SUMMON, c1.Cost, 1, []int{c1.Id})
	err = s.Move(move_summon)
	
	if len(s.Hero().Board) != MAX_BOARD_CARD {
		t.Errorf("Hero action summon monster, got: %d, want: %d.", len(s.Hero().Hand), MAX_BOARD_CARD)
	}
	if len(s.Hero().Hand) != 1 {
		t.Errorf("Hero action summon monster %d, got: %d, want: %d.", c1.Id, len(s.Hero().Hand), 1)
	}
	if err == nil {
		t.Errorf("Monster shouldn't be summon")
	}

}
func TestStateMoveUseItemGreen(t *testing.T) {
	var err error

	
	ref_card := CARDS[40]
	for _, c := range(CARDS) {
		if c.Type != CARD_TYPE_ITEM_GREEN { continue }
		s := initState()

		c1 := ref_card.Copy()
		c1.Id = 61

		c2 := c.Copy()
		c2.Id = 62

		s.Hero().SetMaxMana(c1.Cost)
		s.Hero().ReloadMana()
		s.Hero().Pick(c1)
		s.Hero().HandPlayCard(c1.Id)
		
		s.Hero().SetMaxMana(c2.Cost)
		s.Hero().ReloadMana()
		s.Hero().Pick(c2)

		previous_vilain := s.Vilain().Copy()
		previous_hero := s.Hero().Copy()

		move_use := NewMove(MOVE_USE, c2.Cost, 1, []int{c2.Id, c1.Id})
		err = s.Move(move_use)
	
		expected_abilites := ref_card.Abilities | c2.Abilities
		if c1.Abilities != expected_abilites {
			t.Errorf("Hero action use %d on monster %d. Wrong Abilities got: %s, want: %s.", c2.CardNumber, c1.Id, abilitiesToString(c1.Abilities), abilitiesToString(expected_abilites))
		}
		if c1.Attack != (ref_card.Attack + c2.Attack) {
			t.Errorf("Hero action use %d on monster %d. Wrong Attack, got: %d, want: %d.", c2.CardNumber, c1.Id, c1.Attack, ref_card.Attack + c2.Attack)
		}
		if c1.Defense != (ref_card.Defense + c2.Defense) {
			t.Errorf("Hero action use %d on monster %d. Wrong Defense, got: %d, want: %d.", c2.CardNumber, c1.Id, c1.Defense, ref_card.Defense + c2.Defense)
		}
		if s.Hero().Life != previous_hero.Life + c2.HealthChange {
			t.Errorf("Hero action use %d on monster %d. Wrong Hero Life got: %d, want: %d.", c2.CardNumber, c1.Id, s.Hero().Life, previous_hero.Life + c2.HealthChange)
		}
		if s.Vilain().Life != previous_vilain.Life + c2.OpponentHealthChange {
			t.Errorf("Hero action use %d on monster %d. Wrong Vilain Life got: %d, want: %d.", c2.CardNumber, c1.Id, s.Vilain().Life, previous_vilain.Life + c2.OpponentHealthChange)
		}
		if s.Hero().StackCard != previous_hero.StackCard + c2.CardDraw {
			t.Errorf("Hero action use %d on monster %d. Wrong Stack Card got: %d, want: %d.", c2.CardNumber, c1.Id, s.Hero().StackCard, previous_hero.StackCard + c2.CardDraw)
		}
		if err != nil {
			t.Errorf("Use green item should be able to play. got: %s", err)
		}
	}
}
func TestStateMoveUseItemRed(t *testing.T) {
	var err error

	
	ref_card := CARDS[40]
	for _, c := range(CARDS) {
		if c.Type != CARD_TYPE_ITEM_RED { continue }
		s := initState()

		c1 := ref_card.Copy()
		c1.Id = 61

		c2 := c.Copy()
		c2.Id = 62

		s.Vilain().SetMaxMana(c1.Cost)
		s.Vilain().ReloadMana()
		s.Vilain().Pick(c1)
		s.Vilain().HandPlayCard(c1.Id)
		
		s.Hero().SetMaxMana(c2.Cost)
		s.Hero().ReloadMana()
		s.Hero().Pick(c2)

		previous_vilain := s.Vilain().Copy()
		previous_hero := s.Hero().Copy()

		move_use := NewMove(MOVE_USE, c2.Cost, 1, []int{c2.Id, c1.Id})
		err = s.Move(move_use)
	
		expected_abilites := ref_card.Abilities & ^c2.Abilities
		if c1.Abilities != expected_abilites {
			t.Errorf("Hero action use %d on monster %d. Wrong Abilities got: %s, want: %s.", c2.CardNumber, c1.Id, abilitiesToString(c1.Abilities), abilitiesToString(expected_abilites))
		}
		if c1.Attack != (ref_card.Attack + c2.Attack) {
			t.Errorf("Hero action use %d on monster %d. Wrong Attack, got: %d, want: %d.", c2.CardNumber, c1.Id, c1.Attack, ref_card.Attack + c2.Attack)
		}
		if c1.Defense != (ref_card.Defense + c2.Defense) {
			t.Errorf("Hero action use %d on monster %d. Wrong Defense, got: %d, want: %d.", c2.CardNumber, c1.Id, c1.Defense, ref_card.Defense + c2.Defense)
		}
		if s.Hero().Life != previous_hero.Life + c2.HealthChange {
			t.Errorf("Hero action use %d on monster %d. Wrong Hero Life got: %d, want: %d.", c2.CardNumber, c1.Id, s.Hero().Life, previous_hero.Life + c2.HealthChange)
		}
		if s.Vilain().Life != previous_vilain.Life + c2.OpponentHealthChange {
			t.Errorf("Hero action use %d on monster %d. Wrong Vilain Life got: %d, want: %d.", c2.CardNumber, c1.Id, s.Vilain().Life, previous_vilain.Life + c2.OpponentHealthChange)
		}
		if s.Hero().StackCard != previous_hero.StackCard + c2.CardDraw {
			t.Errorf("Hero action use %d on monster %d. Wrong Stack Card got: %d, want: %d.", c2.CardNumber, c1.Id, s.Hero().StackCard, previous_hero.StackCard + c2.CardDraw)
		}
		if c1.Defense <= c2.Defense && len(s.Vilain().Board) != 0 {
			t.Errorf("Hero action use %d on monster %d. Creature should be dead.")
		}
		if err != nil {
			t.Errorf("Use red item should be able to play. got: %s", err)
		}
	}
	
}
func TestStateMoveUseItemBlue(t *testing.T) {
	var err error

	
	ref_card := CARDS[40]
	for _, c := range(CARDS) {
		if c.Type != CARD_TYPE_ITEM_BLUE { continue }
		s := initState()

		c1 := ref_card.Copy()
		c1.Id = 61

		c2 := c.Copy()
		c2.Id = 62

		s.Vilain().SetMaxMana(c1.Cost)
		s.Vilain().ReloadMana()
		s.Vilain().Pick(c1)
		s.Vilain().HandPlayCard(c1.Id)
		
		s.Hero().SetMaxMana(c2.Cost)
		s.Hero().ReloadMana()
		s.Hero().Pick(c2)

		previous_vilain := s.Vilain().Copy()
		previous_hero := s.Hero().Copy()

		move_use := NewMove(MOVE_USE, c2.Cost, 1, []int{c2.Id, c1.Id})
		err = s.Move(move_use)
	
		expected_abilites := ref_card.Abilities & ^c2.Abilities
		if c1.Abilities != expected_abilites {
			t.Errorf("Hero action use %d on monster %d. Wrong Abilities got: %s, want: %s.", c2.CardNumber, c1.Id, abilitiesToString(c1.Abilities), abilitiesToString(expected_abilites))
		}
		if c1.Attack != (ref_card.Attack + c2.Attack) {
			t.Errorf("Hero action use %d on monster %d. Wrong Attack, got: %d, want: %d.", c2.CardNumber, c1.Id, c1.Attack, ref_card.Attack + c2.Attack)
		}
		if c1.Defense != (ref_card.Defense + c2.Defense) {
			t.Errorf("Hero action use %d on monster %d. Wrong Defense, got: %d, want: %d.", c2.CardNumber, c1.Id, c1.Defense, ref_card.Defense + c2.Defense)
		}
		if s.Hero().Life != previous_hero.Life + c2.HealthChange {
			t.Errorf("Hero action use %d on monster %d. Wrong Hero Life got: %d, want: %d.", c2.CardNumber, c1.Id, s.Hero().Life, previous_hero.Life + c2.HealthChange)
		}
		if s.Vilain().Life != previous_vilain.Life + c2.OpponentHealthChange {
			t.Errorf("Hero action use %d on monster %d. Wrong Vilain Life got: %d, want: %d.", c2.CardNumber, c1.Id, s.Vilain().Life, previous_vilain.Life + c2.OpponentHealthChange)
		}
		if s.Hero().StackCard != previous_hero.StackCard + c2.CardDraw {
			t.Errorf("Hero action use %d on monster %d. Wrong Stack Card got: %d, want: %d.", c2.CardNumber, c1.Id, s.Hero().StackCard, previous_hero.StackCard + c2.CardDraw)
		}
		if c1.Defense <= c2.Defense && len(s.Vilain().Board) != 0 {
			t.Errorf("Hero action use %d on monster %d. Creature should be dead.")
		}
		if err != nil {
			t.Errorf("Use blue item should be able to play. got: %s", err)
		}
	}
	
}





