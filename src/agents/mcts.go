package agent

import (
	g "game"

	"fmt"
	"time"
	"math/rand"
	"math"
	"os"

)


const (
	BIAS_PARAMETER 	= 0.7
	MCTS_ITERATION	= 150
	MCTS_SIMULATION	= 100
	MCTS_TIMEOUT	= 65 
)

/* STARTING */
type AgentMCTS struct {

}

func NewAgentMCTS() *AgentMCTS {
	return &AgentMCTS{}
}

func (a *AgentMCTS) Name() string { return "MCTS" }
func (a *AgentMCTS) Think(s *g.State) []*g.Move {
	node, _ := MonteCarloTimeout(s, BIAS_PARAMETER, MCTS_ITERATION, MCTS_SIMULATION, MCTS_TIMEOUT)
	moves := make([]*g.Move, 0)



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
		moves = append(moves, g.NewMove(g.MOVE_PASS, 0, 0, nil))
	}
	return moves
}



type Node struct {
	Id		 	int
	Parent   		*Node
	Children 		[]*Node
	State    		*g.State
	Outcome     	float64
	Visits   		int
	ByMove			*g.Move
	EndTurn  		bool
	UnexploreMoves	[]*g.Move
}

func NewNode(parent *Node, state *g.State, move *g.Move) *Node {
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
func (n *Node) DeleteUnexploreMoves(m *g.Move) (error) {
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
		n.UnexploreMoves = make([]*g.Move, 0)
	}
	return nil
}

func (n *Node) UpdateScore(score float64) {
	n.Visits++
	n.Outcome += score
}

func MonteCarloTimeout(state *g.State, bias float64, iteration, simulation, timeout int) (*Node, error) {
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
func MonteCarlo(state *g.State, bias float64, iteration, simulation int) (*Node, error) {
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
	if new_state.IsEndTurn() || rmove.Type == g.MOVE_PASS {
		new_node.EndTurn = true
		new_state.NextTurn()
	}

	return new_node
}
func MCSimulation(state *g.State, simulation int) float64 {

	var moves []*g.Move = nil


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
