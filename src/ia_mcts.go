package main

import (
	"errors"
	"fmt"
	"math/rand"
	"reflect"
	"strings"
	"time"
	"math"
	"os"

	//copier "github.com/jinzhu/copier"
)

const (

	MC_MAX_ITERATION 			= 100																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																									
	MC_MAX_SIMULATION			= 1
	MC_MAX_SIMULATION_REPEAT 	= 50

	MC_OUTCOME_WIN		= 50
	MC_OUTCOME_LOSE		= 0

	BIAS_PARAMETER		= 0.6

	MAX_MANA			= 12
	MAX_PLAYERS			= 2
	MIN_PLAYERS			= 2
	MAX_HAND_CARD		= 8
	MAX_BOARD_CARD		= 6

	STARTING_MANA		= 0
	STARTING_LIFE	= 30
	STARTING_CARD	= 30
	STARTING_RUNES	= 5
	DECK_CARDS		= 30

	DRAFT_PICK		= 3

	CARD_TYPE_CREATURE	        = 0
	CARD_TYPE_ITEM_GREEN	    = 1
	CARD_TYPE_ITEM_RED	        = 2
	CARD_TYPE_ITEM_BLUE	        = 3

	CARD_ABILITY_BREAKTHROUGH   = 0
	CARD_ABILITY_CHARGE         = 1
	CARD_ABILITY_DRAIN       	= 2
	CARD_ABILITY_GUARD          = 3
	CARD_ABILITY_LETHAL         = 4
	CARD_ABILITY_WARD           = 5

	MOVE_PASS	= 0
	MOVE_PICK	= 1
	MOVE_SUMMON	= 2
	MOVE_ATTACK = 3
	MOVE_USE	= 4

)

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

var CARDS_COUNT = len(CARDS)

func PickRandomCard() (*Card) {
	source	:= rand.NewSource(time.Now().UnixNano())
	random	:= rand.New(source)
	idx	:= random.Intn(CARDS_COUNT)
	return CARDS[idx]
}

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

type Deck struct {
	Cards []*Card
}

type Card struct {
	CardNumber 				int
	Id    					int
	Location				int
	Type  					int
	Cost  					int
	Attack					int
	Defense					int
	Abilities				[]string
	HealthChange			int
	OpponentHealthChange 	int
	CardDraw				int
	Charge					int
	Attacked				bool
}
type Player struct {
	Id    int
	Deck  *Deck
	Life  int
	Mana  int
	Board []*Card
	Hand  []*Card
	Runes	int

	stack_draw int
	current_mana	int

}
type State struct {
	players     []*Player
	Turn		int
	Draft		[]*Card

	HeroId		int
	AMoves 		[]*Move

}


type Node struct {
	Id		 int
	Parent   *Node
	Children []*Node
	State    *State
	Wins     float64
	Visits   int
	ByMove		*Move
	EndTurn  bool
	UnexploreMoves	[]*Move
}


type Move struct {
	Cost			int
	Type			int
	Params			[]int
	Probability		int
}

func NewMove(type_, cost, probability int, params []int) *Move {
	return &Move {
		Cost: cost,
		Type: type_,
		Probability: probability,
		Params: params,
	}
}

func (m *Move) Copy() *Move {
	move := &Move {
		Cost: m.Cost,
		Type: m.Type,
		Probability: m.Probability,
		Params: make([]int, len(m.Params)),
	}
	copy(move.Params, m.Params)
	return move
}

func (m *Move) toString() string {
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
func NewNode(parent *Node, state *State, move *Move) *Node {
	return &Node{
		Id:		  1,
		Parent:   parent,
		Children: nil,
		State:    state,
		Wins:     0,
		Visits:   0,
		ByMove: move,
		EndTurn: false,
		UnexploreMoves: nil,
	}
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
	f.WriteString("<TD>WINS</TD>")
	f.WriteString("<TD>VISITS</TD>")
	f.WriteString("<TD>ENDTURN</TD>")
	f.WriteString("<TD>UM</TD>")
	f.WriteString("</TR>")
	f.WriteString("<TR>")
	f.WriteString(fmt.Sprintf("<TD>%f</TD>", MCCalculateScore(n)))
	f.WriteString(fmt.Sprintf("<TD>%f</TD>", n.Wins))
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
	
	for _, p := range(n.State.players) {
		f.WriteString("<TR>")
		f.WriteString("<TD>")

		f.WriteString(fmt.Sprintf("%d", p.Id))
		f.WriteString("</TD>")
		for i, v := range(p.Raw()) {
			f.WriteString("<TD>")
			if i == 1 {
				f.WriteString(fmt.Sprintf("%d/%d", p.current_mana, v))
			} else {
				f.WriteString(fmt.Sprintf("%d", v))
			}
			f.WriteString("</TD>")
		}
		f.WriteString("</TR>")
		for _, c := range(p.Hand) {
			f.WriteString("<TR>")
			f.WriteString(fmt.Sprintf("<TD COLSPAN=\"%d\">", len(col_name)))
			card_str := fmt.Sprintf("%d %d %d %d %d %d %d %s %d %d %d",
									c.CardNumber,
									c.Id,
									c.Location,
									c.Type,
									c.Cost,
									c.Attack,
									c.Defense,
									strings.Join(c.Abilities, ""),
									c.HealthChange,
									c.OpponentHealthChange,
									c.CardDraw)
			f.WriteString(fmt.Sprintf("Hand %s", card_str))
			f.WriteString("</TD>")
			f.WriteString("</TR>")
		}
		for _, c := range(p.Board) {
			f.WriteString("<TR>")
			f.WriteString(fmt.Sprintf("<TD COLSPAN=\"%d\">", len(col_name)))
			card_str := fmt.Sprintf("%d %d %d %d %d %d %d %s %d %d %d %d %t",
									c.CardNumber,
									c.Id,
									c.Location,
									c.Type,
									c.Cost,
									c.Attack,
									c.Defense,
									strings.Join(c.Abilities, ""),
									c.HealthChange,
									c.OpponentHealthChange,
									c.CardDraw,
									c.Charge,
									c.Attacked)
			f.WriteString(fmt.Sprintf("Board %s", card_str))
			f.WriteString("</TD>")
			f.WriteString("</TR>")
		}
	}
	f.WriteString("</TABLE>>]\n")


	for i, c := range(n.Children) {
		str := fmt.Sprintf("%d -- %d  [label=\"%s\"]", id, id * 100 + 1 + i, c.ByMove.toString())
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
		if  tmp_m.Type == m.Type &&
			tmp_m.Cost == m.Cost &&
			tmp_m.Probability == m.Probability &&
			reflect.DeepEqual(tmp_m.Params, m.Params) {
			idx = i
			break
		}
	}
	len_moves := len(n.UnexploreMoves)
	if idx != -1 && len_moves > 1 {
		n.UnexploreMoves[idx] = n.UnexploreMoves[len_moves - 1]
		n.UnexploreMoves = n.UnexploreMoves[:len_moves - 1]
	} else if idx != -1 {
		n.UnexploreMoves = make([]*Move, 0)
	}
	return nil
}
func (n *Node) AddChild(node *Node) {
	if n.Children == nil {
		n.Children = make([]*Node, 0)
	}

	n.Children = append(n.Children, node)
}
func (n *Node) UpdateScore(score float64) {
	n.Visits++
	n.Wins += score
}
func (n *Node) Print() {
	
}
func MonteCarloMoves(root_node *Node, timeout int) []*Move {
	
	var n *Node
	if timeout > 5 {
		n = MonteCarloTimeout(root_node, timeout)
	} else {
		n = MonteCarlo(root_node)
	}
	if n == nil {
		return nil
	}
	//fmt.Println("Node count:", CountNode(n))
	//root_node.ExportGraph(fmt.Sprintf("graph-%d", 1))
	
	moves := make([]*Move, 0)
	moves = append(moves, n.ByMove)
	
	//fmt.Fprintln(os.Stderr, "[MCTS] Select move", n.ByMove.toString())
	
	
	for ; n != nil && len(n.Children) > 0 ; {
		
		if n.EndTurn { break }

		var score  float64 = -100

		for _, node := range n.Children {
	//		fmt.Fprintln(os.Stderr, "[MCTS] AMove:", node.ByMove.toString())
			child_score := MCCalculateScore(node)
			if child_score > score {
				score = child_score
				n = node
			}
		}
		//fmt.Fprintln(os.Stderr, "[MCTS] Go Deeper. Select move", n.ByMove.toString())
		moves = append(moves, n.ByMove)	

	}
/*
	fmt.Println("---------------------------------------")
	for _, node := range n.Children {
		fmt.Println(node.ByMove.toString())
		node.State.PrintHero()
		node.State.PrintHeroBoard()
	}
	fmt.Println("---------------------------------------")
*/
	return moves
}

func MonteCarloTimeout(root_node *Node, timeout int) *Node {
	var node *Node
	var candidate_node *Node = nil
	var score  float64 = -100

	var done chan bool 
	go func() {
		
		for {
			
			//fmt.Fprintln(os.Stderr, "0) =========", i, root_node, root_node.Wins, "==========")
			node = MCSelection(root_node)
			node = MCExpansion(node)
	
			/*if len(root_node.UnexploreMoves) == 0 && len(root_node.Children) == 1 {
				return node
			}*/
			score := MCSimulation(node.State)
			MCBackPropagation(node, score)
			
		}
		
		done <- true
	}()

	select {
	case <- time.After(time.Nanosecond * time.Duration(timeout - 5) * (1000000)):
		//fmt.Println("Timeout")
		break
	case <-done:
	}
	
	for _, n := range root_node.Children {
		child_score := MCCalculateScore(n)
	//	fmt.Println(os.Stderr, "[MCTS][SCORE]", n.ByMove.toString(), child_score)
		if child_score > score || candidate_node == nil {
			score = child_score
			candidate_node = n
		}
	}
	
	//fmt.Println("[MCTS] Best node", candidate_node)
	if candidate_node != nil {
		return candidate_node
	}
	
	return nil
}
func MonteCarlo(root_node *Node) *Node {
	var node *Node
	//root_node.ExportGraph(fmt.Sprintf("graph-%d", 0))
	for i := 0; i < MC_MAX_ITERATION; i++ {
		
		root_node.Id = i + 1
		//fmt.Fprintln(os.Stderr, "0) =========", i, root_node, root_node.Wins, "==========")
		node = MCSelection(root_node)
		node = MCExpansion(node)

		/*if len(root_node.UnexploreMoves) == 0 && len(root_node.Children) == 1 {
			return node
		}*/
		score := MCSimulation(node.State)
		MCBackPropagation(node, score)
		//root_node.ExportGraph(fmt.Sprintf("graph-%d", i + 1))
	}

	var candidate_node *Node = nil
	var score  float64 = -100

	for _, n := range root_node.Children {
		child_score := MCCalculateScore(n)
	//	fmt.Println(os.Stderr, "[MCTS][SCORE]", n.ByMove.toString(), child_score)
		if child_score > score || candidate_node == nil {
			score = child_score
			candidate_node = n
		}
	}

	//fmt.Println("[MCTS] Best node", candidate_node)
	if candidate_node != nil {
		return candidate_node
	}
	return nil
}
func MCSelection(node *Node) *Node {
	var candidate_node *Node

	/*
	for _, m := range(node.State.AvailablesMoves()) {
		//fmt.Println("[MCTS][AMOVES]", m.toString())
	}
	*/
	if node.UnexploreMoves == nil {
		node.State.AvailablesMoves()
		node.UnexploreMoves = node.State.CopyAvailablesMoves()
	}


	if len(node.UnexploreMoves) == 0 && node.Children != nil && len(node.Children) > 0 {
		candidate_node = nil
		score := -100.0
		//fmt.Println("[MCTS] Select node with action:", node.ByMove.toString())
		for _, n := range node.Children {
			child_score := MCCalculateScore(n)
			//fmt.Println("[MCTS][SCORE]", child_score)
			if child_score > score {
				score = child_score
				candidate_node = n
			}
		}
		if candidate_node == nil {
			return node
		}
		//fmt.Fprintln(os.Stderr, "[MCTS][SELECT] Select", candidate_node)
		return MCSelection(candidate_node)
	}

/*
	fmt.Println("[MCTS] Select action", *node)
	if node.ByMove != nil {
		fmt.Println("[MCTS] Corresponding to", node.ByMove.toString())
	}
*/
//	node.State.Print()
	//fmt.Fprintln(os.Stderr, "[MCTS][SELECT] Select", node)
	return node
}
func MCCalculateScore(node *Node) float64 {
	if node.Parent == nil {
		return 0
	}
	exploitScore := float64(node.Wins) / float64(node.Visits)
	exploreScore := math.Sqrt(2 * math.Log(float64(node.Parent.Visits)) / float64(node.Visits))
	exploreScore = BIAS_PARAMETER * exploreScore

	return exploitScore + exploreScore
}
func MCExpansion(node *Node) *Node {
	

	if node.UnexploreMoves == nil {
		fmt.Fprintln(os.Stderr, "[MCTS][EXPANSION] UnexploreMove are nil and should not be")
		return node
	}

	if len(node.UnexploreMoves) == 0 {
		return node
	}

	/*
	for _, m := range(node.UnexploreMoves) {
		fmt.Println("[MCTS] Available move", m.toString())
	}
	*/
	
	new_state := node.State.Copy()

	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)
	rmove := node.UnexploreMoves[random.Intn(len(node.UnexploreMoves))]

	new_state.MoveHero(rmove)
	new_node := NewNode(node, new_state, rmove)
	node.DeleteUnexploreMoves(rmove)

	new_node.Parent = node
	//fmt.Println("[MCTS] Expand", *node, node.State.Hero().Mana, "with action:", rmove.toString())
	//fmt.Fprintln(os.Stderr, "[MCTS][EXPANSION] Extend", node, "with", new_node)
	node.Children = append(node.Children, new_node)
		
	if new_state.IsEndTurn() && new_node.State.Hero().Id != node.State.Hero().Id {
		new_node.EndTurn = true
		new_state.NextTurn()
		//fmt.Println("Send end turn")

	} else if new_state.IsEndTurn() {
		//fmt.Println("SHOULD NOT HAPPEN END TURN")
		//fmt.Fprintln(os.Stderr, "[MCTS][EXPANSION] End turn")
		new_state.NextTurnHero(node.State.Hero().Id)
	}


	
	return new_node
}
func MCSimulation(state *State) float64 {

	var avg float64 = 0

	for i := 1 ; i <= MC_MAX_SIMULATION_REPEAT ; i++ {
		simulate_state := state.Copy()

		iteration := 0
		//fmt.Fprintln(os.Stderr, "[MCTS] Start simulation with state")
		//state.Print()
		/*
		for _, m := range(simulate_state.AvailablesMoves()) {
			//fmt.Fprintln(os.Stderr, "[MCTS][SIMULATION] AMoves:", m.toString())
		}
		*/
		for simulate_state.GameOver() == nil && iteration < MC_MAX_SIMULATION {
			/*for _, m := range(simulate_state.AvailablesMoves()) {
				//fmt.Println("[MCTS][SIMULATION] ", m.toString())
			}
			*/
			if simulate_state.GameOver() != nil  {
				break
			}

			simulate_state.RandomMoveHero(simulate_state.Hero().Id)
			//fmt.Fprintln(os.Stderr, "[MCTS] Simulate", iteration, "choose action", m.toString())
	/*
			if simulate_state.IsEndTurn() {
				//fmt.Fprintln(os.Stderr, "[MCTS][SIMULATION] End turn")
				simulate_state.NextTurnHero(simulate_state.Hero().Id)
			} 
	*/	
			iteration++
		}
		avg += MCEvaluation(simulate_state)
	}
	//fmt.Println("[MCTS] Simulation stop at state;")
	//simulate_state.Print()
	return avg / float64(MC_MAX_SIMULATION_REPEAT)
}

func MCEvaluationBoard(s *State) float64 {
	/*
	if s.Hero().Life <= 0 {
		return 0
	} else if s.Vilain().Life <= 0 {
		return 1
	}
	*/
	var score float64 = 0
	var hpow, hdef, vpow, vdef float64
	var hc, vc, hl, vl, hm, vm float64

	hc = float64(len(s.Hero().Hand))
	hl = float64(s.Hero().Life)
	hm = float64(len(s.Hero().Board))
	for _, c := range(s.Hero().Board) {
		hpow += float64(c.Attack)
		hdef += float64(c.Defense)
	}

	vc = float64(len(s.Vilain().Hand))
	vl = float64(s.Vilain().Life)
	vm = float64(len(s.Vilain().Board))
	for _, c := range(s.Vilain().Board) {
		vpow += float64(c.Attack)
		vdef += float64(c.Defense)
	}

	//more_power := (hpow - vdef) / (hpow + vdef + 1)
	more_power := hpow * 1.2 + hdef * 0.8 - (vpow * 1.2 + vdef * 0.8)
    //more_defense := (hdef - vpow) / (vpow + hdef + 1)
	//more_cards := (hc - vc) / (hc + vc + 1)
	more_cards := hc - vc
	//more_life := ((hl - vl) / (hl + vl)) / 2 + 0.5
	more_life := hl - vl
	more_monster := hm - vm
	if vl <= 0 {
		more_life = MC_OUTCOME_WIN
	}
	score += more_life * 5
	score += more_power * 5
	//score += more_defense * 0.05
	score += more_cards * 0.3
	score += more_monster * 1

	if score < 0 {
		//fmt.Fprintln(os.Stderr, "Negative score", score, "=", more_power, more_life, more_cards)
		return MC_OUTCOME_LOSE
	}

	//fmt.Println("Evaluate score", score)
	return score
	
}
func MCEvaluation(s *State) float64 {
	if s != nil {
		return MCEvaluationBoard(s)
		/*
		if s.GameOver() != nil {
			if s.Hero().Life > s.Vilain().Life {
				return 1
			} else if s.Hero().Life < s.Vilain().Life {
				return 0
			} 
		} else {
			
			return MCEvaluationBoard(s)
		}
		*/
	}
	return 0
}
func MCBackPropagation(node *Node, score float64) *Node {

	for node.Parent != nil {
		node.UpdateScore(score)
		node = node.Parent
	}

	node.Visits++

	return node
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
	return &Card{
		CardNumber: 	cardNumber,
		Id:             id,
		Location:		0,
		Type:           type_,
		Cost:           cost,
		Attack:         attack,
		Defense:        defense,
		Abilities:      strings.Split(abilities, ""),
		HealthChange: 	heroHealthChange,
		OpponentHealthChange: opponentHealthChange,
		CardDraw:		cardDraw,
		Charge: 	0,
		Attacked: false,
	}
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
		strings.Join(c.Abilities, ""),
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
		Abilities:		make([]string, len(c.Abilities)),
		OpponentHealthChange: c.OpponentHealthChange,
		CardDraw:		c.CardDraw,
		Charge:			c.Charge,
		Attacked:		c.Attacked,
	}

	copy(new_c.Abilities, c.Abilities)
	return new_c
}


/*************************************/
/************ DECK METHODS ***********/
/*************************************/
func NewDeck() *Deck {
	return &Deck{
		Cards: make([]*Card, 0),
	}
}
func (d *Deck) Count() int {
	return len(d.Cards)
}
func (d *Deck) AddCard(c *Card) {
	d.Cards = append(d.Cards, c)
}
func (d *Deck) FillRandom(n int) {
	d.Cards = make([]*Card, n)
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)

	for i := 0 ; i < n ; i++ {
		random_card := CARDS[random.Intn(CARDS_COUNT)]
		d.Cards[i] = random_card.Copy()
	} 
}
func (d *Deck) Draw() (*Card, error) {
	if len(d.Cards) > 0 {
		c := d.Cards[0]
		d.Cards = d.Cards[1:]
		return c, nil
	} else {
		return nil, errors.New("There is no more card in the deck")
	}
}
func (d *Deck) Copy() *Deck {
	new_deck := NewDeck()
	new_deck.Cards = make([]*Card, len(d.Cards))
	for i, c := range(d.Cards) {
		new_deck.Cards[i] = c.Copy()
	}

	return new_deck
}

func (d *Deck) Shuffle() {
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)
	for i := len(d.Cards) - 1; i > 0; i-- {
		j := random.Intn(i + 1)
		d.Cards[i], d.Cards[j] = d.Cards[j], d.Cards[i]
	}
}

/*************************************/
/********** PLAYERS METHODS **********/
/*************************************/

func NewPlayer(id, life, mana, runes int) *Player {
	return &Player{
		Id:     id,
		Mana:   mana,
		Life:   life,
		Deck: 	NewDeck(),
		Board:  make([]*Card, 0),
		Hand:   make([]*Card, 0),
		Runes: runes,
		current_mana: 0,
		stack_draw: 0,
	}

}
func (p *Player) Copy() *Player {
	new_player := &Player{
		Id:     p.Id,
		Mana:   p.Mana,
		Life:   p.Life,
		Deck: 	p.Deck.Copy(),
		Board:  make([]*Card, len(p.Board)),
		Hand:   make([]*Card, len(p.Hand)),
		Runes: p.Runes,
		stack_draw: p.stack_draw,
		current_mana: p.current_mana,
	}
	for i, _ := range(p.Board) {
		new_player.Board[i] = p.Board[i].Copy()
	}
	for i, _ := range(p.Hand) {
		new_player.Hand[i] = p.Hand[i].Copy()
	}
	return new_player
}
func (p *Player) ReloadMana () {
	p.current_mana = p.Mana
}
func (p *Player) ReloadCreature() {
	for _, c := range(p.Board) {
		c.Attacked 	= false
		c.Charge 	= 1
	}
}
func (p *Player) CurrentMana() int {
	return p.current_mana
}
func (p *Player) Raw() []interface{} {
	return []interface{}{
		p.Life,
		p.Mana,
		p.Deck.Count(),
		p.Runes,
	}
}
func (p *Player) LoseLifeToNextRune() {
	if p.Runes > 0 && p.Life >= STARTING_RUNES * p.Runes {
		var damage int
		damage = p.Life - (STARTING_RUNES  * (p.Runes - 1))
		p.Runes -= 1
		////////fmt.Fprintln(os.Stderr, "[GAME][RUNE] Player", p.Id, "can't draw card. Losing", damage, "damage")
		p.ReceiveDamage(damage)
	}
}
func (p *Player) DrawCard() (error) {
	if len(p.Hand) >= MAX_HAND_CARD {
		//fmt.Println("[GAME][DECK] Maximum card hand reach", MAX_HAND_CARD)
		return fmt.Errorf("[GAME][DECK] Maximum card hand reach %d", MAX_HAND_CARD)
	}
	c, err := p.Deck.Draw()
	if err == nil {
		//fmt.Println("[GAME][DECK]: Player", p.Id, "draw card", c)
		p.Draw(c)
	} else {
		p.LoseLifeToNextRune()
		//fmt.Println("[GAME][INFO]:", err)
		return err
	}

	return nil
}
func (p *Player) DrawCardN(n int) (err error) {
	for i := 0 ; i < n ; i++ {
		err = p.DrawCard()
		if err != nil {
			return err
		}
	}
	return nil
}
func (p *Player) Draw(c *Card) {
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
func (p *Player) SetMana(mana int) {
	if mana <= 12 && mana >= 0 {
		p.Mana = mana
	}
}
func (p *Player) SetLife(life int) {
	p.Life = life
	p.CheckRune()
}
func (p *Player) StackDraw() {
	//fmt.Println("[GAME][RUNE] Player", p.Id, "stacking draw card")
	p.stack_draw++
}
func (p *Player) DrawStackCards() (err error) {
	max := p.stack_draw
	for i := 0 ; i < max ; i++ {
		err = p.DrawCard()
		if err != nil && p.Life <= 0 {
			return err
		}
		p.stack_draw--
	}
	return nil

}
func (p *Player) CheckRune() {
	if p.Life < STARTING_LIFE && p.Life <= (STARTING_RUNES * p.Runes) {		
		for ; p.Life <= (STARTING_RUNES * p.Runes) && p.Runes > 0 ; {
			p.Runes -= 1
			//fmt.Println("[GAME][RUNE] Losing a rune. There are", p.Runes, "left")
			p.StackDraw()
		}
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

	if c.Cost > p.current_mana {
		return fmt.Errorf("[PLAYER] No more mana (%d) for playing card %d with cost %d", p.current_mana, id, c.Cost)
	}

	p.current_mana -= c.Cost

	switch c.Type {
	case CARD_TYPE_CREATURE:
		p.Board = append(p.Board, c)
	case CARD_TYPE_ITEM_BLUE:
	case CARD_TYPE_ITEM_GREEN:
	case CARD_TYPE_ITEM_RED:

	default:
		return fmt.Errorf("Unkow card type %d for player %d", c.Type, p.Id)
	}

	////////fmt.Fprintln(os.Stderr, "[GAME][PLAYER] Player", p.Id, "play card", c)
	p.HandRemoveCard(id)
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
func (p *Player) BoardGetGuardsId() []int {
	ids := make([]int, 0)
	for _, c := range(p.Board) {
		if c.Abilities[CARD_ABILITY_GUARD] != "-" {
			ids = append(ids, c.Id)
		}
	}
	return ids
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
	} else if damage < 0 {
		p.GainLife(-damage)
	}
}

func (p *Player) IncreaseMana() {
	if p.Mana < MAX_MANA {
		p.Mana = p.Mana + 1
		////fmt.Println("Increase Mana")
	}
}

/***************************************/
/************ STATE METHODS ************/
/***************************************/
func NewState() *State {
	return &State{
		players: make([]*Player, 0),
		Turn: 1,
		AMoves: nil,
	}
}
func InitState(p1, p2 *Player) (state *State) {
	state = NewState()
	state.players = append(state.players, []*Player{p1, p2}...)
	return state
}



func (s *State) Copy() (state *State) {
	state = NewState()
	state.players = make([]*Player, len(s.players))
	for i, _ := range(s.players) {
		state.players[i] = s.players[i].Copy()
	}

	if s.AMoves != nil {
		state.AMoves = make([]*Move, len(s.AMoves))
		for i, _ := range(s.AMoves) {
			state.AMoves[i] = s.AMoves[i].Copy()
		}
	}

	return state
}
func (s *State) AddPlayer(p *Player) error {
	if len(s.players) >= MAX_PLAYERS {
		return errors.New("AddPlayer: There is already 2 players")
	}

	s.players = append(s.players, p)
	return nil
}
func (s *State) Hero() (*Player) {
	return s.players[0]
}
func (s *State) Vilain() (*Player) {
	return s.players[1]
}
func (s *State) NextTurnHero(id_hero int) (winner * Player, err error) {
	next_turn_hero := false

	if s.IsEndTurn() && s.Hero().Id == id_hero { 
		s.NextTurn()
		next_turn_hero = s.Hero().Id == id_hero
		for ; ! next_turn_hero ; { 
			s.RandomMove()
			//fmt.Fprintln(os.Stderr, "[MCTS][VILAIN] Play move:", m.toString())
			if s.IsEndTurn() || s.GameOver() != nil {
				//fmt.Fprintln(os.Stderr, "[MCTS][VILAIN] End turn")
				next_turn_hero = true
				break
			}
			//s.Print()

		}
		/*
		if s.GameOver() == nil && s.Hero().Id != id_hero {
			s.NextTurn()
		}
		*/
	} else if s.IsEndTurn() {
		s.NextTurn()
	}
	return s.GameOver(), nil

}

func (s *State) NextTurn() (winner *Player, err error) {
	s.players[0], s.players[1] = s.players[1], s.players[0]
	s.Turn++

	s.Hero().IncreaseMana()
	s.Hero().ReloadMana()
	s.Hero().ReloadCreature()
	s.AMoves = nil

	err = s.Hero().DrawStackCards()

	if err != nil {
		winner = s.Vilain()
		return winner, nil
	}

	err = s.Hero().DrawCard()
	if err != nil && s.Hero().Life <= 0 {
		winner = s.Vilain()
		return winner, nil
	}

	return nil, nil
}
func (s *State) PrintHand(p *Player) {
	if p != nil {
		for _, c := range(p.Hand) {
			fmt.Fprintln(os.Stderr, "[GAME] Hand", c.Cost, c)
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
		fmt.Fprintln(os.Stderr, "[GAME] Player", p.Id, "L:", p.Life, "M:", p.current_mana, "/", p.Mana, "D:", p.Deck.Count(), "R:", p.Runes)
	}
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
func (s *State) CreatureFight(id1, id2 int) (error) {
	c1 := s.Hero().BoardGetCard(id1)
	c2 := s.Vilain().BoardGetCard(id2)

	dead1, err1 := s.DealDamage(c1, c2)
	if err1 != nil { return err1 }

	dead2, err2 := s.DealDamage(c2, c1)
	if err2 != nil { return err2 }

	if dead1 { s.Vilain().BoardRemoveCard(id2) }
	if dead2 { s.Hero().BoardRemoveCard(id1) }

	return nil
}
func (s *State) DealDamage(c1, c2 *Card) (bool, error) {
	var c1a, previous_pv, new_pv, id2 int
	var dead bool

	dead = false
	if id2 > 0 {}

	if c1 == nil {
		return false, fmt.Errorf("[GAME][ERROR] Card doesn't exist in Hand")	
	}

	switch c1.Type {
	case CARD_TYPE_CREATURE:
		c1a = c1.Attack
	default:
		if c1.Defense < 0 {
			c1a = -(c1.Defense)
		} else {
			c1a = c1.Defense
		}
	}

	switch c2 {
	case nil:
		previous_pv = s.Vilain().Life
		s.Vilain().ReceiveDamage(c1a)
		new_pv = s.Vilain().Life
		id2 = -1
	default:

		previous_pv = c2.Defense
		if c2.Abilities[CARD_ABILITY_WARD] == "-" && c1a > 0 {
			c2.Defense -= c1a
		}


		if c2.Defense < 0 && c1.Abilities[CARD_ABILITY_BREAKTHROUGH] != "-" {
			s.Hero().ReceiveDamage(-(c2.Defense))
		}

		if c2.Defense <= 0 {
			dead = true
		} else if c1.Abilities[CARD_ABILITY_LETHAL] != "-" && c2.Abilities[CARD_ABILITY_WARD] == "-" {
			dead = true
		}
		new_pv = c2.Defense
		id2 = c2.Id

		if c1a > 0 {
			//fmt.Println("[GAME][DEFENSE] Creature", c2.Id, "Use Ward protection")
			c2.Abilities[CARD_ABILITY_WARD] = "-"
		}
	}
	if previous_pv - new_pv > 0 {
		//fmt.Println("[GAME][DAMAGE] Card", c1.Id, "deal", previous_pv - new_pv, "to", id2)
	}
	return dead, nil
}

func (s *State) CreatureBoostAttack(c *Card, bonus int) (error) {
	if c == nil || c.Type != CARD_TYPE_CREATURE {
		return fmt.Errorf("[GAME][ERROR] Can't boost attack")
	}
	if bonus == 0 {
		return nil
	}
	c.Attack += bonus
	////////fmt.Fprintln(os.Stderr, "[GAME][CREATURE] Boost Attack by", bonus, "for", c.Id)
	return nil
}
func (s *State) CreatureBoostDefense(c *Card, bonus int) (error) {
	if c == nil || c.Type != CARD_TYPE_CREATURE {
		return fmt.Errorf("[GAME][ERROR] Can't boost defense")
	}
	if bonus == 0 {
		return nil
	}
	c.Defense += bonus
	//////fmt.Fprintln(os.Stderr, "[GAME][CREATURE] Boost Defense by", bonus, "for", c.Id)
	return nil
}
func (s *State) CreatureBoostAbilities(c *Card, bonus []string) (error) {
	if c == nil || c.Type != CARD_TYPE_CREATURE {
		return fmt.Errorf("[GAME][ERROR] Can't boost abilities")
	}
	if len(bonus) != len(c.Abilities) {
		return fmt.Errorf("[GAME][ERROR] Wrong number of abilities")
	}
	for i := 0 ; i < len(bonus) ; i++ {
		if bonus[i] != "-" {
			c.Abilities[i] = bonus[i]
			//fmt.Fprintln(os.Stderr, "[GAME][CREATURE] Boost abilities ", bonus[i], "for", c.Id)
		}
	}
	return nil
}
func (s *State) MoveUse(id1, id2 int) (error) {
	c1 := s.Hero().HandGetCard(id1)
	if c1 == nil {
		return fmt.Errorf("GAME: Use %d. Card doesn't exist in Hand", id1)
	}

	err := s.Hero().HandPlayCard(int(id1))
	if err != nil {
		return err
	}

	////////fmt.Fprintln(os.Stderr, "[GAME][USE] Player", s.Hero().Id, "use card", id1, "on", id2, "for cost", c1.Cost, "(", s.Hero().current_mana, ")")

	switch c1.Type {
	case CARD_TYPE_ITEM_BLUE:
		s.Hero().GainLife(c1.HealthChange)
		s.Vilain().ReceiveDamage(-c1.OpponentHealthChange)
		if id2 != -1 {
			s.DealDamage(c1, nil)
		}

	case CARD_TYPE_ITEM_GREEN:
		c2 := s.Hero().BoardGetCard(id2)
		if c2 == nil {
			////////fmt.Fprintln(os.Stderr, "[GAME][ERROR] Player", s.Hero().Id, "doesn't have card", id2, "on board")
			return fmt.Errorf("[GAME][ERROR]: Use %d. Card doesn't exist in Board on Player %d", id2, s.Hero().Id)	
		}
		s.CreatureBoostAttack(c2, c1.Attack)
		s.CreatureBoostDefense(c2, c1.Defense)
		s.CreatureBoostAbilities(c2, c1.Abilities)
		s.Hero().GainLife(c1.HealthChange)
		s.Vilain().ReceiveDamage(-c1.OpponentHealthChange)

	case CARD_TYPE_ITEM_RED:
		c2 := s.Vilain().BoardGetCard(id2)
		if c2 == nil {
			return fmt.Errorf("[GAME][ERROR]: Use %d. Card doesn't exist in Board on Player %d", id2, s.Vilain().Id)	
		}
		s.CreatureBoostAttack(c2, c1.Attack)
		s.CreatureBoostDefense(c2, c1.Defense)
		s.Hero().GainLife(c1.HealthChange)
		s.Vilain().ReceiveDamage(-c1.OpponentHealthChange)

		for i := 0 ; i < len(c1.Abilities) ; i++ {
			if c2.Abilities[i] != "-" {
				c1.Abilities[i] = "-"
			}
		}
		if c2.Defense <= 0 { s.Vilain().BoardRemoveCard(id2) }
	}

	if c1.CardDraw > 0 {
		s.Hero().DrawCardN(c1.CardDraw)
	}


	return nil
}
func (s *State) MoveSummon(id1 int) error {
	c1 := s.Hero().HandGetCard(id1)
	if c1 == nil {
		return fmt.Errorf("GAME: Use %d. Card doesn't exist in Hand", id1)

	}
	//fmt.Println("Before Play card")
	//s.PrintHero()

	err := s.Hero().HandPlayCard(id1)
	if err != nil {
		return err
	}
	if c1.Type != CARD_TYPE_CREATURE {
		return fmt.Errorf("[GAME][SUMMON]: Can't summon card type %d", c1.Type)
	}

	//fmt.Println("[GAME][SUMMON] Player", s.Hero().Id, "summon creature", id1, "for cost", c1.Cost, "(", s.Hero().current_mana, ")")

	s.Hero().GainLife(c1.HealthChange)
	s.Vilain().ReceiveDamage(-c1.OpponentHealthChange)

	
	s.UpdateAvailablesMoves()
	
	return nil
}
func (s *State) MoveAttackPolicy(id1, id2 int) (error) {
	guards := s.Vilain().BoardGetGuardsId()
	//len_guards := len(guards)

	exist, _ := in_array(id2, guards)
	if len(guards) > 0 && ! exist {
		return fmt.Errorf("[GAME][ATTACK] Move ATTACK %d %d not permitted", id1, id2)
	}
	return nil
}
func (s *State) MoveAttack(id1, id2 int) (err error) {
	var err_str string
	c1 := s.Hero().BoardGetCard(id1)
	if c1 == nil {
		////////fmt.Fprintln(os.Stderr, "[GAME][ATTACK] Create", id1, "not present in Hero board")
		s.PrintHeroBoard()
		err_str = fmt.Sprintf("MoveAttack: Current player %d don't have card %d", s.Hero().Id, id1)
		fmt.Fprintln(os.Stderr, "[STATE] Board before Attack")
		s.PrintBoard(s.Hero())
		return errors.New(err_str)
	}

	err = s.MoveAttackPolicy(id1, id2)
	if err != nil {
		return err
	}

	////////fmt.Fprintln(os.Stderr, "[GAME][ATTACK][", s, "] Player", s.Hero().Id, "attack", id2, "with", id1)
	c1a := c1.Attack

	if id2 == -1 {
		s.Vilain().ReceiveDamage(c1a)
	} else {
		c2 := s.Vilain().BoardGetCard(id2)
		if c2 == nil {
			err_str = fmt.Sprintf("MoveAttack: Current oppoent %d don't have card %d", s.Vilain().Id, id2)
			return errors.New(err_str)
		}
		//fmt.Println("[GAME][FIGHT]", c1.Id, "attack", c2.Id, ". May the force be with them")
		err = s.CreatureFight(c1.Id, c2.Id)
	}
	if c1.Abilities[CARD_ABILITY_DRAIN] != "-" {
		s.Hero().GainLife(c1.Attack)
	}

	c1.Attacked = true
	return nil
}
func (s *State) MovePick(id int) (err error) {
	if len(s.Draft) == 0 {
		return fmt.Errorf("wrong action pick")
	}
	if id < len(s.Draft) {
		card := s.Draft[id].Copy()
		card.Id = s.Hero().Deck.Count() + s.Vilain().Deck.Count() + 1
		s.Hero().Draw(s.Draft[id])
	}
	return nil
}


func (s *State) UpdateAvailablesMoves() {
	
	mb := s.AvailablesMovesBoard()	
	mh := s.AvailablesMovesHand()

	s.AMoves = make([]*Move, 0)


	s.AMoves = append(s.AMoves, mh...)
	if len(s.AMoves) == 0 {
		s.AMoves = append(s.AMoves, mb...)
	}
}
func (s *State) CopyAvailablesMoves() []*Move {
	var moves []*Move = nil

	if s.AvailablesMoves == nil {
		return nil
	}

	moves = make([]*Move, len(s.AMoves))
	for i, m := range(s.AMoves) {
		moves[i] = m.Copy()
	}

	return moves
}
func (s *State) AvailablesMovesBoard() []*Move {
	var move *Move
	moves := make([]*Move, 0)

	guards := s.Vilain().BoardGetGuardsId()
	for _, h := range(s.Hero().Board) {
		if h.Attack <= 0 || h.Attacked || h.Charge <= 0 { continue }

		if len(guards) > 0 {
			for _, vid := range(guards) {
				move = NewMove(MOVE_ATTACK, 0, 1, []int{h.Id, vid})
				moves = append(moves, move)
			}
		} else {
			move = NewMove(MOVE_ATTACK, 0, 1, []int{h.Id, -1})
			moves = append(moves, move)

			for _, v := range(s.Vilain().Board) {
				move = NewMove(MOVE_ATTACK, 0, 1, []int{h.Id, v.Id})
				moves = append(moves, move)
			} 
		}
	}
	return moves
}
func (s *State) AvailablesMovesHand() []*Move {
	var move *Move
	moves := make([]*Move, 0)

	
	for _, h := range(s.Hero().Hand) {
		if h.Cost > s.Hero().CurrentMana() { continue }
		switch h.Type {
		case CARD_TYPE_CREATURE:
			move = NewMove(MOVE_SUMMON, h.Cost, 1, []int{h.Id})
			moves = append(moves, move)
		case CARD_TYPE_ITEM_BLUE:
			move = NewMove(MOVE_USE, h.Cost, 1, []int{h.Id, -1})
			moves = append(moves, move)
		case CARD_TYPE_ITEM_GREEN:
			for _, c := range(s.Hero().Board) {
				move = NewMove(MOVE_USE, 0, 1, []int{h.Id, c.Id})
				moves = append(moves, move)
			}
		case CARD_TYPE_ITEM_RED:
			for _, c := range(s.Vilain().Board) {
				move = NewMove(MOVE_USE, 0, 1, []int{h.Id, c.Id})
				moves = append(moves, move)
			}
		}
	}

	return moves
}
func (s *State) AvailablesMoves() []*Move {
	if s.AMoves != nil {
		return s.AMoves
	}

	//fmt.Fprintln(os.Stderr, "[MCTS] Generate Availables Moves")
	s.UpdateAvailablesMoves()

	if len(s.AMoves) == 0 {
		s.AMoves = append(s.AMoves, NewMove(MOVE_PASS, 0, 0, nil))
	}

	return s.AMoves
}
func (s *State) IsEndTurn() bool {
	s.UpdateAvailablesMoves()
	for _, m := range(s.AvailablesMoves()) {
		if m.Cost <= s.Hero().current_mana && m.Type != MOVE_PASS {
			return false
		}
	}
	return true
}
func (s *State) RandomMoveHero(id_hero int) *Move {
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)
	
	l := len(s.AvailablesMoves())
	if l > 0 {
		n := random.Intn(l)
		move := s.AMoves[n]
		//fmt.Fprintln(os.Stderr, "[MCTS] Random move", move.toString(), "in", s.AvailablesMoves())
		s.MoveHero(move)
		return move
	} 
	return nil
}
func (s *State) RandomMove() *Move {
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)

	l := len(s.AvailablesMoves())
	if l > 0 {
		n := random.Intn(l)
		move := s.AvailablesMoves()[n]
		//fmt.Println("[MCTS] Random move", move.toString(), "in", s.AvailablesMoves())
		s.Move(move)
		return move
	}
	return nil
}
func (s *State) MoveHero(m *Move) bool {
	id_hero := s.Hero().Id
	end_turn := s.Move(m)
	if end_turn {
		s.NextTurnHero(id_hero)
	}
	return end_turn
}
func (s *State) Move(m *Move) bool {
	var err error
	var end_turn bool = false

	if m == nil {
		return true
	}
	//fmt.Println("[STATE] Move", m.toString())
	switch m.Type {
	case MOVE_PASS:
		if len(s.Draft) > 0 {
			s.MovePick(0)
		}
		s.AMoves = nil
		end_turn = true
	case MOVE_PICK:
		err = s.MovePick(m.Params[0])
	case MOVE_SUMMON:
		err = s.MoveSummon(m.Params[0])
	case MOVE_ATTACK:
		err = s.MoveAttack(m.Params[0], m.Params[1])
	case MOVE_USE:
		err = s.MoveUse(m.Params[0], m.Params[1])
	}

	if err != nil {
		//fmt.Fprintln(os.Stderr, "[MCTS][ERROR]", err)
	}
	
	if m.Type != MOVE_PASS {
		//fmt.Fprintln(os.Stderr, "[MCTS] Update Availables Moves")
		s.UpdateAvailablesMoves()
	}
	//s.DeleteAMoves(m)
	////fmt.Fprintln(os.Stderr, "[MCTS] Remove move", m.toString(), "from available moves")

	if s.IsEndTurn() {
		//fmt.Println("[MOVE] End turn reach")	
		end_turn = true
	}
	return end_turn
}
func (s *State) GameOver() *Player {
	if s.Hero().Life <= 0 {
		return s.Vilain()
	} else if s.Vilain().Life <= 0 {
		return s.Hero()
	}
	return nil
}

func InitGame(p1, p2 *Player) {
	p1.Deck.FillRandom(STARTING_CARD)
	p2.Deck.FillRandom(STARTING_CARD)

	for i, c := range(p1.Deck.Cards) { c.Id = i + 1	}
	for i, c := range(p2.Deck.Cards) { c.Id = p1.Deck.Count() + i + 1 } 

	p1.DrawCardN(4)
	p2.DrawCardN(5)

}

func CountNode(n *Node) int {
	count := 1 
	for _, c := range(n.Children) {
		count += CountNode(c)
	}
	return count
}

/*
func main() {
	p1 := NewPlayer(1, STARTING_LIFE, STARTING_MANA, STARTING_RUNES)
	p2 := NewPlayer(2, STARTING_LIFE, STARTING_MANA, STARTING_RUNES)


	InitGame(p1, p2)
	init_state := InitState(p2, p1)
	init_state.NextTurn()

	for i := 0 ; i  < 11 ; i++ {
		//fmt.Fprintln(os.Stderr, "=========== Random Move", i, "================")
		//init_state.Print()
		for _, _ = range(init_state.AvailablesMoves()) {
			//fmt.Fprintln(os.Stderr, "[MCTS][RANDOMMOVE] AMoves", m.toString(), m)
		}
		rmove := init_state.RandomMoveHero(p1.Id)
		if rmove != nil {
			//fmt.Fprintln(os.Stderr, "[MCTS][RANDOMMOVE] Pick", rmove.toString(), rmove)
		}
		if init_state.IsEndTurn() {
			//fmt.Fprintln(os.Stderr, "[MCTS][RANDOMMOVE] End turn")
			init_state.NextTurnHero(p1.Id)
		} else {
			//fmt.Fprintln(os.Stderr, "[MCTS][RANDOMMOVE] There is still available moves", len(init_state.AMoves))
		}
	}

	//fmt.Fprintln(os.Stderr, "==================== END RANDOM =====================")
	//init_state.NextTurn()
	//init_state.Print()
	//fmt.Fprintln(os.Stderr, "==================== END STATE =====================")
	var state *State = init_state
	for i := 0 ; i < 1 ; i++ {
		if state.GameOver() != nil { break }
		fmt.Println("============== BEGIN ===============")
		root_node := NewNode(nil, state, nil)
		state.Print()
		//fmt.Println(root_node)
		//fmt.Println(root_node.State)
		start := time.Now()
		moves := MonteCarloMoves(root_node, -1)
		//root_node.ExportGraph()
		elapsed := time.Since(start)
		fmt.Println("====================================")
		root_node.Print()
		fmt.Println("[MCTS] Nodes", CountNode(root_node))
		fmt.Println("[MCTS] Duration:", elapsed)
		fmt.Println("[MCTS] Moves:", len(moves))
		state = state.Copy()
		//time.Sleep(10 * time.Second)
		for _, m := range(moves) {
			fmt.Println("[MCTS] Suggest move", m.toString())
			state.MoveHero(m)
			if state.IsEndTurn() {
				state.NextTurn()
				//qstate.NextTurnHero(state.Hero().Id)
			} 
		}
		
	}
	if state != nil && state.GameOver() != nil {
		fmt.Println("The winner is", state.GameOver().Id)
	}
}


*/
func main() {

    var hero, vilain, previous_hero, previous_vilain  *Player
	var step_draft bool
	var previous_board []int = make([]int, 0)

    for {

        var players []*Player
 
        players = make([]*Player, 2)
        players[0] = NewPlayer(0, 0, 0, 0)
        players[1] = NewPlayer(1, 0, 0, 0)

        hero    = players[0]
        vilain  = players[1]
        step_draft = false

        for i := 0; i < 2; i++ {
            var playerHealth, playerMana, playerDeck, playerRune int
            fmt.Scan(&playerHealth, &playerMana, &playerDeck, &playerRune)

            players[i].Life = playerHealth
            players[i].Mana = playerMana
			players[i].Runes = playerRune
			players[i].current_mana = playerMana
        }

        if previous_hero != nil && previous_vilain != nil {
            hero.Deck = previous_hero.Deck
            vilain.Deck = previous_vilain.Deck
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
            //fmt.Fprintln(os.Stderr, cardNumber, instanceId, location, cardType, cost, attack, defense, abilities, myHealthChange, opponentHealthChange, cardDraw)
            
            card = NewCard(cardNumber, instanceId, cost, attack, defense, abilities, myHealthChange, opponentHealthChange, cardDraw, cardType)
            switch cardType {
            case CARD_TYPE_CREATURE:
                if card.Abilities[CARD_ABILITY_CHARGE] != "-" {
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
                        if previous_vilain != nil && previous_vilain.BoardGetCard(card.Id) != nil {
                            card.Charge = 1
                        }
                        vilain.Board = append(vilain.Board, card)
                case 0:
                        hero.Draw(card)
                case 1:
                        if len(previous_board) > 0 {
							exist, _ := in_array(card.Id, previous_board)
							if exist {
								card.Charge = 1
							} 
						}
						hero.Board = append(hero.Board, card)
                }
            }
        }

        if step_draft {
			fmt.Println("PASS")
                //hero.Pick(draft)
        } else {
			init_state := InitState(hero, vilain)
			root_node := NewNode(nil, init_state, nil)
			
			hero := root_node.State.Hero()
			fmt.Fprintln(os.Stderr, "[MCTS] Hero", hero.Life, hero.CurrentMana(), hero.Runes, hero.Deck.Count())
			for _, c := range(root_node.State.Hero().Hand) {
				fmt.Fprintln(os.Stderr, "[MCTS] Hand", c.Cost, c)
			}
			for _, c := range(root_node.State.Hero().Board) {
				fmt.Fprintln(os.Stderr, "[MCTS] Board", c.Attack, c.Defense, c)
			}
			moves := MonteCarloMoves(root_node, -1)
			fmt.Fprintln(os.Stderr, "[MCTS] suggest move", moves)
			for _, m := range(moves) {
			    //fmt.Println("[MCTS] Suggest move", m.toString())
			    root_node.State.Move(m)
		    }
			str_moves := make([]string, 0)
			for _, m := range(moves) { str_moves = append(str_moves, m.toString()) }
			fmt.Println(strings.Join(str_moves, ";"))

			previous_board = make([]int, len(root_node.State.Hero().Board))
			for i, c := range(root_node.State.Hero().Board) {
				previous_board[i] = c.Id
			}
        }
        previous_hero   = hero.Copy()
        previous_vilain = vilain.Copy()
        
    }
}

