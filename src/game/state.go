package game

import (
	"fmt"
	"os"
	"errors"
	"strings"
	"math/rand"
	"time"

)
 const (
	WEIGHT_BREAKTHROUGH = 1
	WEIGHT_CHARGE		= 0.5
	WEIGHT_DRAIN		= 1.5
	WEIGHT_GUARD		= 2
	WEIGHT_LETHAL		= 2
	WEIGHT_WARD			= 1


	MOVE_PASS			= 0
	MOVE_PICK			= 1
	MOVE_SUMMON			= 2
	MOVE_ATTACK 		= 3
	MOVE_USE			= 4

	OUTCOME_WIN			= 100
	OUTCOME_LOSE		= -100

 )
/*
type State interface {
	AvailableMoves() 	[]Move
	Copy() 				State
	GameOver()			bool
	Move(m Move)		State
	EvaluationScore()	float64

	IsEndTurn()			bool
	Print()
}

type Move interface {
	Probability()	float64
}
*/

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
			fmt.Fprintln(os.Stderr, "[GAME] Board", c)
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

