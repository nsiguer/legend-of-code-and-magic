package main

import (
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
	"time"

	//copier "github.com/jinzhu/copier"
)

const (
	MC_MAX_ITERATION = 10
	MAX_INT          = int(^uint(0) >> 1)
	MIN_INT          = -MAX_INT - 1

	BIAS_PARAMETER = 0.9
)

type Node struct {
	Parent   *Node
	Children []*Node
	State    *State
	Wins     float64
	Visits   int
	Actions  []string
}

type State struct {
	P1_life int
	P2_life int

	P1_id int
	P2_id int

	P1_mana int
	P2_mana int

	P1_board []Card
	P2_board []Card

	P1_hand []Card
	P2_hand []Card
}

func NewState() *State {
	return &State{}
}

func RootState() (state *State) {
	state = NewState()
	state.Reset()
	return
}

func (s *State) Reset() {
	s.P1_life = STARTING_LIFE
	s.P2_life = STARTING_LIFE

	s.P1_mana = STARTING_MANA
	s.P2_mana = STARTING_MANA

	s.P1_board = make([]Card, 0)
	s.P2_board = make([]Card, 0)

	s.P1_hand = make([]Card, 0)
	s.P2_hand = make([]Card, 0)

	for i := 0; i < STARTING_CARD; i++ {
		s.P1_hand = append(s.P1_hand, PickRandomCard())
		s.P2_hand = append(s.P2_hand, PickRandomCard())
	}
	s.P2_hand = append(s.P2_hand, PickRandomCard())
}

func (s *State) GameOver() bool {
	return true
}

func GetCardFromId(cards []Card, id int) (Card, int) {
	for i, c := range cards {
		if int(c.Id()) == id {
			return c, i
		}
	}
	return nil, -1
}

func (s *State) Action(actions_str string) *Node {
	new_state := NewState()
	copy(new_state, s)

	new_node := NewNode(nil, new_state)
	total_mana := s.P1_mana

	actions := strings.Split(actions_str, ";")
	for _, a := range actions {
		args := strings.Split(a, " ")
		switch args[0] {
		case "SUMMON":
			if len(args) != 2 {
				fmt.Println("Wrong format for SUMMON", args)
				continue
			}
			id, _ := strconv.Atoi(args[1])
			c, idx := GetCardFromId(s.P1_hand, id)

			if int(c.Cost()) > total_mana {
				fmt.Println("Can't summon", c, ". No more mana")
				continue
			}
			s.P1_hand[idx] = s.P1_hand[len(s.P1_hand)-1]
			s.P1_hand = s.P1_hand[:len(s.P1_hand)-2]
			s.P1_board = append(s.P1_board, c)
		case "ATTACK":
			if len(args) != 3 {
				fmt.Println("Wrong format for ATTACK", args)
				continue
			}
			id1, _ := strconv.Atoi(args[1])
			id2, _ := strconv.Atoi(args[2])

			c1, idx1 := GetCardFromId(s.P1_board, id1)
			type1 := reflect.TypeOf(c1)
			if type1.String() != "*main.Monster" {
				fmt.Println("Wront type for Attack", type1)
				continue
			}
			cm1 := c1.(*Monster)

			if id2 == -1 {
				s.P2_life = s.P2_life - int(cm1.Attack)
			} else {
				c2, idx2 := GetCardFromId(s.P2_board, id2)
				type2 := reflect.TypeOf(c2)
				if type2.String() != "*main.Monster" {
					fmt.Println("Wront type for Attack", type2)
					continue
				}

				cm2 := c2.(*Monster)
				cm2.Defense -= int(cm1.Attack)
				cm1.Defense -= int(cm2.Attack)

				if cm1.Defense <= 0 {
					s.P1_board[idx1] = s.P1_board[len(s.P1_board)-1]
					s.P1_board = s.P1_board[:len(s.P1_board)-2]
				}
				if cm2.Defense <= 0 {
					s.P2_board[idx2] = s.P2_board[len(s.P2_board)-2]
					s.P2_board = s.P2_board[:len(s.P2_board)-2]
				}
			}

		default:
			fmt.Println("Unknow action", args[0])
		}

	}
	new_node.State.P1_mana += 1

	return new_node
}

func combinatonSummon(max_mana int, cards []Card) []string {
	combinatons := make([]string, 0)
	for i, _ := range cards {
		current_mana := max_mana
		/*
			if int(c1.Cost()) <= max_mana {
				combinatons = append(combinatons, fmt.Sprintf("SUMMON %d", c1.Id()))
			}
			current_mana -= int(c1.Cost())
		*/
		tmp_combinaton := make([]string, 0)
		for j := i; j < len(cards); j++ {
			c2 := cards[j]
			if int(c2.Cost()) <= current_mana {
				tmp_combinaton = append(tmp_combinaton, fmt.Sprintf("SUMMON %d", c2.Id()))
				current_mana -= int(c2.Cost())
			}

		}
		if len(tmp_combinaton) > 0 {
			combinatons = append(combinatons, strings.Join(tmp_combinaton, ";"))
		}
	}
	return combinatons
}

func combinatonAttack(b1, b2 []Card) []string {
	combinatons := make([]string, 0)
	/*
		for i, c1 := range b1 {
			tmp_combinatons := make([]string, 0)
			tmp_combinatons = append(combinatons, fmt.Sprintf("ATTACK %d -1", c1.Id()))
			for j, c2 := range b2 {
				tmp_combinatons = append(tmp_combinatons, fmt.Sprintf("ATTACK %d %d", c1.Id(), c2.Id()))
			}

			combinatons = append(combinatons, strings.Join(tmp_combinaton, ";"))
		}
	*/
	return combinatons
}

func GenerateActions(s *State) []string {
	return append(combinatonSummon(s.P1_mana, s.P1_hand), combinatonAttack(s.P1_board, s.P2_board)...)
}

func NewNode(parent *Node, state *State) *Node {
	return &Node{
		Parent:   parent,
		Children: nil,
		State:    state,
		Wins:     0,
		Visits:   0,
		Actions:  nil,
	}
}

func (n *Node) AddChild(node *Node) {
	if n.Children == nil {
		n.Children = make([]*Node, 0)
	}

	n.Children = append(n.Children, node)
}

func (n *Node) UpdateScore(score float64) {
	n.Visits++
	n.Wins = n.Wins + score
}

func MonteCarlo(node *Node) *Node {
	for i := 0; i < MC_MAX_ITERATION; i++ {
		node = MCSelection(node)
		node = MCExpansion(node)
		score := MCSimulation(node)
		node = MCBackPropagation(node, score)

	}
	return node
}

func MCSelection(node *Node) *Node {
	var candidate_node *Node
	if node.Actions == nil && node.Children != nil && len(node.Children) > 0 {
		candidate_node = nil
		score := -100.0

		for _, n := range node.Children {
			child_score := MCCalculateScore(node)
			if child_score > score {
				score = child_score
				candidate_node = n
			}
		}

		return MCSelection(candidate_node)

	}
	return node
}

func MCCalculateScore(node *Node) float64 {
	exploitScore := float64(node.Wins) / float64(node.Visits)
	exploreScore := math.Sqrt(math.Log(float64(node.Parent.Visits)) / float64(node.Visits))
	exploreScore = BIAS_PARAMETER * exploreScore

	return exploitScore + exploreScore
}

func MCExpansion(node *Node) *Node {
	if len(node.Actions) == 0 {
		return node
	}

	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)
	idx := random.Intn(len(node.Actions))

	new_node := node.State.Action(node.Actions[idx])
	new_node.Parent = node
	
	node.Actions[idx] = node.Actions[len(node.Actions)-1]
	node.Actions = node.Actions[:len(node.Actions)-2]

	node.Children = append(node.Children, new_node)
	if new_node.State.GameOver() {
		new_node.Actions = GenerateActions(new_node.State)
	}

	return node
}

func MCSimulation(node *Node) float64 {
	simulate_node := NewNode(nil, nil)
	copy(simulate_node, node)

	game_over := false

	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)

	for !game_over {
		if simulate_node.State.GameOver() {
			if simulate_node.State.P1_life <= 0 {
				return 1
			} else {
				return 0
			}
		}

		simulate_node.Actions = GenerateActions(simulate_node.State)
		idx := random.Intn(len(simulate_node.Actions))
		simulate_node = simulate_node.State.Action(node.Actions[idx])
	}

	return 0
}

func MCBackPropagation(node *Node, score float64) *Node {
	current_player_id := node.State.P1_id

	for node.Parent != nil {
		if node.Parent.State.P1_id == current_player_id {
			node.UpdateScore(score)
		} else {
			if score == 0 {
				node.UpdateScore(1)
			}
			if score == 1 {
				node.UpdateScore(0)
			}
		}
		node = node.Parent
	}

	node.Visits++

	return node
}
